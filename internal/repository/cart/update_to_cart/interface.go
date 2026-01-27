// ------------------------------------------------------------
// 📁 File: internal/repository/cart/update_cart_quantity/interface.go
// 🧠 Repository interface for updating cart item quantity.
//     Supports both authenticated and guest users.

package update_cart_quantity

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// type UpdateCartQuantityRepoInterface interface {
// 	// 🔁 Transaction wrapper
// 	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

// 	// 🧑 Fetch customer user by ID
// 	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

// 	// 📦 Fetch variant snapshot by variant ID
// 	GetVariantSnapshotByVariantID(
// 		ctx context.Context,
// 		arg uuid.UUID,
// 	) (sqlc.ProductVariantSnapshot, error)

// 	// 🛒 Find existing cart item (user or guest)
// 	GetCartItemByOwnerAndVariant(
// 		ctx context.Context,
// 		arg sqlc.GetCartItemByOwnerAndVariantParams,
// 	) (sqlc.GetCartItemByOwnerAndVariantRow, error)

// 	// 🔄 Update cart quantity (user or guest)
// 	UpdateCartQuantity(
// 		ctx context.Context,
// 		arg sqlc.UpdateCartQuantityParams,
// 	) error
// }

type UpdateCartQuantityRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	GetVariantSnapshotByVariantID(ctx context.Context, variantID uuid.UUID) (sqlc.ProductVariantSnapshot, error)

	// 👤 User
	GetCartItemByUserAndVariant(ctx context.Context, userID uuid.UUID, variantID uuid.UUID) (sqlc.GetCartItemByUserAndVariantRow, error)
	UpdateCartQuantityForUser(ctx context.Context, arg sqlc.UpdateCartQuantityForUserParams) error

	// 👥 Guest
	GetCartItemByGuestAndVariant(ctx context.Context, guestUserID uuid.UUID, variantID uuid.UUID) (sqlc.GetCartItemByGuestAndVariantRow, error)
	UpdateCartQuantityForGuest(ctx context.Context, arg sqlc.UpdateCartQuantityForGuestParams) error
}
