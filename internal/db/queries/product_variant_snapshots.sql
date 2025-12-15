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
