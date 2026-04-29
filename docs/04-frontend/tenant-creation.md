# Frontend Tenant Creation Guide

> Cap nhat: 2026-04-29
> Pham vi: Angular MFE shell/IAM, tao workspace/tenant moi tu FE.

## Contract

Tenant creation la thao tac tao workspace moi cho current authenticated user. FE goi IAM API:

```http
POST /api/v1/tenants
Content-Type: application/json
Authorization: Bearer <access-token>
```

Body mac dinh cho UI onboarding:

```json
{
  "name": "Acme Bank",
  "slug": "acme-bank",
  "deployment_mode": "SHARED",
  "auth_mode": "SHARED_AUTH"
}
```

Khong can `X-Tenant-ID` khi tao tenant dau tien, vi tenant moi chua ton tai. Backend lay current user tu auth token, tao tenant, owner role, menu/permission mac dinh, roi tao `tenant_users` cho owner.

Entry points hien tai:

- `/workspaces`: man hinh quan ly va tao workspace trong shell layout.
- Header tenant switcher: nut `+` mo nhanh `/workspaces`.
- `/select-workspace`: man hinh chon workspace sau login, van co action tao workspace khi can.

Response:

```json
{
  "id": "10000000-0000-0000-0000-000000000001",
  "name": "Acme Bank",
  "slug": "acme-bank",
  "owner_id": "...",
  "deployment_mode": "SHARED",
  "auth_mode": "SHARED_AUTH"
}
```

## TenantService Pattern

`TenantService.createTenant` phai:

1. Gui `deployment_mode/auth_mode`.
2. Map response snake_case sang camelCase.
3. Chon tenant vua tao.
4. Reload `/api/v1/me/tenants`.

Pattern hien tai:

```ts
createTenant(
  name: string,
  slug: string,
  options: { deploymentMode?: TenantDeploymentMode; authMode?: TenantAuthMode } = {},
): Observable<Tenant> {
  return this.http.post<any>('/api/v1/tenants', {
    name,
    slug,
    deployment_mode: options.deploymentMode ?? 'SHARED',
    auth_mode: options.authMode ?? 'SHARED_AUTH',
  }).pipe(
    map(resp => this.toTenant(resp)),
    tap(created => {
      this.selectTenant(created.id);
      this.loadTenants().subscribe();
    }),
  );
}
```

`toTenant` phai chuan hoa field:

```ts
private toTenant(t: any): Tenant {
  return {
    id: t.id,
    name: t.name,
    slug: t.slug,
    role: t.role ?? 'owner',
    deploymentMode: t.deployment_mode ?? t.deploymentMode,
    authMode: t.auth_mode ?? t.authMode,
  };
}
```

## Workspace UI Flow

Man hinh `workspaces` la noi chinh de quan ly va tao workspace. Man hinh `select-workspace` dung khi user login thanh cong va can chon workspace, dac biet khi chua co tenant.

Flow:

1. User mo `/workspaces` tu tenant switcher hoac menu System.
2. User nhap workspace name/slug.
3. Goi `tenantService.createTenant(name, slug)`.
4. Service select tenant moi.
5. Component navigate ve `/home`.

Flow rut gon tren `select-workspace` van goi cung service nay:

```ts
createFirstWorkspace(): void {
  const name = prompt('Nhap ten Workspace moi:');
  if (!name) return;

  const slug = name.toLowerCase().replace(/\s+/g, '-');

  this.tenantService.createTenant(name, slug).subscribe({
    next: (tenant) => {
      this.tenantService.selectTenant(tenant.id);
      this.router.navigate(['/home']);
    },
    error: (err) => {
      console.error('Failed to create tenant', err);
      alert('Khong the tao workspace. Vui long kiem tra console.');
    },
  });
}
```

## Mode Selection Rule

UI mac dinh khong cho user tu chon dedicated mode.

Default:

```text
deploymentMode = SHARED
authMode = SHARED_AUTH
```

Dedicated mode chi nen nam trong platform/admin provisioning flow:

```text
deploymentMode = DEDICATED
authMode = DEDICATED_AUTH
```

Ly do: dedicated tenant can provisioning DB/schema, secret, migration, backup, monitoring va IdP rieng. Day khong phai thao tac form don gian cua onboarding UI.

## Sau Khi Tao Tenant

Sau khi tao thanh cong, FE nen reload tenant list tu:

```http
GET /api/v1/me/tenants
```

Response membership co them metadata:

```json
{
  "memberships": [
    {
      "tenantId": "...",
      "tenantName": "Acme Bank",
      "tenantSlug": "acme-bank",
      "role": "owner",
      "deploymentMode": "SHARED",
      "authMode": "SHARED_AUTH"
    }
  ]
}
```

`TenantProvider.getTenantId()` tra selected tenant ID cho auth interceptor de gan tenant header cho cac API tenant-scoped sau do.
