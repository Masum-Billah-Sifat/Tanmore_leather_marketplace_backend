// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_edit_review/interface.go
// 🧠 Repository interface for editing a product review.
//     Handles customer, product, and review ownership validation.

package product_edit_review

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductEditReviewRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Validate customer
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Validate product
	GetProductByID(ctx context.Context, id uuid.UUID) (sqlc.Product, error)

	// 🔍 Validate review ownership and moderation
	GetProductReviewByIDAndProductIDAndReviewerID(
		ctx context.Context,
		arg sqlc.GetProductReviewByIDAndProductIDAndReviewerIDParams,
	) (sqlc.ProductReview, error)

	// ✏️ Update review text
	UpdateProductReviewText(
		ctx context.Context,
		arg sqlc.UpdateProductReviewTextParams,
	) error
}
