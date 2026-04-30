CREATE TABLE IF NOT EXISTS notification_provider_configs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT NOT NULL UNIQUE,
  channel TEXT NOT NULL,
  name TEXT NOT NULL,
  priority INT NOT NULL DEFAULT 100,
  rate_limit_per_minute INT NOT NULL DEFAULT 0,
  options_json TEXT NOT NULL DEFAULT '{}',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notification_provider_configs_channel
  ON notification_provider_configs(channel, status, priority);

INSERT INTO notification_provider_configs (code, channel, name, priority, rate_limit_per_minute, options_json, status)
VALUES
  ('IN_APP_STORE', 'IN_APP', 'Arda in-app notification store', 10, 0, '{"adapter":"database"}', 'ACTIVE')
ON CONFLICT (code) DO NOTHING;
