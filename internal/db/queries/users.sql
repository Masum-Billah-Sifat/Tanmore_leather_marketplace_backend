-- name: GetUserByGoogleID :one
SELECT
    id,
    google_id,
    primary_email,
    display_name,
    profile_image_url,
    is_archived,
    is_banned,
    is_muted,
    current_mode,
    is_seller_profile_approved,
    is_seller_profile_created,
    created_at,
    updated_at
FROM users
WHERE google_id = $1;

-- name: InsertUser :one
INSERT INTO users (
  id, google_id, primary_email, display_name, profile_image_url,
  is_archived, is_banned, is_muted,
  current_mode, is_seller_profile_approved, is_seller_profile_created,
  created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5,
  $6, $7, $8,
  $9, $10, $11,
  $12, $13
)
RETURNING id;

-- name: GetUserByID :one
SELECT
    id,
    google_id,
    primary_email,
    display_name,
    profile_image_url,
    is_archived,
    is_banned,
    is_muted,
    current_mode,
    is_seller_profile_approved,
    is_seller_profile_created,
    created_at,
    updated_at
FROM users
WHERE id = $1;


-- name: UpdateUserCurrentMode :exec
UPDATE users
SET current_mode = $2,
    updated_at = $3
WHERE id = $1;