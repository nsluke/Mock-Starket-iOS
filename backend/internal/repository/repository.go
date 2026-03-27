package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luke/mockstarket/internal/model"
	"github.com/shopspring/decimal"
)

// Repo provides database access methods.
type Repo struct {
	pool *pgxpool.Pool
}

// New creates a new repository backed by a connection pool.
func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// ---- Users ----

func (r *Repo) CreateUser(ctx context.Context, firebaseUID, displayName string, isGuest bool) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (firebase_uid, display_name, is_guest)
		 VALUES ($1, $2, $3)
		 RETURNING id, firebase_uid, display_name, avatar_url, is_guest, created_at, updated_at, last_login_at, login_streak, longest_streak`,
		firebaseUID, displayName, isGuest,
	).Scan(&u.ID, &u.FirebaseUID, &u.DisplayName, &u.AvatarURL, &u.IsGuest,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt, &u.LoginStreak, &u.LongestStreak)
	return u, err
}

func (r *Repo) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, firebase_uid, display_name, avatar_url, is_guest, created_at, updated_at, last_login_at, login_streak, longest_streak
		 FROM users WHERE firebase_uid = $1`, firebaseUID,
	).Scan(&u.ID, &u.FirebaseUID, &u.DisplayName, &u.AvatarURL, &u.IsGuest,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt, &u.LoginStreak, &u.LongestStreak)
	return u, err
}

func (r *Repo) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, firebase_uid, display_name, avatar_url, is_guest, created_at, updated_at, last_login_at, login_streak, longest_streak
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.FirebaseUID, &u.DisplayName, &u.AvatarURL, &u.IsGuest,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt, &u.LoginStreak, &u.LongestStreak)
	return u, err
}

func (r *Repo) UpdateUser(ctx context.Context, id uuid.UUID, displayName string, avatarURL *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET display_name = $2, avatar_url = $3, updated_at = NOW() WHERE id = $1`,
		id, displayName, avatarURL)
	return err
}

func (r *Repo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *Repo) UpdateLoginStreak(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET
			last_login_at = NOW(),
			login_streak = CASE
				WHEN last_login_at IS NULL THEN 1
				WHEN last_login_at::date = (NOW() - INTERVAL '1 day')::date THEN login_streak + 1
				WHEN last_login_at::date = NOW()::date THEN login_streak
				ELSE 1
			END,
			longest_streak = GREATEST(longest_streak, login_streak),
			updated_at = NOW()
		WHERE id = $1`, id)
	return err
}

// ---- Stocks ----

