package service

import (
	stderrors "errors"
	"time"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MdmService struct {
	pb.UnimplementedMdmServiceServer
	uc *biz.MdmUsecase
}

func NewMdmService(uc *biz.MdmUsecase) *MdmService {
	return &MdmService{uc: uc}
}

func toServiceError(err error) error {
	if err == nil {
		return nil
	}
	if stderrors.Is(err, biz.ErrNotFound) {
		return kratoserrors.NotFound("MDM_NOT_FOUND", "MDM record not found")
	}
	if stderrors.Is(err, biz.ErrReadOnly) {
		return kratoserrors.Forbidden("MDM_READ_ONLY", "MDM record is read-only")
	}
	if stderrors.Is(err, biz.ErrInvalidArgument) {
		return kratoserrors.BadRequest("MDM_INVALID_ARGUMENT", "invalid MDM request")
	}
	return err
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
