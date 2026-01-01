// ------------------------------------------------------------
// 📁 File: internal/repository/product_variant/interface.go
// 🧠 Repository interface for adding a new variant to an existing product.
//     Contains only DB operations required by the service layer.

package product_variant

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductVariantRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Fetch seller/user by ID (moderation & approval checks)
	GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error)

	// 📦 Verify product ownership
	GetProductByIDAndSellerID(
		ctx context.Context,
		arg sqlc.GetProductByIDAndSellerIDParams,
	) (sqlc.Product, error)

	// 🧩 Insert product variant and return variant ID
	InsertProductVariantReturningID(
		ctx context.Context,
		arg sqlc.InsertProductVariantReturningIDParams,
	) (uuid.UUID, error)

	// 📨 Insert event into events table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error

	// 🧠 Fetch category details
	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (sqlc.Category, error)

	// 🧠 Fetch seller profile metadata
	GetSellerProfileMetadataBySellerID(ctx context.Context, sellerID uuid.UUID) (sqlc.SellerProfileMetadatum, error)

	// 🖼️ Fetch all non-archived medias by type
	GetActiveMediasByProductID(
		ctx context.Context,
		arg sqlc.GetActiveMediasByProductIDParams,
	) ([]sqlc.ProductMedia, error)

	// 🖼️ Fetch primary product image
	GetPrimaryProductImageByProductID(
		ctx context.Context,
		arg sqlc.GetPrimaryProductImageByProductIDParams,
	) (sqlc.ProductMedia, error)
}
