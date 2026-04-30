package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
)

func (s *MdmService) ListCreditInstitutions(ctx context.Context, req *pb.ListCreditInstitutionsRequest) (*pb.ListCreditInstitutionsResponse, error) {
	list, next, err := s.uc.ListCreditInstitutions(ctx, biz.PageFilter{
		Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListCreditInstitutionsResponse{CreditInstitutions: toProtoCreditInstitutions(list), NextPageToken: next}, nil
}

func (s *MdmService) GetCreditInstitution(ctx context.Context, req *pb.GetCreditInstitutionRequest) (*pb.CreditInstitution, error) {
	item, err := s.uc.GetCreditInstitution(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCreditInstitution(item), nil
}

func (s *MdmService) CreateCreditInstitution(ctx context.Context, req *pb.CreateCreditInstitutionRequest) (*pb.CreditInstitution, error) {
	item, err := s.uc.CreateCreditInstitution(ctx, toBizCreditInstitution(req.CreditInstitution))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCreditInstitution(item), nil
}

func (s *MdmService) UpdateCreditInstitution(ctx context.Context, req *pb.UpdateCreditInstitutionRequest) (*pb.CreditInstitution, error) {
	item := toBizCreditInstitution(req.CreditInstitution)
	item.ID = req.Id
	updated, err := s.uc.UpdateCreditInstitution(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCreditInstitution(updated), nil
}

func (s *MdmService) DeleteCreditInstitution(ctx context.Context, req *pb.DeleteCreditInstitutionRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteCreditInstitution(ctx, req.Id))
}

func toBizCreditInstitution(in *pb.CreditInstitution) *biz.CreditInstitution {
	if in == nil {
		return &biz.CreditInstitution{}
	}
	return &biz.CreditInstitution{
		ID:            in.Id,
		Code:          in.Code,
		Name:          in.Name,
		ShortName:     in.ShortName,
		Address:       in.Address,
		Phone:         in.Phone,
		Email:         in.Email,
		LicenseNumber: in.LicenseNumber,
		IssuedDate:    in.IssuedDate,
		TaxCode:       in.TaxCode,
		Website:       in.Website,
		Note:          in.Note,
		Status:        in.Status,
	}
}

func toProtoCreditInstitution(in *biz.CreditInstitution) *pb.CreditInstitution {
	if in == nil {
		return nil
	}
	return &pb.CreditInstitution{
		Id:            in.ID,
		Code:          in.Code,
		Name:          in.Name,
		ShortName:     in.ShortName,
		Address:       in.Address,
		Phone:         in.Phone,
		Email:         in.Email,
		LicenseNumber: in.LicenseNumber,
		IssuedDate:    in.IssuedDate,
		TaxCode:       in.TaxCode,
		Website:       in.Website,
		Note:          in.Note,
		Status:        in.Status,
		CreatedAt:     toTimestamp(in.CreatedAt),
		UpdatedAt:     toTimestamp(in.UpdatedAt),
	}
}

func toProtoCreditInstitutions(in []*biz.CreditInstitution) []*pb.CreditInstitution {
	out := make([]*pb.CreditInstitution, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoCreditInstitution(item))
	}
	return out
}
