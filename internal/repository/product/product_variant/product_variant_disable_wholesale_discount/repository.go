// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_remove_wholesale_discount/repository.go
// 🧠 Concrete implementation of ProductVariantRemoveWholesaleDiscountRepoInterface.
//     Uses SQLC to fetch snapshot, remove wholesale discount, and insert event.

package product_variant_remove_wholesale_discount

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantRemoveWholesaleDiscountRepository implements ProductVariantRemoveWholesaleDiscountRepoInterface
type ProductVariantRemoveWholesaleDiscountRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantRemoveWholesaleDiscountRepository(db *sql.DB) *ProductVariantRemoveWholesaleDiscountRepository {
	return &ProductVariantRemoveWholesaleDiscountRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantRemoveWholesaleDiscountRepository) WithTx(
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

// 🧠 Fetch snapshot
func (r *ProductVariantRemoveWholesaleDiscountRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ❌ Disable wholesale discount
func (r *ProductVariantRemoveWholesaleDiscountRepository) DisableWholesaleDiscount(
	ctx context.Context,
	arg sqlc.DisableWholesaleDiscountParams,
) error {
	return r.q.DisableWholesaleDiscount(ctx, arg)
}

// 📨 Insert outbox event
func (r *ProductVariantRemoveWholesaleDiscountRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
