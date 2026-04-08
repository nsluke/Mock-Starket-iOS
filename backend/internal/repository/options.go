package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/luke/mockstarket/internal/model"
	"github.com/shopspring/decimal"
)

// ---- Option Contracts ----

func (r *Repo) UpsertOptionContract(ctx context.Context, c *model.OptionContract) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO option_contracts (ticker, option_type, strike_price, expiration, contract_symbol, bid_price, ask_price, mark_price, implied_vol, delta, gamma, theta, vega, rho, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		 ON CONFLICT (contract_symbol) DO UPDATE SET
			bid_price = EXCLUDED.bid_price, ask_price = EXCLUDED.ask_price, mark_price = EXCLUDED.mark_price,
			implied_vol = EXCLUDED.implied_vol, delta = EXCLUDED.delta, gamma = EXCLUDED.gamma,
			theta = EXCLUDED.theta, vega = EXCLUDED.vega, rho = EXCLUDED.rho, updated_at = NOW()`,
		c.Ticker, c.OptionType, c.StrikePrice, c.Expiration, c.ContractSymbol,
		c.BidPrice, c.AskPrice, c.MarkPrice, c.ImpliedVol,
		c.Delta, c.Gamma, c.Theta, c.Vega, c.Rho, c.Status)
	return err
}

func (r *Repo) GetOptionChain(ctx context.Context, ticker string, expiration *time.Time) ([]model.OptionContract, error) {
	var query string
	var args []interface{}

	if expiration != nil {
		query = `SELECT id, ticker, option_type, strike_price, expiration, contract_symbol,
				bid_price, ask_price, last_price, mark_price, open_interest, volume, implied_vol,
				delta, gamma, theta, vega, rho, status, created_at, updated_at
			 FROM option_contracts
			 WHERE ticker = $1 AND status = 'active' AND expiration = $2
			 ORDER BY strike_price, option_type`
		args = []interface{}{ticker, *expiration}
	} else {
		query = `SELECT id, ticker, option_type, strike_price, expiration, contract_symbol,
				bid_price, ask_price, last_price, mark_price, open_interest, volume, implied_vol,
				delta, gamma, theta, vega, rho, status, created_at, updated_at
			 FROM option_contracts
			 WHERE ticker = $1 AND status = 'active'
			 ORDER BY expiration, strike_price, option_type`
		args = []interface{}{ticker}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptionContracts(rows)
}

func (r *Repo) GetOptionContract(ctx context.Context, contractID uuid.UUID) (*model.OptionContract, error) {
	c := &model.OptionContract{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, ticker, option_type, strike_price, expiration, contract_symbol,
			bid_price, ask_price, last_price, mark_price, open_interest, volume, implied_vol,
			delta, gamma, theta, vega, rho, status, created_at, updated_at
		 FROM option_contracts WHERE id = $1`, contractID,
	).Scan(&c.ID, &c.Ticker, &c.OptionType, &c.StrikePrice, &c.Expiration, &c.ContractSymbol,
		&c.BidPrice, &c.AskPrice, &c.LastPrice, &c.MarkPrice, &c.OpenInterest, &c.Volume, &c.ImpliedVol,
		&c.Delta, &c.Gamma, &c.Theta, &c.Vega, &c.Rho, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repo) GetOptionExpirations(ctx context.Context, ticker string) ([]time.Time, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT expiration FROM option_contracts
		 WHERE ticker = $1 AND status = 'active' ORDER BY expiration`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expirations []time.Time
	for rows.Next() {
		var t time.Time
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		expirations = append(expirations, t)
	}
	return expirations, nil
}

func (r *Repo) GetActiveContractsByTicker(ctx context.Context, ticker string) ([]model.OptionContract, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, ticker, option_type, strike_price, expiration, contract_symbol,
			bid_price, ask_price, last_price, mark_price, open_interest, volume, implied_vol,
			delta, gamma, theta, vega, rho, status, created_at, updated_at
		 FROM option_contracts WHERE ticker = $1 AND status = 'active'`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptionContracts(rows)
}

func (r *Repo) GetAllActiveContracts(ctx context.Context) ([]model.OptionContract, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, ticker, option_type, strike_price, expiration, contract_symbol,
			bid_price, ask_price, last_price, mark_price, open_interest, volume, implied_vol,
			delta, gamma, theta, vega, rho, status, created_at, updated_at
		 FROM option_contracts WHERE status = 'active'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptionContracts(rows)
}

