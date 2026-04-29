-- Enterprise IAM governance schema.
-- Additive migration: keeps the current RBAC/FGAC runtime behavior intact while
-- preparing for permission catalog, maker-checker, access review, and richer UI.

ALTER TABLE permissions ADD COLUMN IF NOT EXISTS code TEXT;
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS module TEXT NOT NULL DEFAULT 'iam';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS scope TEXT NOT NULL DEFAULT 'tenant';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS risk_level TEXT NOT NULL DEFAULT 'low';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

UPDATE permissions
SET code = resource || ':' || action,
    module = CASE WHEN module = 'iam' THEN resource ELSE module END
WHERE code IS NULL OR code = '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_permissions_tenant_code
    ON permissions (tenant_id, code)
    WHERE code IS NOT NULL AND code <> '';

CREATE INDEX IF NOT EXISTS idx_permissions_module_status
    ON permissions (module, status);

ALTER TABLE roles ADD COLUMN IF NOT EXISTS code TEXT;
ALTER TABLE roles ADD COLUMN IF NOT EXISTS risk_level TEXT NOT NULL DEFAULT 'low';
ALTER TABLE roles ADD COLUMN IF NOT EXISTS approval_required BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE roles ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';
ALTER TABLE roles ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}';

UPDATE roles
SET code = lower(trim(both '_' from regexp_replace(name, '[^a-zA-Z0-9]+', '_', 'g')))
WHERE code IS NULL OR code = '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_tenant_code_active
    ON roles (tenant_id, code)
    WHERE deleted_at IS NULL AND code IS NOT NULL AND code <> '';

CREATE INDEX IF NOT EXISTS idx_roles_tenant_status
    ON roles (tenant_id, status)
    WHERE deleted_at IS NULL;

ALTER TABLE role_permissions ADD COLUMN IF NOT EXISTS effect TEXT NOT NULL DEFAULT 'allow';
ALTER TABLE role_permissions ADD COLUMN IF NOT EXISTS granted_by UUID REFERENCES users(id);
ALTER TABLE role_permissions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE INDEX IF NOT EXISTS idx_role_permissions_permission
    ON role_permissions (permission_id);

ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';
ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS assigned_by UUID REFERENCES users(id);
ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS effective_from TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ;
ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS reason TEXT NOT NULL DEFAULT '';
ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE INDEX IF NOT EXISTS idx_user_roles_active_lookup
    ON user_roles (tenant_id, user_id, status, effective_from, expires_at);

ALTER TABLE resource_permissions ADD COLUMN IF NOT EXISTS effective_from TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE resource_permissions ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ;
ALTER TABLE resource_permissions ADD COLUMN IF NOT EXISTS reason TEXT NOT NULL DEFAULT '';
ALTER TABLE resource_permissions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE INDEX IF NOT EXISTS idx_resource_permissions_active_window
    ON resource_permissions (tenant_id, user_id, resource, action, resource_id, effective_from, expires_at)
    WHERE status = 'active';

ALTER TABLE role_hierarchy ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id);
ALTER TABLE role_hierarchy ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id);
ALTER TABLE role_hierarchy ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now();

UPDATE role_hierarchy rh
SET tenant_id = r.tenant_id
FROM roles r
WHERE rh.parent_role_id = r.id
  AND rh.tenant_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_role_hierarchy_tenant_parent
    ON role_hierarchy (tenant_id, parent_role_id);

ALTER TABLE policies ADD COLUMN IF NOT EXISTS code TEXT;
ALTER TABLE policies ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
ALTER TABLE policies ADD COLUMN IF NOT EXISTS resource TEXT NOT NULL DEFAULT '';
ALTER TABLE policies ADD COLUMN IF NOT EXISTS action TEXT NOT NULL DEFAULT '';
ALTER TABLE policies ADD COLUMN IF NOT EXISTS priority INT NOT NULL DEFAULT 100;
ALTER TABLE policies ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';
ALTER TABLE policies ADD COLUMN IF NOT EXISTS condition_json JSONB NOT NULL DEFAULT '{}';
ALTER TABLE policies ADD COLUMN IF NOT EXISTS effective_from TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE policies ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ;

UPDATE policies
SET code = lower(trim(both '_' from regexp_replace(name, '[^a-zA-Z0-9]+', '_', 'g')))
WHERE code IS NULL OR code = '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_policies_tenant_code
    ON policies (tenant_id, code)
    WHERE code IS NOT NULL AND code <> '';

CREATE INDEX IF NOT EXISTS idx_policies_active_lookup
    ON policies (tenant_id, resource, action, status, priority);

CREATE TABLE IF NOT EXISTS permission_catalog (
    code        TEXT PRIMARY KEY,
    module      TEXT NOT NULL,
    resource    TEXT NOT NULL,
    action      TEXT NOT NULL,
    scope       TEXT NOT NULL DEFAULT 'tenant',
    risk_level  TEXT NOT NULL DEFAULT 'low',
    description TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'active',
    metadata    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO permission_catalog (code, module, resource, action, scope, risk_level, description, status)
SELECT DISTINCT ON (p.code) p.code, p.module, p.resource, p.action, p.scope, p.risk_level, p.description, p.status
FROM permissions p
WHERE p.code IS NOT NULL AND p.code <> ''
ORDER BY p.code, p.created_at
ON CONFLICT (code) DO NOTHING;

CREATE TABLE IF NOT EXISTS tenant_permission_entitlements (
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    permission_code TEXT NOT NULL REFERENCES permission_catalog(code),
    enabled         BOOLEAN NOT NULL DEFAULT true,
    source          TEXT NOT NULL DEFAULT 'seed',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, permission_code)
);

