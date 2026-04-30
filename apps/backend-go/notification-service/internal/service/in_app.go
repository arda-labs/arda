package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/notification-service/api/notification/v1"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
)

func (s *NotificationService) ListInAppNotifications(ctx context.Context, req *pb.ListInAppNotificationsRequest) (*pb.ListInAppNotificationsResponse, error) {
	list, next, err := s.uc.ListInAppNotifications(ctx, biz.InAppFilter{RecipientType: req.RecipientType, RecipientID: req.RecipientId, Status: req.Status, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListInAppNotificationsResponse{Notifications: toProtoInAppNotifications(list), NextPageToken: next}, nil
}

func (s *NotificationService) MarkInAppNotificationRead(ctx context.Context, req *pb.MarkInAppNotificationReadRequest) (*pb.InAppNotification, error) {
	item, err := s.uc.MarkInAppNotificationRead(ctx, req.Id, req.Actor)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoInAppNotification(item), nil
}

func (s *NotificationService) CountUnreadInAppNotifications(ctx context.Context, req *pb.CountUnreadInAppNotificationsRequest) (*pb.CountUnreadInAppNotificationsResponse, error) {
	count, err := s.uc.CountUnreadInAppNotifications(ctx, req.RecipientType, req.RecipientId)
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.CountUnreadInAppNotificationsResponse{Count: int32(count)}, nil
}

func (s *NotificationService) MarkAllInAppNotificationsRead(ctx context.Context, req *pb.MarkAllInAppNotificationsReadRequest) (*pb.MarkAllInAppNotificationsReadResponse, error) {
	updated, err := s.uc.MarkAllInAppNotificationsRead(ctx, req.RecipientType, req.RecipientId, req.Actor)
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.MarkAllInAppNotificationsReadResponse{Updated: int32(updated)}, nil
}

func toProtoInAppNotification(in *biz.InAppNotification) *pb.InAppNotification {
	if in == nil {
		return nil
	}
	return &pb.InAppNotification{
		Id:            in.ID,
		DeliveryId:    in.DeliveryID,
		RecipientType: in.RecipientType,
		RecipientId:   in.RecipientID,
		Title:         in.Title,
		Body:          in.Body,
		DataJson:      in.DataJSON,
		Status:        in.Status,
		ReadAt:        toTimestamp(in.ReadAt),
		CreatedAt:     toTimestamp(in.CreatedAt),
	}
}

func toProtoInAppNotifications(in []*biz.InAppNotification) []*pb.InAppNotification {
	out := make([]*pb.InAppNotification, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoInAppNotification(item))
	}
	return out
}
