package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError + 10}))
}

func TestRequestID_AddsHeader(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	id := rr.Header().Get("X-Request-ID")
	if id == "" {
		t.Error("expected X-Request-ID header to be set")
	}
	if len(id) < 36 { // UUID v4 is 36 chars
		t.Errorf("expected UUID format, got %q", id)
	}
}

func TestRequestID_UniquePerRequest(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rr1 := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()

	handler.ServeHTTP(rr1, httptest.NewRequest("GET", "/", nil))
	handler.ServeHTTP(rr2, httptest.NewRequest("GET", "/", nil))

	id1 := rr1.Header().Get("X-Request-ID")
	id2 := rr2.Header().Get("X-Request-ID")

	if id1 == id2 {
		t.Error("expected different request IDs for different requests")
	}
}

func TestRequestID_SetsContext(t *testing.T) {
	var ctxID interface{}

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxID = r.Context().Value(RequestIDKey)
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	if ctxID == nil {
		t.Error("expected request ID in context")
	}
}

func TestRecoverer_CatchesPanic(t *testing.T) {
	handler := Recoverer(discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 after panic, got %d", rr.Code)
	}
}

func TestRecoverer_PassesThroughNormally(t *testing.T) {
	handler := Recoverer(discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	}))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRateLimiter_AllowsBurst(t *testing.T) {
	handler := RateLimiter(100, 10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Should allow 10 requests immediately (burst)
	for i := 0; i < 10; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))
		if rr.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, rr.Code)
		}
	}
}

func TestRateLimiter_BlocksExcessiveRequests(t *testing.T) {
	handler := RateLimiter(1, 2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Exhaust burst
	for i := 0; i < 2; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))
	}

	// Next request should be rate limited
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 after burst exhausted, got %d", rr.Code)
	}
}

func TestFirebaseAuth_MissingHeader(t *testing.T) {
	handler := FirebaseAuth(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with no auth header, got %d", rr.Code)
	}
}

func TestFirebaseAuth_InvalidPrefix(t *testing.T) {
	handler := FirebaseAuth(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Basic abc123")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with non-Bearer prefix, got %d", rr.Code)
	}
}

func TestFirebaseAuth_DevMode_TreatsTokenAsUID(t *testing.T) {
	var capturedUID string

	handler := FirebaseAuth(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUID = GetFirebaseUID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer test-user-123")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 in dev mode, got %d", rr.Code)
	}
	if capturedUID != "test-user-123" {
		t.Errorf("expected UID 'test-user-123', got %q", capturedUID)
	}
}

type mockVerifier struct {
	uid string
	err error
}

func (m *mockVerifier) VerifyIDToken(_ context.Context, _ string) (string, error) {
	return m.uid, m.err
}

func TestFirebaseAuth_WithVerifier_Success(t *testing.T) {
	var capturedUID string

	verifier := &mockVerifier{uid: "verified-uid-456", err: nil}
	handler := FirebaseAuth(verifier)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUID = GetFirebaseUID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer some-firebase-token")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if capturedUID != "verified-uid-456" {
		t.Errorf("expected UID 'verified-uid-456', got %q", capturedUID)
	}
}

func TestFirebaseAuth_WithVerifier_InvalidToken(t *testing.T) {
	verifier := &mockVerifier{uid: "", err: fmt.Errorf("invalid token")}
	handler := FirebaseAuth(verifier)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for invalid token, got %d", rr.Code)
	}
}

func TestGetFirebaseUID_EmptyContext(t *testing.T) {
	uid := GetFirebaseUID(context.Background())
	if uid != "" {
		t.Errorf("expected empty string from empty context, got %q", uid)
	}
}

func TestResponseWriter_CapturesStatusCode(t *testing.T) {
	rr := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rr, statusCode: http.StatusOK}

	rw.WriteHeader(http.StatusNotFound)

	if rw.statusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rw.statusCode)
	}
}
