// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_edit_wholesale_info/repository.go
// 🧠 Concrete implementation of ProductVariantEditWholesaleInfoRepoInterface.
//     Uses SQLC queries to fetch snapshot, update wholesale info fields, and insert event.

package product_variant_update_wholesale_mode

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantEditWholesaleInfoRepository implements ProductVariantEditWholesaleInfoRepoInterface
type ProductVariantEditWholesaleInfoRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantEditWholesaleInfoRepository(db *sql.DB) *ProductVariantEditWholesaleInfoRepository {
	return &ProductVariantEditWholesaleInfoRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantEditWholesaleInfoRepository) WithTx(
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

// 🧠 Fetch snapshot for validation (seller/product/variant/category)
func (r *ProductVariantEditWholesaleInfoRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ✏️ Update wholesale fields using COALESCE
func (r *ProductVariantEditWholesaleInfoRepository) UpdateWholesaleInfo(
	ctx context.Context,
	arg sqlc.UpdateWholesaleModeParams,
) error {
	return r.q.UpdateWholesaleMode(ctx, arg)
}

// 📨 Insert event into outbox table
func (r *ProductVariantEditWholesaleInfoRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
