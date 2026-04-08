package polygon

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/luke/mockstarket/internal/market"
	"github.com/shopspring/decimal"
)

// FeedConfig configures the MarketFeed behavior.
type FeedConfig struct {
	WSEnabled    bool
	WSURL        string // default: "wss://socket.polygon.io/stocks"
	PollInterval time.Duration
}

// MarketFeed implements market.PriceProvider using Polygon.io data.
type MarketFeed struct {
	client *Client
	config FeedConfig
	logger *slog.Logger

	mu        sync.RWMutex
	stocks    map[string]*market.StockState
	tickers   map[string]tickerInfo // ticker -> metadata
	observers []market.Observer
}

type tickerInfo struct {
	Sector    string
	AssetType string
}

// NewMarketFeed creates a new Polygon.io-backed market data feed.
func NewMarketFeed(client *Client, config FeedConfig, logger *slog.Logger) *MarketFeed {
	if config.PollInterval == 0 {
		config.PollInterval = 30 * time.Second
	}
	if config.WSURL == "" {
		config.WSURL = "wss://socket.polygon.io/stocks"
	}
	return &MarketFeed{
		client:  client,
		config:  config,
		logger:  logger,
		stocks:  make(map[string]*market.StockState),
		tickers: make(map[string]tickerInfo),
	}
}

// TrackTicker adds a ticker to be tracked by this feed.
func (f *MarketFeed) TrackTicker(ticker, sector, assetType string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tickers[ticker] = tickerInfo{Sector: sector, AssetType: assetType}
	if _, exists := f.stocks[ticker]; !exists {
		f.stocks[ticker] = &market.StockState{
			Ticker: ticker,
			Sector: sector,
		}
	}
}

// AddObserver registers a listener for price updates.
func (f *MarketFeed) AddObserver(obs market.Observer) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.observers = append(f.observers, obs)
}

// GetPrice returns the current price for a ticker.
func (f *MarketFeed) GetPrice(ticker string) (decimal.Decimal, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	s, ok := f.stocks[ticker]
	if !ok || s.Price == 0 {
		return decimal.Zero, false
	}
	return decimal.NewFromFloat(s.Price).Round(4), true
}

// GetAllPrices returns current prices for all tracked stocks.
func (f *MarketFeed) GetAllPrices() map[string]decimal.Decimal {
	f.mu.RLock()
	defer f.mu.RUnlock()
	prices := make(map[string]decimal.Decimal, len(f.stocks))
	for ticker, s := range f.stocks {
		if s.Price > 0 {
			prices[ticker] = decimal.NewFromFloat(s.Price).Round(4)
		}
	}
	return prices
}

// GetAllStockStates returns a copy of all stock states.
func (f *MarketFeed) GetAllStockStates() map[string]market.StockState {
	f.mu.RLock()
	defer f.mu.RUnlock()
	states := make(map[string]market.StockState, len(f.stocks))
	for k, v := range f.stocks {
		states[k] = *v
	}
	return states
}

// Run starts the market data feed. Blocks until context is cancelled.
func (f *MarketFeed) Run(ctx context.Context) error {
	f.logger.Info("polygon market feed starting",
		"tickers", len(f.tickers),
		"ws_enabled", f.config.WSEnabled,
		"poll_interval", f.config.PollInterval,
	)

	// Initial snapshot fetch
	if err := f.fetchSnapshots(ctx); err != nil {
		f.logger.Error("initial snapshot fetch failed", "error", err)
		// Continue anyway — will retry on next poll
	}

	// Start WebSocket if enabled
	if f.config.WSEnabled {
		go f.runWebSocket(ctx)
	}

	// Polling loop
	return f.pollLoop(ctx)
}

func (f *MarketFeed) pollLoop(ctx context.Context) error {
	ticker := time.NewTicker(f.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			f.logger.Info("polygon market feed stopped")
			return ctx.Err()
		case <-ticker.C:
			// Reduce frequency outside market hours for equities
			session := GetMarketSession(time.Now())
			if session == SessionClosed {
				// Poll less frequently when market is closed
				time.Sleep(4 * f.config.PollInterval)
			}

			if err := f.fetchSnapshots(ctx); err != nil {
				f.logger.Error("snapshot poll failed", "error", err)
			}
		}
	}
}

func (f *MarketFeed) fetchSnapshots(ctx context.Context) error {
	// Try bulk endpoint first (requires paid plan), fall back to per-ticker
	snapshots, err := f.client.GetAllSnapshots(ctx)
	if err != nil {
		f.logger.Debug("bulk snapshots unavailable, using per-ticker fallback", "error", err)
		return f.fetchSnapshotsPerTicker(ctx)
	}

	// Build a lookup of tracked tickers
	f.mu.RLock()
	tracked := make(map[string]bool, len(f.tickers))
	for t := range f.tickers {
		tracked[t] = true
	}
	f.mu.RUnlock()

	// Filter to only our tracked tickers and build updates
	var updates []market.PriceUpdate

	f.mu.Lock()
	for _, snap := range snapshots {
		if !tracked[snap.Ticker] {
			continue
		}

		state, ok := f.stocks[snap.Ticker]
		if !ok {
			continue
		}

		prevPrice := state.Price
		state.Price = snap.Day.Close
		if state.Price == 0 {
			state.Price = snap.Min.Close // Fall back to last minute bar
		}
		state.DayOpen = snap.Day.Open
		state.DayHigh = snap.Day.High
		state.DayLow = snap.Day.Low
		state.BasePrice = snap.PrevDay.Close
		state.Volume = int64(snap.Day.Volume)

		// Set a default volatility for options pricing (20% annual)
		if state.Volatility == 0 {
			state.Volatility = 0.0010 // ~20% annual at 150 ticks/day * 252 days
		}

		if state.Price > 0 {
			change := state.Price - state.DayOpen
			changePct := 0.0
			if state.DayOpen > 0 {
				changePct = (change / state.DayOpen) * 100
			}

			updates = append(updates, market.PriceUpdate{
				Ticker:    snap.Ticker,
				Price:     decimal.NewFromFloat(state.Price).Round(4),
				Change:    decimal.NewFromFloat(change).Round(4),
				ChangePct: decimal.NewFromFloat(changePct).Round(2),
				Volume:    state.Volume,
				High:      decimal.NewFromFloat(state.DayHigh).Round(4),
				Low:       decimal.NewFromFloat(state.DayLow).Round(4),
			})

			_ = prevPrice // suppress unused warning
		}
	}

	// Copy observers while holding the lock
	observers := make([]market.Observer, len(f.observers))
	copy(observers, f.observers)
	f.mu.Unlock()

	// Notify observers
	if len(updates) > 0 {
		for _, obs := range observers {
			obs.OnPriceBatch(updates)
		}
		f.logger.Debug("polygon price update broadcast", "tickers", len(updates))
	}

	return nil
}

