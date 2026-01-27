// ------------------------------------------------------------
// 📁 File: internal/repository/cart/clear_cart/interface.go
// 🧠 Repository interface for clearing all active items in a customer's cart.

package clear_cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// type ClearCartRepoInterface interface {
// 	// 🔁 Transaction wrapper
// 	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

// 	// 🧑 Fetch customer user by ID
// 	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

// 	// 🧹 Soft delete all active cart items (user or guest)
// 	ClearCartItemsByOwner(
// 		ctx context.Context,
// 		arg sqlc.ClearCartItemsByOwnerParams,
// 	) error
// }

type ClearCartRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	ClearCartItemsByUser(ctx context.Context, arg sqlc.ClearCartItemsByUserParams) error
	ClearCartItemsByGuest(ctx context.Context, arg sqlc.ClearCartItemsByGuestParams) error
}
