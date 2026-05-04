-- Rollback menus, permissions, and entitlements seeded in M00031.

-----------------------------------------------------------------------
-- 1. Remove seeded menus (children + loan root)
-----------------------------------------------------------------------
DELETE FROM menus WHERE slug IN (
  'crm-register', 'crm-list',
  'bpm-inbound', 'bpm-outbound', 'bpm-monitor', 'bpm-search', 'bpm-error',
  'bpm-config', 'bpm-cfg-assignment', 'bpm-cfg-sla', 'bpm-cfg-desc',
  'loan', 'loan-application', 'hrm-onboarding'
);

-- Re-enable old CRM root
UPDATE menus SET enabled = true, updated_at = now()
WHERE slug = 'crm' AND enabled = false;

-----------------------------------------------------------------------
-- 2. Remove role_permissions for bpm/loan/crm-resource permissions
-----------------------------------------------------------------------
DELETE FROM role_permissions
WHERE permission_id IN (
  SELECT p.id FROM permissions p
  WHERE p.resource IN ('bpm', 'loan')
     OR (p.resource = 'crm' AND p.action IN ('read', 'write'))
);

-----------------------------------------------------------------------
-- 3. Remove tenant-permission entitlements
-----------------------------------------------------------------------
DELETE FROM tenant_permission_entitlements
WHERE permission_code IN ('bpm:read', 'bpm:monitor', 'bpm:admin', 'loan:view', 'crm:read', 'crm:write');

-----------------------------------------------------------------------
-- 4. Remove tenant-scoped permissions (old permissions table)
-----------------------------------------------------------------------
DELETE FROM permissions
WHERE resource IN ('bpm', 'loan')
   OR (resource = 'crm' AND action IN ('read', 'write'));

-----------------------------------------------------------------------
-- 5. Remove global permission-catalog entries
-----------------------------------------------------------------------
DELETE FROM permission_catalog
WHERE code IN ('bpm:read', 'bpm:monitor', 'bpm:admin', 'loan:view', 'crm:read', 'crm:write');
