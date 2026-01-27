// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/session_details/interface.go
// 🧠 Repository interface for retrieving a full checkout session's details.
//     Handles user moderation, session validation, address, items, and variant snapshot enrichment.

package session_details

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type CheckoutSessionDetailsRepoInterface interface {
	// 🔁 Optional transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Fetch user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 🧾 Fetch checkout session by ID
	GetCheckoutSessionByID(ctx context.Context, id uuid.UUID) (sqlc.CheckoutSession, error)

	// 📍 Fetch shipping address if present
	GetShippingAddressByIDAndCheckoutID(ctx context.Context, arg sqlc.GetShippingAddressByIDAndCheckoutIDParams) (sqlc.ShippingAddress, error)

	// 📦 Get all checkout items for session
	GetCheckoutItemsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]sqlc.GetCheckoutItemsBySessionIDRow, error)

	// 🔍 Enrich variants from product_variant_snapshots
	GetProductVariantSnapshotsByVariantIDs(ctx context.Context, variantIDs []uuid.UUID) ([]sqlc.GetProductVariantSnapshotsByVariantIDsRow, error)
}
