// ------------------------------------------------------------
// 📁 File: internal/repository/cart/merge_guest_cart/interface.go
// 🧠 Repository interface for merging guest cart into authenticated user's cart.
//     Handles user moderation, guest cart read, item insert, and deprecation.

package merge_guest_cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type MergeGuestCartRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Fetch authenticated customer user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 🛒 Fetch all active, non-deprecated guest cart items
	GetActiveGuestCartItems(ctx context.Context, guestUserID uuid.UUID) ([]sqlc.GetActiveGuestCartItemsRow, error)

	// ➕ Insert a cart item for authenticated user (copy from guest)
	InsertCartItem(ctx context.Context, arg sqlc.InsertCartItemParams) (sqlc.InsertCartItemRow, error)

	// ❌ Deprecate all guest cart items after successful merge
	DeprecateGuestCartItems(ctx context.Context, arg sqlc.DeprecateGuestCartItemsParams) error
}
