DELETE FROM menus
WHERE slug IN (
  'mdm-administrative-units',
  'mdm-area-types',
  'mdm-areas',
  'mdm-code-sets',
  'mdm-code-items',
  'mdm-system-parameters',
  'mdm-banking-reference'
);

DELETE FROM menus
WHERE slug IN ('mdm-geo', 'mdm-reference', 'mdm-banking');

DELETE FROM menus
WHERE slug = 'mdm';

DELETE FROM role_permissions rp
USING permissions p
WHERE rp.permission_id = p.id
  AND p.resource = 'mdm';

DELETE FROM tenant_permission_entitlements
WHERE permission_code IN ('mdm:read', 'mdm:create', 'mdm:update', 'mdm:delete');

DELETE FROM permission_catalog
WHERE code IN ('mdm:read', 'mdm:create', 'mdm:update', 'mdm:delete');

DELETE FROM permissions
WHERE resource = 'mdm';
