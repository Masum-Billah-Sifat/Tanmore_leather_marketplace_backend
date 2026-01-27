// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/confirm_order/interface.go
// 🧠 Repository interface for confirming an order from a checkout session.
//     - Moderation, session, item, snapshot, shipping address, inserts, updates, event.

package confirm_order

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ConfirmOrderRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🔍 Moderation + Checks
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	GetCheckoutSessionByID(ctx context.Context, id uuid.UUID) (sqlc.CheckoutSession, error)
	GetCheckoutItemsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]sqlc.GetCheckoutItemsBySessionIDRow, error)
	GetProductVariantSnapshotsByVariantIDs(ctx context.Context, variantIDs []uuid.UUID) ([]sqlc.GetProductVariantSnapshotsByVariantIDsRow, error)
	GetShippingAddressByIDAndCheckoutID(ctx context.Context, arg sqlc.GetShippingAddressByIDAndCheckoutIDParams) (sqlc.ShippingAddress, error)

	// 📦 Inserts + Updates
	InsertOrderRow(ctx context.Context, arg sqlc.InsertOrderRowParams) (sqlc.Order, error)
	InsertOrderItem(ctx context.Context, arg sqlc.InsertOrderItemParams) error
	UpdateCheckoutSessionStatusToOrderCreated(ctx context.Context, id uuid.UUID, newStatus string) error
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
