package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

type tenantRepo struct {
	data *Data
}

const tenantSelectColumns = `
	t.id,
	t.name,
	t.slug,
	t.owner_id,
	COALESCE(trc.deployment_mode, 'SHARED') AS deployment_mode,
	COALESCE(trc.auth_mode, 'SHARED_AUTH') AS auth_mode,
	t.created_at,
	t.updated_at`

func NewTenantRepo(data *Data) biz.TenantRepo {
	return &tenantRepo{data: data}
}

func (r *tenantRepo) Create(ctx context.Context, t *biz.Tenant) (*biz.Tenant, error) {
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		if err := tx.QueryRow(ctx,
			`INSERT INTO tenants (name, slug, owner_id)
			 VALUES ($1, $2, $3)
			 RETURNING id, created_at, updated_at`,
			t.Name, t.Slug, t.OwnerID,
		).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return err
		}
		if err := r.upsertTenantRuntimeConfig(ctx, tx, t); err != nil {
			return err
		}
		return r.bootstrapTenantDefaults(ctx, tx, t.ID, t.OwnerID)
	})
	return t, err
}

func (r *tenantRepo) GetByID(ctx context.Context, id string) (*biz.Tenant, error) {
	t := &biz.Tenant{}
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT `+tenantSelectColumns+`
			 FROM tenants t
			 LEFT JOIN tenant_runtime_configs trc ON trc.tenant_id = t.id AND trc.deleted_at IS NULL
			 WHERE t.id = $1 AND t.deleted_at IS NULL`, id,
		).Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.DeploymentMode, &t.AuthMode, &t.CreatedAt, &t.UpdatedAt)
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
		`SELECT `+tenantSelectColumns+`
		 FROM tenants t
		 LEFT JOIN tenant_runtime_configs trc ON trc.tenant_id = t.id AND trc.deleted_at IS NULL
		 WHERE t.id IN (%s) AND t.deleted_at IS NULL`,
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
			if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.DeploymentMode, &t.AuthMode, &t.CreatedAt, &t.UpdatedAt); err != nil {
				return err
			}
			list = append(list, t)
		}
		return rows.Err()
	})

	return list, err
}

func (r *tenantRepo) ListAll(ctx context.Context) ([]*biz.Tenant, error) {
	var list []*biz.Tenant
	err := r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT `+tenantSelectColumns+`
			 FROM tenants t
			 LEFT JOIN tenant_runtime_configs trc ON trc.tenant_id = t.id AND trc.deleted_at IS NULL
			 WHERE t.deleted_at IS NULL
			 ORDER BY t.name ASC`,
		)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			t := &biz.Tenant{}
			if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.DeploymentMode, &t.AuthMode, &t.CreatedAt, &t.UpdatedAt); err != nil {
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
			`SELECT `+tenantSelectColumns+`
			 FROM tenants t
			 LEFT JOIN tenant_runtime_configs trc ON trc.tenant_id = t.id AND trc.deleted_at IS NULL
			 WHERE t.slug = $1 AND t.deleted_at IS NULL`, slug,
		).Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.DeploymentMode, &t.AuthMode, &t.CreatedAt, &t.UpdatedAt)
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
		if err := tx.QueryRow(ctx,
			`UPDATE tenants
			 SET name = $2, slug = $3, updated_at = now()
			 WHERE id = $1
			 RETURNING updated_at`,
			t.ID, t.Name, t.Slug,
		).Scan(&t.UpdatedAt); err != nil {
			return err
		}
		return r.upsertTenantRuntimeConfig(ctx, tx, t)
	})
	return t, err
}

func (r *tenantRepo) SoftDelete(ctx context.Context, id string) error {
	return r.data.DB(ctx).ExecInTransaction(ctx, "", func(ctx context.Context, tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE tenants SET deleted_at = now() WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `UPDATE tenant_runtime_configs SET deleted_at = now() WHERE tenant_id = $1 AND deleted_at IS NULL`, id)
		return err
	})
}

func (r *tenantRepo) upsertTenantRuntimeConfig(ctx context.Context, tx pgx.Tx, t *biz.Tenant) error {
	deploymentMode := t.DeploymentMode
	if deploymentMode == "" {
		deploymentMode = biz.TenantDeploymentShared
	}
	authMode := t.AuthMode
	if authMode == "" {
		authMode = biz.TenantAuthShared
	}

	return tx.QueryRow(ctx,
		`INSERT INTO tenant_runtime_configs (tenant_id, deployment_mode, auth_mode)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (tenant_id) DO UPDATE
		   SET deployment_mode = EXCLUDED.deployment_mode,
		       auth_mode = EXCLUDED.auth_mode,
		       updated_at = now(),
		       deleted_at = NULL
		 RETURNING deployment_mode, auth_mode`,
		t.ID, deploymentMode, authMode,
	).Scan(&t.DeploymentMode, &t.AuthMode)
}

