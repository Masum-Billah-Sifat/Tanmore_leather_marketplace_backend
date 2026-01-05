package product_variant_snapshot

import (
	"github.com/google/uuid"
)

type ProductCreatedVariantPayload struct {
	VariantID             uuid.UUID `json:"variant_id"`
	Color                 string    `json:"color"`
	Size                  string    `json:"size"`
	RetailPrice           int64     `json:"retail_price"`
	InStock               bool      `json:"in_stock"`
	StockQuantity         int64     `json:"stock_quantity"`
	RetailDiscount        *int64    `json:"retail_discount"`
	RetailDiscountType    *string   `json:"retail_discount_type"`
	WholesalePrice        *int64    `json:"wholesale_price"`
	MinQtyWholesale       *int64    `json:"min_qty_wholesale"`
	WholesaleDiscount     *int64    `json:"wholesale_discount"`
	WholesaleDiscountType *string   `json:"wholesale_discount_type"`
	WeightGrams           int64     `json:"weight_grams"`
}

type ProductCreatedEvent struct {
	Product struct {
		ID              uuid.UUID `json:"id"`
		Title           string    `json:"title"`
		Description     string    `json:"description"`
		ImageURLs       []string  `json:"image_urls"`
		PrimaryImageURL string    `json:"primary_image_url"`
		PromoVideoURL   *string   `json:"promo_video_url"`
		IsApproved      bool      `json:"is_approved"`
		IsArchived      bool      `json:"is_archived"`
		IsBanned        bool      `json:"is_banned"`
	} `json:"product"`

	Category struct {
		ID         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		IsArchived bool      `json:"is_archived"`
	} `json:"category"`

	Seller struct {
		ID                      uuid.UUID `json:"id"`
		IsArchived              bool      `json:"is_archived"`
		IsBanned                bool      `json:"is_banned"`
		IsSellerProfileApproved bool      `json:"is_seller_profile_approved"`
		SellerStoreName         string    `json:"sellerstorename"`
	} `json:"seller"`

	Variants []ProductCreatedVariantPayload `json:"variants"`
}

type ProductInfoUpdatedEvent struct {
	UserID        uuid.UUID         `json:"user_id"`
	ProductID     uuid.UUID         `json:"product_id"`
	UpdatedFields map[string]string `json:"updated_fields"`
}

// ProductPrimaryImageSetEvent is emitted when a product image is set as primary
type ProductPrimaryImageSetEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	MediaID   uuid.UUID `json:"media_id"`
	MediaURL  string    `json:"media_url"`
}

// VariantCreatedEvent is emitted when a new product variant is created.
type VariantCreatedEvent struct {
	Product struct {
		ID              uuid.UUID `json:"id"`
		Title           string    `json:"title"`
		Description     string    `json:"description"`
		ImageURLs       []string  `json:"image_urls"`
		PrimaryImageURL string    `json:"primary_image_url"`
		PromoVideoURL   *string   `json:"promo_video_url"`
		IsApproved      bool      `json:"is_approved"`
		IsArchived      bool      `json:"is_archived"`
		IsBanned        bool      `json:"is_banned"`
	} `json:"product"`

	Category struct {
		ID         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		IsArchived bool      `json:"is_archived"`
	} `json:"category"`

	Seller struct {
		ID                      uuid.UUID `json:"id"`
		IsArchived              bool      `json:"is_archived"`
		IsBanned                bool      `json:"is_banned"`
		IsSellerProfileApproved bool      `json:"is_seller_profile_approved"`
		SellerStoreName         string    `json:"sellerstorename"`
	} `json:"seller"`

	Variant struct {
		ID                    uuid.UUID `json:"id"`
		Color                 string    `json:"color"`
		Size                  string    `json:"size"`
		RetailPrice           int64     `json:"retail_price"`
		InStock               bool      `json:"in_stock"`
		StockQuantity         int64     `json:"stock_quantity"`
		RetailDiscount        *int64    `json:"retail_discount"`
		RetailDiscountType    *string   `json:"retail_discount_type"`
		WholesalePrice        *int64    `json:"wholesale_price"`
		MinQtyWholesale       *int64    `json:"min_qty_wholesale"`
		WholesaleDiscount     *int64    `json:"wholesale_discount"`
		WholesaleDiscountType *string   `json:"wholesale_discount_type"`
		WeightGrams           int64     `json:"weight_grams"`
		HasRetailDiscount     bool      `json:"has_retail_discount"`
		HasWholesaleDiscount  bool      `json:"has_wholesale_discount"`
		WholesaleEnabled      bool      `json:"wholesale_enabled"`
	} `json:"variant"`
}

// --------------------
// variant.in_stock.updated
// --------------------
type VariantInStockUpdatedEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	VariantID uuid.UUID `json:"variant_id"`
	InStock   bool      `json:"in_stock"`
}

