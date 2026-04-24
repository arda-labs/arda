---
name: postgresql
description: Hỗ trợ thiết kế và tối ưu hóa PostgreSQL database
disable-model-invocation: false

---
# PostgreSQL Skill

Mục đích: Hỗ trợ thiết kế, tối ưu hóa và quản lý PostgreSQL database cho dự án Arda.

## 🎯 Phạm vi

- Viết tối ưu SQL queries
- Thiết kế database schema
- Tạo database indexes
- Thiết kế database views
- Tạo database migrations
- Optimize query performance
- Setup PostgreSQL configuration

## 📦 Database Schema Patterns

### Multi-Tenancy Schema

```sql
-- Multi-tenant tables with tenant_id column
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    CONSTRAINT uk_customers_tenant_code UNIQUE (tenant_id, code)
);

-- Row-Level Security for multi-tenancy
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

CREATE POLICY customers_isolation_policy ON customers
    USING (tenant_id = current_setting('app.current_tenant_id', true)::VARCHAR);

CREATE POLICY customers_insert_policy ON customers
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::VARCHAR);
```

### Audit Trail Pattern

```sql
-- Create audit trigger function
CREATE OR REPLACE FUNCTION audit_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO audit_logs (
        table_name,
        operation,
        old_data,
        new_data,
        changed_by,
        changed_at,
        tenant_id
    ) VALUES (
        TG_TABLE_NAME,
        TG_OP,
        CASE WHEN TG_OP IN ('UPDATE', 'DELETE') THEN to_jsonb(OLD) ELSE NULL END,
        CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN to_jsonb(NEW) ELSE NULL END,
        current_setting('app.current_user_id', true)::VARCHAR,
        NOW(),
        current_setting('app.current_tenant_id', true)::VARCHAR
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name VARCHAR(255) NOT NULL,
    operation VARCHAR(10) NOT NULL,
    old_data JSONB,
    new_data JSONB,
    changed_by VARCHAR(255),
    changed_at TIMESTAMPTZ NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index on audit_logs
CREATE INDEX idx_audit_logs_table_name ON audit_logs(table_name);
CREATE INDEX idx_audit_logs_changed_at ON audit_logs(changed_at);
CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);

-- Apply audit trigger to tables
CREATE TRIGGER customers_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON customers
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();
```

### Soft Delete Pattern

```sql
-- Add deleted_at column for soft delete
ALTER TABLE customers ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Create index on deleted_at
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);

-- Update queries to exclude soft-deleted records
SELECT * FROM customers WHERE deleted_at IS NULL AND tenant_id = 'tenant-123';
```

## 🎨 Indexing Patterns

### Composite Indexes

```sql
-- Composite index for multi-column queries
CREATE INDEX idx_customers_tenant_status_code ON customers(tenant_id, status, code);
```

### Covering Indexes

```sql
-- Covering index to avoid table access
CREATE INDEX idx_customers_covering ON customers(tenant_id, status)
    INCLUDE (name, code, created_at);
```

### Partial Indexes

```sql
-- Partial index for frequently filtered values
CREATE INDEX idx_customers_active ON customers(tenant_id, code)
    WHERE status = 'ACTIVE';
```

### Functional Indexes

```sql
-- Functional index for case-insensitive search
CREATE INDEX idx_customers_name_lower ON customers(LOWER(name));
```

## 📦 View Patterns

### Materialized Views

```sql
-- Create materialized view for reporting
CREATE MATERIALIZED VIEW mv_customer_summary AS
SELECT
    tenant_id,
    status,
    COUNT(*) AS customer_count,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') AS new_customers
FROM customers
WHERE deleted_at IS NULL
GROUP BY tenant_id, status;

-- Create index on materialized view
CREATE INDEX idx_mv_customer_summary_tenant ON mv_customer_summary(tenant_id);

-- Refresh materialized view
REFRESH MATERIALIZED VIEW mv_customer_summary;
```

### Common Table Expressions (CTE)

```sql
WITH customer_transactions AS (
    SELECT
        c.id AS customer_id,
        c.code AS customer_code,
        c.name AS customer_name,
        t.id AS transaction_id,
        t.amount,
        t.transaction_date
    FROM customers c
    LEFT JOIN transactions t ON c.id = t.customer_id
    WHERE c.tenant_id = 'tenant-123'
        AND c.deleted_at IS NULL
)
SELECT * FROM customer_transactions ORDER BY transaction_date DESC LIMIT 100;
```

