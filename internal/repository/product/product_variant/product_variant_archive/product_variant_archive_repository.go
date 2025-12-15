// ------------------------------------------------------------
// 📁 File: internal/repository/product_variant/product_variant_repository.go
// 🧠 Implementation of ProductVariantRepoInterface for variant archival.
//     Uses SQLC-generated queries for DB operations.

package product_variant_archive

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ProductVariantRepository implements ProductVariantRepoInterface
type ProductVariantArchiveRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantArchiveRepository(db *sql.DB) *ProductVariantArchiveRepository {
	return &ProductVariantArchiveRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductVariantArchiveRepository) WithTx(
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

// 📄 Fetch full variant snapshot
func (r *ProductVariantArchiveRepository) GetVariantSnapshot(
	ctx context.Context,
	sellerID, productID, variantID uuid.UUID,
) (sqlc.ProductVariantSnapshot, error) {
	return r.q.GetVariantSnapshot(ctx, sqlc.GetVariantSnapshotParams{
		Sellerid:  sellerID,
		Productid: productID,
		Variantid: variantID,
	})
}

// 🧹 Soft-delete variant by setting is_archived = true
func (r *ProductVariantArchiveRepository) ArchiveProductVariant(
	ctx context.Context,
	arg sqlc.ArchiveProductVariantParams,
) error {
	return r.q.ArchiveProductVariant(ctx, arg)
}

// 📨 Insert archival event
func (r *ProductVariantArchiveRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
