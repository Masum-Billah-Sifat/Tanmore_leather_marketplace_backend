-- name: InsertRefreshToken :exec
INSERT INTO user_refresh_tokens (
  id, user_id, session_id, token_hash,
  deprecated_reason, is_deprecated, deprecated_at,
  expires_at, created_at
) VALUES (
  $1, $2, $3, $4,
  $5, $6, $7,
  $8, $9
);

-- name: GetRefreshTokenByHash :one
SELECT
  id,
  user_id,
  session_id,
  token_hash,
  deprecated_reason,
  is_deprecated,
  deprecated_at,
  expires_at,
  created_at
FROM user_refresh_tokens
WHERE token_hash = $1;


-- name: DeprecateRefreshTokenByID :exec
UPDATE user_refresh_tokens
SET
  is_deprecated     = $2,
  deprecated_reason = $3,
  deprecated_at     = $4
WHERE id = $1;


-- name: DeprecateRefreshTokensBySession :exec
UPDATE user_refresh_tokens
SET
  is_deprecated = $5,
  deprecated_reason = $4,
  deprecated_at = $3
WHERE user_id = $1 AND session_id = $2;