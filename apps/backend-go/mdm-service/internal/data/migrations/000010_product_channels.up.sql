CREATE TABLE IF NOT EXISTS banking_products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  product_type TEXT NOT NULL DEFAULT 'ACCOUNT',
  category TEXT NOT NULL DEFAULT '',
  customer_segment TEXT NOT NULL DEFAULT '',
  currency TEXT NOT NULL DEFAULT '',
  effective_from DATE,
  effective_to DATE,
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS service_channels (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  channel_type TEXT NOT NULL DEFAULT 'DIGITAL',
  availability TEXT NOT NULL DEFAULT '24X7',
  timezone TEXT NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS product_channel_rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_code TEXT NOT NULL,
  channel_code TEXT NOT NULL,
  transaction_type TEXT NOT NULL DEFAULT '',
  enabled BOOLEAN NOT NULL DEFAULT true,
  priority INT NOT NULL DEFAULT 100,
  fee_schedule_code TEXT NOT NULL DEFAULT '',
  limit_profile_code TEXT NOT NULL DEFAULT '',
  effective_from DATE,
  effective_to DATE,
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,
  UNIQUE (product_code, channel_code, transaction_type)
);

INSERT INTO banking_products (code, name, product_type, category, customer_segment, currency, effective_from, description, status)
VALUES
  ('CASA_RETAIL', 'Tài khoản thanh toán cá nhân', 'ACCOUNT', 'CASA', 'RETAIL', 'VND', '2026-01-01', 'Sản phẩm tài khoản thanh toán tiêu chuẩn cho khách hàng cá nhân', 'ACTIVE'),
  ('SAVINGS_TERM', 'Tiền gửi tiết kiệm có kỳ hạn', 'DEPOSIT', 'SAVINGS', 'RETAIL', 'VND', '2026-01-01', 'Tiền gửi có kỳ hạn dùng cho retail và priority banking', 'ACTIVE'),
  ('CARD_DEBIT', 'Thẻ ghi nợ nội địa', 'CARD', 'DEBIT', 'RETAIL', 'VND', '2026-01-01', 'Thẻ ghi nợ liên kết tài khoản thanh toán', 'ACTIVE'),
  ('SME_CURRENT', 'Tài khoản thanh toán doanh nghiệp SME', 'ACCOUNT', 'CASA', 'SME', 'VND', '2026-01-01', 'Tài khoản thanh toán dành cho doanh nghiệp SME', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO service_channels (code, name, channel_type, availability, timezone, description, status)
VALUES
  ('BRANCH', 'Quầy giao dịch', 'BRANCH', 'BUSINESS_HOURS', 'Asia/Ho_Chi_Minh', 'Giao dịch tại chi nhánh/phòng giao dịch', 'ACTIVE'),
  ('INTERNET_BANKING', 'Internet Banking', 'DIGITAL', '24X7', 'Asia/Ho_Chi_Minh', 'Kênh ngân hàng điện tử trên trình duyệt', 'ACTIVE'),
  ('MOBILE_BANKING', 'Mobile Banking', 'DIGITAL', '24X7', 'Asia/Ho_Chi_Minh', 'Kênh ngân hàng số trên thiết bị di động', 'ACTIVE'),
  ('ATM', 'ATM', 'SELF_SERVICE', '24X7', 'Asia/Ho_Chi_Minh', 'Thiết bị tự phục vụ ATM', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO product_channel_rules (product_code, channel_code, transaction_type, enabled, priority, fee_schedule_code, limit_profile_code, effective_from, description, status)
VALUES
  ('CASA_RETAIL', 'BRANCH', 'TRANSFER', true, 10, 'BRANCH_TRANSFER_FEE', 'RETAIL_DAILY_TRANSFER', '2026-01-01', 'Cho phép chuyển khoản tại quầy', 'ACTIVE'),
  ('CASA_RETAIL', 'MOBILE_BANKING', 'TRANSFER', true, 20, 'DIGITAL_TRANSFER_FEE', 'RETAIL_DAILY_TRANSFER', '2026-01-01', 'Cho phép chuyển khoản trên mobile banking', 'ACTIVE'),
  ('CASA_RETAIL', 'INTERNET_BANKING', 'TRANSFER', true, 30, 'DIGITAL_TRANSFER_FEE', 'RETAIL_DAILY_TRANSFER', '2026-01-01', 'Cho phép chuyển khoản trên internet banking', 'ACTIVE'),
  ('CARD_DEBIT', 'ATM', 'WITHDRAWAL', true, 10, 'ATM_WITHDRAWAL_FEE', 'CARD_ATM_DAILY', '2026-01-01', 'Cho phép rút tiền ATM cho thẻ ghi nợ', 'ACTIVE')
ON CONFLICT DO NOTHING;
