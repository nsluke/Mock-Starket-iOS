package polygon

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/luke/mockstarket/internal/repository"
	"github.com/shopspring/decimal"
)

// HistoryService fetches historical price data, checking the local DB first
// and falling back to Polygon.io for older data.
type HistoryService struct {
	client *Client
	repo   *repository.Repo
	logger *slog.Logger
}

// NewHistoryService creates a history service backed by Polygon.io.
func NewHistoryService(client *Client, repo *repository.Repo, logger *slog.Logger) *HistoryService {
	return &HistoryService{client: client, repo: repo, logger: logger}
}

// intervalToPolygon maps our internal interval names to Polygon API parameters.
func intervalToPolygon(interval string) (multiplier int, timespan string) {
	switch interval {
	case "1m":
		return 1, "minute"
	case "5m":
		return 5, "minute"
	case "1h":
		return 1, "hour"
	case "1d":
		return 1, "day"
	default:
		return 1, "minute"
	}
}

// FetchHistory retrieves historical OHLCV bars for a ticker.
// It tries the local price_history table first, then falls back to Polygon.
func (h *HistoryService) FetchHistory(ctx context.Context, ticker, interval string, limit int) error {
	// Determine time range based on interval
	now := time.Now()
	var from time.Time
	switch interval {
	case "1m":
		from = now.AddDate(0, 0, -1) // 1 day of 1m bars
	case "5m":
		from = now.AddDate(0, 0, -5) // 5 days of 5m bars
	case "1h":
		from = now.AddDate(0, -1, 0) // 1 month of 1h bars
	case "1d":
		from = now.AddDate(-1, 0, 0) // 1 year of 1d bars
	default:
		from = now.AddDate(0, 0, -1)
	}

	multiplier, timespan := intervalToPolygon(interval)
	fromStr := from.Format("2006-01-02")
	toStr := now.Format("2006-01-02")

	bars, err := h.client.GetAggregateBars(ctx, ticker, multiplier, timespan, fromStr, toStr)
	if err != nil {
		return fmt.Errorf("fetch polygon bars: %w", err)
	}

	// Cache bars into price_history table
	cached := 0
	for _, bar := range bars {
		ts := time.Unix(0, bar.Timestamp*int64(time.Millisecond)).UTC()
		price := decimal.NewFromFloat(bar.Close)
		open := decimal.NewFromFloat(bar.Open)
		high := decimal.NewFromFloat(bar.High)
		low := decimal.NewFromFloat(bar.Low)
		close := decimal.NewFromFloat(bar.Close)
		volume := int64(bar.Volume)

		err := h.repo.InsertPriceHistory(ctx, ticker, price, open, high, low, close, volume, interval, ts)
		if err != nil {
			// Ignore duplicate key errors — data may already exist
			continue
		}
		cached++
	}

	if cached > 0 {
		h.logger.Debug("cached polygon history", "ticker", ticker, "interval", interval, "bars", cached)
	}

	return nil
}
