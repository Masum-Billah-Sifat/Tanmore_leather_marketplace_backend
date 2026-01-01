-- name: GetCategoryByID :one
SELECT
  id,
  parent_id,
  name,
  slug,
  level,
  is_leaf,
  is_archived,
  created_at
FROM categories
WHERE id = $1;

-- name: GetAllNonArchivedCategories :many
SELECT
    id,
    parent_id,
    name,
    slug,
    level,
    is_leaf
FROM categories
WHERE is_archived = FALSE
ORDER BY level ASC, name ASC;


-- name: GetAllLeafCategoryIDsByRoot :many
WITH RECURSIVE subcategories (id, is_leaf) AS (
    SELECT
        c.id,
        c.is_leaf
    FROM categories c
    WHERE c.id = $1

    UNION ALL

    SELECT
        c.id,
        c.is_leaf
    FROM categories c
    INNER JOIN subcategories sc ON c.parent_id = sc.id
)
SELECT
    sc.id
FROM subcategories sc
WHERE sc.is_leaf = TRUE;



-- name: GetProductVariantIndexesByCategoryIDs :many
SELECT
    id, -- index id (not needed in response)

    -- Category
    categoryid,
    categoryname,
    iscategoryarchived,

    -- Seller
    sellerid,
    sellerstorename,
    issellerapproved,
    issellerarchived,
    issellerbanned,

    -- Product
    productid,
    producttitle,
    productdescription,
    productimages,
    productpromovideourl,
    isproductapproved,
    isproductarchived,
    isproductbanned,

    -- Variant
    variantid,
    color,
    size,
    isvariantinstock,
    stockamount,
    retailprice,
    retaildiscount,
    retaildiscounttype,
    has_retail_discount,
    haswholesaleenabled,
    wholesaleprice,
    wholesaleminquantity,
    wholesalediscount,
    wholesalediscounttype,
    isvariantarchived,
    weight_grams,

    -- Metadata
    views,
    createdat,
    updatedat

FROM product_variant_indexes
WHERE
    categoryid = ANY($1::uuid[])
    AND iscategoryarchived = FALSE
    AND isproductarchived = FALSE
    AND isproductapproved = TRUE
    AND isproductbanned = FALSE
    AND isvariantarchived = FALSE
    AND issellerarchived = FALSE
    AND issellerapproved = TRUE
    AND issellerbanned = FALSE
ORDER BY producttitle, variantid;