func (r *tenantRepo) bootstrapTenantDefaults(ctx context.Context, tx pgx.Tx, tenantID, ownerID string) error {
	if _, err := tx.Exec(ctx, `
		WITH seed(resource, action, module, scope, risk_level, description) AS (
			VALUES
				('dashboard', 'view',   'shell', 'tenant', 'low',    'View dashboard'),
				('iam',       'view',   'iam',   'tenant', 'low',    'View IAM module'),
				('system',    'view',   'shell', 'tenant', 'low',    'View system module'),
				('tenant',    'read',   'shell', 'tenant', 'low',    'View workspace settings'),
				('tenant',    'create', 'shell', 'global', 'medium', 'Create workspace'),
				('tenant',    'update', 'shell', 'tenant', 'medium', 'Update workspace'),
				('tenant',    'delete', 'shell', 'tenant', 'high',   'Delete workspace'),
				('user',      'read',   'iam',   'tenant', 'low',    'View users'),
				('user',      'create', 'iam',   'tenant', 'medium', 'Create users'),
				('user',      'update', 'iam',   'tenant', 'medium', 'Update users'),
				('user',      'delete', 'iam',   'tenant', 'high',   'Delete users'),
				('role',      'read',   'iam',   'tenant', 'low',    'View roles'),
				('role',      'create', 'iam',   'tenant', 'medium', 'Create roles'),
				('role',      'update', 'iam',   'tenant', 'medium', 'Update roles'),
				('role',      'delete', 'iam',   'tenant', 'high',   'Delete roles'),
				('user-group','read',   'iam',   'tenant', 'low',    'View groups'),
				('user-group','create', 'iam',   'tenant', 'medium', 'Create groups'),
				('user-group','update', 'iam',   'tenant', 'medium', 'Update groups'),
				('user-group','delete', 'iam',   'tenant', 'high',   'Delete groups'),
				('menu',      'read',   'iam',   'tenant', 'low',    'View menu configuration'),
				('menu',      'create', 'iam',   'tenant', 'medium', 'Create menu entries'),
				('menu',      'update', 'iam',   'tenant', 'medium', 'Update menu entries'),
				('menu',      'delete', 'iam',   'tenant', 'high',   'Delete menu entries')
		)
		INSERT INTO permissions (tenant_id, resource, action, code, module, scope, risk_level, description)
		SELECT $1, resource, action, resource || ':' || action, module, scope, risk_level, description
		FROM seed
		ON CONFLICT (tenant_id, resource, action) DO UPDATE
		SET code = EXCLUDED.code,
		    module = EXCLUDED.module,
		    scope = EXCLUDED.scope,
		    risk_level = EXCLUDED.risk_level,
		    description = EXCLUDED.description,
		    status = 'active',
		    updated_at = now()`, tenantID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO roles (tenant_id, name, description, is_system, code, risk_level, status)
		VALUES ($1, 'owner', 'Tenant owner', true, 'owner', 'high', 'active')
		ON CONFLICT (tenant_id, name) WHERE deleted_at IS NULL DO UPDATE
		SET description = EXCLUDED.description,
		    is_system = true,
		    code = EXCLUDED.code,
		    risk_level = EXCLUDED.risk_level,
		    status = 'active',
		    updated_at = now()`, tenantID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		WITH owner_role AS (
			SELECT id FROM roles WHERE tenant_id = $1 AND name = 'owner' AND deleted_at IS NULL LIMIT 1
		)
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT owner_role.id, p.id
		FROM owner_role
		JOIN permissions p ON p.tenant_id = $1
		ON CONFLICT (role_id, permission_id) DO NOTHING`, tenantID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		WITH owner_role AS (
			SELECT id FROM roles WHERE tenant_id = $1 AND name = 'owner' AND deleted_at IS NULL LIMIT 1
		)
		INSERT INTO user_roles (user_id, role_id, tenant_id, assigned_by, status)
		SELECT $2, owner_role.id, $1, $2, 'active'
		FROM owner_role
		ON CONFLICT (user_id, role_id, tenant_id) DO UPDATE
		SET status = 'active',
		    assigned_by = EXCLUDED.assigned_by,
		    updated_at = now()`, tenantID, ownerID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		WITH seed(name, slug, icon, route, sort_order, enabled, permission_slug) AS (
			VALUES
				('Dashboard',       'dashboard',       'pi pi-home',      '/home',       1,  true, 'dashboard:view'),
				('IAM',             'iam',             'pi pi-shield',    NULL,          40, true, 'iam:view'),
				('Users',           'users',           'pi pi-users',     '/iam/users',  41, true, 'user:read'),
				('Roles',           'roles',           'pi pi-lock',      '/iam/roles',  42, true, 'role:read'),
				('Nhóm người dùng', 'user-groups',     'pi pi-users',     '/iam/groups', 43, true, 'user-group:read'),
				('Menus',           'menu-management', 'pi pi-bars',      '/iam/menus',  44, true, 'menu:read'),
				('System',          'system',          'pi pi-cog',       NULL,          90, true, 'system:view'),
				('Workspaces',      'workspaces',      'pi pi-building',  '/workspaces', 91, true, 'tenant:read'),
				('Settings',        'settings',        'pi pi-sliders-h', '/settings',   92, true, 'system:view'),
				('Profile',         'profile',         'pi pi-user',      '/iam/profile',93, true, NULL)
		)
		INSERT INTO menus (tenant_id, name, slug, icon, route, sort_order, enabled, permission_slug)
		SELECT $1, name, slug, icon, route, sort_order, enabled, permission_slug
		FROM seed
		ON CONFLICT (tenant_id, slug) DO UPDATE
		SET name = EXCLUDED.name,
		    icon = EXCLUDED.icon,
		    route = EXCLUDED.route,
		    sort_order = EXCLUDED.sort_order,
		    enabled = EXCLUDED.enabled,
		    permission_slug = EXCLUDED.permission_slug,
		    updated_at = now()`, tenantID); err != nil {
		return err
	}

	_, err := tx.Exec(ctx, `
		WITH links(child_slug, parent_slug) AS (
			VALUES
				('users', 'iam'),
				('roles', 'iam'),
				('user-groups', 'iam'),
				('menu-management', 'iam'),
				('workspaces', 'system'),
				('settings', 'system'),
				('profile', 'system')
		)
		UPDATE menus child
		SET parent_id = parent.id,
		    updated_at = now()
		FROM links, menus parent
		WHERE child.tenant_id = $1
		  AND child.slug = links.child_slug
		  AND parent.tenant_id = $1
		  AND parent.slug = links.parent_slug`, tenantID)
	return err
}
