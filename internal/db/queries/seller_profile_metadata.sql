-- name: GetSellerProfileMetadataBySellerID :one
SELECT
    id,
    seller_id,

    sellerstorename,
    sellercontactno,
    sellerwhatsappcontactno,

    sellerwebsitelink,
    sellerfacebookpagename,
    selleremail,

    sellerphysicallocation,

    created_at,
    updated_at
FROM seller_profile_metadata
WHERE seller_id = $1;

-- name: InsertSellerProfileMetadata :one
INSERT INTO seller_profile_metadata (
  seller_id,
  sellerstorename,
  sellercontactno,
  sellerwhatsappcontactno,
  sellerwebsitelink,
  sellerfacebookpagename,
  selleremail,
  sellerphysicallocation,
  created_at,
  updated_at
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id;