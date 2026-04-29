-- Metadata that lets a tenant stay on the shared platform by default, then move
-- to dedicated datastore/auth later without changing the tenant-facing API.
-- Store runtime config outside tenants so the migration does not require table
-- ownership on older installations.

CREATE TABLE IF NOT EXISTS tenant_runtime_configs (
    tenant_id        UUID PRIMARY KEY REFERENCES tenants(id),
    deployment_mode TEXT NOT NULL DEFAULT 'SHARED',
    auth_mode       TEXT NOT NULL DEFAULT 'SHARED_AUTH',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT chk_tenant_runtime_configs_deployment_mode CHECK (deployment_mode IN ('SHARED', 'DEDICATED')),
    CONSTRAINT chk_tenant_runtime_configs_auth_mode CHECK (auth_mode IN ('SHARED_AUTH', 'DEDICATED_AUTH'))
);

INSERT INTO tenant_runtime_configs (tenant_id, deployment_mode, auth_mode)
SELECT id, 'SHARED', 'SHARED_AUTH'
FROM tenants
WHERE deleted_at IS NULL
ON CONFLICT (tenant_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS tenant_datastores (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id),
    kind           TEXT NOT NULL DEFAULT 'POSTGRES',
    dsn_secret_ref TEXT NOT NULL DEFAULT '',
    schema_name    TEXT NOT NULL DEFAULT 'public',
    status         TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT chk_tenant_datastores_kind CHECK (kind IN ('POSTGRES')),
    CONSTRAINT chk_tenant_datastores_status CHECK (status IN ('ACTIVE', 'PROVISIONING', 'DISABLED'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_datastores_active_kind
  ON tenant_datastores (tenant_id, kind)
  WHERE deleted_at IS NULL AND status IN ('ACTIVE', 'PROVISIONING');

CREATE INDEX IF NOT EXISTS idx_tenant_datastores_tenant
  ON tenant_datastores (tenant_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS tenant_identity_providers (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id),
    provider          TEXT NOT NULL DEFAULT 'ZITADEL',
    issuer            TEXT NOT NULL DEFAULT '',
    client_id         TEXT NOT NULL DEFAULT '',
    client_secret_ref TEXT NOT NULL DEFAULT '',
    status            TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT chk_tenant_idp_provider CHECK (provider IN ('ZITADEL', 'OIDC', 'SAML', 'AZURE_AD')),
    CONSTRAINT chk_tenant_idp_status CHECK (status IN ('ACTIVE', 'PROVISIONING', 'DISABLED'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_identity_providers_active
  ON tenant_identity_providers (tenant_id, provider, issuer)
  WHERE deleted_at IS NULL AND status IN ('ACTIVE', 'PROVISIONING');

CREATE INDEX IF NOT EXISTS idx_tenant_identity_providers_tenant
  ON tenant_identity_providers (tenant_id)
  WHERE deleted_at IS NULL;
