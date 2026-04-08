package polygon

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient manages a WebSocket connection to Polygon.io for real-time data.
type WSClient struct {
	url    string
	apiKey string
	logger *slog.Logger

	mu          sync.Mutex
	conn        *websocket.Conn
	subscribed  map[string]bool
	onAggregate func(WSAggregateMessage)

	done chan struct{}
}

// NewWSClient creates a new Polygon.io WebSocket client.
// wsURL is typically "wss://socket.polygon.io/stocks" or "wss://socket.polygon.io/crypto".
func NewWSClient(wsURL, apiKey string, logger *slog.Logger) *WSClient {
	return &WSClient{
		url:        wsURL,
		apiKey:     apiKey,
		logger:     logger,
		subscribed: make(map[string]bool),
		done:       make(chan struct{}),
	}
}

// OnAggregate registers a handler for aggregate messages (AM.* events).
func (w *WSClient) OnAggregate(fn func(WSAggregateMessage)) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.onAggregate = fn
}

// Subscribe adds tickers to the subscription list.
// prefix is "AM" for per-minute aggregates, "A" for per-second.
func (w *WSClient) Subscribe(prefix string, tickers []string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	params := make([]string, len(tickers))
	for i, t := range tickers {
		key := prefix + "." + t
		params[i] = key
		w.subscribed[key] = true
	}

	if w.conn == nil {
		return nil // Will subscribe on connect
	}

	msg := WSMessage{Action: "subscribe", Params: joinParams(params)}
	return w.conn.WriteJSON(msg)
}

// Run connects to the WebSocket and reads messages until context is cancelled.
// Reconnects automatically on disconnect.
func (w *WSClient) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			w.close()
			return ctx.Err()
		default:
		}

		if err := w.connectAndRead(ctx); err != nil {
			w.logger.Warn("polygon websocket disconnected", "error", err)
		}

		// Reconnect backoff
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			w.logger.Info("polygon websocket reconnecting...")
		}
	}
}

func (w *WSClient) connectAndRead(ctx context.Context) error {
	w.logger.Info("connecting to polygon websocket", "url", w.url)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, w.url, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	w.mu.Lock()
	w.conn = conn
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		w.conn = nil
		w.mu.Unlock()
		_ = conn.Close()
	}()

	// Authenticate
	auth := WSMessage{Action: "auth", Params: w.apiKey}
	if err := conn.WriteJSON(auth); err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	// Re-subscribe to all channels
	w.mu.Lock()
	if len(w.subscribed) > 0 {
		params := make([]string, 0, len(w.subscribed))
		for k := range w.subscribed {
			params = append(params, k)
		}
		msg := WSMessage{Action: "subscribe", Params: joinParams(params)}
		if err := conn.WriteJSON(msg); err != nil {
			w.mu.Unlock()
			return fmt.Errorf("subscribe: %w", err)
		}
	}
	w.mu.Unlock()

	w.logger.Info("polygon websocket connected")

	// Read loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		w.handleMessage(message)
	}
}

func (w *WSClient) handleMessage(data []byte) {
	// Polygon sends arrays of messages
	var messages []json.RawMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		// Might be a single message
		w.handleSingleMessage(data)
		return
	}

	for _, msg := range messages {
		w.handleSingleMessage(msg)
	}
}

func (w *WSClient) handleSingleMessage(data []byte) {
	// Peek at the event type
	var peek struct {
		Event string `json:"ev"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return
	}

	switch peek.Event {
	case "AM", "A":
		var agg WSAggregateMessage
		if err := json.Unmarshal(data, &agg); err != nil {
			w.logger.Error("failed to unmarshal aggregate", "error", err)
			return
		}
		w.mu.Lock()
		handler := w.onAggregate
		w.mu.Unlock()
		if handler != nil {
			handler(agg)
		}

	case "status":
		var status WSStatusMessage
		if err := json.Unmarshal(data, &status); err == nil {
			w.logger.Debug("polygon ws status", "status", status.Status, "message", status.Message)
		}
	}
}

func (w *WSClient) close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.conn != nil {
		_ = w.conn.Close()
		w.conn = nil
	}
}

func joinParams(params []string) string {
	result := ""
	for i, p := range params {
		if i > 0 {
			result += ","
		}
		result += p
	}
	return result
}
