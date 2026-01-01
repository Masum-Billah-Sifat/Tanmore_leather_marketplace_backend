-- name: InsertProductMedia :one
INSERT INTO product_medias (
  id,
  product_id,
  media_type,
  media_url,
  is_primary,
  is_archived,
  created_at,
  updated_at
)
VALUES (
  $1,  -- id (UUID)
  $2,  -- product_id (UUID)
  $3,  -- media_type (image or video)
  $4,  -- media_url
  $5,  -- is_primary
  $6,  -- is_archived
  $7,  -- created_at
  $8   -- updated_at
)
RETURNING id;


-- name: GetActiveMediasByProductID :many
SELECT
  id,
  product_id,
  media_type,
  media_url,
  is_primary,
  is_archived,
  created_at,
  updated_at
FROM product_medias
WHERE product_id = $1
  AND media_type = $2
  AND is_archived = $3
ORDER BY created_at ASC;

-- name: GetPrimaryProductImageByProductID :one
SELECT
    id,
    product_id,
    media_type,
    media_url,
    is_primary,
    is_archived,
    created_at,
    updated_at
FROM product_medias
WHERE product_id = $1
  AND media_type = $2
  AND is_primary = $3
  AND is_archived = $4
LIMIT 1;

-- name: GetPromoVideoByProductID :one
SELECT id
FROM product_medias
WHERE product_id = $1
  AND media_type = $2
  AND is_archived = $3
LIMIT 1;

-- name: ArchiveProductMedia :exec
UPDATE product_medias
SET is_archived = $4
WHERE id = $1 AND product_id = $2 AND media_type = $3;


-- name: CountActiveImagesForProduct :one
SELECT COUNT(*) FROM product_medias
WHERE product_id = $1 AND media_type = $2 AND is_archived = $3;

-- name: GetProductMediaByID :one
SELECT
  id,
  product_id,
  media_type,
  media_url,
  is_primary,
  is_archived
FROM product_medias
WHERE id = $1 AND product_id = $2 AND media_type = $3;

-- name: UnsetAllPrimaryImages :exec
UPDATE product_medias
SET is_primary = $4
WHERE product_id = $1 AND media_type = $2 AND is_archived = $3;

-- name: SetAsPrimaryImage :exec
UPDATE product_medias
SET is_primary = $4
WHERE id = $1 AND product_id = $2 AND media_type = $3;
