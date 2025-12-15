-- Retail pricing
ALTER TABLE product_variant_snapshots
    ALTER COLUMN retailprice TYPE BIGINT USING retailprice::BIGINT,
    ALTER COLUMN retaildiscount TYPE BIGINT USING retaildiscount::BIGINT,
    ADD CONSTRAINT chk_variantsnapshot_retaildiscounttype CHECK (retaildiscounttype IN ('flat', 'percentage'));

-- Wholesale pricing
ALTER TABLE product_variant_snapshots
    ALTER COLUMN wholesaleprice TYPE BIGINT USING wholesaleprice::BIGINT,
    ALTER COLUMN wholesalediscount TYPE BIGINT USING wholesalediscount::BIGINT,
    ADD CONSTRAINT chk_variantsnapshot_wholesalediscounttype CHECK (wholesalediscounttype IN ('flat', 'percentage'));
