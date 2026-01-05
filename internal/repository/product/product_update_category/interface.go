// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_update_category/interface.go
// 🧠 Repository interface for updating a product's category.
//     Includes user check, product check, category check, update, and event insert.

package product_update_category

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductUpdateCategoryRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Validate seller identity
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Check product ownership + moderation
	GetProductByIDAndSellerID(ctx context.Context, arg sqlc.GetProductByIDAndSellerIDParams) (sqlc.Product, error)

	// 🌳 Validate new category
	GetCategoryByID(ctx context.Context, id uuid.UUID) (sqlc.Category, error)

	// ✏️ Update category of the product
	UpdateProductCategory(ctx context.Context, arg sqlc.UpdateProductCategoryParams) error

	// 📨 Insert event into outbox
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
