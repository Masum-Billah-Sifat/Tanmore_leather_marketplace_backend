-- name: InsertCartItem :exec
INSERT INTO cart_items (
  id,
  user_id,
  variant_id,
  required_quantity,
  is_active,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);


-- name: ReactivateCartItemByID :exec
UPDATE cart_items
SET
    required_quantity = $1,
    is_active = $2,
    updated_at = $3
WHERE id = $4;


-- name: ListActiveVariantIDsByUser :many
SELECT variant_id
FROM cart_items
WHERE
    is_active = true
    AND user_id = sqlc.arg(user_id);



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
WHERE
    ci.is_active = true
    AND ci.user_id = sqlc.arg(user_id)
    AND ci.variant_id = ANY(sqlc.arg(variant_ids)::UUID[]);




-- name: UpdateCartItemQuantityByID :exec
UPDATE cart_items
SET
    required_quantity = sqlc.arg(required_quantity),
    is_active = sqlc.arg(is_active),
    updated_at = sqlc.arg(updated_at)
WHERE
    id = sqlc.arg(id);


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
WHERE
  user_id = sqlc.arg(user_id)
  AND variant_id = sqlc.arg(variant_id);


-- name: UpdateCartQuantityForUser :exec
UPDATE cart_items
SET
  required_quantity = sqlc.arg(required_quantity),
  updated_at = sqlc.arg(updated_at)
WHERE
  user_id = sqlc.arg(user_id)
  AND variant_id = sqlc.arg(variant_id);


-- name: DeactivateCartItemByUserAndVariant :exec
UPDATE cart_items
SET is_active = FALSE,
    updated_at = $3
WHERE user_id = $1
  AND variant_id = $2
  AND is_active = TRUE;


-- name: ClearCartItemsByUser :exec
UPDATE cart_items
SET is_active = false,
    updated_at = $2
WHERE user_id = $1 AND is_active = true;
