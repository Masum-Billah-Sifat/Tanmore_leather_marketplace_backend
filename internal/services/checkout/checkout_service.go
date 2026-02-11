// ------------------------------------------------------------
// 📁 File: internal/services/checkout/checkout_service.go
// 🧠 FINAL – Checkout Initiation Service
//     - Validates user
//     - Validates variants via snapshot rows
//     - Calculates subtotal + total weight
//     - Selects ONE applicable platform promotion (lowest priority)
//     - Inserts checkout_session + checkout_items
//     - Returns ONLY checkout_session_id + message

package checkout

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/checkout"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"
)

// ------------------------------------------------------------
// Inputs

type CheckoutFromProductInput struct {
	UserID    uuid.UUID
	VariantID uuid.UUID
	Quantity  int32
}

type CheckoutFromCartInput struct {
	UserID     uuid.UUID
	VariantIDs []uuid.UUID
}

// ------------------------------------------------------------
// Output (minimal by design)

type CheckoutMinimalResult struct {
	CheckoutSessionID uuid.UUID `json:"checkout_session_id"`
	Message           string    `json:"message"`
}

// ------------------------------------------------------------
// Service wiring

type CheckoutServiceDeps struct {
	Repo repo.CheckoutRepoInterface
}

type CheckoutService struct {
	Deps CheckoutServiceDeps
}

func NewCheckoutService(deps CheckoutServiceDeps) *CheckoutService {
	return &CheckoutService{Deps: deps}
}

// ============================================================
// FROM PRODUCT
// ============================================================

func (s *CheckoutService) FromProduct(
	ctx context.Context,
	input CheckoutFromProductInput,
) (*CheckoutMinimalResult, error) {
	variantIDs := []uuid.UUID{input.VariantID}

	// 🔐 User moderation
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

	// ✅ Buy-now must not depend on cart rows
	snapshots, err := s.Deps.Repo.GetVariantSnapshotsByVariantIDs(ctx, variantIDs)
	if err != nil {
		return nil, errors.NewServerError("failed to fetch variant snapshot data")
	}

	// quantity map for buy-now
	qtyByVariant := map[uuid.UUID]int32{
		input.VariantID: input.Quantity,
	}

	subtotal, totalWeight, items, invalid := processSnapshotsForCheckout(
		input.UserID,
		variantIDs,
		snapshots,
		qtyByVariant,
		true, // enforce stock cap strictly for buy-now
	)

	if len(items) == 0 {
		reason := "variant unavailable"
		if len(invalid) > 0 {
			reason = invalid[0]["reason"]
		}
		return nil, errors.NewValidationError("variant_id", reason)
	}

	promo, promoAmount := s.selectPlatformPromotion(ctx, subtotal)
	totalPayable := maxInt64(subtotal-promoAmount, 0)

	checkoutSessionID := uuidutil.New()
	now := timeutil.NowUTC()

	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		_, err := q.InsertCheckoutSession(ctx, sqlc.InsertCheckoutSessionParams{
			ID:                            checkoutSessionID,
			UserID:                        input.UserID,
			Subtotal:                      toDecimal(subtotal),
			TotalWeightGrams:              int32(totalWeight),
			DeliveryCharge:                sql.NullString{Valid: false},
			TotalPayable:                  toDecimal(totalPayable),
			ShippingAddressID:             uuid.NullUUID{},
			Status:                        "awaiting_shipping_info",
			PlatformDiscountType:          sqlnull.String(promo.DiscountType),
			PlatformDiscountValue:         sqlnull.String(promo.DiscountValue),
			PlatformDiscountAmountApplied: sqlnull.String(fmt.Sprintf("%d", promoAmount)),
			IsPlatformDiscountApplied:     promo.ID != uuid.Nil,
			PaymentMethod:                 "cod",
			CreatedAt:                     now,
		})
		if err != nil {
			return err
		}

		for i := range items {
			items[i].CheckoutSessionID = checkoutSessionID
			if err := q.InsertCheckoutItem(ctx, items[i]); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, errors.NewServerError("could not create checkout session")
	}

	return &CheckoutMinimalResult{
		CheckoutSessionID: checkoutSessionID,
		Message:           "checkout session created successfully",
	}, nil
}

// ============================================================
// FROM CART (identical flow)
// ============================================================

