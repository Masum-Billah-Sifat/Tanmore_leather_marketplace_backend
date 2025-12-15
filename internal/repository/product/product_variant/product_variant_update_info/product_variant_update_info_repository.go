// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_info/product_variant_update_info_repository.go
// 🧠 Concrete implementation of ProductVariantRepoInterface
//     using SQLC-generated queries.

package product_variant_update_info

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantRepository implements ProductVariantRepoInterface
type ProductVariantUpdateInfoRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantUpdateInfoRepository(db *sql.DB) *ProductVariantUpdateInfoRepository {
	return &ProductVariantUpdateInfoRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductVariantUpdateInfoRepository) WithTx(
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

// 🧠 Fetch snapshot to validate variant-level ownership
func (r *ProductVariantUpdateInfoRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ✏️ Update variant color/size with COALESCE strategy
func (r *ProductVariantUpdateInfoRepository) UpdateVariantColorSize(
	ctx context.Context,
	arg sqlc.UpdateVariantColorSizeParams,
) error {
	return r.q.UpdateVariantColorSize(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantUpdateInfoRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
