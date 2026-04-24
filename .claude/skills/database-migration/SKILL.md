---
name: database-migration
description: Hỗ trợ tạo và chạy database migrations
disable-model-invocation: true

---
# Database Migration Skill

Mục đích: Hỗ trợ tạo, quản lý và chạy database migrations cho dự án Arda sử dụng Flyway.

## 🎯 Phạm vi

- Tạo migration script (SQL)
- Tạo migration rollback script
- Version control migrations
- Setup migration tool (Flyway)
- Run migrations
- Validate migrations

## 📦 Setup Flyway

### build.gradle.kts (Java/Kotlin)

```kotlin
plugins {
    id("org.flywaydb.flyway") version "9.22.3"
}

dependencies {
    implementation("org.flywaydb:flyway-core")
    implementation("org.flywaydb:flyway-database-postgresql")
}

flyway {
    url = "jdbc:postgresql://localhost:5432/arda"
    user = "arda"
    password = "arda_password"
    locations = arrayOf("classpath:db/migration")
    baselineOnMigrate = true
    baselineVersion = "0"
}
```

### build.gradle.kts (Go with sqlc)

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./queries"
    schema: "./migrations"
    gen:
      go:
        package: "db"
        out: "./internal/db"
```

## 📦 Migration File Structure

```
arda-core/services/accounting/src/main/resources/db/migration/
├── V001__create_journals_table.sql
├── V002__create_journal_items_table.sql
├── V003__create_journal_balances_table.sql
├── V004__create_journal_entries_table.sql
└── ...
```

## 📦 Migration Scripts

### Create Tables Migration

```sql
-- V001__create_journals_table.sql

-- Migration: Create journals table
-- Description: Create journals table with RLS for accounting
-- Author: Arda Team
-- Date: 2026-04-25

BEGIN;

-- Create journals table
CREATE TABLE IF NOT EXISTS journals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR(255) NOT NULL,
    journal_no VARCHAR(50) NOT NULL,
    journal_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    posting_date TIMESTAMP WITH TIME ZONE,
    posted_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uk_journals_tenant_journal_no UNIQUE (tenant_id, journal_no),
    CONSTRAINT chk_journals_status CHECK (status IN ('DRAFT', 'POSTED', 'CANCELLED'))
);

-- Create indexes
CREATE INDEX idx_journals_tenant_id ON journals(tenant_id);
CREATE INDEX idx_journals_status ON journals(status);
CREATE INDEX idx_journals_journal_date ON journals(journal_date);
CREATE INDEX idx_journals_deleted_at ON journals(deleted_at);
CREATE INDEX idx_journals_tenant_status_date ON journals(tenant_id, status, journal_date);

-- Enable Row Level Security
ALTER TABLE journals ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
CREATE POLICY journals_isolation_policy ON journals
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::VARCHAR)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::VARCHAR);

COMMIT;
```

### Add Columns Migration

```sql
-- V002__add_journals_attachment_column.sql

-- Migration: Add attachment column to journals
-- Description: Add attachment_id column to journals table
-- Date: 2026-04-25

BEGIN;

-- Add attachment_id column
ALTER TABLE journals ADD COLUMN IF NOT EXISTS attachment_id UUID;

-- Create index for attachment_id
CREATE INDEX IF NOT EXISTS idx_journals_attachment_id ON journals(attachment_id);

COMMIT;
```

### Create Tables with Foreign Keys

```sql
-- V003__create_journal_items_table.sql

-- Migration: Create journal_items table
-- Description: Create journal items table with foreign key to journals
-- Date: 2026-04-25

BEGIN;

-- Create journal_items table
CREATE TABLE IF NOT EXISTS journal_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR(255) NOT NULL,
    journal_id UUID NOT NULL REFERENCES journals(id) ON DELETE CASCADE,
    line_no INTEGER NOT NULL,
    account_code VARCHAR(50) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    description TEXT,
    debit_amount NUMERIC(20, 6) DEFAULT 0,
    credit_amount NUMERIC(20, 6) DEFAULT 0,
    cost_center_id VARCHAR(50),
    department_id VARCHAR(50),
    project_id VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    CONSTRAINT uk_journal_items_journal_line UNIQUE (journal_id, line_no),
    CONSTRAINT chk_journal_items_amounts CHECK (
        (debit_amount = 0 AND credit_amount > 0) OR
        (credit_amount = 0 AND debit_amount > 0)
    )
);

-- Create indexes
CREATE INDEX idx_journal_items_journal_id ON journal_items(journal_id);
CREATE INDEX idx_journal_items_tenant_id ON journal_items(tenant_id);
CREATE INDEX idx_journal_items_account_code ON journal_items(account_code);
CREATE INDEX idx_journal_items_cost_center ON journal_items(cost_center_id);

-- Enable Row Level Security
ALTER TABLE journal_items ENABLE ROW LEVEL SECURITY;

CREATE POLICY journal_items_isolation_policy ON journal_items
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::VARCHAR)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::VARCHAR);

COMMIT;
```

### Seed Data Migration

```sql
-- V004__seed_journal_status_values.sql

-- Migration: Seed journal status values
-- Description: Insert default journal status values for reference
-- Date: 2026-04-25

BEGIN;

