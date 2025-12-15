-- Retail pricing
ALTER TABLE product_variants
    ALTER COLUMN retail_price TYPE BIGINT USING retail_price::BIGINT,
    ALTER COLUMN retaildiscount TYPE BIGINT USING retaildiscount::BIGINT,
    ADD CONSTRAINT chk_variant_retaildiscounttype CHECK (retaildiscounttype IN ('flat', 'percentage'));

-- Wholesale pricing
ALTER TABLE product_variants
    ALTER COLUMN wholesale_price TYPE BIGINT USING wholesale_price::BIGINT,
    ALTER COLUMN wholesalediscount TYPE BIGINT USING wholesalediscount::BIGINT,
    ADD CONSTRAINT chk_variant_wholesalediscounttype CHECK (wholesalediscounttype IN ('flat', 'percentage'));

-- Rename typo + consistency
ALTER TABLE product_variants
    RENAME COLUMN ia_archived TO is_archived;
