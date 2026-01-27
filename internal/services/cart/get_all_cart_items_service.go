// ------------------------------------------------------------
// 📁 File: internal/services/cart/get_all_cart_items_service.go
// 🧠 Handles retrieval of active cart items enriched with variant snapshot data.

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/get_all_cart_items"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// 📥 Input from handler
type GetAllCartItemsInput struct {
	UserID      *uuid.UUID
	GuestUserID *uuid.UUID
}

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

// ❌ Invalid item representation
type InvalidCartItem struct {
	VariantID    uuid.UUID `json:"variant_id"`
	Reason       string    `json:"reason"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductTitle string    `json:"product_title"`
	Color        string    `json:"color"`
	Size         string    `json:"size"`
}

// 📤 Final response
type GetAllCartItemsResult struct {
	ValidItems   []CartGroupedBySeller `json:"valid_items"`
	InvalidItems []InvalidCartItem     `json:"invalid_items"`
}

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

// 🚀 Entrypoint
func (s *GetAllCartItemsService) Start(
	ctx context.Context,
	input GetAllCartItemsInput,
) (*GetAllCartItemsResult, error) {

	// Step 1: Validate user moderation (ONLY if logged in)
	if input.UserID != nil {
		user, err := s.Deps.Repo.GetUserByID(ctx, *input.UserID)
		if err != nil {
			return nil, errors.NewNotFoundError("user")
		}
		if user.IsArchived {
			return nil, errors.NewAuthError("user is archived")
		}
		if user.IsBanned {
			return nil, errors.NewAuthError("user is banned")
		}
	}

	// Step 2: List all active variant IDs
	var variantIDs []uuid.UUID
	var err error

	if input.UserID != nil {
		variantIDs, err = s.Deps.Repo.ListActiveVariantIDsByUser(ctx, *input.UserID)
	} else {
		variantIDs, err = s.Deps.Repo.ListActiveVariantIDsByGuest(ctx, *input.GuestUserID)
	}

	if err != nil {
		return nil, errors.NewTableError("cart_items", "cannot list variant IDs")
	}

	if len(variantIDs) == 0 {
		return &GetAllCartItemsResult{
			ValidItems:   []CartGroupedBySeller{},
			InvalidItems: []InvalidCartItem{},
		}, nil
	}

	// Step 3: Fetch enriched cart snapshot rows
	// var rows []sqlc.GetActiveCartVariantSnapshotsByOwnerAndVariantIDsRow
	var rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow

	if input.UserID != nil {
		rows, err = s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx, sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
			UserID:     sqlnull.UUIDPtr(input.UserID),
			VariantIds: variantIDs,
		})
	} else {
		var guestRows []sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsRow

		guestRows, err = s.Deps.Repo.GetActiveCartVariantSnapshotsByGuestAndVariantIDs(ctx, sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsParams{
			GuestUserID: sqlnull.UUIDPtr(input.GuestUserID),
			VariantIds:  variantIDs,
		})

		// 🪄 Convert guestRows to unified []userRow type so rest of code works
		for _, r := range guestRows {
			rows = append(rows, sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow(r))
		}

	}

	if err != nil {
		return nil, errors.NewServerError("cannot fetch snapshot enriched cart items")
	}

	// Step 4: Group and filter rows
	grouped := make(map[uuid.UUID]*CartGroupedBySeller)
	productMap := make(map[uuid.UUID]map[uuid.UUID]*CartProductItem)
	var invalidItems []InvalidCartItem

	for _, row := range rows {
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

		// Optional fields
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

		productMap[row.Sellerid][row.Productid].Variants = append(productMap[row.Sellerid][row.Productid].Variants, variant)
	}

	// Final result
	var validItems []CartGroupedBySeller
	for _, seller := range grouped {
		validItems = append(validItems, *seller)
	}

	return &GetAllCartItemsResult{
		ValidItems:   validItems,
		InvalidItems: invalidItems,
	}, nil
}
