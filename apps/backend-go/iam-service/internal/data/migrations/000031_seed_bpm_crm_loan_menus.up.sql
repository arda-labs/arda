-- Seed permissions, menus, and roles for BPM, CRM (read/write), Loan, and HRM children.
--
-- Context: M00028 (CRM menu), M00029 (BPM menu) relied on tenant_id which
-- M00030 dropped, and Loan was never seeded.  This migration enlists all
-- four modules in the shared (no tenant_id) menu table and creates the
-- supporting permission rows and role assignments.

-----------------------------------------------------------------------
-- 1. Permission catalog (global)
-----------------------------------------------------------------------
INSERT INTO permission_catalog (code, module, resource, action, scope, risk_level, description, status)
VALUES
  ('bpm:read',    'bpm', 'bpm', 'read',    'tenant', 'low',    'View BPM processes and instances',          'active'),
  ('bpm:monitor', 'bpm', 'bpm', 'monitor', 'tenant', 'low',    'Monitor BPM process execution',             'active'),
  ('bpm:admin',   'bpm', 'bpm', 'admin',   'tenant', 'medium', 'Administer BPM configuration and templates', 'active'),
  ('loan:view',   'loan','loan','view',    'tenant', 'low',    'Access loan module',                         'active'),
  ('crm:read',    'crm', 'crm', 'read',    'tenant', 'low',    'Read CRM customer data',                     'active'),
  ('crm:write',   'crm', 'crm', 'write',   'tenant', 'medium', 'Create/update CRM customer data',            'active')
ON CONFLICT (code) DO UPDATE
  SET module      = EXCLUDED.module,
      resource    = EXCLUDED.resource,
      action      = EXCLUDED.action,
      description = EXCLUDED.description,
      status      = EXCLUDED.status,
      updated_at  = now();

-----------------------------------------------------------------------
-- 2. Tenant-scoped permissions (old permissions table – still used by
--    GetUserMenu → ListByTenant in biz/menu.go)
-----------------------------------------------------------------------
INSERT INTO permissions (tenant_id, resource, action, code, module, scope, risk_level, description, status)
SELECT t.id, v.resource, v.action, v.code, v.module, 'tenant'::text, v.risk_level, v.description, 'active'::text
FROM tenants t
CROSS JOIN (VALUES
  ('bpm',  'read',    'bpm:read',    'bpm', 'low',    'View BPM processes and instances'),
  ('bpm',  'monitor', 'bpm:monitor', 'bpm', 'low',    'Monitor BPM process execution'),
  ('bpm',  'admin',   'bpm:admin',   'bpm', 'medium', 'Administer BPM configuration and templates'),
  ('loan', 'view',    'loan:view',   'loan','low',    'Access loan module'),
  ('crm',  'read',    'crm:read',    'crm', 'low',    'Read CRM customer data'),
  ('crm',  'write',   'crm:write',   'crm', 'medium', 'Create/update CRM customer data')
) AS v(resource, action, code, module, risk_level, description)
WHERE t.deleted_at IS NULL
ON CONFLICT (tenant_id, resource, action) DO UPDATE
  SET code        = EXCLUDED.code,
      module      = EXCLUDED.module,
      scope       = EXCLUDED.scope,
      risk_level  = EXCLUDED.risk_level,
      description = EXCLUDED.description,
      status      = 'active'::text,
      updated_at  = now();

-----------------------------------------------------------------------
-- 3. Tenant permission entitlements
-----------------------------------------------------------------------
INSERT INTO tenant_permission_entitlements (tenant_id, permission_code, source)
SELECT t.id, pc.code, 'seed'
FROM tenants t
CROSS JOIN permission_catalog pc
WHERE t.deleted_at IS NULL
  AND pc.code IN ('bpm:read', 'bpm:monitor', 'bpm:admin', 'loan:view', 'crm:read', 'crm:write')
