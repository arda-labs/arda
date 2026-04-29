-- Menu management is tenant-scoped administration. Keep the permission model
-- as the source of truth; menus only decide navigation visibility.
WITH menu_actions(action, risk_level, description) AS (
  VALUES
    ('read',   'low',    'View tenant menu configuration'),
    ('create', 'medium', 'Create tenant menu entries'),
    ('update', 'medium', 'Update tenant menu entries'),
    ('delete', 'high',   'Delete tenant menu entries')
)
INSERT INTO permissions (tenant_id, resource, action, code, module, scope, risk_level, description)
SELECT t.id, 'menu', ma.action, 'menu:' || ma.action, 'iam', 'tenant', ma.risk_level, ma.description
FROM tenants t
CROSS JOIN menu_actions ma
WHERE t.deleted_at IS NULL
ON CONFLICT (tenant_id, resource, action) DO UPDATE
SET code = EXCLUDED.code,
    module = EXCLUDED.module,
    scope = EXCLUDED.scope,
    risk_level = EXCLUDED.risk_level,
    description = EXCLUDED.description,
    status = 'active',
    updated_at = now();

INSERT INTO permission_catalog (code, module, resource, action, scope, risk_level, description, status)
VALUES
  ('menu:read',   'iam', 'menu', 'read',   'tenant', 'low',    'View tenant menu configuration', 'active'),
  ('menu:create', 'iam', 'menu', 'create', 'tenant', 'medium', 'Create tenant menu entries', 'active'),
  ('menu:update', 'iam', 'menu', 'update', 'tenant', 'medium', 'Update tenant menu entries', 'active'),
  ('menu:delete', 'iam', 'menu', 'delete', 'tenant', 'high',   'Delete tenant menu entries', 'active')
ON CONFLICT (code) DO UPDATE
SET module = EXCLUDED.module,
    resource = EXCLUDED.resource,
    action = EXCLUDED.action,
    scope = EXCLUDED.scope,
    risk_level = EXCLUDED.risk_level,
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    updated_at = now();

INSERT INTO tenant_permission_entitlements (tenant_id, permission_code, source)
SELECT t.id, pc.code, 'seed'
FROM tenants t
JOIN permission_catalog pc ON pc.code IN ('menu:read', 'menu:create', 'menu:update', 'menu:delete')
WHERE t.deleted_at IS NULL
ON CONFLICT (tenant_id, permission_code) DO UPDATE
SET enabled = true,
    updated_at = now();

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.tenant_id = r.tenant_id AND p.resource = 'menu'
WHERE r.deleted_at IS NULL
  AND lower(r.name) IN ('admin', 'super_admin', 'owner')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
VALUES (
  '10000000-0000-0000-0000-000000000001',
  (
    SELECT id
    FROM menus
    WHERE tenant_id = '10000000-0000-0000-0000-000000000001'
      AND slug = 'iam'
    LIMIT 1
  ),
  'Menus',
  'menu-management',
  'pi pi-bars',
  '/app/iam/menus',
  48,
  true,
  'menu:read'
)
ON CONFLICT (tenant_id, slug) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    route = EXCLUDED.route,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    permission_slug = EXCLUDED.permission_slug,
    updated_at = now();
