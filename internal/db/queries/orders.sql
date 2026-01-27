-- name: InsertOrderRow :one
INSERT INTO orders (
    id,
    user_id,
    checkout_session_id,
    order_code,

    subtotal,
    shipping_fee,
    total_amount,
    currency,

    payment_method,
    payment_status,
    latest_payment_id,

    status,

    platform_discount_type,
    platform_discount_value,
    platform_discount_amount_applied,
    is_platform_discount_applied,

    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8,
    $9, $10, $11,
    $12,
    $13, $14, $15, $16,
    $17, $18
)
RETURNING
    id,
    user_id,
    checkout_session_id,
    order_code,
    subtotal,
    shipping_fee,
    total_amount,
    currency,
    payment_method,
    payment_status,
    latest_payment_id,
    status,
    platform_discount_type,
    platform_discount_value,
    platform_discount_amount_applied,
    is_platform_discount_applied,
    created_at,
    updated_at;