// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_in_stock/interface.go
// 🧠 Repository interface for updating in_stock status of a variant.
//     Includes snapshot fetch, mutation, and event insertion.

package product_variant_update_in_stock

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantUpdateInStockRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of product + seller + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// ✅ Update in_stock field of variant
	UpdateVariantInStock(
		ctx context.Context,
		arg sqlc.UpdateVariantInStockParams,
	) error

	// 📨 Insert event into events table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