func (r *Repo) UpsertStock(ctx context.Context, s *model.Stock) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO stocks (ticker, name, sector, base_price, current_price, day_open, day_high, day_low, prev_close, volume, volatility, drift, mean_reversion, description)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		 ON CONFLICT (ticker) DO UPDATE SET
			name = EXCLUDED.name, sector = EXCLUDED.sector, base_price = EXCLUDED.base_price,
			current_price = EXCLUDED.current_price, volatility = EXCLUDED.volatility, drift = EXCLUDED.drift,
			mean_reversion = EXCLUDED.mean_reversion, description = EXCLUDED.description`,
		s.Ticker, s.Name, s.Sector, s.BasePrice, s.CurrentPrice, s.DayOpen, s.DayHigh, s.DayLow,
		s.PrevClose, s.Volume, s.Volatility, s.Drift, s.MeanReversion, s.Description)
	return err
}

func (r *Repo) GetAllStocks(ctx context.Context) ([]model.Stock, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ticker, name, sector, base_price, current_price, day_open, day_high, day_low, prev_close, volume, volatility, drift, mean_reversion, description, logo_url, created_at
		 FROM stocks ORDER BY ticker`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []model.Stock
	for rows.Next() {
		var s model.Stock
		if err := rows.Scan(&s.Ticker, &s.Name, &s.Sector, &s.BasePrice, &s.CurrentPrice,
			&s.DayOpen, &s.DayHigh, &s.DayLow, &s.PrevClose, &s.Volume,
			&s.Volatility, &s.Drift, &s.MeanReversion, &s.Description, &s.LogoURL, &s.CreatedAt); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func (r *Repo) GetStockByTicker(ctx context.Context, ticker string) (*model.Stock, error) {
	s := &model.Stock{}
	err := r.pool.QueryRow(ctx,
		`SELECT ticker, name, sector, base_price, current_price, day_open, day_high, day_low, prev_close, volume, volatility, drift, mean_reversion, description, logo_url, created_at
		 FROM stocks WHERE ticker = $1`, ticker,
	).Scan(&s.Ticker, &s.Name, &s.Sector, &s.BasePrice, &s.CurrentPrice,
		&s.DayOpen, &s.DayHigh, &s.DayLow, &s.PrevClose, &s.Volume,
		&s.Volatility, &s.Drift, &s.MeanReversion, &s.Description, &s.LogoURL, &s.CreatedAt)
	return s, err
}

func (r *Repo) UpdateStockPrices(ctx context.Context, ticker string, price, high, low decimal.Decimal, volume int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE stocks SET current_price = $2, day_high = $3, day_low = $4, volume = $5 WHERE ticker = $1`,
		ticker, price, high, low, volume)
	return err
}

// ---- Portfolios ----

func (r *Repo) CreatePortfolio(ctx context.Context, userID uuid.UUID, startingCash decimal.Decimal) (*model.Portfolio, error) {
	p := &model.Portfolio{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO portfolios (user_id, cash, net_worth) VALUES ($1, $2, $2)
		 RETURNING id, user_id, cash, net_worth, created_at, updated_at`,
		userID, startingCash,
	).Scan(&p.ID, &p.UserID, &p.Cash, &p.NetWorth, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repo) GetPortfolioByUserID(ctx context.Context, userID uuid.UUID) (*model.Portfolio, error) {
	p := &model.Portfolio{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, cash, net_worth, created_at, updated_at
		 FROM portfolios WHERE user_id = $1`, userID,
	).Scan(&p.ID, &p.UserID, &p.Cash, &p.NetWorth, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repo) UpdatePortfolioCash(ctx context.Context, portfolioID uuid.UUID, cash decimal.Decimal) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE portfolios SET cash = $2, updated_at = NOW() WHERE id = $1`,
		portfolioID, cash)
	return err
}

func (r *Repo) UpdatePortfolioNetWorth(ctx context.Context, portfolioID uuid.UUID, netWorth decimal.Decimal) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE portfolios SET net_worth = $2, updated_at = NOW() WHERE id = $1`,
		portfolioID, netWorth)
	return err
}

// ---- Holdings ----

func (r *Repo) GetHoldingsByPortfolioID(ctx context.Context, portfolioID uuid.UUID) ([]model.Holding, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, portfolio_id, ticker, shares, avg_cost, created_at, updated_at
		 FROM holdings WHERE portfolio_id = $1 AND shares > 0 ORDER BY ticker`, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holdings []model.Holding
	for rows.Next() {
		var h model.Holding
		if err := rows.Scan(&h.ID, &h.PortfolioID, &h.Ticker, &h.Shares, &h.AvgCost, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		holdings = append(holdings, h)
	}
	return holdings, nil
}

func (r *Repo) UpsertHolding(ctx context.Context, portfolioID uuid.UUID, ticker string, shares int, avgCost decimal.Decimal) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO holdings (portfolio_id, ticker, shares, avg_cost)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (portfolio_id, ticker) DO UPDATE SET
			shares = EXCLUDED.shares, avg_cost = EXCLUDED.avg_cost, updated_at = NOW()`,
		portfolioID, ticker, shares, avgCost)
	return err
}

func (r *Repo) GetHolding(ctx context.Context, portfolioID uuid.UUID, ticker string) (*model.Holding, error) {
	h := &model.Holding{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, portfolio_id, ticker, shares, avg_cost, created_at, updated_at
		 FROM holdings WHERE portfolio_id = $1 AND ticker = $2`, portfolioID, ticker,
	).Scan(&h.ID, &h.PortfolioID, &h.Ticker, &h.Shares, &h.AvgCost, &h.CreatedAt, &h.UpdatedAt)
	return h, err
}

// ---- Trades ----

