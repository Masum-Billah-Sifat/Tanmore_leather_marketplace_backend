-- name: GetVariantSnapshot :one
SELECT
    id,

    -- Category info
    categoryid,
    iscategoryarchived,
    categoryname,

    -- Seller info
    sellerid,
    issellerapproved,
    issellerarchived,
    issellerbanned,
    sellerstorename,

    -- Product info
    productid,
    isproductapproved,
    isproductarchived,
    isproductbanned,
    producttitle,
    productdescription,
    productprimaryimageurl,

    -- Variant info
    variantid,
    isvariantarchived,
    isvariantinstock,
    stockamount,
    color,
    size,
    retailprice,

    -- Retail discount (optional)
    hasretaildiscount,
    retaildiscounttype,
    retaildiscount,

    -- Wholesale (optional)
    haswholesaleenabled,
    wholesaleprice,
    wholesaleminquantity,
    haswholesalediscount,
    wholesalediscounttype,
    wholesalediscount,

    -- Weight and timestamps
    weight_grams,
    createdat,
    updatedat

FROM product_variant_snapshots
WHERE sellerid = $1 AND productid = $2 AND variantid = $3;


-- name: GetVariantSnapshotByProductIDAndVariantID :one
SELECT
    id,

    -- Category info
    categoryid,
    iscategoryarchived,
    categoryname,

    -- Seller info
    sellerid,
    issellerapproved,
    issellerarchived,
    issellerbanned,
    sellerstorename,

    -- Product info
    productid,
    isproductapproved,
    isproductarchived,
    isproductbanned,
    producttitle,
    productdescription,
    productprimaryimageurl,

    -- Variant info
    variantid,
    isvariantarchived,
    isvariantinstock,
    stockamount,
    color,
    size,
    retailprice,

    -- Retail discount (optional)
    hasretaildiscount,
    retaildiscounttype,
    retaildiscount,

    -- Wholesale (optional)
    haswholesaleenabled,
    wholesaleprice,
    wholesaleminquantity,
    haswholesalediscount,
    wholesalediscounttype,
    wholesalediscount,

    -- Weight and timestamps
    weight_grams,
    createdat,
    updatedat

FROM product_variant_snapshots
WHERE productid = $1 AND variantid = $2;


-- name: GetVariantSnapshotByVariantID :one
SELECT
    id,

    -- Category info
    categoryid,
    iscategoryarchived,
    categoryname,

    -- Seller info
    sellerid,
    issellerapproved,
    issellerarchived,
    issellerbanned,
    sellerstorename,

    -- Product info
    productid,
    isproductapproved,
    isproductarchived,
    isproductbanned,
    producttitle,
    productdescription,
    productprimaryimageurl,

    -- Variant info
    variantid,
    isvariantarchived,
    isvariantinstock,
    stockamount,
    color,
    size,
    retailprice,

    -- Retail discount (optional)
    hasretaildiscount,
    retaildiscounttype,
    retaildiscount,

    -- Wholesale (optional)
    haswholesaleenabled,
    wholesaleprice,
    wholesaleminquantity,
    haswholesalediscount,
    wholesalediscounttype,
    wholesalediscount,

    -- Weight and timestamps
    weight_grams,
    createdat,
    updatedat

FROM product_variant_snapshots
WHERE variantid = $1;



-- name: InsertProductVariantSnapshot :exec
INSERT INTO product_variant_snapshots (
    id,
    categoryid,
    iscategoryarchived,
    categoryname,

    sellerid,
    issellerapproved,
    issellerarchived,
    issellerbanned,
    sellerstorename,

    productid,
    isproductapproved,
    isproductarchived,
    isproductbanned,
    producttitle,
    productdescription,
    productprimaryimageurl,

    variantid,
    isvariantarchived,
    isvariantinstock,
    stockamount,
    color,
    size,
    retailprice,

    hasretaildiscount,
    retaildiscounttype,
    retaildiscount,

    haswholesaleenabled,
    wholesaleprice,
    wholesaleminquantity,
    haswholesalediscount,
    wholesalediscounttype,
    wholesalediscount,

    weight_grams,
    createdat,
    updatedat
)
VALUES (
    $1,  -- id
    $2,  -- categoryid
    $3,  -- iscategoryarchived
    $4,  -- categoryname

    $5,  -- sellerid
    $6,  -- issellerapproved
    $7,  -- issellerarchived
    $8,  -- issellerbanned
    $9,  -- sellerstorename

    $10, -- productid
    $11, -- isproductapproved
    $12, -- isproductarchived
    $13, -- isproductbanned
    $14, -- producttitle
    $15, -- productdescription
    $16, -- productprimaryimageurl

    $17, -- variantid
    $18, -- isvariantarchived
    $19, -- isvariantinstock
    $20, -- stockamount
    $21, -- color
    $22, -- size
    $23, -- retailprice

    $24, -- hasretaildiscount
    $25, -- retaildiscounttype
    $26, -- retaildiscount

    $27, -- haswholesaleenabled
    $28, -- wholesaleprice
    $29, -- wholesaleminquantity
    $30, -- haswholesalediscount
    $31, -- wholesalediscounttype
    $32, -- wholesalediscount

    $33, -- weight_grams
    $34, -- createdat
    $35  -- updatedat
);




