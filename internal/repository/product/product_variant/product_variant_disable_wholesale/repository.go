// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_disable_wholesale/repository.go
// 🧠 Concrete implementation of ProductVariantDisableWholesaleRepoInterface.
//     Uses SQLC to fetch variant snapshot, disable wholesale mode, and insert event.

package product_variant_disable_wholesale

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantDisableWholesaleRepository implements ProductVariantDisableWholesaleRepoInterface
type ProductVariantDisableWholesaleRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantDisableWholesaleRepository(db *sql.DB) *ProductVariantDisableWholesaleRepository {
	return &ProductVariantDisableWholesaleRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantDisableWholesaleRepository) WithTx(
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

// 🧠 Fetch snapshot for validation
func (r *ProductVariantDisableWholesaleRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ❌ Disable wholesale mode (reset fields)
func (r *ProductVariantDisableWholesaleRepository) DisableWholesaleMode(
	ctx context.Context,
	arg sqlc.DisableWholesaleModeParams,
) error {
	return r.q.DisableWholesaleMode(ctx, arg)
}

// 📨 Insert event into outbox
func (r *ProductVariantDisableWholesaleRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
