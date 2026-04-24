package data

import (
	"context"
	"fmt"

	"github.com/arda-labs/arda/pkg/middleware"
	"github.com/arda-labs/arda/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{data: data}
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO users (external_id, email, display_name) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
			user.ExternalID, user.Email, user.DisplayName,
		).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	})
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*biz.User, error) {
	u := &biz.User{}
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT id, external_id, email, display_name, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL`, id,
		).Scan(&u.ID, &u.ExternalID, &u.Email, &u.DisplayName, &u.CreatedAt, &u.UpdatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if u.ID == "" {
		return nil, nil
	}
	return u, err
}

func (r *userRepo) GetByExternalID(ctx context.Context, externalID string) (*biz.User, error) {
	u := &biz.User{}
	// external_id is global, we don't necessarily need RLS for global lookups if users table isn't fully partitioned by tenant
	// But we follow the pattern for consistency if RLS is enabled on the table
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT id, external_id, email, display_name, created_at, updated_at FROM users WHERE external_id = $1 AND deleted_at IS NULL`, externalID,
		).Scan(&u.ID, &u.ExternalID, &u.Email, &u.DisplayName, &u.CreatedAt, &u.UpdatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if u.ID == "" {
		return nil, nil
	}
	return u, err
}

func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`UPDATE users SET email = $2, display_name = $3, updated_at = now() WHERE id = $1 RETURNING updated_at`,
			user.ID, user.Email, user.DisplayName,
		).Scan(&user.UpdatedAt)
	})
	return user, err
}

func (r *userRepo) ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*biz.User, string, error) {
	var users []*biz.User
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT u.id, u.external_id, u.email, u.display_name, u.created_at, u.updated_at
			 FROM users u JOIN memberships m ON u.id = m.user_id
			 WHERE m.tenant_id = $1 AND m.deleted_at IS NULL AND u.deleted_at IS NULL
			 ORDER BY u.created_at DESC LIMIT $2`, tenantID, pageSize+1)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			u := &biz.User{}
			if err := rows.Scan(&u.ID, &u.ExternalID, &u.Email, &u.DisplayName, &u.CreatedAt, &u.UpdatedAt); err != nil {
				return err
			}
			users = append(users, u)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if len(users) > pageSize {
		nextToken = users[pageSize-1].ID
		users = users[:pageSize]
	}
	return users, nextToken, nil
}
