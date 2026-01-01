// ------------------------------------------------------------
// 📁 File: internal/repository/cart/get_all_cart_items/get_all_cart_items_repository.go
// 🧠 Concrete implementation of GetAllCartItemsRepoInterface.
//     Handles customer moderation, variant ID lookup, and snapshot-enriched cart data.

package get_all_cart_items

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 GetAllCartItemsRepository implements GetAllCartItemsRepoInterface
type GetAllCartItemsRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewGetAllCartItemsRepository(db *sql.DB) *GetAllCartItemsRepository {
	return &GetAllCartItemsRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper (future-proofed)
func (r *GetAllCartItemsRepository) WithTx(
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
func (r *GetAllCartItemsRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 🧾 Get all active variant IDs in user's cart
func (r *GetAllCartItemsRepository) ListActiveVariantIDsByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {
	return r.q.ListActiveVariantIDsByUser(ctx, userID)
}

// 🔍 Fetch full cart + snapshot joined rows
func (r *GetAllCartItemsRepository) GetActiveCartVariantSnapshotsByUserAndVariantIDs(
	ctx context.Context,
	arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error) {
	return r.q.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx, arg)
}
