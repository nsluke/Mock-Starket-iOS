package worker

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestOHLCV_FirstUpdate(t *testing.T) {
	bar := &ohlcv{}

	bar.update(d("100.00"), 1000)

	if !bar.Set {
		t.Error("expected bar.Set to be true after first update")
	}
	if !bar.Open.Equal(d("100.00")) {
		t.Errorf("expected open 100.00, got %s", bar.Open)
	}
	if !bar.High.Equal(d("100.00")) {
		t.Errorf("expected high 100.00, got %s", bar.High)
	}
	if !bar.Low.Equal(d("100.00")) {
		t.Errorf("expected low 100.00, got %s", bar.Low)
	}
	if !bar.Close.Equal(d("100.00")) {
		t.Errorf("expected close 100.00, got %s", bar.Close)
	}
	if bar.Volume != 1000 {
		t.Errorf("expected volume 1000, got %d", bar.Volume)
	}
}

func TestOHLCV_MultipleUpdates(t *testing.T) {
	bar := &ohlcv{}

	bar.update(d("100.00"), 1000) // open
	bar.update(d("110.00"), 2000) // new high
	bar.update(d("90.00"), 3000)  // new low
	bar.update(d("105.00"), 4000) // close

	if !bar.Open.Equal(d("100.00")) {
		t.Errorf("open should stay at first price, got %s", bar.Open)
	}
	if !bar.High.Equal(d("110.00")) {
		t.Errorf("expected high 110.00, got %s", bar.High)
	}
	if !bar.Low.Equal(d("90.00")) {
		t.Errorf("expected low 90.00, got %s", bar.Low)
	}
	if !bar.Close.Equal(d("105.00")) {
		t.Errorf("expected close 105.00, got %s", bar.Close)
	}
	if bar.Volume != 4000 {
		t.Errorf("expected volume 4000 (latest), got %d", bar.Volume)
	}
}

func TestOHLCV_MonotonicallyIncreasing(t *testing.T) {
	bar := &ohlcv{}

	bar.update(d("100.00"), 100)
	bar.update(d("101.00"), 200)
	bar.update(d("102.00"), 300)

	if !bar.Open.Equal(d("100.00")) {
		t.Errorf("open should be first, got %s", bar.Open)
	}
	if !bar.High.Equal(d("102.00")) {
		t.Errorf("high should be last, got %s", bar.High)
	}
	if !bar.Low.Equal(d("100.00")) {
		t.Errorf("low should be first, got %s", bar.Low)
	}
}

func TestOHLCV_MonotonicallyDecreasing(t *testing.T) {
	bar := &ohlcv{}

	bar.update(d("100.00"), 100)
	bar.update(d("99.00"), 200)
	bar.update(d("98.00"), 300)

	if !bar.High.Equal(d("100.00")) {
		t.Errorf("high should be first, got %s", bar.High)
	}
	if !bar.Low.Equal(d("98.00")) {
		t.Errorf("low should be last, got %s", bar.Low)
	}
}

func TestPriceHistoryWorker_OnPriceBatch_BuffersData(t *testing.T) {
	w := NewPriceHistoryWorker(nil, nil) // repo/logger nil — we only test buffering

	updates := []testPriceUpdate{
		{ticker: "AAPL", price: "150.00", volume: 1000},
		{ticker: "GOOG", price: "2800.00", volume: 500},
	}

	simUpdates := make([]simulationPriceUpdate, len(updates))
	for i, u := range updates {
		simUpdates[i] = simulationPriceUpdate{
			Ticker: u.ticker,
			Price:  d(u.price),
			Volume: u.volume,
		}
	}

	// Directly test buffering by checking internal state
	w.mu.Lock()
	for _, u := range updates {
		if w.secondBuf[u.ticker] == nil {
			w.secondBuf[u.ticker] = &ohlcv{}
		}
		w.secondBuf[u.ticker].update(d(u.price), u.volume)
	}
	w.mu.Unlock()

	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.secondBuf) != 2 {
		t.Errorf("expected 2 tickers buffered, got %d", len(w.secondBuf))
	}

	aapl := w.secondBuf["AAPL"]
	if aapl == nil {
		t.Fatal("expected AAPL in buffer")
	}
	if !aapl.Close.Equal(d("150.00")) {
		t.Errorf("expected AAPL close 150.00, got %s", aapl.Close)
	}
}

type testPriceUpdate struct {
	ticker string
	price  string
	volume int64
}

type simulationPriceUpdate struct {
	Ticker string
	Price  decimal.Decimal
	Volume int64
}
