CREATE TABLE checkout_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    checkout_session_id UUID NOT NULL 
        REFERENCES checkout_sessions(id) ON DELETE CASCADE,

    user_id UUID NOT NULL 
        REFERENCES users(id) ON DELETE CASCADE,

    seller_id UUID NOT NULL 
        REFERENCES users(id) ON DELETE CASCADE,

    category_id UUID NOT NULL 
        REFERENCES categories(id) ON DELETE RESTRICT,

    category_name TEXT NOT NULL,

    product_id UUID NOT NULL 
        REFERENCES products(id) ON DELETE RESTRICT,

    product_title TEXT NOT NULL,
    product_description TEXT NOT NULL,
    product_primary_image_url TEXT NOT NULL,

    variant_id UUID NOT NULL 
        REFERENCES product_variants(id) ON DELETE RESTRICT,

    color TEXT NOT NULL,
    size TEXT NOT NULL,

    buying_mode TEXT NOT NULL 
        CHECK (buying_mode IN ('retail', 'wholesale')),

    unit_price DECIMAL(10,2) NOT NULL,

    has_discount BOOLEAN NOT NULL,
    discount_type TEXT NOT NULL,
    discount_value DECIMAL(10,2) NOT NULL,

    required_quantity INTEGER NOT NULL,
    weight_grams INTEGER NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_checkout_items_session_id ON checkout_items(checkout_session_id);
CREATE INDEX idx_checkout_items_user_id ON checkout_items(user_id);
CREATE INDEX idx_checkout_items_seller_id ON checkout_items(seller_id);

CREATE INDEX idx_checkout_items_category_id ON checkout_items(category_id);
CREATE INDEX idx_checkout_items_product_id ON checkout_items(product_id);
CREATE INDEX idx_checkout_items_variant_id ON checkout_items(variant_id);

-- for sorting / analytics
CREATE INDEX idx_checkout_items_buying_mode ON checkout_items(buying_mode);
CREATE INDEX idx_checkout_items_created_at ON checkout_items(created_at);

