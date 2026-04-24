---
name: epas-mapping
description: Hỗ trợ mapping EPAS features sang Arda architecture
disable-model-invocation: false

---
# EPAS to Arda Mapping Skill

Mục đích: Hỗ trợ mapping EPAS features sang Arda architecture cho migration process.

## 🎯 Phạm vi

- Mapping EPAS table sang Arda schema
- Mapping EPAS API sang Arda gRPC API
- Mapping EPAS component sang Arda component
- Generate migration plan
- Tạo data mapping document

## 📦 EPAS Overview

### EPAS Architecture (Source)

```
EPAS Architecture:
┌─────────────────────────────────────────┐
│         Frontend (Legacy)              │
│  - fe_host (Angular 15)               │
│  - fe_common (Shared)                 │
│  - fe_accounting (Accounting)          │
│  - fe_loan (Loan)                      │
│  - fe_crm (CRM)                       │
│  - fe_fac (Facility)                   │
│  - fe_bpm (BPM)                       │
├─────────────────────────────────────────┤
│         Backend (Monolithic)             │
│  - epas-be (Spring Boot)              │
│  - epas-api (REST API)                │
├─────────────────────────────────────────┤
│         Database (PostgreSQL)           │
│  - epas_db (Legacy schema)             │
└─────────────────────────────────────────┘
```

## 📦 Arda Architecture (Target)

```
Arda Architecture:
┌─────────────────────────────────────────┐
│         Frontend (Modern)               │
│  - shell (Host - MFE)                 │
│  - common (Shared UI)                  │
│  - accounting (MFE)                   │
│  - loan (MFE)                         │
│  - crm (MFE)                          │
│  - hrm (MFE)                           │
│  - admin (MFE)                         │
├─────────────────────────────────────────┤
│         Backend Go (Operational)         │
│  - iam-service (Kratos + gRPC)        │
│  - crm-service (Kratos + gRPC)        │
│  - hrm-service (Kratos + gRPC)        │
│  - notification-service (Kratos)        │
│  - system-config-service (Kratos)        │
├─────────────────────────────────────────┤
│         Backend Java (Core Banking)      │
│  - accounting-service (Spring + R2DBC)  │
│  - loan-service (Spring + R2DBC)      │
│  - deposit-service (Spring + R2DBC)    │
│  - treasury-service (Spring + R2DBC)   │
├─────────────────────────────────────────┤
│         Infrastructure                  │
│  - K3s (Kubernetes)                  │
│  - ArgoCD (GitOps)                    │
│  - APISIX (API Gateway)               │
│  - Zitadel (Identity)                 │
│  - Redpanda (Events)                  │
│  - Camunda (BPM)                      │
└─────────────────────────────────────────┘
```

## 📦 Table Mapping

### Accounting Tables

```markdown
## EPAS to Arda Table Mapping

### Journals

| EPAS Table | Arda Table | Notes |
|------------|-------------|--------|
| `ac_journal` | `journals` | Simplified structure |
| `ac_journal_no` | `journals.journal_no` | Renamed from ac_journal_no |
| `ac_journal_date` | `journals.journal_date` | TIMESTAMPTZ instead of DATE |
| `ac_description` | `journals.description` | Renamed from ac_description |
| `ac_status` | `journals.status` | ENUM: DRAFT, POSTED, CANCELLED |
| `ac_created_by` | `journals.created_by` | Added audit column |
| `ac_created_at` | `journals.created_at` | TIMESTAMPTZ |
| `ac_updated_by` | `journals.updated_by` | Added audit column |
| `ac_updated_at` | `journals.updated_at` | TIMESTAMPTZ |
| `ac_deleted_at` | `journals.deleted_at` | Added for soft delete |
| - | `journals.tenant_id` | Added for multi-tenancy |

**Column Mapping:**
```sql
-- EPAS
ac_journal_no VARCHAR(50) NOT NULL
ac_journal_date DATE NOT NULL
ac_description TEXT
ac_status VARCHAR(20) DEFAULT 'DRAFT'

