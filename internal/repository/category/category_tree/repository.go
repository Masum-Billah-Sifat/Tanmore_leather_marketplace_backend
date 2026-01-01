// ------------------------------------------------------------
// 📁 File: internal/repository/category/category_tree/category_tree_repository.go
// 🧠 Concrete implementation of CategoryTreeRepoInterface.
//     Provides method to fetch non-archived categories in tree format.

package category_tree

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
)

// 📦 CategoryTreeRepository implements CategoryTreeRepoInterface
type CategoryTreeRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewCategoryTreeRepository(db *sql.DB) *CategoryTreeRepository {
	return &CategoryTreeRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper (for future-proofing)
func (r *CategoryTreeRepository) WithTx(
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

// 🌲 Get all non-archived categories
func (r *CategoryTreeRepository) GetAllNonArchivedCategories(
	ctx context.Context,
) ([]sqlc.GetAllNonArchivedCategoriesRow, error) {
	return r.q.GetAllNonArchivedCategories(ctx)
}
