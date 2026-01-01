// ------------------------------------------------------------
// 📁 File: internal/repository/cart/remove_from_cart/interface.go
// 🧠 Repository interface for removing an item from the cart.
//     Handles user validation, cart lookup, and soft deletion.

package remove_from_cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type RemoveFromCartRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Fetch customer user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 🛒 Fetch cart item for given user and variant
	GetCartItemByUserAndVariant(
		ctx context.Context,
		arg sqlc.GetCartItemByUserAndVariantParams,
	) (sqlc.CartItem, error)

	// ❌ Soft delete (deactivate) cart item
	DeactivateCartItem(
		ctx context.Context,
		arg sqlc.DeactivateCartItemParams,
	) error
}