-- Arda
journals.journal_no VARCHAR(50) NOT NULL
journals.journal_date TIMESTAMPTZ NOT NULL
journals.description TEXT NOT NULL
journals.status VARCHAR(20) DEFAULT 'DRAFT'
journals.tenant_id VARCHAR(255) NOT NULL
```

### Journal Items

| EPAS Table | Arda Table | Notes |
|------------|-------------|--------|
| `ac_journal_item` | `journal_items` | Simplified structure |
| `ac_item_id` | `journal_items.id` | UUID instead of INTEGER |
| `ac_journal_id` | `journal_items.journal_id` | UUID instead of INTEGER |
| `ac_line_no` | `journal_items.line_no` | Same |
| `ac_account_code` | `journal_items.account_code` | Same |
| `ac_account_name` | `journal_items.account_name` | Same |
| `ac_debit_amount` | `journal_items.debit_amount` | NUMERIC(20,6) instead of DECIMAL |
| `ac_credit_amount` | `journal_items.credit_amount` | NUMERIC(20,6) instead of DECIMAL |
| `ac_cost_center` | `journal_items.cost_center_id` | Renamed |
| - | `journal_items.tenant_id` | Added for multi-tenancy |
```

### CRM Tables

```markdown
## CRM Table Mapping

### Customers

| EPAS Table | Arda Table | Notes |
|------------|-------------|--------|
| `m_customer` | `customers` | Renamed |
| `m_cust_id` | `customers.id` | UUID instead of INTEGER |
| `m_cust_no` | `customers.code` | Renamed |
| `m_cust_name` | `customers.name` | Renamed |
| `m_cust_email` | `customers.email` | Same |
| `m_cust_phone` | `customers.phone` | Same |
| `m_cust_address` | `customers.address` | TEXT instead of VARCHAR(255) |
| `m_cust_status` | `customers.status` | ENUM: ACTIVE, INACTIVE, SUSPENDED |
| - | `customers.tenant_id` | Added for multi-tenancy |
| - | `customers.deleted_at` | Added for soft delete |
```

### Users

| EPAS Table | Arda Table | Notes |
|------------|-------------|--------|
| `sys_user` | `users` | Migrated to IAM service |
| `sys_user_id` | `users.id` | UUID |
| `sys_user_name` | `users.full_name` | Renamed |
| `sys_user_email` | `users.email` | Same |
| `sys_user_password` | `users.password` | Hashed |
| `sys_user_status` | `users.status` | Same |
| - | `users.tenant_id` | Added for multi-tenancy |
```

## 📦 API Mapping

### REST to gRPC

```markdown
## EPAS REST API to Arda gRPC Mapping

### Accounting APIs

| EPAS REST Endpoint | Arda gRPC Service | Arda gRPC Method | Notes |
|--------------------|--------------------|-----------------|--------|
| `GET /api/journals` | `AccountingService` | `ListJournals` | Added pagination |
| `GET /api/journals/{id}` | `AccountingService` | `GetJournal` | Same |
| `POST /api/journals` | `AccountingService` | `CreateJournal` | Same |
| `PUT /api/journals/{id}` | `AccountingService` | `UpdateJournal` | Added |
| `DELETE /api/journals/{id}` | `AccountingService` | `DeleteJournal` | Soft delete |
| `POST /api/journals/{id}/post` | `AccountingService` | `PostJournal` | New method |
| `POST /api/journals/{id}/cancel` | `AccountingService` | `CancelJournal` | New method |

### CRM APIs

| EPAS REST Endpoint | Arda gRPC Service | Arda gRPC Method | Notes |
|--------------------|--------------------|-----------------|--------|
| `GET /api/customers` | `CRMService` | `ListCustomers` | Same |
| `GET /api/customers/{id}` | `CRMService` | `GetCustomer` | Same |
| `POST /api/customers` | `CRMService` | `CreateCustomer` | Same |
| `PUT /api/customers/{id}` | `CRMService` | `UpdateCustomer` | Same |
| `DELETE /api/customers/{id}` | `CRMService` | `DeleteCustomer` | Soft delete |
| `GET /api/customers/{id}/contacts` | `CRMService` | `ListCustomerContacts` | Same |
| `POST /api/customers/{id}/contacts` | `CRMService` | `CreateCustomerContact` | Same |

