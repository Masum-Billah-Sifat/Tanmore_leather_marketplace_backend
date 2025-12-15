// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_stock_quantity/product_variant_repository.go
// 🧠 Concrete implementation of ProductVariantUpdateStockQuantityRepoInterface
//     Uses SQLC queries to fetch snapshot, update stock quantity, and insert event.

package product_variant_update_stock_quantity

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantUpdateStockQuantityRepository implements ProductVariantUpdateStockQuantityRepoInterface
type ProductVariantUpdateStockQuantityRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantUpdateStockQuantityRepository(db *sql.DB) *ProductVariantUpdateStockQuantityRepository {
	return &ProductVariantUpdateStockQuantityRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantUpdateStockQuantityRepository) WithTx(
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

// 🧠 Fetch snapshot to validate ownership and access
func (r *ProductVariantUpdateStockQuantityRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// 📦 Update variant stock quantity
func (r *ProductVariantUpdateStockQuantityRepository) UpdateVariantStockQuantity(
	ctx context.Context,
	arg sqlc.UpdateVariantStockQuantityParams,
) error {
	return r.q.UpdateVariantStockQuantity(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantUpdateStockQuantityRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
