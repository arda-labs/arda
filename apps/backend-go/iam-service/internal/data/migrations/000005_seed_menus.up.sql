-- Migration: Create menus table for tenant-based dynamic menu configuration
CREATE TABLE IF NOT EXISTS menus (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    parent_id   UUID REFERENCES menus(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(100) NOT NULL,
    icon        VARCHAR(255),
    route       VARCHAR(500),
    sort_order  INTEGER NOT NULL DEFAULT 0,
    enabled     BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(tenant_id, slug)
);

CREATE INDEX idx_menus_tenant_id ON menus(tenant_id);
CREATE INDEX idx_menus_parent_id ON menus(parent_id);

-- Seed default menus for demo tenant (all disabled by default so permission filter hides them)
INSERT INTO menus (tenant_id, name, slug, icon, route, sort_order, enabled) VALUES
-- Dashboard
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Dashboard', 'dashboard', 'pi pi-home', '/app', 1, true),
-- CRM
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'CRM & Sales', 'crm', 'pi pi-chart-line', NULL, 10, true),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Contacts', 'contacts', 'pi pi-address-book', '/app/crm/contacts', 11, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Leads', 'leads', 'pi pi-user-plus', '/app/crm/leads', 12, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Deals', 'deals', 'pi pi-briefcase', '/app/crm/deals', 13, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Quotations', 'quotations', 'pi pi-file-edit', '/app/crm/quotations', 14, false),
-- HRM
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'HRM', 'hrm', 'pi pi-id-card', NULL, 20, true),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Employees', 'employees', 'pi pi-users', '/app/hrm/employees', 21, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Attendance', 'attendance', 'pi pi-calendar-clock', '/app/hrm/attendance', 22, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Payroll', 'payroll', 'pi pi-wallet', '/app/hrm/payroll', 23, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Recruitment', 'recruitment', 'pi pi-send', '/app/hrm/recruitment', 24, false),
-- Finance
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Finance', 'finance', 'pi pi-dollar', NULL, 30, true),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Invoices', 'invoices', 'pi pi-file', '/app/finance/invoices', 31, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Expenses', 'expenses', 'pi pi-credit-card', '/app/finance/expenses', 32, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Reports', 'reports', 'pi pi-chart-bar', '/app/finance/reports', 33, false),
-- IAM
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'IAM', 'iam', 'pi pi-shield', NULL, 40, true),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Users', 'users', 'pi pi-users', '/app/iam/users', 41, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Members', 'members', 'pi pi-user-edit', '/app/iam/members', 42, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Roles', 'roles', 'pi pi-lock', '/app/iam/roles', 43, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Permissions', 'permissions', 'pi pi-key', '/app/iam/permissions', 44, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Approvals', 'approvals', 'pi pi-check-circle', '/app/iam/approvals', 45, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Tenants', 'tenants', 'pi pi-building', '/app/iam/tenants', 46, false),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Audit Log', 'audit', 'pi pi-history', '/app/iam/audit', 47, false),
-- System
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'System', 'system', 'pi pi-cog', NULL, 50, true),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Settings', 'settings', 'pi pi-sliders-h', '/app/settings', 51, true),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'Profile', 'profile', 'pi pi-user', '/app/profile', 52, true),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'API Logs', 'api-logs', 'pi pi-database', '/app/logs', 53, true);
