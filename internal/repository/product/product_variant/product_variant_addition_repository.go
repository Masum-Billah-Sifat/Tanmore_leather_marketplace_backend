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
