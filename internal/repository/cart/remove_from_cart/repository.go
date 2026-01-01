// ------------------------------------------------------------
// 📁 File: internal/repository/cart/remove_from_cart/remove_from_cart_repository.go
// 🧠 Concrete implementation of RemoveFromCartRepoInterface.
//     Performs customer validation, cart item fetch, and deactivation.

package remove_from_cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 RemoveFromCartRepository implements RemoveFromCartRepoInterface
type RemoveFromCartRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewRemoveFromCartRepository(db *sql.DB) *RemoveFromCartRepository {
	return &RemoveFromCartRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *RemoveFromCartRepository) WithTx(
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
func (r *RemoveFromCartRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 🛒 Get cart item by user and variant
func (r *RemoveFromCartRepository) GetCartItemByUserAndVariant(
	ctx context.Context,
	arg sqlc.GetCartItemByUserAndVariantParams,
) (sqlc.CartItem, error) {
	return r.q.GetCartItemByUserAndVariant(ctx, arg)
}

// ❌ Deactivate cart item
func (r *RemoveFromCartRepository) DeactivateCartItem(
	ctx context.Context,
	arg sqlc.DeactivateCartItemParams,
) error {
	return r.q.DeactivateCartItem(ctx, arg)
}
