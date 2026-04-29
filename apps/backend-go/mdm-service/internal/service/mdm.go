package service

import (
	"context"
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

func (s *MdmService) ListCodeSets(ctx context.Context, req *pb.ListCodeSetsRequest) (*pb.ListCodeSetsResponse, error) {
	list, next, err := s.uc.ListCodeSets(ctx, biz.PageFilter{
		Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken,
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListCodeSetsResponse{CodeSets: toProtoCodeSets(list), NextPageToken: next}, nil
}

func (s *MdmService) GetCodeSet(ctx context.Context, req *pb.GetCodeSetRequest) (*pb.CodeSet, error) {
	item, err := s.uc.GetCodeSet(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCodeSet(item), nil
}

func (s *MdmService) CreateCodeSet(ctx context.Context, req *pb.CreateCodeSetRequest) (*pb.CodeSet, error) {
	item, err := s.uc.CreateCodeSet(ctx, toBizCodeSet(req.CodeSet))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCodeSet(item), nil
}

func (s *MdmService) UpdateCodeSet(ctx context.Context, req *pb.UpdateCodeSetRequest) (*pb.CodeSet, error) {
	item := toBizCodeSet(req.CodeSet)
	item.ID = req.Id
	updated, err := s.uc.UpdateCodeSet(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCodeSet(updated), nil
}

func (s *MdmService) DeleteCodeSet(ctx context.Context, req *pb.DeleteCodeSetRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteCodeSet(ctx, req.Id))
}

func (s *MdmService) ListCodeItems(ctx context.Context, req *pb.ListCodeItemsRequest) (*pb.ListCodeItemsResponse, error) {
	list, next, err := s.uc.ListCodeItems(ctx, biz.CodeItemFilter{
		CodeSetCode: req.CodeSetCode,
		PageFilter:  biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken},
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListCodeItemsResponse{CodeItems: toProtoCodeItems(list), NextPageToken: next}, nil
}

func (s *MdmService) GetCodeItem(ctx context.Context, req *pb.GetCodeItemRequest) (*pb.CodeItem, error) {
	item, err := s.uc.GetCodeItem(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCodeItem(item), nil
}

func (s *MdmService) CreateCodeItem(ctx context.Context, req *pb.CreateCodeItemRequest) (*pb.CodeItem, error) {
	item, err := s.uc.CreateCodeItem(ctx, req.CodeSetCode, toBizCodeItem(req.CodeItem))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCodeItem(item), nil
}

func (s *MdmService) UpdateCodeItem(ctx context.Context, req *pb.UpdateCodeItemRequest) (*pb.CodeItem, error) {
	item := toBizCodeItem(req.CodeItem)
	item.ID = req.Id
	updated, err := s.uc.UpdateCodeItem(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoCodeItem(updated), nil
}

func (s *MdmService) DeleteCodeItem(ctx context.Context, req *pb.DeleteCodeItemRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteCodeItem(ctx, req.Id))
}

func (s *MdmService) ListSystemParameters(ctx context.Context, req *pb.ListSystemParametersRequest) (*pb.ListSystemParametersResponse, error) {
	list, next, err := s.uc.ListSystemParameters(ctx, biz.SystemParameterFilter{
		GroupCode:  req.GroupCode,
		PageFilter: biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken},
	})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListSystemParametersResponse{Parameters: toProtoSystemParameters(list), NextPageToken: next}, nil
}

func (s *MdmService) GetSystemParameter(ctx context.Context, req *pb.GetSystemParameterRequest) (*pb.SystemParameter, error) {
	item, err := s.uc.GetSystemParameter(ctx, req.Key)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoSystemParameter(item), nil
}

func (s *MdmService) CreateSystemParameter(ctx context.Context, req *pb.CreateSystemParameterRequest) (*pb.SystemParameter, error) {
	item, err := s.uc.CreateSystemParameter(ctx, toBizSystemParameter(req.Parameter))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoSystemParameter(item), nil
}

func (s *MdmService) UpdateSystemParameter(ctx context.Context, req *pb.UpdateSystemParameterRequest) (*pb.SystemParameter, error) {
	item := toBizSystemParameter(req.Parameter)
	item.Key = req.Key
	updated, err := s.uc.UpdateSystemParameter(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoSystemParameter(updated), nil
}

func (s *MdmService) DeleteSystemParameter(ctx context.Context, req *pb.DeleteSystemParameterRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteSystemParameter(ctx, req.Key))
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

func toBizCodeSet(in *pb.CodeSet) *biz.CodeSet {
	if in == nil {
		return &biz.CodeSet{}
	}
	return &biz.CodeSet{
		ID:          in.Id,
		Code:        in.Code,
		Name:        in.Name,
		Description: in.Description,
		IsSystem:    in.IsSystem,
		Status:      in.Status,
	}
}

func toProtoCodeSet(in *biz.CodeSet) *pb.CodeSet {
	if in == nil {
		return nil
	}
	return &pb.CodeSet{
		Id:          in.ID,
		Code:        in.Code,
		Name:        in.Name,
		Description: in.Description,
		IsSystem:    in.IsSystem,
		Status:      in.Status,
		CreatedAt:   toTimestamp(in.CreatedAt),
		UpdatedAt:   toTimestamp(in.UpdatedAt),
	}
}

func toProtoCodeSets(in []*biz.CodeSet) []*pb.CodeSet {
	out := make([]*pb.CodeSet, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoCodeSet(item))
	}
	return out
}

func toBizCodeItem(in *pb.CodeItem) *biz.CodeItem {
	if in == nil {
		return &biz.CodeItem{}
	}
	return &biz.CodeItem{
		ID:            in.Id,
		CodeSetID:     in.CodeSetId,
		CodeSetCode:   in.CodeSetCode,
		Code:          in.Code,
		Name:          in.Name,
		Value:         in.Value,
		ParentID:      in.ParentId,
		SortOrder:     int(in.SortOrder),
		Color:         in.Color,
		Icon:          in.Icon,
		MetadataJSON:  in.MetadataJson,
		IsDefault:     in.IsDefault,
		IsSystem:      in.IsSystem,
		Status:        in.Status,
		EffectiveFrom: in.EffectiveFrom,
		EffectiveTo:   in.EffectiveTo,
	}
}

func toProtoCodeItem(in *biz.CodeItem) *pb.CodeItem {
	if in == nil {
		return nil
	}
	return &pb.CodeItem{
		Id:            in.ID,
		CodeSetId:     in.CodeSetID,
		CodeSetCode:   in.CodeSetCode,
		Code:          in.Code,
		Name:          in.Name,
		Value:         in.Value,
		ParentId:      in.ParentID,
		SortOrder:     int32(in.SortOrder),
		Color:         in.Color,
		Icon:          in.Icon,
		MetadataJson:  in.MetadataJSON,
		IsDefault:     in.IsDefault,
		IsSystem:      in.IsSystem,
		Status:        in.Status,
		EffectiveFrom: in.EffectiveFrom,
		EffectiveTo:   in.EffectiveTo,
		CreatedAt:     toTimestamp(in.CreatedAt),
		UpdatedAt:     toTimestamp(in.UpdatedAt),
	}
}

func toProtoCodeItems(in []*biz.CodeItem) []*pb.CodeItem {
	out := make([]*pb.CodeItem, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoCodeItem(item))
	}
	return out
}

func toBizSystemParameter(in *pb.SystemParameter) *biz.SystemParameter {
	if in == nil {
		return &biz.SystemParameter{IsEditable: true}
	}
	return &biz.SystemParameter{
		ID:                 in.Id,
		Key:                in.Key,
		Name:               in.Name,
		GroupCode:          in.GroupCode,
		ValueType:          in.ValueType,
		ValueText:          in.ValueText,
		ValueNumber:        in.ValueNumber,
		ValueBoolean:       in.ValueBoolean,
		ValueJSON:          in.ValueJson,
		DefaultValue:       in.DefaultValue,
		IsSecret:           in.IsSecret,
		IsEditable:         in.IsEditable,
		IsSystem:           in.IsSystem,
		ValidationRuleJSON: in.ValidationRuleJson,
		Description:        in.Description,
		Status:             in.Status,
		UpdatedBy:          in.UpdatedBy,
	}
}

func toProtoSystemParameter(in *biz.SystemParameter) *pb.SystemParameter {
	if in == nil {
		return nil
	}
	return &pb.SystemParameter{
		Id:                 in.ID,
		Key:                in.Key,
		Name:               in.Name,
		GroupCode:          in.GroupCode,
		ValueType:          in.ValueType,
		ValueText:          in.ValueText,
		ValueNumber:        in.ValueNumber,
		ValueBoolean:       in.ValueBoolean,
		ValueJson:          in.ValueJSON,
		DefaultValue:       in.DefaultValue,
		IsSecret:           in.IsSecret,
		IsEditable:         in.IsEditable,
		IsSystem:           in.IsSystem,
		ValidationRuleJson: in.ValidationRuleJSON,
		Description:        in.Description,
		Status:             in.Status,
		UpdatedBy:          in.UpdatedBy,
		CreatedAt:          toTimestamp(in.CreatedAt),
		UpdatedAt:          toTimestamp(in.UpdatedAt),
	}
}

func toProtoSystemParameters(in []*biz.SystemParameter) []*pb.SystemParameter {
	out := make([]*pb.SystemParameter, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoSystemParameter(item))
	}
	return out
}
