// ------------------------------------------------------------
// 📁 File: internal/repository/product/interface.go
// 🧠 This file defines the repository interface for product creation
//     and related DB operations required by the service layer.

package product_creation

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductRepoInterface interface {
	// 🔁 Transaction handler (mandatory for atomic product creation)
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Fetch seller/user by ID (moderation + approval checks)
	GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error)

	// 📦 Insert product row
	InsertProduct(ctx context.Context, arg sqlc.InsertProductParams) (uuid.UUID, error)

	// 🧩 Insert product variant and return variant ID
	InsertProductVariantReturningID(
		ctx context.Context,
		arg sqlc.InsertProductVariantReturningIDParams,
	) (uuid.UUID, error)

	// 📨 Insert event into outbox/events table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error

	// 🧠 Fetch seller profile metadata (to extract sellerstorename)
	GetSellerProfileMetadataBySellerID(ctx context.Context, sellerID uuid.UUID) (sqlc.SellerProfileMetadatum, error)

	// 🗂️ Fetch category by ID to validate leaf + archived status
	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (sqlc.Category, error)

	// 🖼️ Insert image/video into product_medias
	InsertProductMedia(ctx context.Context, arg sqlc.InsertProductMediaParams) (uuid.UUID, error)
}
