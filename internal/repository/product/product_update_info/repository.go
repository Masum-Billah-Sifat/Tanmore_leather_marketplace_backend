// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_update_info/product_update_info_repository.go
// 🧠 Concrete implementation of ProductUpdateInfoRepoInterface
//     using SQLC to fetch, update, and emit product info update events.

package product_update_info

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ProductUpdateInfoRepository implements ProductUpdateInfoRepoInterface
type ProductUpdateInfoRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductUpdateInfoRepository(db *sql.DB) *ProductUpdateInfoRepository {
	return &ProductUpdateInfoRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductUpdateInfoRepository) WithTx(
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

// 🔍 Fetch user by ID
func (r *ProductUpdateInfoRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 🔍 Fetch product by ID and seller ID
func (r *ProductUpdateInfoRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	arg sqlc.GetProductByIDAndSellerIDParams,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, arg)
}

// 📝 Update product title/description
func (r *ProductUpdateInfoRepository) UpdateProductTitleDesc(
	ctx context.Context,
	arg sqlc.UpdateProductTitleDescParams,
) error {
	return r.q.UpdateProductTitleDesc(ctx, arg)
}

// 📨 Insert event into outbox
func (r *ProductUpdateInfoRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
