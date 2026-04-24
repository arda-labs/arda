package data

import (
	"context"
	"encoding/json"

	"github.com/arda-labs/arda/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

type auditRepo struct {
	data *Data
}

func NewAuditRepo(data *Data) biz.AuditRepo {
	return &auditRepo{data: data}
}

func (r *auditRepo) Create(ctx context.Context, log *biz.AuditLog) error {
	metaJSON, _ := json.Marshal(log.Metadata)
	return r.data.DB(ctx).ExecInTransaction(ctx, log.TenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`INSERT INTO audit_logs (actor_id, tenant_id, action, target_type, target_id, metadata) VALUES ($1, $2, $3, $4, $5, $6)`,
			log.ActorID, log.TenantID, log.Action, log.TargetType, log.TargetID, metaJSON)
		return err
	})
}
