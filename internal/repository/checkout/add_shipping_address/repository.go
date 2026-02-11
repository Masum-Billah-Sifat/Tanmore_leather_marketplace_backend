// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/add_shipping_address/add_shipping_address_repository.go
// 🧠 Concrete implementation of AddShippingAddressRepoInterface.
//     Handles validation, address insert, Pathao delivery charge flow, and checkout update.

package add_shipping_address

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 AddShippingAddressRepository implements AddShippingAddressRepoInterface
type AddShippingAddressRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewAddShippingAddressRepository(db *sql.DB) *AddShippingAddressRepository {
	return &AddShippingAddressRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *AddShippingAddressRepository) WithTx(
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

// 🧑 Get customer by ID
func (r *AddShippingAddressRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 🛒 Fetch checkout session
func (r *AddShippingAddressRepository) GetCheckoutSessionByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.CheckoutSession, error) {
	return r.q.GetCheckoutSessionByID(ctx, id)
}

// 📦 Insert shipping address
func (r *AddShippingAddressRepository) InsertShippingAddress(
	ctx context.Context,
	arg sqlc.InsertShippingAddressParams,
) (uuid.UUID, error) {
	return r.q.InsertShippingAddress(ctx, arg)
}

// 🔢 Count checkout items
func (r *AddShippingAddressRepository) CountCheckoutItemsBySessionID(
	ctx context.Context,
	checkoutSessionID uuid.UUID,
) (int64, error) {
	return r.q.CountCheckoutItemsBySessionID(ctx, checkoutSessionID)
}

// 💰 Update checkout session with shipping info
func (r *AddShippingAddressRepository) UpdateCheckoutSessionWithShipping(
	ctx context.Context,
	arg sqlc.UpdateCheckoutSessionWithShippingParams,
) error {
	return r.q.UpdateCheckoutSessionWithShipping(ctx, arg)
}

// internal/repository/checkout/add_shipping_address/add_shipping_address_repository.go

func (r *AddShippingAddressRepository) MarkCheckoutSessionReadyToOrder(
	ctx context.Context,
	arg sqlc.MarkCheckoutSessionReadyToOrderParams,
) (uuid.UUID, error) {
	return r.q.MarkCheckoutSessionReadyToOrder(ctx, arg)
}
