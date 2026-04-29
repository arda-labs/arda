-- Keep the local/bootstrap Zitadel administrator mapped to platform_admin even
-- when the Zitadel subject changes between environments.
INSERT INTO platform_admins (user_id, created_by)
SELECT id, id
FROM users
WHERE deleted_at IS NULL
  AND (
    external_id IN ('369593749817000033', '370594161885970517')
    OR lower(username) IN ('zitadel-admin', 'zitadel-admin@zitadel.auth.arda.io.vn')
    OR lower(email) = 'zitadel-admin@zitadel.auth.arda.io.vn'
  )
ON CONFLICT (user_id) DO UPDATE
SET revoked_at = NULL;
