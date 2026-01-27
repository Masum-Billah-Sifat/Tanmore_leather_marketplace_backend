package remove_from_cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// type RemoveFromCartRepoInterface interface {
// 	// 🔁 Transaction wrapper
// 	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

// 	// 🧑 Fetch customer user by ID (for authenticated only)
// 	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

// 	// 🛒 Get cart item by user OR guest owner
// 	GetCartItemByOwnerAndVariant(
// 		ctx context.Context,
// 		arg sqlc.GetCartItemByOwnerAndVariantParams,
// 	) (sqlc.GetCartItemByOwnerAndVariantRow, error)

// 	// ❌ Deactivate cart item by user OR guest owner
// 	DeactivateCartItemByOwnerAndVariant(
// 		ctx context.Context,
// 		arg sqlc.DeactivateCartItemByOwnerAndVariantParams,
// 	) error
// }

type RemoveFromCartRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 👤 User
	GetCartItemByUserAndVariant(ctx context.Context, userID uuid.UUID, variantID uuid.UUID) (sqlc.GetCartItemByUserAndVariantRow, error)
	DeactivateCartItemByUserAndVariant(ctx context.Context, arg sqlc.DeactivateCartItemByUserAndVariantParams) error

	// 👥 Guest
	GetCartItemByGuestAndVariant(ctx context.Context, guestUserID uuid.UUID, variantID uuid.UUID) (sqlc.GetCartItemByGuestAndVariantRow, error)
	DeactivateCartItemByGuestAndVariant(ctx context.Context, arg sqlc.DeactivateCartItemByGuestAndVariantParams) error
}