func (r *Repo) GetExpiringContracts(ctx context.Context, before time.Time) ([]model.OptionContract, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, ticker, option_type, strike_price, expiration, contract_symbol,
			bid_price, ask_price, last_price, mark_price, open_interest, volume, implied_vol,
			delta, gamma, theta, vega, rho, status, created_at, updated_at
		 FROM option_contracts WHERE expiration <= $1 AND status = 'active'`, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptionContracts(rows)
}

func (r *Repo) UpdateContractStatus(ctx context.Context, contractID uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE option_contracts SET status = $2, updated_at = NOW() WHERE id = $1`,
		contractID, status)
	return err
}

func (r *Repo) BatchUpdateOptionPrices(ctx context.Context, contracts []model.OptionContract) error {
	for _, c := range contracts {
		_, err := r.pool.Exec(ctx,
			`UPDATE option_contracts SET
				bid_price = $2, ask_price = $3, mark_price = $4, implied_vol = $5,
				delta = $6, gamma = $7, theta = $8, vega = $9, rho = $10, updated_at = NOW()
			 WHERE id = $1`,
			c.ID, c.BidPrice, c.AskPrice, c.MarkPrice, c.ImpliedVol,
			c.Delta, c.Gamma, c.Theta, c.Vega, c.Rho)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) UpdateContractLastPrice(ctx context.Context, contractID uuid.UUID, lastPrice decimal.Decimal) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE option_contracts SET last_price = $2, volume = volume + 1, open_interest = open_interest + 1, updated_at = NOW() WHERE id = $1`,
		contractID, lastPrice)
	return err
}

// ---- Option Positions ----

func (r *Repo) GetOptionPositionsByPortfolio(ctx context.Context, portfolioID uuid.UUID) ([]model.OptionPosition, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, portfolio_id, contract_id, quantity, avg_cost, collateral, created_at, updated_at
		 FROM option_positions WHERE portfolio_id = $1 AND quantity != 0 ORDER BY created_at`, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []model.OptionPosition
	for rows.Next() {
		var p model.OptionPosition
		if err := rows.Scan(&p.ID, &p.PortfolioID, &p.ContractID, &p.Quantity, &p.AvgCost, &p.Collateral, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

func (r *Repo) GetOptionPosition(ctx context.Context, portfolioID, contractID uuid.UUID) (*model.OptionPosition, error) {
	p := &model.OptionPosition{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, portfolio_id, contract_id, quantity, avg_cost, collateral, created_at, updated_at
		 FROM option_positions WHERE portfolio_id = $1 AND contract_id = $2`, portfolioID, contractID,
	).Scan(&p.ID, &p.PortfolioID, &p.ContractID, &p.Quantity, &p.AvgCost, &p.Collateral, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repo) UpsertOptionPosition(ctx context.Context, portfolioID, contractID uuid.UUID, quantity int, avgCost, collateral decimal.Decimal) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO option_positions (portfolio_id, contract_id, quantity, avg_cost, collateral)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (portfolio_id, contract_id) DO UPDATE SET
			quantity = EXCLUDED.quantity, avg_cost = EXCLUDED.avg_cost, collateral = EXCLUDED.collateral, updated_at = NOW()`,
		portfolioID, contractID, quantity, avgCost, collateral)
	return err
}

func (r *Repo) DeleteOptionPosition(ctx context.Context, portfolioID, contractID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM option_positions WHERE portfolio_id = $1 AND contract_id = $2`,
		portfolioID, contractID)
	return err
}

func (r *Repo) GetPositionsForContract(ctx context.Context, contractID uuid.UUID) ([]model.OptionPosition, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, portfolio_id, contract_id, quantity, avg_cost, collateral, created_at, updated_at
		 FROM option_positions WHERE contract_id = $1 AND quantity != 0`, contractID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []model.OptionPosition
	for rows.Next() {
		var p model.OptionPosition
		if err := rows.Scan(&p.ID, &p.PortfolioID, &p.ContractID, &p.Quantity, &p.AvgCost, &p.Collateral, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

// ---- Option Trades ----

func (r *Repo) CreateOptionTrade(ctx context.Context, userID, contractID uuid.UUID, side string, quantity int, price, total decimal.Decimal) (*model.OptionTrade, error) {
	t := &model.OptionTrade{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO option_trades (user_id, contract_id, side, quantity, price, total)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, contract_id, side, quantity, price, total, created_at`,
		userID, contractID, side, quantity, price, total,
	).Scan(&t.ID, &t.UserID, &t.ContractID, &t.Side, &t.Quantity, &t.Price, &t.Total, &t.CreatedAt)
	return t, err
}

func (r *Repo) GetOptionTradesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.OptionTrade, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, contract_id, side, quantity, price, total, created_at
		 FROM option_trades WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []model.OptionTrade
	for rows.Next() {
		var t model.OptionTrade
		if err := rows.Scan(&t.ID, &t.UserID, &t.ContractID, &t.Side, &t.Quantity, &t.Price, &t.Total, &t.CreatedAt); err != nil {
			return nil, err
		}
		trades = append(trades, t)
	}
	return trades, nil
}

// ---- Option Orders ----

func (r *Repo) CreateOptionOrder(ctx context.Context, userID, contractID uuid.UUID, side, orderType string, quantity int, limitPrice *decimal.Decimal) (*model.OptionOrder, error) {
	o := &model.OptionOrder{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO option_orders (user_id, contract_id, side, order_type, quantity, limit_price)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, contract_id, side, order_type, quantity, limit_price, status, filled_price, filled_at, created_at, updated_at`,
		userID, contractID, side, orderType, quantity, limitPrice,
	).Scan(&o.ID, &o.UserID, &o.ContractID, &o.Side, &o.OrderType, &o.Quantity,
		&o.LimitPrice, &o.Status, &o.FilledPrice, &o.FilledAt, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *Repo) GetOpenOptionOrders(ctx context.Context, userID uuid.UUID) ([]model.OptionOrder, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, contract_id, side, order_type, quantity, limit_price, status, filled_price, filled_at, created_at, updated_at
		 FROM option_orders WHERE user_id = $1 AND status = 'open' ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OptionOrder
	for rows.Next() {
		var o model.OptionOrder
		if err := rows.Scan(&o.ID, &o.UserID, &o.ContractID, &o.Side, &o.OrderType, &o.Quantity,
			&o.LimitPrice, &o.Status, &o.FilledPrice, &o.FilledAt, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *Repo) GetAllOpenOptionOrders(ctx context.Context) ([]model.OptionOrder, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, contract_id, side, order_type, quantity, limit_price, status, filled_price, filled_at, created_at, updated_at
		 FROM option_orders WHERE status = 'open'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OptionOrder
	for rows.Next() {
		var o model.OptionOrder
		if err := rows.Scan(&o.ID, &o.UserID, &o.ContractID, &o.Side, &o.OrderType, &o.Quantity,
			&o.LimitPrice, &o.Status, &o.FilledPrice, &o.FilledAt, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *Repo) FillOptionOrder(ctx context.Context, orderID uuid.UUID, filledPrice decimal.Decimal) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE option_orders SET status = 'filled', filled_price = $2, filled_at = NOW(), updated_at = NOW() WHERE id = $1`,
		orderID, filledPrice)
	return err
}

func (r *Repo) CancelOptionOrder(ctx context.Context, orderID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE option_orders SET status = 'cancelled', updated_at = NOW() WHERE id = $1 AND user_id = $2 AND status = 'open'`,
		orderID, userID)
	return err
}

// --- helpers ---

type scannable interface {
	Next() bool
	Scan(dest ...interface{}) error
}

func scanOptionContracts(rows scannable) ([]model.OptionContract, error) {
	var contracts []model.OptionContract
	for rows.Next() {
		var c model.OptionContract
		if err := rows.Scan(&c.ID, &c.Ticker, &c.OptionType, &c.StrikePrice, &c.Expiration, &c.ContractSymbol,
			&c.BidPrice, &c.AskPrice, &c.LastPrice, &c.MarkPrice, &c.OpenInterest, &c.Volume, &c.ImpliedVol,
			&c.Delta, &c.Gamma, &c.Theta, &c.Vega, &c.Rho, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		contracts = append(contracts, c)
	}
	return contracts, nil
}
