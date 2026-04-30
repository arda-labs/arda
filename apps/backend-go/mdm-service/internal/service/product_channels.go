package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MdmService) ListBankingProducts(ctx context.Context, req *pb.ListBankingProductsRequest) (*pb.ListBankingProductsResponse, error) {
	list, next, err := s.uc.ListBankingProducts(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListBankingProductsResponse{BankingProducts: toProtoBankingProducts(list), NextPageToken: next}, nil
}

func (s *MdmService) CreateBankingProduct(ctx context.Context, req *pb.CreateBankingProductRequest) (*pb.BankingProduct, error) {
	item, err := s.uc.CreateBankingProduct(ctx, toBizBankingProduct(req.BankingProduct))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoBankingProduct(item), nil
}

func (s *MdmService) UpdateBankingProduct(ctx context.Context, req *pb.UpdateBankingProductRequest) (*pb.BankingProduct, error) {
	item := toBizBankingProduct(req.BankingProduct)
	item.ID = req.Id
	out, err := s.uc.UpdateBankingProduct(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoBankingProduct(out), nil
}

func (s *MdmService) DeleteBankingProduct(ctx context.Context, req *pb.DeleteBankingProductRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteBankingProduct(ctx, req.Id))
}

func (s *MdmService) ListServiceChannels(ctx context.Context, req *pb.ListServiceChannelsRequest) (*pb.ListServiceChannelsResponse, error) {
	list, next, err := s.uc.ListServiceChannels(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListServiceChannelsResponse{ServiceChannels: toProtoServiceChannels(list), NextPageToken: next}, nil
}

func (s *MdmService) CreateServiceChannel(ctx context.Context, req *pb.CreateServiceChannelRequest) (*pb.ServiceChannel, error) {
	item, err := s.uc.CreateServiceChannel(ctx, toBizServiceChannel(req.ServiceChannel))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoServiceChannel(item), nil
}

func (s *MdmService) UpdateServiceChannel(ctx context.Context, req *pb.UpdateServiceChannelRequest) (*pb.ServiceChannel, error) {
	item := toBizServiceChannel(req.ServiceChannel)
	item.ID = req.Id
	out, err := s.uc.UpdateServiceChannel(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoServiceChannel(out), nil
}

func (s *MdmService) DeleteServiceChannel(ctx context.Context, req *pb.DeleteServiceChannelRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteServiceChannel(ctx, req.Id))
}

func (s *MdmService) ListProductChannelRules(ctx context.Context, req *pb.ListProductChannelRulesRequest) (*pb.ListProductChannelRulesResponse, error) {
	list, next, err := s.uc.ListProductChannelRules(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListProductChannelRulesResponse{ProductChannelRules: toProtoProductChannelRules(list), NextPageToken: next}, nil
}

func (s *MdmService) CreateProductChannelRule(ctx context.Context, req *pb.CreateProductChannelRuleRequest) (*pb.ProductChannelRule, error) {
	item, err := s.uc.CreateProductChannelRule(ctx, toBizProductChannelRule(req.ProductChannelRule))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoProductChannelRule(item), nil
}

func (s *MdmService) UpdateProductChannelRule(ctx context.Context, req *pb.UpdateProductChannelRuleRequest) (*pb.ProductChannelRule, error) {
	item := toBizProductChannelRule(req.ProductChannelRule)
	item.ID = req.Id
	out, err := s.uc.UpdateProductChannelRule(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoProductChannelRule(out), nil
}

func (s *MdmService) DeleteProductChannelRule(ctx context.Context, req *pb.DeleteProductChannelRuleRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteProductChannelRule(ctx, req.Id))
}

func toBizBankingProduct(in *pb.BankingProduct) *biz.BankingProduct {
	if in == nil {
		return &biz.BankingProduct{}
	}
	return &biz.BankingProduct{ID: in.Id, Code: in.Code, Name: in.Name, ProductType: in.ProductType, Category: in.Category, CustomerSegment: in.CustomerSegment, Currency: in.Currency, EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo, Description: in.Description, Status: in.Status}
}

func toProtoBankingProduct(in *biz.BankingProduct) *pb.BankingProduct {
	if in == nil {
		return nil
	}
	return &pb.BankingProduct{Id: in.ID, Code: in.Code, Name: in.Name, ProductType: in.ProductType, Category: in.Category, CustomerSegment: in.CustomerSegment, Currency: in.Currency, EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo, Description: in.Description, Status: in.Status, CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt)}
}

func toProtoBankingProducts(in []*biz.BankingProduct) []*pb.BankingProduct {
	out := make([]*pb.BankingProduct, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoBankingProduct(item))
	}
	return out
}

func toBizServiceChannel(in *pb.ServiceChannel) *biz.ServiceChannel {
	if in == nil {
		return &biz.ServiceChannel{}
	}
	return &biz.ServiceChannel{ID: in.Id, Code: in.Code, Name: in.Name, ChannelType: in.ChannelType, Availability: in.Availability, Timezone: in.Timezone, Description: in.Description, Status: in.Status}
}

func toProtoServiceChannel(in *biz.ServiceChannel) *pb.ServiceChannel {
	if in == nil {
		return nil
	}
	return &pb.ServiceChannel{Id: in.ID, Code: in.Code, Name: in.Name, ChannelType: in.ChannelType, Availability: in.Availability, Timezone: in.Timezone, Description: in.Description, Status: in.Status, CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt)}
}

func toProtoServiceChannels(in []*biz.ServiceChannel) []*pb.ServiceChannel {
	out := make([]*pb.ServiceChannel, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoServiceChannel(item))
	}
	return out
}

func toBizProductChannelRule(in *pb.ProductChannelRule) *biz.ProductChannelRule {
	if in == nil {
		return &biz.ProductChannelRule{}
	}
	return &biz.ProductChannelRule{ID: in.Id, ProductCode: in.ProductCode, ChannelCode: in.ChannelCode, TransactionType: in.TransactionType, Enabled: in.Enabled, Priority: int(in.Priority), FeeScheduleCode: in.FeeScheduleCode, LimitProfileCode: in.LimitProfileCode, EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo, Description: in.Description, Status: in.Status}
}

func toProtoProductChannelRule(in *biz.ProductChannelRule) *pb.ProductChannelRule {
	if in == nil {
		return nil
	}
	return &pb.ProductChannelRule{Id: in.ID, ProductCode: in.ProductCode, ChannelCode: in.ChannelCode, TransactionType: in.TransactionType, Enabled: in.Enabled, Priority: int32(in.Priority), FeeScheduleCode: in.FeeScheduleCode, LimitProfileCode: in.LimitProfileCode, EffectiveFrom: in.EffectiveFrom, EffectiveTo: in.EffectiveTo, Description: in.Description, Status: in.Status, CreatedAt: timestamppb.New(in.CreatedAt), UpdatedAt: timestamppb.New(in.UpdatedAt)}
}

func toProtoProductChannelRules(in []*biz.ProductChannelRule) []*pb.ProductChannelRule {
	out := make([]*pb.ProductChannelRule, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoProductChannelRule(item))
	}
	return out
}
