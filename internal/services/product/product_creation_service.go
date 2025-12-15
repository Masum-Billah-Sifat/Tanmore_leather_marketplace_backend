// ------------------------------------------------------------
// 📁 File: internal/service/product/create_product_service.go
// 🧠 Handles seller product creation workflow.
//     - Validates seller moderation & approval
//     - Creates product
//     - Creates one or more variants
//     - Emits product.created event
//     - Returns product_id and variant_ids

package product

import (
	"context"
	"encoding/json"
	"time"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_creation"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📦 Variant input from handler
type CreateProductVariantInput struct {
	Color                 string
	Size                  string
	RetailPrice           int64
	InStock               bool
	StockQuantity         int64
	RetailDiscount        *int64
	RetailDiscountType    *string
	WholesalePrice        *int64
	MinQtyWholesale       *int64
	WholesaleDiscount     *int64
	WholesaleDiscountType *string
	WeightGrams           int64
}

// ------------------------------------------------------------
// 📦 Input from handler
type CreateProductInput struct {
	UserID        uuid.UUID
	CategoryID    uuid.UUID
	Title         string
	Description   string
	ImageURLs     []string
	PromoVideoURL *string
	Variants      []CreateProductVariantInput
}

// ------------------------------------------------------------
// 📦 Result returned to handler
type CreateProductResult struct {
	ProductID  uuid.UUID
	VariantIDs []uuid.UUID
}

// ------------------------------------------------------------
// 📦 Dependencies
type CreateProductServiceDeps struct {
	Repo repo.ProductRepoInterface
}

// ------------------------------------------------------------
// 📦 Service definition
type CreateProductService struct {
	Deps CreateProductServiceDeps
}

// ------------------------------------------------------------
// 🛠️ Constructor
func NewCreateProductService(deps CreateProductServiceDeps) *CreateProductService {
	return &CreateProductService{Deps: deps}
}

// ------------------------------------------------------------
// 🚀 Main Logic
func (s *CreateProductService) Start(
	ctx context.Context,
	input CreateProductInput,
) (*CreateProductResult, error) {

	now := timeutil.NowUTC()

	productID := uuidutil.New()
	variantIDs := make([]uuid.UUID, 0)

	// ------------------------------------------------------------
	// Step 1: Everything inside a transaction
	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {

		// ------------------------------------------------------------
		// Step 1.1: Validate seller
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}

		if user.IsArchived || user.IsBanned {
			return errors.NewAuthError("seller account is not allowed")
		}

		if !user.IsSellerProfileApproved {
			return errors.NewAuthError("seller profile not approved")
		}

		// ------------------------------------------------------------
		// Step 1.2: Insert product
		_, err = q.InsertProduct(ctx, sqlc.InsertProductParams{
			ID:          productID,
			SellerID:    input.UserID,
			CategoryID:  input.CategoryID,
			Title:       input.Title,
			Description: input.Description,
			IsApproved:  false,
			IsBanned:    false,
			IsArchived:  false,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		if err != nil {
			return errors.NewTableError("products.insert", err.Error())
		}

		// ------------------------------------------------------------
		// ------------------------------------------------------------
		// Step 1.3: Insert variants
		for _, v := range input.Variants {

			variantID := uuidutil.New()

			hasRetailDiscount := v.RetailDiscount != nil
			hasWholesaleDiscount := v.WholesaleDiscount != nil
			wholesaleEnabled := v.WholesalePrice != nil

			_, err := q.InsertProductVariantReturningID(
				ctx,
				sqlc.InsertProductVariantReturningIDParams{
					ID:                    variantID,
					ProductID:             productID,
					Color:                 v.Color,
					Size:                  v.Size,
					RetailPrice:           v.RetailPrice,
					Retaildiscounttype:    sqlnull.StringPtr(v.RetailDiscountType),
					Retaildiscount:        sqlnull.Int64Ptr(v.RetailDiscount),
					WholesaleEnabled:      wholesaleEnabled,
					WholesalePrice:        sqlnull.Int64Ptr(v.WholesalePrice),
					MinQtyWholesale:       sqlnull.Int32Ptr(v.MinQtyWholesale),
					Wholesalediscounttype: sqlnull.StringPtr(v.WholesaleDiscountType),
					Wholesalediscount:     sqlnull.Int64Ptr(v.WholesaleDiscount),
					StockQuantity:         int32(v.StockQuantity),
					InStock:               v.InStock,
					WeightGrams:           int32(v.WeightGrams),
					IsArchived:            false,
					CreatedAt:             now,
					UpdatedAt:             now,
					Hasretaildiscount:     hasRetailDiscount,
					Haswholesalediscount:  hasWholesaleDiscount,
				},
			)
			if err != nil {
				return errors.NewTableError("product_variants.insert", err.Error())
			}

			variantIDs = append(variantIDs, variantID)
		}

		// ------------------------------------------------------------
		// Step 1.4: Insert event
		rawPayload, err := json.Marshal(map[string]interface{}{
			"user_id":    input.UserID,
			"product_id": productID,
			"variants":   input.Variants,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "product.created",
			EventPayload: rawPayload,
			DispatchedAt: sqlnull.Time(time.Time{}),
			CreatedAt:    now,
		})
		if err != nil {
			return errors.NewTableError("events.insert", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// Step 2: Prepare response
	result := &CreateProductResult{
		ProductID:  productID,
		VariantIDs: variantIDs,
	}

	return result, nil
}