func (s *CheckoutService) FromCart(
	ctx context.Context,
	input CheckoutFromCartInput,
) (*CheckoutMinimalResult, error) {
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

	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx,
		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
			UserID:     input.UserID,
			VariantIds: input.VariantIDs,
		},
	)
	if err != nil {
		return nil, errors.NewServerError("failed to fetch snapshot data")
	}

	subtotal, totalWeight, items, invalid := processVariants(input.UserID, input.VariantIDs, rows)
	if len(items) == 0 {
		reason := "no valid variants"
		if len(invalid) > 0 {
			reason = invalid[0]["reason"]
		}
		return nil, errors.NewValidationError("variant_ids", reason)
	}

	promo, promoAmount := s.selectPlatformPromotion(ctx, subtotal)
	totalPayable := maxInt64(subtotal-promoAmount, 0)

	checkoutSessionID := uuidutil.New()
	now := timeutil.NowUTC()

	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		_, err := q.InsertCheckoutSession(ctx, sqlc.InsertCheckoutSessionParams{
			ID:                            checkoutSessionID,
			UserID:                        input.UserID,
			Subtotal:                      toDecimal(subtotal),
			TotalWeightGrams:              int32(totalWeight),
			DeliveryCharge:                sql.NullString{Valid: false},
			TotalPayable:                  toDecimal(totalPayable),
			ShippingAddressID:             uuid.NullUUID{},
			Status:                        "awaiting_shipping_info",
			PlatformDiscountType:          sqlnull.String(promo.DiscountType),
			PlatformDiscountValue:         sqlnull.String(promo.DiscountValue),
			PlatformDiscountAmountApplied: sqlnull.String(fmt.Sprintf("%d", promoAmount)),
			IsPlatformDiscountApplied:     promo.ID != uuid.Nil,
			PaymentMethod:                 "cod",
			CreatedAt:                     now,
		})
		if err != nil {
			return err
		}

		for i := range items {
			items[i].CheckoutSessionID = checkoutSessionID
			if err := q.InsertCheckoutItem(ctx, items[i]); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, errors.NewServerError("could not create checkout session")
	}

	return &CheckoutMinimalResult{
		CheckoutSessionID: checkoutSessionID,
		Message:           "checkout session created successfully",
	}, nil
}

// ============================================================
// VARIANT VALIDATION + ITEM BUILD
// ============================================================

func processSnapshotsForCheckout(
	userID uuid.UUID,
	requestedVariantIDs []uuid.UUID,
	snapshots []sqlc.ProductVariantSnapshot,
	qtyByVariant map[uuid.UUID]int32,
	enforceStock bool,
) (subtotal int64, totalWeight int64, items []sqlc.InsertCheckoutItemParams, invalid []map[string]string) {

	found := map[uuid.UUID]bool{}
	now := timeutil.NowUTC()

	for _, s := range snapshots {
		vid := s.Variantid
		found[vid] = true

		qty := qtyByVariant[vid]
		if qty <= 0 {
			invalid = append(invalid, map[string]string{
				"variant_id": vid.String(),
				"reason":     "invalid quantity",
			})
			continue
		}

		// Moderation / availability checks
		if s.Issellerbanned || s.Issellerarchived || !s.Issellerapproved ||
			s.Isproductbanned || s.Isproductarchived || !s.Isproductapproved ||
			s.Iscategoryarchived ||
			s.Isvariantarchived || !s.Isvariantinstock {

			invalid = append(invalid, map[string]string{
				"variant_id": vid.String(),
				"reason":     "variant unavailable or moderated",
			})
			continue
		}

		// Stock cap (if you want strict behavior)
		if enforceStock && int64(qty) > int64(s.Stockamount) {
			invalid = append(invalid, map[string]string{
				"variant_id": vid.String(),
				"reason":     "insufficient stock",
			})
			continue
		}

		// Unit price (retail)
		unitPrice := s.Retailprice

		// Discount metadata (store in checkout item; not altering price unless your rules are finalized)
		hasDiscount := false
		discountType := ""
		discountValue := ""

		// IMPORTANT: Hasretaildiscount is a boolean, but type/value fields may still be NULL.
		// Only set discount fields if they are valid.
		if s.Hasretaildiscount && s.Retaildiscounttype.Valid && s.Retaildiscount.Valid {
			hasDiscount = true
			discountType = s.Retaildiscounttype.String
			discountValue = fmt.Sprintf("%d", s.Retaildiscount.Int64) // keep consistent, cents or percent-int, your call
		}

		line := unitPrice * int64(qty)
		subtotal += line
		totalWeight += int64(qty) * int64(s.WeightGrams)

		items = append(items, sqlc.InsertCheckoutItemParams{
			ID:                     uuidutil.New(),
			CheckoutSessionID:      uuid.UUID{},
			UserID:                 userID,
			SellerID:               s.Sellerid,
			SellerStoreName:        s.Sellerstorename,
			CategoryID:             s.Categoryid,
			CategoryName:           s.Categoryname,
			ProductID:              s.Productid,
			ProductTitle:           s.Producttitle,
			ProductDescription:     s.Productdescription,
			ProductPrimaryImageUrl: s.Productprimaryimageurl,
			VariantID:              vid,
			Color:                  s.Color,
			Size:                   s.Size,
			BuyingMode:             "retail",
			UnitPrice:              toDecimal(unitPrice),

			HasDiscount:   hasDiscount,
			DiscountType:  sqlnull.String(discountType),
			DiscountValue: sqlnull.String(discountValue),

			RequiredQuantity: qty,
			WeightGrams:      int32(s.WeightGrams),
			CreatedAt:        now,
		})
	}

	// Missing snapshot rows => variant not found (or snapshot missing)
	for _, id := range requestedVariantIDs {
		if !found[id] {
			invalid = append(invalid, map[string]string{
				"variant_id": id.String(),
				"reason":     "variant not found",
			})
		}
	}

	return
}

