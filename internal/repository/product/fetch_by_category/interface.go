// ------------------------------------------------------------
// 📁 File: internal/repository/product/fetch_by_category/interface.go
// 🧠 Repository interface for fetching products by category.
//     Includes category validation, leaf discovery, and product variant index fetch.

package fetch_by_category

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type FetchByCategoryRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🔍 Category validation
	GetCategoryByID(ctx context.Context, id uuid.UUID) (sqlc.Category, error)

	// 🌿 Leaf category discovery
	GetAllLeafCategoryIDsByRoot(ctx context.Context, rootID uuid.UUID) ([]uuid.UUID, error)

	// 🛍️ Fetch valid product variant index rows
	GetProductVariantIndexesByCategoryIDs(ctx context.Context, ids []uuid.UUID) ([]sqlc.GetProductVariantIndexesByCategoryIDsRow, error)
}
