// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/edit_shipping_address/edit_shipping_address_repository.go
// 🧠 Concrete implementation of EditShippingAddressRepoInterface.
//     - Performs all DB operations for validating and updating checkout shipping.

package edit_shipping_address

import (
	"context"
	"database/sql"
	"strconv"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// 📦 EditShippingAddressRepository implements EditShippingAddressRepoInterface
type EditShippingAddressRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewEditShippingAddressRepository(db *sql.DB) *EditShippingAddressRepository {
	return &EditShippingAddressRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *EditShippingAddressRepository) WithTx(
	ctx context.Context,
	fn func(q *sqlc.Queries) error,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	qtx := sqlc.New(tx)
	if err := fn(qtx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// 🧑 Fetch customer
func (r *EditShippingAddressRepository) GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 🛒 Fetch checkout session
func (r *EditShippingAddressRepository) GetCheckoutSessionByID(ctx context.Context, id uuid.UUID) (sqlc.CheckoutSession, error) {
	return r.q.GetCheckoutSessionByID(ctx, id)
}

// 📍 Get shipping address by ID + session
func (r *EditShippingAddressRepository) GetShippingAddressByIDAndCheckoutSessionID(
	ctx context.Context,
	addressID, sessionID uuid.UUID,
) (sqlc.ShippingAddress, error) {
	return r.q.GetShippingAddressByIDAndCheckoutID(ctx, sqlc.GetShippingAddressByIDAndCheckoutIDParams{
		ID:                addressID,
		CheckoutSessionID: sessionID,
	})
}

// ✏️ Update address fields conditionally
func (r *EditShippingAddressRepository) UpdateShippingAddressByID(
	ctx context.Context,
	arg sqlc.UpdateShippingAddressByIDParams,
) error {
	return r.q.UpdateShippingAddressByID(ctx, arg)
}

// 💳 Update payment method
func (r *EditShippingAddressRepository) UpdateCheckoutSessionPaymentMethod(
	ctx context.Context,
	method string,
	sessionID uuid.UUID,
) error {
	return r.q.UpdateCheckoutSessionPaymentMethod(ctx, sqlc.UpdateCheckoutSessionPaymentMethodParams{
		PaymentMethod: method,
		ID:            sessionID,
	})
}

// 🔢 Count items
func (r *EditShippingAddressRepository) CountCheckoutItemsBySessionID(
	ctx context.Context,
	sessionID uuid.UUID,
) (int64, error) {
	return r.q.CountCheckoutItemsBySessionID(ctx, sessionID)
}

// 💰 Update pricing info
func (r *EditShippingAddressRepository) UpdateCheckoutSessionDeliveryPricing(
	ctx context.Context,
	deliveryCharge int,
	totalPayable int,
	sessionID uuid.UUID,
) error {
	return r.q.UpdateCheckoutSessionDeliveryPricing(ctx, sqlc.UpdateCheckoutSessionDeliveryPricingParams{
		// DeliveryCharge: int32(deliveryCharge),
		// TotalPayable:   int32(totalPayable),
		DeliveryCharge: sqlnull.String(strconv.Itoa(deliveryCharge)),
		TotalPayable:   strconv.Itoa(totalPayable),

		ID: sessionID,
	})
}
