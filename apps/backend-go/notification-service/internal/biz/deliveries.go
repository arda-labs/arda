package biz

import (
	"context"
	"fmt"
	"strings"
)

func (uc *NotificationUsecase) ListDeliveries(ctx context.Context, filter DeliveryFilter) ([]*NotificationDelivery, string, error) {
	filter.Status = upperDefault(filter.Status, "")
	filter.Channel = upperDefault(filter.Channel, "")
	filter.RecipientID = strings.TrimSpace(filter.RecipientID)
	return uc.repo.ListDeliveries(ctx, filter)
}

func (uc *NotificationUsecase) RetryDelivery(ctx context.Context, id, actor string) (*NotificationDelivery, error) {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		actor = "SYSTEM"
	}
	return uc.repo.RetryDelivery(ctx, strings.TrimSpace(id), actor)
}

func (uc *NotificationUsecase) ProcessDueDeliveries(ctx context.Context, workerID string, limit int) (int, int, error) {
	workerID = strings.TrimSpace(workerID)
	if workerID == "" {
		workerID = "notification-worker"
	}
	limit = intDefault(limit, 20)

	deliveries, err := uc.repo.ClaimDueDeliveries(ctx, workerID, limit)
	if err != nil {
		return 0, 0, err
	}

	failed := 0
	for _, delivery := range deliveries {
		if err := uc.processDelivery(ctx, delivery); err != nil {
			failed++
		}
	}
	return len(deliveries), failed, nil
}

func (uc *NotificationUsecase) processDelivery(ctx context.Context, delivery *NotificationDelivery) error {
	switch delivery.Channel {
	case "IN_APP":
		inbox, err := uc.repo.CreateInAppNotification(ctx, delivery)
		if err != nil {
			_, _ = uc.repo.MarkDeliveryFailed(ctx, delivery.ID, err.Error(), retryAfterSeconds(delivery.AttemptCount))
			return err
		}
		uc.publishInAppNotification(inbox)
		_, err = uc.repo.MarkDeliveryDelivered(ctx, delivery.ID, "IN_APP_STORE", inbox.ID, `{"status":"stored"}`)
		return err
	default:
		message := fmt.Sprintf("provider adapter not configured for channel %s", delivery.Channel)
		_, err := uc.repo.MarkDeliveryFailed(ctx, delivery.ID, message, retryAfterSeconds(delivery.AttemptCount))
		return err
	}
}

func retryAfterSeconds(attempt int) int {
	if attempt <= 1 {
		return 30
	}
	seconds := 30
	for i := 1; i < attempt && seconds < 1800; i++ {
		seconds *= 2
	}
	if seconds > 1800 {
		return 1800
	}
	return seconds
}
