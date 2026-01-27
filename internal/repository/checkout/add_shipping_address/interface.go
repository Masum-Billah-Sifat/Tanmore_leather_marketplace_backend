// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/add_shipping_address/interface.go
// 🧠 Repository interface for adding a shipping address to a checkout session.
//     Handles user validation, session fetch, shipping address insert,
//     order item counting, Pathao API call, and checkout update.

package add_shipping_address

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type AddShippingAddressRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧑 Fetch customer by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 🛒 Fetch checkout session
	GetCheckoutSessionByID(ctx context.Context, id uuid.UUID) (sqlc.CheckoutSession, error)

	// 📦 Insert shipping address row
	InsertShippingAddress(ctx context.Context, arg sqlc.InsertShippingAddressParams) (uuid.UUID, error)

	// 🔢 Count checkout items
	CountCheckoutItemsBySessionID(ctx context.Context, checkoutSessionID uuid.UUID) (int64, error)

	// 💰 Update checkout with address, payment method, and delivery charge
	UpdateCheckoutSessionWithShipping(ctx context.Context, arg sqlc.UpdateCheckoutSessionWithShippingParams) error
}
