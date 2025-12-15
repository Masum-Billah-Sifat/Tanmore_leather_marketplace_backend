// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_weight/product_variant_repository.go
// 🧠 Concrete implementation of ProductVariantUpdateWeightRepoInterface
//     Uses SQLC queries to fetch snapshot, update weight, and insert event.

package product_variant_update_weight

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 ProductVariantUpdateWeightRepository implements ProductVariantUpdateWeightRepoInterface
type ProductVariantUpdateWeightRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantUpdateWeightRepository(db *sql.DB) *ProductVariantUpdateWeightRepository {
	return &ProductVariantUpdateWeightRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductVariantUpdateWeightRepository) WithTx(
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
func (r *ProductVariantUpdateWeightRepository) GetVariantSnapshot(
	ctx context.Context,
	arg sqlc.GetVariantSnapshotParams,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, arg)
}

// ⚖️ Update weight in grams
func (r *ProductVariantUpdateWeightRepository) UpdateVariantWeight(
	ctx context.Context,
	arg sqlc.UpdateVariantWeightParams,
) error {
	return r.q.UpdateVariantWeight(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantUpdateWeightRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
