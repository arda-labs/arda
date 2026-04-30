package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
)

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
