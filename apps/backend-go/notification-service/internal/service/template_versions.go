package service

import (
	"context"

	pb "github.com/arda-labs/arda/arda-be-go/services/notification-service/api/notification/v1"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
)

func (s *NotificationService) ListTemplateVersions(ctx context.Context, req *pb.ListTemplateVersionsRequest) (*pb.ListTemplateVersionsResponse, error) {
	list, err := s.uc.ListTemplateVersions(ctx, biz.TemplateVersionFilter{TemplateID: req.TemplateId, Channel: req.Channel, Language: req.Language, Status: req.Status})
	if err != nil {
		return nil, toServiceError(err)
	}
	return &pb.ListTemplateVersionsResponse{Versions: toProtoTemplateVersions(list)}, nil
}

func (s *NotificationService) CreateTemplateVersion(ctx context.Context, req *pb.CreateTemplateVersionRequest) (*pb.NotificationTemplateVersion, error) {
	item := toBizTemplateVersion(req.Version)
	item.TemplateID = req.TemplateId
	out, err := s.uc.CreateTemplateVersion(ctx, item)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTemplateVersion(out), nil
}

func (s *NotificationService) ApproveTemplateVersion(ctx context.Context, req *pb.ApproveTemplateVersionRequest) (*pb.NotificationTemplateVersion, error) {
	item, err := s.uc.ApproveTemplateVersion(ctx, req.Id, req.Actor, req.Note)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toProtoTemplateVersion(item), nil
}

func toBizTemplateVersion(in *pb.NotificationTemplateVersion) *biz.NotificationTemplateVersion {
	if in == nil {
		return &biz.NotificationTemplateVersion{}
	}
	return &biz.NotificationTemplateVersion{ID: in.Id, TemplateID: in.TemplateId, Version: int(in.Version), Channel: in.Channel, Language: in.Language, Subject: in.Subject, Body: in.Body, PayloadSchemaJSON: in.PayloadSchemaJson, ApprovalStatus: in.ApprovalStatus, ApprovedBy: in.ApprovedBy, ChangeNote: in.ChangeNote, Status: in.Status}
}

func toProtoTemplateVersion(in *biz.NotificationTemplateVersion) *pb.NotificationTemplateVersion {
	if in == nil {
		return nil
	}
	return &pb.NotificationTemplateVersion{Id: in.ID, TemplateId: in.TemplateID, Version: int32(in.Version), Channel: in.Channel, Language: in.Language, Subject: in.Subject, Body: in.Body, PayloadSchemaJson: in.PayloadSchemaJSON, ApprovalStatus: in.ApprovalStatus, ApprovedBy: in.ApprovedBy, ApprovedAt: toTimestamp(in.ApprovedAt), ChangeNote: in.ChangeNote, Status: in.Status, CreatedAt: toTimestamp(in.CreatedAt), UpdatedAt: toTimestamp(in.UpdatedAt)}
}

func toProtoTemplateVersions(in []*biz.NotificationTemplateVersion) []*pb.NotificationTemplateVersion {
	out := make([]*pb.NotificationTemplateVersion, 0, len(in))
	for _, item := range in {
		out = append(out, toProtoTemplateVersion(item))
	}
	return out
}
