// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_add_discount/interface.go
// 🧠 Repository interface for adding retail discount to a variant.
// Includes snapshot fetch, update discount mutation, and event insertion.

package product_variant_add_discount

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantAddDiscountRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of seller + product + variant + category
	GetVariantSnapshot(ctx context.Context, arg sqlc.GetVariantSnapshotParams) (sqlc.ProductVariantSnapshot, error)

	// 💸 Enable retail discount on variant
	EnableRetailDiscount(ctx context.Context, arg sqlc.EnableRetailDiscountParams) error

	// 📨 Insert event into outbox table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
