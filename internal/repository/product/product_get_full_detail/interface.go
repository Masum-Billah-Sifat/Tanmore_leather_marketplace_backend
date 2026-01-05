// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_get_full_detail/interface.go
// 🧠 Repository interface for fetching full product detail for sellers.

package product_get_full_detail

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductGetFullDetailRepoInterface interface {
	// 🔁 Optional transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Validate seller identity
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Validate product ownership & status
	GetProductByIDAndSellerID(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID) (sqlc.Product, error)

	// 🧩 Get all variants from product_variant_indexes
	// GetVariantIndexesByProductAndSeller(ctx context.Context, arg sqlc.GetVariantIndexesByProductAndSellerParams) ([]sqlc.GetVariantIndexesByProductAndSellerParams, error)
	GetVariantIndexesByProductAndSeller(
		ctx context.Context,
		arg sqlc.GetVariantIndexesByProductAndSellerParams,
	) ([]sqlc.ProductVariantIndex, error)

	// 🖼️ Get primary image from product_medias
	GetPrimaryImageForProduct(ctx context.Context, arg sqlc.GetPrimaryProductImageByProductIDParams) (sqlc.ProductMedia, error)
}