// fetchSnapshotsPerTicker uses the previous-close endpoint (available on free tier)
// to fetch prices for tracked tickers. Fetches up to 4 per poll cycle to stay
// within the 5 req/min free tier rate limit.
func (f *MarketFeed) fetchSnapshotsPerTicker(ctx context.Context) error {
	f.mu.RLock()
	tickers := make([]string, 0, len(f.tickers))
	for t := range f.tickers {
		tickers = append(tickers, t)
	}
	f.mu.RUnlock()

	const batchSize = 4
	var updates []market.PriceUpdate

	for i, ticker := range tickers {
		if i >= batchSize {
			break
		}

		bar, err := f.client.GetPreviousClose(ctx, ticker)
		if err != nil {
			f.logger.Debug("prev close fetch failed", "ticker", ticker, "error", err)
			continue
		}

		f.mu.Lock()
		state, ok := f.stocks[ticker]
		if ok {
			state.Price = bar.Close
			state.DayOpen = bar.Open
			state.DayHigh = bar.High
			state.DayLow = bar.Low
			state.BasePrice = bar.Open
			state.Volume = int64(bar.Volume)

			if state.Volatility == 0 {
				state.Volatility = 0.0010
			}

			if state.Price > 0 {
				change := state.Price - state.DayOpen
				changePct := 0.0
				if state.DayOpen > 0 {
					changePct = (change / state.DayOpen) * 100
				}
				updates = append(updates, market.PriceUpdate{
					Ticker:    ticker,
					Price:     decimal.NewFromFloat(state.Price).Round(4),
					Change:    decimal.NewFromFloat(change).Round(4),
					ChangePct: decimal.NewFromFloat(changePct).Round(2),
					Volume:    state.Volume,
					High:      decimal.NewFromFloat(state.DayHigh).Round(4),
					Low:       decimal.NewFromFloat(state.DayLow).Round(4),
				})
			}
		}
		f.mu.Unlock()
	}

	// Notify observers
	if len(updates) > 0 {
		f.mu.RLock()
		observers := make([]market.Observer, len(f.observers))
		copy(observers, f.observers)
		f.mu.RUnlock()

		for _, obs := range observers {
			obs.OnPriceBatch(updates)
		}
		f.logger.Debug("polygon per-ticker update", "tickers", len(updates))
	}

	return nil
}

func (f *MarketFeed) runWebSocket(ctx context.Context) {
	ws := NewWSClient(f.config.WSURL, f.client.apiKey, f.logger)

	// Collect tickers to subscribe
	f.mu.RLock()
	tickers := make([]string, 0, len(f.tickers))
	for t := range f.tickers {
		tickers = append(tickers, t)
	}
	f.mu.RUnlock()

	// Subscribe to per-minute aggregates
	if err := ws.Subscribe("AM", tickers); err != nil {
		f.logger.Error("ws subscribe failed", "error", err)
	}

	// Handle aggregate messages
	ws.OnAggregate(func(agg WSAggregateMessage) {
		f.handleWSAggregate(agg)
	})

	// Run will reconnect automatically
	if err := ws.Run(ctx); err != nil && err != context.Canceled {
		f.logger.Error("polygon websocket error", "error", err)
	}
}

func (f *MarketFeed) handleWSAggregate(agg WSAggregateMessage) {
	f.mu.Lock()
	state, ok := f.stocks[agg.Ticker]
	if !ok {
		f.mu.Unlock()
		return
	}

	state.Price = agg.Close
	if agg.High > state.DayHigh {
		state.DayHigh = agg.High
	}
	if state.DayLow == 0 || agg.Low < state.DayLow {
		state.DayLow = agg.Low
	}
	state.Volume += int64(agg.Volume)

	change := state.Price - state.DayOpen
	changePct := 0.0
	if state.DayOpen > 0 {
		changePct = (change / state.DayOpen) * 100
	}

	update := market.PriceUpdate{
		Ticker:    agg.Ticker,
		Price:     decimal.NewFromFloat(state.Price).Round(4),
		Change:    decimal.NewFromFloat(change).Round(4),
		ChangePct: decimal.NewFromFloat(changePct).Round(2),
		Volume:    state.Volume,
		High:      decimal.NewFromFloat(state.DayHigh).Round(4),
		Low:       decimal.NewFromFloat(state.DayLow).Round(4),
	}

	observers := make([]market.Observer, len(f.observers))
	copy(observers, f.observers)
	f.mu.Unlock()

	for _, obs := range observers {
		obs.OnPriceBatch([]market.PriceUpdate{update})
	}
}
