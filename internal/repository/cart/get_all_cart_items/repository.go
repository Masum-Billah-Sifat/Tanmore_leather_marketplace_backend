package get_all_cart_items

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

type GetAllCartItemsRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewGetAllCartItemsRepository(db *sql.DB) *GetAllCartItemsRepository {
	return &GetAllCartItemsRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper (future-proofed)
func (r *GetAllCartItemsRepository) WithTx(
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

// 🧑 User moderation check
func (r *GetAllCartItemsRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// // 🧾 Get all active variant IDs (supports user or guest)
// func (r *GetAllCartItemsRepository) ListActiveVariantIDsByOwner(
// 	ctx context.Context,
// 	arg sqlc.ListActiveVariantIDsByOwnerParams,
// ) ([]uuid.UUID, error) {
// 	return r.q.ListActiveVariantIDsByOwner(ctx, arg)
// }

// // 🔍 Fetch full cart + snapshot joined rows (supports user or guest)
// func (r *GetAllCartItemsRepository) GetActiveCartVariantSnapshotsByOwnerAndVariantIDs(
// 	ctx context.Context,
// 	arg sqlc.GetActiveCartVariantSnapshotsByOwnerAndVariantIDsParams,
// ) ([]sqlc.GetActiveCartVariantSnapshotsByOwnerAndVariantIDsRow, error) {
// 	return r.q.GetActiveCartVariantSnapshotsByOwnerAndVariantIDs(ctx, arg)
// }

func (r *GetAllCartItemsRepository) ListActiveVariantIDsByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {
	return r.q.ListActiveVariantIDsByUser(ctx, sqlnull.UUID(userID))
}

func (r *GetAllCartItemsRepository) ListActiveVariantIDsByGuest(
	ctx context.Context,
	guestUserID uuid.UUID,
) ([]uuid.UUID, error) {
	return r.q.ListActiveVariantIDsByGuest(ctx, sqlnull.UUID(guestUserID))
}

func (r *GetAllCartItemsRepository) GetActiveCartVariantSnapshotsByUserAndVariantIDs(
	ctx context.Context,
	arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error) {
	return r.q.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx, arg)
}

func (r *GetAllCartItemsRepository) GetActiveCartVariantSnapshotsByGuestAndVariantIDs(
	ctx context.Context,
	arg sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsParams,
) ([]sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsRow, error) {
	return r.q.GetActiveCartVariantSnapshotsByGuestAndVariantIDs(ctx, arg)
}
