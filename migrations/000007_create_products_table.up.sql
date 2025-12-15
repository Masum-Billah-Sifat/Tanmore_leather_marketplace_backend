CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,

    title TEXT NOT NULL,
    description TEXT NOT NULL,

    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned   BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    search_vector tsvector
);


-- Search performance
CREATE INDEX idx_products_search_vector ON products USING GIN (search_vector);

-- Filter by seller
CREATE INDEX idx_products_seller_id ON products(seller_id);

-- Filter by category
CREATE INDEX idx_products_category_id ON products(category_id);

-- Fetch active/approved products quickly
CREATE INDEX idx_products_is_approved ON products(is_approved);
CREATE INDEX idx_products_is_archived ON products(is_archived);
CREATE INDEX idx_products_is_banned   ON products(is_banned);
