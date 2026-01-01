// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_archive/product_archive_repository.go
// 🧠 Concrete implementation of ProductArchiveRepoInterface.
//     Uses SQLC queries to handle seller validation, product check, archiving, and event insertion.

package product_archive

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ProductArchiveRepository implements ProductArchiveRepoInterface
type ProductArchiveRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductArchiveRepository(db *sql.DB) *ProductArchiveRepository {
	return &ProductArchiveRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductArchiveRepository) WithTx(
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

// 👤 Validate seller
func (r *ProductArchiveRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Confirm product ownership
func (r *ProductArchiveRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	productID, sellerID uuid.UUID,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
		ID:       productID,
		SellerID: sellerID,
	})
}

// 🧹 Soft-archive product
func (r *ProductArchiveRepository) ArchiveProduct(
	ctx context.Context,
	arg sqlc.ArchiveProductParams,
) error {
	return r.q.ArchiveProduct(ctx, arg)
}

// 📨 Insert product.archived event
func (r *ProductArchiveRepository) InsertEvent(
	ctx context.Context,
	arg sqlc.InsertEventParams,
) error {
	return r.q.InsertEvent(ctx, arg)
}
