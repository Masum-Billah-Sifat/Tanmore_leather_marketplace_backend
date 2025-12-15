// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_update_weight/interface.go
// 🧠 Repository interface for updating weight (grams) of a variant.
//     Includes snapshot fetch, update mutation, and event insertion.

package product_variant_update_weight

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantUpdateWeightRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of product + seller + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// ⚖️ Update weight (grams) of the variant
	UpdateVariantWeight(
		ctx context.Context,
		arg sqlc.UpdateVariantWeightParams,
	) error

	// 📨 Insert event into events table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
