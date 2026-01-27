// ------------------------------------------------------------
// 📁 File: internal/repository/cart/clear_cart/clear_cart_repository.go
// 🧠 Concrete implementation of ClearCartRepoInterface.
//     Performs user check and deactivates all active cart items.

package clear_cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ClearCartRepository implements ClearCartRepoInterface
type ClearCartRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewClearCartRepository(db *sql.DB) *ClearCartRepository {
	return &ClearCartRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ClearCartRepository) WithTx(
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
func (r *ClearCartRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// // 🧹 Clear all cart items for either user or guest
// func (r *ClearCartRepository) ClearCartItemsByOwner(
// 	ctx context.Context,
// 	arg sqlc.ClearCartItemsByOwnerParams,
// ) error {
// 	return r.q.ClearCartItemsByOwner(ctx, arg)
// }

func (r *ClearCartRepository) ClearCartItemsByUser(
	ctx context.Context,
	arg sqlc.ClearCartItemsByUserParams,
) error {
	return r.q.ClearCartItemsByUser(ctx, arg)
}

func (r *ClearCartRepository) ClearCartItemsByGuest(
	ctx context.Context,
	arg sqlc.ClearCartItemsByGuestParams,
) error {
	return r.q.ClearCartItemsByGuest(ctx, arg)
}
