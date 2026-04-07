package middleware

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type contextKey string

const (
	UserIDKey      contextKey = "user_id"
	FirebaseUIDKey contextKey = "firebase_uid"
	RequestIDKey   contextKey = "request_id"
)

// RequestID adds a unique request ID to each request.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger logs request details using structured logging.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(ww, r)

			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", r.Context().Value(RequestIDKey),
			)
		})
	}
}

// Recoverer catches panics and returns 500.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					logger.Error("panic recovered", "error", rvr, "path", r.URL.Path)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiter applies per-IP rate limiting.
func RateLimiter(rps float64, burst int) func(http.Handler) http.Handler {
	limiters := make(map[string]*rate.Limiter)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
				ip = strings.Split(fwd, ",")[0]
			}

			if _, ok := limiters[ip]; !ok {
				limiters[ip] = rate.NewLimiter(rate.Limit(rps), burst)
			}

			if !limiters[ip].Allow() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// FirebaseAuth verifies Firebase ID tokens.
// In development mode (verifier == nil), it extracts the token as a plain user ID.
type FirebaseAuthVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (uid string, err error)
}

func FirebaseAuth(verifier FirebaseAuthVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			if verifier == nil {
				// Dev mode: treat token as firebase UID directly
				ctx := context.WithValue(r.Context(), FirebaseUIDKey, token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			uid, err := verifier.VerifyIDToken(r.Context(), token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), FirebaseUIDKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetFirebaseUID extracts the Firebase UID from context.
func GetFirebaseUID(ctx context.Context) string {
	if uid, ok := ctx.Value(FirebaseUIDKey).(string); ok {
		return uid
	}
	return ""
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Hijack implements http.Hijacker so WebSocket upgrades work through the logger middleware.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not implement http.Hijacker")
}