func (r *Repo) CreateTrade(ctx context.Context, userID uuid.UUID, ticker, side string, shares int, price, total decimal.Decimal, orderID *uuid.UUID) (*model.Trade, error) {
	t := &model.Trade{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO trades (user_id, ticker, side, shares, price, total, order_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, user_id, ticker, side, shares, price, total, order_id, created_at`,
		userID, ticker, side, shares, price, total, orderID,
	).Scan(&t.ID, &t.UserID, &t.Ticker, &t.Side, &t.Shares, &t.Price, &t.Total, &t.OrderID, &t.CreatedAt)
	return t, err
}

func (r *Repo) GetTradesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Trade, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, ticker, side, shares, price, total, order_id, created_at
		 FROM trades WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []model.Trade
	for rows.Next() {
		var t model.Trade
		if err := rows.Scan(&t.ID, &t.UserID, &t.Ticker, &t.Side, &t.Shares, &t.Price, &t.Total, &t.OrderID, &t.CreatedAt); err != nil {
			return nil, err
		}
		trades = append(trades, t)
	}
	return trades, nil
}

// ---- Orders ----

func (r *Repo) CreateOrder(ctx context.Context, userID uuid.UUID, ticker, side, orderType string, shares int, limitPrice, stopPrice *decimal.Decimal) (*model.Order, error) {
	o := &model.Order{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO orders (user_id, ticker, side, order_type, shares, limit_price, stop_price)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, user_id, ticker, side, order_type, shares, limit_price, stop_price, status, filled_price, filled_at, expires_at, created_at, updated_at`,
		userID, ticker, side, orderType, shares, limitPrice, stopPrice,
	).Scan(&o.ID, &o.UserID, &o.Ticker, &o.Side, &o.OrderType, &o.Shares,
		&o.LimitPrice, &o.StopPrice, &o.Status, &o.FilledPrice, &o.FilledAt, &o.ExpiresAt,
		&o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *Repo) GetOpenOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, ticker, side, order_type, shares, limit_price, stop_price, status, filled_price, filled_at, expires_at, created_at, updated_at
		 FROM orders WHERE user_id = $1 AND status = 'open' ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Ticker, &o.Side, &o.OrderType, &o.Shares,
			&o.LimitPrice, &o.StopPrice, &o.Status, &o.FilledPrice, &o.FilledAt, &o.ExpiresAt,
			&o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *Repo) GetAllOpenOrders(ctx context.Context) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, ticker, side, order_type, shares, limit_price, stop_price, status, filled_price, filled_at, expires_at, created_at, updated_at
		 FROM orders WHERE status = 'open'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Ticker, &o.Side, &o.OrderType, &o.Shares,
			&o.LimitPrice, &o.StopPrice, &o.Status, &o.FilledPrice, &o.FilledAt, &o.ExpiresAt,
			&o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *Repo) FillOrder(ctx context.Context, orderID uuid.UUID, filledPrice decimal.Decimal) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE orders SET status = 'filled', filled_price = $2, filled_at = NOW(), updated_at = NOW() WHERE id = $1`,
		orderID, filledPrice)
	return err
}

func (r *Repo) CancelOrder(ctx context.Context, orderID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE orders SET status = 'cancelled', updated_at = NOW() WHERE id = $1 AND user_id = $2 AND status = 'open'`,
		orderID, userID)
	return err
}

// ---- Price History ----

