-- name: GetCartItemByUserAndVariant :one
SELECT
    id,
    user_id,
    variant_id,
    required_quantity,
    is_active,
    created_at,
    updated_at
FROM cart_items
WHERE user_id = $1 AND variant_id = $2;

-- name: ReactivateCartItemByID :exec
UPDATE cart_items
SET
    required_quantity = $1,
    is_active = $2,
    updated_at = $3
WHERE id = $4;


-- name: InsertCartItem :one
INSERT INTO cart_items (
    user_id,
    variant_id,
    required_quantity,
    is_active,
    created_at,
    updated_at
)
VALUES (
    $1, -- user_id
    $2, -- variant_id
    $3, -- required_quantity
    $4, -- is_active
    $5, -- created_at
    $6  -- updated_at
)
RETURNING
    id,
    user_id,
    variant_id,
    required_quantity,
    is_active,
    created_at,
    updated_at;


-- name: UpdateCartQuantity :exec
UPDATE cart_items
SET
    required_quantity = $1,
    updated_at = $2
WHERE user_id = $3 AND variant_id = $4 AND is_active = $5;

-- name: DeactivateCartItem :exec
UPDATE cart_items
SET
    required_quantity = $5,
    is_active = $4,
    updated_at = $3
WHERE user_id = $1 AND variant_id = $2;

-- name: ClearCartItemsForUser :exec
UPDATE cart_items
SET 
    required_quantity = $4,
    is_active = $3,
    updated_at = $2
WHERE user_id = $1;

-- name: ListActiveVariantIDsByUser :many
SELECT variant_id FROM cart_items
WHERE user_id = sqlc.arg(user_id)
  AND is_active = TRUE;


-- name: GetActiveCartVariantSnapshotsByUserAndVariantIDs :many
SELECT
  -- Cart fields
  ci.id                AS cart_item_id,
  ci.user_id           AS cart_user_id,
  ci.variant_id        AS cart_variant_id,
  ci.required_quantity AS cart_required_quantity,
  ci.is_active         AS cart_is_active,
  ci.created_at        AS cart_created_at,
  ci.updated_at        AS cart_updated_at,

  -- Snapshot fields
  pvs.id                         AS snapshot_id,
  pvs.categoryid,
  pvs.iscategoryarchived,
  pvs.categoryname,
  pvs.sellerid,
  pvs.issellerapproved,
  pvs.issellerarchived,
  pvs.issellerbanned,
  pvs.sellerstorename,
  pvs.productid,
  pvs.isproductapproved,
  pvs.isproductarchived,
  pvs.isproductbanned,
  pvs.producttitle,
  pvs.productdescription,
  pvs.productprimaryimageurl,
  pvs.variantid,
  pvs.isvariantarchived,
  pvs.isvariantinstock,
  pvs.stockamount,
  pvs.color,
  pvs.size,
  pvs.retailprice,
  pvs.hasretaildiscount,
  pvs.retaildiscounttype,
  pvs.retaildiscount,
  pvs.haswholesaleenabled,
  pvs.wholesaleprice,
  pvs.wholesaleminquantity,
  pvs.haswholesalediscount,
  pvs.wholesalediscounttype,
  pvs.wholesalediscount,
  pvs.weight_grams,
  pvs.createdat       AS snapshot_created_at,
  pvs.updatedat       AS snapshot_updated_at

FROM cart_items ci
JOIN product_variant_snapshots pvs
  ON ci.variant_id = pvs.variantid
WHERE ci.user_id = sqlc.arg(user_id)
  AND ci.is_active = TRUE
  AND ci.variant_id = ANY(sqlc.arg(variant_ids)::UUID[]);
