-- name: InsertProductVariantReturningID :one
INSERT INTO product_variants (
    id,
    product_id,
    color,
    size,
    retail_price,
    retaildiscounttype,
    retaildiscount,
    wholesale_enabled,
    wholesale_price,
    min_qty_wholesale,
    wholesalediscounttype,
    wholesalediscount,
    stock_quantity,
    in_stock,
    weight_grams,
    is_archived,
    created_at,
    updated_at,
    hasretaildiscount,
    haswholesalediscount
) VALUES (
    $1,  $2,  $3,  $4,  $5,
    $6,  $7,  $8,  $9,  $10,
    $11, $12, $13, $14, $15,
    $16, $17, $18, $19, $20
)
RETURNING id;


-- name: ArchiveProductVariant :exec
UPDATE product_variants
SET is_archived = $3, updated_at = $4
WHERE id = $1 AND product_id = $2;


-- name: UpdateVariantColorSize :exec
UPDATE product_variants
SET
    color = COALESCE(sqlc.narg(color)::TEXT, color),
    size  = COALESCE(sqlc.narg(size)::TEXT, size),
    updated_at = sqlc.arg(updated_at)
WHERE id = sqlc.arg(variant_id)
  AND product_id = sqlc.arg(product_id);


-- name: UpdateVariantRetailPrice :exec
UPDATE product_variants
SET retail_price = $1, updated_at = $2
WHERE id = $3 AND product_id = $4;

-- name: UpdateVariantInStock :exec
UPDATE product_variants
SET in_stock = $1, updated_at = $2
WHERE id = $3 AND product_id = $4;

-- name: UpdateVariantStockQuantity :exec
UPDATE product_variants
SET stock_quantity = $1, updated_at = $2
WHERE id = $3 AND product_id = $4;

-- name: EnableRetailDiscount :exec
UPDATE product_variants
SET hasretaildiscount = $1,
    retaildiscounttype = $2,
    retaildiscount = $3,
    updated_at = $4
WHERE id = $5 AND product_id = $6;

-- name: UpdateRetailDiscount :exec
UPDATE product_variants
SET
    retaildiscounttype = COALESCE(sqlc.narg(retaildiscounttype)::TEXT, retaildiscounttype),
    retaildiscount     = COALESCE(sqlc.narg(retaildiscount)::BIGINT, retaildiscount),
    updated_at         = sqlc.arg(updated_at)
WHERE id = sqlc.arg(variant_id)
  AND product_id = sqlc.arg(product_id);


-- name: DisableRetailDiscount :exec
UPDATE product_variants
SET hasretaildiscount = $1,
    retaildiscounttype = $2,
    retaildiscount = $3,
    updated_at = $4
WHERE id = $5 AND product_id = $6;


-- name: EnableWholesaleMode :exec
UPDATE product_variants
SET wholesale_enabled = $1,
    wholesale_price = $2,
    min_qty_wholesale = $3,
    haswholesalediscount = $4,
    wholesalediscounttype = $5,
    wholesalediscount = $6,
    updated_at = $7
WHERE id = $8 AND product_id = $9;


-- name: UpdateWholesaleMode :exec
UPDATE product_variants
SET
    wholesale_price    = COALESCE(sqlc.narg(wholesale_price)::BIGINT, wholesale_price),
    min_qty_wholesale  = COALESCE(sqlc.narg(min_qty_wholesale)::INT, min_qty_wholesale),
    updated_at         = sqlc.arg(updated_at)
WHERE id = sqlc.arg(variant_id)
  AND product_id = sqlc.arg(product_id);


-- name: DisableWholesaleMode :exec
UPDATE product_variants
SET wholesale_enabled = $1,
    wholesale_price = $2,
    min_qty_wholesale = $3,
    haswholesalediscount = $4,
    wholesalediscounttype = $5,
    wholesalediscount = $6,
    updated_at = $7
WHERE id = $8 AND product_id = $9;


-- name: EnableWholesaleDiscount :exec
UPDATE product_variants
SET haswholesalediscount = $6,
    wholesalediscounttype = $5,
    wholesalediscount = $4,
    updated_at = $3
WHERE id = $1 AND product_id = $2;


-- name: UpdateWholesaleDiscount :exec
UPDATE product_variants
SET
    wholesalediscounttype = COALESCE(sqlc.narg(wholesalediscounttype)::TEXT, wholesalediscounttype),
    wholesalediscount     = COALESCE(sqlc.narg(wholesalediscount)::BIGINT, wholesalediscount),
    updated_at            = sqlc.arg(updated_at)
WHERE id = sqlc.arg(variant_id)
  AND product_id = sqlc.arg(product_id);



-- name: DisableWholesaleDiscount :exec
UPDATE product_variants
SET haswholesalediscount = $1,
    wholesalediscounttype = $2,
    wholesalediscount = $3,
    updated_at = $4
WHERE id = $5 AND product_id = $6;

-- name: UpdateVariantWeight :exec
UPDATE product_variants
SET weight_grams = $1, updated_at = $2
WHERE id = $3 AND product_id = $4;
