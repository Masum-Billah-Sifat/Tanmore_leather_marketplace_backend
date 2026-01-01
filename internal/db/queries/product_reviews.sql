-- name: InsertProductReview :one
INSERT INTO product_reviews (
    id,
    product_id,
    reviewer_user_id,
    review_text,
    review_image_url,
    is_edited,
    is_archived,
    is_banned,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10
)
RETURNING id;

-- name: UpdateProductReviewText :exec
UPDATE product_reviews
SET
    review_text = $4,
    is_edited = $5,
    updated_at = $6
WHERE id = $1 AND product_id = $2 AND reviewer_user_id = $3;



-- name: GetProductReviewByIDAndProductIDAndReviewerID :one
SELECT
    id,
    product_id,
    reviewer_user_id,
    review_text,
    review_image_url,
    is_edited,
    is_archived,
    is_banned,
    created_at,
    updated_at
FROM product_reviews
WHERE id = $1 AND product_id = $2 AND reviewer_user_id = $3;

-- name: ArchiveProductReview :exec
UPDATE product_reviews
SET
    is_archived = $2,
    updated_at = $3
WHERE id = $1;

-- name: GetProductReviewByID :one
SELECT
    id,
    product_id,
    reviewer_user_id,
    review_text,
    review_image_url,
    is_edited,
    is_archived,
    is_banned,
    created_at,
    updated_at
FROM product_reviews
WHERE id = $1;

-- name: GetAllReviewsByProductID :many
SELECT
    id,
    product_id,
    reviewer_user_id,
    review_text,
    review_image_url,
    is_edited,
    created_at,
    updated_at
FROM product_reviews
WHERE product_id = $1
  AND is_archived = FALSE
  AND is_banned = FALSE
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
