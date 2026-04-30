package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
)

func (s *MdmService) ListAreaTypes(ctx context.Context, req *pb.ListAreaTypesRequest) (*pb.ListAreaTypesResponse, error) {
	list, next, err := s.uc.ListAreaTypes(ctx, biz.PageFilter{
		Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListAreaTypesResponse{AreaTypes: toProtoAreaTypes(list), NextPageToken: next}, nil
}

func (s *MdmService) GetAreaType(ctx context.Context, req *pb.GetAreaTypeRequest) (*pb.AreaType, error) {
	item, err := s.uc.GetAreaType(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoAreaType(item), nil
}

func (s *MdmService) CreateAreaType(ctx context.Context, req *pb.CreateAreaTypeRequest) (*pb.AreaType, error) {
	item, err := s.uc.CreateAreaType(ctx, toBizAreaType(req.AreaType))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoAreaType(item), nil
}

func (s *MdmService) UpdateAreaType(ctx context.Context, req *pb.UpdateAreaTypeRequest) (*pb.AreaType, error) {
	item := toBizAreaType(req.AreaType)
	item.ID = req.Id
	updated, err := s.uc.UpdateAreaType(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoAreaType(updated), nil
}

func (s *MdmService) DeleteAreaType(ctx context.Context, req *pb.DeleteAreaTypeRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteAreaType(ctx, req.Id))
}

func (s *MdmService) ListAreas(ctx context.Context, req *pb.ListAreasRequest) (*pb.ListAreasResponse, error) {
	list, next, err := s.uc.ListAreas(ctx, areaFilter(req))
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListAreasResponse{Areas: toProtoAreas(list), NextPageToken: next}, nil
}

func (s *MdmService) ListAreaTree(ctx context.Context, req *pb.ListAreasRequest) (*pb.ListAreaTreeResponse, error) {
	nodes, err := s.uc.ListAreaTree(ctx, areaFilter(req))
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListAreaTreeResponse{Nodes: toProtoAreaNodes(nodes)}, nil
}

func (s *MdmService) GetArea(ctx context.Context, req *pb.GetAreaRequest) (*pb.Area, error) {
	item, err := s.uc.GetArea(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoArea(item), nil
}

func (s *MdmService) CreateArea(ctx context.Context, req *pb.CreateAreaRequest) (*pb.Area, error) {
	item, err := s.uc.CreateArea(ctx, toBizArea(req.Area))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoArea(item), nil
}

func (s *MdmService) UpdateArea(ctx context.Context, req *pb.UpdateAreaRequest) (*pb.Area, error) {
	item := toBizArea(req.Area)
	item.ID = req.Id
	updated, err := s.uc.UpdateArea(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoArea(updated), nil
}

func (s *MdmService) DeleteArea(ctx context.Context, req *pb.DeleteAreaRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteArea(ctx, req.Id))
}

func (s *MdmService) AssignAreaAdministrativeUnit(ctx context.Context, req *pb.AssignAreaAdministrativeUnitRequest) (*pb.AreaAdministrativeUnit, error) {
	item, err := s.uc.AssignAreaAdministrativeUnit(ctx, &biz.AreaAdministrativeUnit{
		AreaID:               req.AreaId,
		AdministrativeUnitID: req.AdministrativeUnitId,
		ScopeType:            req.ScopeType,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoAreaAdministrativeUnit(item), nil
}

func (s *MdmService) ListAreaAdministrativeUnits(ctx context.Context, req *pb.ListAreaAdministrativeUnitsRequest) (*pb.ListAreaAdministrativeUnitsResponse, error) {
	list, err := s.uc.ListAreaAdministrativeUnits(ctx, req.AreaId)
	if err != nil {
		return nil, toServiceError(err)
	}
	out := make([]*pb.AreaAdministrativeUnit, 0, len(list))
	for _, item := range list {
		out = append(out, toProtoAreaAdministrativeUnit(item))
	}
	return &pb.ListAreaAdministrativeUnitsResponse{Items: out}, nil
}

func (s *MdmService) RemoveAreaAdministrativeUnit(ctx context.Context, req *pb.RemoveAreaAdministrativeUnitRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.RemoveAreaAdministrativeUnit(ctx, req.AreaId, req.AdministrativeUnitId))
}

func toBizAreaType(in *pb.AreaType) *biz.AreaType {
	if in == nil {
		return &biz.AreaType{AllowHierarchy: true}
	}
	return &biz.AreaType{
		ID:             in.Id,
		Code:           in.Code,
		Name:           in.Name,
		Description:    in.Description,
		AllowHierarchy: in.AllowHierarchy,
		Status:         in.Status,
	}
}

func toProtoAreaType(in *biz.AreaType) *pb.AreaType {
	if in == nil {
		return nil
	}
	return &pb.AreaType{
		Id:             in.ID,
		Code:           in.Code,
		Name:           in.Name,
		Description:    in.Description,
		AllowHierarchy: in.AllowHierarchy,
		Status:         in.Status,
		CreatedAt:      toTimestamp(in.CreatedAt),
		UpdatedAt:      toTimestamp(in.UpdatedAt),
	}
}

func toProtoAreaTypes(in []*biz.AreaType) []*pb.AreaType {
	out := make([]*pb.AreaType, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoAreaType(item))
	}
	return out
}

func toBizArea(in *pb.Area) *biz.Area {
	if in == nil {
		return &biz.Area{}
	}
	return &biz.Area{
		ID:            in.Id,
		AreaTypeID:    in.AreaTypeId,
		AreaTypeCode:  in.AreaTypeCode,
		ParentID:      in.ParentId,
		Code:          in.Code,
		Name:          in.Name,
		Description:   in.Description,
		ManagerUserID: in.ManagerUserId,
		Status:        in.Status,
		EffectiveFrom: in.EffectiveFrom,
		EffectiveTo:   in.EffectiveTo,
		MetadataJSON:  in.MetadataJson,
	}
}

func toProtoArea(in *biz.Area) *pb.Area {
	if in == nil {
		return nil
	}
	return &pb.Area{
		Id:            in.ID,
		AreaTypeId:    in.AreaTypeID,
		AreaTypeCode:  in.AreaTypeCode,
		ParentId:      in.ParentID,
		Code:          in.Code,
		Name:          in.Name,
		Description:   in.Description,
		ManagerUserId: in.ManagerUserID,
		Status:        in.Status,
		EffectiveFrom: in.EffectiveFrom,
		EffectiveTo:   in.EffectiveTo,
		MetadataJson:  in.MetadataJSON,
		CreatedAt:     toTimestamp(in.CreatedAt),
		UpdatedAt:     toTimestamp(in.UpdatedAt),
	}
}

func toProtoAreas(in []*biz.Area) []*pb.Area {
	out := make([]*pb.Area, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoArea(item))
	}
	return out
}

func toProtoAreaNodes(in []*biz.AreaNode) []*pb.AreaNode {
	out := make([]*pb.AreaNode, 0, len(in))
	for _, node := range in {
		out = append(out, &pb.AreaNode{Area: toProtoArea(node.Area), Children: toProtoAreaNodes(node.Children)})
	}
	return out
}

func toProtoAreaAdministrativeUnit(in *biz.AreaAdministrativeUnit) *pb.AreaAdministrativeUnit {
	if in == nil {
		return nil
	}
	return &pb.AreaAdministrativeUnit{
		Id:                   in.ID,
		AreaId:               in.AreaID,
		AdministrativeUnitId: in.AdministrativeUnitID,
		ScopeType:            in.ScopeType,
		CreatedAt:            toTimestamp(in.CreatedAt),
	}
}
