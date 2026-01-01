// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_update_info/interface.go
// 🧠 Repository interface for updating product info (title/description).
//     Contains DB operations: user check, product check, update, and event insert.

package product_update_info

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductUpdateInfoRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🔍 Fetch user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 🔍 Fetch product by ID and seller ID
	GetProductByIDAndSellerID(ctx context.Context, arg sqlc.GetProductByIDAndSellerIDParams) (sqlc.Product, error)

	// 📝 Update product title/description
	UpdateProductTitleDesc(ctx context.Context, arg sqlc.UpdateProductTitleDescParams) error

	// 📨 Insert event into outbox
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
