// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_variant/product_variant_edit_wholesale_info/interface.go
// 🧠 Repository interface for editing wholesale info of a variant.
//     Includes snapshot fetch, wholesale update (with COALESCE), and event insertion.

package product_variant_update_wholesale_mode

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
)

type ProductVariantEditWholesaleInfoRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🧠 Fetch snapshot of seller + product + variant + category
	GetVariantSnapshot(
		ctx context.Context,
		arg sqlc.GetVariantSnapshotParams,
	) (sqlc.ProductVariantSnapshot, error)

	// ✏️ Update wholesale fields using COALESCE
	UpdateWholesaleInfo(
		ctx context.Context,
		arg sqlc.UpdateWholesaleModeParams,
	) error

	// 📨 Insert event into outbox table
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
