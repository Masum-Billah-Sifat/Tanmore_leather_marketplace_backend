// ------------------------------------------------------------
// 📁 File: internal/services/checkout/review_checkout_service.go
// 🧠 Service for reviewing a checkout session with item enrichment and validation.
//     - Checks customer moderation
//     - Verifies session ownership
//     - Fetches shipping address, checkout items, and variant snapshots
//     - Splits items into valid and invalid for frontend display

package checkout

import (
	"context"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/checkout/review"
	"tanmore_backend/pkg/errors"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📅 Input

type ReviewCheckoutInput struct {
	UserID            uuid.UUID
	CheckoutSessionID uuid.UUID
}

// ------------------------------------------------------------
// 📄 Output

type ReviewCheckoutResult struct {
	CheckoutSession sqlc.CheckoutSession                  `json:"checkout_session"`
	ShippingAddress *sqlc.ShippingAddress                 `json:"shipping_address,omitempty"`
	ValidItems      []sqlc.GetCheckoutItemsBySessionIDRow `json:"valid_items"`
	InvalidItems    []InvalidCheckoutItem                 `json:"invalid_items"`
}

// ------------------------------------------------------------
// 🧱 Dependencies

type ReviewCheckoutService struct {
	Repo repo.ReviewCheckoutRepoInterface
}

func NewReviewCheckoutService(repo repo.ReviewCheckoutRepoInterface) *ReviewCheckoutService {
	return &ReviewCheckoutService{Repo: repo}
}

// 🚀 Entry Point
func (s *ReviewCheckoutService) Start(
	ctx context.Context,
	input ReviewCheckoutInput,
) (*ReviewCheckoutResult, error) {
	// Step 1: Check user moderation
	user, err := s.Repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, errors.ErrAuthUserNotFound()
	}
	if user.IsArchived {
		return nil, errors.ErrAuthArchivedUser()
	}
	if user.IsBanned {
		return nil, errors.ErrAuthBannedUser()
	}

	// Step 2: Get checkout session
	session, err := s.Repo.GetCheckoutSessionByID(ctx, input.CheckoutSessionID)
	if err != nil {
		return nil, errors.NewNotFoundError("checkout_session")
	}
	if session.UserID != input.UserID {
		return nil, errors.NewAuthError("checkout session does not belong to user")
	}
	// after session fetch + ownership check
	if session.Status != "ready_to_order" {
		return nil, errors.NewValidationError("checkout_session", "invalid session status")
	}

	if !session.ShippingAddressID.Valid {
		return nil, errors.NewServerError("shipping address missing")
	}

	// Step 3: Get shipping address if set
	var shippingAddress *sqlc.ShippingAddress
	if session.ShippingAddressID.Valid {
		addr, err := s.Repo.GetShippingAddressByIDAndCheckoutID(ctx, sqlc.GetShippingAddressByIDAndCheckoutIDParams{
			ID:                session.ShippingAddressID.UUID,
			CheckoutSessionID: session.ID,
		})
		if err != nil {
			return nil, errors.NewNotFoundError("shipping_address")
		}
		shippingAddress = &addr
	}

	// Step 4: Get checkout items
	items, err := s.Repo.GetCheckoutItemsBySessionID(ctx, session.ID)
	if err != nil {
		return nil, errors.NewServerError("could not fetch checkout items")
	}
	if len(items) == 0 {
		return nil, errors.NewServerError("no checkout items found")
	}

	// Step 5: Enrich and validate
	variantIDs := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		variantIDs = append(variantIDs, item.VariantID)
	}

	snapshots, err := s.Repo.GetProductVariantSnapshotsByVariantIDs(ctx, variantIDs)
	if err != nil {
		return nil, errors.NewServerError("could not fetch variant snapshots")
	}

	snapshotMap := make(map[uuid.UUID]sqlc.GetProductVariantSnapshotsByVariantIDsRow)
	for _, snap := range snapshots {
		snapshotMap[snap.Variantid] = snap
	}

	var validItems []sqlc.GetCheckoutItemsBySessionIDRow
	var invalidItems []InvalidCheckoutItem

	for _, item := range items {
		snap, ok := snapshotMap[item.VariantID]
		if !ok {
			invalidItems = append(invalidItems, toInvalid(item.VariantID, "variant_snapshot_not_found"))
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
			validItems = append(validItems, item)
		}
	}

	return &ReviewCheckoutResult{
		CheckoutSession: session,
		ShippingAddress: shippingAddress,
		ValidItems:      validItems,
		InvalidItems:    invalidItems,
	}, nil
}

func toInvalid(variantID uuid.UUID, reason string) InvalidCheckoutItem {
	return InvalidCheckoutItem{
		VariantID: variantID,
		Reason:    reason,
	}
}
