package service

import (
	stderrors "errors"
	"time"

	pb "github.com/arda-labs/arda/arda-be-go/services/notification-service/api/notification/v1"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	uc *biz.NotificationUsecase
}

func NewNotificationService(uc *biz.NotificationUsecase) *NotificationService {
	return &NotificationService{uc: uc}
}

func toServiceError(err error) error {
	if err == nil {
		return nil
	}
	if stderrors.Is(err, biz.ErrNotFound) {
		return kratoserrors.NotFound("NOTIFICATION_NOT_FOUND", "notification record not found")
	}
	if stderrors.Is(err, biz.ErrInvalidArgument) {
		return kratoserrors.BadRequest("NOTIFICATION_INVALID_ARGUMENT", "invalid notification request")
	}
	return err
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
