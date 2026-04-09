package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/luke/mockstarket/internal/model"
	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/shopspring/decimal"
)

// TradeService handles trade execution and validation.
type TradeService struct {
	repo   *repository.Repo
	engine market.PriceProvider

	// Optional achievement callback, set after construction to avoid circular deps
	onTradeExecuted func(ctx context.Context, userID uuid.UUID)
}

// NewTradeService creates a new trade service.
func NewTradeService(repo *repository.Repo, engine market.PriceProvider) *TradeService {
	return &TradeService{repo: repo, engine: engine}
}

// SetOnTradeExecuted registers a callback invoked after each successful trade.
func (s *TradeService) SetOnTradeExecuted(fn func(ctx context.Context, userID uuid.UUID)) {
	s.onTradeExecuted = fn
}

// TradeRequest represents a buy/sell request.
type TradeRequest struct {
	Ticker string `json:"ticker"`
	Side   string `json:"side"`
	Shares int    `json:"shares"`
}

// ExecuteTrade executes a market order.
func (s *TradeService) ExecuteTrade(ctx context.Context, userID uuid.UUID, req TradeRequest) (*model.Trade, error) {
	if req.Shares <= 0 {
		return nil, fmt.Errorf("shares must be positive")
	}
	if req.Side != "buy" && req.Side != "sell" {
		return nil, fmt.Errorf("side must be 'buy' or 'sell'")
	}

	// Get current price from price provider, fall back to DB
	price, ok := s.engine.GetPrice(req.Ticker)
	if !ok {
		stock, err := s.repo.GetStockByTicker(ctx, req.Ticker)
		if err != nil {
			return nil, fmt.Errorf("stock %s not found", req.Ticker)
		}
		price = stock.CurrentPrice
	}

	total := price.Mul(decimal.NewFromInt(int64(req.Shares)))

	// Get user's portfolio
	portfolio, err := s.repo.GetPortfolioByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("portfolio not found: %w", err)
	}

	if req.Side == "buy" {
		if portfolio.Cash.LessThan(total) {
			return nil, fmt.Errorf("insufficient funds: need %s, have %s", total, portfolio.Cash)
		}

		// Deduct cash
		newCash := portfolio.Cash.Sub(total)
		if err := s.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash); err != nil {
			return nil, fmt.Errorf("failed to update cash: %w", err)
		}

		// Update holding
		holding, err := s.repo.GetHolding(ctx, portfolio.ID, req.Ticker)
		if err != nil {
			// New position
			if err := s.repo.UpsertHolding(ctx, portfolio.ID, req.Ticker, req.Shares, price); err != nil {
				return nil, fmt.Errorf("failed to create holding: %w", err)
			}
		} else {
			// Average up/down
			totalCost := holding.AvgCost.Mul(decimal.NewFromInt(int64(holding.Shares))).Add(total)
			newShares := holding.Shares + req.Shares
			newAvgCost := totalCost.Div(decimal.NewFromInt(int64(newShares)))
			if err := s.repo.UpsertHolding(ctx, portfolio.ID, req.Ticker, newShares, newAvgCost); err != nil {
				return nil, fmt.Errorf("failed to update holding: %w", err)
			}
		}
	} else {
		// Sell
		holding, err := s.repo.GetHolding(ctx, portfolio.ID, req.Ticker)
		if err != nil {
			return nil, fmt.Errorf("no position in %s", req.Ticker)
		}
		if holding.Shares < req.Shares {
			return nil, fmt.Errorf("insufficient shares: have %d, selling %d", holding.Shares, req.Shares)
		}

		// Add cash
		newCash := portfolio.Cash.Add(total)
		if err := s.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash); err != nil {
			return nil, fmt.Errorf("failed to update cash: %w", err)
		}

		// Update holding
		newShares := holding.Shares - req.Shares
		if err := s.repo.UpsertHolding(ctx, portfolio.ID, req.Ticker, newShares, holding.AvgCost); err != nil {
			return nil, fmt.Errorf("failed to update holding: %w", err)
		}
	}

	// Record the trade
	trade, err := s.repo.CreateTrade(ctx, userID, req.Ticker, req.Side, req.Shares, price, total, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to record trade: %w", err)
	}

	// Update portfolio net worth and record history
	s.updateNetWorth(ctx, userID, portfolio.ID)

	// Notify achievement evaluator
	if s.onTradeExecuted != nil {
		s.onTradeExecuted(ctx, userID)
	}

	return trade, nil
}

// updateNetWorth recalculates and persists net worth + portfolio history snapshot.
func (s *TradeService) updateNetWorth(ctx context.Context, userID uuid.UUID, portfolioID uuid.UUID) {
	portfolio, err := s.repo.GetPortfolioByUserID(ctx, userID)
	if err != nil {
		return
	}

	holdings, err := s.repo.GetHoldingsByPortfolioID(ctx, portfolioID)
	if err != nil {
		return
	}

	livePrices := s.engine.GetAllPrices()
	investedValue := decimal.Zero
	for _, h := range holdings {
		price := h.AvgCost
		if livePrice, ok := livePrices[h.Ticker]; ok {
			price = livePrice
		}
		investedValue = investedValue.Add(price.Mul(decimal.NewFromInt(int64(h.Shares))))
	}

	netWorth := portfolio.Cash.Add(investedValue)
	_ = s.repo.UpdatePortfolioNetWorth(ctx, portfolioID, netWorth)
	_ = s.repo.InsertPortfolioHistory(ctx, userID, netWorth, portfolio.Cash)
}

// GetTradeHistory returns paginated trade history for a user.
func (s *TradeService) GetTradeHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Trade, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.GetTradesByUserID(ctx, userID, limit, offset)
}
