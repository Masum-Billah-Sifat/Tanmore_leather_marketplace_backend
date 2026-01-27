-- -- name: GetCartItemByOwnerAndVariant :one
-- SELECT
--     id,
--     user_id,
--     guest_user_id,
--     variant_id,
--     required_quantity,
--     is_active,
--     is_deprecated,
--     created_at,
--     updated_at
-- FROM cart_items
-- WHERE
--     variant_id = sqlc.arg(variant_id)
--     AND is_deprecated = false
--     AND (
--         (user_id = sqlc.narg(user_id) AND sqlc.narg(guest_user_id) IS NULL)
--         OR
--         (guest_user_id = sqlc.narg(guest_user_id) AND sqlc.narg(user_id) IS NULL)
--     );


-- name: GetCartItemByUserAndVariant :one
SELECT
    id,
    user_id,
    guest_user_id,
    variant_id,
    required_quantity,
    is_active,
    is_deprecated,
    created_at,
    updated_at
FROM cart_items
WHERE
    variant_id = sqlc.arg(variant_id)
    AND is_deprecated = false
    AND user_id = sqlc.arg(user_id);


-- name: GetCartItemByGuestAndVariant :one
SELECT
    id,
    user_id,
    guest_user_id,
    variant_id,
    required_quantity,
    is_active,
    is_deprecated,
    created_at,
    updated_at
FROM cart_items
WHERE
    variant_id = sqlc.arg(variant_id)
    AND is_deprecated = false
    AND guest_user_id = sqlc.arg(guest_user_id);


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
    guest_user_id,
    variant_id,
    required_quantity,
    is_active,
    is_deprecated,
    created_at,
    updated_at
)
VALUES (
    sqlc.narg(user_id),
    sqlc.narg(guest_user_id),
    sqlc.arg(variant_id),
    sqlc.narg(required_quantity),
    true,
    false,
    sqlc.arg(created_at),
    sqlc.arg(updated_at)
)
RETURNING
    id,
    user_id,
    guest_user_id,
    variant_id,
    required_quantity,
    is_active,
    is_deprecated,
    created_at,
    updated_at;



-- -- name: UpdateCartQuantity :exec
-- UPDATE cart_items
-- SET
--     required_quantity = $1,
--     updated_at = $2
-- WHERE user_id = $3 AND variant_id = $4 AND is_active = $5;


-- -- name: UpdateCartQuantity :exec
-- UPDATE cart_items
-- SET
--     required_quantity = sqlc.narg(required_quantity),
--     updated_at = sqlc.arg(updated_at)
-- WHERE
--     variant_id = sqlc.arg(variant_id)
--     AND is_active = true
--     AND (
--         (user_id = sqlc.narg(user_id) AND sqlc.narg(guest_user_id) IS NULL)
--         OR
--         (guest_user_id = sqlc.narg(guest_user_id) AND sqlc.narg(user_id) IS NULL)
--     );



-- name: UpdateCartQuantityForUser :exec
UPDATE cart_items
SET
    required_quantity = sqlc.narg(required_quantity),
    updated_at = sqlc.arg(updated_at)
WHERE
    variant_id = sqlc.arg(variant_id)
    AND is_active = true
    AND user_id = sqlc.arg(user_id);


-- name: UpdateCartQuantityForGuest :exec
UPDATE cart_items
SET
    required_quantity = sqlc.narg(required_quantity),
    updated_at = sqlc.arg(updated_at)
WHERE
    variant_id = sqlc.arg(variant_id)
    AND is_active = true
    AND guest_user_id = sqlc.arg(guest_user_id);


-- -- name: DeactivateCartItem :exec
-- UPDATE cart_items
-- SET
--     required_quantity = $5,
--     is_active = $4,
--     updated_at = $3
-- WHERE user_id = $1 AND variant_id = $2;

-- -- name: DeactivateCartItemByOwnerAndVariant :exec
-- UPDATE cart_items
-- SET
--     required_quantity = NULL,
--     is_active = false,
--     updated_at = sqlc.arg(updated_at)
-- WHERE
--     variant_id = sqlc.arg(variant_id)
--     AND is_deprecated = false
--     AND (
--         (user_id = sqlc.narg(user_id) AND sqlc.narg(guest_user_id) IS NULL)
--         OR
--         (guest_user_id = sqlc.narg(guest_user_id) AND sqlc.narg(user_id) IS NULL)
--     );


