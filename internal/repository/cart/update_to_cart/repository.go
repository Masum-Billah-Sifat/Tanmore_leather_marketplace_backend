// ------------------------------------------------------------
// 📁 File: internal/repository/cart/update_cart_quantity/update_cart_quantity_repository.go
// 🧠 Concrete implementation of UpdateCartQuantityRepoInterface.
//     Handles moderation, snapshot check, cart lookup, and quantity update.

package update_cart_quantity

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 UpdateCartQuantityRepository implements UpdateCartQuantityRepoInterface
type UpdateCartQuantityRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewUpdateCartQuantityRepository(db *sql.DB) *UpdateCartQuantityRepository {
	return &UpdateCartQuantityRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *UpdateCartQuantityRepository) WithTx(
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
func (r *UpdateCartQuantityRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Fetch variant snapshot by variant ID
func (r *UpdateCartQuantityRepository) GetVariantSnapshotByVariantID(
	ctx context.Context,
	variantID uuid.UUID,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshotByVariantID(ctx, variantID)
}

// 🛒 Find existing cart item
func (r *UpdateCartQuantityRepository) GetCartItemByUserAndVariant(
	ctx context.Context,
	arg sqlc.GetCartItemByUserAndVariantParams,
) (sqlc.CartItem, error) {
	return r.q.GetCartItemByUserAndVariant(ctx, arg)
}

// 🔄 Update quantity
func (r *UpdateCartQuantityRepository) UpdateCartQuantity(
	ctx context.Context,
	arg sqlc.UpdateCartQuantityParams,
) error {
	return r.q.UpdateCartQuantity(ctx, arg)
}
