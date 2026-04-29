-- Workspace lifecycle belongs to shell, not the tenant-scoped IAM screens.
-- Also bootstrap baseline menu/permission/owner role data for existing tenants.

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
SELECT t.id, seed.resource, seed.action, seed.resource || ':' || seed.action, seed.module, seed.scope, seed.risk_level, seed.description
FROM tenants t
CROSS JOIN seed
WHERE t.deleted_at IS NULL
ON CONFLICT (tenant_id, resource, action) DO UPDATE
SET code = EXCLUDED.code,
    module = EXCLUDED.module,
    scope = EXCLUDED.scope,
    risk_level = EXCLUDED.risk_level,
    description = EXCLUDED.description,
    status = 'active',
    updated_at = now();

INSERT INTO roles (tenant_id, name, description, is_system, code, risk_level, status)
SELECT id, 'owner', 'Tenant owner', true, 'owner', 'high', 'active'
FROM tenants
WHERE deleted_at IS NULL
ON CONFLICT (tenant_id, name) WHERE deleted_at IS NULL DO UPDATE
SET description = EXCLUDED.description,
    is_system = true,
    code = EXCLUDED.code,
    risk_level = EXCLUDED.risk_level,
    status = 'active',
    updated_at = now();

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.tenant_id = r.tenant_id
WHERE r.name = 'owner'
  AND r.deleted_at IS NULL
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id, tenant_id, assigned_by, status)
SELECT t.owner_id, r.id, t.id, t.owner_id, 'active'
FROM tenants t
JOIN roles r ON r.tenant_id = t.id AND r.name = 'owner' AND r.deleted_at IS NULL
WHERE t.deleted_at IS NULL
ON CONFLICT (user_id, role_id, tenant_id) DO UPDATE
SET status = 'active',
    assigned_by = EXCLUDED.assigned_by,
    updated_at = now();

INSERT INTO tenant_users (user_id, tenant_id, username, display_name, role, status)
SELECT
  t.owner_id,
  t.id,
  COALESCE(NULLIF(u.username, ''), NULLIF(u.email, ''), 'user-' || left(u.id::text, 8)),
  COALESCE(NULLIF(u.display_name, ''), NULLIF(u.email, ''), ''),
  'owner',
  'ACTIVE'
FROM tenants t
JOIN users u ON u.id = t.owner_id
WHERE t.deleted_at IS NULL
ON CONFLICT (user_id, tenant_id) WHERE deleted_at IS NULL DO NOTHING;

WITH seed(name, slug, icon, route, sort_order, enabled, permission_slug) AS (
  VALUES
    ('Dashboard',       'dashboard',       'pi pi-home',      '/home',        1, true, 'dashboard:view'),
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
SELECT t.id, seed.name, seed.slug, seed.icon, seed.route, seed.sort_order, seed.enabled, seed.permission_slug
FROM tenants t
CROSS JOIN seed
WHERE t.deleted_at IS NULL
ON CONFLICT (tenant_id, slug) DO UPDATE
SET name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    route = EXCLUDED.route,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    permission_slug = EXCLUDED.permission_slug,
    updated_at = now();

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
WHERE child.slug = links.child_slug
  AND parent.tenant_id = child.tenant_id
  AND parent.slug = links.parent_slug;

UPDATE menus
SET enabled = false,
    updated_at = now()
WHERE slug IN ('members', 'permissions', 'approvals', 'tenants', 'audit', 'api-logs')
  AND route IN ('/iam/members', '/iam/permissions', '/iam/approvals', '/iam/tenants', '/iam/audit', '/logs', '/app/iam/members', '/app/iam/permissions', '/app/iam/approvals', '/app/iam/tenants', '/app/iam/audit', '/app/logs');
