package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/notification-service/api/notification/v1"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
)

func (s *NotificationService) ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error) {
	list, next, err := s.uc.ListTemplates(ctx, biz.PageFilter{Status: req.Status, Keyword: req.Keyword, PageSize: int(req.PageSize), PageToken: req.PageToken})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListTemplatesResponse{Templates: toProtoTemplates(list), NextPageToken: next}, nil
}

func (s *NotificationService) GetTemplate(ctx context.Context, req *pb.GetTemplateRequest) (*pb.NotificationTemplate, error) {
	item, err := s.uc.GetTemplate(ctx, req.Id)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTemplate(item), nil
}

func (s *NotificationService) CreateTemplate(ctx context.Context, req *pb.CreateTemplateRequest) (*pb.NotificationTemplate, error) {
	item, err := s.uc.CreateTemplate(ctx, toBizTemplate(req.Template))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTemplate(item), nil
}

func (s *NotificationService) UpdateTemplate(ctx context.Context, req *pb.UpdateTemplateRequest) (*pb.NotificationTemplate, error) {
	item := toBizTemplate(req.Template)
	item.ID = req.Id
	out, err := s.uc.UpdateTemplate(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTemplate(out), nil
}

func (s *NotificationService) DeleteTemplate(ctx context.Context, req *pb.DeleteTemplateRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, toServiceError(s.uc.DeleteTemplate(ctx, req.Id))
}

func toBizTemplate(in *pb.NotificationTemplate) *biz.NotificationTemplate {
	if in == nil {
		return &biz.NotificationTemplate{}
	}
	return &biz.NotificationTemplate{ID: in.Id, Code: in.Code, Name: in.Name, Category: in.Category, DefaultChannel: in.DefaultChannel, Description: in.Description, Status: in.Status}
}

func toProtoTemplate(in *biz.NotificationTemplate) *pb.NotificationTemplate {
	if in == nil {
		return nil
	}
	return &pb.NotificationTemplate{Id: in.ID, Code: in.Code, Name: in.Name, Category: in.Category, DefaultChannel: in.DefaultChannel, Description: in.Description, Status: in.Status, CreatedAt: toTimestamp(in.CreatedAt), UpdatedAt: toTimestamp(in.UpdatedAt)}
}

func toProtoTemplates(in []*biz.NotificationTemplate) []*pb.NotificationTemplate {
	out := make([]*pb.NotificationTemplate, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoTemplate(item))
	}
	return out
}
