// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/review/interface.go
// 🧠 Repository interface for reviewing a checkout session.
//     Includes user check, session/address/item fetch, and variant validation.

package review

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ReviewCheckoutRepoInterface interface {
	// 🔁 Optional transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Fetch customer and moderation flags
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 🧾 Fetch checkout session
	GetCheckoutSessionByID(ctx context.Context, id uuid.UUID) (sqlc.CheckoutSession, error)

	// 🧭 Fetch associated shipping address
	GetShippingAddressByIDAndCheckoutID(ctx context.Context, arg sqlc.GetShippingAddressByIDAndCheckoutIDParams) (sqlc.ShippingAddress, error)

	// 📦 Fetch checkout items
	GetCheckoutItemsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]sqlc.GetCheckoutItemsBySessionIDRow, error)

	// 🔍 Fetch product variant snapshots for validation
	GetProductVariantSnapshotsByVariantIDs(ctx context.Context, variantIDs []uuid.UUID) ([]sqlc.GetProductVariantSnapshotsByVariantIDsRow, error)
}
