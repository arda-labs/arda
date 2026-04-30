DELETE FROM menus WHERE slug IN ('ntf-operations', 'ntf');
DELETE FROM role_permissions
WHERE permission_id IN (SELECT id FROM permissions WHERE resource = 'ntf');
DELETE FROM tenant_permission_entitlements WHERE permission_code IN ('ntf:read', 'ntf:create', 'ntf:update', 'ntf:approve');
DELETE FROM permission_catalog WHERE code IN ('ntf:read', 'ntf:create', 'ntf:update', 'ntf:approve');
DELETE FROM permissions WHERE resource = 'ntf';
