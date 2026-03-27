package websocket

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"
)

// Message is the JSON envelope for all WebSocket communication.
type Message struct {
	Type      string          `json:"type"`
	Channel   string          `json:"channel,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

// Hub manages all WebSocket client connections and message routing.
type Hub struct {
	mu          sync.RWMutex
	clients     map[*Client]bool
	channels    map[string]map[*Client]bool // channel -> set of clients
	maxClients  int
	logger      *slog.Logger
}

// NewHub creates a new WebSocket hub.
func NewHub(maxClients int, logger *slog.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		channels:   make(map[string]map[*Client]bool),
		maxClients: maxClients,
		logger:     logger,
	}
}

// Register adds a client to the hub.
func (h *Hub) Register(c *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.clients) >= h.maxClients {
		h.logger.Warn("max WebSocket clients reached", "max", h.maxClients)
		return false
	}

	h.clients[c] = true
	h.logger.Debug("client registered", "user_id", c.UserID, "total", len(h.clients))
	return true
}

// Unregister removes a client from the hub and all channels.
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c]; !ok {
		return
	}

	delete(h.clients, c)
	for ch, clients := range h.channels {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.channels, ch)
		}
	}

	h.logger.Debug("client unregistered", "user_id", c.UserID, "total", len(h.clients))
}

// Subscribe adds a client to a channel.
func (h *Hub) Subscribe(c *Client, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.channels[channel] == nil {
		h.channels[channel] = make(map[*Client]bool)
	}
	h.channels[channel][c] = true
	h.logger.Debug("client subscribed", "user_id", c.UserID, "channel", channel)
}

// Unsubscribe removes a client from a channel.
func (h *Hub) Unsubscribe(c *Client, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.channels[channel]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.channels, channel)
		}
	}
}

// BroadcastToChannel sends a message to all clients on a channel.
func (h *Hub) BroadcastToChannel(channel string, msg Message) {
	msg.Timestamp = time.Now().UnixMilli()

	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("failed to marshal broadcast message", "error", err)
		return
	}

	h.mu.RLock()
	clients := h.channels[channel]
	targets := make([]*Client, 0, len(clients))
	for c := range clients {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	for _, c := range targets {
		c.Send(data)
	}
}

// BroadcastToAll sends a message to all connected clients.
func (h *Hub) BroadcastToAll(msg Message) {
	msg.Timestamp = time.Now().UnixMilli()

	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("failed to marshal broadcast message", "error", err)
		return
	}

	h.mu.RLock()
	targets := make([]*Client, 0, len(h.clients))
	for c := range h.clients {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	for _, c := range targets {
		c.Send(data)
	}
}

// SendToUser sends a message to a specific user's connections.
func (h *Hub) SendToUser(userID string, msg Message) {
	msg.Timestamp = time.Now().UnixMilli()

	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("failed to marshal user message", "error", err)
		return
	}

	h.mu.RLock()
	targets := make([]*Client, 0)
	for c := range h.clients {
		if c.UserID == userID {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range targets {
		c.Send(data)
	}
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// CleanStale removes clients that haven't pinged recently.
func (h *Hub) CleanStale(maxAge time.Duration) int {
	h.mu.RLock()
	now := time.Now()
	stale := make([]*Client, 0)
	for c := range h.clients {
		if now.Sub(c.LastPing) > maxAge {
			stale = append(stale, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range stale {
		h.Unregister(c)
		c.Close()
	}

	if len(stale) > 0 {
		h.logger.Info("cleaned stale WebSocket clients", "count", len(stale))
	}

	return len(stale)
}
