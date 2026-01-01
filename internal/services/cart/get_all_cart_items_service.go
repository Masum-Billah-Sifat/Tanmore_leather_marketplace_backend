// ------------------------------------------------------------
// 📁 File: internal/services/cart/get_all_cart_items_service.go
// 🧠 Handles retrieval of active cart items enriched with variant snapshot data.
//     - Validates customer (not banned/archived)
//     - Fetches all active variant IDs in user's cart
//     - Runs enriched JOIN query against snapshot table
//     - Filters into valid and invalid items
//     - Groups valid items by seller → product → variants
//     - Returns grouped valid_items and flat invalid_items array

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/get_all_cart_items"
	"tanmore_backend/pkg/errors"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type GetAllCartItemsInput struct {
	UserID uuid.UUID
}

// ------------------------------------------------------------
// 🧤 Variant representation under a product
type CartVariantItem struct {
	VariantID             uuid.UUID `json:"variant_id"`
	Color                 string    `json:"color"`
	Size                  string    `json:"size"`
	RetailPrice           int64     `json:"retail_price"`
	HasRetailDiscount     bool      `json:"has_retail_discount"`
	RetailDiscount        int64     `json:"retail_discount"`
	RetailDiscountType    string    `json:"retail_discount_type"`
	HasWholesaleEnabled   bool      `json:"has_wholesale_enabled"`
	WholesalePrice        int64     `json:"wholesale_price"`
	WholesaleMinQty       int32     `json:"wholesale_min_qty"`
	HasWholesaleDiscount  bool      `json:"has_wholesale_discount"`
	WholesaleDiscount     int64     `json:"wholesale_discount"`
	WholesaleDiscountType string    `json:"wholesale_discount_type"`
	WeightGrams           int32     `json:"weight_grams"`
	QuantityInCart        int32     `json:"quantity_in_cart"`
}

// ------------------------------------------------------------
// 🛍 Product grouping per seller
type CartProductItem struct {
	ProductID           uuid.UUID         `json:"product_id"`
	CategoryName        string            `json:"category_name"`
	ProductTitle        string            `json:"product_title"`
	ProductDescription  string            `json:"product_description"`
	ProductPrimaryImage string            `json:"product_primary_image_url"`
	Variants            []CartVariantItem `json:"variants"`
}

type CartGroupedBySeller struct {
	SellerID  uuid.UUID
	StoreName string
	Products  []*CartProductItem
}

