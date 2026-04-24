-- Migration: seed admin user + Arda tenant với full permissions
-- Chạy sau khi admin đăng nhập lần đầu để có membership

-- 1. Upsert user admin từ Zitadel
-- external_id phải khớp với `sub` claim trong JWT của Zitadel
-- Thay 'admin-zitadel-sub-id' bằng giá trị sub thực từ Zitadel
INSERT INTO users (id, external_id, email, display_name)
VALUES (
  '00000000-0000-0000-0000-000000000099',
  '369593749817000033',         -- Zitadel sub của admin@zitadel.auth.arda.io.vn
  'admin@zitadel.auth.arda.io.vn',
  'Arda Admin'
)
ON CONFLICT (external_id) DO UPDATE
  SET email = EXCLUDED.email,
      display_name = EXCLUDED.display_name;

-- 2. Tạo tenant Arda chính
INSERT INTO tenants (id, name, slug, owner_id)
VALUES (
  '10000000-0000-0000-0000-000000000001',
  'Arda',
  'arda',
  '00000000-0000-0000-0000-000000000099'
)
ON CONFLICT (slug) DO NOTHING;

-- 3. Tạo role admin cho tenant Arda với TẤT CẢ permissions
INSERT INTO permissions (tenant_id, resource, action) VALUES
  ('10000000-0000-0000-0000-000000000001', 'user',       'read'),
  ('10000000-0000-0000-0000-000000000001', 'user',       'create'),
  ('10000000-0000-0000-0000-000000000001', 'user',       'update'),
  ('10000000-0000-0000-0000-000000000001', 'user',       'delete'),
  ('10000000-0000-0000-0000-000000000001', 'member',     'read'),
  ('10000000-0000-0000-0000-000000000001', 'member',     'invite'),
  ('10000000-0000-0000-0000-000000000001', 'member',     'remove'),
  ('10000000-0000-0000-0000-000000000001', 'role',       'read'),
  ('10000000-0000-0000-0000-000000000001', 'role',       'create'),
  ('10000000-0000-0000-0000-000000000001', 'role',       'update'),
  ('10000000-0000-0000-0000-000000000001', 'role',       'delete'),
  ('10000000-0000-0000-0000-000000000001', 'permission', 'read'),
  ('10000000-0000-0000-0000-000000000001', 'permission', 'grant'),
  ('10000000-0000-0000-0000-000000000001', 'approval',   'read'),
  ('10000000-0000-0000-0000-000000000001', 'approval',   'manage'),
  ('10000000-0000-0000-0000-000000000001', 'tenant',     'read'),
  ('10000000-0000-0000-0000-000000000001', 'tenant',     'update'),
  ('10000000-0000-0000-0000-000000000001', 'tenant',     'delete'),
  ('10000000-0000-0000-0000-000000000001', 'audit',      'read'),
  ('10000000-0000-0000-0000-000000000001', 'settings',   'read'),
  ('10000000-0000-0000-0000-000000000001', 'settings',   'update'),
  ('10000000-0000-0000-0000-000000000001', 'lead',       'read'),
  ('10000000-0000-0000-0000-000000000001', 'lead',       'create'),
  ('10000000-0000-0000-0000-000000000001', 'deal',       'read'),
  ('10000000-0000-0000-0000-000000000001', 'deal',       'create'),
  ('10000000-0000-0000-0000-000000000001', 'contact',    'read'),
  ('10000000-0000-0000-0000-000000000001', 'contact',    'create'),
  ('10000000-0000-0000-0000-000000000001', 'employee',   'read'),
  ('10000000-0000-0000-0000-000000000001', 'employee',   'create'),
  ('10000000-0000-0000-0000-000000000001', 'payroll',    'read'),
  ('10000000-0000-0000-0000-000000000001', 'attendance', 'read'),
  ('10000000-0000-0000-0000-000000000001', 'invoice',    'read'),
  ('10000000-0000-0000-0000-000000000001', 'expense',    'read'),
  ('10000000-0000-0000-0000-000000000001', 'report',     'read'),
  ('10000000-0000-0000-0000-000000000001', 'quote',      'read'),
  ('10000000-0000-0000-0000-000000000001', 'me',         'read'),
  ('10000000-0000-0000-0000-000000000001', 'public',     'access')
ON CONFLICT (tenant_id, resource, action) DO NOTHING;

-- 4. Tạo role admin cho tenant Arda
INSERT INTO roles (id, tenant_id, name, description, is_system)
VALUES (
  '10000000-0000-0000-0000-000000000010',
  '10000000-0000-0000-0000-000000000001',
  'admin',
  'Full access administrator',
  true
)
ON CONFLICT (tenant_id, name) WHERE deleted_at IS NULL DO NOTHING;

-- 5. Gán tất cả permissions cho role admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT '10000000-0000-0000-0000-000000000010', id
FROM permissions
WHERE tenant_id = '10000000-0000-0000-0000-000000000001'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 6. Tạo role member (read-only) cho tenant Arda
INSERT INTO roles (id, tenant_id, name, description, is_system)
VALUES (
  '10000000-0000-0000-0000-000000000020',
  '10000000-0000-0000-0000-000000000001',
  'member',
  'Standard member - read only',
  true
)
ON CONFLICT (tenant_id, name) WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT '10000000-0000-0000-0000-000000000020', id
FROM permissions
WHERE tenant_id = '10000000-0000-0000-0000-000000000001'
  AND action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 7. Thêm admin vào tenant Arda
INSERT INTO memberships (user_id, tenant_id, role)
VALUES (
  '00000000-0000-0000-0000-000000000099',
  '10000000-0000-0000-0000-000000000001',
  'admin'
)
ON CONFLICT (user_id, tenant_id) WHERE deleted_at IS NULL DO NOTHING;

-- 8. Gán role admin cho user trong tenant Arda
INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES (
  '00000000-0000-0000-0000-000000000099',
  '10000000-0000-0000-0000-000000000010',
  '10000000-0000-0000-0000-000000000001'
)
ON CONFLICT (user_id, role_id, tenant_id) DO NOTHING;
