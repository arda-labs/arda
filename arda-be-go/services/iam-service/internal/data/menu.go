package data

import (
	"context"
	"fmt"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type menuRepo struct {
	data *Data
	log  *log.Helper
}

func NewMenuRepo(data *Data, logger log.Logger) biz.MenuRepo {
	return &menuRepo{data: data, log: log.NewHelper(logger)}
}

func (r *menuRepo) GetByTenant(ctx context.Context, tenantID string) ([]*biz.Menu, error) {
	var list []*biz.Menu
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		query := `
			SELECT id, tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, created_at, updated_at
			FROM menus
			WHERE tenant_id = $1
			ORDER BY sort_order ASC
		`
		rows, err := tx.Query(ctx, query, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var m biz.Menu
			var parentID, icon, route *string
			if err := rows.Scan(&m.ID, &m.TenantID, &parentID, &m.Name, &m.Slug, &icon, &route, &m.SortOrder, &m.Enabled, &m.CreatedAt, &m.UpdatedAt); err != nil {
				return err
			}
			if parentID != nil {
				m.ParentID = *parentID
			}
			if icon != nil {
				m.Icon = *icon
			}
			if route != nil {
				m.Route = *route
			}
			list = append(list, &m)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, fmt.Errorf("query menus: %w", err)
	}
	return list, nil
}

func (r *menuRepo) GetByID(ctx context.Context, id string) (*biz.Menu, error) {
	var m biz.Menu
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		query := `
			SELECT id, tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, created_at, updated_at
			FROM menus
			WHERE id = $1
		`
		var parentID, icon, route *string
		err := tx.QueryRow(ctx, query, id).Scan(
			&m.ID, &m.TenantID, &parentID, &m.Name, &m.Slug, &icon, &route, &m.SortOrder, &m.Enabled, &m.CreatedAt, &m.UpdatedAt,
		)
		if err == pgx.ErrNoRows {
			return biz.ErrMenuNotFound
		}
		if err != nil {
			return err
		}
		if parentID != nil {
			m.ParentID = *parentID
		}
		if icon != nil {
			m.Icon = *icon
		}
		if route != nil {
			m.Route = *route
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *menuRepo) Create(ctx context.Context, m *biz.Menu) (*biz.Menu, error) {
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		query := `
			INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled)
			VALUES ($1, NULLIF($2, ''), $3, $4, NULLIF($5, ''), NULLIF($6, ''), $7, $8)
			RETURNING id, created_at, updated_at
		`
		return tx.QueryRow(ctx, query, m.TenantID, m.ParentID, m.Name, m.Slug, m.Icon, m.Route, m.SortOrder, m.Enabled).
			Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	})
	if err != nil {
		return nil, fmt.Errorf("insert menu: %w", err)
	}
	return m, nil
}

func (r *menuRepo) Update(ctx context.Context, m *biz.Menu) (*biz.Menu, error) {
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		query := `
			UPDATE menus
			SET parent_id = NULLIF($2, ''), name = $3, slug = $4, icon = NULLIF($5, ''),
			    route = NULLIF($6, ''), sort_order = $7, enabled = $8, updated_at = NOW()
			WHERE id = $1
			RETURNING updated_at
		`
		err := tx.QueryRow(ctx, query, m.ID, m.ParentID, m.Name, m.Slug, m.Icon, m.Route, m.SortOrder, m.Enabled).
			Scan(&m.UpdatedAt)
		if err == pgx.ErrNoRows {
			return biz.ErrMenuNotFound
		}
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("update menu: %w", err)
	}
	return m, nil
}

func (r *menuRepo) Delete(ctx context.Context, id string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		query := `DELETE FROM menus WHERE id = $1`
		ct, err := tx.Exec(ctx, query, id)
		if err != nil {
			return err
		}
		if ct.RowsAffected() == 0 {
			return biz.ErrMenuNotFound
		}
		return nil
	})
}
