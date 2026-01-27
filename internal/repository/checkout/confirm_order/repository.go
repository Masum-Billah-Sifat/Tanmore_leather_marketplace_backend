// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/confirm_order/confirm_order_repository.go
// 🧠 Concrete implementation of ConfirmOrderRepoInterface.

package confirm_order

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ConfirmOrderRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewConfirmOrderRepository(db *sql.DB) *ConfirmOrderRepository {
	return &ConfirmOrderRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ConfirmOrderRepository) WithTx(
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

func (r *ConfirmOrderRepository) GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *ConfirmOrderRepository) GetCheckoutSessionByID(ctx context.Context, id uuid.UUID) (sqlc.CheckoutSession, error) {
	return r.q.GetCheckoutSessionByID(ctx, id)
}

func (r *ConfirmOrderRepository) GetCheckoutItemsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]sqlc.GetCheckoutItemsBySessionIDRow, error) {
	return r.q.GetCheckoutItemsBySessionID(ctx, sessionID)
}

func (r *ConfirmOrderRepository) GetProductVariantSnapshotsByVariantIDs(ctx context.Context, variantIDs []uuid.UUID) ([]sqlc.GetProductVariantSnapshotsByVariantIDsRow, error) {
	return r.q.GetProductVariantSnapshotsByVariantIDs(ctx, variantIDs)
}

func (r *ConfirmOrderRepository) GetShippingAddressByIDAndCheckoutID(ctx context.Context, arg sqlc.GetShippingAddressByIDAndCheckoutIDParams) (sqlc.ShippingAddress, error) {
	return r.q.GetShippingAddressByIDAndCheckoutID(ctx, arg)
}

func (r *ConfirmOrderRepository) InsertOrderRow(ctx context.Context, arg sqlc.InsertOrderRowParams) (sqlc.Order, error) {
	return r.q.InsertOrderRow(ctx, arg)
}

func (r *ConfirmOrderRepository) InsertOrderItem(ctx context.Context, arg sqlc.InsertOrderItemParams) error {
	return r.q.InsertOrderItem(ctx, arg)
}

func (r *ConfirmOrderRepository) UpdateCheckoutSessionStatusToOrderCreated(ctx context.Context, id uuid.UUID, newStatus string) error {
	return r.q.UpdateCheckoutSessionStatusToOrderCreated(ctx, sqlc.UpdateCheckoutSessionStatusToOrderCreatedParams{
		ID:     id,
		Status: newStatus,
	})
}

func (r *ConfirmOrderRepository) InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error {
	return r.q.InsertEvent(ctx, arg)
}
