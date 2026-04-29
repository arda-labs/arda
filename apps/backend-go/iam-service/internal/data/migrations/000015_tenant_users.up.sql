-- Move tenant-specific account data out of the global users table.
-- Keep the legacy memberships/users columns intact so the migration works when
-- the service user can create new objects but is not owner of old tables.

CREATE TABLE IF NOT EXISTS tenant_users (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id),
    tenant_id    UUID NOT NULL REFERENCES tenants(id),
    username     TEXT NOT NULL,
    display_name TEXT NOT NULL DEFAULT '',
    role         TEXT NOT NULL DEFAULT 'member',
    status       TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT chk_tenant_users_status CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED'))
);

INSERT INTO tenant_users (
    id,
    user_id,
    tenant_id,
    username,
    display_name,
    role,
    status,
    created_at,
    updated_at,
    deleted_at
)
SELECT
    m.id,
    m.user_id,
    m.tenant_id,
    COALESCE(NULLIF(u.username, ''), 'user-' || left(u.id::text, 8)),
    COALESCE(NULLIF(u.display_name, ''), NULLIF(u.email, ''), NULLIF(u.username, ''), ''),
    COALESCE(NULLIF(m.role, ''), 'member'),
    'ACTIVE',
    m.created_at,
    now(),
    m.deleted_at
FROM memberships m
JOIN users u ON u.id = m.user_id
ON CONFLICT (id) DO UPDATE SET
    username = EXCLUDED.username,
    display_name = EXCLUDED.display_name,
    role = EXCLUDED.role,
    status = EXCLUDED.status,
    updated_at = now(),
    deleted_at = EXCLUDED.deleted_at;

CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_users_unique_active
  ON tenant_users (user_id, tenant_id)
  WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_users_username_unique_active
  ON tenant_users (tenant_id, lower(username))
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tenant_users_user
  ON tenant_users (user_id)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tenant_users_tenant
  ON tenant_users (tenant_id)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tenant_users_status_active
  ON tenant_users (tenant_id, status)
  WHERE deleted_at IS NULL;
