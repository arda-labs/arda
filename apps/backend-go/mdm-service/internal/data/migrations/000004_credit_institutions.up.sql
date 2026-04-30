CREATE TABLE IF NOT EXISTS credit_institutions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code           TEXT NOT NULL,
    name           TEXT NOT NULL,
    short_name     TEXT NOT NULL DEFAULT '',
    address        TEXT NOT NULL DEFAULT '',
    phone          TEXT NOT NULL DEFAULT '',
    email          TEXT NOT NULL DEFAULT '',
    license_number TEXT NOT NULL DEFAULT '',
    issued_date    DATE,
    tax_code       TEXT NOT NULL DEFAULT '',
    website        TEXT NOT NULL DEFAULT '',
    note           TEXT NOT NULL DEFAULT '',
    status         TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_credit_institutions_code_active
    ON credit_institutions (code)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_credit_institutions_status
    ON credit_institutions (status)
    WHERE deleted_at IS NULL;