-- name: UpdateSnapshotsOnProductInfoUpdate :exec
UPDATE product_variant_snapshots
SET
    producttitle = COALESCE(sqlc.narg(producttitle)::TEXT, producttitle),
    productdescription = COALESCE(sqlc.narg(productdescription)::TEXT, productdescription),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid);



-- name: UpdateSnapshotPrimaryImageURL :exec
UPDATE product_variant_snapshots
SET
    productprimaryimageurl = sqlc.arg(media_url),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid);


-- name: UpdateSnapshotOnInStockUpdate :exec
UPDATE product_variant_snapshots
SET
    isvariantinstock = sqlc.arg(in_stock),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateSnapshotOnRetailPriceUpdate :exec
UPDATE product_variant_snapshots
SET
    retailprice = sqlc.arg(retail_price),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateSnapshotOnStockQuantityUpdate :exec
UPDATE product_variant_snapshots
SET
    stockamount = sqlc.arg(stock_quantity),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: UpdateSnapshotOnWeightUpdate :exec
UPDATE product_variant_snapshots
SET
    weight_grams = sqlc.arg(weight_grams),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);

-- name: ArchiveVariantInSnapshots :exec
UPDATE product_variant_snapshots
SET
    isvariantarchived = sqlc.arg(isvariantarchived),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateSnapshotOnVariantInfoUpdate :exec
UPDATE product_variant_snapshots
SET
    color = COALESCE(sqlc.narg(color)::TEXT, color),
    size = COALESCE(sqlc.narg(size)::TEXT, size),
    updatedat = sqlc.arg(updatedat)
WHERE variantid = sqlc.arg(variantid)
  AND productid = sqlc.arg(productid);


-- name: ApplyRetailDiscountToSnapshot :exec
UPDATE product_variant_snapshots
SET
    hasretaildiscount = sqlc.arg(has_retail_discount),
    retaildiscount = sqlc.arg(retail_discount),
    retaildiscounttype = sqlc.arg(retail_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateRetailDiscountInSnapshot :exec
UPDATE product_variant_snapshots
SET
    retaildiscount = COALESCE(sqlc.narg(retail_discount)::BIGINT, retaildiscount),
    retaildiscounttype = COALESCE(sqlc.narg(retail_discount_type)::TEXT, retaildiscounttype),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: RemoveRetailDiscountFromSnapshot :exec
UPDATE product_variant_snapshots
SET
    hasretaildiscount = sqlc.arg(has_retail_discount),
    retaildiscount = sqlc.arg(retail_discount),
    retaildiscounttype = sqlc.arg(retail_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: EnableWholesaleModeInSnapshots :exec
UPDATE product_variant_snapshots
SET
    wholesaleprice = sqlc.arg(wholesale_price),
    wholesaleminquantity = sqlc.arg(min_qty_wholesale),
    haswholesalediscount = sqlc.arg(has_wholesale_discount),
    wholesalediscount = COALESCE(sqlc.narg(wholesale_discount)::BIGINT, NULL),
    wholesalediscounttype = COALESCE(sqlc.narg(wholesale_discount_type)::TEXT, NULL),
    haswholesaleenabled = sqlc.arg(has_wholesale_enabled),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateWholesaleModeInSnapshots :exec
UPDATE product_variant_snapshots
SET
    wholesaleprice = COALESCE(sqlc.narg(wholesale_price)::BIGINT, wholesaleprice),
    wholesaleminquantity = COALESCE(sqlc.narg(min_qty_wholesale)::INT, wholesaleminquantity),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: DisableWholesaleModeInSnapshots :exec
UPDATE product_variant_snapshots
SET
    haswholesaleenabled = sqlc.arg(has_wholesale_enabled),
    haswholesalediscount = sqlc.arg(has_wholesale_discount),
    wholesaleprice = sqlc.arg(wholesale_price),
    wholesaleminquantity = sqlc.arg(wholesale_min_quantity),
    wholesalediscount = sqlc.arg(wholesale_discount),
    wholesalediscounttype = sqlc.arg(wholesale_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);

-- name: ApplyWholesaleDiscountToSnapshot :exec
UPDATE product_variant_snapshots
SET
    haswholesalediscount = sqlc.arg(has_wholesale_discount),
    wholesalediscount = sqlc.arg(wholesale_discount),
    wholesalediscounttype = sqlc.arg(wholesale_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateWholesaleDiscountInSnapshot :exec
UPDATE product_variant_snapshots
SET
    wholesalediscount = COALESCE(sqlc.narg(wholesale_discount)::BIGINT, wholesalediscount),
    wholesalediscounttype = COALESCE(sqlc.narg(wholesale_discount_type)::TEXT, wholesalediscounttype),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: RemoveWholesaleDiscountFromSnapshot :exec
UPDATE product_variant_snapshots
SET
    haswholesalediscount = sqlc.arg(has_wholesale_discount),
    wholesalediscount = sqlc.arg(wholesale_discount),
    wholesalediscounttype = sqlc.arg(wholesale_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);