// --------------------
// variant.retail_price.updated
// --------------------
type VariantRetailPriceUpdatedEvent struct {
	UserID      uuid.UUID `json:"user_id"`
	ProductID   uuid.UUID `json:"product_id"`
	VariantID   uuid.UUID `json:"variant_id"`
	RetailPrice int64     `json:"retail_price"`
}

// --------------------
// variant.stock_quantity.updated
// --------------------
type VariantStockQuantityUpdatedEvent struct {
	UserID        uuid.UUID `json:"user_id"`
	ProductID     uuid.UUID `json:"product_id"`
	VariantID     uuid.UUID `json:"variant_id"`
	StockQuantity int64     `json:"stock_quantity"`
}

// --------------------
// variant.weight.updated
// --------------------
type VariantWeightUpdatedEvent struct {
	UserID      uuid.UUID `json:"user_id"`
	ProductID   uuid.UUID `json:"product_id"`
	VariantID   uuid.UUID `json:"variant_id"`
	WeightGrams int64     `json:"weight_grams"`
}

// --------------------
// variant.archived
// --------------------
type VariantArchivedEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	VariantID uuid.UUID `json:"variant_id"`
}

// --------------------
// variant.info.updated
// --------------------
type VariantInfoUpdatedEvent struct {
	UserID        uuid.UUID         `json:"user_id"`
	ProductID     uuid.UUID         `json:"product_id"`
	VariantID     uuid.UUID         `json:"variant_id"`
	UpdatedFields map[string]string `json:"updated_fields"` // keys can be "color", "size"
}

// --------------------
// variant.retail_discount.added
// --------------------
type VariantRetailDiscountAddedEvent struct {
	UserID             uuid.UUID `json:"user_id"`
	ProductID          uuid.UUID `json:"product_id"`
	VariantID          uuid.UUID `json:"variant_id"`
	RetailDiscount     int64     `json:"retail_discount"`
	RetailDiscountType string    `json:"retail_discount_type"` // flat | percentage
}

type VariantRetailDiscountUpdatedEvent struct {
	UserID        uuid.UUID              `json:"user_id"`
	ProductID     uuid.UUID              `json:"product_id"`
	VariantID     uuid.UUID              `json:"variant_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"` // may include: "retail_discount", "retail_discount_type"
}

type VariantRetailDiscountRemovedEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	VariantID uuid.UUID `json:"variant_id"`
}

type VariantWholesaleModeEnabledEvent struct {
	EventVersion          int       `json:"event_version"`
	UserID                uuid.UUID `json:"user_id"`
	ProductID             uuid.UUID `json:"product_id"`
	VariantID             uuid.UUID `json:"variant_id"`
	WholesalePrice        int64     `json:"wholesale_price"`
	MinQtyWholesale       int64     `json:"min_qty_wholesale"`
	HasWholesaleDiscount  bool      `json:"has_wholesale_discount"`
	WholesaleDiscount     *int64    `json:"wholesale_discount"`
	WholesaleDiscountType *string   `json:"wholesale_discount_type"`
}

type VariantWholesaleModeUpdatedEvent struct {
	UserID        uuid.UUID              `json:"user_id"`
	ProductID     uuid.UUID              `json:"product_id"`
	VariantID     uuid.UUID              `json:"variant_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
}

type VariantWholesaleModeDisabledEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	VariantID uuid.UUID `json:"variant_id"`
}

type VariantWholesaleDiscountAddedEvent struct {
	UserID                uuid.UUID `json:"user_id"`
	ProductID             uuid.UUID `json:"product_id"`
	VariantID             uuid.UUID `json:"variant_id"`
	WholesaleDiscount     int64     `json:"wholesale_discount"`
	WholesaleDiscountType string    `json:"wholesale_discount_type"`
}

type VariantWholesaleDiscountUpdatedEvent struct {
	UserID        uuid.UUID              `json:"user_id"`
	ProductID     uuid.UUID              `json:"product_id"`
	VariantID     uuid.UUID              `json:"variant_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
}

type VariantWholesaleDiscountRemovedEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	VariantID uuid.UUID `json:"variant_id"`
}

// 🔁 For "product_category_updated" event
type ProductCategoryUpdatedEvent struct {
	ProductID       uuid.UUID `json:"product_id"`
	SellerID        uuid.UUID `json:"seller_id"`
	NewCategoryID   uuid.UUID `json:"new_category_id"`
	NewCategoryName string    `json:"new_category_name"`
}

// 🔁 For "product.archived" event
type ProductArchivedEvent struct {
	SellerID  uuid.UUID `json:"seller_id"`
	ProductID uuid.UUID `json:"product_id"`
}
