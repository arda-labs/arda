package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MdmService) ListFeeSchedules(ctx context.Context, req *pb.ListFeeSchedulesRequest) (*pb.ListFeeSchedulesResponse, error) {
	list, next, err := s.uc.ListFeeSchedules(ctx, biz.PageFilter{
		Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListFeeSchedulesResponse{FeeSchedules: toProtoFeeSchedules(list), NextPageToken: next}, nil
}

func (s *MdmService) GetFeeSchedule(ctx context.Context, req *pb.GetFeeScheduleRequest) (*pb.FeeSchedule, error) {
	item, err := s.uc.GetFeeSchedule(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFeeSchedule(item), nil
}

func (s *MdmService) CreateFeeSchedule(ctx context.Context, req *pb.CreateFeeScheduleRequest) (*pb.FeeSchedule, error) {
	item, err := s.uc.CreateFeeSchedule(ctx, toBizFeeSchedule(req.FeeSchedule))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFeeSchedule(item), nil
}

func (s *MdmService) UpdateFeeSchedule(ctx context.Context, req *pb.UpdateFeeScheduleRequest) (*pb.FeeSchedule, error) {
	item := toBizFeeSchedule(req.FeeSchedule)
	item.ID = req.Id
	updated, err := s.uc.UpdateFeeSchedule(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFeeSchedule(updated), nil
}

func (s *MdmService) DeleteFeeSchedule(ctx context.Context, req *pb.DeleteFeeScheduleRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteFeeSchedule(ctx, req.Id))
}

func (s *MdmService) ApproveFeeSchedule(ctx context.Context, req *pb.ApprovePricingRuleRequest) (*pb.FeeSchedule, error) {
	item, err := s.uc.ApproveFeeSchedule(ctx, req.Id, req.Actor, req.Note)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFeeSchedule(item), nil
}

func (s *MdmService) ListTaxRules(ctx context.Context, req *pb.ListTaxRulesRequest) (*pb.ListTaxRulesResponse, error) {
	list, next, err := s.uc.ListTaxRules(ctx, biz.PageFilter{
		Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListTaxRulesResponse{TaxRules: toProtoTaxRules(list), NextPageToken: next}, nil
}

func (s *MdmService) GetTaxRule(ctx context.Context, req *pb.GetTaxRuleRequest) (*pb.TaxRule, error) {
	item, err := s.uc.GetTaxRule(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTaxRule(item), nil
}

func (s *MdmService) CreateTaxRule(ctx context.Context, req *pb.CreateTaxRuleRequest) (*pb.TaxRule, error) {
	item, err := s.uc.CreateTaxRule(ctx, toBizTaxRule(req.TaxRule))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTaxRule(item), nil
}

func (s *MdmService) UpdateTaxRule(ctx context.Context, req *pb.UpdateTaxRuleRequest) (*pb.TaxRule, error) {
	item := toBizTaxRule(req.TaxRule)
	item.ID = req.Id
	updated, err := s.uc.UpdateTaxRule(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTaxRule(updated), nil
}

func (s *MdmService) DeleteTaxRule(ctx context.Context, req *pb.DeleteTaxRuleRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteTaxRule(ctx, req.Id))
}

func (s *MdmService) ApproveTaxRule(ctx context.Context, req *pb.ApprovePricingRuleRequest) (*pb.TaxRule, error) {
	item, err := s.uc.ApproveTaxRule(ctx, req.Id, req.Actor, req.Note)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTaxRule(item), nil
}

func (s *MdmService) ListStandardLimits(ctx context.Context, req *pb.ListStandardLimitsRequest) (*pb.ListStandardLimitsResponse, error) {
	list, next, err := s.uc.ListStandardLimits(ctx, biz.PageFilter{
		Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListStandardLimitsResponse{StandardLimits: toProtoStandardLimits(list), NextPageToken: next}, nil
}

func (s *MdmService) GetStandardLimit(ctx context.Context, req *pb.GetStandardLimitRequest) (*pb.StandardLimit, error) {
	item, err := s.uc.GetStandardLimit(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoStandardLimit(item), nil
}

func (s *MdmService) CreateStandardLimit(ctx context.Context, req *pb.CreateStandardLimitRequest) (*pb.StandardLimit, error) {
	item, err := s.uc.CreateStandardLimit(ctx, toBizStandardLimit(req.StandardLimit))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoStandardLimit(item), nil
}

func (s *MdmService) UpdateStandardLimit(ctx context.Context, req *pb.UpdateStandardLimitRequest) (*pb.StandardLimit, error) {
	item := toBizStandardLimit(req.StandardLimit)
	item.ID = req.Id
	updated, err := s.uc.UpdateStandardLimit(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoStandardLimit(updated), nil
}

func (s *MdmService) DeleteStandardLimit(ctx context.Context, req *pb.DeleteStandardLimitRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteStandardLimit(ctx, req.Id))
}

func (s *MdmService) ApproveStandardLimit(ctx context.Context, req *pb.ApprovePricingRuleRequest) (*pb.StandardLimit, error) {
	item, err := s.uc.ApproveStandardLimit(ctx, req.Id, req.Actor, req.Note)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoStandardLimit(item), nil
}

func toBizFeeSchedule(in *pb.FeeSchedule) *biz.FeeSchedule {
	if in == nil {
		return &biz.FeeSchedule{}
	}
	return &biz.FeeSchedule{
		ID: in.Id, Code: in.Code, Name: in.Name, FeeType: in.FeeType,
		CalculationMethod: in.CalculationMethod, Currency: in.Currency,
		FixedAmount: in.FixedAmount, RatePercent: in.RatePercent, MinAmount: in.MinAmount,
		MaxAmount: in.MaxAmount, Channel: in.Channel, ProductCode: in.ProductCode,
		EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo,
		Description: in.Description, Status: in.Status,
		ApprovalStatus: in.ApprovalStatus, Version: int(in.Version),
		ApprovedBy: in.ApprovedBy, ChangeNote: in.ChangeNote,
	}
}

func toProtoFeeSchedule(in *biz.FeeSchedule) *pb.FeeSchedule {
	if in == nil {
		return nil
	}
	return &pb.FeeSchedule{
		Id: in.ID, Code: in.Code, Name: in.Name, FeeType: in.FeeType,
		CalculationMethod: in.CalculationMethod, Currency: in.Currency,
		FixedAmount: in.FixedAmount, RatePercent: in.RatePercent, MinAmount: in.MinAmount,
		MaxAmount: in.MaxAmount, Channel: in.Channel, ProductCode: in.ProductCode,
		EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo,
		Description: in.Description, Status: in.Status,
		ApprovalStatus: in.ApprovalStatus, Version: int32(in.Version),
		ApprovedBy: in.ApprovedBy, ApprovedAt: toTimestamp(in.ApprovedAt), ChangeNote: in.ChangeNote,
		CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt),
	}
}

func toProtoFeeSchedules(in []*biz.FeeSchedule) []*pb.FeeSchedule {
	out := make([]*pb.FeeSchedule, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoFeeSchedule(item))
	}
	return out
}

func toBizTaxRule(in *pb.TaxRule) *biz.TaxRule {
	if in == nil {
		return &biz.TaxRule{}
	}
	return &biz.TaxRule{
		ID: in.Id, Code: in.Code, Name: in.Name, TaxType: in.TaxType,
		RatePercent: in.RatePercent, Inclusive: in.Inclusive, Jurisdiction: in.Jurisdiction,
		EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo,
		Description: in.Description, Status: in.Status,
		ApprovalStatus: in.ApprovalStatus, Version: int(in.Version),
		ApprovedBy: in.ApprovedBy, ChangeNote: in.ChangeNote,
	}
}

func toProtoTaxRule(in *biz.TaxRule) *pb.TaxRule {
	if in == nil {
		return nil
	}
	return &pb.TaxRule{
		Id: in.ID, Code: in.Code, Name: in.Name, TaxType: in.TaxType,
		RatePercent: in.RatePercent, Inclusive: in.Inclusive, Jurisdiction: in.Jurisdiction,
		EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo,
		Description: in.Description, Status: in.Status,
		ApprovalStatus: in.ApprovalStatus, Version: int32(in.Version),
		ApprovedBy: in.ApprovedBy, ApprovedAt: toTimestamp(in.ApprovedAt), ChangeNote: in.ChangeNote,
		CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt),
	}
}

func toProtoTaxRules(in []*biz.TaxRule) []*pb.TaxRule {
	out := make([]*pb.TaxRule, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoTaxRule(item))
	}
	return out
}

func toBizStandardLimit(in *pb.StandardLimit) *biz.StandardLimit {
	if in == nil {
		return &biz.StandardLimit{}
	}
	return &biz.StandardLimit{
		ID: in.Id, Code: in.Code, Name: in.Name, LimitType: in.LimitType,
		SubjectType: in.SubjectType, Currency: in.Currency, MinAmount: in.MinAmount,
		PerTxnAmount: in.PerTxnAmount, DailyAmount: in.DailyAmount, MonthlyAmount: in.MonthlyAmount,
		CountLimit: int(in.CountLimit), Channel: in.Channel, ProductCode: in.ProductCode,
		EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo,
		Description: in.Description, Status: in.Status,
		ApprovalStatus: in.ApprovalStatus, Version: int(in.Version),
		ApprovedBy: in.ApprovedBy, ChangeNote: in.ChangeNote,
	}
}

func toProtoStandardLimit(in *biz.StandardLimit) *pb.StandardLimit {
	if in == nil {
		return nil
	}
	return &pb.StandardLimit{
		Id: in.ID, Code: in.Code, Name: in.Name, LimitType: in.LimitType,
		SubjectType: in.SubjectType, Currency: in.Currency, MinAmount: in.MinAmount,
		PerTxnAmount: in.PerTxnAmount, DailyAmount: in.DailyAmount, MonthlyAmount: in.MonthlyAmount,
		CountLimit: int32(in.CountLimit), Channel: in.Channel, ProductCode: in.ProductCode,
		EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo,
		Description: in.Description, Status: in.Status,
		ApprovalStatus: in.ApprovalStatus, Version: int32(in.Version),
		ApprovedBy: in.ApprovedBy, ApprovedAt: toTimestamp(in.ApprovedAt), ChangeNote: in.ChangeNote,
		CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt),
	}
}

func toProtoStandardLimits(in []*biz.StandardLimit) []*pb.StandardLimit {
	out := make([]*pb.StandardLimit, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoStandardLimit(item))
	}
	return out
}