-- name: DeactivateCartItemByUserAndVariant :exec
UPDATE cart_items
SET
    required_quantity = NULL,
    is_active = false,
    updated_at = sqlc.arg(updated_at)
WHERE
    variant_id = sqlc.arg(variant_id)
    AND is_deprecated = false
    AND user_id = sqlc.arg(user_id);


-- name: DeactivateCartItemByGuestAndVariant :exec
UPDATE cart_items
SET
    required_quantity = NULL,
    is_active = false,
    updated_at = sqlc.arg(updated_at)
WHERE
    variant_id = sqlc.arg(variant_id)
    AND is_deprecated = false
    AND guest_user_id = sqlc.arg(guest_user_id);


-- -- name: ClearCartItemsForUser :exec
-- UPDATE cart_items
-- SET 
--     required_quantity = $4,
--     is_active = $3,
--     updated_at = $2
-- WHERE user_id = $1;

-- -- name: ClearCartItemsByOwner :exec
-- UPDATE cart_items
-- SET
--     required_quantity = NULL,
--     is_active = false,
--     updated_at = sqlc.arg(updated_at)
-- WHERE
--     is_active = true
--     AND is_deprecated = false
--     AND (
--         (user_id = sqlc.narg(user_id) AND sqlc.narg(guest_user_id) IS NULL)
--         OR
--         (guest_user_id = sqlc.narg(guest_user_id) AND sqlc.narg(user_id) IS NULL)
--     );


-- name: ClearCartItemsByUser :exec
UPDATE cart_items
SET
    required_quantity = NULL,
    is_active = false,
    updated_at = sqlc.arg(updated_at)
WHERE
    is_active = true
    AND is_deprecated = false
    AND user_id = sqlc.arg(user_id);


-- name: ClearCartItemsByGuest :exec
UPDATE cart_items
SET
    required_quantity = NULL,
    is_active = false,
    updated_at = sqlc.arg(updated_at)
WHERE
    is_active = true
    AND is_deprecated = false
    AND guest_user_id = sqlc.arg(guest_user_id);


-- -- name: ListActiveVariantIDsByUser :many
-- SELECT variant_id FROM cart_items
-- WHERE user_id = sqlc.arg(user_id)
--   AND is_active = TRUE;

-- -- name: ListActiveVariantIDsByOwner :many
-- SELECT variant_id
-- FROM cart_items
-- WHERE
--     is_active = true
--     AND is_deprecated = false
--     AND (
--         (user_id = sqlc.narg(user_id) AND sqlc.narg(guest_user_id) IS NULL)
--         OR
--         (guest_user_id = sqlc.narg(guest_user_id) AND sqlc.narg(user_id) IS NULL)
--     );


-- name: ListActiveVariantIDsByUser :many
SELECT variant_id
FROM cart_items
WHERE
    is_active = true
    AND is_deprecated = false
    AND user_id = sqlc.arg(user_id);


-- name: ListActiveVariantIDsByGuest :many
SELECT variant_id
FROM cart_items
WHERE
    is_active = true
    AND is_deprecated = false
    AND guest_user_id = sqlc.arg(guest_user_id);



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
WHERE   ci.is_active = true
    AND ci.is_deprecated = false
    AND ci.variant_id = ANY(sqlc.arg(variant_ids)::UUID[])
    AND ci.user_id = sqlc.arg(user_id);


-- name: GetActiveCartVariantSnapshotsByGuestAndVariantIDs :many
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
WHERE ci.is_active = true
    AND ci.is_deprecated = false
    AND ci.variant_id = ANY(sqlc.arg(variant_ids)::UUID[])
    AND ci.guest_user_id = sqlc.arg(guest_user_id);



-- name: GetActiveGuestCartItems :many
SELECT
    id,
    guest_user_id,
    variant_id,
    required_quantity,
    created_at,
    updated_at
FROM cart_items
WHERE guest_user_id = $1
  AND is_active = TRUE
  AND is_deprecated = FALSE;


-- name: DeprecateGuestCartItems :exec
UPDATE cart_items
SET
    is_deprecated = $3,
    updated_at = $4
WHERE guest_user_id = $1
  AND is_active = $2
  AND is_deprecated = false;