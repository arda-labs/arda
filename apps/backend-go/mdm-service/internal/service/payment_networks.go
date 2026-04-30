package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
)

func (s *MdmService) ListBankBranches(ctx context.Context, req *pb.ListBankBranchesRequest) (*pb.ListBankBranchesResponse, error) {
	list, next, err := s.uc.ListBankBranches(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListBankBranchesResponse{BankBranches: toProtoBankBranches(list), NextPageToken: next}, nil
}

func (s *MdmService) CreateBankBranch(ctx context.Context, req *pb.CreateBankBranchRequest) (*pb.BankBranch, error) {
	item, err := s.uc.CreateBankBranch(ctx, toBizBankBranch(req.BankBranch))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoBankBranch(item), nil
}

func (s *MdmService) UpdateBankBranch(ctx context.Context, req *pb.UpdateBankBranchRequest) (*pb.BankBranch, error) {
	item := toBizBankBranch(req.BankBranch)
	item.ID = req.Id
	out, err := s.uc.UpdateBankBranch(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoBankBranch(out), nil
}

func (s *MdmService) DeleteBankBranch(ctx context.Context, req *pb.DeleteBankBranchRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteBankBranch(ctx, req.Id))
}

func (s *MdmService) ListPaymentNetworks(ctx context.Context, req *pb.ListPaymentNetworksRequest) (*pb.ListPaymentNetworksResponse, error) {
	list, next, err := s.uc.ListPaymentNetworks(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListPaymentNetworksResponse{PaymentNetworks: toProtoPaymentNetworks(list), NextPageToken: next}, nil
}

func (s *MdmService) CreatePaymentNetwork(ctx context.Context, req *pb.CreatePaymentNetworkRequest) (*pb.PaymentNetwork, error) {
	item, err := s.uc.CreatePaymentNetwork(ctx, toBizPaymentNetwork(req.PaymentNetwork))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoPaymentNetwork(item), nil
}

func (s *MdmService) UpdatePaymentNetwork(ctx context.Context, req *pb.UpdatePaymentNetworkRequest) (*pb.PaymentNetwork, error) {
	item := toBizPaymentNetwork(req.PaymentNetwork)
	item.ID = req.Id
	out, err := s.uc.UpdatePaymentNetwork(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoPaymentNetwork(out), nil
}

func (s *MdmService) DeletePaymentNetwork(ctx context.Context, req *pb.DeletePaymentNetworkRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeletePaymentNetwork(ctx, req.Id))
}

func toBizBankBranch(in *pb.BankBranch) *biz.BankBranch {
	if in == nil {
		return &biz.BankBranch{}
	}
	return &biz.BankBranch{ID: in.Id, InstitutionCode: in.InstitutionCode, Code: in.Code, Name: in.Name, BranchType: in.BranchType, Address: in.Address, ProvinceCode: in.ProvinceCode, Phone: in.Phone, SwiftCode: in.SwiftCode, NapasCode: in.NapasCode, Status: in.Status}
}

func toProtoBankBranch(in *biz.BankBranch) *pb.BankBranch {
	if in == nil {
		return nil
	}
	return &pb.BankBranch{Id: in.ID, InstitutionCode: in.InstitutionCode, Code: in.Code, Name: in.Name, BranchType: in.BranchType, Address: in.Address, ProvinceCode: in.ProvinceCode, Phone: in.Phone, SwiftCode: in.SwiftCode, NapasCode: in.NapasCode, Status: in.Status, CreatedAt: toTimestamp(in.CreatedAt), UpdatedAt: toTimestamp(in.UpdatedAt)}
}

func toProtoBankBranches(in []*biz.BankBranch) []*pb.BankBranch {
	out := make([]*pb.BankBranch, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoBankBranch(item))
	}
	return out
}

func toBizPaymentNetwork(in *pb.PaymentNetwork) *biz.PaymentNetwork {
	if in == nil {
		return &biz.PaymentNetwork{}
	}
	return &biz.PaymentNetwork{ID: in.Id, Code: in.Code, Name: in.Name, NetworkType: in.NetworkType, ClearingMethod: in.ClearingMethod, SettlementCurrency: in.SettlementCurrency, Operator: in.Operator, Availability: in.Availability, Description: in.Description, Status: in.Status}
}

func toProtoPaymentNetwork(in *biz.PaymentNetwork) *pb.PaymentNetwork {
	if in == nil {
		return nil
	}
	return &pb.PaymentNetwork{Id: in.ID, Code: in.Code, Name: in.Name, NetworkType: in.NetworkType, ClearingMethod: in.ClearingMethod, SettlementCurrency: in.SettlementCurrency, Operator: in.Operator, Availability: in.Availability, Description: in.Description, Status: in.Status, CreatedAt: toTimestamp(in.CreatedAt), UpdatedAt: toTimestamp(in.UpdatedAt)}
}

func toProtoPaymentNetworks(in []*biz.PaymentNetwork) []*pb.PaymentNetwork {
	out := make([]*pb.PaymentNetwork, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoPaymentNetwork(item))
	}
	return out
}
