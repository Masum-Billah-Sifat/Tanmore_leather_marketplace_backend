CREATE TABLE product_variant_indexes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Category info
    categoryid UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    iscategoryarchived BOOLEAN NOT NULL DEFAULT FALSE,
    categoryname TEXT NOT NULL,

    -- Seller info
    sellerid UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    issellerapproved BOOLEAN NOT NULL DEFAULT FALSE,
    issellerarchived BOOLEAN NOT NULL DEFAULT FALSE,
    issellerbanned BOOLEAN NOT NULL DEFAULT FALSE,
    sellerstorename TEXT NOT NULL,

    -- Product info
    productid UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    isproductapproved BOOLEAN NOT NULL DEFAULT FALSE,
    isproductarchived BOOLEAN NOT NULL DEFAULT FALSE,
    isproductbanned BOOLEAN NOT NULL DEFAULT FALSE,
    producttitle TEXT NOT NULL,
    productdescription TEXT NOT NULL,

    productimages TEXT[] NOT NULL,
    productpromovideourl TEXT,

    -- Variant info
    variantid UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    isvariantarchived BOOLEAN NOT NULL DEFAULT FALSE,
    isvariantinstock BOOLEAN NOT NULL DEFAULT TRUE,
    stockamount INT NOT NULL DEFAULT 0,
    color TEXT NOT NULL,
    size TEXT NOT NULL,
    retailprice NUMERIC NOT NULL,

    retaildiscounttype TEXT,
    retaildiscount NUMERIC,
    has_discount BOOLEAN NOT NULL DEFAULT FALSE,

    haswholesaleenabled BOOLEAN NOT NULL DEFAULT FALSE,
    wholesaleprice NUMERIC,
    wholesaleminquantity INT,
    wholesalediscounttype TEXT,
    wholesalediscount NUMERIC,

    weight_grams INT NOT NULL,

    search_vector tsvector NOT NULL,
    views BIGINT NOT NULL DEFAULT 0,

    createdat TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updatedat TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_indexes_categoryid ON product_variant_indexes(categoryid);
CREATE INDEX idx_indexes_sellerid ON product_variant_indexes(sellerid);
CREATE INDEX idx_indexes_productid ON product_variant_indexes(productid);
CREATE INDEX idx_indexes_variantid ON product_variant_indexes(variantid);

CREATE INDEX idx_indexes_isvariantinstock ON product_variant_indexes(isvariantinstock);
CREATE INDEX idx_indexes_price ON product_variant_indexes(retailprice);

CREATE INDEX idx_indexes_has_discount ON product_variant_indexes(has_discount);

CREATE INDEX idx_indexes_search_vector ON product_variant_indexes USING GIN (search_vector);

CREATE INDEX idx_indexes_views ON product_variant_indexes(views);
