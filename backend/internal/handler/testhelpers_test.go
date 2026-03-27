package handler

import (
	"log/slog"
	"os"

	ws "github.com/luke/mockstarket/internal/websocket"
)

func newTestHub() *ws.Hub {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError + 10}))
	return ws.NewHub(100, logger)
}
