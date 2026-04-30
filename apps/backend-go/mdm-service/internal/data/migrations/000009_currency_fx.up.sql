CREATE TABLE IF NOT EXISTS currencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    numeric_code TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    minor_unit INTEGER NOT NULL DEFAULT 0,
    symbol TEXT NOT NULL DEFAULT '',
    country_code TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_currencies_code_active
    ON currencies (code)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS fx_rate_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    source_type TEXT NOT NULL DEFAULT 'MANUAL',
    priority INTEGER NOT NULL DEFAULT 100,
    timezone TEXT NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_fx_rate_sources_code_active
    ON fx_rate_sources (code)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS fx_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_currency TEXT NOT NULL,
    quote_currency TEXT NOT NULL,
    source_code TEXT NOT NULL,
    rate_date DATE NOT NULL,
    effective_at TIMESTAMPTZ,
    buy_rate NUMERIC(20, 8) NOT NULL DEFAULT 0,
    sell_rate NUMERIC(20, 8) NOT NULL DEFAULT 0,
    mid_rate NUMERIC(20, 8) NOT NULL DEFAULT 0,
    approval_status TEXT NOT NULL DEFAULT 'DRAFT',
    version INTEGER NOT NULL DEFAULT 1,
    approved_by TEXT NOT NULL DEFAULT '',
    approved_at TIMESTAMPTZ,
    change_note TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_fx_rates_pair_source_date_version_active
    ON fx_rates (base_currency, quote_currency, source_code, rate_date, version)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS ix_fx_rates_lookup
    ON fx_rates (rate_date DESC, base_currency, quote_currency, source_code, status)
    WHERE deleted_at IS NULL;

INSERT INTO currencies (code, numeric_code, name, minor_unit, symbol, country_code, status)
VALUES
    ('VND', '704', 'Vietnamese Dong', 0, '₫', 'VN', 'ACTIVE'),
    ('USD', '840', 'US Dollar', 2, '$', 'US', 'ACTIVE'),
    ('EUR', '978', 'Euro', 2, '€', 'EU', 'ACTIVE'),
    ('JPY', '392', 'Japanese Yen', 0, '¥', 'JP', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO fx_rate_sources (code, name, source_type, priority, timezone, description, status)
VALUES
    ('SBV', 'Ngân hàng Nhà nước Việt Nam', 'CENTRAL_BANK', 10, 'Asia/Ho_Chi_Minh', 'Nguồn tỷ giá tham chiếu chính thức.', 'ACTIVE'),
    ('TREASURY', 'Treasury nội bộ', 'INTERNAL', 20, 'Asia/Ho_Chi_Minh', 'Nguồn tỷ giá vận hành nội bộ.', 'ACTIVE'),
    ('MANUAL', 'Nhập thủ công', 'MANUAL', 100, 'Asia/Ho_Chi_Minh', 'Nguồn tỷ giá nhập tay.', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO fx_rates (
    base_currency, quote_currency, source_code, rate_date, effective_at,
    buy_rate, sell_rate, mid_rate, approval_status, version, approved_by, approved_at,
    change_note, status
)
VALUES
    ('USD', 'VND', 'TREASURY', '2026-01-01', '2026-01-01 00:00:00+07', 25000, 25500, 25250, 'APPROVED', 1, 'SYSTEM', now(), 'Seed tỷ giá tham chiếu.', 'ACTIVE'),
    ('EUR', 'VND', 'TREASURY', '2026-01-01', '2026-01-01 00:00:00+07', 27000, 27800, 27400, 'APPROVED', 1, 'SYSTEM', now(), 'Seed tỷ giá tham chiếu.', 'ACTIVE')
ON CONFLICT DO NOTHING;
