package worker

import (
	"context"
	"log"
	"time"

	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/biz"
)

type CleanupWorker struct {
	uc        *biz.MediaUsecase
	interval  time.Duration
	olderThan time.Duration
	limit     int
}

func NewCleanupWorker(uc *biz.MediaUsecase, interval, olderThan time.Duration, limit int) *CleanupWorker {
	return &CleanupWorker{uc: uc, interval: interval, olderThan: olderThan, limit: limit}
}

func (w *CleanupWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	w.run(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.run(ctx)
		}
	}
}

func (w *CleanupWorker) run(ctx context.Context) {
	before := time.Now().UTC().Add(-w.olderThan)
	cleaned, err := w.uc.CleanupPending(ctx, before, w.limit)
	if err != nil {
		log.Printf("media cleanup failed: %v", err)
		return
	}
	if cleaned > 0 {
		log.Printf("media cleanup removed %d pending upload(s)", cleaned)
	}
}
