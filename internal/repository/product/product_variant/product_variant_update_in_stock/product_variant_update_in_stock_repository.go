// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_in_stock/product_variant_repository.go
// 🧠 Concrete implementation of ProductVariantUpdateInStockRepoInterface
//     Uses SQLC to snapshot, update in_stock field, and insert event.

package product_variant_update_in_stock

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantUpdateInStockRepository implements ProductVariantUpdateInStockRepoInterface
type ProductVariantUpdateInStockRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantUpdateInStockRepository(db *sql.DB) *ProductVariantUpdateInStockRepository {
	return &ProductVariantUpdateInStockRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantUpdateInStockRepository) WithTx(
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

// 🧠 Fetch snapshot of seller + product + variant + category
func (r *ProductVariantUpdateInStockRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ✅ Update in_stock field of variant
func (r *ProductVariantUpdateInStockRepository) UpdateVariantInStock(
	ctx context.Context,
	arg sqlc.UpdateVariantInStockParams,
) error {
	return r.q.UpdateVariantInStock(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantUpdateInStockRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
