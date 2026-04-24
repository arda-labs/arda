package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/ngvgroup/arda/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

type tenantRepo struct {
	data *Data
}

func NewTenantRepo(data *Data) biz.TenantRepo {
	return &tenantRepo{data: data}
}

func (r *tenantRepo) Create(ctx context.Context, t *biz.Tenant) (*biz.Tenant, error) {
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO tenants (name, slug, owner_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
			t.Name, t.Slug, t.OwnerID,
		).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	})
	return t, err
}

func (r *tenantRepo) GetByID(ctx context.Context, id string) (*biz.Tenant, error) {
	t := &biz.Tenant{}
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT id, name, slug, owner_id, created_at, updated_at FROM tenants WHERE id = $1 AND deleted_at IS NULL`, id,
		).Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.CreatedAt, &t.UpdatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if t.ID == "" {
		return nil, nil
	}
	return t, err
}

func (r *tenantRepo) GetByIDs(ctx context.Context, ids []string) ([]*biz.Tenant, error) {
	if len(ids) == 0 {
		return []*biz.Tenant{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	q := fmt.Sprintf(
		`SELECT id, name, slug, owner_id, created_at, updated_at FROM tenants WHERE id IN (%s) AND deleted_at IS NULL`,
		strings.Join(placeholders, ","),
	)

	var list []*biz.Tenant
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx, q, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			t := &biz.Tenant{}
			if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.CreatedAt, &t.UpdatedAt); err != nil {
				return err
			}
			list = append(list, t)
		}
		return rows.Err()
	})

	return list, err
}

func (r *tenantRepo) GetBySlug(ctx context.Context, slug string) (*biz.Tenant, error) {
	t := &biz.Tenant{}
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT id, name, slug, owner_id, created_at, updated_at FROM tenants WHERE slug = $1 AND deleted_at IS NULL`, slug,
		).Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.CreatedAt, &t.UpdatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if t.ID == "" {
		return nil, nil
	}
	return t, err
}

func (r *tenantRepo) Update(ctx context.Context, t *biz.Tenant) (*biz.Tenant, error) {
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`UPDATE tenants SET name = $2, slug = $3, updated_at = now() WHERE id = $1 RETURNING updated_at`,
			t.ID, t.Name, t.Slug,
		).Scan(&t.UpdatedAt)
	})
	return t, err
}

func (r *tenantRepo) SoftDelete(ctx context.Context, id string) error {
	return r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE tenants SET deleted_at = now() WHERE id = $1`, id)
		return err
	})
}
