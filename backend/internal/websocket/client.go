package websocket

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 30 * time.Second
	maxMessageSize = 4096
	sendBufferSize = 256
)

// Client represents a single WebSocket connection.
type Client struct {
	UserID   string
	conn     *websocket.Conn
	hub      *Hub
	send     chan []byte
	done     chan struct{}
	once     sync.Once
	LastPing time.Time
	logger   *slog.Logger
}

// NewClient creates a new WebSocket client.
func NewClient(conn *websocket.Conn, hub *Hub, userID string, logger *slog.Logger) *Client {
	return &Client{
		UserID:   userID,
		conn:     conn,
		hub:      hub,
		send:     make(chan []byte, sendBufferSize),
		done:     make(chan struct{}),
		LastPing: time.Now(),
		logger:   logger,
	}
}

// Send queues a message for delivery to the client.
func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
		// Buffer full, client is slow — disconnect
		c.logger.Warn("client send buffer full, disconnecting", "user_id", c.UserID)
		c.Close()
	}
}

// Close terminates the client connection.
func (c *Client) Close() {
	c.once.Do(func() {
		close(c.done)
		c.conn.Close()
	})
}

// Run starts the read and write pumps. Blocks until the client disconnects.
func (c *Client) Run() {
	go c.writePump()
	c.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister(c)
		c.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.LastPing = time.Now()
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				c.logger.Debug("WebSocket read error", "user_id", c.UserID, "error", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Debug("invalid WebSocket message", "user_id", c.UserID, "error", err)
			continue
		}

		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case <-c.done:
			return
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "subscribe":
		if msg.Channel != "" {
			c.hub.Subscribe(c, msg.Channel)
		}
	case "unsubscribe":
		if msg.Channel != "" {
			c.hub.Unsubscribe(c, msg.Channel)
		}
	case "ping":
		c.LastPing = time.Now()
		pong := Message{Type: "pong"}
		data, _ := json.Marshal(pong)
		c.Send(data)
	default:
		c.logger.Debug("unknown message type", "type", msg.Type, "user_id", c.UserID)
	}
}
