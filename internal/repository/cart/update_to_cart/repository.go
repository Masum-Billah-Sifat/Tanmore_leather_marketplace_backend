package update_cart_quantity

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	// "tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

type UpdateCartQuantityRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewUpdateCartQuantityRepository(db *sql.DB) *UpdateCartQuantityRepository {
	return &UpdateCartQuantityRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

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

func (r *UpdateCartQuantityRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *UpdateCartQuantityRepository) GetVariantSnapshotByVariantID(
	ctx context.Context,
	variantID uuid.UUID,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshotByVariantID(ctx, variantID)
}

func (r *UpdateCartQuantityRepository) GetCartItemByUserAndVariant(
	ctx context.Context,
	userID uuid.UUID,
	variantID uuid.UUID,
) (sqlc.CartItem, error) {
	return r.q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
		UserID:    userID,
		VariantID: variantID,
	})
}

func (r *UpdateCartQuantityRepository) UpdateCartQuantityForUser(
	ctx context.Context,
	arg sqlc.UpdateCartQuantityForUserParams,
) error {
	return r.q.UpdateCartQuantityForUser(ctx, arg)
}
