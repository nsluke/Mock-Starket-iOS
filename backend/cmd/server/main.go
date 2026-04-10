package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luke/mockstarket/internal/config"
	"github.com/luke/mockstarket/internal/handler"
	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/middleware"
	"github.com/luke/mockstarket/internal/polygon"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/server"
	"github.com/luke/mockstarket/internal/service"
	"github.com/luke/mockstarket/internal/simulation"
	ws "github.com/luke/mockstarket/internal/websocket"
	"github.com/luke/mockstarket/internal/worker"
)

func main() {
	// Logger
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Config
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Database
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("connected to database")

	// Repository
	repo := repository.New(pool)

	// Load stocks from DB
	stocks, err := repo.GetAllStocks(context.Background())
	if err != nil {
		logger.Error("failed to load stocks", "error", err)
		os.Exit(1)
	}

	// Initialize market data provider (simulation or polygon)
	var provider market.PriceProvider

	switch cfg.MarketDataSource {
	case "polygon":
		if cfg.PolygonAPIKey == "" {
			logger.Error("POLYGON_API_KEY is required when MARKET_DATA_SOURCE=polygon")
			os.Exit(1)
		}
		polygonClient := polygon.NewClient(cfg.PolygonAPIKey, cfg.PolygonBaseURL, 5, logger)
		feed := polygon.NewMarketFeed(polygonClient, polygon.FeedConfig{
			WSEnabled:    cfg.PolygonWSEnabled,
			PollInterval: time.Duration(cfg.PolygonPollInterval) * time.Millisecond,
		}, logger)
		// Track all stocks from DB — the seed script now only creates real tickers
		tracked := 0
		for _, s := range stocks {
			feed.TrackTicker(s.Ticker, s.Sector, s.AssetType)
			tracked++
		}
		provider = feed
		logger.Info("using polygon.io market data", "tickers", tracked, "ws_enabled", cfg.PolygonWSEnabled)

	default: // "simulation"
		engine := simulation.NewEngine(cfg.SimTickMS, cfg.MarketEventFreq, cfg.SimTicksPerDay, logger)
		for _, s := range stocks {
			engine.AddStock(
				s.Ticker, s.Name, s.Sector,
				s.CurrentPrice.InexactFloat64(),
				s.Volatility.InexactFloat64(),
				s.Drift.InexactFloat64(),
				s.MeanReversion.InexactFloat64(),
			)
		}
		provider = engine
		logger.Info("using simulation engine", "count", len(stocks))
	}

	// WebSocket Hub
	hub := ws.NewHub(cfg.MaxWSClients, logger)

	// ---- Observers (simulation -> workers) ----

	// 1. WebSocket bridge (broadcasts prices to connected clients)
	bridge := worker.NewSimulationBridge(hub, logger)
	provider.AddObserver(bridge)

	// 2. Price history worker (persists OHLCV to database)
	priceHistoryWorker := worker.NewPriceHistoryWorker(repo, logger)
	provider.AddObserver(priceHistoryWorker)

	// 3. Price alert worker (evaluates alert conditions)
	priceAlertWorker := worker.NewPriceAlertWorker(repo, hub, logger)
	provider.AddObserver(priceAlertWorker)

	// Services
	tradeSvc := service.NewTradeService(repo, provider)

	// 4. Achievement service
	achievementSvc := service.NewAchievementService(repo, provider, hub, cfg.StartingCash, logger)
	tradeSvc.SetOnTradeExecuted(achievementSvc.OnTradeExecuted)

	// 5. Order matching worker (evaluates limit/stop orders)
	orderMatchingWorker := worker.NewOrderMatchingWorker(repo, tradeSvc, hub, logger)
	provider.AddObserver(orderMatchingWorker)

	// 6. Leaderboard worker (computes rankings periodically)
	leaderboardWorker := worker.NewLeaderboardWorker(repo, provider, cfg.StartingCash, logger)

	// 7. Stock sync worker (persists live prices back to stocks table)
	stockSyncWorker := worker.NewStockSyncWorker(repo, 5*time.Second, logger)
	provider.AddObserver(stockSyncWorker)

	// 8. Challenge service and worker
	challengeSvc := service.NewChallengeService(repo, provider, logger)
	challengeWorker := worker.NewChallengeWorker(challengeSvc, logger)

	// 9. Options trade service
	optionsTradeSvc := service.NewOptionsTradeService(repo, provider)

	// 10. Options pricing worker (recalculates every 5th tick)
	optionsPricingWorker := worker.NewOptionsPricingWorker(repo, provider, hub, logger)
	provider.AddObserver(optionsPricingWorker)

	// 11. Options chain generator (creates contracts on startup + every sim-day)
	optionsChainWorker := worker.NewOptionsChainWorker(repo, provider, logger)

	// 12. Options expiration worker (settles expired contracts)
	optionsExpirationWorker := worker.NewOptionsExpirationWorker(repo, provider, logger)

	// Firebase Auth
	var authVerifier middleware.FirebaseAuthVerifier
	if !cfg.DevMode {
		verifier, err := middleware.NewFirebaseVerifier(
			context.Background(),
			cfg.FirebaseProjectID,
			cfg.FirebaseCredentialsFile,
		)
		if err != nil {
			logger.Error("failed to initialize Firebase auth", "error", err)
			os.Exit(1)
		}
		authVerifier = verifier
		logger.Info("Firebase auth enabled", "project_id", cfg.FirebaseProjectID)
	} else {
		logger.Warn("running in dev mode — Firebase auth disabled, tokens treated as UIDs")
	}

	// Handlers
	h := handler.New(repo, tradeSvc, provider, hub, cfg.StartingCash)
	h.SetChallengeService(challengeSvc)
	h.SetOptionsTradeService(optionsTradeSvc)

	// Router
	router := server.New(h, hub, authVerifier, cfg.CORSOrigins, logger)

	// HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start market data provider
	go func() {
		if err := provider.Run(ctx); err != nil && err != context.Canceled {
			logger.Error("market data provider error", "error", err)
		}
	}()

	// Start price history worker
	go priceHistoryWorker.Run(ctx)

	// Start leaderboard worker
	go leaderboardWorker.Run(ctx)

	// Start stock sync worker
	go stockSyncWorker.Run(ctx)

	// Start challenge worker
	go challengeWorker.Run(ctx)

	// Start options workers
	go optionsChainWorker.Run(ctx)
	go optionsExpirationWorker.Run(ctx)

	// Start stale WebSocket client cleaner
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				hub.CleanStale(90 * time.Second)
			}
		}
	}()

	// Start HTTP server
	go func() {
		logger.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("shutting down", "signal", sig)

	// Graceful shutdown
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}

	logger.Info("server stopped")
}
