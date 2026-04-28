ALTER TABLE users ADD COLUMN IF NOT EXISTS username TEXT NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_unique_active
  ON users (lower(username))
  WHERE deleted_at IS NULL AND username <> '';

UPDATE users
SET username = 'zitadel-admin',
    email = CASE
      WHEN email = '' OR email = 'admin@zitadel.auth.arda.io.vn' THEN 'zitadel-admin@zitadel.auth.arda.io.vn'
      ELSE email
    END,
    updated_at = now()
WHERE external_id = '369593749817000033'
  AND (username = '' OR username = 'admin');

CREATE TABLE IF NOT EXISTS platform_admins (
  user_id    UUID PRIMARY KEY REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_by UUID REFERENCES users(id),
  revoked_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_platform_admins_active
  ON platform_admins (user_id)
  WHERE revoked_at IS NULL;

INSERT INTO platform_admins (user_id, created_by)
SELECT id, id
FROM users
WHERE external_id = '369593749817000033'
ON CONFLICT (user_id) DO NOTHING;
