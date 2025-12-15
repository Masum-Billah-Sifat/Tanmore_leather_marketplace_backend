ALTER TABLE product_variant_indexes
    RENAME COLUMN has_retail_discount TO has_discount;

ALTER TABLE product_variant_indexes
    DROP CONSTRAINT IF EXISTS chk_variantindex_retaildiscounttype,
    DROP CONSTRAINT IF EXISTS chk_variantindex_wholesalediscounttype,
    ALTER COLUMN retailprice TYPE NUMERIC USING retailprice::NUMERIC,
    ALTER COLUMN retaildiscount TYPE NUMERIC USING retaildiscount::NUMERIC,
    ALTER COLUMN wholesaleprice TYPE NUMERIC USING wholesaleprice::NUMERIC,
    ALTER COLUMN wholesalediscount TYPE NUMERIC USING wholesalediscount::NUMERIC;