INSERT INTO tenant_permission_entitlements (tenant_id, permission_code)
SELECT DISTINCT p.tenant_id, p.code
FROM permissions p
WHERE p.code IS NOT NULL AND p.code <> ''
ON CONFLICT (tenant_id, permission_code) DO NOTHING;

CREATE TABLE IF NOT EXISTS access_requests (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID REFERENCES tenants(id),
    request_type      TEXT NOT NULL,
    status            TEXT NOT NULL DEFAULT 'draft',
    risk_level        TEXT NOT NULL DEFAULT 'low',
    requester_id      UUID NOT NULL REFERENCES users(id),
    subject_user_id   UUID REFERENCES users(id),
    subject_group_id  UUID REFERENCES groups(id),
    maker_id          UUID REFERENCES users(id),
    checker_id        UUID REFERENCES users(id),
    reason            TEXT NOT NULL DEFAULT '',
    decision_reason   TEXT NOT NULL DEFAULT '',
    metadata          JSONB NOT NULL DEFAULT '{}',
    submitted_at      TIMESTAMPTZ,
    decided_at        TIMESTAMPTZ,
    expires_at        TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT access_requests_request_type_chk
        CHECK (request_type IN ('role_assignment', 'role_permission', 'group_membership', 'group_role', 'resource_exception', 'policy_change', 'tenant_admin_change')),
    CONSTRAINT access_requests_status_chk
        CHECK (status IN ('draft', 'pending', 'approved', 'rejected', 'cancelled', 'expired')),
    CONSTRAINT access_requests_risk_chk
        CHECK (risk_level IN ('low', 'medium', 'high', 'critical'))
);

CREATE INDEX IF NOT EXISTS idx_access_requests_tenant_status
    ON access_requests (tenant_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_access_requests_checker
    ON access_requests (checker_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_access_requests_subject_user
    ON access_requests (subject_user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS access_request_items (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id     UUID NOT NULL REFERENCES access_requests(id) ON DELETE CASCADE,
    operation      TEXT NOT NULL,
    target_type    TEXT NOT NULL,
    target_id      TEXT NOT NULL DEFAULT '',
    role_id        UUID REFERENCES roles(id),
    permission_id  UUID REFERENCES permissions(id),
    group_id       UUID REFERENCES groups(id),
    effect         TEXT NOT NULL DEFAULT 'allow',
    before_state   JSONB NOT NULL DEFAULT '{}',
    after_state    JSONB NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT access_request_items_operation_chk
        CHECK (operation IN ('grant', 'revoke', 'update', 'create', 'delete')),
    CONSTRAINT access_request_items_effect_chk
        CHECK (effect IN ('allow', 'deny'))
);

CREATE INDEX IF NOT EXISTS idx_access_request_items_request
    ON access_request_items (request_id);

CREATE TABLE IF NOT EXISTS access_review_campaigns (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id),
    name           TEXT NOT NULL,
    scope          TEXT NOT NULL DEFAULT 'tenant',
    status         TEXT NOT NULL DEFAULT 'draft',
    reviewer_id    UUID REFERENCES users(id),
    due_at         TIMESTAMPTZ,
    created_by     UUID REFERENCES users(id),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT access_review_campaigns_status_chk
        CHECK (status IN ('draft', 'active', 'completed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_access_review_campaigns_tenant_status
    ON access_review_campaigns (tenant_id, status, created_at DESC);

CREATE TABLE IF NOT EXISTS access_review_items (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id   UUID NOT NULL REFERENCES access_review_campaigns(id) ON DELETE CASCADE,
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    subject_user_id UUID REFERENCES users(id),
    role_id       UUID REFERENCES roles(id),
    group_id      UUID REFERENCES groups(id),
    status        TEXT NOT NULL DEFAULT 'pending',
    decision      TEXT NOT NULL DEFAULT '',
    reviewer_id   UUID REFERENCES users(id),
    reason        TEXT NOT NULL DEFAULT '',
    reviewed_at   TIMESTAMPTZ,
    metadata      JSONB NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT access_review_items_status_chk
        CHECK (status IN ('pending', 'approved', 'revoked', 'exception', 'skipped'))
);

CREATE INDEX IF NOT EXISTS idx_access_review_items_campaign
    ON access_review_items (campaign_id, status);
CREATE INDEX IF NOT EXISTS idx_access_review_items_subject
    ON access_review_items (tenant_id, subject_user_id);

CREATE TABLE IF NOT EXISTS segregation_of_duty_rules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID REFERENCES tenants(id),
    code        TEXT NOT NULL,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    severity    TEXT NOT NULL DEFAULT 'high',
    status      TEXT NOT NULL DEFAULT 'active',
    rule_json   JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT segregation_of_duty_rules_severity_chk
        CHECK (severity IN ('medium', 'high', 'critical')),
    CONSTRAINT segregation_of_duty_rules_status_chk
        CHECK (status IN ('active', 'inactive'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sod_rules_tenant_code
    ON segregation_of_duty_rules (tenant_id, code);