func processVariants(
	userID uuid.UUID,
	requested []uuid.UUID,
	rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow,
) (subtotal int64, totalWeight int64, items []sqlc.InsertCheckoutItemParams, invalid []map[string]string) {

	found := map[uuid.UUID]bool{}
	now := timeutil.NowUTC()

	for _, r := range rows {
		vid := r.CartVariantID
		found[vid] = true

		qty := int32(0)
		if r.CartRequiredQuantity.Valid {
			qty = r.CartRequiredQuantity.Int32
		}

		if qty <= 0 {
			invalid = append(invalid, map[string]string{
				"variant_id": vid.String(),
				"reason":     "invalid quantity",
			})
			continue
		}

		// Availability checks (match cart summary + include category archived)
		if !r.Issellerapproved || r.Issellerarchived || r.Issellerbanned ||
			!r.Isproductapproved || r.Isproductarchived || r.Isproductbanned ||
			r.Iscategoryarchived ||
			r.Isvariantarchived || !r.Isvariantinstock {

			invalid = append(invalid, map[string]string{
				"variant_id": vid.String(),
				"reason":     "variant unavailable or moderated",
			})
			continue
		}

		// -------------------------
		// ✅ Pricing logic (same as cart summary)
		unitPrice := r.Retailprice
		buyingMode := "retail"

		hasDiscount := false
		discountType := ""
		discountValue := ""

		// wholesale eligible?
		if r.Haswholesaleenabled &&
			r.Wholesaleminquantity.Valid &&
			qty >= r.Wholesaleminquantity.Int32 &&
			r.Wholesaleprice.Valid {

			buyingMode = "wholesale"
			unitPrice = r.Wholesaleprice.Int64

			if r.Haswholesalediscount &&
				r.Wholesalediscount.Valid &&
				r.Wholesalediscounttype.Valid {

				hasDiscount = true
				discountType = r.Wholesalediscounttype.String
				discountValue = fmt.Sprintf("%d", r.Wholesalediscount.Int64)

				unitPrice = applyDiscount(unitPrice, r.Wholesalediscount.Int64, discountType)
			}

		} else {
			// retail discount
			if r.Hasretaildiscount &&
				r.Retaildiscount.Valid &&
				r.Retaildiscounttype.Valid {

				hasDiscount = true
				discountType = r.Retaildiscounttype.String
				discountValue = fmt.Sprintf("%d", r.Retaildiscount.Int64)

				unitPrice = applyDiscount(unitPrice, r.Retaildiscount.Int64, discountType)
			}
		}

		line := unitPrice * int64(qty)
		subtotal += line
		totalWeight += int64(qty) * int64(r.WeightGrams)

		items = append(items, sqlc.InsertCheckoutItemParams{
			ID:                     uuidutil.New(),
			CheckoutSessionID:      uuid.UUID{},
			UserID:                 userID,
			SellerID:               r.Sellerid,
			SellerStoreName:        r.Sellerstorename,
			CategoryID:             r.Categoryid,
			CategoryName:           r.Categoryname,
			ProductID:              r.Productid,
			ProductTitle:           r.Producttitle,
			ProductDescription:     r.Productdescription,
			ProductPrimaryImageUrl: r.Productprimaryimageurl,
			VariantID:              vid,
			Color:                  r.Color,
			Size:                   r.Size,

			BuyingMode:       buyingMode,
			UnitPrice:        toDecimal(unitPrice),
			HasDiscount:      hasDiscount,
			DiscountType:     sqlnull.String(discountType),
			DiscountValue:    sqlnull.String(discountValue),
			RequiredQuantity: qty,
			WeightGrams:      int32(r.WeightGrams),
			CreatedAt:        now,
		})
	}

	for _, id := range requested {
		if !found[id] {
			invalid = append(invalid, map[string]string{
				"variant_id": id.String(),
				"reason":     "variant not found",
			})
		}
	}

	return
}

