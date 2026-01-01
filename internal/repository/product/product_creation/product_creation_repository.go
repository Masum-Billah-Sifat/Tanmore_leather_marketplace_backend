// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_repository.go
// 🧠 This file provides the concrete implementation of ProductRepoInterface
//     using SQLC-generated queries, following Meta-grade patterns.

package product_creation

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ProductRepository implements ProductRepoInterface
type ProductRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor for ProductRepository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductRepository) WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
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

// 👤 Fetch user by ID (seller validation)
func (r *ProductRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, userID)
}

// 📦 Insert product row
func (r *ProductRepository) InsertProduct(
	ctx context.Context,
	arg sqlc.InsertProductParams,
) (uuid.UUID, error) {
	return r.q.InsertProduct(ctx, arg)
}

// 🧩 Insert product variant and return variant ID
func (r *ProductRepository) InsertProductVariantReturningID(
	ctx context.Context,
	arg sqlc.InsertProductVariantReturningIDParams,
) (uuid.UUID, error) {
	return r.q.InsertProductVariantReturningID(ctx, arg)
}

// 📨 Insert event into events table (outbox pattern)
func (r *ProductRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}

// 🧠 Fetch seller profile metadata (used in event payload)
func (r *ProductRepository) GetSellerProfileMetadataBySellerID(
	ctx context.Context,
	sellerID uuid.UUID,
) (sqlc.SellerProfileMetadatum, error) {
	return r.q.GetSellerProfileMetadataBySellerID(ctx, sellerID)
}

// 🗂️ Fetch category by ID (for validation)
func (r *ProductRepository) GetCategoryByID(
	ctx context.Context,
	categoryID uuid.UUID,
) (sqlc.Category, error) {
	return r.q.GetCategoryByID(ctx, categoryID)
}

// 🖼️ Insert media into product_medias (image/video)
func (r *ProductRepository) InsertProductMedia(
	ctx context.Context,
	arg sqlc.InsertProductMediaParams,
) (uuid.UUID, error) {
	return r.q.InsertProductMedia(ctx, arg)
}
