// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/session_details/session_details_repository.go
// 🧠 Concrete implementation of CheckoutSessionDetailsRepoInterface.

package session_details

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type CheckoutSessionDetailsRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewCheckoutSessionDetailsRepository(db *sql.DB) *CheckoutSessionDetailsRepository {
	return &CheckoutSessionDetailsRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *CheckoutSessionDetailsRepository) WithTx(
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

// 🧑 User moderation check
func (r *CheckoutSessionDetailsRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 🧾 Get checkout session
func (r *CheckoutSessionDetailsRepository) GetCheckoutSessionByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.CheckoutSession, error) {
	return r.q.GetCheckoutSessionByID(ctx, id)
}

// 📍 Get shipping address
func (r *CheckoutSessionDetailsRepository) GetShippingAddressByIDAndCheckoutID(
	ctx context.Context,
	arg sqlc.GetShippingAddressByIDAndCheckoutIDParams,
) (sqlc.ShippingAddress, error) {
	return r.q.GetShippingAddressByIDAndCheckoutID(ctx, arg)
}

// 📦 Get checkout items
func (r *CheckoutSessionDetailsRepository) GetCheckoutItemsBySessionID(
	ctx context.Context,
	sessionID uuid.UUID,
) ([]sqlc.GetCheckoutItemsBySessionIDRow, error) {
	return r.q.GetCheckoutItemsBySessionID(ctx, sessionID)
}

// 🔍 Get variant snapshots
func (r *CheckoutSessionDetailsRepository) GetProductVariantSnapshotsByVariantIDs(
	ctx context.Context,
	variantIDs []uuid.UUID,
) ([]sqlc.GetProductVariantSnapshotsByVariantIDsRow, error) {
	return r.q.GetProductVariantSnapshotsByVariantIDs(ctx, variantIDs)
}
