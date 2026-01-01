// ------------------------------------------------------------
// 📁 File: internal/repository/cart/get_all_cart_items/interface.go
// 🧠 Repository interface for retrieving grouped active cart items.
//     Handles user moderation, active variant lookup, and enriched snapshot query.

package get_all_cart_items

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type GetAllCartItemsRepoInterface interface {
	// 🔁 Transaction wrapper (optional for future expansion)
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Fetch user by ID for moderation checks
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 🧾 List all active variant IDs in user's cart
	ListActiveVariantIDsByUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)

	// 🔍 Fetch full cart item + variant snapshot for valid variant IDs
	GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx context.Context,
		arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
	) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error)
}
