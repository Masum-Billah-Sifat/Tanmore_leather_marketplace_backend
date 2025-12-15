CREATE TABLE product_medias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    media_type TEXT NOT NULL CHECK (media_type IN ('image', 'video')),
    media_url TEXT NOT NULL,

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Find media for a product fast
CREATE INDEX idx_product_medias_product_id ON product_medias(product_id);

-- Quick filtering for primary image
CREATE INDEX idx_product_medias_is_primary ON product_medias(is_primary);

-- Dealing with archival states
CREATE INDEX idx_product_medias_is_archived ON product_medias(is_archived);
