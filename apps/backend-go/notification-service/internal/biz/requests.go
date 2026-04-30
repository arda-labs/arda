package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (uc *NotificationUsecase) CreateNotificationRequest(ctx context.Context, item *NotificationRequest) (*NotificationRequest, error) {
	normalizeNotificationRequest(item)
	if item.SourceService == "" || item.EventType == "" || item.TemplateCode == "" || item.RecipientID == "" {
		return nil, ErrInvalidArgument
	}
	if !json.Valid([]byte(item.PayloadJSON)) {
		return nil, ErrInvalidArgument
	}
	return uc.repo.CreateNotificationRequest(ctx, item)
}

func (uc *NotificationUsecase) GetNotificationRequest(ctx context.Context, id string) (*NotificationRequest, error) {
	return uc.repo.GetNotificationRequest(ctx, strings.TrimSpace(id))
}

func normalizeNotificationRequest(item *NotificationRequest) {
	item.SourceService = upperDefault(item.SourceService, "")
	item.EventType = upperDefault(item.EventType, "")
	item.CorrelationID = strings.TrimSpace(item.CorrelationID)
	item.TemplateCode = upperDefault(item.TemplateCode, "")
	item.RecipientType = upperDefault(item.RecipientType, "USER")
	item.RecipientID = strings.TrimSpace(item.RecipientID)
	item.RecipientAddress = strings.TrimSpace(item.RecipientAddress)
	item.Language = lowerDefault(item.Language, "vi")
	item.PayloadJSON = strings.TrimSpace(item.PayloadJSON)
	if item.PayloadJSON == "" {
		item.PayloadJSON = "{}"
	}
	item.Priority = intDefault(item.Priority, 100)
	item.Status = upperDefault(item.Status, "QUEUED")

	channels := make([]string, 0, len(item.Channels))
	seen := map[string]struct{}{}
	for _, channel := range item.Channels {
		channel = upperDefault(channel, "")
		if channel == "" {
			continue
		}
		if _, ok := seen[channel]; ok {
			continue
		}
		seen[channel] = struct{}{}
		channels = append(channels, channel)
	}
	if len(channels) == 0 {
		channels = []string{"IN_APP"}
	}
	item.Channels = channels

	item.IdempotencyKey = strings.TrimSpace(item.IdempotencyKey)
	if item.IdempotencyKey == "" {
		item.IdempotencyKey = fmt.Sprintf("%s|%s|%s|%s|%s|%s", item.SourceService, item.EventType, item.CorrelationID, item.RecipientType, item.RecipientID, item.TemplateCode)
	}
}

func intDefault(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}
