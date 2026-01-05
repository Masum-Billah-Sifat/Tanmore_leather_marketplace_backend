-- name: InsertProductVariantIndex :exec
INSERT INTO product_variant_indexes (
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
    productimages,
    productpromovideourl,

    variantid,
    isvariantarchived,
    isvariantinstock,
    stockamount,
    color,
    size,
    retailprice,

    retaildiscounttype,
    retaildiscount,
    has_retail_discount,

    haswholesaleenabled,
    wholesaleprice,
    wholesaleminquantity,
    wholesalediscounttype,
    wholesalediscount,

    weight_grams,
    search_vector,
    views,
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
    $16, -- productimages
    $17, -- productpromovideourl

    $18, -- variantid
    $19, -- isvariantarchived
    $20, -- isvariantinstock
    $21, -- stockamount
    $22, -- color
    $23, -- size
    $24, -- retailprice

    $25, -- retaildiscounttype
    $26, -- retaildiscount
    $27, -- has_retail_discount

    $28, -- haswholesaleenabled
    $29, -- wholesaleprice
    $30, -- wholesaleminquantity
    $31, -- wholesalediscounttype
    $32, -- wholesalediscount

    $33, -- weight_grams
    to_tsvector('english', $34), -- search_vector
    $35, -- views
    $36, -- createdat
    $37  -- updatedat
);

-- name: UpdateIndexesOnProductInfoUpdate :exec
UPDATE product_variant_indexes
SET
    producttitle = COALESCE(sqlc.narg(producttitle)::TEXT, producttitle),
    productdescription = COALESCE(sqlc.narg(productdescription)::TEXT, productdescription),
    search_vector = to_tsvector('english', COALESCE(sqlc.narg(producttitle)::TEXT, producttitle) || ' ' || COALESCE(sqlc.narg(productdescription)::TEXT, productdescription)),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid);


-- name: UpdateIndexOnProductImageAdded :exec
UPDATE product_variant_indexes
SET
    productimages = array_append(productimages, sqlc.arg(media_url)),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid);


-- name: UpdateIndexOnPromoVideoAdded :exec
UPDATE product_variant_indexes
SET
    productpromovideourl = sqlc.arg(media_url),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid);


-- name: UpdateIndexOnImageRemoved :exec
UPDATE product_variant_indexes
SET
    productimages = array_remove(productimages, sqlc.arg(media_url)),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid);


-- name: UpdateIndexOnPromoVideoRemoved :exec
UPDATE product_variant_indexes
SET
    productpromovideourl = sqlc.arg(productpromovideourl),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid);


-- name: UpdateIndexOnInStockUpdate :exec
UPDATE product_variant_indexes
SET
    isvariantinstock = sqlc.arg(in_stock),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: UpdateIndexOnRetailPriceUpdate :exec
UPDATE product_variant_indexes
SET
    retailprice = sqlc.arg(retail_price),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateIndexOnStockQuantityUpdate :exec
UPDATE product_variant_indexes
SET
    stockamount = sqlc.arg(stock_quantity),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateIndexOnWeightUpdate :exec
UPDATE product_variant_indexes
SET
    weight_grams = sqlc.arg(weight_grams),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);

-- name: ArchiveVariantInIndexes :exec
UPDATE product_variant_indexes
SET
    isvariantarchived = sqlc.arg(isvariantarchived),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateIndexOnVariantInfoUpdate :exec
UPDATE product_variant_indexes
SET
    color = COALESCE(sqlc.narg(color)::TEXT, color),
    size = COALESCE(sqlc.narg(size)::TEXT, size),
    updatedat = sqlc.arg(updatedat)
WHERE variantid = sqlc.arg(variantid)
  AND productid = sqlc.arg(productid);


