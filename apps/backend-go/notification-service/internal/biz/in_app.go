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

func (uc *NotificationUsecase) CountUnreadInAppNotifications(ctx context.Context, recipientType, recipientID string) (int, error) {
	recipientType = upperDefault(recipientType, "USER")
	recipientID = strings.TrimSpace(recipientID)
	if recipientID == "" {
		return 0, ErrInvalidArgument
	}
	return uc.repo.CountUnreadInAppNotifications(ctx, recipientType, recipientID)
}

func (uc *NotificationUsecase) MarkAllInAppNotificationsRead(ctx context.Context, recipientType, recipientID, actor string) (int, error) {
	recipientType = upperDefault(recipientType, "USER")
	recipientID = strings.TrimSpace(recipientID)
	actor = strings.TrimSpace(actor)
	if actor == "" {
		actor = "SYSTEM"
	}
	if recipientID == "" {
		return 0, ErrInvalidArgument
	}
	return uc.repo.MarkAllInAppNotificationsRead(ctx, recipientType, recipientID, actor)
}
