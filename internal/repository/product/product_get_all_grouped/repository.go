// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_get_all_grouped/repository.go
// 🧠 Concrete implementation of ProductGetAllGroupedRepoInterface

package product_get_all_grouped

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductGetAllGroupedRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductGetAllGroupedRepository(db *sql.DB) *ProductGetAllGroupedRepository {
	return &ProductGetAllGroupedRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Optional transaction wrapper
func (r *ProductGetAllGroupedRepository) WithTx(
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

// 👤 Validate seller identity
func (r *ProductGetAllGroupedRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Fetch all variant indexes by seller
func (r *ProductGetAllGroupedRepository) GetAllProductVariantIndexesBySeller(
	ctx context.Context,
	sellerID uuid.UUID,
) ([]sqlc.ProductVariantIndex, error) {
	return r.q.GetAllProductVariantIndexesBySeller(ctx, sellerID)
}

// 🖼️ Get primary image
func (r *ProductGetAllGroupedRepository) GetPrimaryImageForProduct(
	ctx context.Context,
	arg sqlc.GetPrimaryProductImageByProductIDParams,
) (sqlc.ProductMedia, error) {
	return r.q.GetPrimaryProductImageByProductID(ctx, arg)
}
