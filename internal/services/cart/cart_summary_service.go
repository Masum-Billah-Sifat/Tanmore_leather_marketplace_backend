// ------------------------------------------------------------
// 📁 File: internal/services/cart/cart_summary_service.go
// 🧠 Handles POST /api/cart/summary with full wholesale and discount logic.

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/cart_summary"
	"tanmore_backend/pkg/errors"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type CartSummaryInput struct {
	UserID     uuid.UUID
	VariantIDs []uuid.UUID
}

// ------------------------------------------------------------
// 📤 Output returned to handler
// type CartSummaryResult struct {
// 	TotalPrice int64 `json:"total_price"`
// }

type CartSummaryResult struct {
	TotalPrice   int64                    `json:"total_price"`
	InvalidItems []InvalidCartSummaryItem `json:"invalid_items"`
}

type InvalidCartSummaryItem struct {
	VariantID    uuid.UUID `json:"variant_id"`
	Reason       string    `json:"reason"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductTitle string    `json:"product_title"`
	Color        string    `json:"color"`
	Size         string    `json:"size"`
}

// ------------------------------------------------------------
// 🧱 Dependencies
type CartSummaryServiceDeps struct {
	Repo repo.CartSummaryRepoInterface
}

// 🛠️ Service Definition
type CartSummaryService struct {
	Deps CartSummaryServiceDeps
}

// 🚀 Constructor
func NewCartSummaryService(deps CartSummaryServiceDeps) *CartSummaryService {
	return &CartSummaryService{Deps: deps}
}

func (s *CartSummaryService) Start(
	ctx context.Context,
	input CartSummaryInput,
) (*CartSummaryResult, error) {

	// Step 1: Validate user
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

	// Step 2: Fetch snapshot-enriched cart items
	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
			UserID:     input.UserID,
			VariantIds: input.VariantIDs,
		})
	if err != nil {
		return nil, errors.NewServerError("failed to fetch cart variant snapshots")
	}

	var totalPrice int64 = 0
	var invalidItems []InvalidCartSummaryItem

	// ✅ Track which variants were actually found
	foundVariantIDs := make(map[uuid.UUID]bool)
	for _, row := range rows {
		foundVariantIDs[row.Variantid] = true
	}

	// Step 3: Per-item processing
	for _, row := range rows {

		// Extract quantity
		var quantity int32
		if row.CartRequiredQuantity.Valid {
			quantity = row.CartRequiredQuantity.Int32
		}
		if quantity == 0 {
			continue
		}

		// Moderation / availability checks
		if !row.Issellerapproved || row.Issellerarchived || row.Issellerbanned ||
			!row.Isproductapproved || row.Isproductarchived || row.Isproductbanned ||
			row.Isvariantarchived || !row.Isvariantinstock {

			invalidItems = append(invalidItems, InvalidCartSummaryItem{
				VariantID:    row.Variantid,
				Reason:       "variant unavailable due to moderation or stock",
				ProductID:    row.Productid,
				ProductTitle: row.Producttitle,
				Color:        row.Color,
				Size:         row.Size,
			})
			continue
		}

		// Pricing logic
		unitPrice := row.Retailprice

		// Wholesale path
		if row.Haswholesaleenabled &&
			row.Wholesaleminquantity.Valid &&
			quantity >= row.Wholesaleminquantity.Int32 {

			if row.Wholesaleprice.Valid {
				unitPrice = row.Wholesaleprice.Int64

				if row.Haswholesalediscount &&
					row.Wholesalediscount.Valid &&
					row.Wholesalediscounttype.Valid {

					switch row.Wholesalediscounttype.String {
					case "flat":
						unitPrice -= row.Wholesalediscount.Int64
					case "percentage":
						unitPrice -= (unitPrice * row.Wholesalediscount.Int64) / 100
					}
					if unitPrice < 0 {
						unitPrice = 0
					}
				}
			}

		} else {
			// Retail path
			if row.Hasretaildiscount &&
				row.Retaildiscount.Valid &&
				row.Retaildiscounttype.Valid {

				switch row.Retaildiscounttype.String {
				case "flat":
					unitPrice -= row.Retaildiscount.Int64
				case "percentage":
					unitPrice -= (unitPrice * row.Retaildiscount.Int64) / 100
				}
				if unitPrice < 0 {
					unitPrice = 0
				}
			}
		}

		totalPrice += unitPrice * int64(quantity)
	}

	// ✅ Add variants that were requested but not found in DB
	for _, inputID := range input.VariantIDs {
		if !foundVariantIDs[inputID] {
			invalidItems = append(invalidItems, InvalidCartSummaryItem{
				VariantID: inputID,
				Reason:    "variant not found in system",
			})
		}
	}

	return &CartSummaryResult{
		TotalPrice:   totalPrice,
		InvalidItems: invalidItems,
	}, nil
}

// func (s *CartSummaryService) Start(
// 	ctx context.Context,
// 	input CartSummaryInput,
// ) (*CartSummaryResult, error) {
// 	// Step 1: Validate user
// 	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
// 	if err != nil {
// 		return nil, errors.NewNotFoundError("user")
// 	}
// 	if user.IsArchived {
// 		return nil, errors.NewAuthError("user is archived")
// 	}
// 	if user.IsBanned {
// 		return nil, errors.NewAuthError("user is banned")
// 	}

// 	// Step 2: Fetch snapshot-enriched cart items
// 	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
// 		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
// 			UserID:     input.UserID,
// 			VariantIds: input.VariantIDs,
// 		})
// 	if err != nil {
// 		return nil, errors.NewServerError("failed to fetch cart variant snapshots")
// 	}

// 	var totalPrice int64 = 0
// 	var invalidItems []InvalidCartSummaryItem

// 	// Step 3: Per-item processing
// 	for _, row := range rows {
// 		// Extract quantity
// 		var quantity int32
// 		if row.CartRequiredQuantity.Valid {
// 			quantity = row.CartRequiredQuantity.Int32
// 		}
// 		if quantity == 0 {
// 			continue
// 		}

// 		// Moderation check
// 		if !row.Issellerapproved || row.Issellerarchived || row.Issellerbanned ||
// 			!row.Isproductapproved || row.Isproductarchived || row.Isproductbanned ||
// 			row.Isvariantarchived || !row.Isvariantinstock {
// 			invalidItems = append(invalidItems, InvalidCartSummaryItem{
// 				VariantID:    row.Variantid,
// 				Reason:       "variant unavailable due to moderation or stock",
// 				ProductID:    row.Productid,
// 				ProductTitle: row.Producttitle,
// 				Color:        row.Color,
// 				Size:         row.Size,
// 			})
// 			continue
// 		}

// 		// Pricing logic
// 		unitPrice := row.Retailprice

// 		// Wholesale path
// 		if row.Haswholesaleenabled && row.Wholesaleminquantity.Valid && quantity >= row.Wholesaleminquantity.Int32 {
// 			if row.Wholesaleprice.Valid {
// 				unitPrice = row.Wholesaleprice.Int64
// 				if row.Haswholesalediscount && row.Wholesalediscount.Valid && row.Wholesalediscounttype.Valid {
// 					switch row.Wholesalediscounttype.String {
// 					case "flat":
// 						unitPrice -= row.Wholesalediscount.Int64
// 					case "percentage":
// 						unitPrice -= (unitPrice * row.Wholesalediscount.Int64) / 100
// 					}
// 					if unitPrice < 0 {
// 						unitPrice = 0
// 					}
// 				}
// 			}
// 		} else {
// 			// Retail path
// 			if row.Hasretaildiscount && row.Retaildiscount.Valid && row.Retaildiscounttype.Valid {
// 				switch row.Retaildiscounttype.String {
// 				case "flat":
// 					unitPrice -= row.Retaildiscount.Int64
// 				case "percentage":
// 					unitPrice -= (unitPrice * row.Retaildiscount.Int64) / 100
// 				}
// 				if unitPrice < 0 {
// 					unitPrice = 0
// 				}
// 			}
// 		}

// 		itemTotal := unitPrice * int64(quantity)
// 		totalPrice += itemTotal
// 	}

// 	return &CartSummaryResult{
// 		TotalPrice:   totalPrice,
// 		InvalidItems: invalidItems,
// 	}, nil
// }

// // 🚀 Entrypoint
// func (s *CartSummaryService) Start(
// 	ctx context.Context,
// 	input CartSummaryInput,
// ) (*CartSummaryResult, error) {
// 	// Step 1: Validate user
// 	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
// 	if err != nil {
// 		return nil, errors.NewNotFoundError("user")
// 	}
// 	if user.IsArchived {
// 		return nil, errors.NewAuthError("user is archived")
// 	}
// 	if user.IsBanned {
// 		return nil, errors.NewAuthError("user is banned")
// 	}

// 	// Step 2: Fetch snapshot-enriched cart items
// 	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
// 		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
// 			UserID:     input.UserID,
// 			VariantIds: input.VariantIDs,
// 		})
// 	if err != nil {
// 		return nil, errors.NewServerError("failed to fetch cart variant snapshots")
// 	}

// 	var totalPrice int64 = 0

// 	// Step 3: Per-item total price calculation
// 	for _, row := range rows {
// 		// Skip invalid or moderated variants
// 		if !row.Issellerapproved || row.Issellerarchived || row.Issellerbanned ||
// 			!row.Isproductapproved || row.Isproductarchived || row.Isproductbanned ||
// 			row.Isvariantarchived || !row.Isvariantinstock {
// 			continue
// 		}

// 		// Extract quantity
// 		var quantity int32
// 		if row.CartRequiredQuantity.Valid {
// 			quantity = row.CartRequiredQuantity.Int32
// 		}
// 		if quantity == 0 {
// 			continue
// 		}

// 		// Determine per-unit price
// 		unitPrice := row.Retailprice

// 		if row.Haswholesaleenabled && row.Wholesaleminquantity.Valid && quantity >= row.Wholesaleminquantity.Int32 {
// 			// Wholesale pricing path
// 			if row.Wholesaleprice.Valid {
// 				unitPrice = row.Wholesaleprice.Int64

// 				// Apply wholesale discount if available
// 				if row.Haswholesalediscount && row.Wholesalediscount.Valid && row.Wholesalediscounttype.Valid {
// 					switch row.Wholesalediscounttype.String {
// 					case "flat":
// 						unitPrice -= row.Wholesalediscount.Int64
// 						if unitPrice < 0 {
// 							unitPrice = 0
// 						}
// 					case "percentage":
// 						discount := (unitPrice * row.Wholesalediscount.Int64) / 100
// 						unitPrice -= discount
// 						if unitPrice < 0 {
// 							unitPrice = 0
// 						}
// 					}
// 				}
// 			}
// 		} else {
// 			// Retail pricing path
// 			if row.Hasretaildiscount && row.Retaildiscount.Valid && row.Retaildiscounttype.Valid {
// 				switch row.Retaildiscounttype.String {
// 				case "flat":
// 					unitPrice -= row.Retaildiscount.Int64
// 					if unitPrice < 0 {
// 						unitPrice = 0
// 					}
// 				case "percentage":
// 					discount := (unitPrice * row.Retaildiscount.Int64) / 100
// 					unitPrice -= discount
// 					if unitPrice < 0 {
// 						unitPrice = 0
// 					}
// 				}
// 			}
// 		}

// 		itemTotal := unitPrice * int64(quantity)
// 		totalPrice += itemTotal
// 	}

// 	return &CartSummaryResult{
// 		TotalPrice: totalPrice,
// 	}, nil
// }
