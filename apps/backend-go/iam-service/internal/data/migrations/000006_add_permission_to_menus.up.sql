-- Migration: Add permission_slug to menus table
ALTER TABLE menus ADD COLUMN IF NOT EXISTS permission_slug VARCHAR(100);

-- Update existing menus with some default permissions for testing
UPDATE menus SET permission_slug = 'dashboard:view' WHERE slug = 'dashboard';
UPDATE menus SET permission_slug = 'crm:view' WHERE slug = 'crm';
UPDATE menus SET permission_slug = 'hrm:view' WHERE slug = 'hrm';
UPDATE menus SET permission_slug = 'finance:view' WHERE slug = 'finance';
UPDATE menus SET permission_slug = 'iam:view' WHERE slug = 'iam';
UPDATE menus SET permission_slug = 'system:view' WHERE slug = 'system';
