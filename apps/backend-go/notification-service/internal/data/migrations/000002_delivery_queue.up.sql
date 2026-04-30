CREATE TABLE IF NOT EXISTS notification_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source_service TEXT NOT NULL,
  event_type TEXT NOT NULL,
  correlation_id TEXT NOT NULL DEFAULT '',
  idempotency_key TEXT NOT NULL UNIQUE,
  template_code TEXT NOT NULL,
  recipient_type TEXT NOT NULL DEFAULT 'USER',
  recipient_id TEXT NOT NULL,
  recipient_address TEXT NOT NULL DEFAULT '',
  channels TEXT[] NOT NULL DEFAULT ARRAY['IN_APP']::TEXT[],
  language TEXT NOT NULL DEFAULT 'vi',
  payload_json TEXT NOT NULL DEFAULT '{}',
  priority INT NOT NULL DEFAULT 100,
  status TEXT NOT NULL DEFAULT 'QUEUED',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notification_requests_recipient
  ON notification_requests(recipient_type, recipient_id, created_at DESC);

CREATE TABLE IF NOT EXISTS notification_deliveries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id UUID NOT NULL REFERENCES notification_requests(id),
  template_version_id UUID NOT NULL REFERENCES notification_template_versions(id),
  channel TEXT NOT NULL,
  recipient_type TEXT NOT NULL,
  recipient_id TEXT NOT NULL,
  recipient_address TEXT NOT NULL DEFAULT '',
  subject TEXT NOT NULL DEFAULT '',
  body TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'QUEUED',
  attempt_count INT NOT NULL DEFAULT 0,
  max_attempts INT NOT NULL DEFAULT 5,
  next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  locked_by TEXT NOT NULL DEFAULT '',
  locked_at TIMESTAMPTZ,
  provider_code TEXT NOT NULL DEFAULT '',
  provider_message_id TEXT NOT NULL DEFAULT '',
  provider_response_json TEXT NOT NULL DEFAULT '{}',
  error_message TEXT NOT NULL DEFAULT '',
  priority INT NOT NULL DEFAULT 100,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notification_deliveries_due
  ON notification_deliveries(status, next_attempt_at, priority, created_at);

CREATE INDEX IF NOT EXISTS idx_notification_deliveries_request_id
  ON notification_deliveries(request_id);

CREATE INDEX IF NOT EXISTS idx_notification_deliveries_recipient
  ON notification_deliveries(recipient_type, recipient_id, created_at DESC);

CREATE TABLE IF NOT EXISTS in_app_notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  delivery_id UUID NOT NULL UNIQUE REFERENCES notification_deliveries(id),
  recipient_type TEXT NOT NULL,
  recipient_id TEXT NOT NULL,
  title TEXT NOT NULL DEFAULT '',
  body TEXT NOT NULL DEFAULT '',
  data_json TEXT NOT NULL DEFAULT '{}',
  status TEXT NOT NULL DEFAULT 'UNREAD',
  read_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_in_app_notifications_recipient
  ON in_app_notifications(recipient_type, recipient_id, status, created_at DESC);
