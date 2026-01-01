// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_get_all_reviews/interface.go
// 🧠 Repository interface for fetching reviews (and replies) for a product.

package product_get_all_reviews

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductGetAllReviewsRepoInterface interface {
	// 🔁 Transaction wrapper (optional, for future extensibility)
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🛄 Validate product
	GetProductByID(ctx context.Context, id uuid.UUID) (sqlc.Product, error)

	// interface.go
	GetAllReviewsByProductID(ctx context.Context, arg sqlc.GetAllReviewsByProductIDParams) ([]sqlc.GetAllReviewsByProductIDRow, error)
	GetRepliesByReviewIDs(ctx context.Context, reviewIDs []uuid.UUID) ([]sqlc.GetRepliesByReviewIDsRow, error)
}
