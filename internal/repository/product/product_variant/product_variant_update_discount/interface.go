// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_discount/interface.go
// 🧠 Repository interface for updating retail discount of a variant.
//     Includes snapshot fetch, discount update (with COALESCE), and event insertion.

package product_variant_update_discount

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantUpdateDiscountRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of product + seller + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// ✏️ Update retail discount using COALESCE
	UpdateRetailDiscount(
		ctx context.Context,
		arg sqlc.UpdateRetailDiscountParams,
	) error

	// 📨 Insert event into events table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
