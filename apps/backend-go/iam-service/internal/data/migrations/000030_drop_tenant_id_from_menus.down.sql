-- Rollback: Restore tenant_id to menus
ALTER TABLE menus ADD COLUMN tenant_id UUID;
DROP INDEX IF EXISTS idx_menus_slug;

UPDATE menus SET tenant_id = '10000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;

ALTER TABLE menus ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_menus_tenant_id ON menus(tenant_id);
ALTER TABLE menus ADD CONSTRAINT menus_tenant_id_slug_key UNIQUE(tenant_id, slug);
