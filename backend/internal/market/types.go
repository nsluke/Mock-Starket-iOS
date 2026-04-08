package market

import "github.com/shopspring/decimal"

// StockState holds the live state for a single stock.
type StockState struct {
	Ticker             string
	Price              float64
	BasePrice          float64
	Volatility         float64
	Drift              float64
	MeanReversionSpeed float64
	Sector             string
	DayOpen            float64
	DayHigh            float64
	DayLow             float64
	Volume             int64
}

// PriceUpdate is broadcast after each price tick.
type PriceUpdate struct {
	Ticker    string          `json:"ticker"`
	Price     decimal.Decimal `json:"price"`
	Change    decimal.Decimal `json:"change"`
	ChangePct decimal.Decimal `json:"change_pct"`
	Volume    int64           `json:"volume"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
}

// MarketEvent represents an event that affects prices.
type MarketEvent struct {
	Type      string `json:"event"`
	Ticker    string `json:"ticker,omitempty"`
	Sector    string `json:"sector,omitempty"`
	Headline  string `json:"headline"`
	Impact    string `json:"impact"`
	Magnitude string `json:"magnitude"`
}

// Observer receives price updates and market events.
type Observer interface {
	OnPriceBatch(updates []PriceUpdate)
	OnMarketEvent(event MarketEvent)
}
