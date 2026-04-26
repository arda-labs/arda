package biz

import (
	"context"
	"time"
)

type AuditLog struct {
	ID         string
	ActorID    string
	TenantID   string
	Action     string
	TargetType string
	TargetID   string
	Metadata   map[string]interface{}
	CreatedAt  time.Time
}

type AuditRepo interface {
	Create(ctx context.Context, log *AuditLog) error
	ListByActor(ctx context.Context, actorID string, pageSize int, cursor string) ([]*AuditLog, string, error)
}

type AuditUsecase struct {
	repo AuditRepo
}

func NewAuditUsecase(repo AuditRepo) *AuditUsecase {
	return &AuditUsecase{repo: repo}
}

func (uc *AuditUsecase) ListMyLogs(ctx context.Context, actorID string, pageSize int, cursor string) ([]*AuditLog, string, error) {
	return uc.repo.ListByActor(ctx, actorID, pageSize, cursor)
}

func (uc *AuditUsecase) Log(ctx context.Context, actorID, tenantID, action, targetType, targetID string, metadata map[string]interface{}) {
	meta := metadata
	if meta == nil {
		meta = make(map[string]interface{})
	}
	_ = uc.repo.Create(ctx, &AuditLog{
		ActorID:    actorID,
		TenantID:   tenantID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Metadata:   meta,
	})
}