-- name: ApplyRetailDiscountToIndex :exec
UPDATE product_variant_indexes
SET
    has_retail_discount = sqlc.arg(has_retail_discount),
    retaildiscount = sqlc.arg(retail_discount),
    retaildiscounttype = sqlc.arg(retail_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: UpdateRetailDiscountInIndex :exec
UPDATE product_variant_indexes
SET
    retaildiscount = COALESCE(sqlc.narg(retail_discount)::BIGINT, retaildiscount),
    retaildiscounttype = COALESCE(sqlc.narg(retail_discount_type)::TEXT, retaildiscounttype),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: RemoveRetailDiscountFromIndex :exec
UPDATE product_variant_indexes
SET
    has_retail_discount = sqlc.arg(has_retail_discount),
    retaildiscount = sqlc.arg(retail_discount),
    retaildiscounttype = sqlc.arg(retail_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: EnableWholesaleModeInIndexes :exec
UPDATE product_variant_indexes
SET
    wholesaleprice = sqlc.arg(wholesale_price),
    wholesaleminquantity = sqlc.arg(min_qty_wholesale),
    wholesalediscount = COALESCE(sqlc.narg(wholesale_discount)::BIGINT, NULL),
    wholesalediscounttype = COALESCE(sqlc.narg(wholesale_discount_type)::TEXT, NULL),
    haswholesaleenabled = sqlc.arg(has_wholesale_enabled),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: UpdateWholesaleModeInIndexes :exec
UPDATE product_variant_indexes
SET
    wholesaleprice = COALESCE(sqlc.narg(wholesale_price)::BIGINT, wholesaleprice),
    wholesaleminquantity = COALESCE(sqlc.narg(min_qty_wholesale)::INT, wholesaleminquantity),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: DisableWholesaleModeInIndexes :exec
UPDATE product_variant_indexes
SET
    haswholesaleenabled = sqlc.arg(has_wholesale_enabled),
    wholesaleprice = sqlc.arg(wholesale_price),
    wholesaleminquantity = sqlc.arg(wholesale_min_quantity),
    wholesalediscount = sqlc.arg(wholesale_discount),
    wholesalediscounttype = sqlc.arg(wholesale_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: ApplyWholesaleDiscountToIndex :exec
UPDATE product_variant_indexes
SET
    wholesalediscount = sqlc.arg(wholesale_discount),
    wholesalediscounttype = sqlc.arg(wholesale_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: UpdateWholesaleDiscountInIndex :exec
UPDATE product_variant_indexes
SET
    wholesalediscount = COALESCE(sqlc.narg(wholesale_discount)::BIGINT, wholesalediscount),
    wholesalediscounttype = COALESCE(sqlc.narg(wholesale_discount_type)::TEXT, wholesalediscounttype),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);



-- name: RemoveWholesaleDiscountFromIndex :exec
UPDATE product_variant_indexes
SET
    wholesalediscount = sqlc.arg(wholesale_discount),
    wholesalediscounttype = sqlc.arg(wholesale_discount_type),
    updatedat = sqlc.arg(updatedat)
WHERE productid = sqlc.arg(productid)
  AND variantid = sqlc.arg(variantid);


-- name: GetVariantIndexesByProductAndSeller :many
SELECT
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
  productimages,
  productpromovideourl,
  variantid,
  isvariantarchived,
  isvariantinstock,
  stockamount,
  color,
  size,
  retailprice,
  retaildiscounttype,
  retaildiscount,
  has_retail_discount,
  haswholesaleenabled,
  wholesaleprice,
  wholesaleminquantity,
  wholesalediscounttype,
  wholesalediscount,
  weight_grams,
  search_vector,
  views,
  createdat,
  updatedat
FROM product_variant_indexes
WHERE productid = $1
  AND sellerid = $2;


-- name: GetAllProductVariantIndexesBySeller :many
SELECT
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
  productimages,
  productpromovideourl,
  variantid,
  isvariantarchived,
  isvariantinstock,
  stockamount,
  color,
  size,
  retailprice,
  retaildiscounttype,
  retaildiscount,
  has_retail_discount,
  haswholesaleenabled,
  wholesaleprice,
  wholesaleminquantity,
  wholesalediscounttype,
  wholesalediscount,
  weight_grams,
  search_vector,
  views,
  createdat,
  updatedat
FROM product_variant_indexes
WHERE sellerid = $1
  AND issellerapproved = true
  AND issellerarchived = false
  AND issellerbanned = false;



-- name: UpdateCategoryInProductVariantIndexes :exec
UPDATE product_variant_indexes
SET
    categoryid = $1,
    categoryname = $2,
    updatedat = $3
WHERE productid = $4;


-- name: MarkProductArchivedInProductVariantIndexes :exec
UPDATE product_variant_indexes
SET
    isproductarchived = true,
    updatedat = $1
WHERE productid = $2;
