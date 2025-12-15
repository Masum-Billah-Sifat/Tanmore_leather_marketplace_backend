CREATE TABLE shipping_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    recipient_name TEXT NOT NULL,
    recipient_phone TEXT NOT NULL,

    recipient_email TEXT,        -- optional

    address_line TEXT NOT NULL,

    delivery_note TEXT,          -- optional

    city_id INT NOT NULL,
    zone_id INT NOT NULL,
    area_id INT NOT NULL,

    latitude DECIMAL(10,7),      -- optional
    longitude DECIMAL(10,7),     -- optional

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipping_addresses_city_zone_area ON shipping_addresses(city_id, zone_id, area_id);