// func processVariants(
// 	userID uuid.UUID,
// 	requested []uuid.UUID,
// 	rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow,
// ) (subtotal int64, totalWeight int64, items []sqlc.InsertCheckoutItemParams, invalid []map[string]string) {

// 	found := map[uuid.UUID]bool{}
// 	now := timeutil.NowUTC()

// 	for _, r := range rows {
// 		found[r.CartVariantID] = true

// 		qty := int32(0)
// 		if r.CartRequiredQuantity.Valid {
// 			qty = r.CartRequiredQuantity.Int32
// 		}

// 		if qty == 0 || r.Issellerbanned || r.Issellerarchived || !r.Issellerapproved ||
// 			r.Isproductbanned || r.Isproductarchived || !r.Isproductapproved ||
// 			r.Isvariantarchived || !r.Isvariantinstock {

// 			invalid = append(invalid, map[string]string{
// 				"variant_id": r.CartVariantID.String(),
// 				"reason":     "variant unavailable or moderated",
// 			})
// 			continue
// 		}

// 		line := r.Retailprice * int64(qty)
// 		subtotal += line
// 		totalWeight += int64(qty) * int64(r.WeightGrams)

// 		items = append(items, sqlc.InsertCheckoutItemParams{
// 			ID:                     uuidutil.New(),
// 			CheckoutSessionID:      uuid.UUID{},
// 			UserID:                 userID,
// 			SellerID:               r.Sellerid,
// 			SellerStoreName:        r.Sellerstorename,
// 			CategoryID:             r.Categoryid,
// 			CategoryName:           r.Categoryname,
// 			ProductID:              r.Productid,
// 			ProductTitle:           r.Producttitle,
// 			ProductDescription:     r.Productdescription,
// 			ProductPrimaryImageUrl: r.Productprimaryimageurl,
// 			VariantID:              r.CartVariantID,
// 			Color:                  r.Color,
// 			Size:                   r.Size,
// 			BuyingMode:             "retail",
// 			UnitPrice:              toDecimal(r.Retailprice),
// 			HasDiscount:            false,
// 			DiscountType:           sqlnull.String(""),
// 			DiscountValue:          sqlnull.String(""),
// 			RequiredQuantity:       qty,
// 			WeightGrams:            int32(r.WeightGrams),
// 			CreatedAt:              now,
// 		})
// 	}

// 	for _, id := range requested {
// 		if !found[id] {
// 			invalid = append(invalid, map[string]string{
// 				"variant_id": id.String(),
// 				"reason":     "variant not found",
// 			})
// 		}
// 	}

// 	return
// }

// ============================================================
// PLATFORM PROMOTION SELECTION
// ============================================================

