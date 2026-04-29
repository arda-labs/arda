-- Add workspace update/delete permissions for installations that already ran
-- 000017 before those permissions were added.

WITH seed(resource, action, module, scope, risk_level, description) AS (
  VALUES
    ('tenant', 'update', 'shell', 'tenant', 'medium', 'Update workspace'),
    ('tenant', 'delete', 'shell', 'tenant', 'high',   'Delete workspace')
)
INSERT INTO permissions (tenant_id, resource, action, code, module, scope, risk_level, description)
SELECT t.id, seed.resource, seed.action, seed.resource || ':' || seed.action, seed.module, seed.scope, seed.risk_level, seed.description
FROM tenants t
CROSS JOIN seed
WHERE t.deleted_at IS NULL
ON CONFLICT (tenant_id, resource, action) DO UPDATE
SET code = EXCLUDED.code,
    module = EXCLUDED.module,
    scope = EXCLUDED.scope,
    risk_level = EXCLUDED.risk_level,
    description = EXCLUDED.description,
    status = 'active',
    updated_at = now();

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.tenant_id = r.tenant_id
WHERE r.name = 'owner'
  AND r.deleted_at IS NULL
  AND p.resource = 'tenant'
  AND p.action IN ('update', 'delete')
ON CONFLICT (role_id, permission_id) DO NOTHING;
