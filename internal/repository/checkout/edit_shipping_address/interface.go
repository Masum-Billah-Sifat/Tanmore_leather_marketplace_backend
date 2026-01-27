// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/edit_shipping_address/interface.go
// 🧠 Repository interface for editing an existing shipping address in a checkout session.
//     - Validates customer and session
//     - Checks address-session match
//     - Conditionally updates address and payment method
//     - Recalculates delivery fee if necessary

package edit_shipping_address

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type EditShippingAddressRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Validate customer
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 🛒 Get checkout session
	GetCheckoutSessionByID(ctx context.Context, id uuid.UUID) (sqlc.CheckoutSession, error)

	// 📍 Get shipping address by ID and session
	GetShippingAddressByIDAndCheckoutSessionID(ctx context.Context, addressID, sessionID uuid.UUID) (sqlc.ShippingAddress, error)

	// ✏️ Update shipping address conditionally
	UpdateShippingAddressByID(ctx context.Context, arg sqlc.UpdateShippingAddressByIDParams) error

	// 💳 Update payment method if provided
	UpdateCheckoutSessionPaymentMethod(ctx context.Context, method string, sessionID uuid.UUID) error

	// 🔢 Count items for delivery calculation
	CountCheckoutItemsBySessionID(ctx context.Context, sessionID uuid.UUID) (int64, error)

	// 💰 Update delivery charge and total payable
	UpdateCheckoutSessionDeliveryPricing(ctx context.Context, deliveryCharge int, totalPayable int, sessionID uuid.UUID) error
}
