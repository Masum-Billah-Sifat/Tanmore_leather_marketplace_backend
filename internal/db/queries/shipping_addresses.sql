-- name: GetShippingAddressByIDAndCheckoutID :one
SELECT
    id,
    checkout_session_id,
    recipient_name,
    recipient_phone,
    recipient_email,
    address_line,
    delivery_note,
    city_id,
    zone_id,
    area_id,
    latitude,
    longitude,
    created_at
FROM shipping_addresses
WHERE id = $1
  AND checkout_session_id = $2;


-- name: InsertShippingAddress :one
INSERT INTO shipping_addresses (
    id,
    checkout_session_id,
    recipient_name,
    recipient_phone,
    recipient_email,
    address_line,
    delivery_note,
    city_id,
    zone_id,
    area_id,
    latitude,
    longitude,
    created_at
)
VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING id;


-- name: UpdateShippingAddressByID :exec
UPDATE shipping_addresses
SET
    recipient_name   = COALESCE(sqlc.narg(recipient_name), recipient_name),
    recipient_phone  = COALESCE(sqlc.narg(recipient_phone), recipient_phone),
    recipient_email  = COALESCE(sqlc.narg(recipient_email), recipient_email),
    address_line     = COALESCE(sqlc.narg(address_line), address_line),
    delivery_note    = COALESCE(sqlc.narg(delivery_note), delivery_note),
    city_id          = COALESCE(sqlc.narg(city_id), city_id),
    zone_id          = COALESCE(sqlc.narg(zone_id), zone_id),
    area_id          = COALESCE(sqlc.narg(area_id), area_id),
    latitude         = COALESCE(sqlc.narg(latitude), latitude),
    longitude        = COALESCE(sqlc.narg(longitude), longitude)
WHERE id = sqlc.arg(id);