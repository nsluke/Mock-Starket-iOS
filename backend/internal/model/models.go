package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	FirebaseUID   string    `json:"firebase_uid" db:"firebase_uid"`
	DisplayName   string    `json:"display_name" db:"display_name"`
	AvatarURL     *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	IsGuest       bool      `json:"is_guest" db:"is_guest"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	LoginStreak   int       `json:"login_streak" db:"login_streak"`
	LongestStreak int       `json:"longest_streak" db:"longest_streak"`
}

type Stock struct {
	Ticker         string          `json:"ticker" db:"ticker"`
	Name           string          `json:"name" db:"name"`
	Sector         string          `json:"sector" db:"sector"`
	AssetType      string          `json:"asset_type" db:"asset_type"` // stock, etf, crypto, commodity
	BasePrice      decimal.Decimal `json:"base_price" db:"base_price"`
	CurrentPrice   decimal.Decimal `json:"current_price" db:"current_price"`
	DayOpen        decimal.Decimal `json:"day_open" db:"day_open"`
	DayHigh        decimal.Decimal `json:"day_high" db:"day_high"`
	DayLow         decimal.Decimal `json:"day_low" db:"day_low"`
	PrevClose      decimal.Decimal `json:"prev_close" db:"prev_close"`
	Volume         int64           `json:"volume" db:"volume"`
	Volatility     decimal.Decimal `json:"volatility" db:"volatility"`
	Drift          decimal.Decimal `json:"drift" db:"drift"`
	MeanReversion  decimal.Decimal `json:"mean_reversion" db:"mean_reversion"`
	Description    *string         `json:"description,omitempty" db:"description"`
	LogoURL        *string         `json:"logo_url,omitempty" db:"logo_url"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

type ETFHolding struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	ETFTicker     string          `json:"etf_ticker" db:"etf_ticker"`
	HoldingTicker string          `json:"holding_ticker" db:"holding_ticker"`
	Weight        decimal.Decimal `json:"weight" db:"weight"`
}

type Portfolio struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	Cash      decimal.Decimal `json:"cash" db:"cash"`
	NetWorth  decimal.Decimal `json:"net_worth" db:"net_worth"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type Holding struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	PortfolioID uuid.UUID       `json:"portfolio_id" db:"portfolio_id"`
	Ticker      string          `json:"ticker" db:"ticker"`
	Shares      int             `json:"shares" db:"shares"`
	AvgCost     decimal.Decimal `json:"avg_cost" db:"avg_cost"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type Trade struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	Ticker    string          `json:"ticker" db:"ticker"`
	Side      string          `json:"side" db:"side"`
	Shares    int             `json:"shares" db:"shares"`
	Price     decimal.Decimal `json:"price" db:"price"`
	Total     decimal.Decimal `json:"total" db:"total"`
	OrderID   *uuid.UUID      `json:"order_id,omitempty" db:"order_id"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

type Order struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	UserID      uuid.UUID        `json:"user_id" db:"user_id"`
	Ticker      string           `json:"ticker" db:"ticker"`
	Side        string           `json:"side" db:"side"`
	OrderType   string           `json:"order_type" db:"order_type"`
	Shares      int              `json:"shares" db:"shares"`
	LimitPrice  *decimal.Decimal `json:"limit_price,omitempty" db:"limit_price"`
	StopPrice   *decimal.Decimal `json:"stop_price,omitempty" db:"stop_price"`
	Status      string           `json:"status" db:"status"`
	FilledPrice *decimal.Decimal `json:"filled_price,omitempty" db:"filled_price"`
	FilledAt    *time.Time       `json:"filled_at,omitempty" db:"filled_at"`
	ExpiresAt   *time.Time       `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

type PricePoint struct {
	ID         int64           `json:"id" db:"id"`
	Ticker     string          `json:"ticker" db:"ticker"`
	Price      decimal.Decimal `json:"price" db:"price"`
	Open       decimal.Decimal `json:"open" db:"open"`
	High       decimal.Decimal `json:"high" db:"high"`
	Low        decimal.Decimal `json:"low" db:"low"`
	Close      decimal.Decimal `json:"close" db:"close"`
	Volume     int64           `json:"volume" db:"volume"`
	Interval   string          `json:"interval" db:"interval"`
	RecordedAt time.Time       `json:"recorded_at" db:"recorded_at"`
}

type LeaderboardEntry struct {
	ID          int64           `json:"id" db:"id"`
	UserID      uuid.UUID       `json:"user_id" db:"user_id"`
	DisplayName string          `json:"display_name" db:"display_name"`
	NetWorth    decimal.Decimal `json:"net_worth" db:"net_worth"`
	TotalReturn decimal.Decimal `json:"total_return" db:"total_return"`
	Rank        int             `json:"rank" db:"rank"`
	Period      string          `json:"period" db:"period"`
	ComputedAt  time.Time       `json:"computed_at" db:"computed_at"`
}

