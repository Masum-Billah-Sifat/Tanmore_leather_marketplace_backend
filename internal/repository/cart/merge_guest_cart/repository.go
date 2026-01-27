// ------------------------------------------------------------
// 📁 File: internal/repository/cart/merge_guest_cart/merge_guest_cart_repository.go
// 🧠 Concrete implementation of MergeGuestCartRepoInterface.
//     Handles user moderation, guest cart retrieval, inserting into user cart, and deprecating guest cart items.

package merge_guest_cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

type MergeGuestCartRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewMergeGuestCartRepository(db *sql.DB) *MergeGuestCartRepository {
	return &MergeGuestCartRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *MergeGuestCartRepository) WithTx(
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

// 🧑 Get user by ID
func (r *MergeGuestCartRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 🛒 Get active, non-deprecated guest cart items
func (r *MergeGuestCartRepository) GetActiveGuestCartItems(
	ctx context.Context,
	guestUserID uuid.UUID,
) ([]sqlc.GetActiveGuestCartItemsRow, error) {
	return r.q.GetActiveGuestCartItems(ctx, sqlnull.UUID(guestUserID))
}

// ➕ Insert item into authenticated user cart
func (r *MergeGuestCartRepository) InsertCartItem(
	ctx context.Context,
	arg sqlc.InsertCartItemParams,
) (sqlc.InsertCartItemRow, error) {
	return r.q.InsertCartItem(ctx, arg)
}

// ❌ Deprecate all guest items
func (r *MergeGuestCartRepository) DeprecateGuestCartItems(
	ctx context.Context,
	arg sqlc.DeprecateGuestCartItemsParams,
) error {
	return r.q.DeprecateGuestCartItems(ctx, arg)
}
