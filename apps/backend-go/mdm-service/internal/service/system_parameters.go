package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/mdm-service/api/mdm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
)

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
