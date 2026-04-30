package biz

import (
	"context"
	"strings"
)

func (uc *NotificationUsecase) ListInAppNotifications(ctx context.Context, filter InAppFilter) ([]*InAppNotification, string, error) {
	filter.RecipientType = upperDefault(filter.RecipientType, "USER")
	filter.RecipientID = strings.TrimSpace(filter.RecipientID)
	filter.Status = upperDefault(filter.Status, "")
	if filter.RecipientID == "" {
		return nil, "", ErrInvalidArgument
	}
	return uc.repo.ListInAppNotifications(ctx, filter)
}

func (uc *NotificationUsecase) MarkInAppNotificationRead(ctx context.Context, id, actor string) (*InAppNotification, error) {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		actor = "SYSTEM"
	}
	return uc.repo.MarkInAppNotificationRead(ctx, strings.TrimSpace(id), actor)
}
