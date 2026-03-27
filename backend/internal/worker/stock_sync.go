package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/simulation"
	"github.com/shopspring/decimal"
)

// stockSnapshot holds the latest price data for a single stock.
type stockSnapshot struct {
	Price  decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Volume int64
}

// StockSyncWorker periodically writes live simulation prices back to the
// stocks table so that REST API responses (which read from DB) stay current.
type StockSyncWorker struct {
	repo   *repository.Repo
	logger *slog.Logger

	mu       sync.Mutex
	latest   map[string]stockSnapshot
	interval time.Duration
}

// NewStockSyncWorker creates a worker that syncs simulation prices to the DB.
// flushInterval controls how often prices are persisted (e.g. 5s).
func NewStockSyncWorker(repo *repository.Repo, flushInterval time.Duration, logger *slog.Logger) *StockSyncWorker {
	return &StockSyncWorker{
		repo:     repo,
		logger:   logger,
		latest:   make(map[string]stockSnapshot),
		interval: flushInterval,
	}
}

// OnPriceBatch captures the latest prices from each simulation tick.
func (w *StockSyncWorker) OnPriceBatch(updates []simulation.PriceUpdate) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, u := range updates {
		w.latest[u.Ticker] = stockSnapshot{
			Price:  u.Price,
			High:   u.High,
			Low:    u.Low,
			Volume: u.Volume,
		}
	}
}

// OnMarketEvent is a no-op for this worker.
func (w *StockSyncWorker) OnMarketEvent(_ simulation.MarketEvent) {}

// Run starts the periodic flush loop. Blocks until context is cancelled.
func (w *StockSyncWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info("stock sync worker started", "interval", w.interval)

	for {
		select {
		case <-ctx.Done():
			// Final flush on shutdown
			w.flush(context.Background())
			w.logger.Info("stock sync worker stopped")
			return
		case <-ticker.C:
			w.flush(ctx)
		}
	}
}

// flush writes all buffered prices to the stocks table.
func (w *StockSyncWorker) flush(ctx context.Context) {
	w.mu.Lock()
	snapshot := make(map[string]stockSnapshot, len(w.latest))
	for k, v := range w.latest {
		snapshot[k] = v
	}
	w.mu.Unlock()

	if len(snapshot) == 0 {
		return
	}

	synced := 0
	for ticker, snap := range snapshot {
		err := w.repo.UpdateStockPrices(ctx, ticker, snap.Price, snap.High, snap.Low, snap.Volume)
		if err != nil {
			w.logger.Error("stock sync: failed to update",
				"ticker", ticker, "error", err)
			continue
		}
		synced++
	}

	w.logger.Debug("stock prices synced to database", "count", synced)
}
