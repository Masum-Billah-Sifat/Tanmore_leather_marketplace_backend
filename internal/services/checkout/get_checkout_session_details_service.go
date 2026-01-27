// ------------------------------------------------------------
// 📁 File: internal/services/checkout/get_checkout_session_details_service.go
// 🧠 Retrieves full checkout session details with enrichment + validation.
//     - Validates customer (not banned/archived)
//     - Verifies checkout session ownership + status
//     - Loads shipping address if available
//     - Fetches checkout items
//     - Enriches items via variant snapshots
//     - Filters into valid and invalid items (with failure reasons)

package checkout

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/checkout/session_details"
	"tanmore_backend/pkg/errors"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input

type GetCheckoutSessionDetailsInput struct {
	UserID            uuid.UUID
	CheckoutSessionID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Output

type CheckoutSessionDetailsResult struct {
	CheckoutSession sqlc.CheckoutSession                  `json:"checkout_session"`
	ShippingAddress *sqlc.ShippingAddress                 `json:"shipping_address,omitempty"`
	ValidItems      []sqlc.GetCheckoutItemsBySessionIDRow `json:"valid_items"`
	InvalidItems    []InvalidCheckoutItem                 `json:"invalid_items"`
}

type InvalidCheckoutItem struct {
	VariantID uuid.UUID `json:"variant_id"`
	Reason    string    `json:"failure_reason"`
}

// ------------------------------------------------------------
// 🧱 Dependencies

type GetCheckoutSessionDetailsServiceDeps struct {
	Repo repo.CheckoutSessionDetailsRepoInterface
}

// 🛠️ Service Definition
type GetCheckoutSessionDetailsService struct {
	Deps GetCheckoutSessionDetailsServiceDeps
}

func NewGetCheckoutSessionDetailsService(deps GetCheckoutSessionDetailsServiceDeps) *GetCheckoutSessionDetailsService {
	return &GetCheckoutSessionDetailsService{Deps: deps}
}

// 🚀 Entry Point
func (s *GetCheckoutSessionDetailsService) Start(
	ctx context.Context,
	input GetCheckoutSessionDetailsInput,
) (*CheckoutSessionDetailsResult, error) {
	// Step 1: Validate customer
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

	// Step 2: Get checkout session
	session, err := s.Deps.Repo.GetCheckoutSessionByID(ctx, input.CheckoutSessionID)
	if err != nil {
		return nil, errors.NewNotFoundError("checkout_session")
	}
	if session.UserID != input.UserID {
		return nil, errors.NewServerError("checkout session does not belong to user")
	}
	if session.Status != "ready_to_order" && session.Status != "awaiting_shipping_info" {
		return nil, errors.NewServerError("invalid checkout session status")
	}

	// Step 3: Load shipping address (if present)
	var shippingAddress *sqlc.ShippingAddress
	if session.ShippingAddressID.Valid {
		addr, err := s.Deps.Repo.GetShippingAddressByIDAndCheckoutID(ctx, sqlc.GetShippingAddressByIDAndCheckoutIDParams{
			ID:                session.ShippingAddressID.UUID,
			CheckoutSessionID: session.ID,
		})
		if err != nil {
			return nil, errors.NewNotFoundError("shipping_address")
		}
		shippingAddress = &addr
	}

	// Step 4: Get checkout items
	items, err := s.Deps.Repo.GetCheckoutItemsBySessionID(ctx, input.CheckoutSessionID)
	if err != nil {
		return nil, errors.NewTableError("checkout_items", "could not fetch")
	}
	if len(items) == 0 {
		return nil, errors.NewServerError("no checkout items found")
	}

	// Step 5: Enrich + Validate
	variantIDs := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		variantIDs = append(variantIDs, item.VariantID)
	}

	snapshots, err := s.Deps.Repo.GetProductVariantSnapshotsByVariantIDs(ctx, variantIDs)
	if err != nil {
		return nil, errors.NewServerError("could not enrich checkout items")
	}

	// Build snapshot lookup map
	snapshotMap := make(map[uuid.UUID]sqlc.GetProductVariantSnapshotsByVariantIDsRow)
	for _, snap := range snapshots {
		snapshotMap[snap.Variantid] = snap
	}

	var validItems []sqlc.GetCheckoutItemsBySessionIDRow
	var invalidItems []InvalidCheckoutItem

	for _, item := range items {
		snap, ok := snapshotMap[item.VariantID]
		if !ok {
			invalidItems = append(invalidItems, InvalidCheckoutItem{
				VariantID: item.VariantID,
				Reason:    "variant_snapshot_not_found",
			})
			continue
		}

		switch {
		case snap.Issellerbanned:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "seller_banned"))

		case snap.Issellerarchived:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "seller_archived"))

		case !snap.Issellerapproved:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "seller_not_approved"))

		case snap.Isproductbanned:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "product_banned"))

		case snap.Isproductarchived:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "product_archived"))

		case !snap.Isproductapproved:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "product_not_approved"))

		case snap.Iscategoryarchived:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "category_archived"))

		case snap.Isvariantarchived:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "variant_archived"))

		case !snap.Isvariantinstock || snap.Stockamount < item.RequiredQuantity:
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "variant_out_of_stock"))

		default:
			// ✅ item is already the correct SQLC row type
			validItems = append(validItems, item)
		}
	}

	return &CheckoutSessionDetailsResult{
		CheckoutSession: session,
		ShippingAddress: shippingAddress,
		ValidItems:      validItems,
		InvalidItems:    invalidItems,
	}, nil
}

// func toInvalid(variantID uuid.UUID, reason string) InvalidCheckoutItem {
// 	return InvalidCheckoutItem{
// 		VariantID: variantID,
// 		Reason:    reason,
// 	}
// }
