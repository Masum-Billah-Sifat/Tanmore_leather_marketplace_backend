// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_info/interface.go
// 🧠 Repository interface for adding a new variant to an existing product.
//     Contains only DB operations required by the service layer.

package product_variant_update_info

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantUpdateInfoRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of product + seller + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// ✏️ Update color and/or size of a variant
	UpdateVariantColorSize(
		ctx context.Context,
		arg sqlc.UpdateVariantColorSizeParams,
	) error

	// 📨 Insert event into events table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
