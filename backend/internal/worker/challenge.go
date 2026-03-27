package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/luke/mockstarket/internal/service"
)

// ChallengeWorker generates daily challenges and keeps them up to date.
type ChallengeWorker struct {
	challengeSvc *service.ChallengeService
	logger       *slog.Logger
}

// NewChallengeWorker creates a worker that generates daily challenges.
func NewChallengeWorker(challengeSvc *service.ChallengeService, logger *slog.Logger) *ChallengeWorker {
	return &ChallengeWorker{
		challengeSvc: challengeSvc,
		logger:       logger,
	}
}

// Run starts the challenge generation loop. Blocks until context is cancelled.
func (w *ChallengeWorker) Run(ctx context.Context) {
	// Generate immediately on start
	w.generate(ctx)

	// Check every 5 minutes (handles server restarts mid-day and midnight rollover)
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	w.logger.Info("challenge worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("challenge worker stopped")
			return
		case <-ticker.C:
			w.generate(ctx)
		}
	}
}

func (w *ChallengeWorker) generate(ctx context.Context) {
	_, err := w.challengeSvc.GenerateToday(ctx)
	if err != nil {
		w.logger.Error("challenge worker: failed to generate challenge", "error", err)
	}
}
