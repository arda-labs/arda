WITH mdm_actions(action, risk_level, description) AS (
  VALUES
    ('read',   'low',    'View MDM master data'),
    ('create', 'medium', 'Create MDM master data'),
    ('update', 'medium', 'Update MDM master data'),
    ('delete', 'high',   'Delete MDM master data')
)
INSERT INTO permissions (tenant_id, resource, action, code, module, scope, risk_level, description)
SELECT t.id, 'mdm', ma.action, 'mdm:' || ma.action, 'mdm', 'tenant', ma.risk_level, ma.description
FROM tenants t
CROSS JOIN mdm_actions ma
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
  ('mdm:read',   'mdm', 'mdm', 'read',   'tenant', 'low',    'View MDM master data', 'active'),
  ('mdm:create', 'mdm', 'mdm', 'create', 'tenant', 'medium', 'Create MDM master data', 'active'),
  ('mdm:update', 'mdm', 'mdm', 'update', 'tenant', 'medium', 'Update MDM master data', 'active'),
  ('mdm:delete', 'mdm', 'mdm', 'delete', 'tenant', 'high',   'Delete MDM master data', 'active')
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
JOIN permission_catalog pc ON pc.code IN ('mdm:read', 'mdm:create', 'mdm:update', 'mdm:delete')
WHERE t.deleted_at IS NULL
ON CONFLICT (tenant_id, permission_code) DO UPDATE
SET enabled = true,
    updated_at = now();

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.tenant_id = r.tenant_id AND p.resource = 'mdm'
WHERE r.deleted_at IS NULL
  AND lower(r.name) IN ('admin', 'super_admin', 'owner')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT t.id, NULL, 'MDM', 'mdm', 'pi pi-database', NULL, 35, true, 'mdm:read'
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

WITH mdm_parent AS (
  SELECT id, tenant_id FROM menus WHERE slug = 'mdm'
),
groups(name, slug, icon, sort_order) AS (
  VALUES
    ('Địa lý hành chính', 'mdm-geo', 'pi pi-map', 10),
    ('Danh mục & tham số', 'mdm-reference', 'pi pi-list-check', 20),
    ('Tài chính ngân hàng', 'mdm-banking', 'pi pi-building-columns', 30)
)
INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT p.tenant_id, p.id, g.name, g.slug, g.icon, NULL, g.sort_order, true, 'mdm:read'
FROM mdm_parent p
CROSS JOIN groups g
ON CONFLICT (tenant_id, slug) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    route = EXCLUDED.route,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    permission_slug = EXCLUDED.permission_slug,
    updated_at = now();

WITH groups AS (
  SELECT id, tenant_id, slug FROM menus WHERE slug IN ('mdm-geo', 'mdm-reference', 'mdm-banking')
),
leaves(group_slug, name, slug, icon, route, sort_order) AS (
  VALUES
    ('mdm-geo',       'Tỉnh, phường/xã',    'mdm-administrative-units', 'pi pi-map-marker',      '/app/mdm/geo/administrative-units', 10),
    ('mdm-geo',       'Loại khu vực',       'mdm-area-types',           'pi pi-tags',            '/app/mdm/geo/area-types',          20),
    ('mdm-geo',       'Khu vực quản lý',    'mdm-areas',                'pi pi-sitemap',         '/app/mdm/geo/areas',               30),
    ('mdm-reference', 'Bộ danh mục',        'mdm-code-sets',            'pi pi-list-check',      '/app/mdm/catalog/code-sets',       10),
    ('mdm-reference', 'Giá trị danh mục',   'mdm-code-items',           'pi pi-objects-column',  '/app/mdm/catalog/code-items',      20),
    ('mdm-reference', 'Tham số hệ thống',   'mdm-system-parameters',    'pi pi-sliders-h',       '/app/mdm/system/parameters',       30),
    ('mdm-banking',   'Gợi ý mở rộng',      'mdm-banking-reference',    'pi pi-lightbulb',       '/app/mdm/banking/reference',       10)
)
INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT g.tenant_id, g.id, l.name, l.slug, l.icon, l.route, l.sort_order, true, 'mdm:read'
FROM groups g
JOIN leaves l ON l.group_slug = g.slug
ON CONFLICT (tenant_id, slug) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    route = EXCLUDED.route,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    permission_slug = EXCLUDED.permission_slug,
    updated_at = now();
