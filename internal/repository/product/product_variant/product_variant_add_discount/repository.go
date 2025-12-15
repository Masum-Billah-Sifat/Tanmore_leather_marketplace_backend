// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_add_discount/repository.go
// 🧠 Concrete implementation of ProductVariantAddDiscountRepoInterface.
// Uses SQLC to fetch variant snapshot, apply discount update, and insert event.

package product_variant_add_discount

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📆 ProductVariantAddDiscountRepository implements ProductVariantAddDiscountRepoInterface
type ProductVariantAddDiscountRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantAddDiscountRepository(db *sql.DB) *ProductVariantAddDiscountRepository {
	return &ProductVariantAddDiscountRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantAddDiscountRepository) WithTx(
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
func (r *ProductVariantAddDiscountRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// 💸 Enable retail discount
func (r *ProductVariantAddDiscountRepository) EnableRetailDiscount(
	ctx context.Context,
	arg sqlc.EnableRetailDiscountParams,
) error {
	return r.q.EnableRetailDiscount(ctx, arg)
}

// 📨 Insert outbox event
func (r *ProductVariantAddDiscountRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
