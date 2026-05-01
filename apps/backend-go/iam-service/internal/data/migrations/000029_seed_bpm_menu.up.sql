-- BPM Root Menu
INSERT INTO menus (tenant_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT DISTINCT tenant_id, 'Quy trình (BPM)', 'bpm-root', 'pi pi-sitemap', '', 40, true, 'bpm:read'
FROM menus
ON CONFLICT (tenant_id, slug) DO UPDATE
SET name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    updated_at = now();

-- BPM Top-level Leaf Menus
WITH groups AS (
  SELECT id, tenant_id FROM menus WHERE slug = 'bpm-root'
),
leaves(name, slug, icon, route, sort_order, permission) AS (
  VALUES
    ('Giao dịch đến', 'bpm-inbound', 'pi pi-download', '/app/bpm/inbound', 10, 'bpm:read'),
    ('Giao dịch đi', 'bpm-outbound', 'pi pi-upload', '/app/bpm/outbound', 20, 'bpm:read'),
    ('Giám sát vận hành', 'bpm-monitor', 'pi pi-chart-bar', '/app/bpm/monitor', 30, 'bpm:monitor'),
    ('Tra cứu giao dịch', 'bpm-search', 'pi pi-search', '/app/bpm/search', 40, 'bpm:read'),
    ('Xử lý lỗi (Hospital)', 'bpm-error', 'pi pi-heart-fill', '/app/bpm/error-hospital', 60, 'bpm:admin')
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

-- BPM Configuration Sub-menu
INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT g.tenant_id, g.id, 'Cấu hình hệ thống', 'bpm-config', 'pi pi-cog', '', 50, true, 'bpm:admin'
FROM menus g WHERE g.slug = 'bpm-root'
ON CONFLICT (tenant_id, slug) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled,
    updated_at = now();

-- BPM Configuration Leaf Menus
WITH groups AS (
  SELECT id, tenant_id FROM menus WHERE slug = 'bpm-config'
),
leaves(name, slug, icon, route, sort_order, permission) AS (
  VALUES
    ('Quy tắc chia bài', 'bpm-cfg-assignment', 'pi pi-users', '/app/bpm/config/assignment', 10, 'bpm:admin'),
    ('Cấu hình SLA', 'bpm-cfg-sla', 'pi pi-clock', '/app/bpm/config/sla', 20, 'bpm:admin'),
    ('Cấu trúc diễn giải', 'bpm-cfg-desc', 'pi pi-comment', '/app/bpm/config/description', 30, 'bpm:admin')
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
