// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_add_media/product_add_media_repository.go
// 🧠 Concrete implementation of ProductAddMediaRepoInterface.
//     Uses SQLC queries to validate user/product, insert media, and emit event.

package product_add_media

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ProductAddMediaRepository implements ProductAddMediaRepoInterface
type ProductAddMediaRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductAddMediaRepository(db *sql.DB) *ProductAddMediaRepository {
	return &ProductAddMediaRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductAddMediaRepository) WithTx(
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
func (r *ProductAddMediaRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Fetch product by ID and seller ID
func (r *ProductAddMediaRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	arg sqlc.GetProductByIDAndSellerIDParams,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, arg)
}

// 📽️ Check if promo video exists
func (r *ProductAddMediaRepository) GetPromoVideoByProductID(
	ctx context.Context,
	arg sqlc.GetPromoVideoByProductIDParams,
) (uuid.UUID, error) {
	return r.q.GetPromoVideoByProductID(ctx, arg)
}

// 📥 Insert product media
func (r *ProductAddMediaRepository) InsertProductMedia(
	ctx context.Context,
	arg sqlc.InsertProductMediaParams,
) (uuid.UUID, error) {
	return r.q.InsertProductMedia(ctx, arg)
}

// 📨 Insert event
func (r *ProductAddMediaRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
