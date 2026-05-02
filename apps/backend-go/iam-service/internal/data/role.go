package data

import (
	"context"
	"fmt"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

type roleRepo struct {
	data *Data
}

func NewRoleRepo(data *Data) biz.RoleRepo {
	return &roleRepo{data: data}
}

func (r *roleRepo) Create(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO roles (tenant_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
			role.TenantID, role.Name, role.Description,
		).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
	})
	return role, err
}

func (r *roleRepo) GetByID(ctx context.Context, id string) (*biz.Role, error) {
	role := &biz.Role{}
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT id, tenant_id, name, description, is_system, created_at, updated_at FROM roles WHERE id = $1 AND deleted_at IS NULL`, id,
		).Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if role.ID == "" {
		return nil, nil
	}
	return role, err
}

func (r *roleRepo) GetByName(ctx context.Context, tenantID, name string) (*biz.Role, error) {
	role := &biz.Role{}
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT id, tenant_id, name, description, is_system, created_at, updated_at
			 FROM roles
			 WHERE tenant_id = $1
			   AND lower(name) = lower($2)
			   AND deleted_at IS NULL
			 LIMIT 1`, tenantID, name,
		).Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if role.ID == "" {
		return nil, nil
	}
	return role, err
}

func (r *roleRepo) ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*biz.Role, string, error) {
	var list []*biz.Role
	page := pagination.Normalize(pageSize, cursor)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT id, tenant_id, name, description, is_system, created_at, updated_at FROM roles
			 WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC, id DESC LIMIT $2 OFFSET $3`,
			tenantID, page.Limit+1, page.Offset,
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
			list = append(list, role)
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

func (r *roleRepo) Update(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`UPDATE roles SET name = $2, description = $3, updated_at = now() WHERE id = $1 RETURNING updated_at`,
			role.ID, role.Name, role.Description,
		).Scan(&role.UpdatedAt)
	})
	return role, err
}

func (r *roleRepo) SoftDelete(ctx context.Context, id string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE roles SET deleted_at = now() WHERE id = $1`, id)
		return err
	})
}

func (r *roleRepo) AssignRole(ctx context.Context, userID, roleID, tenantID string) error {
	return r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`INSERT INTO user_roles (user_id, role_id, tenant_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			userID, roleID, tenantID)
		return err
	})
}

func (r *roleRepo) RevokeRole(ctx context.Context, userID, roleID, tenantID string) error {
	return r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2 AND tenant_id = $3`,
			userID, roleID, tenantID)
		return err
	})
}

func (r *roleRepo) GetUserRoles(ctx context.Context, userID, tenantID string) ([]*biz.Role, error) {
	var roles []*biz.Role
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT r.id, r.tenant_id, r.name, r.description, r.is_system, r.created_at, r.updated_at
			 FROM roles r
			 JOIN user_roles ur ON r.id = ur.role_id
			 JOIN users u ON u.id = ur.user_id
			 LEFT JOIN tenant_users tu
			   ON tu.user_id = u.id
			  AND tu.tenant_id = $2
			  AND tu.deleted_at IS NULL
			 WHERE (u.id::text = $1 OR u.external_id = $1 OR lower(u.email) = lower($1) OR lower(tu.username) = lower($1))
			   AND ur.tenant_id = $2
			   AND r.deleted_at IS NULL
			   AND u.deleted_at IS NULL`, userID, tenantID)
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
		return rows.Err()
	})
	return roles, err
}

func (r *roleRepo) GetGroupRoles(ctx context.Context, userID, tenantID string) ([]*biz.Role, error) {
	var roles []*biz.Role
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT r.id, r.tenant_id, r.name, r.description, r.is_system, r.created_at, r.updated_at
			 FROM roles r
			 JOIN group_roles gr ON r.id = gr.role_id
			 JOIN group_members gm ON gr.group_id = gm.group_id
			 JOIN users u ON u.id = gm.user_id
			 LEFT JOIN tenant_users tu
			   ON tu.user_id = u.id
			  AND tu.tenant_id = $2
			  AND tu.deleted_at IS NULL
			 WHERE (u.id::text = $1 OR u.external_id = $1 OR lower(u.email) = lower($1) OR lower(tu.username) = lower($1))
			   AND gr.tenant_id = $2
			   AND r.deleted_at IS NULL
			   AND u.deleted_at IS NULL`, userID, tenantID)
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
		return rows.Err()
	})
	return roles, err
}

func (r *roleRepo) GetRolePermissions(ctx context.Context, roleID string) ([]*biz.Permission, error) {
	var perms []*biz.Permission
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT p.id, p.tenant_id, p.resource, p.action
			 FROM permissions p JOIN role_permissions rp ON p.id = rp.permission_id
			 WHERE rp.role_id = $1`, roleID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			p := &biz.Permission{}
			if err := rows.Scan(&p.ID, &p.TenantID, &p.Resource, &p.Action); err != nil {
				return err
			}
			perms = append(perms, p)
		}
		return rows.Err()
	})
	return perms, err
}

func (r *roleRepo) SetRolePermissions(ctx context.Context, roleID string, permIDs []string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID)
		if err != nil {
			return err
		}
		if len(permIDs) == 0 {
			return nil
		}
		// permIDs chứa permission codes dạng "resource:action" → resolve sang UUID từ bảng permissions
		for _, code := range permIDs {
			_, err = tx.Exec(ctx,
				`INSERT INTO role_permissions (role_id, permission_id)
				 SELECT $1, p.id FROM permissions p
				 WHERE p.tenant_id = $2 AND (p.resource || ':' || p.action) = $3
				 ON CONFLICT DO NOTHING`,
				roleID, tenantID, code,
			)
			if err != nil {
				return fmt.Errorf("inserting role_permission for code %q: %w", code, err)
			}
		}
		return nil
	})
}
