// 📁 File: internal/services/product/product_variant_index_feed_or_search/feed_query_service.go
// 🧠 Handles product feed and search queries dynamically with optional filters and Amazon-style ranking

package product_variant_index_feed_or_search

import (
	"context"
	"database/sql"
	"fmt"

	feedquery "tanmore_backend/internal/repository/product/product_variant/product_variant_index_feed_or_search"
	"tanmore_backend/pkg/errors"

	pkg_feed_query "tanmore_backend/pkg/feed_query"

	"github.com/google/uuid"
)

// 80 84 100 120 149 in this lines errors show up

// 📦 Input structure from handler
type FeedQueryInput struct {
	Page              int
	PerPage           int
	Query             *string
	CategoryID        *uuid.UUID
	MinPrice          *int64
	MaxPrice          *int64
	MinWeight         *int
	MaxWeight         *int
	Color             *string
	Size              *string
	InStock           *bool
	HasRetailDiscount *bool
	OnlyWholesale     *bool
	Sort              *string
}

// type FeedQueryResult struct {
// 	Page       int           `json:"page"`
// 	PerPage    int           `json:"per_page"`
// 	TotalItems int           `json:"total_items"`
// 	Sellers    []SellerGroup `json:"sellers"`
// }

type FeedQueryResult struct {
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalItems int            `json:"total_items"`
	Products   []ProductGroup `json:"products"` // 🆕 updated
}

type SellerGroup struct {
	SellerID        uuid.UUID                    `json:"seller_id"`
	SellerStoreName string                       `json:"seller_store_name"`
	CategoriesList  []CategoryGroup              `json:"categories"`
	Categories      map[uuid.UUID]*CategoryGroup `json:"-"` // internal only
}

type CategoryGroup struct {
	CategoryID   uuid.UUID                   `json:"category_id"`
	CategoryName string                      `json:"category_name"`
	ProductsList []ProductGroup              `json:"products"`
	Products     map[uuid.UUID]*ProductGroup `json:"-"` // internal only
}

// type ProductGroup struct {
// 	ProductID     uuid.UUID      `json:"product_id"`
// 	Title         string         `json:"product_title"`
// 	Description   string         `json:"product_description"`
// 	Images        []string       `json:"product_images"`
// 	PromoVideoURL *string        `json:"product_promo_video_url,omitempty"`
// 	Variants      []VariantGroup `json:"variants"`
// }

type ProductGroup struct {
	ProductID     uuid.UUID `json:"product_id"`
	Title         string    `json:"product_title"`
	Description   string    `json:"product_description"`
	Images        []string  `json:"product_images"`
	PromoVideoURL *string   `json:"product_promo_video_url,omitempty"`

	SellerID        uuid.UUID `json:"seller_id"`
	SellerStoreName string    `json:"seller_store_name"`
	CategoryID      uuid.UUID `json:"category_id"`
	CategoryName    string    `json:"category_name"`

	Variants []VariantGroup `json:"variants"`
}

type VariantGroup struct {
	VariantID          uuid.UUID `json:"variant_id"`
	Color              string    `json:"color"`
	Size               string    `json:"size"`
	RetailPrice        int64     `json:"retail_price"`
	RetailDiscount     *int64    `json:"retail_discount,omitempty"`
	RetailDiscountType *string   `json:"retail_discount_type,omitempty"`
	HasRetailDiscount  bool      `json:"has_retail_discount"`

	WholesaleEnabled      bool    `json:"wholesale_enabled"`
	WholesalePrice        *int64  `json:"wholesale_price,omitempty"`
	WholesaleMinQty       *int    `json:"wholesale_min_qty,omitempty"`
	WholesaleDiscount     *int64  `json:"wholesale_discount,omitempty"`
	WholesaleDiscountType *string `json:"wholesale_discount_type,omitempty"`

	WeightGrams    int      `json:"weight_grams"`
	RelevanceScore *float64 `json:"relevance_score,omitempty"`
}

// 🧠 Service struct
type FeedQueryService struct {
	Repo *feedquery.ProductVariantIndexFeedOrSearchRepository
}

// 🏗️ Constructor
func NewFeedQueryService(repo *feedquery.ProductVariantIndexFeedOrSearchRepository) *FeedQueryService {
	return &FeedQueryService{Repo: repo}
}

