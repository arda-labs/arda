package data

import (
	"context"
	"encoding/json"

	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"
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

func (r *auditRepo) ListByActor(ctx context.Context, actorID string, pageSize int, cursor string) ([]*biz.AuditLog, string, error) {
	var list []*biz.AuditLog
	// Simple implementation without complex pagination for now
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		query := `
			SELECT id, actor_id, tenant_id, action, target_type, target_id, metadata, created_at
			FROM audit_logs
			WHERE actor_id = $1
			ORDER BY created_at DESC
			LIMIT $2
		`
		rows, err := tx.Query(ctx, query, actorID, pageSize)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var l biz.AuditLog
			var meta []byte
			if err := rows.Scan(&l.ID, &l.ActorID, &l.TenantID, &l.Action, &l.TargetType, &l.TargetID, &meta, &l.CreatedAt); err != nil {
				return err
			}
			if len(meta) > 0 {
				_ = json.Unmarshal(meta, &l.Metadata)
			}
			list = append(list, &l)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, "", err
	}
	return list, "", nil
}
