package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/market"
	ws "github.com/luke/mockstarket/internal/websocket"
	"github.com/shopspring/decimal"
)

// PriceAlertWorker evaluates untriggered price alerts against live prices
// and notifies users via WebSocket when conditions are met.
type PriceAlertWorker struct {
	repo   *repository.Repo
	hub    *ws.Hub
	logger *slog.Logger

	mu     sync.Mutex
	prices map[string]decimal.Decimal
}

// NewPriceAlertWorker creates a new price alert evaluation worker.
func NewPriceAlertWorker(repo *repository.Repo, hub *ws.Hub, logger *slog.Logger) *PriceAlertWorker {
	return &PriceAlertWorker{
		repo:   repo,
		hub:    hub,
		logger: logger,
		prices: make(map[string]decimal.Decimal),
	}
}

// OnPriceBatch updates prices and checks alerts.
func (w *PriceAlertWorker) OnPriceBatch(updates []market.PriceUpdate) {
	w.mu.Lock()
	for _, u := range updates {
		w.prices[u.Ticker] = u.Price
	}
	priceSnap := make(map[string]decimal.Decimal, len(w.prices))
	for k, v := range w.prices {
		priceSnap[k] = v
	}
	w.mu.Unlock()

	w.checkAlerts(priceSnap)
}

// OnMarketEvent is a no-op for this worker.
func (w *PriceAlertWorker) OnMarketEvent(_ market.MarketEvent) {}

func (w *PriceAlertWorker) checkAlerts(prices map[string]decimal.Decimal) {
	ctx := context.Background()

	alerts, err := w.repo.GetUntriggeredAlerts(ctx)
	if err != nil {
		w.logger.Error("price alerts: failed to fetch alerts", "error", err)
		return
	}

	for _, alert := range alerts {
		price, ok := prices[alert.Ticker]
		if !ok {
			continue
		}

		triggered := false
		switch alert.Condition {
		case "above":
			triggered = price.GreaterThanOrEqual(alert.TargetPrice)
		case "below":
			triggered = price.LessThanOrEqual(alert.TargetPrice)
		}

		if !triggered {
			continue
		}

		if err := w.repo.TriggerAlert(ctx, alert.ID); err != nil {
			w.logger.Error("price alerts: failed to trigger alert",
				"alert_id", alert.ID, "error", err)
			continue
		}

		w.logger.Info("price alert triggered",
			"alert_id", alert.ID,
			"ticker", alert.Ticker,
			"condition", alert.Condition,
			"target", alert.TargetPrice,
			"current", price,
		)

		// Notify the user via WebSocket
		w.hub.SendToUser(alert.UserID.String(), ws.Message{
			Type: "alert_triggered",
			Data: mustMarshal(map[string]interface{}{
				"alert_id":     alert.ID,
				"ticker":       alert.Ticker,
				"condition":    alert.Condition,
				"target_price": alert.TargetPrice,
				"current_price": price,
			}),
		})
	}
}