ON CONFLICT (tenant_id, permission_code) DO UPDATE
  SET enabled   = true,
      updated_at = now();

-----------------------------------------------------------------------
-- 4. Assign permissions to admin + owner roles
-----------------------------------------------------------------------
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.tenant_id = r.tenant_id
  AND p.resource IN ('bpm', 'loan', 'crm')
  AND p.action IN ('read', 'monitor', 'admin', 'write', 'view')
WHERE r.deleted_at IS NULL
  AND lower(r.name) IN ('admin', 'owner')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-----------------------------------------------------------------------
-- 5. Menus (shared, no tenant_id)
-----------------------------------------------------------------------

-- 5a. CRM children (under crm-root)
WITH crm_parent AS (
  SELECT id FROM menus WHERE slug = 'crm-root'
), crm_leaves(slug, name, icon, route, sort_order, permission_slug) AS (
  VALUES
    ('crm-register', 'Đăng ký mới',           'pi pi-user-plus',  '/crm/register/init',      10, 'crm:write'),
    ('crm-list',     'Danh sách khách hàng',   'pi pi-list',       '/crm/info/customer-list', 20, 'crm:read')
)
INSERT INTO menus (parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT cp.id, l.name, l.slug, l.icon, l.route, l.sort_order, true, l.permission_slug
FROM crm_parent cp, crm_leaves l
ON CONFLICT (slug) DO UPDATE
  SET parent_id      = (SELECT id FROM menus WHERE slug = 'crm-root'),
      name           = EXCLUDED.name,
      icon           = EXCLUDED.icon,
      route          = EXCLUDED.route,
      sort_order     = EXCLUDED.sort_order,
      enabled        = EXCLUDED.enabled,
      permission_slug = EXCLUDED.permission_slug,
      updated_at     = now();

-- 5b. BPM children (under bpm-root)
WITH bpm_parent AS (
  SELECT id FROM menus WHERE slug = 'bpm-root'
), bpm_leaves(slug, name, icon, route, sort_order, permission_slug) AS (
  VALUES
    ('bpm-inbound',       'Giao dịch đến',         'pi pi-download',   '/bpm/inbound',            10, 'bpm:read'),
    ('bpm-outbound',      'Giao dịch đi',          'pi pi-upload',     '/bpm/outbound',           20, 'bpm:read'),
    ('bpm-monitor',       'Giám sát vận hành',    'pi pi-chart-bar',  '/bpm/monitor',            30, 'bpm:monitor'),
    ('bpm-search',        'Tra cứu giao dịch',    'pi pi-search',     '/bpm/search',             40, 'bpm:read'),
    ('bpm-error',         'Xử lý lỗi (Hospital)', 'pi pi-heart-fill', '/bpm/error-hospital',     60, 'bpm:admin')
)
INSERT INTO menus (parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT bp.id, l.name, l.slug, l.icon, l.route, l.sort_order, true, l.permission_slug
FROM bpm_parent bp, bpm_leaves l
ON CONFLICT (slug) DO UPDATE
  SET parent_id      = (SELECT id FROM menus WHERE slug = 'bpm-root'),
      name           = EXCLUDED.name,
      icon           = EXCLUDED.icon,
      route          = EXCLUDED.route,
      sort_order     = EXCLUDED.sort_order,
      enabled        = EXCLUDED.enabled,
      permission_slug = EXCLUDED.permission_slug,
      updated_at     = now();

-- 5c. BPM config parent (under bpm-root)
INSERT INTO menus (parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT id, 'Cấu hình hệ thống', 'bpm-config', 'pi pi-cog', '', 50, true, 'bpm:admin'
FROM menus WHERE slug = 'bpm-root'
ON CONFLICT (slug) DO UPDATE
  SET parent_id      = (SELECT id FROM menus WHERE slug = 'bpm-root'),
      name           = EXCLUDED.name,
      icon           = EXCLUDED.icon,
      sort_order     = EXCLUDED.sort_order,
      enabled        = EXCLUDED.enabled,
      permission_slug = EXCLUDED.permission_slug,
      updated_at     = now();

-- 5d. BPM config leaf children (under bpm-config)
WITH cfg_parent AS (
  SELECT id FROM menus WHERE slug = 'bpm-config'
), cfg_leaves(slug, name, icon, route, sort_order, permission_slug) AS (
  VALUES
    ('bpm-cfg-assignment', 'Quy tắc chia bài',    'pi pi-users',   '/bpm/config/assignment', 10, 'bpm:admin'),
    ('bpm-cfg-sla',        'Cấu hình SLA',        'pi pi-clock',   '/bpm/config/sla',        20, 'bpm:admin'),
    ('bpm-cfg-desc',       'Cấu trúc diễn giải', 'pi pi-comment', '/bpm/config/description', 30, 'bpm:admin')
)
INSERT INTO menus (parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT cp.id, l.name, l.slug, l.icon, l.route, l.sort_order, true, l.permission_slug
FROM cfg_parent cp, cfg_leaves l
ON CONFLICT (slug) DO UPDATE
  SET parent_id      = (SELECT id FROM menus WHERE slug = 'bpm-config'),
      name           = EXCLUDED.name,
      icon           = EXCLUDED.icon,
      route          = EXCLUDED.route,
      sort_order     = EXCLUDED.sort_order,
      enabled        = EXCLUDED.enabled,
      permission_slug = EXCLUDED.permission_slug,
      updated_at     = now();

-- 5e. Loan root + leaf
INSERT INTO menus (name, slug, icon, route, sort_order, enabled, permission_slug)
VALUES ('Khoản vay (Loan)', 'loan', 'pi pi-money-bill', '', 50, true, 'loan:view')
ON CONFLICT (slug) DO UPDATE
  SET name           = EXCLUDED.name,
      icon           = EXCLUDED.icon,
      route          = EXCLUDED.route,
      sort_order     = EXCLUDED.sort_order,
      enabled        = EXCLUDED.enabled,
      permission_slug = EXCLUDED.permission_slug,
      updated_at     = now();

INSERT INTO menus (parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT id, 'Hồ sơ vay', 'loan-application', 'pi pi-file', '/loan/application', 10, true, 'loan:view'
FROM menus WHERE slug = 'loan'
ON CONFLICT (slug) DO UPDATE
  SET parent_id      = (SELECT id FROM menus WHERE slug = 'loan'),
      name           = EXCLUDED.name,
      icon           = EXCLUDED.icon,
      route          = EXCLUDED.route,
      sort_order     = EXCLUDED.sort_order,
      enabled        = EXCLUDED.enabled,
      permission_slug = EXCLUDED.permission_slug,
      updated_at     = now();

-- 5f. HRM children (employee sub-menus were migrated as standalone roots
--     from M00005 and kept disabled.  Add proper HRM children here.)
INSERT INTO menus (parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
SELECT id, 'Onboarding', 'hrm-onboarding', 'pi pi-user-plus', '/hrm/onboarding', 10, true, 'hrm:view'
FROM menus WHERE slug = 'hrm'
ON CONFLICT (slug) DO UPDATE
  SET parent_id      = (SELECT id FROM menus WHERE slug = 'hrm'),
      name           = EXCLUDED.name,
      icon           = EXCLUDED.icon,
      route          = EXCLUDED.route,
      sort_order     = EXCLUDED.sort_order,
      enabled        = EXCLUDED.enabled,
      permission_slug = EXCLUDED.permission_slug,
      updated_at     = now();

-- 5g. Remove conflicting standalone root menu "crm" (CRM & Sales, slug = crm)
--     now that we have crm-root with proper children.  Disable the old one
--     to avoid confusion in role-permission check.
UPDATE menus
SET enabled = false, updated_at = now()
WHERE slug = 'crm' AND enabled = true;
