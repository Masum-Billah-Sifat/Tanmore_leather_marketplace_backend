ALTER TABLE product_variant_snapshots
    DROP CONSTRAINT IF EXISTS chk_variantsnapshot_retaildiscounttype,
    DROP CONSTRAINT IF EXISTS chk_variantsnapshot_wholesalediscounttype,
    ALTER COLUMN retailprice TYPE NUMERIC USING retailprice::NUMERIC,
    ALTER COLUMN retaildiscount TYPE NUMERIC USING retaildiscount::NUMERIC,
    ALTER COLUMN wholesaleprice TYPE NUMERIC USING wholesaleprice::NUMERIC,
    ALTER COLUMN wholesalediscount TYPE NUMERIC USING wholesalediscount::NUMERIC;
