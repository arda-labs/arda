-- Migration: Remove tenant_id from menus (shared menu tree across tenants)
ALTER TABLE menus DROP CONSTRAINT IF EXISTS menus_tenant_id_slug_key;
DROP INDEX IF EXISTS idx_menus_tenant_id;

ALTER TABLE menus DROP COLUMN tenant_id;

CREATE UNIQUE INDEX idx_menus_slug ON menus(slug);
