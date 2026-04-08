package simulation

import (
	"context"
	"log/slog"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// StockState holds the live simulation state for a single stock.
type StockState struct {
	Ticker             string
	Price              float64
	BasePrice          float64
	Volatility         float64
	Drift              float64
	MeanReversionSpeed float64
	Sector             string
	DayOpen            float64
	DayHigh            float64
	DayLow             float64
	Volume             int64
}

// PriceUpdate is broadcast after each simulation tick.
type PriceUpdate struct {
	Ticker    string          `json:"ticker"`
	Price     decimal.Decimal `json:"price"`
	Change    decimal.Decimal `json:"change"`
	ChangePct decimal.Decimal `json:"change_pct"`
	Volume    int64           `json:"volume"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
}

// MarketEvent represents a simulation event that affects prices.
type MarketEvent struct {
	Type      string `json:"event"`
	Ticker    string `json:"ticker,omitempty"`
	Sector    string `json:"sector,omitempty"`
	Headline  string `json:"headline"`
	Impact    string `json:"impact"`
	Magnitude string `json:"magnitude"`
}

// Observer receives simulation updates.
type Observer interface {
	OnPriceBatch(updates []PriceUpdate)
	OnMarketEvent(event MarketEvent)
}

// Engine drives the stock price simulation.
type Engine struct {
	mu            sync.RWMutex
	stocks        map[string]*StockState
	sectorFactors map[string]float64
	marketFactor  float64
	tickInterval  time.Duration
	ticksPerDay   int
	eventFreq     int
	tickCount     int64
	observers     []Observer
	rng           *rand.Rand
	logger        *slog.Logger
}

// NewEngine creates a simulation engine with the given tick interval.
// ticksPerDay controls simulation speed: volatility/drift params represent
// daily magnitudes, and dt = 1/ticksPerDay normalizes each tick accordingly.
// With default 150 ticks/day at 2s/tick, one simulated day ≈ 5 real minutes.
func NewEngine(tickMS int, eventFreq int, ticksPerDay int, logger *slog.Logger) *Engine {
	if ticksPerDay <= 0 {
		ticksPerDay = 150
	}
	return &Engine{
		stocks:        make(map[string]*StockState),
		sectorFactors: make(map[string]float64),
		tickInterval:  time.Duration(tickMS) * time.Millisecond,
		ticksPerDay:   ticksPerDay,
		eventFreq:     eventFreq,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		logger:        logger,
	}
}

// AddStock registers a stock in the simulation.
func (e *Engine) AddStock(ticker, name, sector string, price, volatility, drift, meanReversion float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.stocks[ticker] = &StockState{
		Ticker:             ticker,
		Price:              price,
		BasePrice:          price,
		Volatility:         volatility,
		Drift:              drift,
		MeanReversionSpeed: meanReversion,
		Sector:             sector,
		DayOpen:            price,
		DayHigh:            price,
		DayLow:             price,
		Volume:             0,
	}
}

// AddObserver registers a listener for price updates and events.
func (e *Engine) AddObserver(obs Observer) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.observers = append(e.observers, obs)
}

// GetPrice returns the current price for a ticker.
func (e *Engine) GetPrice(ticker string) (decimal.Decimal, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	s, ok := e.stocks[ticker]
	if !ok {
		return decimal.Zero, false
	}
	return decimal.NewFromFloat(s.Price).Round(4), true
}

// GetStockState returns a copy of the simulation state for a ticker.
func (e *Engine) GetStockState(ticker string) (*StockState, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	s, ok := e.stocks[ticker]
	if !ok {
		return nil, false
	}
	cp := *s
	return &cp, true
}

// GetAllStockStates returns a copy of all stock states.
func (e *Engine) GetAllStockStates() map[string]StockState {
	e.mu.RLock()
	defer e.mu.RUnlock()

	states := make(map[string]StockState, len(e.stocks))
	for k, v := range e.stocks {
		states[k] = *v
	}
	return states
}

// GetTicksPerDay returns the number of ticks per simulated day.
func (e *Engine) GetTicksPerDay() int {
	return e.ticksPerDay
}

// GetTickCount returns the current tick count.
func (e *Engine) GetTickCount() int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.tickCount
}

// GetAllPrices returns current prices for all stocks.
func (e *Engine) GetAllPrices() map[string]decimal.Decimal {
	e.mu.RLock()
	defer e.mu.RUnlock()

	prices := make(map[string]decimal.Decimal, len(e.stocks))
	for ticker, s := range e.stocks {
		prices[ticker] = decimal.NewFromFloat(s.Price).Round(4)
	}
	return prices
}

// Run starts the simulation loop. Blocks until context is cancelled.
func (e *Engine) Run(ctx context.Context) error {
	e.logger.Info("simulation engine started", "tick_interval", e.tickInterval)

	ticker := time.NewTicker(e.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.logger.Info("simulation engine stopped")
			return ctx.Err()
		case <-ticker.C:
			e.tick()
		}
	}
}

// tick performs one simulation step.
func (e *Engine) tick() {
	e.mu.Lock()

	e.tickCount++

	dt := 1.0 / float64(e.ticksPerDay)

	// Generate market-wide and sector noise (small shared factors)
	e.marketFactor = e.rng.NormFloat64() * 0.3
	sectors := make(map[string]bool)
	for _, s := range e.stocks {
		sectors[s.Sector] = true
	}
	for sector := range sectors {
		e.sectorFactors[sector] = e.rng.NormFloat64() * 0.2
	}

	updates := make([]PriceUpdate, 0, len(e.stocks))

	for _, s := range e.stocks {
		prevPrice := s.Price

		// Simple per-tick model:
		// - Volatility is the per-tick standard deviation as a fraction of price
		//   e.g. 0.0001 = 0.01% per tick
		// - Noise is a blend of individual + market + sector factors
		// - Mean reversion gently pulls price back toward base
		individualNoise := e.rng.NormFloat64()
		combinedNoise := 0.5*individualNoise + 0.3*e.marketFactor + 0.2*e.sectorFactors[s.Sector]

		noise := s.Volatility * s.Price * combinedNoise
		meanReversion := s.MeanReversionSpeed * (s.BasePrice - s.Price) * dt
		drift := s.Drift * s.Price * dt

		dp := noise + meanReversion + drift

		// Cap per-tick change at ±0.5%
		maxChange := s.Price * 0.005
		if dp > maxChange {
			dp = maxChange
		} else if dp < -maxChange {
			dp = -maxChange
		}

		s.Price += dp

		// Ensure price stays positive (floor at $0.01)
		if s.Price < 0.01 {
			s.Price = 0.01
		}

		// Update day statistics
		if s.Price > s.DayHigh {
			s.DayHigh = s.Price
		}
		if s.Price < s.DayLow {
			s.DayLow = s.Price
		}

		// Simulate volume (random, correlated with price movement magnitude)
		priceMovePct := math.Abs(s.Price-prevPrice) / prevPrice
		baseVolume := 100 + int64(priceMovePct*50000)
		s.Volume += baseVolume + int64(e.rng.Intn(200))

		change := s.Price - s.DayOpen
		changePct := 0.0
		if s.DayOpen > 0 {
			changePct = (change / s.DayOpen) * 100
		}

		updates = append(updates, PriceUpdate{
			Ticker:    s.Ticker,
			Price:     decimal.NewFromFloat(s.Price).Round(4),
			Change:    decimal.NewFromFloat(change).Round(4),
			ChangePct: decimal.NewFromFloat(changePct).Round(2),
			Volume:    s.Volume,
			High:      decimal.NewFromFloat(s.DayHigh).Round(4),
			Low:       decimal.NewFromFloat(s.DayLow).Round(4),
		})
	}

	// Copy observers under lock
	observers := make([]Observer, len(e.observers))
	copy(observers, e.observers)
	e.mu.Unlock()

	// Notify observers outside lock
	for _, obs := range observers {
		obs.OnPriceBatch(updates)
	}

	// Periodic market events
	if e.eventFreq > 0 && e.tickCount%int64(e.eventFreq) == 0 {
		if event, ok := e.generateEvent(); ok {
			for _, obs := range observers {
				obs.OnMarketEvent(event)
			}
		}
	}
}

// generateEvent creates a random market event.
func (e *Engine) generateEvent() (MarketEvent, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	roll := e.rng.Float64()

	switch {
	case roll < 0.3:
		// Stock-specific earnings event — temporary price shock, not permanent base shift
		stock := e.randomStock()
		if stock == nil {
			return MarketEvent{}, false
		}
		positive := e.rng.Float64() > 0.4
		magnitude := 0.005 + e.rng.Float64()*0.015 // 0.5-2% price shock
		impact := "positive"
		headline := stock.Ticker + " reports strong quarterly earnings"
		if !positive {
			magnitude = -magnitude
			impact = "negative"
			headline = stock.Ticker + " misses earnings expectations"
		}
		stock.Price *= (1 + magnitude)
		return MarketEvent{
			Type:      "earnings_surprise",
			Ticker:    stock.Ticker,
			Headline:  headline,
			Impact:    impact,
			Magnitude: magnitudeLabel(math.Abs(magnitude)),
		}, true

	case roll < 0.5:
		// Sector event — temporary price shock
		sectors := []string{"Tech", "Consumer", "Defense", "Food", "Industrial"}
		sector := sectors[e.rng.Intn(len(sectors))]
		positive := e.rng.Float64() > 0.5
		shift := 0.002 + e.rng.Float64()*0.008 // 0.2-1% price shock
		impact := "positive"
		headline := sector + " sector surges on strong demand"
		if !positive {
			shift = -shift
			impact = "negative"
			headline = sector + " sector drops amid regulatory concerns"
		}
		for _, s := range e.stocks {
			if s.Sector == sector {
				s.Price *= (1 + shift)
			}
		}
		return MarketEvent{
			Type:      "sector_event",
			Sector:    sector,
			Headline:  headline,
			Impact:    impact,
			Magnitude: magnitudeLabel(math.Abs(shift)),
		}, true

	case roll < 0.65:
		// Market-wide event — temporary price shock
		positive := e.rng.Float64() > 0.5
		shift := 0.001 + e.rng.Float64()*0.004 // 0.1-0.5% price shock
		impact := "positive"
		headline := "Federal Reserve signals accommodative policy"
		if !positive {
			shift = -shift
			impact = "negative"
			headline = "Global markets tumble on economic uncertainty"
		}
		for _, s := range e.stocks {
			s.Price *= (1 + shift)
		}
		return MarketEvent{
			Type:      "market_event",
			Headline:  headline,
			Impact:    impact,
			Magnitude: magnitudeLabel(math.Abs(shift)),
		}, true

	default:
		return MarketEvent{}, false
	}
}

func (e *Engine) randomStock() *StockState {
	keys := make([]string, 0, len(e.stocks))
	for k := range e.stocks {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return nil
	}
	return e.stocks[keys[e.rng.Intn(len(keys))]]
}

func magnitudeLabel(pct float64) string {
	switch {
	case pct > 0.15:
		return "extreme"
	case pct > 0.10:
		return "high"
	case pct > 0.05:
		return "medium"
	default:
		return "low"
	}
}
