-- Migration: Force full permissions for the primary admin user
-- User: admin@zitadel.auth.arda.io.vn (sub: 369593749817000033)

-- 1. Đảm bảo User tồn tại với đúng thông tin
INSERT INTO users (id, external_id, email, display_name)
VALUES (
  '00000000-0000-0000-0000-000000000099',
  '369593749817000033',
  'admin@zitadel.auth.arda.io.vn',
  'ZITADEL Admin'
)
ON CONFLICT (external_id) DO UPDATE
SET email = EXCLUDED.email, display_name = EXCLUDED.display_name;

-- 2. Đảm bảo Tenant 'arda' tồn tại
INSERT INTO tenants (id, name, slug, owner_id)
VALUES (
  '10000000-0000-0000-0000-000000000001',
  'Arda',
  'arda',
  '00000000-0000-0000-0000-000000000099'
)
ON CONFLICT (slug) DO NOTHING;

-- 3. Gán tất cả các permission hiện có cho Tenant Arda
-- (Bao gồm các quyền IAM, CRM, HRM và Menu View mới thêm)
INSERT INTO role_permissions (role_id, permission_id)
SELECT '10000000-0000-0000-0000-000000000010', id
FROM permissions
WHERE tenant_id = '10000000-0000-0000-0000-000000000001'
ON CONFLICT DO NOTHING;

-- 4. Đảm bảo User là Member (Admin role) của Tenant
INSERT INTO memberships (user_id, tenant_id, role)
VALUES (
  '00000000-0000-0000-0000-000000000099',
  '10000000-0000-0000-0000-000000000001',
  'admin'
)
ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = 'admin';

-- 5. Gán Role Admin thực thụ cho User
INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES (
  '00000000-0000-0000-0000-000000000099',
  '10000000-0000-0000-0000-000000000010',
  '10000000-0000-0000-0000-000000000001'
)
ON CONFLICT DO NOTHING;
