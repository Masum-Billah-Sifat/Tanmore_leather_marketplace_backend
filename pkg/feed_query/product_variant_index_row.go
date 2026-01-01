// 🧱 Struct representing a single row from product_variant_indexes
package feedquery

import (
	"database/sql"
	// feedquery "tanmore_backend/pkg/feed_query"
	// feedquery "tanmore_backend/pkg/feed_query"
	// feedquery "tanmore_backend/pkg/feed_query"
	"time"

	"github.com/google/uuid"
)

type ProductVariantIndexFeedOrSearchRow struct {
	// Primary key
	ID uuid.UUID

	// Category
	CategoryID         uuid.UUID
	IsCategoryArchived bool
	CategoryName       string

	// Seller
	SellerID         uuid.UUID
	IsSellerApproved bool
	IsSellerArchived bool
	IsSellerBanned   bool
	SellerStoreName  string

	// Product
	ProductID          uuid.UUID
	IsProductApproved  bool
	IsProductArchived  bool
	IsProductBanned    bool
	ProductTitle       string
	ProductDescription string

	// before
	// ProductImages []string
	// after
	ProductImages StringArray

	ProductPromoVideoURL sql.NullString

	// Variant
	VariantID         uuid.UUID
	IsVariantArchived bool
	IsVariantInStock  bool
	StockAmount       int32
	Color             string
	Size              string
	RetailPrice       int64

	RetailDiscountType sql.NullString
	RetailDiscount     sql.NullInt64
	HasRetailDiscount  bool

	HasWholesaleEnabled   bool
	WholesalePrice        sql.NullInt64
	WholesaleMinQuantity  sql.NullInt32
	WholesaleDiscountType sql.NullString
	WholesaleDiscount     sql.NullInt64

	WeightGrams int32

	// Feed/Search metadata
	Views          int64
	RelevanceScore sql.NullFloat64

	CreatedAt time.Time
	UpdatedAt time.Time
}
