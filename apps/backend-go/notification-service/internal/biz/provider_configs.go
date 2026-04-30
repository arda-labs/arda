package biz

import (
	"context"
	"encoding/json"
	"strings"
)

func (uc *NotificationUsecase) ListProviderConfigs(ctx context.Context, filter ProviderConfigFilter) ([]*ProviderConfig, error) {
	filter.Channel = upperDefault(filter.Channel, "")
	filter.Status = upperDefault(filter.Status, "")
	return uc.repo.ListProviderConfigs(ctx, filter)
}

func (uc *NotificationUsecase) UpsertProviderConfig(ctx context.Context, item *ProviderConfig) (*ProviderConfig, error) {
	normalizeProviderConfig(item)
	if item.Code == "" || item.Channel == "" || item.Name == "" {
		return nil, ErrInvalidArgument
	}
	if !json.Valid([]byte(item.OptionsJSON)) {
		return nil, ErrInvalidArgument
	}
	return uc.repo.UpsertProviderConfig(ctx, item)
}

func normalizeProviderConfig(item *ProviderConfig) {
	item.Code = upperDefault(item.Code, "")
	item.Channel = upperDefault(item.Channel, "")
	item.Name = strings.TrimSpace(item.Name)
	item.Priority = intDefault(item.Priority, 100)
	item.RateLimitPerMinute = intDefault(item.RateLimitPerMinute, 0)
	item.OptionsJSON = strings.TrimSpace(item.OptionsJSON)
	if item.OptionsJSON == "" {
		item.OptionsJSON = "{}"
	}
	item.Status = upperDefault(item.Status, "ACTIVE")
}
