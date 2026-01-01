// ------------------------------------------------------------
// 📁 File: internal/repository/category/category_tree/interface.go
// 🧠 Repository interface for fetching full category tree data.

package category_tree

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type CategoryTreeRepoInterface interface {
	// 🔁 Transaction wrapper (if needed in future extensions)
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🌲 Get all non-archived categories for tree building
	GetAllNonArchivedCategories(ctx context.Context) ([]sqlc.GetAllNonArchivedCategoriesRow, error)
}
