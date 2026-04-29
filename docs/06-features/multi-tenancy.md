# Multi-Tenancy Model

> Cap nhat: 2026-04-29
> Pham vi: tenant data isolation, tenant account, shared/dedicated deployment boundary.

## Muc Tieu

Arda mac dinh chay theo mo hinh SaaS shared platform, nhung phai co boundary de sau nay tach rieng DB hoac auth cho enterprise tenant ma khong phai sua lai toan bo business code.

Thiet ke hien tai tach 4 lop:

| Lop | Vai tro |
| --- | --- |
| `users` | Global identity: external ID, email, display name. Khong chua tenant username. |
| `tenant_users` | Tenant account: username, display name, role, status trong tung tenant. |
| Business data | Du lieu nghiep vu co `tenant_id NOT NULL`. |
| Deployment metadata | `tenant_runtime_configs`, `tenant_datastores`, `tenant_identity_providers`. |

## Identity Va Tenant Account

`users` la identity dung chung de login/SSO/audit. `tenant_users` la tai khoan cua identity do trong mot tenant cu the.

Vi du cung mot email co the xuat hien o nhieu tenant voi username khac nhau:

```text
users
  id = U1
  email = minh@example.com

tenant_users
  tenant A: user_id = U1, username = minh, role = ADMIN
  tenant B: user_id = U1, username = m.nguyen, role = VIEWER
```

Username chi unique trong tenant:

```sql
UNIQUE (tenant_id, lower(username)) WHERE deleted_at IS NULL
```

## Tenant Deployment Mode

`tenant_runtime_configs.deployment_mode` co 2 gia tri:

| Mode | Y nghia |
| --- | --- |
| `SHARED` | Tenant dung chung platform DB/schema. Tach du lieu bang `tenant_id` va RLS. |
| `DEDICATED` | Tenant co DB/schema rieng. Hien tai moi la metadata va boundary, chua provisioning DB rieng. |

`tenant_runtime_configs.auth_mode` co 2 gia tri:

| Mode | Y nghia |
| --- | --- |
| `SHARED_AUTH` | Tenant dung auth chung. |
| `DEDICATED_AUTH` | Tenant co IdP/Zitadel/OIDC/SAML rieng. Hien tai moi la metadata va boundary. |

## Boundary De Tach DB Sau Nay

Tat ca repo tenant-scoped nen di qua:

```go
db, err := data.DBForTenant(ctx, tenantID)
```

hoac:

```go
err := data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
    // query tenant data
    return nil
})
```

Hien tai `DBForTenant` tra ve shared DB. Sau nay khi `deployment_mode = DEDICATED`, resolver se:

1. Doc `tenant_runtime_configs` theo `tenant_id`.
2. Doc `tenant_datastores` theo `tenant_id`.
3. Lay DSN qua `dsn_secret_ref`.
4. Tao/cache pool DB rieng.
5. Chay migration/provisioning rieng neu can.
6. Tra pool dedicated cho repo.

Repo nghiep vu khong duoc hard-code shared DB neu query du lieu theo tenant.

## Boundary De Tach Auth Sau Nay

Bang `tenant_identity_providers` luu metadata IdP rieng:

```text
tenant_id
provider          -- ZITADEL | OIDC | SAML | AZURE_AD
issuer
client_id
client_secret_ref
status
```

Hien tai frontend/backend dung shared auth. Khi co enterprise tenant, auth resolver se dua tren tenant/host/login_hint de chon IdP tuong ung.

## Quy Tac Cho Bang Nghiep Vu

Moi bang du lieu thuoc tenant phai co:

```sql
tenant_id UUID NOT NULL REFERENCES tenants(id)
```

Index phai bat dau bang `tenant_id` cho cac query pho bien:

```sql
CREATE INDEX idx_orders_tenant_created
ON orders (tenant_id, created_at DESC);

CREATE UNIQUE INDEX idx_customers_tenant_code
ON customers (tenant_id, code)
WHERE deleted_at IS NULL;
```

Query phai luon co tenant scope:

```sql
SELECT *
FROM orders
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2;
```

Nen bat PostgreSQL RLS cho bang tenant-owned:

```sql
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_orders ON orders
USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);
```

## Khi Nao Dung Dedicated Tenant

Khong nen cho user tu tao dedicated tenant tu man hinh onboarding mac dinh. Dedicated tenant nen la provisioning flow rieng cho platform admin/enterprise onboarding, vi can:

- tao DB/schema rieng;
- quan ly secret DSN;
- chay migration;
- thiet lap backup/monitoring;
- cau hinh IdP rieng;
- kiem tra routing/auth/gateway.

UI mac dinh nen tao:

```text
deployment_mode = SHARED
auth_mode = SHARED_AUTH
```

Platform console sau nay moi cho phep chon:

```text
deployment_mode = DEDICATED
auth_mode = DEDICATED_AUTH
```
