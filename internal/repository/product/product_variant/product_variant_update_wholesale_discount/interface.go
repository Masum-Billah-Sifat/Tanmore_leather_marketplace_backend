// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_wholesale_discount/interface.go
// 🧠 Repository interface for updating wholesale discount of a variant.
//     Includes snapshot fetch, COALESCE-based update, and event insertion.

package product_variant_update_wholesale_discount

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantUpdateWholesaleDiscountRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of seller + product + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// ✏️ Update wholesale discount using COALESCE
	UpdateWholesaleDiscount(
		ctx context.Context,
		arg sqlc.UpdateWholesaleDiscountParams,
	) error

	// 📨 Insert event into outbox table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
