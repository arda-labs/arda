-- Table for User Groups
CREATE TABLE groups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_groups_tenant ON groups (tenant_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_groups_unique_name ON groups (tenant_id, name) WHERE deleted_at IS NULL;

-- Table for Group Membership (Users in Groups)
CREATE TABLE group_members (
    group_id    UUID NOT NULL REFERENCES groups(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, user_id)
);
CREATE INDEX idx_group_members_user ON group_members (user_id);

-- Table for Group Roles (Roles assigned to Groups)
CREATE TABLE group_roles (
    group_id    UUID NOT NULL REFERENCES groups(id),
    role_id     UUID NOT NULL REFERENCES roles(id),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, role_id)
);
CREATE INDEX idx_group_roles_group ON group_roles (group_id);
CREATE INDEX idx_group_roles_role  ON group_roles (role_id);

-- Add Menu Item for User Groups under IAM
DO $$
DECLARE
    v_tenant_id UUID := '10000000-0000-0000-0000-000000000001'; -- Arda Tenant
    v_iam_id UUID;
BEGIN
    -- Get IAM menu ID
    SELECT id INTO v_iam_id FROM menus WHERE slug = 'iam' AND tenant_id = v_tenant_id;

    IF v_iam_id IS NOT NULL THEN
        INSERT INTO menus (tenant_id, parent_id, name, slug, icon, route, sort_order, enabled, permission_slug)
        VALUES (
            v_tenant_id,
            v_iam_id,
            'Nhóm người dùng',
            'user-groups',
            'pi pi-users-group',
            '/iam/groups',
            48,
            true,
            'user-group:read'
        ) ON CONFLICT (tenant_id, slug) DO NOTHING;
    END IF;

    -- Add basic permissions for groups
    INSERT INTO permissions (tenant_id, resource, action) VALUES
      (v_tenant_id, 'user-group', 'create'),
      (v_tenant_id, 'user-group', 'read'),
      (v_tenant_id, 'user-group', 'update'),
      (v_tenant_id, 'user-group', 'delete')
    ON CONFLICT (tenant_id, resource, action) DO NOTHING;

    -- Grant group permissions to admin role
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT '10000000-0000-0000-0000-000000000010', id
    FROM permissions
    WHERE tenant_id = v_tenant_id AND resource = 'user-group'
    ON CONFLICT DO NOTHING;
END $$;
