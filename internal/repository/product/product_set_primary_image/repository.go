// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_set_primary_image/product_set_primary_image_repository.go
// 🧠 Concrete implementation of ProductSetPrimaryImageRepoInterface.
//     Handles seller check, product validation, image update, and emits event.

package product_set_primary_image

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ProductSetPrimaryImageRepository implements ProductSetPrimaryImageRepoInterface
type ProductSetPrimaryImageRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductSetPrimaryImageRepository(db *sql.DB) *ProductSetPrimaryImageRepository {
	return &ProductSetPrimaryImageRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductSetPrimaryImageRepository) WithTx(
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

// 👤 Fetch seller by ID
func (r *ProductSetPrimaryImageRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Fetch product by ID and seller ID
func (r *ProductSetPrimaryImageRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	arg sqlc.GetProductByIDAndSellerIDParams,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, arg)
}

// 🚫 Clear all previous primary images
func (r *ProductSetPrimaryImageRepository) UnsetAllPrimaryImages(
	ctx context.Context,
	arg sqlc.UnsetAllPrimaryImagesParams,
) error {
	return r.q.UnsetAllPrimaryImages(ctx, arg)
}

// ✅ Set selected image as primary
func (r *ProductSetPrimaryImageRepository) SetAsPrimaryImage(
	ctx context.Context,
	arg sqlc.SetAsPrimaryImageParams,
) error {
	return r.q.SetAsPrimaryImage(ctx, arg)
}

// 📨 Insert event
func (r *ProductSetPrimaryImageRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}

// 🖼️ Fetch media by ID and product ID
func (r *ProductSetPrimaryImageRepository) GetProductMediaByID(
	ctx context.Context,
	arg sqlc.GetProductMediaByIDParams,
) (sqlc.GetProductMediaByIDRow, error) {
	return r.q.GetProductMediaByID(ctx, arg)
}
