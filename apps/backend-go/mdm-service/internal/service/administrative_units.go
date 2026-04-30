package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
)

func (s *MdmService) ListAdministrativeUnits(ctx context.Context, req *pb.ListAdministrativeUnitsRequest) (*pb.ListAdministrativeUnitsResponse, error) {
	list, next, err := s.uc.ListAdministrativeUnits(ctx, administrativeUnitFilter(req))
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListAdministrativeUnitsResponse{Units: toProtoAdministrativeUnits(list), NextPageToken: next}, nil
}

func (s *MdmService) ListAdministrativeUnitTree(ctx context.Context, req *pb.ListAdministrativeUnitsRequest) (*pb.ListAdministrativeUnitTreeResponse, error) {
	nodes, err := s.uc.ListAdministrativeUnitTree(ctx, administrativeUnitFilter(req))
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListAdministrativeUnitTreeResponse{Nodes: toProtoAdministrativeUnitNodes(nodes)}, nil
}

func (s *MdmService) ListProvinces(ctx context.Context, req *pb.ListAdministrativeUnitsRequest) (*pb.ListAdministrativeUnitsResponse, error) {
	list, next, err := s.uc.ListProvinces(ctx, administrativeUnitFilter(req))
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListAdministrativeUnitsResponse{Units: toProtoAdministrativeUnits(list), NextPageToken: next}, nil
}

func (s *MdmService) ListWards(ctx context.Context, req *pb.ListWardsRequest) (*pb.ListAdministrativeUnitsResponse, error) {
	list, next, err := s.uc.ListWards(ctx, req.ProvinceId, biz.PageFilter{
		Status:    req.Status,
		Keyword:   req.Keyword,
		PageSize:  int(req.PageSize),
		PageToken: req.PageToken,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListAdministrativeUnitsResponse{Units: toProtoAdministrativeUnits(list), NextPageToken: next}, nil
}

func (s *MdmService) GetAdministrativeUnit(ctx context.Context, req *pb.GetAdministrativeUnitRequest) (*pb.AdministrativeUnit, error) {
	unit, err := s.uc.GetAdministrativeUnit(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoAdministrativeUnit(unit), nil
}

func (s *MdmService) CreateAdministrativeUnit(ctx context.Context, req *pb.CreateAdministrativeUnitRequest) (*pb.AdministrativeUnit, error) {
	unit, err := s.uc.CreateAdministrativeUnit(ctx, toBizAdministrativeUnit(req.Unit))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoAdministrativeUnit(unit), nil
}

func (s *MdmService) UpdateAdministrativeUnit(ctx context.Context, req *pb.UpdateAdministrativeUnitRequest) (*pb.AdministrativeUnit, error) {
	unit := toBizAdministrativeUnit(req.Unit)
	unit.ID = req.Id
	updated, err := s.uc.UpdateAdministrativeUnit(ctx, unit)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoAdministrativeUnit(updated), nil
}

func (s *MdmService) DeleteAdministrativeUnit(ctx context.Context, req *pb.DeleteAdministrativeUnitRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteAdministrativeUnit(ctx, req.Id))
}

func (s *MdmService) SyncAdministrativeUnitsFromAddressKit(ctx context.Context, req *pb.SyncAdministrativeUnitsFromAddressKitRequest) (*pb.SyncAdministrativeUnitsFromAddressKitResponse, error) {
	result, err := s.uc.SyncAdministrativeUnitsFromAddressKit(ctx)
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.SyncAdministrativeUnitsFromAddressKitResponse{
		ProvinceCount: int32(result.ProvinceCount),
		WardCount:     int32(result.WardCount),
		EffectiveDate: result.EffectiveDate,
		Source:        result.Source,
	}, nil
}

func administrativeUnitFilter(req *pb.ListAdministrativeUnitsRequest) biz.AdministrativeUnitFilter {
	return biz.AdministrativeUnitFilter{
		ParentID: req.ParentId,
		Level:    req.Level,
		PageFilter: biz.PageFilter{
			Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
		},
	}
}

func areaFilter(req *pb.ListAreasRequest) biz.AreaFilter {
	return biz.AreaFilter{
		AreaTypeID: req.AreaTypeId,
		ParentID:   req.ParentId,
		PageFilter: biz.PageFilter{
			Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
		},
	}
}

func toBizAdministrativeUnit(in *pb.AdministrativeUnit) *biz.AdministrativeUnit {
	if in == nil {
		return &biz.AdministrativeUnit{}
	}
	return &biz.AdministrativeUnit{
		ID:            in.Id,
		Code:          in.Code,
		Name:          in.Name,
		FullName:      in.FullName,
		ShortName:     in.ShortName,
		Level:         in.Level,
		UnitType:      in.UnitType,
		ParentID:      in.ParentId,
		Path:          in.Path,
		SortOrder:     int(in.SortOrder),
		Latitude:      in.Latitude,
		Longitude:     in.Longitude,
		Status:        in.Status,
		EffectiveFrom: in.EffectiveFrom,
		EffectiveTo:   in.EffectiveTo,
		Source:        in.Source,
		MetadataJSON:  in.MetadataJson,
	}
}

func toProtoAdministrativeUnit(in *biz.AdministrativeUnit) *pb.AdministrativeUnit {
	if in == nil {
		return nil
	}
	return &pb.AdministrativeUnit{
		Id:            in.ID,
		Code:          in.Code,
		Name:          in.Name,
		FullName:      in.FullName,
		ShortName:     in.ShortName,
		Level:         in.Level,
		UnitType:      in.UnitType,
		ParentId:      in.ParentID,
		Path:          in.Path,
		SortOrder:     int32(in.SortOrder),
		Latitude:      in.Latitude,
		Longitude:     in.Longitude,
		Status:        in.Status,
		EffectiveFrom: in.EffectiveFrom,
		EffectiveTo:   in.EffectiveTo,
		Source:        in.Source,
		MetadataJson:  in.MetadataJSON,
		CreatedAt:     toTimestamp(in.CreatedAt),
		UpdatedAt:     toTimestamp(in.UpdatedAt),
	}
}

func toProtoAdministrativeUnits(in []*biz.AdministrativeUnit) []*pb.AdministrativeUnit {
	out := make([]*pb.AdministrativeUnit, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoAdministrativeUnit(item))
	}
	return out
}

func toProtoAdministrativeUnitNodes(in []*biz.AdministrativeUnitNode) []*pb.AdministrativeUnitNode {
	out := make([]*pb.AdministrativeUnitNode, 0, len(in))
	for _, node := range in {
		out = append(out, &pb.AdministrativeUnitNode{
			Unit:     toProtoAdministrativeUnit(node.Unit),
			Children: toProtoAdministrativeUnitNodes(node.Children),
		})
	}
	return out
}
