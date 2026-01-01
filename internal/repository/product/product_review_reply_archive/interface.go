// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_archive_reply/interface.go
// 🧠 Repository interface for archiving a seller's reply to a review.

package product_review_reply_archive

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductReviewReplyArchiveRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Validate seller
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Validate product ownership
	GetProductByIDAndSellerID(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID) (sqlc.Product, error)

	// 💬 Validate review
	GetProductReviewByID(ctx context.Context, id uuid.UUID) (sqlc.ProductReview, error)

	// 🔎 Validate existing reply
	GetReviewReplyByReviewIDAndSellerID(
		ctx context.Context,
		reviewID uuid.UUID,
		sellerID uuid.UUID,
	) (sqlc.ProductReviewReply, error)

	// 🗑️ Archive reply
	ArchiveReviewReply(ctx context.Context, arg sqlc.ArchiveReviewReplyParams) error
}
