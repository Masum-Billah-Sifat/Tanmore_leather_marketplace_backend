// ------------------------------------------------------------
// 📁 File: internal/services/cart/cart_summary_service.go
// 🧠 Handles POST /api/cart/summary with full wholesale and discount logic.

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/cart_summary"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type CartSummaryInput struct {
	UserID      *uuid.UUID
	GuestUserID *uuid.UUID
	VariantIDs  []uuid.UUID
}

// 📤 Output returned to handler
type CartSummaryResult struct {
	TotalPrice   int64                    `json:"total_price"`
	InvalidItems []InvalidCartSummaryItem `json:"invalid_items"`
}

type InvalidCartSummaryItem struct {
	VariantID    uuid.UUID `json:"variant_id"`
	Reason       string    `json:"reason"`
	ProductID    uuid.UUID `json:"product_id,omitempty"`
	ProductTitle string    `json:"product_title,omitempty"`
	Color        string    `json:"color,omitempty"`
	Size         string    `json:"size,omitempty"`
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

// 🚀 Entrypoint
func (s *CartSummaryService) Start(
	ctx context.Context,
	input CartSummaryInput,
) (*CartSummaryResult, error) {

	// Step 1: Validate user (if authenticated)
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

	// // Step 2: Fetch variant snapshots based on owner type
	// var rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow

	// if input.UserID != nil {
	// 	rows, err = s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx, sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
	// 		UserID:     sqlnull.UUIDPtr(input.UserID),
	// 		VariantIds: input.VariantIDs,
	// 	})
	// } else {
	// 	var guestRows []sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsRow

	// 	guestRows, err = s.Deps.Repo.GetActiveCartVariantSnapshotsByGuestAndVariantIDs(ctx, sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsParams{
	// 		GuestUserID: sqlnull.UUIDPtr(input.GuestUserID),
	// 		VariantIds:  input.VariantIDs,
	// 	})

	// 	// 🪄 Convert guestRows to unified []userRow type so rest of code works
	// 	for _, r := range guestRows {
	// 		rows = append(rows, sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow(r))
	// 	}

	// }

	// Step 2: Fetch variant snapshots based on owner type
	var (
		rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow
		err  error
	)

	if input.UserID != nil {
		rows, err = s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx, sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
			UserID:     sqlnull.UUIDPtr(input.UserID),
			VariantIds: input.VariantIDs,
		})
	} else {
		var guestRows []sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsRow
		guestRows, err = s.Deps.Repo.GetActiveCartVariantSnapshotsByGuestAndVariantIDs(ctx, sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsParams{
			GuestUserID: sqlnull.UUIDPtr(input.GuestUserID),
			VariantIds:  input.VariantIDs,
		})

		// 🪄 Cast guestRows → userRow slice
		for _, r := range guestRows {
			rows = append(rows, sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow(r))
		}
	}

	if err != nil {
		return nil, errors.NewServerError("failed to fetch cart variant snapshots")
	}

	var totalPrice int64
	var invalidItems []InvalidCartSummaryItem

	// ✅ Track which variants were found
	foundVariantIDs := make(map[uuid.UUID]bool)
	for _, row := range rows {
		foundVariantIDs[row.Variantid] = true
	}

	// Step 3: Process each variant
	for _, row := range rows {
		quantity := int32(0)
		if row.CartRequiredQuantity.Valid {
			quantity = row.CartRequiredQuantity.Int32
		}
		if quantity == 0 {
			continue
		}

		// Moderation/availability checks
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

		// --------------------------------------------
		// 🧮 Pricing logic
		unitPrice := row.Retailprice

		if row.Haswholesaleenabled &&
			row.Wholesaleminquantity.Valid &&
			quantity >= row.Wholesaleminquantity.Int32 {

			// Use wholesale price if valid
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
			// Fallback to retail discount
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

	// Step 4: Include missing variant IDs as invalid
	for _, variantID := range input.VariantIDs {
		if !foundVariantIDs[variantID] {
			invalidItems = append(invalidItems, InvalidCartSummaryItem{
				VariantID: variantID,
				Reason:    "variant not found in system",
			})
		}
	}

	return &CartSummaryResult{
		TotalPrice:   totalPrice,
		InvalidItems: invalidItems,
	}, nil
}
