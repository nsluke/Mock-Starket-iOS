package worker

import (
	"testing"

	"github.com/luke/mockstarket/internal/market"
)

func TestStockSyncWorker_OnPriceBatch_BuffersLatest(t *testing.T) {
	w := NewStockSyncWorker(nil, 0, nil) // repo/logger nil — only testing buffering

	batch1 := []market.PriceUpdate{
		{Ticker: "AAPL", Price: d("150.00"), High: d("155.00"), Low: d("148.00"), Volume: 1000},
		{Ticker: "GOOG", Price: d("2800.00"), High: d("2850.00"), Low: d("2790.00"), Volume: 500},
	}
	w.OnPriceBatch(batch1)

	w.mu.Lock()
	if len(w.latest) != 2 {
		t.Errorf("expected 2 tickers, got %d", len(w.latest))
	}
	aapl := w.latest["AAPL"]
	if !aapl.Price.Equal(d("150.00")) {
		t.Errorf("expected AAPL price 150.00, got %s", aapl.Price)
	}
	w.mu.Unlock()

	// Second batch overwrites with newer prices
	batch2 := []market.PriceUpdate{
		{Ticker: "AAPL", Price: d("152.00"), High: d("156.00"), Low: d("147.00"), Volume: 1200},
	}
	w.OnPriceBatch(batch2)

	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.latest) != 2 {
		t.Errorf("expected still 2 tickers, got %d", len(w.latest))
	}
	aapl = w.latest["AAPL"]
	if !aapl.Price.Equal(d("152.00")) {
		t.Errorf("expected AAPL price updated to 152.00, got %s", aapl.Price)
	}
	// GOOG should still be there from batch1
	goog := w.latest["GOOG"]
	if !goog.Price.Equal(d("2800.00")) {
		t.Errorf("expected GOOG price 2800.00, got %s", goog.Price)
	}
}

func TestStockSyncWorker_OnPriceBatch_CapturesAllFields(t *testing.T) {
	w := NewStockSyncWorker(nil, 0, nil)

	w.OnPriceBatch([]market.PriceUpdate{
		{Ticker: "TEST", Price: d("42.50"), High: d("45.00"), Low: d("40.00"), Volume: 9999},
	})

	w.mu.Lock()
	defer w.mu.Unlock()

	snap := w.latest["TEST"]
	if !snap.Price.Equal(d("42.50")) {
		t.Errorf("price: got %s, want 42.50", snap.Price)
	}
	if !snap.High.Equal(d("45.00")) {
		t.Errorf("high: got %s, want 45.00", snap.High)
	}
	if !snap.Low.Equal(d("40.00")) {
		t.Errorf("low: got %s, want 40.00", snap.Low)
	}
	if snap.Volume != 9999 {
		t.Errorf("volume: got %d, want 9999", snap.Volume)
	}
}
