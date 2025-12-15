// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_stock_quantity/interface.go
// 🧠 Repository interface for updating stock quantity of a variant.
//     Includes snapshot fetch, update mutation, and event insertion.

package product_variant_update_stock_quantity

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantUpdateStockQuantityRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of product + seller + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// 📦 Update stock quantity of a variant
	UpdateVariantStockQuantity(
		ctx context.Context,
		arg sqlc.UpdateVariantStockQuantityParams,
	) error

	// 📨 Insert event into events table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
