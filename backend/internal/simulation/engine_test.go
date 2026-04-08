package simulation

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/luke/mockstarket/internal/market"
	"github.com/shopspring/decimal"
)

type mockObserver struct {
	mu       sync.Mutex
	batches  [][]market.PriceUpdate
	events   []market.MarketEvent
}

func (m *mockObserver) OnPriceBatch(updates []market.PriceUpdate) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.batches = append(m.batches, updates)
}

func (m *mockObserver) OnMarketEvent(event market.MarketEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

func (m *mockObserver) batchCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.batches)
}

func newTestEngine() *Engine {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	e := NewEngine(100, 0, 150, logger) // 100ms tick, no events
	e.AddStock("TEST", "Test Corp", "Tech", 100.0, 0.02, 0.0, 0.1)
	return e
}

func TestEngineAddStock(t *testing.T) {
	e := newTestEngine()

	price, ok := e.GetPrice("TEST")
	if !ok {
		t.Fatal("expected stock TEST to exist")
	}
	if !price.Equal(decimal.NewFromFloat(100.0)) {
		t.Errorf("expected price 100.0, got %s", price)
	}
}

func TestEngineGetAllPrices(t *testing.T) {
	e := newTestEngine()
	e.AddStock("FOO", "Foo Inc", "Consumer", 50.0, 0.03, 0.0, 0.1)

	prices := e.GetAllPrices()
	if len(prices) != 2 {
		t.Errorf("expected 2 stocks, got %d", len(prices))
	}
}

func TestEngineTickProducesPriceUpdates(t *testing.T) {
	e := newTestEngine()
	obs := &mockObserver{}
	e.AddObserver(obs)

	// Run a single tick
	e.tick()

	if obs.batchCount() != 1 {
		t.Errorf("expected 1 batch after 1 tick, got %d", obs.batchCount())
	}

	batch := obs.batches[0]
	if len(batch) != 1 {
		t.Errorf("expected 1 update in batch, got %d", len(batch))
	}
	if batch[0].Ticker != "TEST" {
		t.Errorf("expected ticker TEST, got %s", batch[0].Ticker)
	}
}

func TestEnginePricesStayPositive(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	e := NewEngine(10, 0, 150, logger)
	// Start at very low price with high volatility
	e.AddStock("PENNY", "Penny Stock", "Tech", 0.05, 0.50, -0.01, 0.1)

	for i := 0; i < 1000; i++ {
		e.tick()
	}

	price, _ := e.GetPrice("PENNY")
	if price.LessThanOrEqual(decimal.Zero) {
		t.Errorf("price should never go to zero or negative, got %s", price)
	}
}

func TestEngineMeanReversion(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	e := NewEngine(10, 0, 150, logger)
	e.AddStock("MEAN", "Mean Corp", "Tech", 100.0, 0.001, 0.0, 0.5) // Strong reversion, low volatility

	// Manually push price far from base
	e.stocks["MEAN"].Price = 200.0

	// Run many ticks — price should trend back toward 100
	for i := 0; i < 500; i++ {
		e.tick()
	}

	price, _ := e.GetPrice("MEAN")
	pf := price.InexactFloat64()

	// Should be closer to 100 than 200
	if pf > 160 {
		t.Errorf("expected mean reversion toward 100, but price is %.2f", pf)
	}
}

func TestEngineRunAndStop(t *testing.T) {
	e := newTestEngine()
	obs := &mockObserver{}
	e.AddObserver(obs)

	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	err := e.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded, got %v", err)
	}

	// Should have received at least 2 ticks in 350ms with 100ms interval
	count := obs.batchCount()
	if count < 2 {
		t.Errorf("expected at least 2 batches, got %d", count)
	}
}

func TestEngineMultipleStocksCorrelation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	e := NewEngine(10, 0, 150, logger)

	// Add stocks in same sector
	e.AddStock("A", "Stock A", "Tech", 100.0, 0.02, 0.0, 0.1)
	e.AddStock("B", "Stock B", "Tech", 100.0, 0.02, 0.0, 0.1)
	e.AddStock("C", "Stock C", "Food", 100.0, 0.02, 0.0, 0.1)

	obs := &mockObserver{}
	e.AddObserver(obs)

	e.tick()

	if obs.batchCount() != 1 {
		t.Fatal("expected 1 batch")
	}

	batch := obs.batches[0]
	if len(batch) != 3 {
		t.Errorf("expected 3 updates, got %d", len(batch))
	}
}
