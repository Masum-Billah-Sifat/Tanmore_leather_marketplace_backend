// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_get_all_grouped/interface.go
// 🧠 Repository interface for fetching ALL seller products with grouped variants.

package product_get_all_grouped

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductGetAllGroupedRepoInterface interface {
	// 🔁 Optional transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Validate seller identity
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Fetch all variant indexes grouped by product
	GetAllProductVariantIndexesBySeller(ctx context.Context, sellerID uuid.UUID) ([]sqlc.ProductVariantIndex, error)

	// 🖼️ Get primary image for product
	GetPrimaryImageForProduct(ctx context.Context, arg sqlc.GetPrimaryProductImageByProductIDParams) (sqlc.ProductMedia, error)
}
