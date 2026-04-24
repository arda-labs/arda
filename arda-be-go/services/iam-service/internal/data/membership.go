package data

import (
	"context"
	"fmt"

	"github.com/arda-labs/arda/pkg/middleware"
	"github.com/arda-labs/arda/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

type membershipRepo struct {
	data *Data
}

func NewMembershipRepo(data *Data) biz.MembershipRepo {
	return &membershipRepo{data: data}
}

func (r *membershipRepo) Create(ctx context.Context, m *biz.Membership) (*biz.Membership, error) {
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1, $2, $3) RETURNING id, created_at`,
			m.UserID, m.TenantID, m.Role,
		).Scan(&m.ID, &m.CreatedAt)
	})
	return m, err
}

func (r *membershipRepo) ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*biz.Membership, string, error) {
	var list []*biz.Membership
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		q := `SELECT m.id, m.user_id, m.tenant_id, m.role, m.created_at
			  FROM memberships m WHERE m.tenant_id = $1 AND m.deleted_at IS NULL
			  ORDER BY m.created_at DESC LIMIT $2`
		rows, err := tx.Query(ctx, q, tenantID, pageSize+1)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			m := &biz.Membership{}
			if err := rows.Scan(&m.ID, &m.UserID, &m.TenantID, &m.Role, &m.CreatedAt); err != nil {
				return err
			}
			list = append(list, m)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, "", err
	}
	var next string
	if len(list) > pageSize {
		next = list[pageSize-1].ID
		list = list[:pageSize]
	}
	return list, next, nil
}

func (r *membershipRepo) ListByUser(ctx context.Context, userID string) ([]*biz.Membership, error) {
	var list []*biz.Membership
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		q := `SELECT m.id, m.user_id, m.tenant_id, m.role, m.created_at
			  FROM memberships m WHERE m.user_id = $1 AND m.deleted_at IS NULL
			  ORDER BY m.created_at ASC`
		rows, err := tx.Query(ctx, q, userID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			m := &biz.Membership{}
			if err := rows.Scan(&m.ID, &m.UserID, &m.TenantID, &m.Role, &m.CreatedAt); err != nil {
				return err
			}
			list = append(list, m)
		}
		return rows.Err()
	})
	return list, err
}

func (r *membershipRepo) GetByUserAndTenant(ctx context.Context, userID, tenantID string) (*biz.Membership, error) {
	m := &biz.Membership{}
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT id, user_id, tenant_id, role, created_at FROM memberships WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL`,
			userID, tenantID,
		).Scan(&m.ID, &m.UserID, &m.TenantID, &m.Role, &m.CreatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if m.ID == "" {
		return nil, nil
	}
	return m, err
}

func (r *membershipRepo) SoftDelete(ctx context.Context, userID, tenantID string) error {
	return r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`UPDATE memberships SET deleted_at = now() WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL`,
			userID, tenantID)
		return err
	})
}