// 🚀 Executes the feed/search logic (refined version)
func (s *FeedQueryService) Start(ctx context.Context, input FeedQueryInput) (*FeedQueryResult, error) {
	// 1️⃣ Fallback pagination
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 10
	}
	offset := (input.Page - 1) * input.PerPage

	// 2️⃣ Build dynamic filter params
	params := pkg_feed_query.FeedQueryParams{
		Q:                 input.Query,
		CategoryID:        input.CategoryID,
		MinPrice:          input.MinPrice,
		MaxPrice:          input.MaxPrice,
		MinWeight:         intToInt64Ptr(input.MinWeight),
		MaxWeight:         intToInt64Ptr(input.MaxWeight),
		Color:             input.Color,
		Size:              input.Size,
		InStock:           input.InStock,
		HasRetailDiscount: input.HasRetailDiscount,
		OnlyWholesale:     input.OnlyWholesale,
		Sort:              input.Sort,
		Limit:             input.PerPage,
		Offset:            offset,
	}

	// 3️⃣ Run query
	rows, err := s.Repo.RunFeedOrSearchQuery(ctx, params)
	if err != nil {
		return nil, errors.NewServerError(fmt.Sprintf("feed_query.repo: %s", err.Error()))
	}

	// ✅ Debug log
	fmt.Printf("🧪 [DEBUG] FeedQueryService: %d rows returned from query\n", len(rows))

	// 4️⃣ Group by product → variants
	grouped := map[uuid.UUID]*ProductGroup{}
	totalVariants := 0

	for _, row := range rows {
		totalVariants++

		product, ok := grouped[row.ProductID]
		if !ok {
			product = &ProductGroup{
				ProductID:     row.ProductID,
				Title:         row.ProductTitle,
				Description:   row.ProductDescription,
				Images:        row.ProductImages,
				PromoVideoURL: nullTextToPtr(row.ProductPromoVideoURL),
				Variants:      []VariantGroup{},

				// 🆕 Enriched metadata for feed display
				SellerID:        row.SellerID,
				SellerStoreName: row.SellerStoreName,
				CategoryID:      row.CategoryID,
				CategoryName:    row.CategoryName,
			}
			grouped[row.ProductID] = product
		}

		// Add variant to product
		variant := VariantGroup{
			VariantID:          row.VariantID,
			Color:              row.Color,
			Size:               row.Size,
			RetailPrice:        row.RetailPrice,
			RetailDiscount:     nullIntToPtr(row.RetailDiscount),
			RetailDiscountType: nullTextToPtr(row.RetailDiscountType),
			HasRetailDiscount:  row.HasRetailDiscount,

			WholesaleEnabled:      row.HasWholesaleEnabled,
			WholesalePrice:        nullIntToPtr(row.WholesalePrice),
			WholesaleMinQty:       nullInt32ToIntPtr(row.WholesaleMinQuantity),
			WholesaleDiscount:     nullIntToPtr(row.WholesaleDiscount),
			WholesaleDiscountType: nullTextToPtr(row.WholesaleDiscountType),

			WeightGrams:    int(row.WeightGrams),
			RelevanceScore: nullFloatToPtr(row.RelevanceScore),
		}
		product.Variants = append(product.Variants, variant)
	}

	// 5️⃣ Flatten for final response
	products := []ProductGroup{}
	for _, prod := range grouped {
		products = append(products, *prod)
	}

	return &FeedQueryResult{
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalItems: totalVariants,
		Products:   products, // ⚠️ Rename `Sellers` field if needed
	}, nil
}

// // 🚀 Executes the feed/search logic
// func (s *FeedQueryService) Start(ctx context.Context, input FeedQueryInput) (*FeedQueryResult, error) {
// 	// 1️⃣ Fallback pagination
// 	if input.Page <= 0 {
// 		input.Page = 1
// 	}
// 	if input.PerPage <= 0 {
// 		input.PerPage = 10
// 	}
// 	offset := (input.Page - 1) * input.PerPage

// 	// 2️⃣ Build dynamic filter params
// 	params := pkg_feed_query.FeedQueryParams{
// 		Q:                 input.Query,
// 		CategoryID:        input.CategoryID,
// 		MinPrice:          input.MinPrice,
// 		MaxPrice:          input.MaxPrice,
// 		MinWeight:         intToInt64Ptr(input.MinWeight),
// 		MaxWeight:         intToInt64Ptr(input.MaxWeight),
// 		Color:             input.Color,
// 		Size:              input.Size,
// 		InStock:           input.InStock,
// 		HasRetailDiscount: input.HasRetailDiscount,
// 		OnlyWholesale:     input.OnlyWholesale,
// 		Sort:              input.Sort,
// 		Limit:             input.PerPage,
// 		Offset:            offset,
// 	}

