package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/notification-service/api/notification/v1"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
)

func (s *NotificationService) CreateNotificationRequest(ctx context.Context, req *pb.CreateNotificationRequestRequest) (*pb.NotificationRequest, error) {
	item, err := s.uc.CreateNotificationRequest(ctx, toBizNotificationRequest(req.Request))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoNotificationRequest(item), nil
}

func (s *NotificationService) GetNotificationRequest(ctx context.Context, req *pb.GetNotificationRequestRequest) (*pb.NotificationRequest, error) {
	item, err := s.uc.GetNotificationRequest(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoNotificationRequest(item), nil
}

func toBizNotificationRequest(in *pb.NotificationRequest) *biz.NotificationRequest {
	if in == nil {
		return &biz.NotificationRequest{}
	}
	return &biz.NotificationRequest{
		ID:               in.Id,
		SourceService:    in.SourceService,
		EventType:        in.EventType,
		CorrelationID:    in.CorrelationId,
		IdempotencyKey:   in.IdempotencyKey,
		TemplateCode:     in.TemplateCode,
		RecipientType:    in.RecipientType,
		RecipientID:      in.RecipientId,
		RecipientAddress: in.RecipientAddress,
		Channels:         append([]string(nil), in.Channels...),
		Language:         in.Language,
		PayloadJSON:      in.PayloadJson,
		Priority:         int(in.Priority),
		Status:           in.Status,
	}
}

func toProtoNotificationRequest(in *biz.NotificationRequest) *pb.NotificationRequest {
	if in == nil {
		return nil
	}
	return &pb.NotificationRequest{
		Id:               in.ID,
		SourceService:    in.SourceService,
		EventType:        in.EventType,
		CorrelationId:    in.CorrelationID,
		IdempotencyKey:   in.IdempotencyKey,
		TemplateCode:     in.TemplateCode,
		RecipientType:    in.RecipientType,
		RecipientId:      in.RecipientID,
		RecipientAddress: in.RecipientAddress,
		Channels:         append([]string(nil), in.Channels...),
		Language:         in.Language,
		PayloadJson:      in.PayloadJSON,
		Priority:         int32(in.Priority),
		Status:           in.Status,
		CreatedAt:        toTimestamp(in.CreatedAt),
		UpdatedAt:        toTimestamp(in.UpdatedAt),
		Deliveries:       toProtoDeliveries(in.Deliveries),
	}
}
