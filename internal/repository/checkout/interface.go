// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/interface.go
// 🧠 Repository interface for unified checkout flow (cart + product).
//     - Validates user moderation
//     - Fetches variant snapshot data
//     - Inserts checkout session and items (exec-only)
//     - Fetches platform promotions

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

	// 🛒 Fetch enriched cart + snapshot join
	GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx context.Context,
		arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
	) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error)

	// 🧾 Insert checkout session (exec-only)
	InsertCheckoutSession(
		ctx context.Context,
		arg sqlc.InsertCheckoutSessionParams,
	) error

	// 📄 Insert checkout item (exec-only)
	InsertCheckoutItem(
		ctx context.Context,
		arg sqlc.InsertCheckoutItemParams,
	) error

	// 🎁 Get active platform-level promotions
	GetActivePlatformPromotions(
		ctx context.Context,
	) ([]sqlc.GetActivePlatformPromotionsRow, error)

	// GetActivePlatformPromotions(
	// 	ctx context.Context,
	// ) ([]sqlc.PlatformPromotion, error)
}
