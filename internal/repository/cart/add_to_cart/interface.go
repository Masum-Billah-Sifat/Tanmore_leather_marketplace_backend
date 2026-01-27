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

// type AddToCartRepoInterface interface {
// 	// 🔁 Transaction wrapper
// 	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

// 	// 🧑 Fetch customer user by ID
// 	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

// 	// 📦 Fetch snapshot for given product + variant
// 	GetVariantSnapshotByProductIDAndVariantID(
// 		ctx context.Context,
// 		arg sqlc.GetVariantSnapshotByProductIDAndVariantIDParams,
// 	) (sqlc.ProductVariantSnapshot, error)

// 	GetCartItemByOwnerAndVariant(
// 		ctx context.Context,
// 		arg sqlc.GetCartItemByOwnerAndVariantParams,
// 	) (sqlc.GetCartItemByOwnerAndVariantRow, error)

// 	// ♻️ Reactivate a cart item (if already exists but inactive)
// 	ReactivateCartItemByID(
// 		ctx context.Context,
// 		arg sqlc.ReactivateCartItemByIDParams,
// 	) error

// 	InsertCartItem(
// 		ctx context.Context,
// 		arg sqlc.InsertCartItemParams,
// 	) (sqlc.InsertCartItemRow, error)
// }

type AddToCartRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	GetVariantSnapshotByProductIDAndVariantID(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotByProductIDAndVariantIDParams,
	) (sqlc.ProductVariantSnapshot, error)

	// 👤 User
	GetCartItemByUserAndVariant(
		ctx context.Context,
		userID uuid.UUID,
		variantID uuid.UUID,
	) (sqlc.GetCartItemByUserAndVariantRow, error)

	// 👥 Guest
	GetCartItemByGuestAndVariant(
		ctx context.Context,
		guestUserID uuid.UUID,
		variantID uuid.UUID,
	) (sqlc.GetCartItemByGuestAndVariantRow, error)

	// 	// In your interface file:
	// GetCartItemByGuestAndVariant(ctx context.Context, guestUserID, variantID uuid.UUID) (sqlc.GetCartItemByGuestAndVariantRow, error)

	ReactivateCartItemByID(
		ctx context.Context,
		arg sqlc.ReactivateCartItemByIDParams,
	) error

	InsertCartItem(
		ctx context.Context,
		arg sqlc.InsertCartItemParams,
	) (sqlc.InsertCartItemRow, error)
}
