// internal/jobs/engine/router/router.go
package router

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/internal/jobs/product_variant_index"
	"tanmore_backend/internal/jobs/product_variant_snapshot"
)

func Route(ctx context.Context, event sqlc.Event, q *sqlc.Queries) error {
	switch event.EventType {

	// Product Events
	case "product.created",
		"product.info.updated",
		"product.image.added",
		"product.image.removed",
		"product.image.set_primary",
		"product.promo_video.added",
		"product.promo_video.removed",
		"product.category_updated",
		"product.archived":
		if err := product_variant_snapshot.ProcessSnapshotEvent(ctx, event, q); err != nil {
			return err
		}
		if err := product_variant_index.ProcessIndexEvent(ctx, event, q); err != nil {
			return err
		}

	// Variant Events
	case "variant.created",
		"variant.archived",
		"variant.info.updated",
		"variant.in_stock.updated",
		"variant.stock_quantity.updated",
		"variant.retail_price.updated",
		"variant.retail_discount.added",
		"variant.retail_discount.updated",
		"variant.retail_discount.removed",
		"variant.wholesale_mode.enabled",
		"variant.wholesale_mode.updated",
		"variant.wholesale_mode.disabled",
		"variant.wholesale_discount.added",
		"variant.wholesale_discount.updated",
		"variant.wholesale_discount.removed":
		if err := product_variant_snapshot.ProcessSnapshotEvent(ctx, event, q); err != nil {
			return err
		}
		if err := product_variant_index.ProcessIndexEvent(ctx, event, q); err != nil {
			return err
		}
	}

	return nil
}

// // internal/jobs/engine/router/router.go
// package router

// import (
// 	"context"

// 	"tanmore_backend/internal/db/sqlc"
// 	"tanmore_backend/internal/jobs/product_variant_index"
// 	"tanmore_backend/internal/jobs/product_variant_snapshot"
// )

// func Route(ctx context.Context, event sqlc.Event, q *sqlc.Queries) error {
// 	switch event.EventType {

// 	// Product Events
// 	case "product.created",
// 		"product.info.updated",
// 		"product.image.added",
// 		"product.image.removed",
// 		"product.image.set_primary",
// 		"product.promo_video.added",
// 		"product.promo_video.removed":
// 		if err := product_variant_snapshot.ProcessSnapshotEvent(ctx, event, q); err != nil {
// 			return err
// 		}
// 		if err := product_variant_index.ProcessIndexEvent(ctx, event, q); err != nil {
// 			return err
// 		}

// 	// Variant Events
// 	case "variant.created",
// 		"variant.archived",
// 		"variant.info.updated",
// 		"variant.in_stock.updated",
// 		"variant.stock_quantity.updated",
// 		"variant.retail_price.updated",
// 		"variant.retail_discount.added",
// 		"variant.retail_discount.updated",
// 		"variant.retail_discount.removed",
// 		"variant.wholesale_mode.enabled",
// 		"variant.wholesale_mode.updated",
// 		"variant.wholesale_mode.disabled",
// 		"variant.wholesale_discount.added",
// 		"variant.wholesale_discount.updated",
// 		"variant.wholesale_discount.removed":
// 		if err := product_variant_snapshot.ProcessSnapshotEvent(ctx, event, q); err != nil {
// 			return err
// 		}
// 		if err := product_variant_index.ProcessIndexEvent(ctx, event, q); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
