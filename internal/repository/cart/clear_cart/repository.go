// 📁 File: internal/repository/cart/clear_cart/clear_cart_repository.go

package clear_cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ClearCartRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewClearCartRepository(db *sql.DB) *ClearCartRepository {
	return &ClearCartRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

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

func (r *ClearCartRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *ClearCartRepository) ClearCartItemsByUser(
	ctx context.Context,
	arg sqlc.ClearCartItemsByUserParams,
) error {
	return r.q.ClearCartItemsByUser(ctx, arg)
}
