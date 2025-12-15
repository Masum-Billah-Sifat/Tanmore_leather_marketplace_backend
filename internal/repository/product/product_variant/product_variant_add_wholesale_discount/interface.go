// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_add_wholesale_discount/interface.go
// 🧠 Repository interface for adding wholesale discount to a variant.
//     Includes snapshot fetch, wholesale discount update, and event insertion.

package product_variant_add_wholesale_discount

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantAddWholesaleDiscountRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of seller + product + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// ➕ Add wholesale discount fields to variant
	EnableWholesaleDiscount(
		ctx context.Context,
		arg sqlc.EnableWholesaleDiscountParams,
	) error

	// 📨 Insert event into outbox table
	InsertEvent(
		ctx context.Context,
		arg sqlc.InsertEventParams,
	) error
}
