package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MdmService) ListBusinessCalendars(ctx context.Context, req *pb.ListBusinessCalendarsRequest) (*pb.ListBusinessCalendarsResponse, error) {
	list, next, err := s.uc.ListBusinessCalendars(ctx, biz.PageFilter{
		Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListBusinessCalendarsResponse{BusinessCalendars: toProtoBusinessCalendars(list), NextPageToken: next}, nil
}

func (s *MdmService) GetBusinessCalendar(ctx context.Context, req *pb.GetBusinessCalendarRequest) (*pb.BusinessCalendar, error) {
	item, err := s.uc.GetBusinessCalendar(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoBusinessCalendar(item), nil
}

func (s *MdmService) CreateBusinessCalendar(ctx context.Context, req *pb.CreateBusinessCalendarRequest) (*pb.BusinessCalendar, error) {
	item, err := s.uc.CreateBusinessCalendar(ctx, toBizBusinessCalendar(req.BusinessCalendar))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoBusinessCalendar(item), nil
}

func (s *MdmService) UpdateBusinessCalendar(ctx context.Context, req *pb.UpdateBusinessCalendarRequest) (*pb.BusinessCalendar, error) {
	item := toBizBusinessCalendar(req.BusinessCalendar)
	item.ID = req.Id
	updated, err := s.uc.UpdateBusinessCalendar(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoBusinessCalendar(updated), nil
}

func (s *MdmService) DeleteBusinessCalendar(ctx context.Context, req *pb.DeleteBusinessCalendarRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteBusinessCalendar(ctx, req.Id))
}

func (s *MdmService) ListWorkingHours(ctx context.Context, req *pb.ListWorkingHoursRequest) (*pb.ListWorkingHoursResponse, error) {
	list, err := s.uc.ListWorkingHours(ctx, req.CalendarId)
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListWorkingHoursResponse{WorkingHours: toProtoWorkingHours(list)}, nil
}

func (s *MdmService) CreateWorkingHour(ctx context.Context, req *pb.CreateWorkingHourRequest) (*pb.WorkingHour, error) {
	item := toBizWorkingHour(req.WorkingHour)
	item.CalendarID = req.CalendarId
	created, err := s.uc.CreateWorkingHour(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoWorkingHour(created), nil
}

func (s *MdmService) UpdateWorkingHour(ctx context.Context, req *pb.UpdateWorkingHourRequest) (*pb.WorkingHour, error) {
	item := toBizWorkingHour(req.WorkingHour)
	item.ID = req.Id
	updated, err := s.uc.UpdateWorkingHour(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoWorkingHour(updated), nil
}

func (s *MdmService) DeleteWorkingHour(ctx context.Context, req *pb.DeleteWorkingHourRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteWorkingHour(ctx, req.Id))
}

func (s *MdmService) ListCalendarExceptions(ctx context.Context, req *pb.ListCalendarExceptionsRequest) (*pb.ListCalendarExceptionsResponse, error) {
	list, err := s.uc.ListCalendarExceptions(ctx, biz.CalendarExceptionFilter{
		CalendarID: req.CalendarId, FromDate: req.FromDate, ToDate: req.ToDate, ExceptionType: req.ExceptionType,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListCalendarExceptionsResponse{CalendarExceptions: toProtoCalendarExceptions(list)}, nil
}

func (s *MdmService) CreateCalendarException(ctx context.Context, req *pb.CreateCalendarExceptionRequest) (*pb.CalendarException, error) {
	item := toBizCalendarException(req.CalendarException)
	item.CalendarID = req.CalendarId
	created, err := s.uc.CreateCalendarException(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCalendarException(created), nil
}

func (s *MdmService) UpdateCalendarException(ctx context.Context, req *pb.UpdateCalendarExceptionRequest) (*pb.CalendarException, error) {
	item := toBizCalendarException(req.CalendarException)
	item.ID = req.Id
	updated, err := s.uc.UpdateCalendarException(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCalendarException(updated), nil
}

func (s *MdmService) DeleteCalendarException(ctx context.Context, req *pb.DeleteCalendarExceptionRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteCalendarException(ctx, req.Id))
}

func (s *MdmService) CalculateBusinessDay(ctx context.Context, req *pb.CalculateBusinessDayRequest) (*pb.CalculateBusinessDayResponse, error) {
	result, err := s.uc.CalculateBusinessDay(ctx, biz.BusinessDayCalculation{
		CalendarID: req.CalendarId, CalendarCode: req.CalendarCode, StartDate: req.StartDate,
		OffsetDays: int(req.OffsetDays), AdjustmentRule: req.AdjustmentRule,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.CalculateBusinessDayResponse{
		StartDate: result.StartDate, ResultDate: result.ResultDate, CalendarDays: int32(result.CalendarDays),
		SkippedDates: result.SkippedDates, IsBusinessDay: result.IsBusinessDay,
	}, nil
}

func toBizBusinessCalendar(in *pb.BusinessCalendar) *biz.BusinessCalendar {
	if in == nil {
		return &biz.BusinessCalendar{}
	}
	return &biz.BusinessCalendar{
		ID:           in.Id,
		Code:         in.Code,
		Name:         in.Name,
		Timezone:     in.Timezone,
		CalendarType: in.CalendarType,
		Description:  in.Description,
		Status:       in.Status,
	}
}

func toProtoBusinessCalendar(in *biz.BusinessCalendar) *pb.BusinessCalendar {
	if in == nil {
		return nil
	}
	return &pb.BusinessCalendar{
		Id:           in.ID,
		Code:         in.Code,
		Name:         in.Name,
		Timezone:     in.Timezone,
		CalendarType: in.CalendarType,
		Description:  in.Description,
		Status:       in.Status,
		CreatedAt:    timestamppb.New(in.CreatedAt),
		UpdatedAt:    timestamppb.New(in.UpdatedAt),
	}
}

func toProtoBusinessCalendars(in []*biz.BusinessCalendar) []*pb.BusinessCalendar {
	out := make([]*pb.BusinessCalendar, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoBusinessCalendar(item))
	}
	return out
}

func toBizWorkingHour(in *pb.WorkingHour) *biz.WorkingHour {
	if in == nil {
		return &biz.WorkingHour{}
	}
	return &biz.WorkingHour{
		ID:           in.Id,
		CalendarID:   in.CalendarId,
		DayOfWeek:    int(in.DayOfWeek),
		IsWorkingDay: in.IsWorkingDay,
		StartTime:    in.StartTime,
		EndTime:      in.EndTime,
		CutoffTime:   in.CutoffTime,
		SessionName:  in.SessionName,
		SortOrder:    int(in.SortOrder),
	}
}

func toProtoWorkingHour(in *biz.WorkingHour) *pb.WorkingHour {
	if in == nil {
		return nil
	}
	return &pb.WorkingHour{
		Id:           in.ID,
		CalendarId:   in.CalendarID,
		DayOfWeek:    int32(in.DayOfWeek),
		IsWorkingDay: in.IsWorkingDay,
		StartTime:    in.StartTime,
		EndTime:      in.EndTime,
		CutoffTime:   in.CutoffTime,
		SessionName:  in.SessionName,
		SortOrder:    int32(in.SortOrder),
		CreatedAt:    timestamppb.New(in.CreatedAt),
		UpdatedAt:    timestamppb.New(in.UpdatedAt),
	}
}

func toProtoWorkingHours(in []*biz.WorkingHour) []*pb.WorkingHour {
	out := make([]*pb.WorkingHour, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoWorkingHour(item))
	}
	return out
}

func toBizCalendarException(in *pb.CalendarException) *biz.CalendarException {
	if in == nil {
		return &biz.CalendarException{}
	}
	return &biz.CalendarException{
		ID:            in.Id,
		CalendarID:    in.CalendarId,
		Date:          in.Date,
		ExceptionType: in.ExceptionType,
		Name:          in.Name,
		IsWorkingDay:  in.IsWorkingDay,
		StartTime:     in.StartTime,
		EndTime:       in.EndTime,
		CutoffTime:    in.CutoffTime,
		Source:        in.Source,
		Note:          in.Note,
	}
}

func toProtoCalendarException(in *biz.CalendarException) *pb.CalendarException {
	if in == nil {
		return nil
	}
	return &pb.CalendarException{
		Id:            in.ID,
		CalendarId:    in.CalendarID,
		Date:          in.Date,
		ExceptionType: in.ExceptionType,
		Name:          in.Name,
		IsWorkingDay:  in.IsWorkingDay,
		StartTime:     in.StartTime,
		EndTime:       in.EndTime,
		CutoffTime:    in.CutoffTime,
		Source:        in.Source,
		Note:          in.Note,
		CreatedAt:     timestamppb.New(in.CreatedAt),
		UpdatedAt:     timestamppb.New(in.UpdatedAt),
	}
}

func toProtoCalendarExceptions(in []*biz.CalendarException) []*pb.CalendarException {
	out := make([]*pb.CalendarException, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoCalendarException(item))
	}
	return out
}
