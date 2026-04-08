package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/luke/mockstarket/internal/handler"
	mw "github.com/luke/mockstarket/internal/middleware"
	ws "github.com/luke/mockstarket/internal/websocket"
)

// New creates the HTTP router with all routes and middleware.
func New(h *handler.Handler, hub *ws.Hub, authVerifier mw.FirebaseAuthVerifier, corsOrigins []string, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(mw.RequestID)
	r.Use(mw.Logger(logger))
	r.Use(mw.Recoverer(logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(mw.RateLimiter(100, 200))

	// Health check (no auth)
	r.Get("/api/v1/system/health", h.HealthCheck)

	// WebSocket endpoint
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for WebSocket
		},
	}

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("WebSocket upgrade failed", "error", err)
			return
		}

		// Extract user ID from query param or first message
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			userID = "anonymous"
		}

		client := ws.NewClient(conn, hub, userID, logger)
		if !hub.Register(client) {
			_ = conn.Close()
			return
		}

		client.Run()
	})

	// Public API routes (no auth required)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/guest", h.CreateGuest)
	})

	// Authenticated API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.FirebaseAuth(authVerifier))

		// Auth
		r.Get("/auth/me", h.GetMe)
		r.Put("/auth/me", h.UpdateMe)
		r.Delete("/auth/me", h.DeleteMe)

		// Market Status
		r.Get("/market/status", h.GetMarketStatus)

		// Stocks
		r.Get("/stocks", h.ListStocks)
		r.Get("/stocks/market-summary", h.GetMarketSummary)
		r.Get("/stocks/{ticker}", h.GetStock)
		r.Get("/stocks/{ticker}/history", h.GetStockHistory)
		r.Get("/stocks/{ticker}/holdings", h.GetETFHoldings)

		// Trading
		r.Post("/trades", h.ExecuteTrade)
		r.Get("/trades", h.GetTradeHistory)

		// Orders
		r.Post("/orders", h.CreateOrder)
		r.Get("/orders", h.ListOrders)
		r.Delete("/orders/{id}", h.CancelOrder)

		// Portfolio
		r.Get("/portfolio", h.GetPortfolio)
		r.Get("/portfolio/history", h.GetPortfolioHistory)

		// Leaderboard
		r.Get("/leaderboard", h.GetLeaderboard)

		// Alerts
		r.Post("/alerts", h.CreateAlert)
		r.Get("/alerts", h.ListAlerts)
		r.Delete("/alerts/{id}", h.DeleteAlert)

		// Achievements
		r.Get("/achievements", h.ListAchievements)
		r.Get("/achievements/me", h.GetMyAchievements)

		// Challenges
		r.Get("/challenges/today", h.GetTodaysChallenge)
		r.Post("/challenges/check", h.CheckChallenge)
		r.Post("/challenges/{id}/claim", h.ClaimChallengeReward)

		// Watchlist
		r.Get("/watchlist", h.GetWatchlist)
		r.Post("/watchlist", h.AddToWatchlist)
		r.Delete("/watchlist/{ticker}", h.RemoveFromWatchlist)

		// Options
		r.Get("/stocks/{ticker}/options", h.GetOptionChain)
		r.Get("/stocks/{ticker}/options/expirations", h.GetOptionExpirations)
		r.Get("/options/{id}", h.GetOptionContractDetail)
		r.Post("/options/trades", h.ExecuteOptionsTrade)
		r.Get("/options/trades", h.GetOptionsTradeHistory)
		r.Get("/options/positions", h.GetOptionsPositions)
		r.Post("/options/orders", h.CreateOptionsOrder)
		r.Get("/options/orders", h.ListOptionsOrders)
		r.Delete("/options/orders/{id}", h.CancelOptionsOrder)
	})

	return r
}