// 	// 3️⃣ Run query
// 	rows, err := s.Repo.RunFeedOrSearchQuery(ctx, params)
// 	if err != nil {
// 		return nil, errors.NewServerError(fmt.Sprintf("feed_query.repo: %s", err.Error()))
// 	}

// 	// ✅ Debug log
// 	fmt.Printf("🧪 [DEBUG] FeedQueryService: %d rows returned from query\n", len(rows))

// 	// 4️⃣ Group data: Seller → Category → Product → Variants
// 	grouped := map[uuid.UUID]*SellerGroup{}
// 	totalVariants := 0

// 	for _, row := range rows {
// 		totalVariants++

// 		// --- Seller level
// 		seller, ok := grouped[row.SellerID]
// 		if !ok {
// 			seller = &SellerGroup{
// 				SellerID:        row.SellerID,
// 				SellerStoreName: row.SellerStoreName,
// 				Categories:      map[uuid.UUID]*CategoryGroup{},
// 			}
// 			grouped[row.SellerID] = seller
// 		}

// 		// --- Category level
// 		category, ok := seller.Categories[row.CategoryID]
// 		if !ok {
// 			category = &CategoryGroup{
// 				CategoryID:   row.CategoryID,
// 				CategoryName: row.CategoryName,
// 				Products:     map[uuid.UUID]*ProductGroup{},
// 			}
// 			seller.Categories[row.CategoryID] = category
// 		}

// 		// --- Product level
// 		product, ok := category.Products[row.ProductID]
// 		if !ok {
// 			product = &ProductGroup{
// 				ProductID:     row.ProductID,
// 				Title:         row.ProductTitle,
// 				Description:   row.ProductDescription,
// 				Images:        row.ProductImages,
// 				PromoVideoURL: nullTextToPtr(row.ProductPromoVideoURL),
// 				Variants:      []VariantGroup{},
// 			}
// 			category.Products[row.ProductID] = product
// 		}

// 		// --- Variant level
// 		variant := VariantGroup{
// 			VariantID:          row.VariantID,
// 			Color:              row.Color,
// 			Size:               row.Size,
// 			RetailPrice:        row.RetailPrice,
// 			RetailDiscount:     nullIntToPtr(row.RetailDiscount),
// 			RetailDiscountType: nullTextToPtr(row.RetailDiscountType),
// 			HasRetailDiscount:  row.HasRetailDiscount,

// 			WholesaleEnabled:      row.HasWholesaleEnabled,
// 			WholesalePrice:        nullIntToPtr(row.WholesalePrice),
// 			WholesaleMinQty:       nullInt32ToIntPtr(row.WholesaleMinQuantity),
// 			WholesaleDiscount:     nullIntToPtr(row.WholesaleDiscount),
// 			WholesaleDiscountType: nullTextToPtr(row.WholesaleDiscountType),

// 			WeightGrams:    int(row.WeightGrams),
// 			RelevanceScore: nullFloatToPtr(row.RelevanceScore),
// 		}

// 		product.Variants = append(product.Variants, variant)
// 	}

// 	// 5️⃣ Flatten final grouped structure for JSON response
// 	sellers := []SellerGroup{}
// 	for _, seller := range grouped {
// 		cats := []CategoryGroup{}
// 		for _, cat := range seller.Categories {
// 			prods := []ProductGroup{}
// 			for _, prod := range cat.Products {
// 				prods = append(prods, *prod)
// 			}
// 			cat.ProductsList = prods
// 			cats = append(cats, *cat)
// 		}
// 		seller.CategoriesList = cats
// 		sellers = append(sellers, *seller)
// 	}

// 	return &FeedQueryResult{
// 		Page:       input.Page,
// 		PerPage:    input.PerPage,
// 		TotalItems: totalVariants,
// 		Sellers:    sellers,
// 	}, nil
// }

func nullTextToPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func nullIntToPtr(i sql.NullInt64) *int64 {
	if i.Valid {
		return &i.Int64
	}
	return nil
}

func nullFloatToPtr(f sql.NullFloat64) *float64 {
	if f.Valid {
		return &f.Float64
	}
	return nil
}

// func uuidPtrToStr(id *uuid.UUID) *string {
// 	if id == nil {
// 		return nil
// 	}
// 	str := id.String()
// 	return &str
// }

func nullInt32ToIntPtr(i sql.NullInt32) *int {
	if i.Valid {
		v := int(i.Int32)
		return &v
	}
	return nil
}

func intToInt64Ptr(v *int) *int64 {
	if v == nil {
		return nil
	}
	val := int64(*v)
	return &val
}
