package service

import (
	"context"
	"time"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MdmService) ListCurrencies(ctx context.Context, req *pb.ListCurrenciesRequest) (*pb.ListCurrenciesResponse, error) {
	list, next, err := s.uc.ListCurrencies(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListCurrenciesResponse{Currencies: toProtoCurrencies(list), NextPageToken: next}, nil
}
func (s *MdmService) CreateCurrency(ctx context.Context, req *pb.CreateCurrencyRequest) (*pb.Currency, error) {
	item, err := s.uc.CreateCurrency(ctx, toBizCurrency(req.Currency))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCurrency(item), nil
}
func (s *MdmService) UpdateCurrency(ctx context.Context, req *pb.UpdateCurrencyRequest) (*pb.Currency, error) {
	item := toBizCurrency(req.Currency)
	item.ID = req.Id
	out, err := s.uc.UpdateCurrency(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCurrency(out), nil
}
func (s *MdmService) DeleteCurrency(ctx context.Context, req *pb.DeleteCurrencyRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteCurrency(ctx, req.Id))
}

func (s *MdmService) ListFxRateSources(ctx context.Context, req *pb.ListFxRateSourcesRequest) (*pb.ListFxRateSourcesResponse, error) {
	list, next, err := s.uc.ListFxRateSources(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListFxRateSourcesResponse{Sources: toProtoFxRateSources(list), NextPageToken: next}, nil
}
func (s *MdmService) CreateFxRateSource(ctx context.Context, req *pb.CreateFxRateSourceRequest) (*pb.FxRateSource, error) {
	item, err := s.uc.CreateFxRateSource(ctx, toBizFxRateSource(req.Source))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFxRateSource(item), nil
}
func (s *MdmService) UpdateFxRateSource(ctx context.Context, req *pb.UpdateFxRateSourceRequest) (*pb.FxRateSource, error) {
	item := toBizFxRateSource(req.Source)
	item.ID = req.Id
	out, err := s.uc.UpdateFxRateSource(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFxRateSource(out), nil
}
func (s *MdmService) DeleteFxRateSource(ctx context.Context, req *pb.DeleteFxRateSourceRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteFxRateSource(ctx, req.Id))
}

func (s *MdmService) ListFxRates(ctx context.Context, req *pb.ListFxRatesRequest) (*pb.ListFxRatesResponse, error) {
	list, next, err := s.uc.ListFxRates(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListFxRatesResponse{FxRates: toProtoFxRates(list), NextPageToken: next}, nil
}
func (s *MdmService) CreateFxRate(ctx context.Context, req *pb.CreateFxRateRequest) (*pb.FxRate, error) {
	item, err := s.uc.CreateFxRate(ctx, toBizFxRate(req.FxRate))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFxRate(item), nil
}
func (s *MdmService) UpdateFxRate(ctx context.Context, req *pb.UpdateFxRateRequest) (*pb.FxRate, error) {
	item := toBizFxRate(req.FxRate)
	item.ID = req.Id
	out, err := s.uc.UpdateFxRate(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFxRate(out), nil
}
func (s *MdmService) DeleteFxRate(ctx context.Context, req *pb.DeleteFxRateRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteFxRate(ctx, req.Id))
}
func (s *MdmService) ApproveFxRate(ctx context.Context, req *pb.ApprovePricingRuleRequest) (*pb.FxRate, error) {
	item, err := s.uc.ApproveFxRate(ctx, req.Id, req.Actor, req.Note)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoFxRate(item), nil
}

func toBizCurrency(in *pb.Currency) *biz.Currency {
	if in == nil {
		return &biz.Currency{}
	}
	return &biz.Currency{ID: in.Id, Code: in.Code, NumericCode: in.NumericCode, Name: in.Name, MinorUnit: int(in.MinorUnit), Symbol: in.Symbol, CountryCode: in.CountryCode, Status: in.Status}
}
func toProtoCurrency(in *biz.Currency) *pb.Currency {
	if in == nil {
		return nil
	}
	return &pb.Currency{Id: in.ID, Code: in.Code, NumericCode: in.NumericCode, Name: in.Name, MinorUnit: int32(in.MinorUnit), Symbol: in.Symbol, CountryCode: in.CountryCode, Status: in.Status, CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt)}
}
func toProtoCurrencies(in []*biz.Currency) []*pb.Currency {
	out := make([]*pb.Currency, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoCurrency(item))
	}
	return out
}

func toBizFxRateSource(in *pb.FxRateSource) *biz.FxRateSource {
	if in == nil {
		return &biz.FxRateSource{}
	}
	return &biz.FxRateSource{ID: in.Id, Code: in.Code, Name: in.Name, SourceType: in.SourceType, Priority: int(in.Priority), Timezone: in.Timezone, Description: in.Description, Status: in.Status}
}
func toProtoFxRateSource(in *biz.FxRateSource) *pb.FxRateSource {
	if in == nil {
		return nil
	}
	return &pb.FxRateSource{Id: in.ID, Code: in.Code, Name: in.Name, SourceType: in.SourceType, Priority: int32(in.Priority), Timezone: in.Timezone, Description: in.Description, Status: in.Status, CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt)}
}
func toProtoFxRateSources(in []*biz.FxRateSource) []*pb.FxRateSource {
	out := make([]*pb.FxRateSource, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoFxRateSource(item))
	}
	return out
}

func toBizFxRate(in *pb.FxRate) *biz.FxRate {
	if in == nil {
		return &biz.FxRate{}
	}
	return &biz.FxRate{ID: in.Id, BaseCurrency: in.BaseCurrency, QuoteCurrency: in.QuoteCurrency, SourceCode: in.SourceCode, RateDate: in.RateDate, EffectiveAt: timestampTime(in.EffectiveAt), BuyRate: in.BuyRate, SellRate: in.SellRate, MidRate: in.MidRate, ApprovalStatus: in.ApprovalStatus, Version: int(in.Version), ApprovedBy: in.ApprovedBy, ChangeNote: in.ChangeNote, Status: in.Status}
}
func toProtoFxRate(in *biz.FxRate) *pb.FxRate {
	if in == nil {
		return nil
	}
	return &pb.FxRate{Id: in.ID, BaseCurrency: in.BaseCurrency, QuoteCurrency: in.QuoteCurrency, SourceCode: in.SourceCode, RateDate: in.RateDate, EffectiveAt: toTimestamp(in.EffectiveAt), BuyRate: in.BuyRate, SellRate: in.SellRate, MidRate: in.MidRate, ApprovalStatus: in.ApprovalStatus, Version: int32(in.Version), ApprovedBy: in.ApprovedBy, ApprovedAt: toTimestamp(in.ApprovedAt), ChangeNote: in.ChangeNote, Status: in.Status, CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt)}
}
func toProtoFxRates(in []*biz.FxRate) []*pb.FxRate {
	out := make([]*pb.FxRate, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoFxRate(item))
	}
	return out
}

func timestampTime(in *timestamppb.Timestamp) time.Time {
	if in == nil {
		return time.Time{}
	}
	return in.AsTime()
}
