// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_disable_wholesale/interface.go
// 🧠 Repository interface for disabling wholesale mode on a variant.
//     Includes snapshot fetch, wholesale field reset, and event insertion.

package product_variant_disable_wholesale

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantDisableWholesaleRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of seller + product + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// ❌ Disable wholesale mode on variant (set fields to false/null)
	DisableWholesaleMode(
		ctx context.Context,
		arg sqlc.DisableWholesaleModeParams,
	) error

	// 📨 Insert event into outbox table
	InsertEvent(
		ctx context.Context,
		arg sqlc.InsertEventParams,
	) error
}
