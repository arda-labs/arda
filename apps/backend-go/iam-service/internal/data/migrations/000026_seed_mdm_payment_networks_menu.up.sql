INSERT INTO menus (id, parent_id, code, name, path, icon, sort_order, status, created_at, updated_at)
SELECT gen_random_uuid(), parent.id, 'mdm-payment-networks', 'Chi nhánh & mạng thanh toán', '/app/mdm/banking/payment-networks', 'pi pi-share-alt', 70, 'ACTIVE', now(), now()
FROM menus parent
WHERE parent.code = 'mdm-banking'
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
    path = EXCLUDED.path,
    icon = EXCLUDED.icon,
    sort_order = EXCLUDED.sort_order,
    status = EXCLUDED.status,
    updated_at = now();