// ------------------------------------------------------------
// ❌ Invalid item representation
type InvalidCartItem struct {
	VariantID    uuid.UUID `json:"variant_id"`
	Reason       string    `json:"reason"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductTitle string    `json:"product_title"`
	Color        string    `json:"color"`
	Size         string    `json:"size"`
}

// ------------------------------------------------------------
// 📤 Final response
type GetAllCartItemsResult struct {
	ValidItems   []CartGroupedBySeller `json:"valid_items"`
	InvalidItems []InvalidCartItem     `json:"invalid_items"`
}

// ------------------------------------------------------------
// 🧱 Dependencies
type GetAllCartItemsServiceDeps struct {
	Repo repo.GetAllCartItemsRepoInterface
}

// 🛠️ Service Definition
type GetAllCartItemsService struct {
	Deps GetAllCartItemsServiceDeps
}

// 🚀 Constructor
func NewGetAllCartItemsService(deps GetAllCartItemsServiceDeps) *GetAllCartItemsService {
	return &GetAllCartItemsService{Deps: deps}
}

func (s *GetAllCartItemsService) Start(
	ctx context.Context,
	input GetAllCartItemsInput,
) (*GetAllCartItemsResult, error) {
	// Step 1: Validate user moderation
	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}
	if user.IsArchived {
		return nil, errors.NewAuthError("user is archived")
	}
	if user.IsBanned {
		return nil, errors.NewAuthError("user is banned")
	}

	// Step 2: Fetch active variant IDs in cart
	activeVariantIDs, err := s.Deps.Repo.ListActiveVariantIDsByUser(ctx, input.UserID)
	if err != nil {
		return nil, errors.NewTableError("cart_items", "cannot list variant IDs")
	}
	if len(activeVariantIDs) == 0 {
		return &GetAllCartItemsResult{
			ValidItems:   []CartGroupedBySeller{},
			InvalidItems: []InvalidCartItem{},
		}, nil
	}

	// Step 3: Enriched join query
	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
			UserID:     input.UserID,
			VariantIds: activeVariantIDs,
		},
	)
	if err != nil {
		return nil, errors.NewServerError("cannot fetch snapshot enriched cart items")
	}

	// Step 4: Process rows
	grouped := make(map[uuid.UUID]*CartGroupedBySeller)
	productMap := make(map[uuid.UUID]map[uuid.UUID]*CartProductItem)
	var invalidItems []InvalidCartItem

	for _, row := range rows {
		// Moderation filters
		if !row.Issellerapproved || row.Issellerarchived || row.Issellerbanned ||
			!row.Isproductapproved || row.Isproductarchived || row.Isproductbanned ||
			row.Isvariantarchived || !row.Isvariantinstock {
			invalidItems = append(invalidItems, InvalidCartItem{
				VariantID:    row.CartVariantID,
				Reason:       "variant unavailable due to moderation or stock",
				ProductID:    row.Productid,
				ProductTitle: row.Producttitle,
				Color:        row.Color,
				Size:         row.Size,
			})
			continue
		}

		// Create seller group if not exists
		if _, ok := grouped[row.Sellerid]; !ok {
			grouped[row.Sellerid] = &CartGroupedBySeller{
				SellerID:  row.Sellerid,
				StoreName: row.Sellerstorename,
				Products:  []*CartProductItem{},
			}
			productMap[row.Sellerid] = make(map[uuid.UUID]*CartProductItem)
		}

		// Create product group if not exists
		if _, ok := productMap[row.Sellerid][row.Productid]; !ok {
			product := &CartProductItem{
				ProductID:           row.Productid,
				CategoryName:        row.Categoryname,
				ProductTitle:        row.Producttitle,
				ProductDescription:  row.Productdescription,
				ProductPrimaryImage: row.Productprimaryimageurl,
				Variants:            []CartVariantItem{},
			}
			productMap[row.Sellerid][row.Productid] = product
			// 🔁 FIXED: Append pointer instead of dereferenced copy
			grouped[row.Sellerid].Products = append(grouped[row.Sellerid].Products, product)
		}

		// ---- Retail discount
		var retailDiscount int64
		var retailDiscountType string
		if row.Retaildiscount.Valid {
			retailDiscount = row.Retaildiscount.Int64
		}
		if row.Retaildiscounttype.Valid {
			retailDiscountType = row.Retaildiscounttype.String
		}

		// ---- Wholesale fields
		var wholesalePrice int64
		var wholesaleMinQty int32
		var wholesaleDiscount int64
		var wholesaleDiscountType string
		if row.Wholesaleprice.Valid {
			wholesalePrice = row.Wholesaleprice.Int64
		}
		if row.Wholesaleminquantity.Valid {
			wholesaleMinQty = row.Wholesaleminquantity.Int32
		}
		if row.Wholesalediscount.Valid {
			wholesaleDiscount = row.Wholesalediscount.Int64
		}
		if row.Wholesalediscounttype.Valid {
			wholesaleDiscountType = row.Wholesalediscounttype.String
		}

		// ---- Cart quantity
		var quantityInCart int32
		if row.CartRequiredQuantity.Valid {
			quantityInCart = row.CartRequiredQuantity.Int32
		}

		variant := CartVariantItem{
			VariantID:             row.Variantid,
			Color:                 row.Color,
			Size:                  row.Size,
			RetailPrice:           row.Retailprice,
			HasRetailDiscount:     row.Hasretaildiscount,
			RetailDiscount:        retailDiscount,
			RetailDiscountType:    retailDiscountType,
			HasWholesaleEnabled:   row.Haswholesaleenabled,
			WholesalePrice:        wholesalePrice,
			WholesaleMinQty:       wholesaleMinQty,
			HasWholesaleDiscount:  row.Haswholesalediscount,
			WholesaleDiscount:     wholesaleDiscount,
			WholesaleDiscountType: wholesaleDiscountType,
			WeightGrams:           row.WeightGrams,
			QuantityInCart:        quantityInCart,
		}

		// ✅ Add variant to product (via pointer)
		prodPtr := productMap[row.Sellerid][row.Productid]
		prodPtr.Variants = append(prodPtr.Variants, variant)
	}

	// Collect final grouped valid items
	var validItems []CartGroupedBySeller
	for _, seller := range grouped {
		validItems = append(validItems, *seller)
	}

	return &GetAllCartItemsResult{
		ValidItems:   validItems,
		InvalidItems: invalidItems,
	}, nil
}
