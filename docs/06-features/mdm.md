# MDM — Dữ liệu nền

Updated: 2026-04-30

Status: Implemented first release. Backend, migrations, IAM menu seed, frontend
MFE, and GitOps manifests exist.

## Scope

`mdm-service` owns master data that is reused across IAM, CRM, HRM, accounting, banking, payments, and reporting. It replaces the old `common-service` skeleton because `common` is too broad and tends to become an unclear catch-all.

Runtime contract:

- Backend service: `apps/backend-go/mdm-service`
- Database: `mdm`
- Database user: `mdm`
- API prefix: `/api/v1/mdm/*`
- Frontend MFE: `apps/frontend-micro/projects/mdm`
- MFE asset route: `/mfe-mdm/*`
- Shell route: `/app/mdm/*`
- Permissions: `mdm:read`, `mdm:create`, `mdm:update`, `mdm:delete`

## Current CRUD

The first MDM release covers these entities:

| Area | Entity | Purpose |
| --- | --- | --- |
| Geographic | Administrative units | Province/city and ward/commune/special-zone hierarchy |
| Geographic | Area types | Type definitions for business regions, delivery zones, branch regions, risk zones |
| Geographic | Areas | Business-defined regions that can map to administrative units |
| Reference data | Code sets | Named catalogs such as currency, payment channel, risk rating |
| Reference data | Code items | Values inside a code set, including ordering, flags, metadata, and status |
| System | System parameters | Runtime parameters grouped by security, UI, banking, finance, risk, pricing |

## Frontend Menu

IAM seeds the MDM menu as a parent with three groups:

| Parent | Child routes |
| --- | --- |
| Địa lý hành chính | `/app/mdm/geo/administrative-units`, `/app/mdm/geo/area-types`, `/app/mdm/geo/areas` |
| Danh mục & tham số | `/app/mdm/catalog/code-sets`, `/app/mdm/catalog/code-items`, `/app/mdm/system/parameters` |
| Tài chính ngân hàng | `/app/mdm/banking/reference` |

The banking page is currently a proposal/reference view. CRUD should only be added once the owning domain is clear.

## Banking Reference Data Candidates

Good MDM candidates for a financial/banking platform:

- `CURRENCY`, `COUNTRY`, `BUSINESS_CALENDAR`
- `BANK`, `BANK_BRANCH`, `PAYMENT_NETWORK`
- `BANK_ACCOUNT_TYPE`, `PAYMENT_METHOD`, `PAYMENT_CHANNEL`
- `TRANSACTION_TYPE`, `FEE_TYPE`, `TAX_CODE`
- `CUSTOMER_SEGMENT`, `DOCUMENT_TYPE`, `OCCUPATION`, `ECONOMIC_SECTOR`
- `RISK_RATING`, `AML_LIST_SOURCE`, `SANCTION_SOURCE`
- `INTEREST_RATE_TYPE`, `COLLATERAL_TYPE`, `LIMIT_PROFILE`

Do not put product-owned rules into MDM just because they are lists. Product pricing, loan policy, accounting rules, and approval workflow definitions should stay in their domain services unless they are truly cross-domain reference data.

## Local And Dev Database

The standalone service should use its own database and login:

```text
postgres://mdm:mdm%40123@thinkcenter:5432/mdm?sslmode=disable
```

Using the IAM database user for MDM is only acceptable for short-lived bootstrap work. The committed service config and GitOps secret script use the dedicated `mdm` user.
