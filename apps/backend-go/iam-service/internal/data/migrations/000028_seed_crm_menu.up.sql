-- CRM Root Menu
INSERT INTO menus (tenant_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT DISTINCT tenant_id, 'Khách hàng (CRM)', 'crm-root', 'pi pi-users', '', 30, true, 'crm:read'
FROM menus
ON CONFLICT (tenant_id, slug) DO UPDATE
SET name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    updated_at = now();

-- CRM Leaf Menus
WITH groups AS (
  SELECT id, tenant_id FROM menus WHERE slug = 'crm-root'
),
leaves(name, slug, icon, route, sort_order, permission) AS (
  VALUES
    ('Đăng ký mới', 'crm-register', 'pi pi-user-plus', '/app/crm/register/init', 10, 'crm:write'),
    ('Danh sách khách hàng', 'crm-list', 'pi pi-list', '/app/crm/info/customer-list', 20, 'crm:read')
)
INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT g.tenant_id, g.id, l.name, l.slug, l.icon, l.route, l.sort_order, true, l.permission
FROM groups g
JOIN leaves l ON true
ON CONFLICT (tenant_id, slug) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    route = EXCLUDED.route,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    permission_slug = EXCLUDED.permission_slug,
    updated_at = now();
