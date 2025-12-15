CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    color TEXT NOT NULL,
    size  TEXT NOT NULL,

    retail_price NUMERIC NOT NULL,

    retaildiscounttype TEXT,   -- optional
    retaildiscount     NUMERIC, -- optional

    wholesale_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    wholesale_price   NUMERIC, -- optional
    min_qty_wholesale INT,     -- optional
    wholesalediscounttype TEXT, -- optional
    wholesalediscount     NUMERIC, -- optional

    stock_quantity INT NOT NULL DEFAULT 0,
    in_stock BOOLEAN NOT NULL DEFAULT TRUE,

    weight_grams INT NOT NULL,

    ia_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    hasretaildiscount BOOLEAN NOT NULL DEFAULT FALSE,
    haswholesalediscount BOOLEAN NOT NULL DEFAULT FALSE
);


-- Most common lookup: find all variants for a product
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);

-- Inventory-related filtering
CREATE INDEX idx_product_variants_in_stock ON product_variants(in_stock);

-- Retail price filtering
CREATE INDEX idx_product_variants_retail_price ON product_variants(retail_price);

-- Weight for shipping calculations
CREATE INDEX idx_product_variants_weight ON product_variants(weight_grams);

-- Archival state for admin operations
CREATE INDEX idx_product_variants_ia_archived ON product_variants(ia_archived);

