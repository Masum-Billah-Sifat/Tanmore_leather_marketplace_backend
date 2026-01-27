// ------------------------------------------------------------
// 📁 File: internal/repository/cart/update_cart_quantity/update_cart_quantity_repository.go
// 🧠 Concrete implementation of UpdateCartQuantityRepoInterface.
//     Supports user and guest cart item updates.

package update_cart_quantity

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/pkg/sqlnull"

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

// // 🛒 Fetch cart item by user or guest + variant
// func (r *UpdateCartQuantityRepository) GetCartItemByOwnerAndVariant(
// 	ctx context.Context,
// 	arg sqlc.GetCartItemByOwnerAndVariantParams,
// ) (sqlc.GetCartItemByOwnerAndVariantRow, error) {
// 	return r.q.GetCartItemByOwnerAndVariant(ctx, arg)
// }

// // 🔄 Update cart quantity for user or guest
// func (r *UpdateCartQuantityRepository) UpdateCartQuantity(
// 	ctx context.Context,
// 	arg sqlc.UpdateCartQuantityParams,
// ) error {
// 	return r.q.UpdateCartQuantity(ctx, arg)
// }

func (r *UpdateCartQuantityRepository) GetCartItemByUserAndVariant(
	ctx context.Context,
	userID uuid.UUID,
	variantID uuid.UUID,
) (sqlc.GetCartItemByUserAndVariantRow, error) {
	return r.q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
		UserID: sqlnull.UUID(userID), VariantID: variantID,
	})
}

func (r *UpdateCartQuantityRepository) GetCartItemByGuestAndVariant(
	ctx context.Context,
	guestUserID uuid.UUID,
	variantID uuid.UUID,
) (sqlc.GetCartItemByGuestAndVariantRow, error) {
	return r.q.GetCartItemByGuestAndVariant(ctx, sqlc.GetCartItemByGuestAndVariantParams{
		GuestUserID: sqlnull.UUID(guestUserID), VariantID: variantID,
	})
}

func (r *UpdateCartQuantityRepository) UpdateCartQuantityForUser(
	ctx context.Context,
	arg sqlc.UpdateCartQuantityForUserParams,
) error {
	return r.q.UpdateCartQuantityForUser(ctx, arg)
}

func (r *UpdateCartQuantityRepository) UpdateCartQuantityForGuest(
	ctx context.Context,
	arg sqlc.UpdateCartQuantityForGuestParams,
) error {
	return r.q.UpdateCartQuantityForGuest(ctx, arg)
}