### Authentication APIs

| EPAS REST Endpoint | Arda gRPC Service | Notes |
|--------------------|--------------------|--------|
| `POST /api/auth/login` | `IAMService` (Zitadel OIDC) | Use OIDC instead |
| `POST /api/auth/logout` | `IAMService` (Zitadel OIDC) | Use OIDC instead |
| `GET /api/auth/user` | `IAMService` | `GetUser` |
| `POST /api/auth/refresh` | `IAMService` | `RefreshToken` |
```

### Request/Response Mapping

```markdown
## Request/Response Mapping

### Create Journal

**EPAS REST Request:**
```json
{
  "journalNo": "JNL-2024-001",
  "journalDate": "2024-04-25",
  "description": "Monthly closing journal"
}
```

**Arda gRPC Request:**
```protobuf
message CreateJournalRequest {
  string tenant_id = 1;
  string journal_no = 2;
  int64 journal_date = 3;  // Unix timestamp
  string description = 4;
}
```

**Arda gRPC Response:**
```protobuf
message CreateJournalResponse {
  Journal journal = 1;
}

message Journal {
  string id = 1;
  string tenant_id = 2;
  string journal_no = 3;
  int64 journal_date = 4;
  string description = 5;
  string status = 6;
}
```

### List Journals

**EPAS REST Response:**
```json
{
  "data": [
    {
      "id": 1,
      "journalNo": "JNL-2024-001",
      "description": "Monthly closing journal"
    }
  ],
  "total": 100
}
```

**Arda gRPC Request:**
```protobuf
message ListJournalsRequest {
  string tenant_id = 1;
  int32 page = 2;
  int32 page_size = 3;
  string status = 4;
}
```

**Arda gRPC Response:**
```protobuf
message ListJournalsResponse {
  repeated Journal journals = 1;
  int32 total_count = 2;
  int32 page = 3;
  int32 page_size = 4;
  int32 total_pages = 5;
}
```
```

## 📦 Component Mapping

```markdown
## Component Mapping

### Accounting Components

| EPAS Component | Arda MFE | Notes |
|---------------|-----------|--------|
| `JournalListComponent` | `accounting/journals/list` | Migrated to PrimeNG table |
| `JournalFormComponent` | `accounting/journals/form` | Migrated to PrimeNG form |
| `JournalDetailComponent` | `accounting/journals/detail` | Same functionality |
| `JournalBalanceComponent` | `accounting/balances/list` | New component |
| `JournalReportComponent` | `accounting/reports/list` | New component |

### CRM Components

| EPAS Component | Arda MFE | Notes |
|---------------|-----------|--------|
| `CustomerListComponent` | `crm/customers/list` | Migrated to PrimeNG table |
| `CustomerFormComponent` | `crm/customers/form` | Migrated to PrimeNG form |
| `CustomerDetailComponent` | `crm/customers/detail` | Same functionality |
| `ContactListComponent` | `crm/contacts/list` | New component |
| `LeadListComponent` | `crm/leads/list` | New component |

### Common Components

| EPAS Component | Arda MFE | Notes |
|---------------|-----------|--------|
| `SidebarComponent` | `common/sidebar` | Migrated to shell app |
| `HeaderComponent` | `common/header` | Migrated to shell app |
| `NotificationComponent` | `common/notifications` | Enhanced with toast |
| `AuthGuard` | `auth/guard` | Updated for OIDC |
| `ErrorBoundaryComponent` | `common/error-boundary` | Enhanced |
```

## 📦 Migration Plan

```markdown
## Migration Plan

### Phase 1: Infrastructure Setup (2 weeks)

**Tasks:**
1. Setup K3s cluster on thinkcenter
2. Configure APISIX gateway
3. Setup Zitadel identity provider
4. Setup Redpanda event broker
5. Setup Camunda BPM engine
6. Configure ArgoCD for GitOps

