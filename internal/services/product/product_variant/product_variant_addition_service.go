// ------------------------------------------------------------
// 📁 File: internal/service/product_variant/add_product_variant_service.go
// 🧠 Handles adding a new variant to an existing product.
//     - Validates seller
//     - Fetches seller profile (for store name)
//     - Confirms product ownership
//     - Validates product category
//     - Fetches product media assets
//     - Inserts new variant
//     - Emits variant.created event
//     - Returns variant_id and product_id

package product_variant

import (
	"context"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type AddProductVariantInput struct {
	UserID                uuid.UUID
	ProductID             uuid.UUID
	Color                 string
	Size                  string
	RetailPrice           int64
	InStock               bool
	StockQuantity         int64
	RetailDiscount        *int64
	RetailDiscountType    *string
	WholesalePrice        *int64
	MinQtyWholesale       *int64
	WholesaleDiscount     *int64
	WholesaleDiscountType *string
	WeightGrams           int64
}

// ------------------------------------------------------------
// 📤 Result to return
type AddProductVariantResult struct {
	ProductID uuid.UUID
	VariantID uuid.UUID
}

// ------------------------------------------------------------
// 🧱 Dependencies
type AddProductVariantServiceDeps struct {
	Repo repo.ProductVariantRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type AddProductVariantService struct {
	Deps AddProductVariantServiceDeps
}

// 🏗️ Constructor
func NewAddProductVariantService(deps AddProductVariantServiceDeps) *AddProductVariantService {
	return &AddProductVariantService{Deps: deps}
}

// 🚀 Main entrypoint
func (s *AddProductVariantService) Start(
	ctx context.Context,
	input AddProductVariantInput,
) (*AddProductVariantResult, error) {

	now := timeutil.NowUTC()
	variantID := uuidutil.New()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {

		// ------------------------------------------------------------
		// Step 1: Validate seller
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsBanned || user.IsArchived {
			return errors.NewAuthError("seller account is not allowed")
		}
		if !user.IsSellerProfileApproved {
			return errors.NewAuthError("seller profile not approved")
		}

		// ------------------------------------------------------------
		// Step 2: Fetch seller profile metadata
		sellerMeta, err := q.GetSellerProfileMetadataBySellerID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("seller_profile_metadata")
		}

		// ------------------------------------------------------------
		// Step 3: Confirm product ownership
		product, err := q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
			ID:       input.ProductID,
			SellerID: input.UserID,
		})
		if err != nil {
			return errors.NewNotFoundError("product")
		}

		// ------------------------------------------------------------
		// Step 4: Fetch & validate category
		category, err := q.GetCategoryByID(ctx, product.CategoryID)
		if err != nil {
			return errors.NewNotFoundError("category")
		}
		if category.IsArchived {
			return errors.NewValidationError("category", "category is archived")
		}
		if !category.IsLeaf {
			return errors.NewValidationError("category", "category must be a leaf node")
		}

		// // ------------------------------------------------------------
		// // Step 5: Fetch promo video URL
		// promoVideos, err := q.GetActiveMediasByProductID(ctx, sqlc.GetActiveMediasByProductIDParams{
		// 	ProductID:  product.ID,
		// 	MediaType:  "video",
		// 	IsArchived: false,
		// })
		// if err != nil {
		// 	return errors.NewTableError("product_medias.video", err.Error())
		// }
		// var promoVideoURL *string
		// if len(promoVideos) > 0 {
		// 	promoVideoURL = &promoVideos[0].MediaUrl
		// }
		// ✅ Safe even if no video exists — optional field

		var promoVideoURL *string
		promoVideos, err := q.GetActiveMediasByProductID(ctx, sqlc.GetActiveMediasByProductIDParams{
			ProductID:  product.ID,
			MediaType:  "video",
			IsArchived: false,
		})
		if err != nil {
			return errors.NewTableError("product_medias.video", err.Error())
		}
		if len(promoVideos) > 0 {
			url := promoVideos[0].MediaUrl
			promoVideoURL = &url
		}

		// ------------------------------------------------------------
		// Step 6: Fetch product image URLs
		imageMedias, err := q.GetActiveMediasByProductID(ctx, sqlc.GetActiveMediasByProductIDParams{
			ProductID:  product.ID,
			MediaType:  "image",
			IsArchived: false,
		})
		if err != nil {
			return errors.NewTableError("product_medias.images", err.Error())
		}
		imageURLs := make([]string, 0)
		for _, media := range imageMedias {
			imageURLs = append(imageURLs, media.MediaUrl)
		}

		// ------------------------------------------------------------
		// Step 7: Fetch primary image
		primaryMedia, err := q.GetPrimaryProductImageByProductID(ctx, sqlc.GetPrimaryProductImageByProductIDParams{
			ProductID:  product.ID,
			MediaType:  "image",
			IsPrimary:  true,
			IsArchived: false,
		})
		var primaryImageURL string
		if err == nil {
			primaryImageURL = primaryMedia.MediaUrl
		}

		// ------------------------------------------------------------
		// Step 8: Insert variant
		hasRetailDiscount := input.RetailDiscount != nil
		hasWholesaleDiscount := input.WholesaleDiscount != nil
		wholesaleEnabled := input.WholesalePrice != nil

		_, err = q.InsertProductVariantReturningID(ctx, sqlc.InsertProductVariantReturningIDParams{
			ID:                    variantID,
			ProductID:             input.ProductID,
			Color:                 input.Color,
			Size:                  input.Size,
			RetailPrice:           input.RetailPrice,
			Retaildiscounttype:    sqlnull.StringPtr(input.RetailDiscountType),
			Retaildiscount:        sqlnull.Int64Ptr(input.RetailDiscount),
			WholesaleEnabled:      wholesaleEnabled,
			WholesalePrice:        sqlnull.Int64Ptr(input.WholesalePrice),
			MinQtyWholesale:       sqlnull.Int32Ptr(input.MinQtyWholesale),
			Wholesalediscounttype: sqlnull.StringPtr(input.WholesaleDiscountType),
			Wholesalediscount:     sqlnull.Int64Ptr(input.WholesaleDiscount),
			StockQuantity:         int32(input.StockQuantity),
			InStock:               input.InStock,
			WeightGrams:           int32(input.WeightGrams),
			IsArchived:            false,
			CreatedAt:             now,
			UpdatedAt:             now,
			Hasretaildiscount:     hasRetailDiscount,
			Haswholesalediscount:  hasWholesaleDiscount,
		})
		if err != nil {
			return errors.NewTableError("product_variants.insert", err.Error())
		}

		// ------------------------------------------------------------
		// Step 9: Emit enriched variant.created event
		payload := map[string]interface{}{
			"product": map[string]interface{}{
				"id":                product.ID,
				"title":             product.Title,
				"description":       product.Description,
				"image_urls":        imageURLs,
				"primary_image_url": primaryImageURL,
				"promo_video_url":   promoVideoURL,
				"is_approved":       product.IsApproved,
				"is_archived":       product.IsArchived,
				"is_banned":         product.IsBanned,
			},
			"category": map[string]interface{}{
				"id":          category.ID,
				"name":        category.Name,
				"is_archived": category.IsArchived,
			},
			"seller": map[string]interface{}{
				"id":                         user.ID,
				"is_archived":                user.IsArchived,
				"is_banned":                  user.IsBanned,
				"is_seller_profile_approved": user.IsSellerProfileApproved,
				"sellerstorename":            sellerMeta.Sellerstorename,
			},
			"variant": map[string]interface{}{
				"id":                      variantID,
				"color":                   input.Color,
				"size":                    input.Size,
				"retail_price":            input.RetailPrice,
				"in_stock":                input.InStock,
				"stock_quantity":          input.StockQuantity,
				"retail_discount":         input.RetailDiscount,
				"retail_discount_type":    input.RetailDiscountType,
				"wholesale_price":         input.WholesalePrice,
				"min_qty_wholesale":       input.MinQtyWholesale,
				"wholesale_discount":      input.WholesaleDiscount,
				"wholesale_discount_type": input.WholesaleDiscountType,
				"weight_grams":            input.WeightGrams,
				"has_retail_discount":     hasRetailDiscount,
				"has_wholesale_discount":  hasWholesaleDiscount,
				"wholesale_enabled":       wholesaleEnabled,
			},
		}

		rawPayload, err := json.Marshal(payload)
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.created",
			EventPayload: rawPayload,
			DispatchedAt: sqlnull.TimePtr(nil),
			CreatedAt:    now,
		})
		if err != nil {
			return errors.NewTableError("events.insert", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AddProductVariantResult{
		ProductID: input.ProductID,
		VariantID: variantID,
	}, nil
}