## 📦 Migration Patterns

### Versioned Migrations

```sql
-- migrations/V001__create_customers_table.sql

-- Migration: V001__create_customers_table
-- Description: Create customers table with RLS
-- Date: 2026-04-25

BEGIN;

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uk_customers_tenant_code UNIQUE (tenant_id, code)
);

CREATE INDEX idx_customers_tenant_id ON customers(tenant_id);
CREATE INDEX idx_customers_status ON customers(status);
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);

ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

CREATE POLICY customers_isolation_policy ON customers
    USING (tenant_id = current_setting('app.current_tenant_id', true)::VARCHAR);

CREATE POLICY customers_insert_policy ON customers
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::VARCHAR);

COMMIT;
```

### Migration Rollback

```sql
-- migrations/rollback/R001__drop_customers_table.sql

BEGIN;

DROP POLICY IF EXISTS customers_insert_policy ON customers;
DROP POLICY IF EXISTS customers_isolation_policy ON customers;
DROP TABLE IF EXISTS customers;

COMMIT;
```

## 🎯 Performance Tips

### Query Optimization

```sql
-- Use EXPLAIN ANALYZE to analyze query performance
EXPLAIN ANALYZE
SELECT * FROM customers
WHERE tenant_id = 'tenant-123' AND status = 'ACTIVE'
ORDER BY created_at DESC LIMIT 100;
```

### Batch Operations

```sql
-- Use COPY for bulk inserts
COPY customers (tenant_id, code, name, status, created_by)
FROM '/tmp/customers.csv' DELIMITER ',' CSV HEADER;

-- Use batching for updates
UPDATE customers
SET status = 'INACTIVE'
WHERE tenant_id = 'tenant-123' AND id = ANY($1);
```

### Connection Pooling

```sql
-- Use connection pooling for better performance
-- Configure in application (pgx)
pool.MaxConns = 25
pool.MinConns = 5
pool.MaxConnLifetime = 1 * time.Hour
pool.MaxConnIdleTime = 30 * time.Minute
pool.HealthCheckPeriod = 1 * time.Minute
```

## 📦 Configuration

### postgresql.conf Key Settings

```ini
# Memory Settings
shared_buffers = 4GB
effective_cache_size = 12GB
maintenance_work_mem = 1GB
work_mem = 256MB

# WAL Settings
wal_buffers = 16MB
max_wal_size = 4GB
min_wal_size = 1GB

# Query Tuning
random_page_cost = 1.1
effective_io_concurrency = 200

# Connection Settings
max_connections = 200
```

## 🎯 Usage Examples

```
/postgresql "Tạo schema cho module customers"

Usage:
/postgresql "Tạo schema cho module customers với các table: customers, customer_addresses, customer_contacts"

Sẽ:
1. Tạo SQL DDL cho các tables
2. Setup Row-Level Security
3. Tạo indexes
4. Tạo audit triggers
```

```
/postgresql "Optimize query"

Usage:
/postgresql "Optimize query: SELECT * FROM customers WHERE LOWER(name) LIKE '%john%' AND tenant_id = 'tenant-123'"

Sẽ:
1. Phân tích query hiện tại
2. Đề xuất indexes
3. Cải thiện query
```

```
/postgresql "Tạo migration"

Usage:
/postgresql "Tạo migration thêm column phone vào customers table"

Sẽ:
1. Tạo migration SQL file
2. Tạo rollback script
3. Validate migration
```

## 📦 Best Practices

### Naming Conventions

- Tables: snake_case, plural (customers, transactions)
- Columns: snake_case, lowercase (customer_id, created_at)
- Indexes: idx_[table]_[columns] (idx_customers_tenant_id)
- Constraints: uk_[table]_[columns] (uk_customers_tenant_code)
- Views: v_[description] or mv_[description] (mv_customer_summary)

### Data Types

- Use UUID for primary keys
- Use TIMESTAMPTZ for timestamps
- Use TEXT for variable-length strings
- Use NUMERIC for monetary values
- Use JSONB for flexible data
- Use VARCHAR with length for codes/identifiers

### Security

- Enable Row-Level Security for multi-tenancy
- Use prepared statements to prevent SQL injection
- Implement audit logging
- Use least privilege principle for database users

---

*Last Updated: 2026-04-25*
