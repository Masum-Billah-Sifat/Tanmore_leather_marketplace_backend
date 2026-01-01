// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_archive/interface.go
// 🧠 Repository interface for product archival.
//     Includes seller validation, product ownership check, archiving, and event insertion.

package product_archive

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductArchiveRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Validate seller (must be active and approved)
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Check if product exists and belongs to seller
	GetProductByIDAndSellerID(ctx context.Context, productID, sellerID uuid.UUID) (sqlc.Product, error)

	// 🧹 Soft-archive the product
	ArchiveProduct(ctx context.Context, arg sqlc.ArchiveProductParams) error

	// 📨 Insert product.archived event
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
