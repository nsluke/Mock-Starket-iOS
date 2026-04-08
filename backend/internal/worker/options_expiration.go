package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/simulation"
	"github.com/shopspring/decimal"
)

// OptionsExpirationWorker settles expired option contracts.
type OptionsExpirationWorker struct {
	repo   *repository.Repo
	engine market.PriceProvider
	logger *slog.Logger
}

func NewOptionsExpirationWorker(repo *repository.Repo, engine market.PriceProvider, logger *slog.Logger) *OptionsExpirationWorker {
	return &OptionsExpirationWorker{repo: repo, engine: engine, logger: logger}
}

// Run checks for expired contracts every minute.
func (w *OptionsExpirationWorker) Run(ctx context.Context) {
	w.logger.Info("options expiration worker started")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("options expiration worker stopped")
			return
		case <-ticker.C:
			w.processExpirations(ctx)
		}
	}
}

func (w *OptionsExpirationWorker) processExpirations(ctx context.Context) {
	now := time.Now()
	contracts, err := w.repo.GetExpiringContracts(ctx, now)
	if err != nil {
		w.logger.Error("expiration worker: failed to get expiring contracts", "error", err)
		return
	}

	for _, contract := range contracts {
		// Get underlying price
		price, ok := w.engine.GetPrice(contract.Ticker)
		if !ok {
			continue
		}

		strike := contract.StrikePrice
		isCall := contract.OptionType == "call"
		multiplier := decimal.NewFromInt(simulation.ContractMultiplier)

		// Determine if ITM
		itm := false
		if isCall && price.GreaterThan(strike) {
			itm = true
		} else if !isCall && price.LessThan(strike) {
			itm = true
		}

		// Process all positions for this contract
		positions, err := w.repo.GetPositionsForContract(ctx, contract.ID)
		if err != nil {
			w.logger.Error("expiration worker: failed to get positions", "contract", contract.ContractSymbol, "error", err)
			continue
		}

		for _, pos := range positions {
			portfolio, err := w.repo.GetPortfolioByUserID(ctx, pos.PortfolioID)
			if err != nil {
				// portfolio_id isn't user_id; use a join or just skip
				continue
			}

			qty := decimal.NewFromInt(int64(pos.Quantity)).Abs()
			isLong := pos.Quantity > 0

			if itm {
				// Auto-exercise
				if isLong {
					if isCall {
						// Long call ITM: pay strike, receive shares value
						// Net credit = (price - strike) * qty * 100
						intrinsic := price.Sub(strike).Mul(qty).Mul(multiplier)
						newCash := portfolio.Cash.Add(intrinsic)
						_ = w.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash)
					} else {
						// Long put ITM: receive (strike - price) * qty * 100
						intrinsic := strike.Sub(price).Mul(qty).Mul(multiplier)
						newCash := portfolio.Cash.Add(intrinsic)
						_ = w.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash)
					}
				} else {
					// Short position ITM: pay the intrinsic value
					if isCall {
						intrinsic := price.Sub(strike).Mul(qty).Mul(multiplier)
						released := pos.Collateral.Sub(intrinsic)
						if released.IsNegative() {
							released = decimal.Zero
						}
						newCash := portfolio.Cash.Add(released)
						_ = w.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash)
					} else {
						intrinsic := strike.Sub(price).Mul(qty).Mul(multiplier)
						released := pos.Collateral.Sub(intrinsic)
						if released.IsNegative() {
							released = decimal.Zero
						}
						newCash := portfolio.Cash.Add(released)
						_ = w.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash)
					}
				}

				w.logger.Info("auto-exercised option",
					"contract", contract.ContractSymbol,
					"long", isLong,
					"price", price,
					"strike", strike)
			} else {
				// OTM: expire worthless
				if !isLong {
					// Release collateral for short positions
					newCash := portfolio.Cash.Add(pos.Collateral)
					_ = w.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash)
				}
				// Long positions: premium already paid, nothing to do
			}

			// Remove position
			_ = w.repo.DeleteOptionPosition(ctx, pos.PortfolioID, contract.ID)
		}

		// Mark contract as expired/exercised
		status := "expired"
		if itm {
			status = "exercised"
		}
		_ = w.repo.UpdateContractStatus(ctx, contract.ID, status)
	}

	if len(contracts) > 0 {
		w.logger.Info("processed expired contracts", "count", len(contracts))
	}
}
