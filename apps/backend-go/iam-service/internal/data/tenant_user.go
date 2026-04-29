package data

import (
	"context"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

type tenantUserRepo struct {
	data *Data
}

func NewTenantUserRepo(data *Data) biz.TenantUserRepo {
	return &tenantUserRepo{data: data}
}

func (r *tenantUserRepo) Create(ctx context.Context, tu *biz.TenantUser) (*biz.TenantUser, error) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		tenantID = tu.TenantID
	}
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO tenant_users (user_id, tenant_id, username, display_name, role, status)
			 SELECT u.id,
			        $2,
			        COALESCE(NULLIF($3, ''), NULLIF(u.email, ''), 'user-' || left(u.id::text, 8)),
			        COALESCE(NULLIF($4, ''), u.display_name, ''),
			        COALESCE(NULLIF($5, ''), 'MEMBER'),
			        COALESCE(NULLIF($6, ''), 'ACTIVE')
			 FROM users u
			 WHERE u.id = $1 AND u.deleted_at IS NULL
			 RETURNING id, user_id, tenant_id, username, display_name, role, status, created_at, updated_at`,
			tu.UserID, tu.TenantID, tu.Username, tu.DisplayName, tu.Role, tu.Status,
		).Scan(&tu.ID, &tu.UserID, &tu.TenantID, &tu.Username, &tu.DisplayName, &tu.Role, &tu.Status, &tu.CreatedAt, &tu.UpdatedAt)
	})
	return tu, err
}

func (r *tenantUserRepo) ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*biz.TenantUser, string, error) {
	var list []*biz.TenantUser
	page := pagination.Normalize(pageSize, cursor)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT tu.id, tu.user_id, tu.tenant_id, tu.username, u.email, tu.display_name,
			        tu.role, tu.status, tu.created_at, tu.updated_at
			 FROM tenant_users tu
			 JOIN users u ON u.id = tu.user_id
			 WHERE tu.tenant_id = $1
			   AND tu.deleted_at IS NULL
			   AND u.deleted_at IS NULL
			 ORDER BY tu.created_at DESC, tu.id DESC
			 LIMIT $2 OFFSET $3`,
			tenantID, page.Limit+1, page.Offset,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			tu := &biz.TenantUser{}
			if err := rows.Scan(
				&tu.ID,
				&tu.UserID,
				&tu.TenantID,
				&tu.Username,
				&tu.Email,
				&tu.DisplayName,
				&tu.Role,
				&tu.Status,
				&tu.CreatedAt,
				&tu.UpdatedAt,
			); err != nil {
				return err
			}
			list = append(list, tu)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *tenantUserRepo) ListByUser(ctx context.Context, userID string) ([]*biz.TenantUser, error) {
	var list []*biz.TenantUser
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT tu.id, tu.user_id, tu.tenant_id, tu.username, u.email, tu.display_name,
			        tu.role, tu.status, tu.created_at, tu.updated_at
			 FROM tenant_users tu
			 JOIN users u ON u.id = tu.user_id
			 WHERE tu.user_id = $1
			   AND tu.deleted_at IS NULL
			   AND u.deleted_at IS NULL
			 ORDER BY tu.created_at ASC`, userID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			tu := &biz.TenantUser{}
			if err := rows.Scan(
				&tu.ID,
				&tu.UserID,
				&tu.TenantID,
				&tu.Username,
				&tu.Email,
				&tu.DisplayName,
				&tu.Role,
				&tu.Status,
				&tu.CreatedAt,
				&tu.UpdatedAt,
			); err != nil {
				return err
			}
			list = append(list, tu)
		}
		return rows.Err()
	})
	return list, err
}

func (r *tenantUserRepo) GetByUserAndTenant(ctx context.Context, userID, tenantID string) (*biz.TenantUser, error) {
	tu := &biz.TenantUser{}
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT tu.id, tu.user_id, tu.tenant_id, tu.username, u.email, tu.display_name,
			        tu.role, tu.status, tu.created_at, tu.updated_at
			 FROM tenant_users tu
			 JOIN users u ON u.id = tu.user_id
			 WHERE tu.user_id = $1
			   AND tu.tenant_id = $2
			   AND tu.deleted_at IS NULL
			   AND u.deleted_at IS NULL`,
			userID, tenantID,
		).Scan(
			&tu.ID,
			&tu.UserID,
			&tu.TenantID,
			&tu.Username,
			&tu.Email,
			&tu.DisplayName,
			&tu.Role,
			&tu.Status,
			&tu.CreatedAt,
			&tu.UpdatedAt,
		)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if tu.ID == "" {
		return nil, nil
	}
	return tu, err
}

func (r *tenantUserRepo) SoftDelete(ctx context.Context, userID, tenantID string) error {
	return r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`UPDATE tenant_users
			 SET deleted_at = now(), updated_at = now()
			 WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL`,
			userID, tenantID)
		return err
	})
}
