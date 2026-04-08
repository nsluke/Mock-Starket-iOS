package market

import (
	"context"

	"github.com/shopspring/decimal"
)

// PriceProvider is the interface for any market data source (simulation or live).
type PriceProvider interface {
	// GetPrice returns the current price for a ticker.
	GetPrice(ticker string) (decimal.Decimal, bool)

	// GetAllPrices returns current prices for all tracked stocks.
	GetAllPrices() map[string]decimal.Decimal

	// GetAllStockStates returns a copy of all stock states.
	GetAllStockStates() map[string]StockState

	// AddObserver registers a listener for price updates and events.
	AddObserver(obs Observer)

	// Run starts the data feed. Blocks until context is cancelled.
	Run(ctx context.Context) error
}