func (s *CheckoutService) selectPlatformPromotion(
	ctx context.Context,
	subtotal int64, // subtotal is in BDT int
) (sqlc.GetActivePlatformPromotionsRow, int64) {

	promos, err := s.Deps.Repo.GetActivePlatformPromotions(ctx)
	if err != nil || len(promos) == 0 {
		return sqlc.GetActivePlatformPromotionsRow{}, 0
	}

	now := timeutil.NowUTC()
	var chosen *sqlc.GetActivePlatformPromotionsRow

	for i := range promos {
		p := &promos[i]

		if now.Before(p.StartTime) || now.After(p.EndTime) {
			continue
		}

		if p.MinCartValue.Valid && subtotal < toBDTInt(p.MinCartValue.String) {
			continue
		}
		if p.MaxCartValue.Valid && subtotal > toBDTInt(p.MaxCartValue.String) {
			continue
		}

		if chosen == nil || p.Priority < chosen.Priority {
			chosen = p
		}
	}

	if chosen == nil {
		return sqlc.GetActivePlatformPromotionsRow{}, 0
	}

	var amount int64

	switch chosen.DiscountType {
	case "flat":
		amount = toBDTInt(chosen.DiscountValue)

	case "percentage":
		// DiscountValue should be like "5" or "10"
		pct := toBDTInt(chosen.DiscountValue)
		if pct < 0 {
			pct = 0
		}
		amount = (subtotal * pct) / 100
	}

	if chosen.MaxDiscountCap.Valid {
		cap := toBDTInt(chosen.MaxDiscountCap.String)
		if amount > cap {
			amount = cap
		}
	}

	return *chosen, amount
}

// func (s *CheckoutService) selectPlatformPromotion(
// 	ctx context.Context,
// 	subtotal int64,
// ) (sqlc.GetActivePlatformPromotionsRow, int64) {

// 	promos, err := s.Deps.Repo.GetActivePlatformPromotions(ctx)
// 	if err != nil || len(promos) == 0 {
// 		return sqlc.GetActivePlatformPromotionsRow{}, 0
// 	}

// 	now := timeutil.NowUTC()
// 	var chosen *sqlc.GetActivePlatformPromotionsRow

// 	for i := range promos {
// 		p := &promos[i]

// 		// ⏰ Time window check
// 		if now.Before(p.StartTime) || now.After(p.EndTime) {
// 			continue
// 		}

// 		// 💰 Cart value constraints
// 		if p.MinCartValue.Valid && subtotal < toCents(p.MinCartValue.String) {
// 			continue
// 		}
// 		if p.MaxCartValue.Valid && subtotal > toCents(p.MaxCartValue.String) {
// 			continue
// 		}

// 		// 🥇 Lowest priority wins
// 		if chosen == nil || p.Priority < chosen.Priority {
// 			chosen = p
// 		}
// 	}

// 	if chosen == nil {
// 		return sqlc.GetActivePlatformPromotionsRow{}, 0
// 	}

// 	// 💸 Calculate discount amount
// 	var amount int64

// 	switch chosen.DiscountType {
// 	case "flat":
// 		amount = toCents(chosen.DiscountValue)

// 	case "percentage":
// 		amount = (subtotal * toCents(chosen.DiscountValue)) / 10000
// 	}

// 	// 🧢 Max discount cap
// 	if chosen.MaxDiscountCap.Valid {
// 		cap := toCents(chosen.MaxDiscountCap.String)
// 		if amount > cap {
// 			amount = cap
// 		}
// 	}

// 	return *chosen, amount
// }

// ============================================================
// HELPERS
// ============================================================

// func toDecimal(val int64) string {
// 	return fmt.Sprintf("%.2f", float64(val)/100)
// }

// func toCents(val string) int64 {
// 	var i int64
// 	fmt.Sscan(val, &i)
// 	return i * 100
// }

// func maxInt64(a, b int64) int64 {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// BDT int -> decimal string for DB (e.g., 19000 -> "19000.00")
func toDecimal(val int64) string {
	return fmt.Sprintf("%.2f", float64(val))
}

// Parse DB string like "5000" or "5000.00" into BDT int64 (5000)
func toBDTInt(val string) int64 {
	v := strings.TrimSpace(val)
	if v == "" {
		return 0
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f))
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func applyDiscount(unitPrice int64, discount int64, typ string) int64 {
	switch typ {
	case "flat":
		unitPrice -= discount
	case "percentage":
		unitPrice -= (unitPrice * discount) / 100
	}
	if unitPrice < 0 {
		unitPrice = 0
	}
	return unitPrice
}
