package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/simulation"
	ws "github.com/luke/mockstarket/internal/websocket"
	"github.com/shopspring/decimal"
)

// OptionsPricingWorker recalculates option prices and greeks on each price tick.
type OptionsPricingWorker struct {
	repo   *repository.Repo
	engine market.PriceProvider
	hub    *ws.Hub
	logger *slog.Logger
	mu     sync.Mutex
	count  int64
}

func NewOptionsPricingWorker(repo *repository.Repo, engine market.PriceProvider, hub *ws.Hub, logger *slog.Logger) *OptionsPricingWorker {
	return &OptionsPricingWorker{repo: repo, engine: engine, hub: hub, logger: logger}
}

func (w *OptionsPricingWorker) OnPriceBatch(_ []market.PriceUpdate) {
	w.mu.Lock()
	w.count++
	tick := w.count
	w.mu.Unlock()

	// Only reprice every 5th tick to reduce DB load
	if tick%5 != 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	contracts, err := w.repo.GetAllActiveContracts(ctx)
	if err != nil {
		w.logger.Error("options pricing: failed to load contracts", "error", err)
		return
	}

	if len(contracts) == 0 {
		return
	}

	now := time.Now()
	states := w.engine.GetAllStockStates()

	// Group updates by ticker for broadcasting
	tickerUpdates := make(map[string][]map[string]interface{})

	for i := range contracts {
		c := &contracts[i]
		state, ok := states[c.Ticker]
		if !ok || state.Price <= 0 {
			continue
		}

		strike := c.StrikePrice.InexactFloat64()
		isCall := c.OptionType == "call"
		t := c.Expiration.Sub(now).Hours() / (24 * 365)
		if t <= 0 {
			t = 0
		}

		iv := simulation.SimulatedIV(state.Volatility, state.Price, strike, isCall, t)
		mark := simulation.BlackScholes(state.Price, strike, t, simulation.RiskFreeRate, iv, isCall)
		if mark < 0.01 {
			mark = 0.01
		}

		greeks := simulation.CalculateGreeks(state.Price, strike, t, simulation.RiskFreeRate, iv, isCall)

		c.MarkPrice = decimal.NewFromFloat(mark).Round(4)
		c.BidPrice = decimal.NewFromFloat(mark * 0.95).Round(4)
		c.AskPrice = decimal.NewFromFloat(mark * 1.05).Round(4)
		c.ImpliedVol = decimal.NewFromFloat(iv).Round(6)
		c.Delta = decimal.NewFromFloat(greeks.Delta).Round(6)
		c.Gamma = decimal.NewFromFloat(greeks.Gamma).Round(6)
		c.Theta = decimal.NewFromFloat(greeks.Theta).Round(6)
		c.Vega = decimal.NewFromFloat(greeks.Vega).Round(6)
		c.Rho = decimal.NewFromFloat(greeks.Rho).Round(6)

		tickerUpdates[c.Ticker] = append(tickerUpdates[c.Ticker], map[string]interface{}{
			"id":         c.ID,
			"mark_price": c.MarkPrice,
			"bid_price":  c.BidPrice,
			"ask_price":  c.AskPrice,
			"delta":      c.Delta,
			"gamma":      c.Gamma,
			"theta":      c.Theta,
			"vega":       c.Vega,
			"rho":        c.Rho,
			"implied_vol": c.ImpliedVol,
		})
	}

	// Batch update DB
	if err := w.repo.BatchUpdateOptionPrices(ctx, contracts); err != nil {
		w.logger.Error("options pricing: failed to batch update", "error", err)
	}

	// Broadcast via WebSocket
	for ticker, updates := range tickerUpdates {
		data, _ := json.Marshal(updates)
		w.hub.BroadcastToChannel("options:"+ticker, ws.Message{
			Type: "options_price_batch",
			Data: data,
		})
	}
}

func (w *OptionsPricingWorker) OnMarketEvent(_ market.MarketEvent) {}