-- Create reference table for journal status
CREATE TABLE IF NOT EXISTS journal_status (
    code VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Insert default status values
INSERT INTO journal_status (code, name, description) VALUES
    ('DRAFT', 'Draft', 'Journal entry in draft status'),
    ('POSTED', 'Posted', 'Journal entry posted to ledger'),
    ('CANCELLED', 'Cancelled', 'Journal entry cancelled')
ON CONFLICT (code) DO NOTHING;

COMMIT;
```

### Rollback Migration

```sql
-- R003__drop_journal_items_table.sql

-- Rollback: Drop journal_items table
-- Description: Rollback V003__create_journal_items_table.sql
-- Date: 2026-04-25

BEGIN;

-- Drop RLS policies
DROP POLICY IF EXISTS journal_items_isolation_policy ON journal_items;

-- Drop indexes
DROP INDEX IF EXISTS idx_journal_items_cost_center;
DROP INDEX IF EXISTS idx_journal_items_account_code;
DROP INDEX IF EXISTS idx_journal_items_tenant_id;
DROP INDEX IF EXISTS idx_journal_items_journal_id;

-- Drop table
DROP TABLE IF EXISTS journal_items;

COMMIT;
```

## 📦 Migration Commands

### Run Migrations (Flyway)

```bash
# Run migrations
./gradlew flywayMigrate

# Run migrations with specific config
./gradlew flywayMigrate -Dflyway.configFiles=flyway.conf

# Validate migrations
./gradlew flywayValidate

# Show migration history
./gradlew flywayInfo

# Clean database (DANGEROUS - removes all data)
./gradlew flywayClean

# Repair migration history (if corrupted)
./gradlew flywayRepair
```

### Run Migrations (Flyway CLI)

```bash
# Run migrations
flyway -url=jdbc:postgresql://localhost:5432/arda \
       -user=arda \
       -password=arda_password \
       -locations=classpath:db/migration \
       migrate

# Show migration status
flyway -url=jdbc:postgresql://localhost:5432/arda \
       -user=arda \
       -password=arda_password \
       info

# Validate migrations
flyway -url=jdbc:postgresql://localhost:5432/arda \
       -user=arda \
       -password=arda_password \
       validate

# Baseline existing database
flyway -url=jdbc:postgresql://localhost:5432/arda \
       -user=arda \
       -password=arda_password \
       baseline
```

## 📦 Configuration

### flyway.conf

```
flyway.url=jdbc:postgresql://localhost:5432/arda
flyway.user=arda
flyway.password=arda_password
flyway.locations=classpath:db/migration
flyway.baselineOnMigrate=true
flyway.baselineVersion=0
flyway.validateOnMigrate=true
flyway.outOfOrder=false
flyway.table=schema_version
```

### application.properties (Spring Boot)

```properties
# Flyway Configuration
spring.flyway.enabled=true
spring.flyway.url=jdbc:postgresql://localhost:5432/arda
spring.flyway.user=arda
spring.flyway.password=arda_password
spring.flyway.locations=classpath:db/migration
spring.flyway.baseline-on-migrate=true
spring.flyway.baseline-version=0
spring.flyway.validate-on-migrate=true
spring.flyway.out-of-order=false
```

## 📦 Naming Convention

### Migration File Names

```
V{version}__{description}.sql

Examples:
- V001__create_journals_table.sql
- V002__add_journals_attachment_column.sql
- V003__create_journal_items_table.sql
- V004__seed_journal_status_values.sql
```

### Rollback File Names

```
R{version}__{description}.sql

Examples:
- R003__drop_journal_items_table.sql
- R002__remove_journals_attachment_column.sql
- R001__drop_journals_table.sql
```

## 📦 Best Practices

### Migration Guidelines

- Use version numbers sequentially (V001, V002, etc.)
- Keep migrations reversible when possible
- Use transactions (BEGIN; COMMIT;)
- Test migrations on staging first
- Document changes with comments
- Don't modify existing migrations
- Use idempotent SQL (IF NOT EXISTS)
- Keep transactions short

### Database Design

- Use UUID for primary keys
- Use TIMESTAMPTZ for timestamps
- Use NUMERIC for monetary values
- Add indexes for foreign keys
- Use proper constraints (CHECK, UNIQUE)
- Enable Row Level Security for multi-tenancy
- Add audit columns (created_at, updated_at, etc.)

### Performance

- Add indexes for frequently queried columns
- Use composite indexes for multi-column queries
- Use partial indexes for filtered queries
- Keep indexes selective
- Analyze query performance with EXPLAIN ANALYZE

## 🎯 Usage Examples

```
/database-migration "Tạo migration mới"

Usage:
/database-migration "Tạo migration cho accounts table với các columns: id, code, name, type"

Sẽ:
1. Tạo migration file
2. Write DDL SQL
3. Tạo indexes
4. Setup RLS
```

```
/database-migration "Add column vào table"

Usage:
/database-migration "Tạo migration thêm column currency vào journals table"

Sẽ:
1. Tạo migration file
2. Add ALTER TABLE statement
3. Tạo index cho column mới
```

```
/database-migration "Tạo rollback"

Usage:
/database-migration "Tạo rollback script cho migration V003__create_journal_items_table.sql"

Sẽ:
1. Tạo rollback file
2. Write DROP statements
3. Remove indexes và constraints
```

## 📦 Migration Checklist

Before creating a migration:
- [ ] Review existing schema
- [ ] Test SQL in development database
- [ ] Use appropriate version number
- [ ] Add descriptive comments
- [ ] Include RLS policies
- [ ] Create necessary indexes
- [ ] Add constraints
- [ ] Write rollback script
- [ ] Test rollback

After running migration:
- [ ] Verify table created correctly
- [ ] Check indexes created
- [ ] Verify RLS policies working
- [ ] Test application with new schema
- [ ] Document changes

---

*Last Updated: 2026-04-25*
