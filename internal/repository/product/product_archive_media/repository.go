// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_archive_media/product_archive_media_repository.go
// 🧠 Concrete implementation for archiving media (image or promo video) from a product.

package product_archive_media

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 Implements ProductArchiveMediaRepoInterface
type ProductArchiveMediaRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductArchiveMediaRepository(db *sql.DB) *ProductArchiveMediaRepository {
	return &ProductArchiveMediaRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *ProductArchiveMediaRepository) WithTx(
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

// 👤 Fetch seller
func (r *ProductArchiveMediaRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Fetch product ownership
func (r *ProductArchiveMediaRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	arg sqlc.GetProductByIDAndSellerIDParams,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, arg)
}

// // 🖼️ Get product media row
// func (r *ProductArchiveMediaRepository) GetProductMediaByID(
// 	ctx context.Context,
// 	arg sqlc.GetProductMediaByIDParams,
// ) (sqlc.ProductMedia, error) {
// 	return r.q.GetProductMediaByID(ctx, arg)
// }

// 🖼️ Get product media row
func (r *ProductArchiveMediaRepository) GetProductMediaByID(
	ctx context.Context,
	arg sqlc.GetProductMediaByIDParams,
) (sqlc.GetProductMediaByIDRow, error) {
	return r.q.GetProductMediaByID(ctx, arg)
}

// 🔢 Count active images
func (r *ProductArchiveMediaRepository) CountActiveImagesForProduct(
	ctx context.Context,
	arg sqlc.CountActiveImagesForProductParams,
) (int64, error) {
	return r.q.CountActiveImagesForProduct(ctx, arg)
}

// 🗑️ Archive media
func (r *ProductArchiveMediaRepository) ArchiveProductMedia(
	ctx context.Context,
	arg sqlc.ArchiveProductMediaParams,
) error {
	return r.q.ArchiveProductMedia(ctx, arg)
}

// 📨 Emit event
func (r *ProductArchiveMediaRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
