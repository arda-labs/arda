-- Migration: Add menu permissions and link to admin role
-- This handles permissions in format "resource:action" used by the Menu filtering logic

-- 1. Add new permissions for menu viewing
INSERT INTO permissions (tenant_id, resource, action) VALUES
  ('10000000-0000-0000-0000-000000000001', 'dashboard', 'view'),
  ('10000000-0000-0000-0000-000000000001', 'crm',       'view'),
  ('10000000-0000-0000-0000-000000000001', 'hrm',       'view'),
  ('10000000-0000-0000-0000-000000000001', 'finance',   'view'),
  ('10000000-0000-0000-0000-000000000001', 'iam',       'view'),
  ('10000000-0000-0000-0000-000000000001', 'system',    'view')
ON CONFLICT (tenant_id, resource, action) DO NOTHING;

-- 2. Grant these new permissions to the Admin role of Arda tenant
INSERT INTO role_permissions (role_id, permission_id)
SELECT '10000000-0000-0000-0000-000000000010', id
FROM permissions
WHERE tenant_id = '10000000-0000-0000-0000-000000000001'
  AND action = 'view'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 3. Also ensure user admin is assigned to this role (in case previous migration was skipped or partially failed)
INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES (
  '00000000-0000-0000-0000-000000000099',
  '10000000-0000-0000-0000-000000000010',
  '10000000-0000-0000-0000-000000000001'
)
ON CONFLICT (user_id, role_id, tenant_id) DO NOTHING;
