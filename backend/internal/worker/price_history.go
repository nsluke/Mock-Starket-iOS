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

// ohlcv tracks open-high-low-close-volume for a single stock in a single interval window.
type ohlcv struct {
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume int64
	Set    bool
}

func (o *ohlcv) update(price decimal.Decimal, volume int64) {
	if !o.Set {
		o.Open = price
		o.High = price
		o.Low = price
		o.Set = true
	}
	if price.GreaterThan(o.High) {
		o.High = price
	}
	if price.LessThan(o.Low) {
		o.Low = price
	}
	o.Close = price
	o.Volume = volume
}

// PriceHistoryWorker records simulation ticks into the price_history table
// at 1-second granularity, then periodically flushes aggregated OHLCV bars
// for 1m, 5m, and 1h intervals.
type PriceHistoryWorker struct {
	repo   *repository.Repo
	logger *slog.Logger

	mu       sync.Mutex
	secondBuf map[string]*ohlcv // ticker -> current 1s bar
	minuteBuf map[string]*ohlcv // ticker -> current 1m bar
	fiveMinBuf map[string]*ohlcv // ticker -> current 5m bar
	hourBuf   map[string]*ohlcv // ticker -> current 1h bar
}

// NewPriceHistoryWorker creates a worker that persists price data.
func NewPriceHistoryWorker(repo *repository.Repo, logger *slog.Logger) *PriceHistoryWorker {
	return &PriceHistoryWorker{
		repo:       repo,
		logger:     logger,
		secondBuf:  make(map[string]*ohlcv),
		minuteBuf:  make(map[string]*ohlcv),
		fiveMinBuf: make(map[string]*ohlcv),
		hourBuf:    make(map[string]*ohlcv),
	}
}

// OnPriceBatch receives every simulation tick and buffers the data.
func (w *PriceHistoryWorker) OnPriceBatch(updates []simulation.PriceUpdate) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, u := range updates {
		// Buffer into all intervals
		if w.secondBuf[u.Ticker] == nil {
			w.secondBuf[u.Ticker] = &ohlcv{}
		}
		w.secondBuf[u.Ticker].update(u.Price, u.Volume)

		if w.minuteBuf[u.Ticker] == nil {
			w.minuteBuf[u.Ticker] = &ohlcv{}
		}
		w.minuteBuf[u.Ticker].update(u.Price, u.Volume)

		if w.fiveMinBuf[u.Ticker] == nil {
			w.fiveMinBuf[u.Ticker] = &ohlcv{}
		}
		w.fiveMinBuf[u.Ticker].update(u.Price, u.Volume)

		if w.hourBuf[u.Ticker] == nil {
			w.hourBuf[u.Ticker] = &ohlcv{}
		}
		w.hourBuf[u.Ticker].update(u.Price, u.Volume)
	}
}

// OnMarketEvent is a no-op for this worker.
func (w *PriceHistoryWorker) OnMarketEvent(_ simulation.MarketEvent) {}

// Run starts the flush loops. Blocks until context is cancelled.
func (w *PriceHistoryWorker) Run(ctx context.Context) {
	// Align flush times to clean boundaries
	secondTicker := time.NewTicker(1 * time.Second)
	minuteTicker := time.NewTicker(1 * time.Minute)
	fiveMinTicker := time.NewTicker(5 * time.Minute)
	hourTicker := time.NewTicker(1 * time.Hour)

	defer secondTicker.Stop()
	defer minuteTicker.Stop()
	defer fiveMinTicker.Stop()
	defer hourTicker.Stop()

	w.logger.Info("price history worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("price history worker stopped")
			return
		case <-secondTicker.C:
			w.flush(ctx, &w.secondBuf, "1s")
		case <-minuteTicker.C:
			w.flush(ctx, &w.minuteBuf, "1m")
		case <-fiveMinTicker.C:
			w.flush(ctx, &w.fiveMinBuf, "5m")
		case <-hourTicker.C:
			w.flush(ctx, &w.hourBuf, "1h")
		}
	}
}

// flush writes buffered OHLCV data to the database and resets the buffer.
func (w *PriceHistoryWorker) flush(ctx context.Context, buf *map[string]*ohlcv, interval string) {
	w.mu.Lock()
	snapshot := *buf
	*buf = make(map[string]*ohlcv)
	w.mu.Unlock()

	now := time.Now().UTC().Truncate(time.Second)
	written := 0

	for ticker, bar := range snapshot {
		if !bar.Set {
			continue
		}
		err := w.repo.InsertPriceHistory(ctx, ticker,
			bar.Close, bar.Open, bar.High, bar.Low, bar.Close,
			bar.Volume, interval, now)
		if err != nil {
			w.logger.Error("failed to write price history",
				"ticker", ticker, "interval", interval, "error", err)
			continue
		}
		written++
	}

	if written > 0 {
		w.logger.Debug("flushed price history", "interval", interval, "count", written)
	}
}
