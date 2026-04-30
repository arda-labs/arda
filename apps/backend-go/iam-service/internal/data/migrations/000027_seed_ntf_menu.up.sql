WITH ntf_actions(action, risk_level, description) AS (
  VALUES
    ('read',   'low',    'View notification templates, queue, and providers'),
    ('create', 'medium', 'Create notification templates and provider configs'),
    ('update', 'medium', 'Update notification templates, versions, and provider configs'),
    ('approve','high',   'Approve customer-facing notification template versions')
)
INSERT INTO permissions (tenant_id, resource, action, code, module, scope, risk_level, description)
SELECT t.id, 'ntf', na.action, 'ntf:' || na.action, 'ntf', 'tenant', na.risk_level, na.description
FROM tenants t
CROSS JOIN ntf_actions na
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
  ('ntf:read',    'ntf', 'ntf', 'read',    'tenant', 'low',    'View notification templates, queue, and providers', 'active'),
  ('ntf:create',  'ntf', 'ntf', 'create',  'tenant', 'medium', 'Create notification templates and provider configs', 'active'),
  ('ntf:update',  'ntf', 'ntf', 'update',  'tenant', 'medium', 'Update notification templates, versions, and provider configs', 'active'),
  ('ntf:approve', 'ntf', 'ntf', 'approve', 'tenant', 'high',   'Approve customer-facing notification template versions', 'active')
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
JOIN permission_catalog pc ON pc.code IN ('ntf:read', 'ntf:create', 'ntf:update', 'ntf:approve')
WHERE t.deleted_at IS NULL
ON CONFLICT (tenant_id, permission_code) DO UPDATE
SET enabled = true,
    updated_at = now();

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.tenant_id = r.tenant_id AND p.resource = 'ntf'
WHERE r.deleted_at IS NULL
  AND lower(r.name) IN ('admin', 'super_admin', 'owner')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT t.id, NULL, 'Thông báo', 'ntf', 'pi pi-bell', NULL, 45, true, 'ntf:read'
FROM tenants t
WHERE t.deleted_at IS NULL
  AND t.slug <> 'system'
ON CONFLICT (tenant_id, slug) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    route = EXCLUDED.route,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    permission_slug = EXCLUDED.permission_slug,
    updated_at = now();

WITH ntf_parent AS (
  SELECT id, tenant_id FROM menus WHERE slug = 'ntf'
),
leaves(name, slug, icon, route, sort_order) AS (
  VALUES
    ('Vận hành thông báo', 'ntf-operations', 'pi pi-send', '/app/ntf/operations', 10)
)
INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT p.tenant_id, p.id, l.name, l.slug, l.icon, l.route, l.sort_order, true, 'ntf:read'
FROM ntf_parent p
CROSS JOIN leaves l
ON CONFLICT (tenant_id, slug) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    route = EXCLUDED.route,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    permission_slug = EXCLUDED.permission_slug,
    updated_at = now();
