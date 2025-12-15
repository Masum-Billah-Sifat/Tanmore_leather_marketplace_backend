// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_enable_wholesale/interface.go
// 🧠 Repository interface for enabling wholesale mode on a variant.
//     Includes snapshot fetch, wholesale mode update, and event insertion.

package product_variant_enable_wholesale_mode

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantEnableWholesaleRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of seller + product + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// 🏷️ Enable wholesale mode on variant
	EnableWholesaleMode(
		ctx context.Context,
		arg sqlc.EnableWholesaleModeParams,
	) error

	// 📨 Insert event into outbox table
	InsertEvent(
		ctx context.Context,
		arg sqlc.InsertEventParams,
	) error
}
