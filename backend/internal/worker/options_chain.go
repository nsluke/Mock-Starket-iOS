package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/luke/mockstarket/internal/model"
	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/simulation"
	"github.com/shopspring/decimal"
)

// OptionsChainWorker generates option contracts for all stocks periodically.
type OptionsChainWorker struct {
	repo   *repository.Repo
	engine market.PriceProvider
	logger *slog.Logger
}

func NewOptionsChainWorker(repo *repository.Repo, engine market.PriceProvider, logger *slog.Logger) *OptionsChainWorker {
	return &OptionsChainWorker{repo: repo, engine: engine, logger: logger}
}

// Run generates option chains on startup and then every simulated day.
func (w *OptionsChainWorker) Run(ctx context.Context) {
	w.logger.Info("options chain worker started")

	// Generate immediately on startup
	w.generateAll(ctx)

	// Then regenerate periodically (every 5 real minutes ≈ 1 sim day)
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("options chain worker stopped")
			return
		case <-ticker.C:
			w.generateAll(ctx)
		}
	}
}

func (w *OptionsChainWorker) generateAll(ctx context.Context) {
	states := w.engine.GetAllStockStates()
	now := time.Now()
	expirations := simulation.GenerateExpirations(now)
	count := 0

	for _, state := range states {
		// Skip stocks with no price data (e.g., Polygon hasn't fetched yet)
		if state.Price <= 0 {
			continue
		}
		// Only generate options for regular stocks
		if state.Sector == "Crypto" || state.Sector == "Commodities" || state.Sector == "ETF" {
			continue
		}

		strikes := simulation.GenerateStrikes(state.Price)

		for _, exp := range expirations {
			for _, strike := range strikes {
				for _, optType := range []string{"call", "put"} {
					isCall := optType == "call"

					// Calculate initial pricing
					t := exp.Sub(now).Hours() / (24 * 365)
					if t <= 0 {
						continue
					}

					iv := simulation.SimulatedIV(state.Volatility, state.Price, strike, isCall, t)
					mark := simulation.BlackScholes(state.Price, strike, t, simulation.RiskFreeRate, iv, isCall)
					if mark < 0.01 {
						mark = 0.01
					}

					greeks := simulation.CalculateGreeks(state.Price, strike, t, simulation.RiskFreeRate, iv, isCall)

					symbol := simulation.BuildContractSymbol(state.Ticker, exp, optType, strike)

					contract := &model.OptionContract{
						Ticker:         state.Ticker,
						OptionType:     optType,
						StrikePrice:    decimal.NewFromFloat(strike),
						Expiration:     exp,
						ContractSymbol: symbol,
						BidPrice:       decimal.NewFromFloat(mark * 0.95),
						AskPrice:       decimal.NewFromFloat(mark * 1.05),
						MarkPrice:      decimal.NewFromFloat(mark),
						ImpliedVol:     decimal.NewFromFloat(iv),
						Delta:          decimal.NewFromFloat(greeks.Delta),
						Gamma:          decimal.NewFromFloat(greeks.Gamma),
						Theta:          decimal.NewFromFloat(greeks.Theta),
						Vega:           decimal.NewFromFloat(greeks.Vega),
						Rho:            decimal.NewFromFloat(greeks.Rho),
						Status:         "active",
					}

					if err := w.repo.UpsertOptionContract(ctx, contract); err != nil {
						w.logger.Error("failed to upsert option contract", "symbol", symbol, "error", err)
						continue
					}
					count++
				}
			}
		}
	}

	if count > 0 {
		w.logger.Info("generated option contracts", "count", count)
	}
}
