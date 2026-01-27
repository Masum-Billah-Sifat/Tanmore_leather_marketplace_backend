CREATE TABLE shipping_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ✅ Linked to checkout session
    checkout_session_id UUID NOT NULL REFERENCES checkout_sessions(id) ON DELETE CASCADE,

    -- 🧍 Recipient Info
    recipient_name TEXT NOT NULL,
    recipient_phone TEXT NOT NULL CHECK (char_length(recipient_phone) = 11),
    recipient_email TEXT,

    -- 📍 Address Details
    address_line TEXT NOT NULL,
    delivery_note TEXT,

    -- 🌐 Pathao Integration Fields
    city_id INT NOT NULL,
    zone_id INT NOT NULL,
    area_id INT NOT NULL,

    -- 📌 Geo-coordinates (optional)
    latitude DECIMAL(10, 7),
    longitude DECIMAL(10, 7),

    created_at TIMESTAMP DEFAULT NOW()
);

-- 📦 Lookup by session
CREATE INDEX idx_shipping_addresses_checkout_session_id ON shipping_addresses(checkout_session_id);

-- 🚚 Optimized delivery zone queries
CREATE INDEX idx_shipping_addresses_location ON shipping_addresses(city_id, zone_id, area_id);

-- 🔍 Phone lookup for fraud check / quick search
CREATE INDEX idx_shipping_addresses_recipient_phone ON shipping_addresses(recipient_phone);
