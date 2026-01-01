// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_archive_media/interface.go
// 🧠 Repository interface for archiving media (image or promo video) of a product.
//     Handles user/product/media validation, active image count, archive update, and event emission.

package product_archive_media

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductArchiveMediaRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Validate seller
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Validate product ownership
	GetProductByIDAndSellerID(ctx context.Context, arg sqlc.GetProductByIDAndSellerIDParams) (sqlc.Product, error)

	// 🖼️ Fetch media row
	// GetProductMediaByID(ctx context.Context, arg sqlc.GetProductMediaByIDParams) (sqlc.ProductMedia, error)
	// 🖼️ Fetch media row
	GetProductMediaByID(ctx context.Context, arg sqlc.GetProductMediaByIDParams) (sqlc.GetProductMediaByIDRow, error)

	// 🔢 Count active media for given type
	CountActiveImagesForProduct(ctx context.Context, arg sqlc.CountActiveImagesForProductParams) (int64, error)

	// 🗑️ Archive the media
	ArchiveProductMedia(ctx context.Context, arg sqlc.ArchiveProductMediaParams) error

	// 📨 Emit event
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
