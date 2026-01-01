// ------------------------------------------------------------
// 📁 File: internal/repository/product_variant/product_variant_repository.go
// 🧠 Concrete implementation of ProductVariantRepoInterface
//     using SQLC-generated queries.

package product_variant

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ProductVariantRepository implements ProductVariantRepoInterface
type ProductVariantRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductVariantRepository(db *sql.DB) *ProductVariantRepository {
	return &ProductVariantRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductVariantRepository) WithTx(
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

// 👤 Fetch user by ID
func (r *ProductVariantRepository) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, userID)
}

// 📦 Fetch product by ID and seller ID (ownership check)
func (r *ProductVariantRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	arg sqlc.GetProductByIDAndSellerIDParams,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, arg)
}

// 🧩 Insert product variant and return ID
func (r *ProductVariantRepository) InsertProductVariantReturningID(
	ctx context.Context,
	arg sqlc.InsertProductVariantReturningIDParams,
) (uuid.UUID, error) {
	return r.q.InsertProductVariantReturningID(ctx, arg)
}

// 📨 Insert event (outbox pattern)
func (r *ProductVariantRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}

// 🧠 Fetch category info by ID
func (r *ProductVariantRepository) GetCategoryByID(
	ctx context.Context,
	categoryID uuid.UUID,
) (sqlc.Category, error) {
	return r.q.GetCategoryByID(ctx, categoryID)
}

// 🧠 Fetch seller profile metadata
func (r *ProductVariantRepository) GetSellerProfileMetadataBySellerID(
	ctx context.Context,
	sellerID uuid.UUID,
) (sqlc.SellerProfileMetadatum, error) {
	return r.q.GetSellerProfileMetadataBySellerID(ctx, sellerID)
}

// 🖼️ Get all active medias for product
func (r *ProductVariantRepository) GetActiveMediasByProductID(
	ctx context.Context,
	arg sqlc.GetActiveMediasByProductIDParams,
) ([]sqlc.ProductMedia, error) {
	return r.q.GetActiveMediasByProductID(ctx, arg)
}

// 🖼️ Get primary image for product
func (r *ProductVariantRepository) GetPrimaryProductImageByProductID(
	ctx context.Context,
	arg sqlc.GetPrimaryProductImageByProductIDParams,
) (sqlc.ProductMedia, error) {
	return r.q.GetPrimaryProductImageByProductID(ctx, arg)
}
