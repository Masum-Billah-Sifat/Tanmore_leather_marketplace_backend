package remove_from_cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/pkg/sqlnull"

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

// // 🛒 Get cart item by user or guest owner
// func (r *RemoveFromCartRepository) GetCartItemByOwnerAndVariant(
// 	ctx context.Context,
// 	arg sqlc.GetCartItemByOwnerAndVariantParams,
// ) (sqlc.GetCartItemByOwnerAndVariantRow, error) {
// 	return r.q.GetCartItemByOwnerAndVariant(ctx, arg)
// }

func (r *RemoveFromCartRepository) GetCartItemByUserAndVariant(
	ctx context.Context,
	userID uuid.UUID,
	variantID uuid.UUID,
) (sqlc.GetCartItemByUserAndVariantRow, error) {
	return r.q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
		UserID: sqlnull.UUID(userID), VariantID: variantID,
	})
}

func (r *RemoveFromCartRepository) GetCartItemByGuestAndVariant(
	ctx context.Context,
	guestUserID uuid.UUID,
	variantID uuid.UUID,
) (sqlc.GetCartItemByGuestAndVariantRow, error) {
	return r.q.GetCartItemByGuestAndVariant(ctx, sqlc.GetCartItemByGuestAndVariantParams{
		GuestUserID: sqlnull.UUID(guestUserID), VariantID: variantID,
	})
}

// // ❌ Deactivate cart item by user or guest owner
// func (r *RemoveFromCartRepository) DeactivateCartItemByOwnerAndVariant(
// 	ctx context.Context,
// 	arg sqlc.DeactivateCartItemByOwnerAndVariantParams,
// ) error {
// 	return r.q.DeactivateCartItemByOwnerAndVariant(ctx, arg)
// }

func (r *RemoveFromCartRepository) DeactivateCartItemByUserAndVariant(
	ctx context.Context,
	arg sqlc.DeactivateCartItemByUserAndVariantParams,
) error {
	return r.q.DeactivateCartItemByUserAndVariant(ctx, arg)
}

func (r *RemoveFromCartRepository) DeactivateCartItemByGuestAndVariant(
	ctx context.Context,
	arg sqlc.DeactivateCartItemByGuestAndVariantParams,
) error {
	return r.q.DeactivateCartItemByGuestAndVariant(ctx, arg)
}
