// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_add_review/interface.go
// 🧠 Repository interface for adding a review to a product.
//     Handles customer + product validation and review insert.

package product_add_review

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductAddReviewRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Validate customer
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Validate product existence and moderation
	GetProductByID(ctx context.Context, id uuid.UUID) (sqlc.Product, error)

	// ✍️ Insert product review
	InsertProductReview(
		ctx context.Context,
		arg sqlc.InsertProductReviewParams,
	) (uuid.UUID, error)
}
