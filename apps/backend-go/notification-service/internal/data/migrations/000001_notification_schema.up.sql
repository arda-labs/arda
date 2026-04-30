CREATE TABLE IF NOT EXISTS notification_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT '',
  default_channel TEXT NOT NULL DEFAULT 'IN_APP',
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS notification_template_versions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  template_id UUID NOT NULL REFERENCES notification_templates(id),
  version INT NOT NULL,
  channel TEXT NOT NULL,
  language TEXT NOT NULL DEFAULT 'vi',
  subject TEXT NOT NULL DEFAULT '',
  body TEXT NOT NULL,
  payload_schema_json TEXT NOT NULL DEFAULT '{}',
  approval_status TEXT NOT NULL DEFAULT 'DRAFT',
  approved_by TEXT NOT NULL DEFAULT '',
  approved_at TIMESTAMPTZ,
  change_note TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,
  UNIQUE (template_id, version, channel, language)
);

CREATE INDEX IF NOT EXISTS idx_notification_template_versions_template_id
  ON notification_template_versions(template_id)
  WHERE deleted_at IS NULL;

INSERT INTO notification_templates (code, name, category, default_channel, description, status)
VALUES
  ('IAM_SECURITY_LOGIN', 'Cảnh báo đăng nhập', 'SECURITY', 'IN_APP', 'Thông báo đăng nhập bảo mật cho người dùng nội bộ', 'ACTIVE'),
  ('SYSTEM_JOB_FAILED', 'Cảnh báo job hệ thống lỗi', 'OPERATIONS', 'EMAIL', 'Thông báo vận hành khi job hệ thống lỗi', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO notification_template_versions (template_id, version, channel, language, subject, body, payload_schema_json, approval_status, approved_by, approved_at, change_note, status)
SELECT t.id, 1, item.channel, item.language, item.subject, item.body, item.payload_schema_json, 'APPROVED', 'SYSTEM', now(), 'Seed template', 'ACTIVE'
FROM notification_templates t
JOIN (
  VALUES
    ('IAM_SECURITY_LOGIN', 'IN_APP', 'vi', 'Cảnh báo đăng nhập', 'Tài khoản của bạn vừa đăng nhập lúc {{login_time}} từ {{ip_address}}.', '{"login_time":"string","ip_address":"string"}'),
    ('SYSTEM_JOB_FAILED', 'EMAIL', 'vi', '[Arda] Job hệ thống lỗi', 'Job {{job_name}} lỗi lúc {{failed_at}}. Mã lỗi: {{error_code}}.', '{"job_name":"string","failed_at":"string","error_code":"string"}')
) AS item(template_code, channel, language, subject, body, payload_schema_json)
  ON t.code = item.template_code
ON CONFLICT DO NOTHING;
