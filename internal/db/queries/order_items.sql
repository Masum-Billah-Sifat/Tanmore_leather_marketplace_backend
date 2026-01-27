-- name: InsertOrderItem :exec
INSERT INTO order_items (
    id,
    order_id,

    customer_id,
    seller_id,

    category_id,
    category_name,

    product_id,
    product_title,
    product_description,
    product_primary_image_url,

    variant_id,
    color,
    size,

    buying_mode,
    unit_price,
    quantity,
    total_price,

    has_discount,
    discount_type,
    discount_value,

    weight_grams_per_unit,
    total_weight_grams,

    seller_store_name,
    created_at
) VALUES (
    $1, $2,
    $3, $4,
    $5, $6,
    $7, $8, $9, $10,
    $11, $12, $13,
    $14, $15, $16, $17,
    $18, $19, $20,
    $21, $22,
    $23, $24
);