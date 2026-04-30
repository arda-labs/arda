package worker

import (
	"context"
	"time"

	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type DeliveryWorker struct {
	uc        *biz.NotificationUsecase
	logger    *log.Helper
	workerID  string
	interval  time.Duration
	batchSize int
}

func NewDeliveryWorker(uc *biz.NotificationUsecase, logger log.Logger, workerID string, interval time.Duration, batchSize int) *DeliveryWorker {
	return &DeliveryWorker{
		uc:        uc,
		logger:    log.NewHelper(logger),
		workerID:  workerID,
		interval:  interval,
		batchSize: batchSize,
	}
}

func (w *DeliveryWorker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		w.runOnce(ctx)
		for {
			select {
			case <-ctx.Done():
				w.logger.Info("notification delivery worker stopped")
				return
			case <-ticker.C:
				w.runOnce(ctx)
			}
		}
	}()
}

func (w *DeliveryWorker) runOnce(ctx context.Context) {
	processed, failed, err := w.uc.ProcessDueDeliveries(ctx, w.workerID, w.batchSize)
	if err != nil {
		w.logger.Errorf("processing notification deliveries failed: %v", err)
		return
	}
	if processed > 0 || failed > 0 {
		w.logger.Infof("notification deliveries processed=%d failed=%d", processed, failed)
	}
}
