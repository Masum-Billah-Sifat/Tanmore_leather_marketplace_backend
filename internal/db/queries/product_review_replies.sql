-- name: InsertReviewReply :one
INSERT INTO product_review_replies (
    id,
    review_id,
    seller_user_id,
    reply_text,
    reply_image_url,
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

-- name: GetReviewReplyByReviewID :one
SELECT
    id,
    review_id,
    seller_user_id,
    reply_text,
    reply_image_url,
    is_edited,
    is_archived,
    is_banned,
    created_at,
    updated_at
FROM product_review_replies
WHERE review_id = $1;

-- name: UpdateReviewReplyText :exec
UPDATE product_review_replies
SET
    reply_text = $3,
    is_edited = $4,
    updated_at = $5
WHERE review_id = $1 AND seller_user_id = $2;

-- name: GetReviewReplyByReviewIDAndSellerID :one
SELECT
    id,
    review_id,
    seller_user_id,
    reply_text,
    reply_image_url,
    is_edited,
    is_archived,
    is_banned,
    created_at,
    updated_at
FROM product_review_replies
WHERE review_id = $1 AND seller_user_id = $2;

-- name: ArchiveReviewReply :exec
UPDATE product_review_replies
SET
    is_archived = $3,
    updated_at = $4
WHERE review_id = $1 AND seller_user_id = $2;

-- name: GetRepliesByReviewIDs :many
SELECT
    id,
    review_id,
    seller_user_id,
    reply_text,
    reply_image_url,
    is_edited,
    created_at,
    updated_at
FROM product_review_replies
WHERE review_id = ANY($1::uuid[])
  AND is_archived = FALSE
  AND is_banned = FALSE;