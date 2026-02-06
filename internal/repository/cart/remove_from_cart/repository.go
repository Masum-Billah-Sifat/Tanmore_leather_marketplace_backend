package remove_from_cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

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

// 🧑 Get authenticated user by ID
func (r *RemoveFromCartRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *RemoveFromCartRepository) GetCartItemByUserAndVariant(
	ctx context.Context,
	userID uuid.UUID,
	variantID uuid.UUID,
) (sqlc.CartItem, error) {
	return r.q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
		UserID: userID, VariantID: variantID,
	})
}

func (r *RemoveFromCartRepository) DeactivateCartItemByUserAndVariant(
	ctx context.Context,
	arg sqlc.DeactivateCartItemByUserAndVariantParams,
) error {
	return r.q.DeactivateCartItemByUserAndVariant(ctx, arg)
}
