-- name: InsertUserSession :one
INSERT INTO user_sessions (
  id, user_id, ip_address, user_agent, device_fingerprint,
  is_revoked, is_archived, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5,
  $6, $7, $8, $9
)
RETURNING id;


-- name: GetSessionByIDAndUserID :one
SELECT
  id,
  user_id,
  ip_address,
  user_agent,
  device_fingerprint,
  is_revoked,
  is_archived,
  created_at,
  updated_at
FROM user_sessions
WHERE id = $1
  AND user_id = $2;


-- name: RevokeUserSession :exec
UPDATE user_sessions
SET is_revoked = $4,
    updated_at = $3
WHERE id = $1 AND user_id = $2;