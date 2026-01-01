// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_set_primary_image/interface.go
// 🧠 Repository interface for setting an existing image as primary.
//     Handles seller validation, product validation, image update, and event insertion.

package product_set_primary_image

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductSetPrimaryImageRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Fetch seller by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Fetch product by ID and seller ID
	GetProductByIDAndSellerID(
		ctx context.Context,
		arg sqlc.GetProductByIDAndSellerIDParams,
	) (sqlc.Product, error)

	// 🚫 Clear all previous primary flags (is_primary = false)
	UnsetAllPrimaryImages(
		ctx context.Context,
		arg sqlc.UnsetAllPrimaryImagesParams,
	) error

	// ✅ Set selected image as primary (is_primary = true)
	SetAsPrimaryImage(
		ctx context.Context,
		arg sqlc.SetAsPrimaryImageParams,
	) error

	// 📨 Insert event into outbox
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error

	// // 🖼️ Fetch a single product media row by ID and product ID
	// GetProductMediaByID(
	// 	ctx context.Context,
	// 	arg sqlc.GetProductMediaByIDParams,
	// ) (sqlc.ProductMedia, error)

	GetProductMediaByID(ctx context.Context, arg sqlc.GetProductMediaByIDParams) (sqlc.GetProductMediaByIDRow, error)
}
