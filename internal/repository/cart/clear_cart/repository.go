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

// // 🧹 Clear all active cart items for user
// func (r *ClearCartRepository) ClearCartItemsForUser(
// 	ctx context.Context,
// 	arg sqlc.ClearCartItemsForUserParams,
// ) (int64, error) {
// 	return r.q.ClearCartItemsForUser(ctx, arg)
// }

func (r *ClearCartRepository) ClearCartItemsForUser(
	ctx context.Context,
	arg sqlc.ClearCartItemsForUserParams,
) error {
	return r.q.ClearCartItemsForUser(ctx, arg)
}
