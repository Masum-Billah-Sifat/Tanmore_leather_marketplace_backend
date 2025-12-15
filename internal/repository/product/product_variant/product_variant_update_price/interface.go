// ------------------------------------------------------------
// 📁 File: internal/repository/product_variant/product_variant_update_price/interface.go
// 🧠 Repository interface for updating retail price of a variant.
//     Includes snapshot fetch, update mutation, and event insertion.

package product_variant_update_price

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantUpdatePriceRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of product + seller + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// 💵 Update retail price of a variant
	UpdateVariantRetailPrice(
		ctx context.Context,
		arg sqlc.UpdateVariantRetailPriceParams,
	) error

	// 📨 Insert event into events table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
