package product_variant_index

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

	searchText := event.Product.Title + " " + event.Product.Description

	for _, v := range event.Variants {

		hasRetailDiscount := v.RetailDiscount != nil
		hasWholesaleEnabled := v.WholesalePrice != nil

		err := q.InsertProductVariantIndex(ctx, sqlc.InsertProductVariantIndexParams{
			ID:                 uuidutil.New(),
			Categoryid:         event.Category.ID,
			Iscategoryarchived: event.Category.IsArchived,
			Categoryname:       event.Category.Name,

			Sellerid:         event.Seller.ID,
			Issellerapproved: event.Seller.IsSellerProfileApproved,
			Issellerarchived: event.Seller.IsArchived,
			Issellerbanned:   event.Seller.IsBanned,
			Sellerstorename:  event.Seller.SellerStoreName,

			Productid:            event.Product.ID,
			Isproductapproved:    event.Product.IsApproved,
			Isproductarchived:    event.Product.IsArchived,
			Isproductbanned:      event.Product.IsBanned,
			Producttitle:         event.Product.Title,
			Productdescription:   event.Product.Description,
			Productimages:        event.Product.ImageURLs,
			Productpromovideourl: sqlnull.StringPtr(event.Product.PromoVideoURL),

			Variantid:         v.VariantID,
			Isvariantarchived: false,
			Isvariantinstock:  v.InStock,
			Stockamount:       int32(v.StockQuantity),
			Color:             v.Color,
			Size:              v.Size,
			Retailprice:       v.RetailPrice,

			Retaildiscounttype: sqlnull.StringPtr(v.RetailDiscountType),
			Retaildiscount:     sqlnull.Int64Ptr(v.RetailDiscount),
			HasRetailDiscount:  hasRetailDiscount,

			Haswholesaleenabled:   hasWholesaleEnabled,
			Wholesaleprice:        sqlnull.Int64Ptr(v.WholesalePrice),
			Wholesaleminquantity:  sqlnull.Int32Ptr(v.MinQtyWholesale),
			Wholesalediscounttype: sqlnull.StringPtr(v.WholesaleDiscountType),
			Wholesalediscount:     sqlnull.Int64Ptr(v.WholesaleDiscount),

			WeightGrams: int32(v.WeightGrams),
			// SearchVector: searchText,
			ToTsvector: json.RawMessage(`"` + searchText + `"`),
			Views:      0,
			Createdat:  now,
			Updatedat:  now,
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

	return q.UpdateIndexesOnProductInfoUpdate(ctx, sqlc.UpdateIndexesOnProductInfoUpdateParams{
		Productid:          event.ProductID,
		Producttitle:       sqlnull.StringPtr(titlePtr),
		Productdescription: sqlnull.StringPtr(descPtr),
		Updatedat:          now,
	})
}

func processProductImageAdded(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event ProductMediaAddedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateIndexOnProductImageAdded(ctx, sqlc.UpdateIndexOnProductImageAddedParams{
		Productid: event.ProductID,
		MediaUrl:  event.MediaURL,
		Updatedat: timeutil.NowUTC(),
	})
}

func processProductImageRemoved(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event ProductMediaRemovedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateIndexOnImageRemoved(ctx, sqlc.UpdateIndexOnImageRemovedParams{
		Productid: event.ProductID,
		MediaUrl:  sqlnull.String(event.MediaURL),
		Updatedat: timeutil.NowUTC(),
	})
}

// func processProductPrimaryImageSet(ctx context.Context, payload []byte, q *sqlc.Queries) error {
// 	return nil
// }

func processProductPromoVideoAdded(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event ProductMediaAddedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateIndexOnPromoVideoAdded(ctx, sqlc.UpdateIndexOnPromoVideoAddedParams{
		Productid: event.ProductID,
		MediaUrl:  sqlnull.String(event.MediaURL),
		Updatedat: timeutil.NowUTC(),
	})
}

func processProductPromoVideoRemoved(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event ProductMediaRemovedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateIndexOnPromoVideoRemoved(ctx, sqlc.UpdateIndexOnPromoVideoRemovedParams{
		Productid:            event.ProductID,
		Productpromovideourl: sqlnull.String(""), // This will result in `NULL`
		Updatedat:            timeutil.NowUTC(),
	})
}

// ───────── Variant processors ─────────
func processVariantCreated(ctx context.Context, payload []byte, q *sqlc.Queries) error {
	var event VariantCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	now := timeutil.NowUTC()

	return q.InsertProductVariantIndex(ctx, sqlc.InsertProductVariantIndexParams{
		ID:                 uuidutil.New(),
		Categoryid:         event.Category.ID,
		Iscategoryarchived: event.Category.IsArchived,
		Categoryname:       event.Category.Name,

		Sellerid:         event.Seller.ID,
		Issellerapproved: event.Seller.IsSellerProfileApproved,
		Issellerarchived: event.Seller.IsArchived,
		Issellerbanned:   event.Seller.IsBanned,
		Sellerstorename:  event.Seller.SellerStoreName,

		Productid:            event.Product.ID,
		Isproductapproved:    event.Product.IsApproved,
		Isproductarchived:    event.Product.IsArchived,
		Isproductbanned:      event.Product.IsBanned,
		Producttitle:         event.Product.Title,
		Productdescription:   event.Product.Description,
		Productimages:        event.Product.ImageURLs,
		Productpromovideourl: sqlnull.StringPtr(event.Product.PromoVideoURL),

		Variantid:         event.Variant.ID,
		Isvariantarchived: false,
		Isvariantinstock:  event.Variant.InStock,
		Stockamount:       int32(event.Variant.StockQuantity),
		Color:             event.Variant.Color,
		Size:              event.Variant.Size,
		Retailprice:       event.Variant.RetailPrice,

		Retaildiscounttype: sqlnull.StringPtr(event.Variant.RetailDiscountType),
		Retaildiscount:     sqlnull.Int64Ptr(event.Variant.RetailDiscount),
		HasRetailDiscount:  event.Variant.HasRetailDiscount,

		Haswholesaleenabled:   event.Variant.WholesaleEnabled,
		Wholesaleprice:        sqlnull.Int64Ptr(event.Variant.WholesalePrice),
		Wholesaleminquantity:  sqlnull.Int32Ptr(event.Variant.MinQtyWholesale),
		Wholesalediscounttype: sqlnull.StringPtr(event.Variant.WholesaleDiscountType),
		Wholesalediscount:     sqlnull.Int64Ptr(event.Variant.WholesaleDiscount),

		WeightGrams: int32(event.Variant.WeightGrams),
		// SearchVector: event.Product.Title + " " + event.Product.Description,
		ToTsvector: json.RawMessage(`"` + event.Product.Title + " " + event.Product.Description + `"`), // ✅ fixed

		Views:     0,
		Createdat: now,
		Updatedat: now,
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

	return q.ArchiveVariantInSnapshots(ctx, sqlc.ArchiveVariantInSnapshotsParams{
		Productid:         event.ProductID,
		Variantid:         event.VariantID,
		Isvariantarchived: true,
		Updatedat:         timeutil.NowUTC(),
	})
}

func processVariantInfoUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantInfoUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

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

	return q.UpdateSnapshotOnVariantInfoUpdate(ctx, sqlc.UpdateSnapshotOnVariantInfoUpdateParams{
		Productid: event.ProductID,
		Variantid: event.VariantID,
		Color:     sqlnull.StringPtr(colorPtr),
		Size:      sqlnull.StringPtr(sizePtr),
		Updatedat: timeutil.NowUTC(),
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

	return q.UpdateIndexOnInStockUpdate(ctx, sqlc.UpdateIndexOnInStockUpdateParams{
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

	return q.UpdateIndexOnRetailPriceUpdate(ctx, sqlc.UpdateIndexOnRetailPriceUpdateParams{
		Productid:   event.ProductID,
		Variantid:   event.VariantID,
		RetailPrice: event.RetailPrice,
		Updatedat:   timeutil.NowUTC(),
	})
}

// more on top like this
func processVariantStockQuantityUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event VariantStockQuantityUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateIndexOnStockQuantityUpdate(ctx, sqlc.UpdateIndexOnStockQuantityUpdateParams{
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

	return q.UpdateIndexOnWeightUpdate(ctx, sqlc.UpdateIndexOnWeightUpdateParams{
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

	return q.ApplyRetailDiscountToIndex(ctx, sqlc.ApplyRetailDiscountToIndexParams{
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
		if val, ok := raw.(float64); ok {
			i64 := int64(val)
			discountPtr = &i64
		}
	}

	if raw, ok := event.UpdatedFields["retail_discount_type"]; ok {
		if val, ok := raw.(string); ok {
			discountTypePtr = &val
		}
	}

	return q.UpdateRetailDiscountInIndex(ctx, sqlc.UpdateRetailDiscountInIndexParams{
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

	return q.RemoveRetailDiscountFromIndex(ctx, sqlc.RemoveRetailDiscountFromIndexParams{
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

	return q.EnableWholesaleModeInIndexes(ctx, sqlc.EnableWholesaleModeInIndexesParams{
		Productid: event.ProductID,
		Variantid: event.VariantID,
		// ✅ FIXED
		WholesalePrice:  sqlnull.Int64(event.WholesalePrice),
		MinQtyWholesale: sqlnull.Int32(event.MinQtyWholesale),

		WholesaleDiscount:     sqlnull.Int64Ptr(event.WholesaleDiscount),
		WholesaleDiscountType: sqlnull.StringPtr(event.WholesaleDiscountType),
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

	return q.UpdateWholesaleModeInIndexes(ctx, sqlc.UpdateWholesaleModeInIndexesParams{
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

	return q.DisableWholesaleModeInIndexes(ctx, sqlc.DisableWholesaleModeInIndexesParams{
		Productid:             event.ProductID,
		Variantid:             event.VariantID,
		HasWholesaleEnabled:   false,
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

	return q.ApplyWholesaleDiscountToIndex(ctx, sqlc.ApplyWholesaleDiscountToIndexParams{
		Productid:             event.ProductID,
		Variantid:             event.VariantID,
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

	return q.UpdateWholesaleDiscountInIndex(ctx, sqlc.UpdateWholesaleDiscountInIndexParams{
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

	return q.RemoveWholesaleDiscountFromIndex(ctx, sqlc.RemoveWholesaleDiscountFromIndexParams{
		Productid:             event.ProductID,
		Variantid:             event.VariantID,
		WholesaleDiscount:     sqlnull.Int64Ptr(nil),
		WholesaleDiscountType: sqlnull.StringPtr(nil),
		Updatedat:             timeutil.NowUTC(),
	})
}

func processProductCategoryUpdated(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event ProductCategoryUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.UpdateCategoryInProductVariantIndexes(ctx, sqlc.UpdateCategoryInProductVariantIndexesParams{
		Categoryid:   event.NewCategoryID,
		Categoryname: event.NewCategoryName,
		Updatedat:    timeutil.NowUTC(),
		Productid:    event.ProductID,
	})
}

func processProductArchived(
	ctx context.Context,
	payload []byte,
	q *sqlc.Queries,
) error {
	var event ProductArchivedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return q.MarkProductArchivedInProductVariantIndexes(ctx, sqlc.MarkProductArchivedInProductVariantIndexesParams{
		Updatedat: timeutil.NowUTC(),
		Productid: event.ProductID,
	})
}
