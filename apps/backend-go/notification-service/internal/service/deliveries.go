package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/notification-service/api/notification/v1"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
)

func (s *NotificationService) ListDeliveries(ctx context.Context, req *pb.ListDeliveriesRequest) (*pb.ListDeliveriesResponse, error) {
	list, next, err := s.uc.ListDeliveries(ctx, biz.DeliveryFilter{Status: req.Status, Channel: req.Channel, RecipientID: req.RecipientId, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListDeliveriesResponse{Deliveries: toProtoDeliveries(list), NextPageToken: next}, nil
}

func (s *NotificationService) RetryDelivery(ctx context.Context, req *pb.RetryDeliveryRequest) (*pb.NotificationDelivery, error) {
	item, err := s.uc.RetryDelivery(ctx, req.Id, req.Actor)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoDelivery(item), nil
}

func (s *NotificationService) RunDeliveryWorkerOnce(ctx context.Context, req *pb.RunDeliveryWorkerOnceRequest) (*pb.RunDeliveryWorkerOnceResponse, error) {
	processed, failed, err := s.uc.ProcessDueDeliveries(ctx, req.WorkerId, int(req.BatchSize))
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.RunDeliveryWorkerOnceResponse{Processed: int32(processed), Failed: int32(failed)}, nil
}

func toProtoDelivery(in *biz.NotificationDelivery) *pb.NotificationDelivery {
	if in == nil {
		return nil
	}
	return &pb.NotificationDelivery{
		Id:                   in.ID,
		RequestId:            in.RequestID,
		TemplateVersionId:    in.TemplateVersionID,
		Channel:              in.Channel,
		RecipientType:        in.RecipientType,
		RecipientId:          in.RecipientID,
		RecipientAddress:     in.RecipientAddress,
		Subject:              in.Subject,
		Body:                 in.Body,
		Status:               in.Status,
		AttemptCount:         int32(in.AttemptCount),
		MaxAttempts:          int32(in.MaxAttempts),
		NextAttemptAt:        toTimestamp(in.NextAttemptAt),
		LockedBy:             in.LockedBy,
		LockedAt:             toTimestamp(in.LockedAt),
		ProviderCode:         in.ProviderCode,
		ProviderMessageId:    in.ProviderMessageID,
		ProviderResponseJson: in.ProviderResponseJSON,
		ErrorMessage:         in.ErrorMessage,
		Priority:             int32(in.Priority),
		CreatedAt:            toTimestamp(in.CreatedAt),
		UpdatedAt:            toTimestamp(in.UpdatedAt),
	}
}

func toProtoDeliveries(in []*biz.NotificationDelivery) []*pb.NotificationDelivery {
	out := make([]*pb.NotificationDelivery, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoDelivery(item))
	}
	return out
}
