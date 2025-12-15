// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_add_wholesale_discount/repository.go
// 🧠 Concrete implementation of ProductVariantAddWholesaleDiscountRepoInterface.
//     Uses SQLC to fetch variant snapshot, add wholesale discount, and insert event.

package product_variant_add_wholesale_discount

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantAddWholesaleDiscountRepository implements ProductVariantAddWholesaleDiscountRepoInterface
type ProductVariantAddWholesaleDiscountRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantAddWholesaleDiscountRepository(db *sql.DB) *ProductVariantAddWholesaleDiscountRepository {
	return &ProductVariantAddWholesaleDiscountRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantAddWholesaleDiscountRepository) WithTx(
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
func (r *ProductVariantAddWholesaleDiscountRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ➕ Add wholesale discount to variant
func (r *ProductVariantAddWholesaleDiscountRepository) EnableWholesaleDiscount(
	ctx context.Context,
	arg sqlc.EnableWholesaleDiscountParams,
) error {
	return r.q.EnableWholesaleDiscount(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantAddWholesaleDiscountRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
