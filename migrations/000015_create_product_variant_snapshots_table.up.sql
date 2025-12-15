CREATE TABLE product_variant_snapshots (
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
    productprimaryimageurl TEXT NOT NULL,

    -- Variant info
    variantid UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    isvariantarchived BOOLEAN NOT NULL DEFAULT FALSE,
    isvariantinstock BOOLEAN NOT NULL DEFAULT TRUE,
    stockamount INT NOT NULL DEFAULT 0,
    color TEXT NOT NULL,
    size TEXT NOT NULL,
    retailprice NUMERIC NOT NULL,

    -- Retail discount (optional)
    hasretaildiscount BOOLEAN NOT NULL DEFAULT FALSE,
    retaildiscounttype TEXT,
    retaildiscount NUMERIC,

    -- Wholesale (optional)
    haswholesaleenabled BOOLEAN NOT NULL DEFAULT FALSE,
    wholesaleprice NUMERIC,
    wholesaleminquantity INT,
    haswholesalediscount BOOLEAN NOT NULL DEFAULT FALSE,
    wholesalediscounttype TEXT,
    wholesalediscount NUMERIC,

    weight_grams INT NOT NULL,

    createdat TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updatedat TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



CREATE INDEX idx_snapshots_categoryid ON product_variant_snapshots(categoryid);
CREATE INDEX idx_snapshots_sellerid ON product_variant_snapshots(sellerid);
CREATE INDEX idx_snapshots_productid ON product_variant_snapshots(productid);
CREATE INDEX idx_snapshots_variantid ON product_variant_snapshots(variantid);

CREATE INDEX idx_snapshots_isvariantinstock ON product_variant_snapshots(isvariantinstock);
CREATE INDEX idx_snapshots_retailprice ON product_variant_snapshots(retailprice);
CREATE INDEX idx_snapshots_hasretaildiscount ON product_variant_snapshots(hasretaildiscount);
