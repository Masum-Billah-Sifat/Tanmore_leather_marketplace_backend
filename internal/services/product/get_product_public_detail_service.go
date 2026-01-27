package product

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_get_public_detail"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler

type GetProductPublicDetailInput struct {
	ProductID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Variant Response

type VariantDetail struct {
	VariantID             uuid.UUID `json:"variant_id"`
	Color                 string    `json:"color"`
	Size                  string    `json:"size"`
	InStock               bool      `json:"in_stock"`
	StockAmount           int32     `json:"stock_amount"`
	RetailPrice           int64     `json:"retail_price"`
	HasRetailDiscount     bool      `json:"has_retail_discount"`
	RetailDiscount        *int64    `json:"retail_discount,omitempty"`
	RetailDiscountType    *string   `json:"retail_discount_type,omitempty"`
	WholesaleEnabled      bool      `json:"wholesale_enabled"`
	WholesalePrice        *int64    `json:"wholesale_price,omitempty"`
	WholesaleMinQuantity  *int32    `json:"wholesale_min_quantity,omitempty"`
	WholesaleDiscount     *int64    `json:"wholesale_discount,omitempty"`
	WholesaleDiscountType *string   `json:"wholesale_discount_type,omitempty"`
	WeightGrams           int32     `json:"weight_grams"`
}

// ------------------------------------------------------------
// 📤 Final Output

type GetProductPublicDetailResult struct {
	ProductID       uuid.UUID       `json:"product_id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	CategoryID      uuid.UUID       `json:"category_id"`
	CategoryName    string          `json:"category_name"`
	SellerID        uuid.UUID       `json:"seller_id"`
	SellerStoreName string          `json:"seller_store_name"`
	Images          []string        `json:"images"`
	PromoVideoURL   *string         `json:"promo_video_url,omitempty"`
	Variants        []VariantDetail `json:"variants"`
}

// ------------------------------------------------------------
// 🧱 Dependencies

type GetProductPublicDetailServiceDeps struct {
	Repo repo.ProductGetPublicDetailRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service

type GetProductPublicDetailService struct {
	Deps GetProductPublicDetailServiceDeps
}

// 🚀 Constructor

func NewGetProductPublicDetailService(
	deps GetProductPublicDetailServiceDeps,
) *GetProductPublicDetailService {
	return &GetProductPublicDetailService{Deps: deps}
}

// ------------------------------------------------------------
// 🚀 Entrypoint

func (s *GetProductPublicDetailService) Start(
	ctx context.Context,
	input GetProductPublicDetailInput,
) (*GetProductPublicDetailResult, error) {

	// Step 1: Fetch all variant rows
	rows, err := s.Deps.Repo.GetProductDetailByProductID(ctx, input.ProductID)
	if err != nil {
		return nil, errors.NewTableError("product_variant_indexes.select", err.Error())
	}
	if len(rows) == 0 {
		return nil, errors.NewNotFoundError("product")
	}

	// Step 2: Validate moderation flags across all rows
	var r *sqlc.GetProductDetailByProductIDRow

	for _, row := range rows {
		if r == nil {
			r = &row
		}
		if row.Isproductbanned || row.Isproductarchived || !row.Isproductapproved {
			return nil, errors.NewServerError("Product not available")
		}
		if row.Issellerbanned || row.Issellerarchived || !row.Issellerapproved {
			return nil, errors.NewServerError("Seller not available")
		}
		if row.Iscategoryarchived {
			return nil, errors.NewServerError("Category not available")
		}
	}

	// Step 3: Filter variants
	var validVariants []VariantDetail
	for _, row := range rows {
		if row.Isvariantarchived {
			continue
		}
		v := VariantDetail{
			VariantID:             row.Variantid,
			Color:                 row.Color,
			Size:                  row.Size,
			InStock:               row.Isvariantinstock,
			StockAmount:           row.Stockamount,
			RetailPrice:           row.Retailprice,
			HasRetailDiscount:     row.HasRetailDiscount,
			RetailDiscount:        sqlnull.ToInt64Ptr(row.Retaildiscount),
			RetailDiscountType:    sqlnull.ToStringPtr(row.Retaildiscounttype),
			WholesaleEnabled:      row.Haswholesaleenabled,
			WholesalePrice:        sqlnull.ToInt64Ptr(row.Wholesaleprice),
			WholesaleMinQuantity:  sqlnull.ToInt32Ptr(row.Wholesaleminquantity),
			WholesaleDiscount:     sqlnull.ToInt64Ptr(row.Wholesalediscount),
			WholesaleDiscountType: sqlnull.ToStringPtr(row.Wholesalediscounttype),
			WeightGrams:           row.WeightGrams,
		}
		validVariants = append(validVariants, v)
	}

	if len(validVariants) == 0 {
		return nil, errors.NewNotFoundError("No available variants")
	}

	// Step 4: Return response
	product := &GetProductPublicDetailResult{
		ProductID:       r.Productid,
		Title:           r.Producttitle,
		Description:     r.Productdescription,
		CategoryID:      r.Categoryid,
		CategoryName:    r.Categoryname,
		SellerID:        r.Sellerid,
		SellerStoreName: r.Sellerstorename,
		Images:          r.Productimages,
		PromoVideoURL:   sqlnull.ToStringPtr(r.Productpromovideourl),
		Variants:        validVariants,
	}

	return product, nil
}
