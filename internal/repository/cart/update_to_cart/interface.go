// ------------------------------------------------------------
// 📁 File: internal/repository/cart/update_cart_quantity/interface.go
// 🧠 Repository interface for updating cart item quantity.
//     Handles user validation, snapshot fetch, cart fetch, and quantity update.

package update_cart_quantity

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type UpdateCartQuantityRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Fetch customer user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Fetch variant snapshot by variant ID
	GetVariantSnapshotByVariantID(
		ctx context.Context,
		arg uuid.UUID,
	) (sqlc.ProductVariantSnapshot, error)

	// 🛒 Find existing cart item
	GetCartItemByUserAndVariant(
		ctx context.Context,
		arg sqlc.GetCartItemByUserAndVariantParams,
	) (sqlc.CartItem, error)

	// 🔄 Update cart quantity
	UpdateCartQuantity(
		ctx context.Context,
		arg sqlc.UpdateCartQuantityParams,
	) error
}
