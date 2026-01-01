// ------------------------------------------------------------
// 📁 File: internal/repository/cart/cart_summary/interface.go
// 🧠 Repository interface for calculating cart order summary.
//     - Validates user moderation
//     - Retrieves enriched cart + snapshot rows for selected variant IDs

package cart_summary

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type CartSummaryRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Fetch customer user by ID
	GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error)

	// 📦 Fetch active cart + snapshot for selected variants
	GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx context.Context,
		arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
	) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error)
}
