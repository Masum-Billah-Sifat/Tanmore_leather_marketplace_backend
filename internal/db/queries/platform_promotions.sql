-- name: GetActivePlatformPromotions :many
SELECT
    id,
    title,
    description,
    discount_type,
    discount_value,
    min_cart_value,
    max_cart_value,
    max_discount_cap,
    start_time,
    end_time,
    priority
FROM platform_promotions
WHERE is_active = true
  AND is_archived = false;
