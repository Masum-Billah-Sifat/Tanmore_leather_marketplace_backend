// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_get_full_detail/product_get_full_detail_repository.go
// 🧠 Concrete implementation of ProductGetFullDetailRepoInterface

package product_get_full_detail

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductGetFullDetailRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductGetFullDetailRepository(db *sql.DB) *ProductGetFullDetailRepository {
	return &ProductGetFullDetailRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Optional transaction wrapper
func (r *ProductGetFullDetailRepository) WithTx(
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
func (r *ProductGetFullDetailRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Validate product ownership & moderation
func (r *ProductGetFullDetailRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	productID uuid.UUID,
	sellerID uuid.UUID,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
		ID:       productID,
		SellerID: sellerID,
	})
}

// // 🧩 Get all variant index rows for this product/seller
// func (r *ProductGetFullDetailRepository) GetVariantIndexesByProductAndSeller(
// 	ctx context.Context,
// 	arg sqlc.GetVariantIndexesByProductAndSellerParams,
// ) ([]sqlc.GetVariantIndexesByProductAndSellerParams, error) {
// 	return r.q.GetVariantIndexesByProductAndSeller(ctx, arg)
// }

// func (r *ProductGetFullDetailRepository) GetVariantIndexesByProductAndSeller(
// 	ctx context.Context,
// 	arg sqlc.GetVariantIndexesByProductAndSellerParams,
// ) ([]sqlc.GetVariantIndexesByProductAndSellerParams, error) {
// 	return r.q.GetVariantIndexesByProductAndSeller(ctx, arg)
// }

func (r *ProductGetFullDetailRepository) GetVariantIndexesByProductAndSeller(
	ctx context.Context,
	arg sqlc.GetVariantIndexesByProductAndSellerParams,
) ([]sqlc.ProductVariantIndex, error) {
	return r.q.GetVariantIndexesByProductAndSeller(ctx, arg)
}

// 🖼️ Get primary image
func (r *ProductGetFullDetailRepository) GetPrimaryImageForProduct(
	ctx context.Context,
	arg sqlc.GetPrimaryProductImageByProductIDParams,
) (sqlc.ProductMedia, error) {
	return r.q.GetPrimaryProductImageByProductID(ctx, arg)
}
