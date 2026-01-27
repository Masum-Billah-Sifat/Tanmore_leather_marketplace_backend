CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- 📦 Parent order reference
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

    -- 👥 Identifiers
    customer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    category_id UUID NOT NULL REFERENCES categories(id),
    category_name TEXT NOT NULL,

    product_id UUID NOT NULL REFERENCES products(id),
    product_title TEXT NOT NULL,
    product_description TEXT NOT NULL,
    product_primary_image_url TEXT NOT NULL,

    variant_id UUID NOT NULL REFERENCES product_variants(id),
    color TEXT NOT NULL,
    size TEXT NOT NULL,

    -- 💰 Pricing & Quantity
    buying_mode TEXT NOT NULL CHECK (buying_mode IN ('retail', 'wholesale')),
    unit_price DECIMAL(10,2) NOT NULL,
    quantity INTEGER NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,

    -- 🎁 Discount
    has_discount BOOLEAN NOT NULL,
    discount_type TEXT,
    discount_value DECIMAL(10,2),

    -- ⚖️ Weight
    weight_grams_per_unit INTEGER NOT NULL,
    total_weight_grams INTEGER NOT NULL,

    -- 🏬 Seller Info
    seller_store_name TEXT NOT NULL,

    -- 🕒 Timestamp
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- ✅ Constraint
    CHECK (
        (has_discount = false AND discount_type IS NULL AND discount_value IS NULL)
        OR
        (has_discount = true AND discount_type IS NOT NULL AND discount_value IS NOT NULL)
    )
);

-- 🔍 Indexes
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_seller_id ON order_items(seller_id);
CREATE INDEX idx_order_items_user_id ON order_items(customer_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_items_variant_id ON order_items(variant_id);
