// ------------------------------------------------------------
// 📁 File: internal/repository/cart/add_to_cart/add_to_cart_repository.go
// 🧠 Concrete implementation of AddToCartRepoInterface.
//     Performs user checks, snapshot reads, cart logic, and event insertions.

package add_to_cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// 📦 AddToCartRepository implements AddToCartRepoInterface
type AddToCartRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewAddToCartRepository(db *sql.DB) *AddToCartRepository {
	return &AddToCartRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *AddToCartRepository) WithTx(
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
func (r *AddToCartRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Fetch variant snapshot by product + variant
func (r *AddToCartRepository) GetVariantSnapshotByProductIDAndVariantID(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotByProductIDAndVariantIDParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshotByProductIDAndVariantID(ctx, arg)
}

// // 🛒 Find cart item by owner (user_id or guest_user_id) and variant
// func (r *AddToCartRepository) GetCartItemByOwnerAndVariant(
// 	ctx context.Context,
// 	arg sqlc.GetCartItemByOwnerAndVariantParams,
// ) (sqlc.GetCartItemByOwnerAndVariantRow, error) {
// 	return r.q.GetCartItemByOwnerAndVariant(ctx, arg)
// }

func (r *AddToCartRepository) GetCartItemByUserAndVariant(
	ctx context.Context,
	userID uuid.UUID,
	variantID uuid.UUID,
) (sqlc.GetCartItemByUserAndVariantRow, error) {
	return r.q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
		UserID:    sqlnull.UUID(userID),
		VariantID: variantID,
	})
}

func (r *AddToCartRepository) GetCartItemByGuestAndVariant(
	ctx context.Context,
	guestUserID uuid.UUID,
	variantID uuid.UUID,
) (sqlc.GetCartItemByGuestAndVariantRow, error) {
	return r.q.GetCartItemByGuestAndVariant(ctx, sqlc.GetCartItemByGuestAndVariantParams{
		GuestUserID: sqlnull.UUID(guestUserID),
		VariantID:   variantID,
	})
}

// ♻️ Reactivate an inactive cart item
func (r *AddToCartRepository) ReactivateCartItemByID(
	ctx context.Context,
	arg sqlc.ReactivateCartItemByIDParams,
) error {
	return r.q.ReactivateCartItemByID(ctx, arg)
}

// ➕ Insert a new cart item
func (r *AddToCartRepository) InsertCartItem(
	ctx context.Context,
	arg sqlc.InsertCartItemParams,
) (sqlc.InsertCartItemRow, error) {
	return r.q.InsertCartItem(ctx, arg)
}
