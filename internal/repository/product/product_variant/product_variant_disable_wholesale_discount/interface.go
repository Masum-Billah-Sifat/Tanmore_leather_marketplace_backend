// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_remove_wholesale_discount/interface.go
// 🧠 Repository interface for disabling/removing wholesale discount from a variant.
//     Includes snapshot fetch, discount field removal, and event insertion.

package product_variant_remove_wholesale_discount

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantRemoveWholesaleDiscountRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of seller + product + variant + category
	GetVariantSnapshot(ctx context.Context, arg sqlc.GetVariantSnapshotParams) (sqlc.ProductVariantSnapshot, error)

	// ❌ Disable wholesale discount on variant
	DisableWholesaleDiscount(ctx context.Context, arg sqlc.DisableWholesaleDiscountParams) error

	// 📨 Insert event into outbox table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
