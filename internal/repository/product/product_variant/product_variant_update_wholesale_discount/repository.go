// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_wholesale_discount/repository.go
// 🧠 Concrete implementation of ProductVariantUpdateWholesaleDiscountRepoInterface.
//     Uses SQLC to fetch snapshot, update discount using COALESCE, and insert event.

package product_variant_update_wholesale_discount

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantUpdateWholesaleDiscountRepository implements ProductVariantUpdateWholesaleDiscountRepoInterface
type ProductVariantUpdateWholesaleDiscountRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantUpdateWholesaleDiscountRepository(db *sql.DB) *ProductVariantUpdateWholesaleDiscountRepository {
	return &ProductVariantUpdateWholesaleDiscountRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantUpdateWholesaleDiscountRepository) WithTx(
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
func (r *ProductVariantUpdateWholesaleDiscountRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ✏️ Update wholesale discount using COALESCE
func (r *ProductVariantUpdateWholesaleDiscountRepository) UpdateWholesaleDiscount(
	ctx context.Context,
	arg sqlc.UpdateWholesaleDiscountParams,
) error {
	return r.q.UpdateWholesaleDiscount(ctx, arg)
}

// 📨 Insert event into outbox
func (r *ProductVariantUpdateWholesaleDiscountRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
