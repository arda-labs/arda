-- Keep tenant_users.role and user_roles consistent for users created before
-- role assignment was wired into tenant user creation.

INSERT INTO user_roles (user_id, role_id, tenant_id, status)
SELECT tu.user_id, r.id, tu.tenant_id, 'active'
FROM tenant_users tu
JOIN roles r
  ON r.tenant_id = tu.tenant_id
 AND lower(r.name) = lower(COALESCE(NULLIF(tu.role, ''), 'member'))
 AND r.deleted_at IS NULL
WHERE tu.deleted_at IS NULL
ON CONFLICT (user_id, role_id, tenant_id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id, tenant_id, status)
SELECT tu.user_id, r.id, tu.tenant_id, 'active'
FROM tenant_users tu
JOIN roles r
  ON r.tenant_id = tu.tenant_id
 AND lower(r.name) = 'member'
 AND r.deleted_at IS NULL
WHERE tu.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1
    FROM user_roles ur
    WHERE ur.user_id = tu.user_id
      AND ur.tenant_id = tu.tenant_id
  )
ON CONFLICT (user_id, role_id, tenant_id) DO NOTHING;
