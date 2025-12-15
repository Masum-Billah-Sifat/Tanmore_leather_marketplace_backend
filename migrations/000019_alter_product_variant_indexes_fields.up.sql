-- Retail pricing
ALTER TABLE product_variant_indexes
    ALTER COLUMN retailprice TYPE BIGINT USING retailprice::BIGINT,
    ALTER COLUMN retaildiscount TYPE BIGINT USING retaildiscount::BIGINT,
    ADD CONSTRAINT chk_variantindex_retaildiscounttype CHECK (retaildiscounttype IN ('flat', 'percentage'));

-- Wholesale pricing
ALTER TABLE product_variant_indexes
    ALTER COLUMN wholesaleprice TYPE BIGINT USING wholesaleprice::BIGINT,
    ALTER COLUMN wholesalediscount TYPE BIGINT USING wholesalediscount::BIGINT,
    ADD CONSTRAINT chk_variantindex_wholesalediscounttype CHECK (wholesalediscounttype IN ('flat', 'percentage'));

-- Rename for consistency
ALTER TABLE product_variant_indexes
    RENAME COLUMN has_discount TO has_retail_discount;