func (r *Repo) InsertPriceHistory(ctx context.Context, ticker string, price, openP, high, low, closeP decimal.Decimal, volume int64, interval string, recordedAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO price_history (ticker, price, open_price, high, low, close_price, volume, interval, recorded_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (ticker, interval, recorded_at) DO NOTHING`,
		ticker, price, openP, high, low, closeP, volume, interval, recordedAt)
	return err
}

func (r *Repo) GetPriceHistory(ctx context.Context, ticker, interval string, from, to time.Time, limit int) ([]model.PricePoint, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, ticker, price, open_price, high, low, close_price, volume, interval, recorded_at
		 FROM price_history
		 WHERE ticker = $1 AND interval = $2 AND recorded_at >= $3 AND recorded_at <= $4
		 ORDER BY recorded_at DESC LIMIT $5`,
		ticker, interval, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []model.PricePoint
	for rows.Next() {
		var p model.PricePoint
		if err := rows.Scan(&p.ID, &p.Ticker, &p.Price, &p.Open, &p.High, &p.Low, &p.Close, &p.Volume, &p.Interval, &p.RecordedAt); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}

// ---- Leaderboard ----

func (r *Repo) GetLeaderboard(ctx context.Context, period string, limit, offset int) ([]model.LeaderboardEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, display_name, net_worth, total_return, rank, period, computed_at
		 FROM (
			SELECT DISTINCT ON (user_id) id, user_id, display_name, net_worth, total_return, rank, period, computed_at
			FROM leaderboard
			WHERE period = $1
			ORDER BY user_id, computed_at DESC
		 ) AS latest
		 ORDER BY rank ASC
		 LIMIT $2 OFFSET $3`,
		period, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []model.LeaderboardEntry
	for rows.Next() {
		var e model.LeaderboardEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.DisplayName, &e.NetWorth, &e.TotalReturn, &e.Rank, &e.Period, &e.ComputedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *Repo) InsertLeaderboardEntry(ctx context.Context, userID uuid.UUID, displayName string, netWorth, totalReturn decimal.Decimal, rank int, period string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO leaderboard (user_id, display_name, net_worth, total_return, rank, period)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, displayName, netWorth, totalReturn, rank, period)
	return err
}

// ---- Price Alerts ----

func (r *Repo) CreatePriceAlert(ctx context.Context, userID uuid.UUID, ticker, condition string, targetPrice decimal.Decimal) (*model.PriceAlert, error) {
	a := &model.PriceAlert{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO price_alerts (user_id, ticker, condition, target_price)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, ticker, condition, target_price, triggered, triggered_at, created_at`,
		userID, ticker, condition, targetPrice,
	).Scan(&a.ID, &a.UserID, &a.Ticker, &a.Condition, &a.TargetPrice, &a.Triggered, &a.TriggeredAt, &a.CreatedAt)
	return a, err
}

func (r *Repo) GetAlertsByUserID(ctx context.Context, userID uuid.UUID) ([]model.PriceAlert, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, ticker, condition, target_price, triggered, triggered_at, created_at
		 FROM price_alerts WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []model.PriceAlert
	for rows.Next() {
		var a model.PriceAlert
		if err := rows.Scan(&a.ID, &a.UserID, &a.Ticker, &a.Condition, &a.TargetPrice, &a.Triggered, &a.TriggeredAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *Repo) GetUntriggeredAlerts(ctx context.Context) ([]model.PriceAlert, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, ticker, condition, target_price, triggered, triggered_at, created_at
		 FROM price_alerts WHERE triggered = FALSE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []model.PriceAlert
	for rows.Next() {
		var a model.PriceAlert
		if err := rows.Scan(&a.ID, &a.UserID, &a.Ticker, &a.Condition, &a.TargetPrice, &a.Triggered, &a.TriggeredAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *Repo) TriggerAlert(ctx context.Context, alertID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE price_alerts SET triggered = TRUE, triggered_at = NOW() WHERE id = $1`,
		alertID)
	return err
}

func (r *Repo) DeleteAlert(ctx context.Context, alertID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM price_alerts WHERE id = $1 AND user_id = $2`,
		alertID, userID)
	return err
}

// ---- Achievements ----

func (r *Repo) GetAllAchievements(ctx context.Context) ([]model.Achievement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, icon, category, criteria_json FROM achievements ORDER BY category, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []model.Achievement
	for rows.Next() {
		var a model.Achievement
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Icon, &a.Category, &a.CriteriaJSON); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

func (r *Repo) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]model.UserAchievement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, achievement_id, earned_at FROM user_achievements WHERE user_id = $1 ORDER BY earned_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []model.UserAchievement
	for rows.Next() {
		var a model.UserAchievement
		if err := rows.Scan(&a.ID, &a.UserID, &a.AchievementID, &a.EarnedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

func (r *Repo) GrantAchievement(ctx context.Context, userID uuid.UUID, achievementID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_achievements (user_id, achievement_id)
		 VALUES ($1, $2) ON CONFLICT (user_id, achievement_id) DO NOTHING`,
		userID, achievementID)
	return err
}

// ---- Watchlist ----

func (r *Repo) AddToWatchlist(ctx context.Context, userID uuid.UUID, ticker string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO watchlist (user_id, ticker) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, ticker)
	return err
}

func (r *Repo) RemoveFromWatchlist(ctx context.Context, userID uuid.UUID, ticker string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM watchlist WHERE user_id = $1 AND ticker = $2`, userID, ticker)
	return err
}

func (r *Repo) GetWatchlist(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT ticker FROM watchlist WHERE user_id = $1 ORDER BY created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickers []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		tickers = append(tickers, t)
	}
	return tickers, nil
}

// ---- Portfolio History ----

func (r *Repo) InsertPortfolioHistory(ctx context.Context, userID uuid.UUID, netWorth, cash decimal.Decimal) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO portfolio_history (user_id, net_worth, cash) VALUES ($1, $2, $3)`,
		userID, netWorth, cash)
	return err
}

func (r *Repo) GetPortfolioHistory(ctx context.Context, userID uuid.UUID, limit int) ([]model.PortfolioHistory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, net_worth, cash, recorded_at
		 FROM portfolio_history WHERE user_id = $1 ORDER BY recorded_at DESC LIMIT $2`,
		userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.PortfolioHistory
	for rows.Next() {
		var h model.PortfolioHistory
		if err := rows.Scan(&h.ID, &h.UserID, &h.NetWorth, &h.Cash, &h.RecordedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

// ---- Daily Challenges ----

func (r *Repo) GetTodaysChallenge(ctx context.Context) (*model.DailyChallenge, error) {
	c := &model.DailyChallenge{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, date, challenge_type, description, target_json, reward_cash, created_at
		 FROM daily_challenges WHERE date = CURRENT_DATE`,
	).Scan(&c.ID, &c.Date, &c.ChallengeType, &c.Description, &c.TargetJSON, &c.RewardCash, &c.CreatedAt)
	return c, err
}

func (r *Repo) CreateDailyChallenge(ctx context.Context, challengeType, description, targetJSON string, rewardCash decimal.Decimal) (*model.DailyChallenge, error) {
	c := &model.DailyChallenge{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO daily_challenges (date, challenge_type, description, target_json, reward_cash)
		 VALUES (CURRENT_DATE, $1, $2, $3::jsonb, $4)
		 ON CONFLICT (date) DO NOTHING
		 RETURNING id, date, challenge_type, description, target_json, reward_cash, created_at`,
		challengeType, description, targetJSON, rewardCash,
	).Scan(&c.ID, &c.Date, &c.ChallengeType, &c.Description, &c.TargetJSON, &c.RewardCash, &c.CreatedAt)
	return c, err
}

func (r *Repo) GetUserChallenge(ctx context.Context, userID, challengeID uuid.UUID) (*model.UserChallenge, error) {
	uc := &model.UserChallenge{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, challenge_id, completed, completed_at, claimed
		 FROM user_challenges WHERE user_id = $1 AND challenge_id = $2`,
		userID, challengeID,
	).Scan(&uc.ID, &uc.UserID, &uc.ChallengeID, &uc.Completed, &uc.CompletedAt, &uc.Claimed)
	return uc, err
}

func (r *Repo) UpsertUserChallenge(ctx context.Context, userID, challengeID uuid.UUID, completed bool) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_challenges (user_id, challenge_id, completed, completed_at)
		 VALUES ($1, $2, $3, CASE WHEN $3 THEN NOW() ELSE NULL END)
		 ON CONFLICT (user_id, challenge_id) DO UPDATE SET
			completed = EXCLUDED.completed,
			completed_at = CASE WHEN EXCLUDED.completed AND NOT user_challenges.completed THEN NOW() ELSE user_challenges.completed_at END`,
		userID, challengeID, completed)
	return err
}

func (r *Repo) ClaimChallengeReward(ctx context.Context, userID, challengeID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE user_challenges SET claimed = TRUE
		 WHERE user_id = $1 AND challenge_id = $2 AND completed = TRUE AND claimed = FALSE`,
		userID, challengeID)
	return err
}

// ---- Counts & Aggregations ----

func (r *Repo) CountTradesByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM trades WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *Repo) CountStocks(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM stocks`).Scan(&count)
	return count, err
}

func (r *Repo) CountTodaysTradesByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM trades WHERE user_id = $1 AND created_at::date = CURRENT_DATE`,
		userID).Scan(&count)
	return count, err
}

func (r *Repo) HasTodaysTrade(ctx context.Context, userID uuid.UUID, ticker, side string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM trades WHERE user_id = $1 AND ticker = $2 AND side = $3 AND created_at::date = CURRENT_DATE)`,
		userID, ticker, side).Scan(&exists)
	return exists, err
}

func (r *Repo) SumTodaysSharesByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(shares), 0) FROM trades WHERE user_id = $1 AND created_at::date = CURRENT_DATE`,
		userID).Scan(&total)
	return total, err
}

// ---- Helpers ----

// GetAllPortfolios returns all portfolios (for leaderboard computation).
func (r *Repo) GetAllPortfolios(ctx context.Context) ([]model.Portfolio, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, cash, net_worth, created_at, updated_at FROM portfolios`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []model.Portfolio
	for rows.Next() {
		var p model.Portfolio
		if err := rows.Scan(&p.ID, &p.UserID, &p.Cash, &p.NetWorth, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		portfolios = append(portfolios, p)
	}
	return portfolios, nil
}
