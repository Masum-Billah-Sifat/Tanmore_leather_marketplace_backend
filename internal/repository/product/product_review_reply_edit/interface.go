// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_edit_reply/interface.go
// 🧠 Repository interface for editing a seller's reply to a review.

package product_edit_reply

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductReviewReplyEditRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Validate seller
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Validate product ownership
	GetProductByIDAndSellerID(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID) (sqlc.Product, error)

	// 💬 Validate review
	GetProductReviewByID(ctx context.Context, id uuid.UUID) (sqlc.ProductReview, error)

	// 🔎 Validate existing reply
	GetReviewReplyByReviewIDAndSellerID(ctx context.Context, reviewID uuid.UUID, sellerID uuid.UUID) (sqlc.ProductReviewReply, error)

	// ✏️ Update reply text
	UpdateReviewReplyText(ctx context.Context, arg sqlc.UpdateReviewReplyTextParams) error
}
