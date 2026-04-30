# MDM - Master Data Management

Updated: 2026-04-30

Status: Active. Backend service, database migrations, IAM menu seed, frontend
MFE pages, and GitOps manifests exist. The service now covers the core shared
reference data needed before building more transactional banking services.

## Runtime Contract

- Backend service: `apps/backend-go/mdm-service`
- Database: `mdm`
- Database user: `mdm`
- Native API prefix: `/v1/mdm/*`
- Gateway API prefix: `/api/v1/mdm/*`
- Frontend MFE: `apps/frontend-micro/projects/mdm`
- MFE asset route: `/mfe-mdm/*`
- Shell route: `/app/mdm/*`
- Permissions: `mdm:read`, `mdm:create`, `mdm:update`, `mdm:delete`

## Implemented Domains

| Group | Entities | UI route |
| --- | --- | --- |
| Geographic | Administrative units, area types, areas | `/app/mdm/geo/*` |
| Reference data | Code sets, code items | `/app/mdm/catalog/*` |
| System | System parameters | `/app/mdm/system/parameters` |
| Banking counterparties | Credit institutions | `/app/mdm/banking/credit-institutions` |
| Business time | Business calendars, working hours, holidays/exceptions | `/app/mdm/banking/business-calendars` |
| Pricing governance | Fee schedules, tax rules, standard limits | `/app/mdm/banking/pricing-rules` |
| Treasury reference | Currencies, FX rate sources, FX rates | `/app/mdm/banking/currency-fx` |
| Product/channel setup | Banking products, service channels, product-channel rules | `/app/mdm/banking/product-channels` |
| Payment routing | Bank branches, SWIFT/NAPAS codes, payment networks | `/app/mdm/banking/payment-networks` |

## Frontend Menu

IAM seeds MDM under a parent menu with these groups:

| Parent | Child routes |
| --- | --- |
| Địa lý hành chính | Tỉnh/phường xã, loại khu vực, khu vực quản lý |
| Danh mục & tham số | Bộ danh mục, giá trị danh mục, tham số hệ thống |
| Tài chính ngân hàng | Gợi ý mở rộng, tổ chức tín dụng, lịch làm việc, biểu phí/thuế/hạn mức, tiền tệ/tỷ giá, sản phẩm/kênh, chi nhánh/mạng thanh toán |

Important migration note: IAM `menus` uses `tenant_id`, `slug`, `route`,
`enabled`, and `permission_slug`. New MDM menu migrations should follow the
existing `000021`-`000026` style and must not use `code`, `path`, or `status`
columns.

## External Data Sync

Vietnam administrative unit data is synced from CASSO AddressKit:

```text
https://production.cas.so/address-kit/latest/provinces
https://production.cas.so/address-kit/latest/communes
```

The MDM sync endpoint replaces current administrative units in one transaction:

```powershell
curl -X POST http://localhost:8001/api/v1/mdm/administrative-units/sync-addresskit
```

The frontend exposes this as the `Đồng bộ AddressKit` action on the
administrative unit page.

## Governance Rules

- MDM owns data that is shared across services and stable enough to be reused.
- Domain-specific transaction rules should stay in their owning service unless
  they are truly cross-domain reference data.
- Pricing, tax, limit, FX, product-channel, and payment routing data should be
  effective-dated where downstream services need historical behavior.
- Operational records should use soft delete with status changes so dependent
  systems do not break on missing references.

## Verification

Backend:

```powershell
cd apps\backend-go\mdm-service
go test ./...
```

Frontend:

```powershell
cd apps\frontend-micro
npx ng build mdm --configuration development
```

IAM menu migrations:

```powershell
cd apps\backend-go\iam-service
go test ./...
```

## Next MDM Candidates

These are useful later, but should not block notification service design:

- Customer/KYC reference lists: occupation, economic sector, document subtype,
  residency, AML list source, sanction source.
- Risk/compliance reference data: risk grades, review frequency, watchlist
  categories, screening result codes.
- Accounting reference data: chart-of-account mapping keys, posting event
  types, cost/profit center catalogs.
- Localization data: supported languages, message keys, country-specific
  formatting rules.
