package polygon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Client is a Polygon.io REST API client with rate limiting and caching.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
	limiter *rate.Limiter
	logger  *slog.Logger

	mu    sync.RWMutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	data      []byte
	expiresAt time.Time
}

// NewClient creates a Polygon.io REST client.
// rateLimit is requests per minute (5 for free tier).
func NewClient(apiKey, baseURL string, rateLimit int, logger *slog.Logger) *Client {
	if rateLimit <= 0 {
		rateLimit = 5
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
		logger:  logger,
		cache:   make(map[string]cacheEntry),
	}
}

// doRequest performs a rate-limited, cached GET request.
func (c *Client) doRequest(ctx context.Context, path string, cacheTTL time.Duration) ([]byte, error) {
	cacheKey := path

	// Check cache
	if cacheTTL > 0 {
		c.mu.RLock()
		if entry, ok := c.cache[cacheKey]; ok && time.Now().Before(entry.expiresAt) {
			c.mu.RUnlock()
			return entry.data, nil
		}
		c.mu.RUnlock()
	}

	// Wait for rate limiter
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	url := c.baseURL + path
	if strings.Contains(path, "?") {
		url += "&apiKey=" + c.apiKey
	} else {
		url += "?apiKey=" + c.apiKey
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.logger.Debug("polygon API request", "path", path)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("polygon rate limit exceeded (HTTP 429)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("polygon API error: status %d, body: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
	}

	// Cache the response
	if cacheTTL > 0 {
		c.mu.Lock()
		c.cache[cacheKey] = cacheEntry{data: body, expiresAt: time.Now().Add(cacheTTL)}
		c.mu.Unlock()
	}

	return body, nil
}

// GetAllSnapshots fetches current snapshots for all US stock tickers.
func (c *Client) GetAllSnapshots(ctx context.Context) ([]Snapshot, error) {
	data, err := c.doRequest(ctx, "/v2/snapshot/locale/us/markets/stocks/tickers", 30*time.Second)
	if err != nil {
		return nil, err
	}
	var resp SnapshotResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal snapshots: %w", err)
	}
	return resp.Tickers, nil
}

// GetSnapshot fetches the current snapshot for a single ticker.
func (c *Client) GetSnapshot(ctx context.Context, ticker string) (*Snapshot, error) {
	path := fmt.Sprintf("/v2/snapshot/locale/us/markets/stocks/tickers/%s", ticker)
	data, err := c.doRequest(ctx, path, 15*time.Second)
	if err != nil {
		return nil, err
	}
	var resp SingleSnapshotResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &resp.Ticker, nil
}

// GetAggregateBars fetches historical OHLCV bars.
// timespan: "minute", "hour", "day", "week", "month"
func (c *Client) GetAggregateBars(ctx context.Context, ticker string, multiplier int, timespan string, from, to string) ([]AggBar, error) {
	path := fmt.Sprintf("/v2/aggs/ticker/%s/range/%d/%s/%s/%s?adjusted=true&sort=asc&limit=5000",
		ticker, multiplier, timespan, from, to)
	data, err := c.doRequest(ctx, path, 5*time.Minute)
	if err != nil {
		return nil, err
	}
	var resp AggregateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal aggregates: %w", err)
	}
	return resp.Results, nil
}

// GetTickerDetails fetches reference data for a ticker.
func (c *Client) GetTickerDetails(ctx context.Context, ticker string) (*TickerDetail, error) {
	path := fmt.Sprintf("/v3/reference/tickers/%s", ticker)
	data, err := c.doRequest(ctx, path, 1*time.Hour)
	if err != nil {
		return nil, err
	}
	var resp TickerDetailResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal ticker details: %w", err)
	}
	return &resp.Results, nil
}

// GetPreviousClose fetches the previous day's OHLCV for a ticker.
func (c *Client) GetPreviousClose(ctx context.Context, ticker string) (*PrevCloseBar, error) {
	path := fmt.Sprintf("/v2/aggs/ticker/%s/prev?adjusted=true", ticker)
	data, err := c.doRequest(ctx, path, 5*time.Minute)
	if err != nil {
		return nil, err
	}
	var resp PreviousCloseResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal previous close: %w", err)
	}
	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("no previous close data for %s", ticker)
	}
	return &resp.Results[0], nil
}

// GetMarketStatus fetches the current market status.
func (c *Client) GetMarketStatus(ctx context.Context) (*MarketStatusResponse, error) {
	data, err := c.doRequest(ctx, "/v1/marketstatus/now", 60*time.Second)
	if err != nil {
		return nil, err
	}
	var resp MarketStatusResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal market status: %w", err)
	}
	return &resp, nil
}

// ClearCache removes all cached entries.
func (c *Client) ClearCache() {
	c.mu.Lock()
	c.cache = make(map[string]cacheEntry)
	c.mu.Unlock()
}
