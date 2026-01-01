// ------------------------------------------------------------
// 📁 File: internal/repository/checkout/checkout_repository.go
// 🧠 Concrete implementation of CheckoutRepoInterface.
//     - Handles moderation, snapshot fetch, session + item inserts

package checkout

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 CheckoutRepository implements CheckoutRepoInterface
type CheckoutRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewCheckoutRepository(db *sql.DB) *CheckoutRepository {
	return &CheckoutRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *CheckoutRepository) WithTx(
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

// 🧑 Fetch customer by ID
func (r *CheckoutRepository) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, userID)
}

// 📦 Fetch snapshot for single variant (Buy Now)
func (r *CheckoutRepository) GetVariantSnapshotByVariantID(
	ctx context.Context,
	variantID uuid.UUID,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshotByVariantID(ctx, variantID)
}

// 🛒 Fetch enriched cart + snapshot join (Bulk cart)
func (r *CheckoutRepository) GetActiveCartVariantSnapshotsByUserAndVariantIDs(
	ctx context.Context,
	arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error) {
	return r.q.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx, arg)
}

// 🧾 Insert checkout session row
func (r *CheckoutRepository) InsertCheckoutSession(
	ctx context.Context,
	arg sqlc.InsertCheckoutSessionParams,
) (sqlc.CheckoutSession, error) {
	return r.q.InsertCheckoutSession(ctx, arg)
}

// // 📄 Insert one checkout item row
// func (r *CheckoutRepository) InsertCheckoutItem(
// 	ctx context.Context,
// 	arg sqlc.InsertCheckoutItemParams,
// ) (sqlc.CheckoutItem, error) {
// 	return r.q.InsertCheckoutItem(ctx, arg)
// }

func (r *CheckoutRepository) InsertCheckoutItem(
	ctx context.Context,
	arg sqlc.InsertCheckoutItemParams,
) (sqlc.InsertCheckoutItemRow, error) {
	return r.q.InsertCheckoutItem(ctx, arg)
}
