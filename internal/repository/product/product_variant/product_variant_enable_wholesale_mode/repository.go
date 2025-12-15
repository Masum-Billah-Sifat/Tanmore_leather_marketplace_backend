// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_enable_wholesale/repository.go
// 🧠 Concrete implementation of ProductVariantEnableWholesaleRepoInterface.
//     Uses SQLC to fetch variant snapshot, enable wholesale mode, and insert event.

package product_variant_enable_wholesale_mode

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantEnableWholesaleRepository implements ProductVariantEnableWholesaleRepoInterface
type ProductVariantEnableWholesaleRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantEnableWholesaleRepository(db *sql.DB) *ProductVariantEnableWholesaleRepository {
	return &ProductVariantEnableWholesaleRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantEnableWholesaleRepository) WithTx(
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

// 🧠 Fetch snapshot to validate ownership and moderation
func (r *ProductVariantEnableWholesaleRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// 🏷️ Enable wholesale mode
func (r *ProductVariantEnableWholesaleRepository) EnableWholesaleMode(
	ctx context.Context,
	arg sqlc.EnableWholesaleModeParams,
) error {
	return r.q.EnableWholesaleMode(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantEnableWholesaleRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
