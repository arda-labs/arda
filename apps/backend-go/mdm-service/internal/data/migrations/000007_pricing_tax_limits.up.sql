CREATE TABLE IF NOT EXISTS fee_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    fee_type TEXT NOT NULL DEFAULT 'SERVICE_FEE',
    calculation_method TEXT NOT NULL DEFAULT 'FIXED',
    currency TEXT NOT NULL DEFAULT 'VND',
    fixed_amount NUMERIC(20, 4) NOT NULL DEFAULT 0,
    rate_percent NUMERIC(9, 6) NOT NULL DEFAULT 0,
    min_amount NUMERIC(20, 4) NOT NULL DEFAULT 0,
    max_amount NUMERIC(20, 4) NOT NULL DEFAULT 0,
    channel TEXT NOT NULL DEFAULT '',
    product_code TEXT NOT NULL DEFAULT '',
    effective_from DATE,
    effective_to DATE,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_fee_schedules_code_active
    ON fee_schedules (code)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS ix_fee_schedules_lookup
    ON fee_schedules (status, fee_type, channel, product_code, currency)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS tax_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    tax_type TEXT NOT NULL DEFAULT 'VAT',
    rate_percent NUMERIC(9, 6) NOT NULL DEFAULT 0,
    inclusive BOOLEAN NOT NULL DEFAULT false,
    jurisdiction TEXT NOT NULL DEFAULT 'VN',
    effective_from DATE,
    effective_to DATE,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_tax_rules_code_active
    ON tax_rules (code)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS ix_tax_rules_lookup
    ON tax_rules (status, tax_type, jurisdiction)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS standard_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    limit_type TEXT NOT NULL DEFAULT 'TRANSACTION_AMOUNT',
    subject_type TEXT NOT NULL DEFAULT 'CUSTOMER',
    currency TEXT NOT NULL DEFAULT 'VND',
    min_amount NUMERIC(20, 4) NOT NULL DEFAULT 0,
    per_txn_amount NUMERIC(20, 4) NOT NULL DEFAULT 0,
    daily_amount NUMERIC(20, 4) NOT NULL DEFAULT 0,
    monthly_amount NUMERIC(20, 4) NOT NULL DEFAULT 0,
    count_limit INTEGER NOT NULL DEFAULT 0,
    channel TEXT NOT NULL DEFAULT '',
    product_code TEXT NOT NULL DEFAULT '',
    effective_from DATE,
    effective_to DATE,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_standard_limits_code_active
    ON standard_limits (code)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS ix_standard_limits_lookup
    ON standard_limits (status, limit_type, subject_type, channel, product_code, currency)
    WHERE deleted_at IS NULL;

INSERT INTO fee_schedules (
    code, name, fee_type, calculation_method, currency, fixed_amount, rate_percent,
    min_amount, max_amount, channel, product_code, effective_from, description, status
)
VALUES
    ('FEE_INTERNAL_TRANSFER_VND', 'Phí chuyển khoản nội bộ VND', 'TRANSFER_FEE', 'FIXED', 'VND', 0, 0, 0, 0, 'DIGITAL', 'CASA', '2026-01-01', 'Biểu phí mặc định cho chuyển khoản nội bộ.', 'ACTIVE'),
    ('FEE_INTERBANK_TRANSFER_VND', 'Phí chuyển khoản liên ngân hàng VND', 'TRANSFER_FEE', 'FIXED', 'VND', 3300, 0, 0, 0, 'DIGITAL', 'CASA', '2026-01-01', 'Biểu phí mặc định cho chuyển khoản liên ngân hàng.', 'ACTIVE'),
    ('FEE_CASH_WITHDRAWAL_ATM_VND', 'Phí rút tiền ATM VND', 'WITHDRAWAL_FEE', 'FIXED', 'VND', 1100, 0, 0, 0, 'ATM', 'DEBIT_CARD', '2026-01-01', 'Biểu phí rút tiền ATM nội địa mặc định.', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO tax_rules (
    code, name, tax_type, rate_percent, inclusive, jurisdiction, effective_from, description, status
)
VALUES
    ('TAX_VAT_STANDARD_VN', 'Thuế VAT tiêu chuẩn Việt Nam', 'VAT', 10, false, 'VN', '2026-01-01', 'Quy tắc thuế VAT mặc định áp dụng cho phí dịch vụ chịu thuế.', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO standard_limits (
    code, name, limit_type, subject_type, currency, min_amount, per_txn_amount, daily_amount,
    monthly_amount, count_limit, channel, product_code, effective_from, description, status
)
VALUES
    ('LIMIT_DIGITAL_TRANSFER_RETAIL_VND', 'Hạn mức chuyển khoản số khách hàng cá nhân', 'TRANSFER_AMOUNT', 'RETAIL_CUSTOMER', 'VND', 10000, 500000000, 2000000000, 20000000000, 0, 'DIGITAL', 'CASA', '2026-01-01', 'Hạn mức chuẩn tham chiếu cho chuyển khoản số khách hàng cá nhân.', 'ACTIVE'),
    ('LIMIT_ATM_WITHDRAWAL_DEBIT_VND', 'Hạn mức rút tiền ATM thẻ ghi nợ', 'WITHDRAWAL_AMOUNT', 'CARD', 'VND', 50000, 5000000, 50000000, 0, 20, 'ATM', 'DEBIT_CARD', '2026-01-01', 'Hạn mức chuẩn tham chiếu cho rút tiền ATM.', 'ACTIVE')
ON CONFLICT DO NOTHING;
