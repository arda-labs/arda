package data

import (
	"context"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
)

type groupRepo struct {
	data *Data
	log  *log.Helper
}

func NewGroupRepo(data *Data, logger log.Logger) biz.GroupRepo {
	return &groupRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *groupRepo) Create(ctx context.Context, g *biz.Group) (*biz.Group, error) {
	err := r.data.DB(ctx).ExecInTransaction(ctx, g.TenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO groups (tenant_id, name, description) VALUES ($1, $2, $3)
			 RETURNING id, created_at, updated_at`,
			g.TenantID, g.Name, g.Description,
		).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
	})
	return g, err
}

func (r *groupRepo) GetByID(ctx context.Context, id string) (*biz.Group, error) {
	g := &biz.Group{}
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`SELECT id, tenant_id, name, description, created_at, updated_at FROM groups WHERE id = $1 AND deleted_at IS NULL`,
			id,
		).Scan(&g.ID, &g.TenantID, &g.Name, &g.Description, &g.CreatedAt, &g.UpdatedAt)
	})
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return g, err
}

func (r *groupRepo) Update(ctx context.Context, g *biz.Group) (*biz.Group, error) {
	err := r.data.DB(ctx).ExecInTransaction(ctx, g.TenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`UPDATE groups SET name = $2, description = $3, updated_at = now() WHERE id = $1 RETURNING updated_at`,
			g.ID, g.Name, g.Description,
		).Scan(&g.UpdatedAt)
	})
	return g, err
}

func (r *groupRepo) Delete(ctx context.Context, id string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE groups SET deleted_at = now() WHERE id = $1`, id)
		return err
	})
}

func (r *groupRepo) ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*biz.Group, string, error) {
	var groups []*biz.Group
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT id, tenant_id, name, description, created_at, updated_at FROM groups
			 WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2`,
			tenantID, pageSize,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			g := &biz.Group{}
			if err := rows.Scan(&g.ID, &g.TenantID, &g.Name, &g.Description, &g.CreatedAt, &g.UpdatedAt); err != nil {
				return err
			}
			groups = append(groups, g)
		}
		return nil
	})
	return groups, "", err
}

func (r *groupRepo) AddMember(ctx context.Context, groupID, userID string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, groupID, userID)
		return err
	})
}

func (r *groupRepo) RemoveMember(ctx context.Context, groupID, userID string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`, groupID, userID)
		return err
	})
}

func (r *groupRepo) ListMembers(ctx context.Context, groupID string) ([]*biz.User, error) {
	var users []*biz.User
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT u.id, u.external_id, u.email, u.display_name, u.created_at, u.updated_at
			 FROM users u JOIN group_members gm ON u.id = gm.user_id
			 WHERE gm.group_id = $1 AND u.deleted_at IS NULL`,
			groupID,
		)
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
		return nil
	})
	return users, err
}

func (r *groupRepo) AssignRole(ctx context.Context, groupID, roleID, tenantID string) error {
	return r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO group_roles (group_id, role_id, tenant_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, groupID, roleID, tenantID)
		return err
	})
}

func (r *groupRepo) RevokeRole(ctx context.Context, groupID, roleID string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `DELETE FROM group_roles WHERE group_id = $1 AND role_id = $2`, groupID, roleID)
		return err
	})
}

func (r *groupRepo) ListRoles(ctx context.Context, groupID string) ([]*biz.Role, error) {
	var roles []*biz.Role
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.DB(ctx).ExecInTransaction(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT r.id, r.tenant_id, r.name, r.description, r.is_system, r.created_at, r.updated_at
			 FROM roles r JOIN group_roles gr ON r.id = gr.role_id
			 WHERE gr.group_id = $1 AND r.deleted_at IS NULL`,
			groupID,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			role := &biz.Role{}
			if err := rows.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt); err != nil {
				return err
			}
			roles = append(roles, role)
		}
		return nil
	})
	return roles, err
}
