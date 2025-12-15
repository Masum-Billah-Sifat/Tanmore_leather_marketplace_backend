// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_discount/product_variant_update_discount_repository.go
// 🧠 Concrete implementation of ProductVariantUpdateDiscountRepoInterface
//     Uses SQLC queries to fetch snapshot, update discount fields, and insert event.

package product_variant_update_discount

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantUpdateDiscountRepository implements ProductVariantUpdateDiscountRepoInterface
type ProductVariantUpdateDiscountRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantUpdateDiscountRepository(db *sql.DB) *ProductVariantUpdateDiscountRepository {
	return &ProductVariantUpdateDiscountRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantUpdateDiscountRepository) WithTx(
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

// 🧠 Fetch snapshot to validate seller/product/variant/category
func (r *ProductVariantUpdateDiscountRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ✏️ Update variant retail discount using COALESCE
func (r *ProductVariantUpdateDiscountRepository) UpdateRetailDiscount(
	ctx context.Context,
	arg sqlc.UpdateRetailDiscountParams,
) error {
	return r.q.UpdateRetailDiscount(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantUpdateDiscountRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
