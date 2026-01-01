// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_reply_review/interface.go
// 🧠 Repository interface for replying to a review (seller only).
//     - Validates seller
//     - Validates product ownership
//     - Validates review
//     - Checks for existing reply
//     - Inserts reply

package product_reply_review

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductReplyReviewRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// ✅ Validate seller
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Validate product ownership
	GetProductByIDAndSellerID(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID) (sqlc.Product, error)

	// 🔍 Validate review
	GetProductReviewByID(ctx context.Context, id uuid.UUID) (sqlc.ProductReview, error)

	// 🔄 Check existing reply
	GetReviewReplyByReviewID(ctx context.Context, reviewID uuid.UUID) (sqlc.ProductReviewReply, error)

	// ✍️ Insert reply
	InsertReviewReply(ctx context.Context, arg sqlc.InsertReviewReplyParams) (uuid.UUID, error)
}