**Deliverables:**
- Infrastructure documentation
- Deployment scripts
- GitOps repository structure

### Phase 2: Data Migration (3 weeks)

**Tasks:**
1. Create Arda database schema
2. Develop data migration scripts
3. Migrate accounting data
4. Migrate CRM data
5. Migrate user accounts to Zitadel
6. Validate data integrity

**Deliverables:**
- Migration scripts
- Data validation reports
- Rollback scripts

### Phase 3: Backend Services (4 weeks)

**Tasks:**
1. Develop IAM service (Kratos + gRPC)
2. Develop CRM service (Kratos + gRPC)
3. Develop HRM service (Kratos + gRPC)
4. Develop notification service (Kratos + gRPC)
5. Develop system-config service (Kratos + gRPC)
6. Develop accounting service (Spring + R2DBC)

**Deliverables:**
- Service implementations
- gRPC proto definitions
- Unit tests

### Phase 4: Frontend MFEs (4 weeks)

**Tasks:**
1. Migrate shell app (host)
2. Migrate accounting MFE
3. Migrate CRM MFE
4. Migrate HRM MFE
5. Migrate admin MFE
6. Update common UI library

**Deliverables:**
- MFE applications
- Module Federation configuration
- Unit tests

### Phase 5: Testing & Validation (2 weeks)

**Tasks:**
1. Write integration tests
2. Write E2E tests
3. Performance testing
4. Security testing
5. User acceptance testing

**Deliverables:**
- Test suite
- Test reports
- Bug fixes

### Phase 6: Cutover (1 week)

**Tasks:**
1. Final data migration
2. Switch DNS to new system
3. Monitor production
4. Handle issues
5. Decommission EPAS

**Deliverables:**
- Production deployment
- Monitoring setup
- Decommission plan
```

## 📦 Data Migration Scripts

### Accounting Data Migration

```sql
-- EPAS to Arda: Journals Migration Script

-- Step 1: Migrate journals
INSERT INTO arda.journals (
    id,
    tenant_id,
    journal_no,
    journal_date,
    description,
    status,
    created_at,
    updated_at,
    created_by,
    updated_by,
    deleted_at
)
SELECT
    gen_random_uuid(),                                          -- Generate UUID
    'tenant-1',                                                  -- Set default tenant (will be updated later)
    ac_journal_no,
    (ac_journal_date || 'T00:00:00Z')::TIMESTAMPTZ,         -- Convert to TIMESTAMPTZ
    ac_description,
    ac_status,
    ac_created_at,
    ac_updated_at,
    ac_created_by,
    ac_updated_by,
    NULL                                                         -- No soft delete in EPAS
FROM epas.ac_journal
WHERE ac_deleted_at IS NULL                                    -- Only active records
ORDER BY ac_created_at;

-- Step 2: Migrate journal items
INSERT INTO arda.journal_items (
    id,
    tenant_id,
    journal_id,
    line_no,
    account_code,
    account_name,
    debit_amount,
    credit_amount,
    cost_center_id,
    department_id,
    project_id,
    created_at,
    updated_at,
    created_by
)
SELECT
    gen_random_uuid(),
    j.tenant_id,
    j.id,
    ji.ac_line_no,
    ji.ac_account_code,
    ji.ac_account_name,
    ji.ac_debit_amount,
    ji.ac_credit_amount,
    ji.ac_cost_center,
    ji.ac_department,
    ji.ac_project,
    ji.ac_created_at,
    ji.ac_updated_at,
    ji.ac_created_by
FROM epas.ac_journal_item ji
INNER JOIN arda.journals j ON j.journal_no = ji.ac_journal_no
WHERE ji.ac_deleted_at IS NULL;

-- Step 3: Update tenant_id based on branch/organization
UPDATE arda.journals
SET tenant_id = 
    CASE 
        WHEN ac_journal_no LIKE 'BR1-%' THEN 'tenant-branch-1'
        WHEN ac_journal_no LIKE 'BR2-%' THEN 'tenant-branch-2'
        ELSE 'tenant-main'
    END
