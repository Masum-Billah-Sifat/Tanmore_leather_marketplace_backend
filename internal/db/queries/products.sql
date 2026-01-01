-- name: InsertProduct :one
INSERT INTO products (
    id,
    seller_id,
    category_id,
    title,
    description,
    is_approved,
    is_banned,
    is_archived,
    created_at,
    updated_at
) VALUES (
    $1,  $2,  $3,  $4,  $5,
    $6,  $7,  $8,  $9,  $10
)
RETURNING id;

-- name: GetProductByIDAndSellerID :one
SELECT
    id,
    seller_id,
    category_id,
    title,
    description,
    is_approved,
    is_banned,
    is_archived,
    created_at,
    updated_at
FROM products
WHERE id = $1
  AND seller_id = $2;

-- name: UpdateProductTitleDesc :exec
UPDATE products
SET
    title       = COALESCE(sqlc.narg(title)::TEXT, title),
    description = COALESCE(sqlc.narg(description)::TEXT, description),
    updated_at  = sqlc.arg(updated_at)
WHERE id = sqlc.arg(product_id)
  AND seller_id = sqlc.arg(seller_id);


-- name: ArchiveProduct :exec
UPDATE products
SET is_archived = $3,
    updated_at = $4
WHERE id = $1 AND seller_id = $2;


-- name: GetProductByID :one
SELECT
    id,
    seller_id,
    category_id,
    title,
    description,
    is_approved,
    is_banned,
    is_archived,
    created_at,
    updated_at
FROM products
WHERE id = $1;