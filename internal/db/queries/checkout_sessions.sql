-- name: InsertCheckoutSession :one
INSERT INTO checkout_sessions (
    id,
    user_id,
    subtotal,
    total_weight_grams,
    delivery_charge,
    total_payable,
    shipping_address_id,
    created_at
)
VALUES (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(subtotal),
    sqlc.arg(total_weight_grams),
    sqlc.arg(delivery_charge),
    sqlc.arg(total_payable),
    sqlc.arg(shipping_address_id),
    sqlc.arg(created_at)
)
RETURNING
    id,
    user_id,
    subtotal,
    total_weight_grams,
    delivery_charge,
    total_payable,
    shipping_address_id,
    created_at;
