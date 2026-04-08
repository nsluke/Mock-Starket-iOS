package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/luke/mockstarket/internal/model"
	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/simulation"
	"github.com/shopspring/decimal"
)

var (
	contractMultiplier = decimal.NewFromInt(simulation.ContractMultiplier)
)

// OptionsTradeService handles options trade execution and validation.
type OptionsTradeService struct {
	repo   *repository.Repo
	engine market.PriceProvider
}

// NewOptionsTradeService creates a new options trade service.
func NewOptionsTradeService(repo *repository.Repo, engine market.PriceProvider) *OptionsTradeService {
	return &OptionsTradeService{repo: repo, engine: engine}
}

// OptionsTradeRequest represents an options buy/sell request.
type OptionsTradeRequest struct {
	ContractID uuid.UUID `json:"contract_id"`
	Side       string    `json:"side"`     // buy_to_open, sell_to_open, buy_to_close, sell_to_close
	Quantity   int       `json:"quantity"` // number of contracts
}

// ExecuteOptionsTrade executes an options market order.
func (s *OptionsTradeService) ExecuteOptionsTrade(ctx context.Context, userID uuid.UUID, req OptionsTradeRequest) (*model.OptionTrade, error) {
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

	validSides := map[string]bool{
		"buy_to_open": true, "sell_to_open": true,
		"buy_to_close": true, "sell_to_close": true,
	}
	if !validSides[req.Side] {
		return nil, fmt.Errorf("invalid side: must be buy_to_open, sell_to_open, buy_to_close, or sell_to_close")
	}

	// Get the contract
	contract, err := s.repo.GetOptionContract(ctx, req.ContractID)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}
	if contract.Status != "active" {
		return nil, fmt.Errorf("contract is %s, cannot trade", contract.Status)
	}

	// Get user's portfolio
	portfolio, err := s.repo.GetPortfolioByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("portfolio not found: %w", err)
	}

	qty := decimal.NewFromInt(int64(req.Quantity))
	var tradePrice decimal.Decimal
	var total decimal.Decimal

	switch req.Side {
	case "buy_to_open":
		// Buy contracts at ask price
		tradePrice = contract.AskPrice
		total = tradePrice.Mul(qty).Mul(contractMultiplier)

		if portfolio.Cash.LessThan(total) {
			return nil, fmt.Errorf("insufficient funds: need %s, have %s", total, portfolio.Cash)
		}

		// Deduct cash
		newCash := portfolio.Cash.Sub(total)
		if err := s.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash); err != nil {
			return nil, fmt.Errorf("failed to update cash: %w", err)
		}

		// Create or update position (long)
		existing, err := s.repo.GetOptionPosition(ctx, portfolio.ID, contract.ID)
		if err != nil {
			// New position
			if err := s.repo.UpsertOptionPosition(ctx, portfolio.ID, contract.ID, req.Quantity, tradePrice, decimal.Zero); err != nil {
				return nil, fmt.Errorf("failed to create position: %w", err)
			}
		} else {
			// Average cost
			oldTotal := existing.AvgCost.Mul(decimal.NewFromInt(int64(existing.Quantity)).Abs())
			newQty := existing.Quantity + req.Quantity
			newAvg := oldTotal.Add(total.Div(contractMultiplier)).Div(decimal.NewFromInt(int64(newQty)).Abs())
			if err := s.repo.UpsertOptionPosition(ctx, portfolio.ID, contract.ID, newQty, newAvg, existing.Collateral); err != nil {
				return nil, fmt.Errorf("failed to update position: %w", err)
			}
		}

	case "sell_to_open":
		// Write/sell contracts at bid price — requires collateral
		tradePrice = contract.BidPrice
		total = tradePrice.Mul(qty).Mul(contractMultiplier)

		// Calculate required collateral
		collateral, err := s.calculateCollateral(ctx, contract, portfolio, req.Quantity)
		if err != nil {
			return nil, err
		}

		// Check sufficient funds/shares for collateral
		availableCash := portfolio.Cash
		if collateral.GreaterThan(availableCash) {
			return nil, fmt.Errorf("insufficient collateral: need %s, have %s cash available", collateral, availableCash)
		}

		// Reserve collateral (deduct from cash) but credit premium
		newCash := portfolio.Cash.Sub(collateral).Add(total)
		if err := s.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash); err != nil {
			return nil, fmt.Errorf("failed to update cash: %w", err)
		}

		// Create short position (negative quantity)
		existing, err := s.repo.GetOptionPosition(ctx, portfolio.ID, contract.ID)
		if err != nil {
			if err := s.repo.UpsertOptionPosition(ctx, portfolio.ID, contract.ID, -req.Quantity, tradePrice, collateral); err != nil {
				return nil, fmt.Errorf("failed to create position: %w", err)
			}
		} else {
			newQty := existing.Quantity - req.Quantity
			newCollateral := existing.Collateral.Add(collateral)
			if err := s.repo.UpsertOptionPosition(ctx, portfolio.ID, contract.ID, newQty, tradePrice, newCollateral); err != nil {
				return nil, fmt.Errorf("failed to update position: %w", err)
			}
		}

	case "buy_to_close":
		// Close a short position by buying back
		tradePrice = contract.AskPrice
		total = tradePrice.Mul(qty).Mul(contractMultiplier)

		existing, err := s.repo.GetOptionPosition(ctx, portfolio.ID, contract.ID)
		if err != nil {
			return nil, fmt.Errorf("no short position to close")
		}
		if existing.Quantity >= 0 {
			return nil, fmt.Errorf("no short position to close (position is long)")
		}
		absExisting := -existing.Quantity
		if req.Quantity > absExisting {
			return nil, fmt.Errorf("cannot close %d contracts, only %d short", req.Quantity, absExisting)
		}

		if portfolio.Cash.LessThan(total) {
			return nil, fmt.Errorf("insufficient funds: need %s, have %s", total, portfolio.Cash)
		}

		// Release proportional collateral
		releaseFraction := decimal.NewFromInt(int64(req.Quantity)).Div(decimal.NewFromInt(int64(absExisting)))
		released := existing.Collateral.Mul(releaseFraction)

		newCash := portfolio.Cash.Sub(total).Add(released)
		if err := s.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash); err != nil {
			return nil, fmt.Errorf("failed to update cash: %w", err)
		}

		newQty := existing.Quantity + req.Quantity // becomes less negative
		if newQty == 0 {
			if err := s.repo.DeleteOptionPosition(ctx, portfolio.ID, contract.ID); err != nil {
				return nil, fmt.Errorf("failed to delete position: %w", err)
			}
		} else {
			newCollateral := existing.Collateral.Sub(released)
			if err := s.repo.UpsertOptionPosition(ctx, portfolio.ID, contract.ID, newQty, existing.AvgCost, newCollateral); err != nil {
				return nil, fmt.Errorf("failed to update position: %w", err)
			}
		}

	case "sell_to_close":
		// Close a long position by selling
		tradePrice = contract.BidPrice
		total = tradePrice.Mul(qty).Mul(contractMultiplier)

		existing, err := s.repo.GetOptionPosition(ctx, portfolio.ID, contract.ID)
		if err != nil {
			return nil, fmt.Errorf("no long position to close")
		}
		if existing.Quantity <= 0 {
			return nil, fmt.Errorf("no long position to close (position is short)")
		}
		if req.Quantity > existing.Quantity {
			return nil, fmt.Errorf("cannot close %d contracts, only %d held", req.Quantity, existing.Quantity)
		}

		// Credit cash
		newCash := portfolio.Cash.Add(total)
		if err := s.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash); err != nil {
			return nil, fmt.Errorf("failed to update cash: %w", err)
		}

		newQty := existing.Quantity - req.Quantity
		if newQty == 0 {
			if err := s.repo.DeleteOptionPosition(ctx, portfolio.ID, contract.ID); err != nil {
				return nil, fmt.Errorf("failed to delete position: %w", err)
			}
		} else {
			if err := s.repo.UpsertOptionPosition(ctx, portfolio.ID, contract.ID, newQty, existing.AvgCost, decimal.Zero); err != nil {
				return nil, fmt.Errorf("failed to update position: %w", err)
			}
		}
	}

	// Update last traded price
	_ = s.repo.UpdateContractLastPrice(ctx, contract.ID, tradePrice)

	// Record the trade
	trade, err := s.repo.CreateOptionTrade(ctx, userID, contract.ID, req.Side, req.Quantity, tradePrice, total)
	if err != nil {
		return nil, fmt.Errorf("failed to record trade: %w", err)
	}

	return trade, nil
}

// calculateCollateral determines margin requirement for writing options.
func (s *OptionsTradeService) calculateCollateral(ctx context.Context, contract *model.OptionContract, portfolio *model.Portfolio, quantity int) (decimal.Decimal, error) {
	qty := decimal.NewFromInt(int64(quantity))

	if contract.OptionType == "call" {
		// Covered calls: must own underlying shares
		holding, err := s.repo.GetHolding(ctx, portfolio.ID, contract.Ticker)
		if err != nil || holding.Shares < quantity*simulation.ContractMultiplier {
			return decimal.Zero, fmt.Errorf("covered calls require %d shares of %s (you have %d). Naked calls are not allowed in this simulator",
				quantity*simulation.ContractMultiplier, contract.Ticker, 0)
		}
		// Collateral = current value of the shares pledged
		price, ok := s.engine.GetPrice(contract.Ticker)
		if !ok {
			return decimal.Zero, fmt.Errorf("cannot price underlying %s", contract.Ticker)
		}
		return price.Mul(qty).Mul(contractMultiplier), nil

	} else {
		// Cash-secured puts: must have strike * 100 * quantity in cash
		collateral := contract.StrikePrice.Mul(qty).Mul(contractMultiplier)
		return collateral, nil
	}
}
