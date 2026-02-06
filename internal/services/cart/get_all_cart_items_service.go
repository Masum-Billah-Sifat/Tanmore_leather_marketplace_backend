// ------------------------------------------------------------
// 📁 File: internal/services/cart/get_all_cart_items_service.go
// 🧠 Handles retrieval of active cart items enriched with variant snapshot data.
// 🔒 AUTHENTICATED USERS ONLY

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/get_all_cart_items"
	"tanmore_backend/pkg/errors"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler (AUTH ONLY)
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
	SellerID  uuid.UUID          `json:"seller_id"`
	StoreName string             `json:"store_name"`
	Products  []*CartProductItem `json:"products"`
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

// ------------------------------------------------------------
// 🛠️ Service Definition
type GetAllCartItemsService struct {
	Deps GetAllCartItemsServiceDeps
}

// 🚀 Constructor
func NewGetAllCartItemsService(deps GetAllCartItemsServiceDeps) *GetAllCartItemsService {
	return &GetAllCartItemsService{Deps: deps}
}

// ------------------------------------------------------------
// 🚀 Entrypoint
func (s *GetAllCartItemsService) Start(
	ctx context.Context,
	input GetAllCartItemsInput,
) (*GetAllCartItemsResult, error) {

	// ------------------------------------------------------------
	// Step 1: Validate user moderation (MANDATORY)
	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, errors.ErrAuthUserNotFound()
	}
	if user.IsArchived {
		return nil, errors.ErrAuthArchivedUser()
	}
	if user.IsBanned {
		return nil, errors.ErrAuthBannedUser()
	}
	// ------------------------------------------------------------
	// Step 2: List active variant IDs in cart
	variantIDs, err := s.Deps.Repo.ListActiveVariantIDsByUser(ctx, input.UserID)
	if err != nil {
		return nil, errors.NewTableError("cart_items", "cannot list variant IDs")
	}

	if len(variantIDs) == 0 {
		return &GetAllCartItemsResult{
			ValidItems:   []CartGroupedBySeller{},
			InvalidItems: []InvalidCartItem{},
		}, nil
	}

	// ------------------------------------------------------------
	// Step 3: Fetch enriched snapshot rows
	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx,
		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
			UserID:     input.UserID,
			VariantIds: variantIDs,
		},
	)
	if err != nil {
		return nil, errors.NewServerError("cannot fetch snapshot enriched cart items")
	}

	// ------------------------------------------------------------
	// Step 4: Group + filter
	grouped := make(map[uuid.UUID]*CartGroupedBySeller)
	productMap := make(map[uuid.UUID]map[uuid.UUID]*CartProductItem)
	var invalidItems []InvalidCartItem

	for _, row := range rows {

		// ❌ Moderation / availability failure
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

		// Seller group
		if _, ok := grouped[row.Sellerid]; !ok {
			grouped[row.Sellerid] = &CartGroupedBySeller{
				SellerID:  row.Sellerid,
				StoreName: row.Sellerstorename,
				Products:  []*CartProductItem{},
			}
			productMap[row.Sellerid] = make(map[uuid.UUID]*CartProductItem)
		}

		// Product group
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
			grouped[row.Sellerid].Products = append(grouped[row.Sellerid].Products, product)
		}

		// Variant
		variant := CartVariantItem{
			VariantID:             row.Variantid,
			Color:                 row.Color,
			Size:                  row.Size,
			RetailPrice:           row.Retailprice,
			HasRetailDiscount:     row.Hasretaildiscount,
			RetailDiscount:        row.Retaildiscount.Int64,
			RetailDiscountType:    row.Retaildiscounttype.String,
			HasWholesaleEnabled:   row.Haswholesaleenabled,
			WholesalePrice:        row.Wholesaleprice.Int64,
			WholesaleMinQty:       row.Wholesaleminquantity.Int32,
			HasWholesaleDiscount:  row.Haswholesalediscount,
			WholesaleDiscount:     row.Wholesalediscount.Int64,
			WholesaleDiscountType: row.Wholesalediscounttype.String,
			WeightGrams:           row.WeightGrams,
			QuantityInCart:        row.CartRequiredQuantity.Int32,
		}

		productMap[row.Sellerid][row.Productid].Variants =
			append(productMap[row.Sellerid][row.Productid].Variants, variant)
	}

	// ------------------------------------------------------------
	// Step 5: Final response
	var validItems []CartGroupedBySeller
	for _, seller := range grouped {
		validItems = append(validItems, *seller)
	}

	return &GetAllCartItemsResult{
		ValidItems:   validItems,
		InvalidItems: invalidItems,
	}, nil
}
