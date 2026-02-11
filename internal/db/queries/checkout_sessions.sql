-- -- name: InsertCheckoutSession :one
-- INSERT INTO checkout_sessions (
--     id,
--     user_id,
--     subtotal,
--     total_weight_grams,
--     delivery_charge,
--     total_payable,
--     shipping_address_id,
--     created_at
-- )
-- VALUES (
--     sqlc.arg(id),
--     sqlc.arg(user_id),
--     sqlc.arg(subtotal),
--     sqlc.arg(total_weight_grams),
--     sqlc.arg(delivery_charge),
--     sqlc.arg(total_payable),
--     sqlc.arg(shipping_address_id),
--     sqlc.arg(created_at)
-- )
-- RETURNING
--     id,
--     user_id,
--     subtotal,
--     total_weight_grams,
--     delivery_charge,
--     total_payable,
--     shipping_address_id,
--     created_at;


-- name: InsertCheckoutSession :one
INSERT INTO checkout_sessions (
    id,
    user_id,
    subtotal,
    total_weight_grams,
    delivery_charge,
    total_payable,
    shipping_address_id,
    status,
    platform_discount_type,
    platform_discount_value,
    platform_discount_amount_applied,
    is_platform_discount_applied,
    payment_method,
    created_at
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14
)
RETURNING id;


-- name: GetCheckoutSessionByID :one
SELECT
    id,
    user_id,
    subtotal,
    total_weight_grams,
    delivery_charge,
    total_payable,
    shipping_address_id,
    created_at,
    status,
    platform_discount_type,
    platform_discount_value,
    platform_discount_amount_applied,
    is_platform_discount_applied,
    payment_method
FROM checkout_sessions
WHERE id = $1;


-- name: UpdateCheckoutSessionWithShipping :exec
UPDATE checkout_sessions
SET
    shipping_address_id = $1,
    payment_method = $2,
    delivery_charge = $3,
    total_payable = $4
WHERE id = $5;


-- name: UpdateCheckoutSessionPaymentMethod :exec
UPDATE checkout_sessions
SET
    payment_method = $1
WHERE id = $2;

  -- name: UpdateCheckoutSessionDeliveryPricing :exec
UPDATE checkout_sessions
SET
    delivery_charge = $1,
    total_payable = $2
WHERE id = $3;

-- name: UpdateCheckoutSessionStatusToOrderCreated :exec
UPDATE checkout_sessions
SET
    status = $2
WHERE id = $1;

-- name: UpdateCheckoutSessionDeliveryPricing :exec
UPDATE checkout_sessions
SET
    delivery_charge = $1,
    total_payable = $2
WHERE id = $3;


-- name: MarkCheckoutSessionReadyToOrder :one
UPDATE checkout_sessions
SET status = 'ready_to_order'
WHERE id = $1
  AND user_id = $2
  AND status = 'awaiting_shipping_info'
RETURNING id;
