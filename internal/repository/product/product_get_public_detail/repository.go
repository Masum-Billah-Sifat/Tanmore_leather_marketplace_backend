// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_get_public_detail/product_get_public_detail_repository.go
// 🧠 Concrete implementation of ProductGetPublicDetailRepoInterface

package product_get_public_detail

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductGetPublicDetailRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductGetPublicDetailRepository(db *sql.DB) *ProductGetPublicDetailRepository {
	return &ProductGetPublicDetailRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Optional transaction wrapper
func (r *ProductGetPublicDetailRepository) WithTx(
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

// 🌐 Core product + variants fetch for customer view
func (r *ProductGetPublicDetailRepository) GetProductDetailByProductID(
	ctx context.Context,
	productID uuid.UUID,
) ([]sqlc.GetProductDetailByProductIDRow, error) {
	return r.q.GetProductDetailByProductID(ctx, productID)
}
