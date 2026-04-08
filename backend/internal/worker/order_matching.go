package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/service"
	"github.com/luke/mockstarket/internal/market"
	ws "github.com/luke/mockstarket/internal/websocket"
	"github.com/shopspring/decimal"
)

// OrderMatchingWorker evaluates open limit/stop/stop-limit orders against
// live prices and executes trades when conditions are met.
type OrderMatchingWorker struct {
	repo     *repository.Repo
	tradeSvc *service.TradeService
	hub      *ws.Hub
	logger   *slog.Logger

	mu     sync.Mutex
	prices map[string]decimal.Decimal // latest price snapshot
}

// NewOrderMatchingWorker creates the order matching worker.
func NewOrderMatchingWorker(repo *repository.Repo, tradeSvc *service.TradeService, hub *ws.Hub, logger *slog.Logger) *OrderMatchingWorker {
	return &OrderMatchingWorker{
		repo:     repo,
		tradeSvc: tradeSvc,
		hub:      hub,
		logger:   logger,
		prices:   make(map[string]decimal.Decimal),
	}
}

// OnPriceBatch updates the latest prices and triggers order matching.
func (w *OrderMatchingWorker) OnPriceBatch(updates []market.PriceUpdate) {
	w.mu.Lock()
	for _, u := range updates {
		w.prices[u.Ticker] = u.Price
	}
	priceSnap := make(map[string]decimal.Decimal, len(w.prices))
	for k, v := range w.prices {
		priceSnap[k] = v
	}
	w.mu.Unlock()

	w.matchOrders(priceSnap)
}

// OnMarketEvent is a no-op for this worker.
func (w *OrderMatchingWorker) OnMarketEvent(_ market.MarketEvent) {}

// matchOrders fetches all open orders and checks if any should be filled.
func (w *OrderMatchingWorker) matchOrders(prices map[string]decimal.Decimal) {
	ctx := context.Background()

	orders, err := w.repo.GetAllOpenOrders(ctx)
	if err != nil {
		w.logger.Error("order matching: failed to fetch open orders", "error", err)
		return
	}

	for _, order := range orders {
		price, ok := prices[order.Ticker]
		if !ok {
			continue
		}

		if !w.shouldFill(order.OrderType, order.Side, price, order.LimitPrice, order.StopPrice) {
			continue
		}

		// Execute the trade through the trade service
		trade, err := w.tradeSvc.ExecuteTrade(ctx, order.UserID, service.TradeRequest{
			Ticker: order.Ticker,
			Side:   order.Side,
			Shares: order.Shares,
		})
		if err != nil {
			w.logger.Warn("order matching: failed to execute order",
				"order_id", order.ID, "ticker", order.Ticker, "error", err)
			continue
		}

		// Mark the order as filled
		if err := w.repo.FillOrder(ctx, order.ID, price); err != nil {
			w.logger.Error("order matching: failed to mark order filled",
				"order_id", order.ID, "error", err)
			continue
		}

		w.logger.Info("order filled",
			"order_id", order.ID,
			"ticker", order.Ticker,
			"side", order.Side,
			"type", order.OrderType,
			"price", price,
			"trade_id", trade.ID,
		)

		// Notify user via WebSocket
		w.hub.SendToUser(order.UserID.String(), ws.Message{
			Type: "order_filled",
			Data: mustMarshal(map[string]interface{}{
				"order_id": order.ID,
				"trade_id": trade.ID,
				"ticker":   order.Ticker,
				"side":     order.Side,
				"shares":   order.Shares,
				"price":    price,
			}),
		})
	}
}

// shouldFill determines if an order's conditions are met at the current price.
func (w *OrderMatchingWorker) shouldFill(orderType, side string, currentPrice decimal.Decimal, limitPrice, stopPrice *decimal.Decimal) bool {
	switch orderType {
	case "limit":
		if limitPrice == nil {
			return false
		}
		if side == "buy" {
			// Buy limit: execute when price drops to or below limit
			return currentPrice.LessThanOrEqual(*limitPrice)
		}
		// Sell limit: execute when price rises to or above limit
		return currentPrice.GreaterThanOrEqual(*limitPrice)

	case "stop":
		if stopPrice == nil {
			return false
		}
		if side == "buy" {
			// Buy stop: execute when price rises to or above stop (breakout buy)
			return currentPrice.GreaterThanOrEqual(*stopPrice)
		}
		// Sell stop (stop-loss): execute when price drops to or below stop
		return currentPrice.LessThanOrEqual(*stopPrice)

	case "stop_limit":
		if stopPrice == nil || limitPrice == nil {
			return false
		}
		if side == "buy" {
			// Stop triggered (price >= stop) AND limit condition met (price <= limit)
			return currentPrice.GreaterThanOrEqual(*stopPrice) && currentPrice.LessThanOrEqual(*limitPrice)
		}
		// Stop triggered (price <= stop) AND limit condition met (price >= limit)
		return currentPrice.LessThanOrEqual(*stopPrice) && currentPrice.GreaterThanOrEqual(*limitPrice)

	default:
		return false
	}
}
