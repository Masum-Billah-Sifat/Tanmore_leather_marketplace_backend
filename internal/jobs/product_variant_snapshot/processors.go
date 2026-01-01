package product_variant_snapshot

import (
	"context"
	"encoding/json"
	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"
)

// ───────── Product processors ─────────
func processProductCreated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {

	var event ProductCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	now := timeutil.NowUTC()

	for _, v := range event.Variants {

		hasRetailDiscount := v.RetailDiscount != nil
		hasWholesaleEnabled := v.WholesalePrice != nil
		hasWholesaleDiscount := v.WholesaleDiscount != nil

		err := q.InsertProductVariantSnapshot(ctx, sqlc.InsertProductVariantSnapshotParams{
			ID:                 uuidutil.New(),
			Categoryid:         event.Category.ID,
			Iscategoryarchived: event.Category.IsArchived,
			Categoryname:       event.Category.Name,

			Sellerid:         event.Seller.ID,
			Issellerapproved: event.Seller.IsSellerProfileApproved,
			Issellerarchived: event.Seller.IsArchived,
			Issellerbanned:   event.Seller.IsBanned,
			Sellerstorename:  event.Seller.SellerStoreName,

			Productid:              event.Product.ID,
			Isproductapproved:      event.Product.IsApproved,
			Isproductarchived:      event.Product.IsArchived,
			Isproductbanned:        event.Product.IsBanned,
			Producttitle:           event.Product.Title,
			Productdescription:     event.Product.Description,
			Productprimaryimageurl: event.Product.PrimaryImageURL,

			Variantid:         v.VariantID,
			Isvariantarchived: false,
			Isvariantinstock:  v.InStock,
			Stockamount:       int32(v.StockQuantity),
			Color:             v.Color,
			Size:              v.Size,
			Retailprice:       v.RetailPrice,

			Hasretaildiscount:  hasRetailDiscount,
			Retaildiscounttype: sqlnull.StringPtr(v.RetailDiscountType),
			Retaildiscount:     sqlnull.Int64Ptr(v.RetailDiscount),

			Haswholesaleenabled:   hasWholesaleEnabled,
			Wholesaleprice:        sqlnull.Int64Ptr(v.WholesalePrice),
			Wholesaleminquantity:  sqlnull.Int32Ptr(v.MinQtyWholesale),
			Haswholesalediscount:  hasWholesaleDiscount,
			Wholesalediscounttype: sqlnull.StringPtr(v.WholesaleDiscountType),
			Wholesalediscount:     sqlnull.Int64Ptr(v.WholesaleDiscount),

			WeightGrams: int32(v.WeightGrams),
			Createdat:   now,
			Updatedat:   now,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func processProductInfoUpdated(ctx context.Context, payload []byte, q *sqlc.Queries) error {
	var event ProductInfoUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	now := timeutil.NowUTC()

	title, hasTitle := event.UpdatedFields["title"]
	description, hasDesc := event.UpdatedFields["description"]

	var titlePtr *string
	var descPtr *string
	if hasTitle {
		titlePtr = &title
	}
	if hasDesc {
		descPtr = &description
	}

	return q.UpdateSnapshotsOnProductInfoUpdate(ctx, sqlc.UpdateSnapshotsOnProductInfoUpdateParams{
		Productid:          event.ProductID,
		Producttitle:       sqlnull.StringPtr(titlePtr),
		Productdescription: sqlnull.StringPtr(descPtr),
		Updatedat:          now,
	})
}

func processPrimaryImageSet(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event ProductPrimaryImageSetEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateSnapshotPrimaryImageURL(ctx, sqlc.UpdateSnapshotPrimaryImageURLParams{
		Productid: event.ProductID,
		MediaUrl:  event.MediaURL,
		Updatedat: timeutil.NowUTC(),
	})
}

// func processProductImageRemoved(ctx context.Context, payload []byte, q *sqlc.Queries) error {
// 	return nil
// }
// func processProductPromoVideoRemoved(ctx context.Context, payload []byte, q *sqlc.Queries) error {
// 	return nil
// }

// ───────── Variant processors ─────────
func processVariantCreated(ctx context.Context, payload []byte, q *sqlc.Queries) error {
	var event VariantCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	now := timeutil.NowUTC()

	return q.InsertProductVariantSnapshot(ctx, sqlc.InsertProductVariantSnapshotParams{
		ID:                 uuidutil.New(),
		Categoryid:         event.Category.ID,
		Iscategoryarchived: event.Category.IsArchived,
		Categoryname:       event.Category.Name,

		Sellerid:         event.Seller.ID,
		Issellerapproved: event.Seller.IsSellerProfileApproved,
		Issellerarchived: event.Seller.IsArchived,
		Issellerbanned:   event.Seller.IsBanned,
		Sellerstorename:  event.Seller.SellerStoreName,

		Productid:              event.Product.ID,
		Isproductapproved:      event.Product.IsApproved,
		Isproductarchived:      event.Product.IsArchived,
		Isproductbanned:        event.Product.IsBanned,
		Producttitle:           event.Product.Title,
		Productdescription:     event.Product.Description,
		Productprimaryimageurl: event.Product.PrimaryImageURL,

		Variantid:         event.Variant.ID,
		Isvariantarchived: false,
		Isvariantinstock:  event.Variant.InStock,
		Stockamount:       int32(event.Variant.StockQuantity),
		Color:             event.Variant.Color,
		Size:              event.Variant.Size,
		Retailprice:       event.Variant.RetailPrice,

		Hasretaildiscount:  event.Variant.HasRetailDiscount,
		Retaildiscounttype: sqlnull.StringPtr(event.Variant.RetailDiscountType),
		Retaildiscount:     sqlnull.Int64Ptr(event.Variant.RetailDiscount),

		Haswholesaleenabled:   event.Variant.WholesaleEnabled,
		Wholesaleprice:        sqlnull.Int64Ptr(event.Variant.WholesalePrice),
		Wholesaleminquantity:  sqlnull.Int32Ptr(event.Variant.MinQtyWholesale),
		Haswholesalediscount:  event.Variant.HasWholesaleDiscount,
		Wholesalediscounttype: sqlnull.StringPtr(event.Variant.WholesaleDiscountType),
		Wholesalediscount:     sqlnull.Int64Ptr(event.Variant.WholesaleDiscount),

		WeightGrams: int32(event.Variant.WeightGrams),
		Createdat:   now,
		Updatedat:   now,
	})
}

func processVariantArchived(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantArchivedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.ArchiveVariantInIndexes(ctx, sqlc.ArchiveVariantInIndexesParams{
		Productid:         event.ProductID,
		Variantid:         event.VariantID,
		Isvariantarchived: true,
		Updatedat:         timeutil.NowUTC(),
	})
}

func processVariantInfoUpdated(ctx context.Context, payload []byte, q *sqlc.Queries) error {
	var event VariantInfoUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	now := timeutil.NowUTC()

	color, hasColor := event.UpdatedFields["color"]
	size, hasSize := event.UpdatedFields["size"]

	var colorPtr *string
	var sizePtr *string

	if hasColor {
		colorPtr = &color
	}
	if hasSize {
		sizePtr = &size
	}

	return q.UpdateIndexOnVariantInfoUpdate(ctx, sqlc.UpdateIndexOnVariantInfoUpdateParams{
		Productid: event.ProductID,
		Variantid: event.VariantID,
		Color:     sqlnull.StringPtr(colorPtr),
		Size:      sqlnull.StringPtr(sizePtr),
		Updatedat: now,
	})
}

func processVariantInStockUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantInStockUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateSnapshotOnInStockUpdate(ctx, sqlc.UpdateSnapshotOnInStockUpdateParams{
		Productid: event.ProductID,
		Variantid: event.VariantID,
		InStock:   event.InStock,
		Updatedat: timeutil.NowUTC(),
	})
}

func processVariantRetailPriceUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantRetailPriceUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateSnapshotOnRetailPriceUpdate(ctx, sqlc.UpdateSnapshotOnRetailPriceUpdateParams{
		Productid:   event.ProductID,
		Variantid:   event.VariantID,
		RetailPrice: event.RetailPrice,
		Updatedat:   timeutil.NowUTC(),
	})
}

func processVariantStockQuantityUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantStockQuantityUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateSnapshotOnStockQuantityUpdate(ctx, sqlc.UpdateSnapshotOnStockQuantityUpdateParams{
		Productid:     event.ProductID,
		Variantid:     event.VariantID,
		StockQuantity: int32(event.StockQuantity),
		Updatedat:     timeutil.NowUTC(),
	})
}

func processVariantWeightUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantWeightUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateSnapshotOnWeightUpdate(ctx, sqlc.UpdateSnapshotOnWeightUpdateParams{
		Productid:   event.ProductID,
		Variantid:   event.VariantID,
		WeightGrams: int32(event.WeightGrams),
		Updatedat:   timeutil.NowUTC(),
	})
}

