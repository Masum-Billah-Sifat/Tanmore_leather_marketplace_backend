// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_price/product_variant_repository.go
// 🧠 Concrete implementation of ProductVariantUpdatePriceRepoInterface
//     Uses SQLC queries to fetch snapshot, update retail price, and insert event.

package product_variant_update_price

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantUpdatePriceRepository implements ProductVariantUpdatePriceRepoInterface
type ProductVariantUpdatePriceRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantUpdatePriceRepository(db *sql.DB) *ProductVariantUpdatePriceRepository {
	return &ProductVariantUpdatePriceRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantUpdatePriceRepository) WithTx(
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
func (r *ProductVariantUpdatePriceRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// 💵 Update variant retail price
func (r *ProductVariantUpdatePriceRepository) UpdateVariantRetailPrice(
	ctx context.Context,
	arg sqlc.UpdateVariantRetailPriceParams,
) error {
	return r.q.UpdateVariantRetailPrice(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantUpdatePriceRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
