package product_variant_snapshot

import (
	"context"
	"tanmore_backend/internal/db/sqlc"
)

func ProcessSnapshotEvent(ctx context.Context, event sqlc.Event, q *sqlc.Queries) error {
	switch event.EventType {

	// ───────── Product events ─────────
	case "product.created":
		return processProductCreated(ctx, event.EventPayload, q)

	case "product.info.updated":
		return processProductInfoUpdated(ctx, event.EventPayload, q)

	case "product.image.set_primary":
		return processPrimaryImageSet(ctx, event.EventPayload, q)

	// case "product.image.removed":
	// 	return processProductImageRemoved(ctx, event.EventPayload, q)

	// case "product.promo_video.removed":
	// 	return processProductPromoVideoRemoved(ctx, event.EventPayload, q)

	// ───────── Variant events ─────────
	case "variant.created":
		return processVariantCreated(ctx, event.EventPayload, q)

	case "variant.archived":
		return processVariantArchived(ctx, event.EventPayload, q)

	case "variant.info.updated":
		return processVariantInfoUpdated(ctx, event.EventPayload, q)

	case "variant.in_stock.updated":
		return processVariantInStockUpdated(ctx, event.EventPayload, q)

	case "variant.stock_quantity.updated":
		return processVariantStockQuantityUpdated(ctx, event.EventPayload, q)

	case "variant.weight.updated":
		return processVariantWeightUpdated(ctx, event.EventPayload, q)

	case "variant.retail_price.updated":
		return processVariantRetailPriceUpdated(ctx, event.EventPayload, q)

	case "variant.retail_discount.added":
		return processVariantRetailDiscountAdded(ctx, event.EventPayload, q)

	case "variant.retail_discount.updated":
		return processVariantRetailDiscountUpdated(ctx, event.EventPayload, q)

	case "variant.retail_discount.removed":
		return processVariantRetailDiscountRemoved(ctx, event.EventPayload, q)

	case "variant.wholesale_mode.enabled":
		return processVariantWholesaleModeEnabled(ctx, event.EventPayload, q)

	case "variant.wholesale_mode.updated":
		return processVariantWholesaleModeUpdated(ctx, event.EventPayload, q)

	case "variant.wholesale_mode.disabled":
		return processVariantWholesaleModeDisabled(ctx, event.EventPayload, q)

	case "variant.wholesale_discount.added":
		return processVariantWholesaleDiscountAdded(ctx, event.EventPayload, q)

	case "variant.wholesale_discount.updated":
		return processVariantWholesaleDiscountUpdated(ctx, event.EventPayload, q)

	case "variant.wholesale_discount.removed":
		return processVariantWholesaleDiscountRemoved(ctx, event.EventPayload, q)

	default:
		return nil
	}
}
