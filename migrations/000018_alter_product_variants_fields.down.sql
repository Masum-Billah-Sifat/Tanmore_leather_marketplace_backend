ALTER TABLE product_variants
    RENAME COLUMN is_archived TO ia_archived;

ALTER TABLE product_variants
    DROP CONSTRAINT IF EXISTS chk_variant_retaildiscounttype,
    DROP CONSTRAINT IF EXISTS chk_variant_wholesalediscounttype,
    ALTER COLUMN retail_price TYPE NUMERIC USING retail_price::NUMERIC,
    ALTER COLUMN retaildiscount TYPE NUMERIC USING retaildiscount::NUMERIC,
    ALTER COLUMN wholesale_price TYPE NUMERIC USING wholesale_price::NUMERIC,
    ALTER COLUMN wholesalediscount TYPE NUMERIC USING wholesalediscount::NUMERIC;
