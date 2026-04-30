WITH groups AS (
  SELECT id, tenant_id FROM menus WHERE slug = 'mdm-banking'
),
leaves(group_slug, name, slug, icon, route, sort_order) AS (
  VALUES
    ('mdm-banking', 'Tổ chức tín dụng', 'mdm-credit-institutions', 'pi pi-building-columns', '/app/mdm/banking/credit-institutions', 20)
)
INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT g.tenant_id, g.id, l.name, l.slug, l.icon, l.route, l.sort_order, true, 'mdm:read'
FROM groups g
JOIN leaves l ON l.group_slug = 'mdm-banking'
ON CONFLICT (tenant_id, slug) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    route = EXCLUDED.route,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    permission_slug = EXCLUDED.permission_slug,
    updated_at = now();
