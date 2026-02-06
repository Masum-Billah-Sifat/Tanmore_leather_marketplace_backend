// ------------------------------------------------------------
// 📁 File: internal/repository/cart/add_to_cart/interface.go
// 🧠 Repository interface for adding an item to the cart.
//     Handles user validation, snapshot fetch, cart lookup/update/insert, and event logging.

package add_to_cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type AddToCartRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	GetVariantSnapshotByProductIDAndVariantID(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotByProductIDAndVariantIDParams,
	) (sqlc.ProductVariantSnapshot, error)

	GetCartItemByUserAndVariant(
		ctx context.Context,
		userID uuid.UUID,
		variantID uuid.UUID,
	) (sqlc.CartItem, error)

	ReactivateCartItemByID(
		ctx context.Context,
		arg sqlc.ReactivateCartItemByIDParams,
	) error

	InsertCartItem(
		ctx context.Context,
		arg sqlc.InsertCartItemParams,
	) error
}
