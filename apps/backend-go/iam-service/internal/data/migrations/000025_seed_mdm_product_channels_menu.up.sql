INSERT INTO menus (id, parent_id, code, name, path, icon, sort_order, status, created_at, updated_at)
SELECT gen_random_uuid(), parent.id, 'mdm-product-channels', 'Sản phẩm & kênh', '/app/mdm/banking/product-channels', 'pi pi-sitemap', 60, 'ACTIVE', now(), now()
FROM menus parent
WHERE parent.code = 'mdm-banking'
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
    path = EXCLUDED.path,
    icon = EXCLUDED.icon,
    sort_order = EXCLUDED.sort_order,
    status = EXCLUDED.status,
    updated_at = now();
