// ------------------------------------------------------------
// 📁 File: internal/services/product/get_all_grouped_products_service.go
// 🧠 Handles fetching all grouped products + variants for a seller
//     - Validates seller moderation
//     - Fetches all variant indexes by seller
//     - Fetches primary product images
//     - Groups by product_id -> valid, archived, banned
//     - Formats final grouped response

package product

import (
	"context"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_get_all_grouped"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler

type GetAllProductsBySellerInput struct {
	UserID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Final Response

type GroupedProduct struct {
	ProductID        uuid.UUID                `json:"product_id"`
	Title            string                   `json:"title"`
	Description      string                   `json:"description"`
	CategoryID       uuid.UUID                `json:"category_id"`
	CategoryName     string                   `json:"category_name"`
	ImageURLs        []string                 `json:"image_urls"`
	PromoVideoURL    *string                  `json:"promo_video_url"`
	PrimaryImageURL  *string                  `json:"primary_image_url"`
	ValidVariants    []ProductVariantResponse `json:"valid_variants"`
	ArchivedVariants []ProductVariantResponse `json:"archived_variants"`
}

type GetAllProductsBySellerResult struct {
	SellerID                 uuid.UUID        `json:"seller_id"`
	SellerStoreName          string           `json:"seller_store_name"`
	ValidApprovedProducts    []GroupedProduct `json:"valid_products"`
	ValidNonApprovedProducts []GroupedProduct `json:"valid_non_approved_products"`
	ArchivedProducts         []GroupedProduct `json:"archived_products"`
	BannedProducts           []GroupedProduct `json:"banned_products"`
}

// ------------------------------------------------------------
// 🧱 Dependencies

type GetAllProductsBySellerServiceDeps struct {
	Repo repo.ProductGetAllGroupedRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service

type GetAllProductsBySellerService struct {
	Deps GetAllProductsBySellerServiceDeps
}

func NewGetAllProductsBySellerService(deps GetAllProductsBySellerServiceDeps) *GetAllProductsBySellerService {
	return &GetAllProductsBySellerService{Deps: deps}
}

// ------------------------------------------------------------
// 🚀 Entrypoint

func (s *GetAllProductsBySellerService) Start(
	ctx context.Context,
	input GetAllProductsBySellerInput,
) (*GetAllProductsBySellerResult, error) {

	// Step 1: Validate seller
	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, errors.NewNotFoundError("seller")
	}
	if user.IsArchived || user.IsBanned || !user.IsSellerProfileApproved || !user.IsSellerProfileCreated {
		return nil, errors.NewValidationError("seller", "not allowed")
	}

	// Step 2: Fetch all product_variant_index rows
	variantRows, err := s.Deps.Repo.GetAllProductVariantIndexesBySeller(ctx, input.UserID)
	if err != nil {
		return nil, errors.NewTableError("product_variant_indexes.select", err.Error())
	}

	// Step 3: Build productID to primaryImageURL map
	primaryImages := make(map[uuid.UUID]*string)
	seen := make(map[uuid.UUID]bool)
	for _, row := range variantRows {
		if !seen[row.Productid] {
			seen[row.Productid] = true
			media, err := s.Deps.Repo.GetPrimaryImageForProduct(ctx, sqlc.GetPrimaryProductImageByProductIDParams{
				ProductID:  row.Productid,
				MediaType:  "image",
				IsPrimary:  true,
				IsArchived: false,
			})
			if err == nil {
				primaryImages[row.Productid] = &media.MediaUrl
			}
		}
	}

	// Step 4: Group variants by product_id
	type productGroup struct {
		Meta   GroupedProduct
		Status string // valid, archived, banned, non-approved
	}

	grouped := make(map[uuid.UUID]*productGroup)

	for _, v := range variantRows {
		pg, ok := grouped[v.Productid]
		if !ok {
			pg = &productGroup{
				Meta: GroupedProduct{
					ProductID:        v.Productid,
					Title:            v.Producttitle,
					Description:      v.Productdescription,
					CategoryID:       v.Categoryid,
					CategoryName:     v.Categoryname,
					ImageURLs:        v.Productimages,
					PromoVideoURL:    sqlnull.ToStringPtr(v.Productpromovideourl),
					PrimaryImageURL:  primaryImages[v.Productid],
					ValidVariants:    []ProductVariantResponse{},
					ArchivedVariants: []ProductVariantResponse{},
				},
				Status: func() string {
					if v.Isproductbanned {
						return "banned"
					} else if v.Isproductarchived {
						return "archived"
					} else if !v.Isproductapproved {
						return "nonapproved"
					}
					return "valid"
				}(),
			}
			grouped[v.Productid] = pg
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
			pg.Meta.ArchivedVariants = append(pg.Meta.ArchivedVariants, variant)
		} else {
			pg.Meta.ValidVariants = append(pg.Meta.ValidVariants, variant)
		}
	}

	// Step 5: Group products by status into final response
	var result GetAllProductsBySellerResult
	result.SellerID = user.ID

	for _, pg := range grouped {
		result.SellerStoreName = variantRows[0].Sellerstorename // consistent for all
		switch pg.Status {
		case "valid":
			result.ValidApprovedProducts = append(result.ValidApprovedProducts, pg.Meta)
		case "nonapproved":
			result.ValidNonApprovedProducts = append(result.ValidNonApprovedProducts, pg.Meta)
		case "archived":
			result.ArchivedProducts = append(result.ArchivedProducts, pg.Meta)
		case "banned":
			result.BannedProducts = append(result.BannedProducts, pg.Meta)
		}
	}

	return &result, nil
}
