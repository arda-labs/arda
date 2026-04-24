INSERT INTO users (id, external_id, email, display_name)
VALUES ('00000000-0000-0000-0000-000000000000', 'system@internal', 'system@internal', 'System Owner');

INSERT INTO tenants (id, name, slug, owner_id)
VALUES ('00000000-0000-0000-0000-000000000001', 'System', 'system', '00000000-0000-0000-0000-000000000000');

INSERT INTO permissions (tenant_id, resource, action) VALUES
  ('00000000-0000-0000-0000-000000000001', 'project', 'create'),
  ('00000000-0000-0000-0000-000000000001', 'project', 'read'),
  ('00000000-0000-0000-0000-000000000001', 'project', 'update'),
  ('00000000-0000-0000-0000-000000000001', 'project', 'delete'),
  ('00000000-0000-0000-0000-000000000001', 'member', 'invite'),
  ('00000000-0000-0000-0000-000000000001', 'member', 'read'),
  ('00000000-0000-0000-0000-000000000001', 'member', 'remove'),
  ('00000000-0000-0000-0000-000000000001', 'role', 'create'),
  ('00000000-0000-0000-0000-000000000001', 'role', 'read'),
  ('00000000-0000-0000-0000-000000000001', 'role', 'update'),
  ('00000000-0000-0000-0000-000000000001', 'role', 'delete'),
  ('00000000-0000-0000-0000-000000000001', 'tenant', 'read'),
  ('00000000-0000-0000-0000-000000000001', 'tenant', 'update'),
  ('00000000-0000-0000-0000-000000000001', 'tenant', 'delete');

INSERT INTO roles (id, tenant_id, name, description, is_system)
VALUES ('00000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000001', 'admin', 'Tenant administrator', true);

INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000010', id FROM permissions
WHERE tenant_id = '00000000-0000-0000-0000-000000000001';

INSERT INTO roles (id, tenant_id, name, description, is_system)
VALUES ('00000000-0000-0000-0000-000000000020', '00000000-0000-0000-0000-000000000001', 'member', 'Tenant member', true);

INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000020', p.id
FROM permissions p
WHERE p.tenant_id = '00000000-0000-0000-0000-000000000001'
  AND p.action = 'read'
  AND p.resource IN ('project', 'member', 'role', 'tenant');
