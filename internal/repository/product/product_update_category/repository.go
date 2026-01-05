// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_update_category/product_update_category_repository.go
// 🧠 Concrete implementation of ProductUpdateCategoryRepoInterface

package product_update_category

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductUpdateCategoryRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductUpdateCategoryRepository(db *sql.DB) *ProductUpdateCategoryRepository {
	return &ProductUpdateCategoryRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductUpdateCategoryRepository) WithTx(
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
func (r *ProductUpdateCategoryRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Validate product ownership + moderation
func (r *ProductUpdateCategoryRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	arg sqlc.GetProductByIDAndSellerIDParams,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, arg)
}

// 🌳 Validate new category
func (r *ProductUpdateCategoryRepository) GetCategoryByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.Category, error) {
	return r.q.GetCategoryByID(ctx, id)
}

// ✏️ Update product category
func (r *ProductUpdateCategoryRepository) UpdateProductCategory(
	ctx context.Context,
	arg sqlc.UpdateProductCategoryParams,
) error {
	return r.q.UpdateProductCategory(ctx, arg)
}

// 📨 Insert category update event
func (r *ProductUpdateCategoryRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
