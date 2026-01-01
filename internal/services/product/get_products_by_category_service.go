// ------------------------------------------------------------
// 📁 File: internal/services/product/get_products_by_category_service.go
// 🧠 Handles GET /api/category-products
//     - Validates category existence and status
//     - Resolves all leaf categories under given root (if needed)
//     - Fetches product variant indexes
//     - Groups result by product
//     - Returns formatted response for handler

package product

import (
	"context"

	repo "tanmore_backend/internal/repository/product/fetch_by_category"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type GetProductsByCategoryInput struct {
	CategoryID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Variant Struct
type CategoryProductVariant struct {
	VariantID             uuid.UUID `json:"variant_id"`
	Color                 string    `json:"color"`
	Size                  string    `json:"size"`
	IsInStock             bool      `json:"is_in_stock"`
	StockAmount           int32     `json:"stock_amount"`
	RetailPrice           int64     `json:"retail_price"`
	RetailDiscount        int64     `json:"retail_discount"`
	RetailDiscountType    string    `json:"retail_discount_type"`
	HasRetailDiscount     bool      `json:"has_retail_discount"`
	HasWholesaleEnabled   bool      `json:"has_wholesale_enabled"`
	WholesalePrice        int64     `json:"wholesale_price"`
	WholesaleMinQuantity  int32     `json:"wholesale_min_quantity"`
	WholesaleDiscount     int64     `json:"wholesale_discount"`
	WholesaleDiscountType string    `json:"wholesale_discount_type"`
	WeightGrams           int32     `json:"weight_grams"`
}

// ------------------------------------------------------------
// 📤 Product Struct
type CategoryProductResponse struct {
	ProductID       uuid.UUID                `json:"product_id"`
	CategoryID      uuid.UUID                `json:"category_id"`
	CategoryName    string                   `json:"category_name"`
	SellerID        uuid.UUID                `json:"seller_id"`
	SellerStoreName string                   `json:"seller_store_name"`
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	ImageURLs       []string                 `json:"image_urls"`
	PromoVideoURL   string                   `json:"promo_video_url"`
	Variants        []CategoryProductVariant `json:"variants"`
}

// ------------------------------------------------------------
// 📤 Final Result
type GetProductsByCategoryResult struct {
	Status string                    `json:"status"`
	Data   []CategoryProductResponse `json:"data"`
}

// ------------------------------------------------------------
// 🧱 Dependencies
type GetProductsByCategoryServiceDeps struct {
	Repo repo.FetchByCategoryRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type GetProductsByCategoryService struct {
	Deps GetProductsByCategoryServiceDeps
}

// 🚀 Constructor
func NewGetProductsByCategoryService(deps GetProductsByCategoryServiceDeps) *GetProductsByCategoryService {
	return &GetProductsByCategoryService{Deps: deps}
}

// 🚀 Entrypoint
func (s *GetProductsByCategoryService) Start(
	ctx context.Context,
	input GetProductsByCategoryInput,
) (*GetProductsByCategoryResult, error) {

	// Step 1: Validate category
	category, err := s.Deps.Repo.GetCategoryByID(ctx, input.CategoryID)
	if err != nil {
		return nil, errors.NewNotFoundError("category")
	}
	if category.IsArchived {
		return nil, errors.NewValidationError("category", "is archived")
	}

	// Step 2: Resolve leaf category IDs
	var leafCategoryIDs []uuid.UUID
	if category.IsLeaf {
		leafCategoryIDs = []uuid.UUID{category.ID}
	} else {
		leafCategoryIDs, err = s.Deps.Repo.GetAllLeafCategoryIDsByRoot(ctx, category.ID)
		if err != nil {
			return nil, errors.NewTableError("category.leaf_query", err.Error())
		}
		if len(leafCategoryIDs) == 0 {
			return &GetProductsByCategoryResult{
				Status: "success",
				Data:   []CategoryProductResponse{},
			}, nil
		}
	}

	// Step 3: Fetch product variant index data
	rows, err := s.Deps.Repo.GetProductVariantIndexesByCategoryIDs(ctx, leafCategoryIDs)
	if err != nil {
		return nil, errors.NewTableError("product_variant_indexes.select", err.Error())
	}

	// Step 4: Group by product_id
	productMap := make(map[uuid.UUID]*CategoryProductResponse)

	for _, row := range rows {
		productID := row.Productid

		if _, exists := productMap[productID]; !exists {
			productMap[productID] = &CategoryProductResponse{
				ProductID:       productID,
				CategoryID:      row.Categoryid,
				CategoryName:    row.Categoryname,
				SellerID:        row.Sellerid,
				SellerStoreName: row.Sellerstorename,
				Title:           row.Producttitle,
				Description:     row.Productdescription,
				ImageURLs:       row.Productimages,
				PromoVideoURL:   sqlnull.StringOrEmpty(row.Productpromovideourl), // ✅ null-safe
				Variants:        []CategoryProductVariant{},
			}
		}

		productMap[productID].Variants = append(productMap[productID].Variants, CategoryProductVariant{
			VariantID:             row.Variantid,
			Color:                 row.Color,
			Size:                  row.Size,
			IsInStock:             row.Isvariantinstock,
			StockAmount:           row.Stockamount,
			RetailPrice:           row.Retailprice,
			RetailDiscount:        sqlnull.Int64OrZero(row.Retaildiscount),
			RetailDiscountType:    sqlnull.StringOrEmpty(row.Retaildiscounttype),
			HasRetailDiscount:     row.HasRetailDiscount,
			HasWholesaleEnabled:   row.Haswholesaleenabled,
			WholesalePrice:        sqlnull.Int64OrZero(row.Wholesaleprice),
			WholesaleMinQuantity:  sqlnull.Int32OrZero(row.Wholesaleminquantity),
			WholesaleDiscount:     sqlnull.Int64OrZero(row.Wholesalediscount),
			WholesaleDiscountType: sqlnull.StringOrEmpty(row.Wholesalediscounttype),
			WeightGrams:           row.WeightGrams,
		})
	}

	// Step 5: Convert map to flat slice
	var products []CategoryProductResponse
	for _, p := range productMap {
		products = append(products, *p)
	}

	return &GetProductsByCategoryResult{
		Status: "success",
		Data:   products,
	}, nil
}
