// ------------------------------------------------------------
// 📁 File: internal/service/product/create_product_service.go
// 🧠 Handles seller product creation workflow.
//     - Validates seller moderation & approval
//     - Validates seller profile metadata
//     - Validates category (leaf + not archived)
//     - Creates product
//     - Creates product medias (images + promo video)
//     - Creates variants
//     - Emits product.created event
//     - Returns product_id and variant_ids

package product

import (
	"context"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_creation"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📦 Variant input from handler
type CreateProductVariantInput struct {
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
// 📦 Input from handler
type CreateProductInput struct {
	UserID        uuid.UUID
	CategoryID    uuid.UUID
	Title         string
	Description   string
	ImageURLs     []string
	PromoVideoURL *string
	Variants      []CreateProductVariantInput
}

// ------------------------------------------------------------
// 📦 Result returned to handler
type CreateProductResult struct {
	ProductID  uuid.UUID
	VariantIDs []uuid.UUID
}

// ------------------------------------------------------------
// 📦 Dependencies
type CreateProductServiceDeps struct {
	Repo repo.ProductRepoInterface
}

// ------------------------------------------------------------
// 📦 Service definition
type CreateProductService struct {
	Deps CreateProductServiceDeps
}

// ------------------------------------------------------------
// 🛠️ Constructor
func NewCreateProductService(deps CreateProductServiceDeps) *CreateProductService {
	return &CreateProductService{Deps: deps}
}

// ------------------------------------------------------------
// 🚀 Main Logic
func (s *CreateProductService) Start(
	ctx context.Context,
	input CreateProductInput,
) (*CreateProductResult, error) {

	now := timeutil.NowUTC()
	productID := uuidutil.New()
	variantIDs := make([]uuid.UUID, 0)

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {

		// ------------------------------------------------------------
		// Step 1: Validate seller
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}

		if user.IsArchived || user.IsBanned {
			return errors.NewAuthError("seller account is not allowed")
		}

		if !user.IsSellerProfileApproved {
			return errors.NewAuthError("seller profile not approved")
		}

		// ------------------------------------------------------------
		// Step 2: Fetch seller profile metadata (for event payload)
		sellerMeta, err := q.GetSellerProfileMetadataBySellerID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("seller_profile_metadata")
		}

		// ------------------------------------------------------------
		// Step 3: Validate category
		category, err := q.GetCategoryByID(ctx, input.CategoryID)
		if err != nil {
			return errors.NewNotFoundError("category")
		}

		if category.IsArchived {
			return errors.NewValidationError("category", "category is archived")
		}

		if !category.IsLeaf {
			return errors.NewValidationError("category", "category must be a leaf node")
		}

		// ------------------------------------------------------------
		// Step 4: Insert product
		_, err = q.InsertProduct(ctx, sqlc.InsertProductParams{
			ID:          productID,
			SellerID:    input.UserID,
			CategoryID:  input.CategoryID,
			Title:       input.Title,
			Description: input.Description,
			IsApproved:  false,
			IsBanned:    false,
			IsArchived:  false,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		if err != nil {
			return errors.NewTableError("products.insert", err.Error())
		}

		// ------------------------------------------------------------
		// Step 5: Insert product images (Meta-grade explicit loop)
		var primaryImageURL string

		for i, url := range input.ImageURLs {
			isPrimary := i == 0
			if isPrimary {
				primaryImageURL = url
			}

			_, err := q.InsertProductMedia(ctx, sqlc.InsertProductMediaParams{
				ID:         uuidutil.New(),
				ProductID:  productID,
				MediaType:  "image",
				MediaUrl:   url,
				IsPrimary:  isPrimary,
				IsArchived: false,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
			if err != nil {
				return errors.NewTableError("product_medias.insert.image", err.Error())
			}
		}

		// ------------------------------------------------------------
		// Step 6: Insert promo video (if provided)
		if input.PromoVideoURL != nil {
			_, err := q.InsertProductMedia(ctx, sqlc.InsertProductMediaParams{
				ID:         uuidutil.New(),
				ProductID:  productID,
				MediaType:  "video",
				MediaUrl:   *input.PromoVideoURL,
				IsPrimary:  false,
				IsArchived: false,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
			if err != nil {
				return errors.NewTableError("product_medias.insert.video", err.Error())
			}
		}

		// ------------------------------------------------------------
		// Step 7: Insert variants (UNCHANGED – already correct)
		variantPayloads := make([]map[string]interface{}, 0, len(input.Variants))

		for _, v := range input.Variants {

			variantID := uuidutil.New()

			hasRetailDiscount := v.RetailDiscount != nil
			hasWholesaleDiscount := v.WholesaleDiscount != nil
			wholesaleEnabled := v.WholesalePrice != nil

			_, err := q.InsertProductVariantReturningID(
				ctx,
				sqlc.InsertProductVariantReturningIDParams{
					ID:                    variantID,
					ProductID:             productID,
					Color:                 v.Color,
					Size:                  v.Size,
					RetailPrice:           v.RetailPrice,
					Retaildiscounttype:    sqlnull.StringPtr(v.RetailDiscountType),
					Retaildiscount:        sqlnull.Int64Ptr(v.RetailDiscount),
					WholesaleEnabled:      wholesaleEnabled,
					WholesalePrice:        sqlnull.Int64Ptr(v.WholesalePrice),
					MinQtyWholesale:       sqlnull.Int32Ptr(v.MinQtyWholesale),
					Wholesalediscounttype: sqlnull.StringPtr(v.WholesaleDiscountType),
					Wholesalediscount:     sqlnull.Int64Ptr(v.WholesaleDiscount),
					StockQuantity:         int32(v.StockQuantity),
					InStock:               v.InStock,
					WeightGrams:           int32(v.WeightGrams),
					IsArchived:            false,
					CreatedAt:             now,
					UpdatedAt:             now,
					Hasretaildiscount:     hasRetailDiscount,
					Haswholesalediscount:  hasWholesaleDiscount,
				},
			)
			if err != nil {
				return errors.NewTableError("product_variants.insert", err.Error())
			}

			variantIDs = append(variantIDs, variantID)

			variantPayloads = append(variantPayloads, map[string]interface{}{
				"variant_id":              variantID,
				"color":                   v.Color,
				"size":                    v.Size,
				"retail_price":            v.RetailPrice,
				"in_stock":                v.InStock,
				"stock_quantity":          v.StockQuantity,
				"retail_discount":         v.RetailDiscount,
				"retail_discount_type":    v.RetailDiscountType,
				"wholesale_price":         v.WholesalePrice,
				"min_qty_wholesale":       v.MinQtyWholesale,
				"wholesale_discount":      v.WholesaleDiscount,
				"wholesale_discount_type": v.WholesaleDiscountType,
				"weight_grams":            v.WeightGrams,
			})
		}

		// ------------------------------------------------------------
		// Step 8: Emit product.created event (FULL payload)
		eventPayload := map[string]interface{}{
			"product": map[string]interface{}{
				"id":                productID,
				"title":             input.Title,
				"description":       input.Description,
				"image_urls":        input.ImageURLs,
				"primary_image_url": primaryImageURL,
				"promo_video_url":   input.PromoVideoURL,
				"is_approved":       false,
				"is_archived":       false,
				"is_banned":         false,
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
			"variants": variantPayloads,
		}

		rawPayload, err := json.Marshal(eventPayload)
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "product.created",
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

	return &CreateProductResult{
		ProductID:  productID,
		VariantIDs: variantIDs,
	}, nil
}
