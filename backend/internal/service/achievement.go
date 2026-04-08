package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/repository"
	ws "github.com/luke/mockstarket/internal/websocket"
	"github.com/shopspring/decimal"
)

// AchievementService evaluates and grants achievements based on user actions.
type AchievementService struct {
	repo   *repository.Repo
	engine market.PriceProvider
	hub    *ws.Hub
	logger *slog.Logger

	startingCash decimal.Decimal
}

// NewAchievementService creates a new achievement evaluator.
func NewAchievementService(repo *repository.Repo, engine market.PriceProvider, hub *ws.Hub, startingCash float64, logger *slog.Logger) *AchievementService {
	return &AchievementService{
		repo:         repo,
		engine:       engine,
		hub:          hub,
		logger:       logger,
		startingCash: decimal.NewFromFloat(startingCash),
	}
}

// OnTradeExecuted should be called after every trade to check trading
// and portfolio achievements.
func (s *AchievementService) OnTradeExecuted(ctx context.Context, userID uuid.UUID) {
	s.checkTradingAchievements(ctx, userID)
	s.checkPortfolioAchievements(ctx, userID)
}

// OnLogin should be called on each user login to check streak achievements.
func (s *AchievementService) OnLogin(ctx context.Context, userID uuid.UUID) {
	s.checkStreakAchievements(ctx, userID)
}

// OnLeaderboardComputed should be called after leaderboard computation
// to check social/ranking achievements.
func (s *AchievementService) OnLeaderboardComputed(ctx context.Context, userID uuid.UUID, rank int) {
	if rank <= 10 {
		s.grant(ctx, userID, "top_ten")
	}
	if rank <= 3 {
		s.grant(ctx, userID, "top_three")
	}
	if rank == 1 {
		s.grant(ctx, userID, "number_one")
	}
}

func (s *AchievementService) checkTradingAchievements(ctx context.Context, userID uuid.UUID) {
	count, err := s.repo.CountTradesByUserID(ctx, userID)
	if err != nil {
		return
	}

	if count >= 1 {
		s.grant(ctx, userID, "first_trade")
	}
	if count >= 10 {
		s.grant(ctx, userID, "ten_trades")
	}
	if count >= 100 {
		s.grant(ctx, userID, "hundred_trades")
	}
}

func (s *AchievementService) checkPortfolioAchievements(ctx context.Context, userID uuid.UUID) {
	portfolio, err := s.repo.GetPortfolioByUserID(ctx, userID)
	if err != nil {
		return
	}

	holdings, err := s.repo.GetHoldingsByPortfolioID(ctx, portfolio.ID)
	if err != nil {
		return
	}

	// Calculate net worth with live prices
	livePrices := s.engine.GetAllPrices()
	investedValue := decimal.Zero
	activeHoldings := 0

	for _, h := range holdings {
		if h.Shares <= 0 {
			continue
		}
		activeHoldings++
		price := h.AvgCost
		if livePrice, ok := livePrices[h.Ticker]; ok {
			price = livePrice
		}
		investedValue = investedValue.Add(price.Mul(decimal.NewFromInt(int64(h.Shares))))
	}

	netWorth := portfolio.Cash.Add(investedValue)

	// First profit
	if netWorth.GreaterThan(s.startingCash) {
		s.grant(ctx, userID, "first_profit")
	}

	// Double up ($200k)
	doubleTarget := s.startingCash.Mul(decimal.NewFromInt(2))
	if netWorth.GreaterThanOrEqual(doubleTarget) {
		s.grant(ctx, userID, "double_up")
	}

	// Millionaire
	million := decimal.NewFromInt(1_000_000)
	if netWorth.GreaterThanOrEqual(million) {
		s.grant(ctx, userID, "millionaire")
	}

	// Diversified (10+ different stocks)
	if activeHoldings >= 10 {
		s.grant(ctx, userID, "diversified")
	}

	// Collector (own every stock) - count total stocks available
	totalStocks, err := s.repo.CountStocks(ctx)
	if err == nil && totalStocks > 0 && activeHoldings >= totalStocks {
		s.grant(ctx, userID, "collector")
	}

	// All In (90%+ in one stock)
	if !investedValue.IsZero() && len(holdings) > 0 {
		for _, h := range holdings {
			if h.Shares <= 0 {
				continue
			}
			price := h.AvgCost
			if livePrice, ok := livePrices[h.Ticker]; ok {
				price = livePrice
			}
			holdingValue := price.Mul(decimal.NewFromInt(int64(h.Shares)))
			ratio := holdingValue.Div(netWorth)
			if ratio.GreaterThanOrEqual(decimal.NewFromFloat(0.90)) {
				s.grant(ctx, userID, "all_in")
				break
			}
		}
	}
}

func (s *AchievementService) checkStreakAchievements(ctx context.Context, userID uuid.UUID) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return
	}

	if user.LoginStreak >= 3 {
		s.grant(ctx, userID, "streak_3")
	}
	if user.LoginStreak >= 7 {
		s.grant(ctx, userID, "streak_7")
	}
	if user.LoginStreak >= 30 {
		s.grant(ctx, userID, "streak_30")
	}
}

// grant idempotently awards an achievement and notifies via WebSocket.
func (s *AchievementService) grant(ctx context.Context, userID uuid.UUID, achievementID string) {
	// GrantAchievement uses ON CONFLICT DO NOTHING, so this is idempotent
	if err := s.repo.GrantAchievement(ctx, userID, achievementID); err != nil {
		s.logger.Error("failed to grant achievement",
			"user_id", userID, "achievement_id", achievementID, "error", err)
		return
	}

	// Only log and notify — the ON CONFLICT means we can't easily distinguish
	// new vs already-earned, but the notification is harmless if duplicated
	// since the client can dedupe on achievement_id.
	s.logger.Debug("achievement granted", "user_id", userID, "achievement_id", achievementID)
}
