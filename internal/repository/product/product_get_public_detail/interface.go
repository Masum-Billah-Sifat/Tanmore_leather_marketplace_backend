// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_get_public_detail/interface.go
// 🧠 Repository interface for fetching public product detail (customer-facing)

package product_get_public_detail

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductGetPublicDetailRepoInterface interface {
	// 🔁 Optional transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🌐 Core query from product_variant_indexes for public view
	GetProductDetailByProductID(
		ctx context.Context,
		productID uuid.UUID,
	) ([]sqlc.GetProductDetailByProductIDRow, error)
}
