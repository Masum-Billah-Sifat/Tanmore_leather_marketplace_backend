// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_remove_discount/repository.go
// 🧠 Concrete implementation of ProductVariantRemoveDiscountRepoInterface.
// Uses SQLC to fetch variant snapshot, remove discount fields, and insert event.

package product_variant_remove_retail_discount

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📆 ProductVariantRemoveDiscountRepository implements ProductVariantRemoveDiscountRepoInterface
type ProductVariantRemoveDiscountRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantRemoveDiscountRepository(db *sql.DB) *ProductVariantRemoveDiscountRepository {
	return &ProductVariantRemoveDiscountRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantRemoveDiscountRepository) WithTx(
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
func (r *ProductVariantRemoveDiscountRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ❌ Disable retail discount
func (r *ProductVariantRemoveDiscountRepository) DisableRetailDiscount(
	ctx context.Context,
	arg sqlc.DisableRetailDiscountParams,
) error {
	return r.q.DisableRetailDiscount(ctx, arg)
}

// 📨 Insert outbox event
func (r *ProductVariantRemoveDiscountRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
