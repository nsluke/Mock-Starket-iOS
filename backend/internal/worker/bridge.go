package worker

import (
	"encoding/json"
	"log/slog"

	"github.com/luke/mockstarket/internal/simulation"
	ws "github.com/luke/mockstarket/internal/websocket"
)

// SimulationBridge connects the simulation engine to the WebSocket hub.
type SimulationBridge struct {
	hub    *ws.Hub
	logger *slog.Logger
}

// NewSimulationBridge creates a bridge between simulation and WebSocket.
func NewSimulationBridge(hub *ws.Hub, logger *slog.Logger) *SimulationBridge {
	return &SimulationBridge{hub: hub, logger: logger}
}

// OnPriceBatch broadcasts price updates to the "market" channel.
func (b *SimulationBridge) OnPriceBatch(updates []simulation.PriceUpdate) {
	data, err := json.Marshal(updates)
	if err != nil {
		b.logger.Error("failed to marshal price batch", "error", err)
		return
	}

	b.hub.BroadcastToChannel("market", ws.Message{
		Type: "price_batch",
		Data: data,
	})
}

// OnMarketEvent broadcasts market events to all connected clients.
func (b *SimulationBridge) OnMarketEvent(event simulation.MarketEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		b.logger.Error("failed to marshal market event", "error", err)
		return
	}

	b.hub.BroadcastToAll(ws.Message{
		Type: "market_event",
		Data: data,
	})
}