type Achievement struct {
	ID           string  `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	Description  string  `json:"description" db:"description"`
	Icon         string  `json:"icon" db:"icon"`
	Category     string  `json:"category" db:"category"`
	CriteriaJSON *string `json:"criteria_json,omitempty" db:"criteria_json"`
}

type UserAchievement struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	AchievementID string    `json:"achievement_id" db:"achievement_id"`
	EarnedAt      time.Time `json:"earned_at" db:"earned_at"`
}

type PriceAlert struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	UserID      uuid.UUID       `json:"user_id" db:"user_id"`
	Ticker      string          `json:"ticker" db:"ticker"`
	Condition   string          `json:"condition" db:"condition"`
	TargetPrice decimal.Decimal `json:"target_price" db:"target_price"`
	Triggered   bool            `json:"triggered" db:"triggered"`
	TriggeredAt *time.Time      `json:"triggered_at,omitempty" db:"triggered_at"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

type DailyChallenge struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	Date          time.Time       `json:"date" db:"date"`
	ChallengeType string          `json:"challenge_type" db:"challenge_type"`
	Description   string          `json:"description" db:"description"`
	TargetJSON    string          `json:"target_json" db:"target_json"`
	RewardCash    decimal.Decimal `json:"reward_cash" db:"reward_cash"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

type UserChallenge struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	ChallengeID uuid.UUID  `json:"challenge_id" db:"challenge_id"`
	Completed   bool       `json:"completed" db:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	Claimed     bool       `json:"claimed" db:"claimed"`
}

type WatchlistItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Ticker    string    `json:"ticker" db:"ticker"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PortfolioHistory struct {
	ID         int64           `json:"id" db:"id"`
	UserID     uuid.UUID       `json:"user_id" db:"user_id"`
	NetWorth   decimal.Decimal `json:"net_worth" db:"net_worth"`
	Cash       decimal.Decimal `json:"cash" db:"cash"`
	RecordedAt time.Time       `json:"recorded_at" db:"recorded_at"`
}

// ---- Options ----

type OptionContract struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	Ticker         string          `json:"ticker" db:"ticker"`
	OptionType     string          `json:"option_type" db:"option_type"`
	StrikePrice    decimal.Decimal `json:"strike_price" db:"strike_price"`
	Expiration     time.Time       `json:"expiration" db:"expiration"`
	ContractSymbol string          `json:"contract_symbol" db:"contract_symbol"`
	BidPrice       decimal.Decimal `json:"bid_price" db:"bid_price"`
	AskPrice       decimal.Decimal `json:"ask_price" db:"ask_price"`
	LastPrice      decimal.Decimal `json:"last_price" db:"last_price"`
	MarkPrice      decimal.Decimal `json:"mark_price" db:"mark_price"`
	OpenInterest   int             `json:"open_interest" db:"open_interest"`
	Volume         int             `json:"volume" db:"volume"`
	ImpliedVol     decimal.Decimal `json:"implied_vol" db:"implied_vol"`
	Delta          decimal.Decimal `json:"delta" db:"delta"`
	Gamma          decimal.Decimal `json:"gamma" db:"gamma"`
	Theta          decimal.Decimal `json:"theta" db:"theta"`
	Vega           decimal.Decimal `json:"vega" db:"vega"`
	Rho            decimal.Decimal `json:"rho" db:"rho"`
	Status         string          `json:"status" db:"status"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

type OptionPosition struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	PortfolioID uuid.UUID       `json:"portfolio_id" db:"portfolio_id"`
	ContractID  uuid.UUID       `json:"contract_id" db:"contract_id"`
	Quantity    int             `json:"quantity" db:"quantity"`
	AvgCost     decimal.Decimal `json:"avg_cost" db:"avg_cost"`
	Collateral  decimal.Decimal `json:"collateral" db:"collateral"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type OptionTrade struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	UserID     uuid.UUID       `json:"user_id" db:"user_id"`
	ContractID uuid.UUID       `json:"contract_id" db:"contract_id"`
	Side       string          `json:"side" db:"side"`
	Quantity   int             `json:"quantity" db:"quantity"`
	Price      decimal.Decimal `json:"price" db:"price"`
	Total      decimal.Decimal `json:"total" db:"total"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

type OptionOrder struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	UserID      uuid.UUID        `json:"user_id" db:"user_id"`
	ContractID  uuid.UUID        `json:"contract_id" db:"contract_id"`
	Side        string           `json:"side" db:"side"`
	OrderType   string           `json:"order_type" db:"order_type"`
	Quantity    int              `json:"quantity" db:"quantity"`
	LimitPrice  *decimal.Decimal `json:"limit_price,omitempty" db:"limit_price"`
	Status      string           `json:"status" db:"status"`
	FilledPrice *decimal.Decimal `json:"filled_price,omitempty" db:"filled_price"`
	FilledAt    *time.Time       `json:"filled_at,omitempty" db:"filled_at"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}
