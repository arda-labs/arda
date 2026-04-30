package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/notification-service/api/notification/v1"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
)

func (s *NotificationService) ListProviderConfigs(ctx context.Context, req *pb.ListProviderConfigsRequest) (*pb.ListProviderConfigsResponse, error) {
	list, err := s.uc.ListProviderConfigs(ctx, biz.ProviderConfigFilter{Channel: req.Channel, Status: req.Status})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListProviderConfigsResponse{ProviderConfigs: toProtoProviderConfigs(list)}, nil
}

func (s *NotificationService) UpsertProviderConfig(ctx context.Context, req *pb.UpsertProviderConfigRequest) (*pb.ProviderConfig, error) {
	item := toBizProviderConfig(req.ProviderConfig)
	item.Code = req.Code
	out, err := s.uc.UpsertProviderConfig(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoProviderConfig(out), nil
}

func toBizProviderConfig(in *pb.ProviderConfig) *biz.ProviderConfig {
	if in == nil {
		return &biz.ProviderConfig{}
	}
	return &biz.ProviderConfig{
		ID:                 in.Id,
		Code:               in.Code,
		Channel:            in.Channel,
		Name:               in.Name,
		Priority:           int(in.Priority),
		RateLimitPerMinute: int(in.RateLimitPerMinute),
		OptionsJSON:        in.OptionsJson,
		Status:             in.Status,
	}
}

func toProtoProviderConfig(in *biz.ProviderConfig) *pb.ProviderConfig {
	if in == nil {
		return nil
	}
	return &pb.ProviderConfig{
		Id:                 in.ID,
		Code:               in.Code,
		Channel:            in.Channel,
		Name:               in.Name,
		Priority:           int32(in.Priority),
		RateLimitPerMinute: int32(in.RateLimitPerMinute),
		OptionsJson:        in.OptionsJSON,
		Status:             in.Status,
		CreatedAt:          toTimestamp(in.CreatedAt),
		UpdatedAt:          toTimestamp(in.UpdatedAt),
	}
}

func toProtoProviderConfigs(in []*biz.ProviderConfig) []*pb.ProviderConfig {
	out := make([]*pb.ProviderConfig, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoProviderConfig(item))
	}
	return out
}
