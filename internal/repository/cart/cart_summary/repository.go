package cart_summary

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 CartSummaryRepository implements CartSummaryRepoInterface
type CartSummaryRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewCartSummaryRepository(db *sql.DB) *CartSummaryRepository {
	return &CartSummaryRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *CartSummaryRepository) WithTx(
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

// 🧑 Get customer user by ID
func (r *CartSummaryRepository) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, userID)
}

// // 🛒 Fetch cart + snapshot rows for selected variant IDs (user OR guest)
// func (r *CartSummaryRepository) GetActiveCartVariantSnapshotsByOwnerAndVariantIDs(
// 	ctx context.Context,
// 	arg sqlc.GetActiveCartVariantSnapshotsByOwnerAndVariantIDsParams,
// ) ([]sqlc.GetActiveCartVariantSnapshotsByOwnerAndVariantIDsRow, error) {
// 	return r.q.GetActiveCartVariantSnapshotsByOwnerAndVariantIDs(ctx, arg)
// }

func (r *CartSummaryRepository) GetActiveCartVariantSnapshotsByUserAndVariantIDs(
	ctx context.Context,
	arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error) {
	return r.q.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx, arg)
}

func (r *CartSummaryRepository) GetActiveCartVariantSnapshotsByGuestAndVariantIDs(
	ctx context.Context,
	arg sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsParams,
) ([]sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsRow, error) {
	return r.q.GetActiveCartVariantSnapshotsByGuestAndVariantIDs(ctx, arg)
}
