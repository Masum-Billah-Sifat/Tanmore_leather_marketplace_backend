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