// ------------------------------------------------------------
// 📁 File: internal/repository/product_variant/interface.go
// 🧠 Repository interface for product variant archival.
//     Defines methods required for soft-deletion and validation.

package product_variant_archive

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductVariantArchiveRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 📄 Fetch variant snapshot with all validations (ownership, status)
	GetVariantSnapshot(
		ctx context.Context,
		sellerID, productID, variantID uuid.UUID,
	) (sqlc.ProductVariantSnapshot, error)

	// 🧹 Soft-delete the variant by marking is_archived = true
	ArchiveProductVariant(
		ctx context.Context,
		arg sqlc.ArchiveProductVariantParams,
	) error

	// 📨 Insert outbox event for archiving
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
