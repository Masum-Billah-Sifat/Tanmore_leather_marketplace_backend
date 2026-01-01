// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/interface.go
// 🧠 Repository interface for unified checkout flow (cart + product).
//     - Validates user moderation
//     - Fetches variant snapshot data
//     - Fetches active cart variant snapshot joins
//     - Inserts checkout session and item(s)

package checkout

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type CheckoutRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Validate user moderation
	GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error)

	// 📦 Fetch snapshot for single product variant
	GetVariantSnapshotByVariantID(ctx context.Context, variantID uuid.UUID) (sqlc.ProductVariantSnapshot, error)

	// 🛒 Fetch enriched cart + snapshot for multiple variants
	GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx context.Context,
		arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
	) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error)

	// 🧾 Insert new checkout session row
	InsertCheckoutSession(
		ctx context.Context,
		arg sqlc.InsertCheckoutSessionParams,
	) (sqlc.CheckoutSession, error)

	InsertCheckoutItem(
		ctx context.Context,
		arg sqlc.InsertCheckoutItemParams,
	) (sqlc.InsertCheckoutItemRow, error)
}
