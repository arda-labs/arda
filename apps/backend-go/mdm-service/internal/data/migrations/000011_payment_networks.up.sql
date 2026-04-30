CREATE TABLE IF NOT EXISTS bank_branches (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  institution_code TEXT NOT NULL,
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  branch_type TEXT NOT NULL DEFAULT 'BRANCH',
  address TEXT NOT NULL DEFAULT '',
  province_code TEXT NOT NULL DEFAULT '',
  phone TEXT NOT NULL DEFAULT '',
  swift_code TEXT NOT NULL DEFAULT '',
  napas_code TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS payment_networks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  network_type TEXT NOT NULL DEFAULT 'DOMESTIC',
  clearing_method TEXT NOT NULL DEFAULT '',
  settlement_currency TEXT NOT NULL DEFAULT 'VND',
  operator TEXT NOT NULL DEFAULT '',
  availability TEXT NOT NULL DEFAULT '24X7',
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

INSERT INTO bank_branches (institution_code, code, name, branch_type, address, province_code, phone, swift_code, napas_code, status)
VALUES
  ('VCB', 'VCB_HO', 'Vietcombank Hội sở chính', 'HEAD_OFFICE', '198 Trần Quang Khải, Hoàn Kiếm, Hà Nội', '01', '1900545413', 'BFTVVNVX', '970436', 'ACTIVE'),
  ('VCB', 'VCB_HCM', 'Vietcombank Chi nhánh TP Hồ Chí Minh', 'BRANCH', '5 Công trường Mê Linh, Quận 1, TP Hồ Chí Minh', '79', '1900545413', 'BFTVVNVX007', '970436', 'ACTIVE'),
  ('BIDV', 'BIDV_HO', 'BIDV Hội sở chính', 'HEAD_OFFICE', '35 Hàng Vôi, Hoàn Kiếm, Hà Nội', '01', '19009247', 'BIDVVNVX', '970418', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO payment_networks (code, name, network_type, clearing_method, settlement_currency, operator, availability, description, status)
VALUES
  ('NAPAS', 'Napas', 'DOMESTIC', 'NET_CLEARING', 'VND', 'NAPAS', '24X7', 'Mạng chuyển mạch thẻ và thanh toán nội địa Việt Nam', 'ACTIVE'),
  ('SWIFT', 'SWIFT', 'INTERNATIONAL', 'CORRESPONDENT', 'USD', 'SWIFT', 'BUSINESS_HOURS', 'Mạng điện tài chính quốc tế dùng cho thanh toán cross-border', 'ACTIVE'),
  ('IBFT', 'Chuyển tiền nhanh liên ngân hàng', 'DOMESTIC', 'REAL_TIME', 'VND', 'NAPAS', '24X7', 'Dịch vụ chuyển tiền nhanh liên ngân hàng', 'ACTIVE')
ON CONFLICT DO NOTHING;
