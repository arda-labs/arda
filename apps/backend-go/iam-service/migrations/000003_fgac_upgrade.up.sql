-- Thêm bảng policies để hỗ trợ ABAC (Attribute-based Access Control)
CREATE TABLE policies (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID REFERENCES tenants(id),
    name        TEXT NOT NULL,
    effect      TEXT NOT NULL CHECK (effect IN ('allow', 'deny')),
    expression  TEXT, -- Biểu hiện logic (Ví dụ: "resource.amount < 500000000")
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Bảng lưu trữ quan hệ phân cấp của Role (Role Hierarchy)
CREATE TABLE role_hierarchy (
    parent_role_id  UUID NOT NULL REFERENCES roles(id),
    child_role_id   UUID NOT NULL REFERENCES roles(id),
    PRIMARY KEY (parent_role_id, child_role_id)
);

-- Mở rộng bảng resource_permissions để hỗ trợ Maker-Checker metadata
ALTER TABLE resource_permissions ADD COLUMN IF NOT EXISTS maker_id UUID REFERENCES users(id);
ALTER TABLE resource_permissions ADD COLUMN IF NOT EXISTS checker_id UUID REFERENCES users(id);
ALTER TABLE resource_permissions ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'active'; -- active, pending_approval, revoked

-- Index tối ưu cho forward-auth check cực nhanh
CREATE INDEX IF NOT EXISTS idx_fgac_check ON resource_permissions (tenant_id, resource, action, resource_id) WHERE status = 'active';