func processVariantRetailDiscountAdded(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantRetailDiscountAddedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.ApplyRetailDiscountToSnapshot(ctx, sqlc.ApplyRetailDiscountToSnapshotParams{
		Productid:          event.ProductID,
		Variantid:          event.VariantID,
		HasRetailDiscount:  true,
		RetailDiscount:     sqlnull.Int64(event.RetailDiscount),
		RetailDiscountType: sqlnull.String(event.RetailDiscountType),
		Updatedat:          timeutil.NowUTC(),
	})
}

func processVariantRetailDiscountUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantRetailDiscountUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	var discountPtr *int64
	var discountTypePtr *string

	if raw, ok := event.UpdatedFields["retail_discount"]; ok {
		if val, ok := raw.(float64); ok { // JSON numbers are float64 by default
			i64 := int64(val)
			discountPtr = &i64
		}
	}

	if raw, ok := event.UpdatedFields["retail_discount_type"]; ok {
		if val, ok := raw.(string); ok {
			discountTypePtr = &val
		}
	}

	return q.UpdateRetailDiscountInSnapshot(ctx, sqlc.UpdateRetailDiscountInSnapshotParams{
		Productid:          event.ProductID,
		Variantid:          event.VariantID,
		RetailDiscount:     sqlnull.Int64Ptr(discountPtr),
		RetailDiscountType: sqlnull.StringPtr(discountTypePtr),
		Updatedat:          timeutil.NowUTC(),
	})
}

