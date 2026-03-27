package worker

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/simulation"
	"github.com/shopspring/decimal"
)

// LeaderboardWorker periodically computes leaderboard rankings based on
// portfolio net worth (cash + holdings at live prices).
type LeaderboardWorker struct {
	repo   *repository.Repo
	engine *simulation.Engine
	logger *slog.Logger

	startingCash decimal.Decimal
}

// NewLeaderboardWorker creates a new leaderboard computation worker.
func NewLeaderboardWorker(repo *repository.Repo, engine *simulation.Engine, startingCash float64, logger *slog.Logger) *LeaderboardWorker {
	return &LeaderboardWorker{
		repo:         repo,
		engine:       engine,
		logger:       logger,
		startingCash: decimal.NewFromFloat(startingCash),
	}
}

type rankedUser struct {
	userID      string
	displayName string
	netWorth    decimal.Decimal
	totalReturn decimal.Decimal
}

// Run starts the leaderboard computation loop. Blocks until context is cancelled.
func (w *LeaderboardWorker) Run(ctx context.Context) {
	// Compute immediately on start, then on interval
	w.computeAll(ctx)

	// Alltime: every 5 minutes
	alltimeTicker := time.NewTicker(5 * time.Minute)
	// Daily/weekly: every 15 minutes
	periodicTicker := time.NewTicker(15 * time.Minute)

	defer alltimeTicker.Stop()
	defer periodicTicker.Stop()

	w.logger.Info("leaderboard worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("leaderboard worker stopped")
			return
		case <-alltimeTicker.C:
			w.compute(ctx, "alltime")
		case <-periodicTicker.C:
			w.compute(ctx, "daily")
			w.compute(ctx, "weekly")
		}
	}
}

func (w *LeaderboardWorker) computeAll(ctx context.Context) {
	for _, period := range []string{"alltime", "daily", "weekly"} {
		w.compute(ctx, period)
	}
}

func (w *LeaderboardWorker) compute(ctx context.Context, period string) {
	portfolios, err := w.repo.GetAllPortfolios(ctx)
	if err != nil {
		w.logger.Error("leaderboard: failed to fetch portfolios", "period", period, "error", err)
		return
	}

	if len(portfolios) == 0 {
		return
	}

	livePrices := w.engine.GetAllPrices()
	ranked := make([]rankedUser, 0, len(portfolios))

	for _, p := range portfolios {
		// Get holdings to compute invested value at live prices
		holdings, err := w.repo.GetHoldingsByPortfolioID(ctx, p.ID)
		if err != nil {
			w.logger.Error("leaderboard: failed to fetch holdings",
				"portfolio_id", p.ID, "error", err)
			continue
		}

		investedValue := decimal.Zero
		for _, h := range holdings {
			price := h.AvgCost
			if livePrice, ok := livePrices[h.Ticker]; ok {
				price = livePrice
			}
			investedValue = investedValue.Add(price.Mul(decimal.NewFromInt(int64(h.Shares))))
		}

		netWorth := p.Cash.Add(investedValue)
		totalReturn := decimal.Zero
		if !w.startingCash.IsZero() {
			totalReturn = netWorth.Sub(w.startingCash).Div(w.startingCash).Mul(decimal.NewFromInt(100)).Round(2)
		}

		// Get user display name
		user, err := w.repo.GetUserByID(ctx, p.UserID)
		if err != nil {
			continue
		}

		ranked = append(ranked, rankedUser{
			userID:      p.UserID.String(),
			displayName: user.DisplayName,
			netWorth:    netWorth,
			totalReturn: totalReturn,
		})

		// Update portfolio net_worth while we have the computed value
		_ = w.repo.UpdatePortfolioNetWorth(ctx, p.ID, netWorth)
	}

	// Sort by net worth descending
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].netWorth.GreaterThan(ranked[j].netWorth)
	})

	// Persist rankings
	for i, r := range ranked {
		uid, err := parseUUID(r.userID)
		if err != nil {
			continue
		}
		err = w.repo.InsertLeaderboardEntry(ctx, uid, r.displayName, r.netWorth, r.totalReturn, i+1, period)
		if err != nil {
			w.logger.Error("leaderboard: failed to insert entry",
				"user_id", r.userID, "period", period, "error", err)
		}
	}

	w.logger.Debug("leaderboard computed", "period", period, "entries", len(ranked))
}