FROM epas.ac_journal e
WHERE arda.journals.journal_no = e.ac_journal_no;
```

### CRM Data Migration

```sql
-- EPAS to Arda: Customers Migration Script

-- Step 1: Migrate customers
INSERT INTO arda.customers (
    id,
    tenant_id,
    code,
    name,
    email,
    phone,
    address,
    status,
    created_at,
    updated_at,
    deleted_at
)
SELECT
    gen_random_uuid(),
    'tenant-1',                                      -- Set default tenant
    m_cust_no,
    m_cust_name,
    m_cust_email,
    m_cust_phone,
    m_cust_address,
    m_cust_status,
    m_created_at,
    m_updated_at,
    NULL
FROM epas.m_customer
WHERE m_deleted_at IS NULL
ORDER BY m_created_at;

-- Step 2: Migrate customer contacts
INSERT INTO arda.customer_contacts (
    id,
    tenant_id,
    customer_id,
    name,
    email,
    phone,
    position,
    created_at
)
SELECT
    gen_random_uuid(),
    c.tenant_id,
    c.id,
    mc.mc_name,
    mc.mc_email,
    mc.mc_phone,
    mc.mc_position,
    mc.mc_created_at
FROM epas.m_customer mc
INNER JOIN arda.customers c ON c.code = mc.mc_cust_no
WHERE mc.mc_deleted_at IS NULL;
```

## 📦 Migration Validation

### Data Validation Queries

```sql
-- Validate journal counts
SELECT 
    'Total EPAS journals' AS source,
    COUNT(*) AS count
FROM epas.ac_journal
WHERE ac_deleted_at IS NULL

UNION ALL

SELECT 
    'Total Arda journals' AS source,
    COUNT(*) AS count
FROM arda.journals
WHERE deleted_at IS NULL;

-- Validate journal item counts
SELECT 
    'Total EPAS journal items' AS source,
    COUNT(*) AS count
FROM epas.ac_journal_item
WHERE ac_deleted_at IS NULL

UNION ALL

SELECT 
    'Total Arda journal items' AS source,
    COUNT(*) AS count
FROM arda.journal_items
WHERE deleted_at IS NULL;

-- Validate customer counts
SELECT 
    'Total EPAS customers' AS source,
    COUNT(*) AS count
FROM epas.m_customer
WHERE m_deleted_at IS NULL

UNION ALL

SELECT 
    'Total Arda customers' AS source,
    COUNT(*) AS count
FROM arda.customers
WHERE deleted_at IS NULL;

-- Validate balance totals
SELECT 
    'EPAS debit total' AS metric,
    SUM(ac_debit_amount) AS total
FROM epas.ac_journal_item
WHERE ac_deleted_at IS NULL

UNION ALL

SELECT 
    'Arda debit total' AS metric,
    SUM(debit_amount) AS total
FROM arda.journal_items
WHERE deleted_at IS NULL;
```

## 📦 Rollback Scripts

```sql
-- Arda to EPAS Rollback Script

-- Rollback journals
DELETE FROM arda.journals;

-- Rollback journal items
DELETE FROM arda.journal_items;

-- Rollback customers
DELETE FROM arda.customers;

-- Rollback customer contacts
DELETE FROM arda.customer_contacts;

-- Reset sequences (if any)
-- Note: Arda uses UUIDs, no sequences to reset
```

## 🎯 Usage Examples

```
/epas-mapping "Tạo table mapping"

Usage:
/epas-mapping "Tạo mapping cho accounting tables: journals, journal_items, balances"

Sẽ:
1. Analyze EPAS table structure
2. Design Arda table structure
3. Create column mapping
```

```
/epas-mapping "Tạo API mapping"

Usage:
/epas-mapping "Map REST APIs sang gRPC cho accounting module"

Sẽ:
1. Analyze EPAS REST endpoints
2. Design gRPC services
3. Create method mapping
```

```
/epas-mapping "Tạo migration plan"

Usage:
/epas-mapping "Tạo migration plan cho accounting module migration"

Sẽ:
1. Define migration phases
2. Create task list
3. Set timeline
```

---
*Last Updated: 2026-04-25*