func processVariantRetailDiscountRemoved(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantRetailDiscountRemovedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.RemoveRetailDiscountFromSnapshot(ctx, sqlc.RemoveRetailDiscountFromSnapshotParams{
		Productid:          event.ProductID,
		Variantid:          event.VariantID,
		HasRetailDiscount:  false,
		RetailDiscount:     sqlnull.Int64Ptr(nil),
		RetailDiscountType: sqlnull.StringPtr(nil),
		Updatedat:          timeutil.NowUTC(),
	})
}

func processVariantWholesaleModeEnabled(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantWholesaleModeEnabledEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.EnableWholesaleModeInSnapshots(ctx, sqlc.EnableWholesaleModeInSnapshotsParams{
		Productid: event.ProductID,
		Variantid: event.VariantID,
		// ✅ FIXED
		WholesalePrice:        sqlnull.Int64(event.WholesalePrice),
		MinQtyWholesale:       sqlnull.Int32(event.MinQtyWholesale),
		WholesaleDiscount:     sqlnull.Int64Ptr(event.WholesaleDiscount),
		WholesaleDiscountType: sqlnull.StringPtr(event.WholesaleDiscountType),
		HasWholesaleDiscount:  event.HasWholesaleDiscount,
		HasWholesaleEnabled:   true,
		Updatedat:             timeutil.NowUTC(),
	})
}

func processVariantWholesaleModeUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantWholesaleModeUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	var wholesalePricePtr *int64
	var minQtyPtr *int64

	if raw, ok := event.UpdatedFields["wholesale_price"]; ok {
		if val, ok := raw.(float64); ok {
			i64 := int64(val)
			wholesalePricePtr = &i64
		}
	}
	if raw, ok := event.UpdatedFields["min_qty_wholesale"]; ok {
		if val, ok := raw.(float64); ok {
			i64 := int64(val)
			minQtyPtr = &i64
		}
	}

	return q.UpdateWholesaleModeInSnapshots(ctx, sqlc.UpdateWholesaleModeInSnapshotsParams{
		Productid:       event.ProductID,
		Variantid:       event.VariantID,
		WholesalePrice:  sqlnull.Int64Ptr(wholesalePricePtr),
		MinQtyWholesale: sqlnull.Int32Ptr(minQtyPtr),
		Updatedat:       timeutil.NowUTC(),
	})
}

func processVariantWholesaleModeDisabled(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantWholesaleModeDisabledEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.DisableWholesaleModeInSnapshots(ctx, sqlc.DisableWholesaleModeInSnapshotsParams{
		Productid:             event.ProductID,
		Variantid:             event.VariantID,
		HasWholesaleEnabled:   false,
		HasWholesaleDiscount:  false,
		WholesalePrice:        sqlnull.Int64Ptr(nil),
		WholesaleMinQuantity:  sqlnull.Int32Ptr(nil),
		WholesaleDiscount:     sqlnull.Int64Ptr(nil),
		WholesaleDiscountType: sqlnull.StringPtr(nil),
		Updatedat:             timeutil.NowUTC(),
	})
}

func processVariantWholesaleDiscountAdded(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantWholesaleDiscountAddedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.ApplyWholesaleDiscountToSnapshot(ctx, sqlc.ApplyWholesaleDiscountToSnapshotParams{
		Productid:             event.ProductID,
		Variantid:             event.VariantID,
		HasWholesaleDiscount:  true,
		WholesaleDiscount:     sqlnull.Int64(event.WholesaleDiscount),
		WholesaleDiscountType: sqlnull.String(event.WholesaleDiscountType),
		Updatedat:             timeutil.NowUTC(),
	})
}

func processVariantWholesaleDiscountUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantWholesaleDiscountUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	var discountPtr *int64
	var discountTypePtr *string

	if raw, ok := event.UpdatedFields["wholesale_discount"]; ok {
		if v, ok := raw.(float64); ok {
			i := int64(v)
			discountPtr = &i
		}
	}

	if raw, ok := event.UpdatedFields["wholesale_discount_type"]; ok {
		if v, ok := raw.(string); ok {
			discountTypePtr = &v
		}
	}

	return q.UpdateWholesaleDiscountInSnapshot(ctx, sqlc.UpdateWholesaleDiscountInSnapshotParams{
		Productid:             event.ProductID,
		Variantid:             event.VariantID,
		WholesaleDiscount:     sqlnull.Int64Ptr(discountPtr),
		WholesaleDiscountType: sqlnull.StringPtr(discountTypePtr),
		Updatedat:             timeutil.NowUTC(),
	})
}

func processVariantWholesaleDiscountRemoved(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantWholesaleDiscountRemovedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.RemoveWholesaleDiscountFromSnapshot(ctx, sqlc.RemoveWholesaleDiscountFromSnapshotParams{
		Productid:             event.ProductID,
		Variantid:             event.VariantID,
		HasWholesaleDiscount:  false,
		WholesaleDiscount:     sqlnull.Int64Ptr(nil),
		WholesaleDiscountType: sqlnull.StringPtr(nil),
		Updatedat:             timeutil.NowUTC(),
	})
}
