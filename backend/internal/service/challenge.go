package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/luke/mockstarket/internal/model"
	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/shopspring/decimal"
)

// ChallengeService handles daily challenge generation, evaluation, and claiming.
type ChallengeService struct {
	repo   *repository.Repo
	engine market.PriceProvider
	logger *slog.Logger
	rng    *rand.Rand
}

// NewChallengeService creates a new challenge service.
func NewChallengeService(repo *repository.Repo, engine market.PriceProvider, logger *slog.Logger) *ChallengeService {
	return &ChallengeService{
		repo:   repo,
		engine: engine,
		logger: logger,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// challengeTemplate defines a type of challenge that can be generated.
type challengeTemplate struct {
	Type        string
	Description string
	TargetJSON  string
	Reward      decimal.Decimal
}

// GenerateToday creates today's challenge if one doesn't exist yet.
func (s *ChallengeService) GenerateToday(ctx context.Context) (*model.DailyChallenge, error) {
	// Check if today's challenge already exists
	existing, err := s.repo.GetTodaysChallenge(ctx)
	if err == nil && existing.ID != uuid.Nil {
		return existing, nil
	}

	// Pick a random challenge template
	template := s.pickTemplate()

	challenge, err := s.repo.CreateDailyChallenge(ctx,
		template.Type, template.Description, template.TargetJSON, template.Reward)
	if err != nil {
		return nil, fmt.Errorf("failed to create daily challenge: %w", err)
	}

	s.logger.Info("daily challenge generated",
		"type", template.Type,
		"description", template.Description,
	)

	return challenge, nil
}

func (s *ChallengeService) pickTemplate() challengeTemplate {
	// Get a random ticker for stock-specific challenges
	prices := s.engine.GetAllPrices()
	tickers := make([]string, 0, len(prices))
	for t := range prices {
		tickers = append(tickers, t)
	}

	var randomTicker string
	if len(tickers) > 0 {
		randomTicker = tickers[s.rng.Intn(len(tickers))]
	}

	tradeCount := 3 + s.rng.Intn(8)   // 3-10
	profitTarget := 500 + s.rng.Intn(2000) // $500 - $2500

	templates := []challengeTemplate{
		{
			Type:        "trade_count",
			Description: fmt.Sprintf("Execute %d trades today", tradeCount),
			TargetJSON:  fmt.Sprintf(`{"count": %d}`, tradeCount),
			Reward:      decimal.NewFromInt(int64(1000 + tradeCount*100)),
		},
		{
			Type:        "buy_stock",
			Description: fmt.Sprintf("Buy shares of %s", randomTicker),
			TargetJSON:  fmt.Sprintf(`{"ticker": "%s", "side": "buy"}`, randomTicker),
			Reward:      decimal.NewFromInt(500),
		},
		{
			Type:        "sell_stock",
			Description: fmt.Sprintf("Sell shares of %s", randomTicker),
			TargetJSON:  fmt.Sprintf(`{"ticker": "%s", "side": "sell"}`, randomTicker),
			Reward:      decimal.NewFromInt(500),
		},
		{
			Type:        "profit_target",
			Description: fmt.Sprintf("Earn $%d in profit today", profitTarget),
			TargetJSON:  fmt.Sprintf(`{"profit": %d}`, profitTarget),
			Reward:      decimal.NewFromInt(int64(profitTarget)),
		},
		{
			Type:        "diversify",
			Description: "Own shares in at least 5 different stocks",
			TargetJSON:  `{"unique_holdings": 5}`,
			Reward:      decimal.NewFromInt(2000),
		},
		{
			Type:        "volume_trader",
			Description: fmt.Sprintf("Trade at least %d total shares today", 50+s.rng.Intn(100)),
			TargetJSON:  fmt.Sprintf(`{"total_shares": %d}`, 50+s.rng.Intn(100)),
			Reward:      decimal.NewFromInt(1500),
		},
	}

	return templates[s.rng.Intn(len(templates))]
}

// GetTodaysChallenge returns today's challenge and the user's progress on it.
func (s *ChallengeService) GetTodaysChallenge(ctx context.Context, userID uuid.UUID) (*model.DailyChallenge, *model.UserChallenge, error) {
	challenge, err := s.repo.GetTodaysChallenge(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("no challenge today: %w", err)
	}

	// Get or init user's progress
	userChallenge, err := s.repo.GetUserChallenge(ctx, userID, challenge.ID)
	if err != nil {
		// User hasn't started this challenge yet — return nil progress
		return challenge, nil, nil
	}

	return challenge, userChallenge, nil
}

// EvaluateForUser checks if the user has completed today's challenge.
func (s *ChallengeService) EvaluateForUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	challenge, err := s.repo.GetTodaysChallenge(ctx)
	if err != nil {
		return false, nil // No challenge today
	}

	var target map[string]interface{}
	if err := json.Unmarshal([]byte(challenge.TargetJSON), &target); err != nil {
		return false, fmt.Errorf("invalid target JSON: %w", err)
	}

	completed := false

	switch challenge.ChallengeType {
	case "trade_count":
		count, err := s.repo.CountTodaysTradesByUserID(ctx, userID)
		if err != nil {
			return false, err
		}
		if targetCount, ok := target["count"].(float64); ok {
			completed = count >= int(targetCount)
		}

	case "buy_stock", "sell_stock":
		ticker, _ := target["ticker"].(string)
		side, _ := target["side"].(string)
		if ticker != "" && side != "" {
			hasTrade, err := s.repo.HasTodaysTrade(ctx, userID, ticker, side)
			if err != nil {
				return false, err
			}
			completed = hasTrade
		}

	case "diversify":
		portfolio, err := s.repo.GetPortfolioByUserID(ctx, userID)
		if err != nil {
			return false, nil
		}
		holdings, err := s.repo.GetHoldingsByPortfolioID(ctx, portfolio.ID)
		if err != nil {
			return false, err
		}
		activeCount := 0
		for _, h := range holdings {
			if h.Shares > 0 {
				activeCount++
			}
		}
		if targetHoldings, ok := target["unique_holdings"].(float64); ok {
			completed = activeCount >= int(targetHoldings)
		}

	case "volume_trader":
		volume, err := s.repo.SumTodaysSharesByUserID(ctx, userID)
		if err != nil {
			return false, err
		}
		if targetShares, ok := target["total_shares"].(float64); ok {
			completed = volume >= int(targetShares)
		}

	case "profit_target":
		// Simplified: check if user's net worth increased by the target today
		// This is approximate — a full implementation would track intraday P&L
		completed = false // Too complex to evaluate without intraday snapshots
	}

	if completed {
		_ = s.repo.UpsertUserChallenge(ctx, userID, challenge.ID, true)
	}

	return completed, nil
}

// ClaimReward lets a user claim their reward for a completed challenge.
func (s *ChallengeService) ClaimReward(ctx context.Context, userID uuid.UUID, challengeID uuid.UUID) error {
	challenge, err := s.repo.GetTodaysChallenge(ctx)
	if err != nil {
		return fmt.Errorf("challenge not found")
	}
	if challenge.ID != challengeID {
		return fmt.Errorf("can only claim today's challenge")
	}

	uc, err := s.repo.GetUserChallenge(ctx, userID, challengeID)
	if err != nil {
		return fmt.Errorf("challenge not started")
	}
	if !uc.Completed {
		return fmt.Errorf("challenge not completed")
	}
	if uc.Claimed {
		return fmt.Errorf("reward already claimed")
	}

	// Grant reward cash
	portfolio, err := s.repo.GetPortfolioByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("portfolio not found")
	}

	newCash := portfolio.Cash.Add(challenge.RewardCash)
	if err := s.repo.UpdatePortfolioCash(ctx, portfolio.ID, newCash); err != nil {
		return fmt.Errorf("failed to grant reward: %w", err)
	}

	if err := s.repo.ClaimChallengeReward(ctx, userID, challengeID); err != nil {
		return fmt.Errorf("failed to mark claimed: %w", err)
	}

	s.logger.Info("challenge reward claimed",
		"user_id", userID,
		"challenge_id", challengeID,
		"reward", challenge.RewardCash,
	)

	return nil
}
