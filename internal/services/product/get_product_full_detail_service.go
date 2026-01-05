// ------------------------------------------------------------
// 📁 File: internal/services/product/get_product_full_detail_service.go
// 🧠 Handles fetching full product detail for seller
//     - Validates seller moderation
//     - Validates product ownership & moderation
//     - Fetches all variant indexes
//     - Fetches primary product image
//     - Groups variants into valid / archived
//     - Formats final response

package product

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_get_full_detail"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler

type GetProductFullDetailInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Variant Response

type ProductVariantResponse struct {
	VariantID             uuid.UUID `json:"variant_id"`
	Color                 string    `json:"color"`
	Size                  string    `json:"size"`
	RetailPrice           int64     `json:"retail_price"`
	HasRetailDiscount     bool      `json:"has_retail_discount"`
	RetailDiscount        *int64    `json:"retail_discount,omitempty"`
	RetailDiscountType    *string   `json:"retail_discount_type,omitempty"`
	IsInStock             bool      `json:"is_in_stock"`
	StockQuantity         int32     `json:"stock_quantity"`
	HasWholesaleEnabled   bool      `json:"has_wholesale_enabled"`
	WholesalePrice        *int64    `json:"wholesale_price,omitempty"`
	WholesaleMinQuantity  *int32    `json:"wholesale_min_quantity,omitempty"`
	WholesaleDiscount     *int64    `json:"wholesale_discount,omitempty"`
	WholesaleDiscountType *string   `json:"wholesale_discount_type,omitempty"`
	WeightGrams           int32     `json:"weight_grams"`
	IsVariantArchived     bool      `json:"is_variant_archived"`
}

// ------------------------------------------------------------
// 📤 Final Response

type GetProductFullDetailResult struct {
	ProductID         uuid.UUID `json:"product_id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	CategoryID        uuid.UUID `json:"category_id"`
	CategoryName      string    `json:"category_name"`
	SellerID          uuid.UUID `json:"seller_id"`
	SellerStoreName   string    `json:"seller_store_name"`
	ImageURLs         []string  `json:"image_urls"`
	PrimaryImageURL   *string   `json:"primary_image_url"`
	PromoVideoURL     *string   `json:"promo_video_url"`
	IsProductApproved bool      `json:"is_product_approved"`

	ValidVariants    []ProductVariantResponse `json:"valid_variants"`
	ArchivedVariants []ProductVariantResponse `json:"archived_variants"`
}

// ------------------------------------------------------------
// 🧱 Dependencies

type GetProductFullDetailServiceDeps struct {
	Repo repo.ProductGetFullDetailRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service

type GetProductFullDetailService struct {
	Deps GetProductFullDetailServiceDeps
}

// 🚀 Constructor
func NewGetProductFullDetailService(
	deps GetProductFullDetailServiceDeps,
) *GetProductFullDetailService {
	return &GetProductFullDetailService{Deps: deps}
}

// ------------------------------------------------------------
// 🚀 Entrypoint

func (s *GetProductFullDetailService) Start(
	ctx context.Context,
	input GetProductFullDetailInput,
) (*GetProductFullDetailResult, error) {

	// ------------------------------------------------------------
	// Step 1: Validate seller identity

	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, errors.NewNotFoundError("seller")
	}

	if user.IsArchived || user.IsBanned || !user.IsSellerProfileApproved || !user.IsSellerProfileCreated {
		return nil, errors.NewValidationError("seller", "not allowed")
	}

	// ------------------------------------------------------------
	// Step 2: Validate product ownership & moderation

	product, err := s.Deps.Repo.GetProductByIDAndSellerID(ctx, input.ProductID, input.UserID)
	if err != nil {
		return nil, errors.NewNotFoundError("product")
	}

	if product.IsBanned || product.IsArchived {
		return nil, errors.NewValidationError("product", "banned or archived")
	}

	// ------------------------------------------------------------
	// Step 3: Fetch variant indexes

	variants, err := s.Deps.Repo.GetVariantIndexesByProductAndSeller(
		ctx,
		sqlc.GetVariantIndexesByProductAndSellerParams{
			Productid: input.ProductID,
			Sellerid:  input.UserID,
		},
	)
	if err != nil {
		return nil, errors.NewTableError("product_variant_indexes.select", err.Error())
	}

	if len(variants) == 0 {
		return nil, errors.NewValidationError("variants", "no variants found")
	}

	// ------------------------------------------------------------
	// Step 4: Fetch primary image (optional)

	var primaryImageURL *string

	media, err := s.Deps.Repo.GetPrimaryImageForProduct(
		ctx,
		sqlc.GetPrimaryProductImageByProductIDParams{
			ProductID:  input.ProductID,
			MediaType:  "image",
			IsPrimary:  true,
			IsArchived: false,
		},
	)
	if err == nil {
		primaryImageURL = &media.MediaUrl
	}

	// ------------------------------------------------------------
	// Step 5: Group variants

	var validVariants []ProductVariantResponse
	var archivedVariants []ProductVariantResponse

	for _, v := range variants {

		// safety checks (as per your spec)
		if v.Issellerbanned || v.Issellerarchived || !v.Issellerapproved {
			return nil, errors.NewValidationError("seller", "invalid seller state")
		}
		if v.Isproductbanned || v.Isproductarchived {
			return nil, errors.NewValidationError("product", "invalid product state")
		}

		variant := ProductVariantResponse{
			VariantID:             v.Variantid,
			Color:                 v.Color,
			Size:                  v.Size,
			RetailPrice:           v.Retailprice,
			HasRetailDiscount:     v.HasRetailDiscount,
			RetailDiscount:        sqlnull.ToInt64Ptr(v.Retaildiscount),
			RetailDiscountType:    sqlnull.ToStringPtr(v.Retaildiscounttype),
			IsInStock:             v.Isvariantinstock,
			StockQuantity:         v.Stockamount,
			HasWholesaleEnabled:   v.Haswholesaleenabled,
			WholesalePrice:        sqlnull.ToInt64Ptr(v.Wholesaleprice),
			WholesaleMinQuantity:  sqlnull.ToInt32Ptr(v.Wholesaleminquantity),
			WholesaleDiscount:     sqlnull.ToInt64Ptr(v.Wholesalediscount),
			WholesaleDiscountType: sqlnull.ToStringPtr(v.Wholesalediscounttype),
			WeightGrams:           v.WeightGrams,
			IsVariantArchived:     v.Isvariantarchived,
		}

		if v.Isvariantarchived {
			archivedVariants = append(archivedVariants, variant)
		} else {
			validVariants = append(validVariants, variant)
		}
	}

	// ------------------------------------------------------------
	// Step 6: Final response

	first := variants[0]

	return &GetProductFullDetailResult{
		ProductID:         product.ID,
		Title:             product.Title,
		Description:       product.Description,
		CategoryID:        first.Categoryid,
		CategoryName:      first.Categoryname,
		SellerID:          user.ID,
		SellerStoreName:   first.Sellerstorename,
		ImageURLs:         first.Productimages,
		PrimaryImageURL:   primaryImageURL,
		PromoVideoURL:     sqlnull.ToStringPtr(first.Productpromovideourl),
		IsProductApproved: product.IsApproved,

		ValidVariants:    validVariants,
		ArchivedVariants: archivedVariants,
	}, nil
}
