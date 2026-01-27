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
	if err != nil || user.IsBanned || user.IsArchived {
		return nil, errors.NewAuthError("unauthorized user")
	}

	// 📸 Snapshot rows
	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx,
		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
			UserID:     sqlnull.UUID(input.UserID),
			VariantIds: variantIDs,
		},
	)
	if err != nil {
		return nil, errors.NewServerError("failed to fetch snapshot data")
	}

	for i := range rows {
		rows[i].CartRequiredQuantity = sqlnull.Int32(int64(input.Quantity))
	}

	subtotal, totalWeight, items, invalid := processVariants(input.UserID, variantIDs, rows)
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
		// 🧾 Insert checkout session
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

		// 🧾 Insert items
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
	if err != nil || user.IsBanned || user.IsArchived {
		return nil, errors.NewAuthError("unauthorized user")
	}

	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx,
		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
			UserID:     sqlnull.UUID(input.UserID),
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

func processVariants(
	userID uuid.UUID,
	requested []uuid.UUID,
	rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow,
) (subtotal int64, totalWeight int64, items []sqlc.InsertCheckoutItemParams, invalid []map[string]string) {

	found := map[uuid.UUID]bool{}
	now := timeutil.NowUTC()

	for _, r := range rows {
		found[r.CartVariantID] = true

		qty := int32(0)
		if r.CartRequiredQuantity.Valid {
			qty = r.CartRequiredQuantity.Int32
		}

		if qty == 0 || r.Issellerbanned || r.Issellerarchived || !r.Issellerapproved ||
			r.Isproductbanned || r.Isproductarchived || !r.Isproductapproved ||
			r.Isvariantarchived || !r.Isvariantinstock {

			invalid = append(invalid, map[string]string{
				"variant_id": r.CartVariantID.String(),
				"reason":     "variant unavailable or moderated",
			})
			continue
		}

		line := r.Retailprice * int64(qty)
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
			VariantID:              r.CartVariantID,
			Color:                  r.Color,
			Size:                   r.Size,
			BuyingMode:             "retail",
			UnitPrice:              toDecimal(r.Retailprice),
			HasDiscount:            false,
			DiscountType:           sqlnull.String(""),
			DiscountValue:          sqlnull.String(""),
			RequiredQuantity:       qty,
			WeightGrams:            int32(r.WeightGrams),
			CreatedAt:              now,
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

// ============================================================
// PLATFORM PROMOTION SELECTION
// ============================================================

func (s *CheckoutService) selectPlatformPromotion(
	ctx context.Context,
	subtotal int64,
) (sqlc.GetActivePlatformPromotionsRow, int64) {

	promos, err := s.Deps.Repo.GetActivePlatformPromotions(ctx)
	if err != nil || len(promos) == 0 {
		return sqlc.GetActivePlatformPromotionsRow{}, 0
	}

	now := timeutil.NowUTC()
	var chosen *sqlc.GetActivePlatformPromotionsRow

	for i := range promos {
		p := &promos[i]

		// ⏰ Time window check
		if now.Before(p.StartTime) || now.After(p.EndTime) {
			continue
		}

		// 💰 Cart value constraints
		if p.MinCartValue.Valid && subtotal < toCents(p.MinCartValue.String) {
			continue
		}
		if p.MaxCartValue.Valid && subtotal > toCents(p.MaxCartValue.String) {
			continue
		}

		// 🥇 Lowest priority wins
		if chosen == nil || p.Priority < chosen.Priority {
			chosen = p
		}
	}

	if chosen == nil {
		return sqlc.GetActivePlatformPromotionsRow{}, 0
	}

	// 💸 Calculate discount amount
	var amount int64

	switch chosen.DiscountType {
	case "flat":
		amount = toCents(chosen.DiscountValue)

	case "percentage":
		amount = (subtotal * toCents(chosen.DiscountValue)) / 10000
	}

	// 🧢 Max discount cap
	if chosen.MaxDiscountCap.Valid {
		cap := toCents(chosen.MaxDiscountCap.String)
		if amount > cap {
			amount = cap
		}
	}

	return *chosen, amount
}

// func (s *CheckoutService) selectPlatformPromotion(ctx context.Context, subtotal int64) (sqlc.PlatformPromotion, int64) {
// 	promos, err := s.Deps.Repo.GetActivePlatformPromotions(ctx)
// 	if err != nil || len(promos) == 0 {
// 		return sqlc.PlatformPromotion{}, 0
// 	}

// 	now := timeutil.NowUTC()
// 	var chosen *sqlc.PlatformPromotion

// 	for i := range promos {
// 		p := &promos[i]
// 		if now.Before(p.StartTime) || now.After(p.EndTime) {
// 			continue
// 		}

// 		if p.MinCartValue.Valid && subtotal < toCents(p.MinCartValue.String) {
// 			continue
// 		}
// 		if p.MaxCartValue.Valid && subtotal > toCents(p.MaxCartValue.String) {
// 			continue
// 		}

// 		if chosen == nil || p.Priority < chosen.Priority {
// 			chosen = p
// 		}
// 	}

// 	if chosen == nil {
// 		return sqlc.PlatformPromotion{}, 0
// 	}

// 	amount := int64(0)
// 	switch chosen.DiscountType {
// 	case "flat":
// 		amount = toCents(chosen.DiscountValue)
// 	case "percentage":
// 		amount = (subtotal * toCents(chosen.DiscountValue)) / 10000
// 	}

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

func toDecimal(val int64) string {
	return fmt.Sprintf("%.2f", float64(val)/100)
}

func toCents(val string) int64 {
	var i int64
	fmt.Sscan(val, &i)
	return i * 100
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// // ------------------------------------------------------------
// // 📁 File: internal/services/checkout/checkout_service.go
// // 🧠 Finalized version of POST /api/checkout/initiate
// //     - Handles platform promotion selection
// //     - Inserts enriched checkout session & items
// //     - Returns only checkout_session_id and message

// package checkout

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"

// 	"github.com/google/uuid"

// 	"tanmore_backend/internal/db/sqlc"
// 	repo "tanmore_backend/internal/repository/checkout"
// 	"tanmore_backend/pkg/errors"
// 	"tanmore_backend/pkg/sqlnull"
// 	"tanmore_backend/pkg/timeutil"
// 	uuidutil "tanmore_backend/pkg/uuid"
// )

// type CheckoutFromProductInput struct {
// 	UserID    uuid.UUID
// 	VariantID uuid.UUID
// 	Quantity  int32
// }

// type CheckoutFromCartInput struct {
// 	UserID     uuid.UUID
// 	VariantIDs []uuid.UUID
// }

// type CheckoutMinimalResult struct {
// 	CheckoutSessionID uuid.UUID `json:"checkout_session_id"`
// 	Message           string    `json:"message"`
// }

// type CheckoutServiceDeps struct {
// 	Repo repo.CheckoutRepoInterface
// }

// type CheckoutService struct {
// 	Deps CheckoutServiceDeps
// }

// func NewCheckoutService(deps CheckoutServiceDeps) *CheckoutService {
// 	return &CheckoutService{Deps: deps}
// }

// func (s *CheckoutService) FromProduct(
// 	ctx context.Context,
// 	input CheckoutFromProductInput,
// ) (*CheckoutMinimalResult, error) {
// 	variantIDs := []uuid.UUID{input.VariantID}

// 	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
// 	if err != nil || user.IsBanned || user.IsArchived {
// 		return nil, errors.NewAuthError("unauthorized user")
// 	}

// 	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
// 		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
// 			UserID:     sqlnull.UUID(input.UserID),
// 			VariantIds: variantIDs,
// 		})
// 	if err != nil {
// 		return nil, errors.NewServerError("failed to fetch snapshot data")
// 	}

// 	for i := range rows {
// 		rows[i].CartRequiredQuantity = sqlnull.Int32(int64(input.Quantity))
// 	}

// 	subtotal, totalWeight, validItems, invalidItems := s.simplifiedProcessVariants(input.UserID, variantIDs, rows)

// 	if len(validItems) == 0 {
// 		reason := "variant unavailable"
// 		if len(invalidItems) > 0 {
// 			if msg, ok := invalidItems[0]["reason"]; ok {
// 				reason = msg
// 			}
// 		}
// 		return nil, errors.NewValidationError("variant_id", reason)
// 	}

// 	// 🔍 Attempt platform promo
// 	promo, promoAppliedAmount := s.Deps.Repo.FindApplicablePlatformPromo(ctx, subtotal)

// 	checkoutSessionID := uuidutil.New()
// 	now := timeutil.NowUTC()
// 	totalPayable := subtotal - promoAppliedAmount
// 	if totalPayable < 0 {
// 		totalPayable = 0
// 	}

// 	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
// 		sessionParams := sqlc.InsertCheckoutSessionParams{
// 			ID:                            checkoutSessionID,
// 			UserID:                        input.UserID,
// 			Subtotal:                      toDecimal(subtotal),
// 			TotalWeightGrams:              int32(totalWeight),
// 			DeliveryCharge:                sql.NullString{Valid: false},
// 			TotalPayable:                  toDecimal(totalPayable),
// 			ShippingAddressID:             uuid.NullUUID{},
// 			Status:                        "awaiting_shipping_info",
// 			PlatformDiscountType:          sqlnull.StringPtr(promo.DiscountType),
// 			PlatformDiscountValue:         sqlnull.StringPtr(promo.DiscountValue),
// 			PlatformDiscountAmountApplied: sqlnull.String(promoAppliedAmount),
// 			IsPlatformDiscountApplied:     promo.ID != uuid.Nil,
// 			PaymentMethod:                 "cod",
// 			CreatedAt:                     now,
// 		}

// 		if _, err := q.InsertCheckoutSession(ctx, sessionParams); err != nil {
// 			return err
// 		}
// 		for i := range validItems {
// 			validItems[i].CheckoutSessionID = checkoutSessionID
// 			if _, err := q.InsertCheckoutItem(ctx, validItems[i]); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, errors.NewServerError("could not create checkout session")
// 	}

// 	return &CheckoutMinimalResult{
// 		CheckoutSessionID: checkoutSessionID,
// 		Message:           "checkout session created successfully",
// 	}, nil
// }

// // ⚠️ FromCart() will be updated the same way (logic identical)
// // We can refactor to reuse if needed

// // 🧠 The `FindApplicablePlatformPromo` logic will:
// // - Return most prioritized active promo within time range
// // - Apply only if subtotal fits promo min/max rules
// // - Return promo struct and applied discount value
// //
// // 🔍 Example query:
// // -- name: FindActivePlatformPromotions :many
// // SELECT * FROM platform_promotions
// // WHERE is_active = TRUE AND is_archived = FALSE
// // ORDER BY priority ASC, created_at DESC;
// //
// // ✅ Applied inside Go logic:
// // - Check current time within (start_time, end_time)
// // - Check subtotal against min_cart and max_cart
// // - Calculate promoAppliedAmount = min(discount_value, max_discount_cap)
// // - Return the promo + applied value

// func (s *CheckoutService) simplifiedProcessVariants(
// 	userID uuid.UUID,
// 	requestedVariantIDs []uuid.UUID,
// 	rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow,
// ) (int64, int64, []sqlc.InsertCheckoutItemParams, []map[string]string) {
// 	var (
// 		subtotal     int64
// 		totalWeight  int64
// 		validItems   []sqlc.InsertCheckoutItemParams
// 		invalidItems []map[string]string
// 		now          = timeutil.NowUTC()
// 		foundMap     = make(map[uuid.UUID]bool)
// 	)

// 	for _, row := range rows {
// 		foundMap[row.CartVariantID] = true
// 		quantity := int32(0)
// 		if row.CartRequiredQuantity.Valid {
// 			quantity = row.CartRequiredQuantity.Int32
// 		}

// 		if !row.Issellerapproved || row.Issellerarchived || row.Issellerbanned ||
// 			!row.Isproductapproved || row.Isproductarchived || row.Isproductbanned ||
// 			row.Isvariantarchived || !row.Isvariantinstock || quantity == 0 {
// 			invalidItems = append(invalidItems, map[string]string{
// 				"variant_id": row.CartVariantID.String(),
// 				"reason":     "variant unavailable or moderated",
// 			})
// 			continue
// 		}

// 		price := row.Retailprice
// 		total := price * int64(quantity)
// 		subtotal += total
// 		totalWeight += int64(quantity) * int64(row.WeightGrams)

// 		validItems = append(validItems, sqlc.InsertCheckoutItemParams{
// 			ID:                     uuidutil.New(),
// 			CheckoutSessionID:      uuid.UUID{},
// 			UserID:                 userID,
// 			SellerID:               row.Sellerid,
// 			SellerStoreName:        row.Sellerstorename,
// 			CategoryID:             row.Categoryid,
// 			CategoryName:           row.Categoryname,
// 			ProductID:              row.Productid,
// 			ProductTitle:           row.Producttitle,
// 			ProductDescription:     row.Productdescription,
// 			ProductPrimaryImageUrl: row.Productprimaryimageurl,
// 			VariantID:              row.CartVariantID,
// 			Color:                  row.Color,
// 			Size:                   row.Size,
// 			BuyingMode:             "retail",
// 			UnitPrice:              toDecimal(price).String(),
// 			HasDiscount:            false,
// 			DiscountType:           "",
// 			DiscountValue:          "0.00",
// 			RequiredQuantity:       quantity,
// 			WeightGrams:            int32(row.WeightGrams),
// 			CreatedAt:              now,
// 		})
// 	}

// 	for _, id := range requestedVariantIDs {
// 		if !foundMap[id] {
// 			invalidItems = append(invalidItems, map[string]string{
// 				"variant_id": id.String(),
// 				"reason":     "variant not found",
// 			})
// 		}
// 	}

// 	return subtotal, totalWeight, validItems, invalidItems
// }

// func toDecimal(val int64) string {
// 	return fmt.Sprintf("%.2f", float64(val)/100)
// }

// // // ------------------------------------------------------------
// // // 📁 File: internal/services/checkout/checkout_service.go
// // // 🧠 Refactored version of POST /api/checkout/initiate
// // //     - Still validates user moderation
// // //     - Still fetches variant snapshot data
// // //     - Still inserts checkout session + items
// // //     - ✅ But now returns ONLY checkout_session_id + message

// // package checkout

// // import (
// // 	"context"
// // 	"database/sql"
// // 	"fmt"

// // 	"github.com/google/uuid"

// // 	"tanmore_backend/internal/db/sqlc"
// // 	repo "tanmore_backend/internal/repository/checkout"
// // 	"tanmore_backend/pkg/errors"
// // 	"tanmore_backend/pkg/sqlnull"
// // 	"tanmore_backend/pkg/timeutil"
// // 	uuidutil "tanmore_backend/pkg/uuid"
// // )

// // type CheckoutFromProductInput struct {
// // 	UserID    uuid.UUID
// // 	VariantID uuid.UUID
// // 	Quantity  int32
// // }

// // type CheckoutFromCartInput struct {
// // 	UserID     uuid.UUID
// // 	VariantIDs []uuid.UUID
// // }

// // type CheckoutMinimalResult struct {
// // 	CheckoutSessionID uuid.UUID `json:"checkout_session_id"`
// // 	Message           string    `json:"message"`
// // }

// // type CheckoutServiceDeps struct {
// // 	Repo repo.CheckoutRepoInterface
// // }

// // type CheckoutService struct {
// // 	Deps CheckoutServiceDeps
// // }

// // func NewCheckoutService(deps CheckoutServiceDeps) *CheckoutService {
// // 	return &CheckoutService{Deps: deps}
// // }

// // func (s *CheckoutService) FromProduct(
// // 	ctx context.Context,
// // 	input CheckoutFromProductInput,
// // ) (*CheckoutMinimalResult, error) {
// // 	variantIDs := []uuid.UUID{input.VariantID}

// // 	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
// // 	if err != nil || user.IsBanned || user.IsArchived {
// // 		return nil, errors.NewAuthError("unauthorized user")
// // 	}

// // 	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
// // 		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
// // 			UserID:     sqlnull.UUID(input.UserID),
// // 			VariantIds: variantIDs,
// // 		})
// // 	if err != nil {
// // 		return nil, errors.NewServerError("failed to fetch snapshot data")
// // 	}

// // 	for i := range rows {
// // 		rows[i].CartRequiredQuantity = sqlnull.Int32(int64(input.Quantity))
// // 	}

// // 	subtotal, _, validItems, invalidItems := s.simplifiedProcessVariants(input.UserID, variantIDs, rows)

// // 	if len(validItems) == 0 {
// // 		reason := "variant unavailable"
// // 		if len(invalidItems) > 0 {
// // 			if msg, ok := invalidItems[0]["reason"]; ok {
// // 				reason = msg
// // 			}
// // 		}
// // 		return nil, errors.NewValidationError("variant_id", reason)
// // 	}

// // 	checkoutSessionID := uuidutil.New()
// // 	now := timeutil.NowUTC()

// // 	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
// // 		sessionParams := sqlc.InsertCheckoutSessionParams{
// // 			ID:                checkoutSessionID,
// // 			UserID:            input.UserID,
// // 			Subtotal:          toDecimal(subtotal).String(),
// // 			TotalPayable:      toDecimal(subtotal).String(),
// // 			DeliveryCharge:    sql.NullString{Valid: false},
// // 			ShippingAddressID: uuid.NullUUID{},
// // 			CreatedAt:         now,
// // 		}
// // 		if _, err := q.InsertCheckoutSession(ctx, sessionParams); err != nil {
// // 			return err
// // 		}
// // 		for i := range validItems {
// // 			validItems[i].CheckoutSessionID = checkoutSessionID
// // 			if _, err := q.InsertCheckoutItem(ctx, validItems[i]); err != nil {
// // 				return err
// // 			}
// // 		}
// // 		return nil
// // 	})
// // 	if err != nil {
// // 		return nil, errors.NewServerError("could not create checkout session")
// // 	}

// // 	return &CheckoutMinimalResult{
// // 		CheckoutSessionID: checkoutSessionID,
// // 		Message:           "checkout session created successfully",
// // 	}, nil
// // }

// // func (s *CheckoutService) FromCart(
// // 	ctx context.Context,
// // 	input CheckoutFromCartInput,
// // ) (*CheckoutMinimalResult, error) {
// // 	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
// // 	if err != nil || user.IsBanned || user.IsArchived {
// // 		return nil, errors.NewAuthError("unauthorized user")
// // 	}

// // 	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
// // 		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
// // 			UserID:     sqlnull.UUID(input.UserID),
// // 			VariantIds: input.VariantIDs,
// // 		})
// // 	if err != nil {
// // 		return nil, errors.NewServerError("failed to fetch snapshot data")
// // 	}

// // 	subtotal, _, validItems, invalidItems := s.simplifiedProcessVariants(input.UserID, input.VariantIDs, rows)

// // 	if len(validItems) == 0 {
// // 		return nil, errors.NewValidationError("variant_ids", "no valid variants to checkout")
// // 	}

// // 	checkoutSessionID := uuidutil.New()
// // 	now := timeutil.NowUTC()

// // 	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
// // 		sessionParams := sqlc.InsertCheckoutSessionParams{
// // 			ID:                checkoutSessionID,
// // 			UserID:            input.UserID,
// // 			Subtotal:          toDecimal(subtotal).String(),
// // 			TotalPayable:      toDecimal(subtotal).String(),
// // 			DeliveryCharge:    sql.NullString{Valid: false},
// // 			ShippingAddressID: uuid.NullUUID{},
// // 			CreatedAt:         now,
// // 		}
// // 		if _, err := q.InsertCheckoutSession(ctx, sessionParams); err != nil {
// // 			return err
// // 		}
// // 		for i := range validItems {
// // 			validItems[i].CheckoutSessionID = checkoutSessionID
// // 			if _, err := q.InsertCheckoutItem(ctx, validItems[i]); err != nil {
// // 				return err
// // 			}
// // 		}
// // 		return nil
// // 	})
// // 	if err != nil {
// // 		return nil, errors.NewServerError("could not create checkout session")
// // 	}

// // 	return &CheckoutMinimalResult{
// // 		CheckoutSessionID: checkoutSessionID,
// // 		Message:           "checkout session created successfully",
// // 	}, nil
// // }

// // func (s *CheckoutService) simplifiedProcessVariants(
// // 	userID uuid.UUID,
// // 	requestedVariantIDs []uuid.UUID,
// // 	rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow,
// // ) (int64, int64, []sqlc.InsertCheckoutItemParams, []map[string]string) {
// // 	var (
// // 		subtotal     int64
// // 		totalWeight  int64
// // 		validItems   []sqlc.InsertCheckoutItemParams
// // 		invalidItems []map[string]string
// // 		now          = timeutil.NowUTC()
// // 		foundMap     = make(map[uuid.UUID]bool)
// // 	)

// // 	for _, row := range rows {
// // 		foundMap[row.CartVariantID] = true
// // 		quantity := int32(0)
// // 		if row.CartRequiredQuantity.Valid {
// // 			quantity = row.CartRequiredQuantity.Int32
// // 		}

// // 		if !row.Issellerapproved || row.Issellerarchived || row.Issellerbanned ||
// // 			!row.Isproductapproved || row.Isproductarchived || row.Isproductbanned ||
// // 			row.Isvariantarchived || !row.Isvariantinstock || quantity == 0 {
// // 			invalidItems = append(invalidItems, map[string]string{
// // 				"variant_id": row.CartVariantID.String(),
// // 				"reason":     "variant unavailable or moderated",
// // 			})
// // 			continue
// // 		}

// // 		price := row.Retailprice
// // 		total := price * int64(quantity)
// // 		subtotal += total
// // 		totalWeight += int64(quantity) * int64(row.WeightGrams)

// // 		validItems = append(validItems, sqlc.InsertCheckoutItemParams{
// // 			ID:                     uuidutil.New(),
// // 			CheckoutSessionID:      uuid.UUID{},
// // 			UserID:                 userID,
// // 			SellerID:               row.Sellerid,
// // 			SellerStoreName:        row.Sellerstorename,
// // 			CategoryID:             row.Categoryid,
// // 			CategoryName:           row.Categoryname,
// // 			ProductID:              row.Productid,
// // 			ProductTitle:           row.Producttitle,
// // 			ProductDescription:     row.Productdescription,
// // 			ProductPrimaryImageUrl: row.Productprimaryimageurl,
// // 			VariantID:              row.CartVariantID,
// // 			Color:                  row.Color,
// // 			Size:                   row.Size,
// // 			BuyingMode:             "retail",
// // 			UnitPrice:              toDecimal(price).String(),
// // 			HasDiscount:            false,
// // 			DiscountType:           "",
// // 			DiscountValue:          "0.00",
// // 			RequiredQuantity:       quantity,
// // 			WeightGrams:            int32(row.WeightGrams),
// // 			CreatedAt:              now,
// // 		})
// // 	}

// // 	for _, id := range requestedVariantIDs {
// // 		if !foundMap[id] {
// // 			invalidItems = append(invalidItems, map[string]string{
// // 				"variant_id": id.String(),
// // 				"reason":     "variant not found",
// // 			})
// // 		}
// // 	}

// // 	return subtotal, totalWeight, validItems, invalidItems
// // }

// // func toDecimal(val int64) sql.NullString {
// // 	return sql.NullString{Valid: true, String: fmt.Sprintf("%.2f", float64(val)/100)}
// // }

// // // // ------------------------------------------------------------
// // // // 📁 File: internal/services/checkout/checkout_service.go
// // // // 🧠 Handles POST /api/checkout/initiate
// // // //     - Validates user moderation
// // // //     - Fetches snapshot-enriched variant rows
// // // //     - Applies retail vs wholesale pricing logic
// // // //     - Creates checkout session and item rows
// // // //     - Returns valid + invalid items + checkout_session_id

// // // package checkout

// // // import (
// // // 	"context"
// // // 	"database/sql"
// // // 	"fmt"
// // // 	"time"

// // // 	"github.com/google/uuid"

// // // 	"tanmore_backend/internal/db/sqlc"
// // // 	repo "tanmore_backend/internal/repository/checkout"
// // // 	"tanmore_backend/pkg/errors"
// // // 	"tanmore_backend/pkg/sqlnull"
// // // 	"tanmore_backend/pkg/timeutil"
// // // 	uuidutil "tanmore_backend/pkg/uuid"

// // // 	"github.com/shopspring/decimal"
// // // )

// // // // ------------------------------------------------------------
// // // // 📥 Input from handler
// // // type CheckoutFromProductInput struct {
// // // 	UserID    uuid.UUID
// // // 	VariantID uuid.UUID
// // // 	Quantity  int32
// // // }

// // // type CheckoutFromCartInput struct {
// // // 	UserID     uuid.UUID
// // // 	VariantIDs []uuid.UUID
// // // }

// // // // newly added
// // // type GroupedCheckoutItemVariant struct {
// // // 	VariantID        uuid.UUID `json:"variant_id"`
// // // 	Color            string    `json:"color"`
// // // 	Size             string    `json:"size"`
// // // 	BuyingMode       string    `json:"buying_mode"`
// // // 	UnitPrice        string    `json:"unit_price"`
// // // 	HasDiscount      bool      `json:"has_discount"`
// // // 	DiscountType     string    `json:"discount_type"`
// // // 	DiscountValue    string    `json:"discount_value"`
// // // 	RequiredQuantity int32     `json:"required_quantity"`
// // // 	WeightGrams      int32     `json:"weight_grams"`
// // // 	CreatedAt        time.Time `json:"created_at"`
// // // }

// // // type GroupedCheckoutItemProduct struct {
// // // 	ProductID              uuid.UUID                    `json:"product_id"`
// // // 	ProductTitle           string                       `json:"product_title"`
// // // 	ProductDescription     string                       `json:"product_description"`
// // // 	ProductPrimaryImageURL string                       `json:"product_primary_image_url"`
// // // 	Variants               []GroupedCheckoutItemVariant `json:"variants"`
// // // }

// // // type GroupedCheckoutItemCategory struct {
// // // 	CategoryID   uuid.UUID                    `json:"category_id"`
// // // 	CategoryName string                       `json:"category_name"`
// // // 	Products     []GroupedCheckoutItemProduct `json:"products"`
// // // }

// // // type GroupedCheckoutItemSeller struct {
// // // 	SellerID        uuid.UUID                     `json:"seller_id"`
// // // 	SellerStoreName string                        `json:"seller_store_name"`
// // // 	Categories      []GroupedCheckoutItemCategory `json:"categories"`
// // // }

// // // // ------------------------------------------------------------
// // // // 📤 Output to handler
// // // // type CheckoutResult struct {
// // // // 	CheckoutSessionID uuid.UUID           `json:"checkout_session_id,omitempty"`
// // // // 	ValidItems        []sqlc.CheckoutItem `json:"valid_items"`
// // // // 	InvalidItems      []map[string]string `json:"invalid_items"` // variant_id, reason
// // // // }

// // // type CheckoutResult struct {
// // // 	CheckoutSessionID uuid.UUID                   `json:"checkout_session_id"`
// // // 	ValidItems        []CheckoutItemResponse      `json:"valid_items"`
// // // 	InvalidItems      []map[string]string         `json:"invalid_items"`
// // // 	ValidItemsGrouped []GroupedCheckoutItemSeller `json:"valid_items_grouped"`
// // // }

// // // // one more stuct added
// // // type CheckoutItemResponse struct {
// // // 	ID                     uuid.UUID `json:"id"`
// // // 	CheckoutSessionID      uuid.UUID `json:"checkout_session_id"`
// // // 	UserID                 uuid.UUID `json:"user_id"`
// // // 	SellerID               uuid.UUID `json:"seller_id"`
// // // 	SellerStoreName        string    `json:"seller_store_name"`
// // // 	CategoryID             uuid.UUID `json:"category_id"`
// // // 	CategoryName           string    `json:"category_name"`
// // // 	ProductID              uuid.UUID `json:"product_id"`
// // // 	ProductTitle           string    `json:"product_title"`
// // // 	ProductDescription     string    `json:"product_description"`
// // // 	ProductPrimaryImageUrl string    `json:"product_primary_image_url"`
// // // 	VariantID              uuid.UUID `json:"variant_id"`
// // // 	Color                  string    `json:"color"`
// // // 	Size                   string    `json:"size"`
// // // 	BuyingMode             string    `json:"buying_mode"`
// // // 	UnitPrice              string    `json:"unit_price"`
// // // 	HasDiscount            bool      `json:"has_discount"`
// // // 	DiscountType           string    `json:"discount_type"`
// // // 	DiscountValue          string    `json:"discount_value"`
// // // 	RequiredQuantity       int32     `json:"required_quantity"`
// // // 	WeightGrams            int32     `json:"weight_grams"`
// // // 	CreatedAt              time.Time `json:"created_at"`
// // // }

// // // // ------------------------------------------------------------
// // // // 🧱 Service dependencies
// // // type CheckoutServiceDeps struct {
// // // 	Repo repo.CheckoutRepoInterface
// // // }

// // // // 🛠️ Service struct
// // // type CheckoutService struct {
// // // 	Deps CheckoutServiceDeps
// // // }

// // // // 🚀 Constructor
// // // func NewCheckoutService(deps CheckoutServiceDeps) *CheckoutService {
// // // 	return &CheckoutService{Deps: deps}
// // // }

// // // func (s *CheckoutService) FromProduct(
// // // 	ctx context.Context,
// // // 	input CheckoutFromProductInput,
// // // ) (*CheckoutResult, error) {
// // // 	variantIDs := []uuid.UUID{input.VariantID}

// // // 	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
// // // 	if err != nil {
// // // 		return nil, errors.NewNotFoundError("user")
// // // 	}
// // // 	if user.IsBanned || user.IsArchived {
// // // 		return nil, errors.NewAuthError("user is banned or archived")
// // // 	}

// // // 	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
// // // 		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
// // // 			UserID:     sqlnull.UUID(input.UserID),
// // // 			VariantIds: variantIDs,
// // // 		})
// // // 	if err != nil {
// // // 		return nil, errors.NewServerError("failed to fetch snapshot data")
// // // 	}

// // // 	// Inject quantity manually (not from cart)
// // // 	for i := range rows {
// // // 		rows[i].CartRequiredQuantity = sqlnull.Int32(int64(input.Quantity))
// // // 	}

// // // 	// added for debugging
// // // 	for i, row := range rows {
// // // 		fmt.Printf("Row %d → SellerID: %s, SellerStoreName: '%s'\n", i, row.Sellerid.String(), row.Sellerstorename)
// // // 	}

// // // 	// Process variants
// // // 	subtotal, _, validItemsToInsert, validItemsForResponse, invalidItems := s.processVariants(input.UserID, variantIDs, rows)

// // // 	if len(validItemsToInsert) == 0 {
// // // 		reason := "invalid product variant"
// // // 		if len(invalidItems) > 0 {
// // // 			if msg, ok := invalidItems[0]["reason"]; ok {
// // // 				reason = msg
// // // 			}
// // // 		}
// // // 		return nil, errors.NewValidationError("variant_id", reason)
// // // 	}

// // // 	checkoutSessionID := uuidutil.New()
// // // 	now := timeutil.NowUTC()

// // // 	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
// // // 		// 🧾 Step 1: Insert checkout session
// // // 		sessionParams := sqlc.InsertCheckoutSessionParams{
// // // 			ID:           checkoutSessionID,
// // // 			UserID:       input.UserID,
// // // 			Subtotal:     toDecimal(subtotal).String(),
// // // 			TotalPayable: toDecimal(subtotal).String(),
// // // 			// DeliveryCharge:    sqlnull.Float64Ptr(nil), // Will be filled later after shipping
// // // 			DeliveryCharge:    sql.NullString{Valid: false},
// // // 			ShippingAddressID: uuid.NullUUID{}, // Will be updated after shipping selection
// // // 			CreatedAt:         now,
// // // 		}

// // // 		// Debug log (temporary)
// // // 		fmt.Printf("[DEBUG] InsertCheckoutSessionParams: %+v\n", sessionParams)

// // // 		_, err := q.InsertCheckoutSession(ctx, sessionParams)
// // // 		if err != nil {
// // // 			// Optional: log to console for debugging
// // // 			fmt.Printf("[ERROR] Failed to insert checkout session: %v\n", err)
// // // 			return err
// // // 		}

// // // 		// 🧾 Step 2: Insert valid checkout items
// // // 		for i := range validItemsToInsert {
// // // 			validItemsToInsert[i].CheckoutSessionID = checkoutSessionID

// // // 			// Debug each item
// // // 			fmt.Printf("[DEBUG] InsertCheckoutItemParams #%d: %+v\n", i, validItemsToInsert[i])

// // // 			if _, err := q.InsertCheckoutItem(ctx, validItemsToInsert[i]); err != nil {
// // // 				fmt.Printf("[ERROR] Failed to insert checkout item #%d: %v\n", i, err)
// // // 				return err
// // // 			}
// // // 		}
// // // 		return nil
// // // 	})

// // // 	if err != nil {
// // // 		return nil, errors.NewServerError("failed to insert checkout session or items")
// // // 	}
// // // 	return &CheckoutResult{
// // // 		CheckoutSessionID: checkoutSessionID,
// // // 		ValidItems:        validItemsForResponse,
// // // 		InvalidItems:      invalidItems,
// // // 		ValidItemsGrouped: GroupValidCheckoutItems(validItemsForResponse),
// // // 	}, nil

// // // }

// // // func (s *CheckoutService) FromCart(
// // // 	ctx context.Context,
// // // 	input CheckoutFromCartInput,
// // // ) (*CheckoutResult, error) {
// // // 	user, err := s.Deps.Repo.GetUserByID(ctx, input.UserID)
// // // 	if err != nil {
// // // 		return nil, errors.NewNotFoundError("user")
// // // 	}
// // // 	if user.IsBanned || user.IsArchived {
// // // 		return nil, errors.NewAuthError("user is banned or archived")
// // // 	}

// // // 	rows, err := s.Deps.Repo.GetActiveCartVariantSnapshotsByUserAndVariantIDs(ctx,
// // // 		sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams{
// // // 			UserID:     sqlnull.UUID(input.UserID),
// // // 			VariantIds: input.VariantIDs,
// // // 		})
// // // 	if err != nil {
// // // 		return nil, errors.NewServerError("failed to fetch snapshot data")
// // // 	}

// // // 	// added for debugging
// // // 	for i, row := range rows {
// // // 		fmt.Printf("Row %d → SellerID: %s, SellerStoreName: '%s'\n", i, row.Sellerid.String(), row.Sellerstorename)
// // // 	}

// // // 	// subtotal, _, validItems, invalidItems := s.processVariants(input.UserID, input.VariantIDs, rows)
// // // 	subtotal, _, validItemsToInsert, validItemsForResponse, invalidItems := s.processVariants(input.UserID, input.VariantIDs, rows)

// // // 	if len(validItemsToInsert) == 0 {
// // // 		return &CheckoutResult{
// // // 			InvalidItems: invalidItems,
// // // 			// ValidItems:   convertToCheckoutItems(nil), // empty slice
// // // 			ValidItems: []CheckoutItemResponse{},
// // // 		}, nil
// // // 	}

// // // 	checkoutSessionID := uuidutil.New()
// // // 	now := timeutil.NowUTC()

// // // 	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
// // // 		// 🧾 Step 1: Insert checkout session
// // // 		sessionParams := sqlc.InsertCheckoutSessionParams{
// // // 			ID:           checkoutSessionID,
// // // 			UserID:       input.UserID,
// // // 			Subtotal:     toDecimal(subtotal).String(),
// // // 			TotalPayable: toDecimal(subtotal).String(),
// // // 			// DeliveryCharge:    "0.00", // Will be filled later after shipping
// // // 			DeliveryCharge:    sql.NullString{Valid: false},
// // // 			ShippingAddressID: uuid.NullUUID{}, // Will be updated after shipping selection
// // // 			CreatedAt:         now,
// // // 		}

// // // 		// Debug log (temporary)
// // // 		fmt.Printf("[DEBUG] InsertCheckoutSessionParams: %+v\n", sessionParams)

// // // 		_, err := q.InsertCheckoutSession(ctx, sessionParams)
// // // 		if err != nil {
// // // 			// Optional: log to console for debugging
// // // 			fmt.Printf("[ERROR] Failed to insert checkout session: %v\n", err)
// // // 			return err
// // // 		}

// // // 		// 🧾 Step 2: Insert valid checkout items
// // // 		for i := range validItemsToInsert {
// // // 			validItemsToInsert[i].CheckoutSessionID = checkoutSessionID

// // // 			// Debug each item
// // // 			fmt.Printf("[DEBUG] InsertCheckoutItemParams #%d: %+v\n", i, validItemsToInsert[i])

// // // 			if _, err := q.InsertCheckoutItem(ctx, validItemsToInsert[i]); err != nil {
// // // 				fmt.Printf("[ERROR] Failed to insert checkout item #%d: %v\n", i, err)
// // // 				return err
// // // 			}
// // // 		}
// // // 		return nil
// // // 	})

// // // 	if err != nil {
// // // 		return nil, errors.NewServerError("failed to insert checkout session or items")
// // // 	}

// // // 	return &CheckoutResult{
// // // 		CheckoutSessionID: checkoutSessionID,
// // // 		ValidItems:        validItemsForResponse,
// // // 		InvalidItems:      invalidItems,
// // // 		ValidItemsGrouped: GroupValidCheckoutItems(validItemsForResponse),
// // // 	}, nil

// // // }

// // // func (s *CheckoutService) processVariants(
// // // 	userID uuid.UUID,
// // // 	requestedVariantIDs []uuid.UUID,
// // // 	rows []sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow,
// // // ) (int64, int64, []sqlc.InsertCheckoutItemParams, []CheckoutItemResponse, []map[string]string) {

// // // 	var (
// // // 		subtotal           int64
// // // 		totalWeight        int64
// // // 		validItems         []sqlc.InsertCheckoutItemParams
// // // 		validItemsResponse []CheckoutItemResponse
// // // 		invalidItems       []map[string]string
// // // 		now                = timeutil.NowUTC()
// // // 	)

// // // 	foundMap := make(map[uuid.UUID]bool)

// // // 	for idx, row := range rows {
// // // 		foundMap[row.CartVariantID] = true

// // // 		// 🧪 Debug: Check incoming row store name
// // // 		fmt.Printf("[ROW %d] seller_store_name from row: '%s'\n", idx, row.Sellerstorename)

// // // 		quantity := int32(0)
// // // 		if row.CartRequiredQuantity.Valid {
// // // 			quantity = row.CartRequiredQuantity.Int32
// // // 		}

// // // 		if !row.Issellerapproved || row.Issellerarchived || row.Issellerbanned ||
// // // 			!row.Isproductapproved || row.Isproductarchived || row.Isproductbanned ||
// // // 			row.Isvariantarchived || !row.Isvariantinstock || quantity == 0 {

// // // 			invalidItems = append(invalidItems, map[string]string{
// // // 				"variant_id":    row.CartVariantID.String(),
// // // 				"reason":        "variant unavailable due to moderation or stock",
// // // 				"product_id":    row.Productid.String(),
// // // 				"product_title": row.Producttitle,
// // // 				"color":         row.Color,
// // // 				"size":          row.Size,
// // // 			})
// // // 			continue
// // // 		}

// // // 		// [pricing logic remains unchanged...]
// // // 		buyingMode := "retail"
// // // 		unitPrice := row.Retailprice
// // // 		hasDiscount := false
// // // 		discountType := ""
// // // 		discountValue := int64(0)

// // // 		if row.Haswholesaleenabled && row.Wholesaleminquantity.Valid && quantity >= row.Wholesaleminquantity.Int32 {
// // // 			buyingMode = "wholesale"
// // // 			if row.Wholesaleprice.Valid {
// // // 				unitPrice = row.Wholesaleprice.Int64
// // // 			}
// // // 			if row.Haswholesalediscount && row.Wholesalediscount.Valid && row.Wholesalediscounttype.Valid {
// // // 				hasDiscount = true
// // // 				discountType = row.Wholesalediscounttype.String
// // // 				discountValue = row.Wholesalediscount.Int64
// // // 				switch discountType {
// // // 				case "flat":
// // // 					unitPrice -= discountValue
// // // 				case "percentage":
// // // 					unitPrice -= (unitPrice * discountValue) / 100
// // // 				}
// // // 				if unitPrice < 0 {
// // // 					unitPrice = 0
// // // 				}
// // // 			}
// // // 		} else {
// // // 			if row.Hasretaildiscount && row.Retaildiscount.Valid && row.Retaildiscounttype.Valid {
// // // 				hasDiscount = true
// // // 				discountType = row.Retaildiscounttype.String
// // // 				discountValue = row.Retaildiscount.Int64
// // // 				switch discountType {
// // // 				case "flat":
// // // 					unitPrice -= discountValue
// // // 				case "percentage":
// // // 					unitPrice -= (unitPrice * discountValue) / 100
// // // 				}
// // // 				if unitPrice < 0 {
// // // 					unitPrice = 0
// // // 				}
// // // 			}
// // // 		}

// // // 		itemTotal := unitPrice * int64(quantity)
// // // 		subtotal += itemTotal
// // // 		totalWeight += int64(quantity) * int64(row.WeightGrams)

// // // 		item := sqlc.InsertCheckoutItemParams{
// // // 			ID:                     uuidutil.New(),
// // // 			CheckoutSessionID:      uuid.UUID{},
// // // 			UserID:                 userID,
// // // 			SellerID:               row.Sellerid,
// // // 			SellerStoreName:        row.Sellerstorename,
// // // 			CategoryID:             row.Categoryid,
// // // 			CategoryName:           row.Categoryname,
// // // 			ProductID:              row.Productid,
// // // 			ProductTitle:           row.Producttitle,
// // // 			ProductDescription:     row.Productdescription,
// // // 			ProductPrimaryImageUrl: row.Productprimaryimageurl,
// // // 			VariantID:              row.CartVariantID,
// // // 			Color:                  row.Color,
// // // 			Size:                   row.Size,
// // // 			BuyingMode:             buyingMode,
// // // 			UnitPrice:              toDecimal(unitPrice).String(),
// // // 			HasDiscount:            hasDiscount,
// // // 			DiscountType:           discountType,
// // // 			DiscountValue:          toDecimal(discountValue).String(),
// // // 			RequiredQuantity:       quantity,
// // // 			WeightGrams:            int32(row.WeightGrams),
// // // 			CreatedAt:              now,
// // // 		}

// // // 		// 🧪 Debug: Check right after creating item struct
// // // 		fmt.Printf("[INSERT PARAMS %d] seller_store_name before append: '%s'\n", idx, item.SellerStoreName)

// // // 		validItems = append(validItems, item)

// // // 		// 🧪 Debug: Final check in response
// // // 		fmt.Printf("[RESPONSE BUILD %d] seller_store_name before response append: '%s'\n", idx, item.SellerStoreName)

// // // 		validItemsResponse = append(validItemsResponse, CheckoutItemResponse{
// // // 			ID:                     item.ID,
// // // 			CheckoutSessionID:      item.CheckoutSessionID,
// // // 			UserID:                 item.UserID,
// // // 			SellerID:               item.SellerID,
// // // 			SellerStoreName:        item.SellerStoreName,
// // // 			CategoryID:             item.CategoryID,
// // // 			CategoryName:           item.CategoryName,
// // // 			ProductID:              item.ProductID,
// // // 			ProductTitle:           item.ProductTitle,
// // // 			ProductDescription:     item.ProductDescription,
// // // 			ProductPrimaryImageUrl: item.ProductPrimaryImageUrl,
// // // 			VariantID:              item.VariantID,
// // // 			Color:                  item.Color,
// // // 			Size:                   item.Size,
// // // 			BuyingMode:             item.BuyingMode,
// // // 			UnitPrice:              item.UnitPrice,
// // // 			HasDiscount:            item.HasDiscount,
// // // 			DiscountType:           item.DiscountType,
// // // 			DiscountValue:          item.DiscountValue,
// // // 			RequiredQuantity:       item.RequiredQuantity,
// // // 			WeightGrams:            item.WeightGrams,
// // // 			CreatedAt:              item.CreatedAt,
// // // 		})
// // // 	}

// // // 	// Handle not-found variants
// // // 	for _, id := range requestedVariantIDs {
// // // 		if !foundMap[id] {
// // // 			invalidItems = append(invalidItems, map[string]string{
// // // 				"variant_id":    id.String(),
// // // 				"reason":        "variant not found in system",
// // // 				"product_id":    "00000000-0000-0000-0000-000000000000",
// // // 				"product_title": "",
// // // 				"color":         "",
// // // 				"size":          "",
// // // 			})
// // // 		}
// // // 	}

// // // 	return subtotal, totalWeight, validItems, validItemsResponse, invalidItems
// // // }

// // // // Convert int64 price to DECIMAL(10,2)
// // // func toDecimal(val int64) decimal.Decimal {
// // // 	return decimal.NewFromInt(val).Div(decimal.NewFromInt(100))
// // // }

// // // // Convert params to public-facing CheckoutItem struct
// // // func convertToCheckoutItems(params []CheckoutItemResponse) []sqlc.CheckoutItem {
// // // 	items := make([]sqlc.CheckoutItem, 0, len(params))
// // // 	for _, p := range params {
// // // 		items = append(items, sqlc.CheckoutItem{
// // // 			ID:                     p.ID,
// // // 			CheckoutSessionID:      p.CheckoutSessionID,
// // // 			UserID:                 p.UserID,
// // // 			SellerID:               p.SellerID,
// // // 			CategoryID:             p.CategoryID,
// // // 			CategoryName:           p.CategoryName,
// // // 			ProductID:              p.ProductID,
// // // 			ProductTitle:           p.ProductTitle,
// // // 			ProductDescription:     p.ProductDescription,
// // // 			ProductPrimaryImageUrl: p.ProductPrimaryImageUrl,
// // // 			VariantID:              p.VariantID,
// // // 			Color:                  p.Color,
// // // 			Size:                   p.Size,
// // // 			BuyingMode:             p.BuyingMode,
// // // 			UnitPrice:              p.UnitPrice,
// // // 			HasDiscount:            p.HasDiscount,
// // // 			DiscountType:           p.DiscountType,
// // // 			DiscountValue:          p.DiscountValue,
// // // 			RequiredQuantity:       p.RequiredQuantity,
// // // 			WeightGrams:            p.WeightGrams,
// // // 			CreatedAt:              p.CreatedAt,
// // // 		})
// // // 	}
// // // 	return items
// // // }

// // // // newly added

// // // func GroupValidCheckoutItems(
// // // 	items []CheckoutItemResponse,
// // // ) []GroupedCheckoutItemSeller {

// // // 	sellerMap := make(map[uuid.UUID]*GroupedCheckoutItemSeller)

// // // 	for _, item := range items {
// // // 		// 🧱 Seller level
// // // 		seller, ok := sellerMap[item.SellerID]
// // // 		if !ok {
// // // 			seller = &GroupedCheckoutItemSeller{
// // // 				SellerID:        item.SellerID,
// // // 				SellerStoreName: item.SellerStoreName,
// // // 				Categories:      []GroupedCheckoutItemCategory{},
// // // 			}
// // // 			sellerMap[item.SellerID] = seller
// // // 		}

// // // 		// 🧱 Category level
// // // 		var category *GroupedCheckoutItemCategory
// // // 		for i := range seller.Categories {
// // // 			if seller.Categories[i].CategoryID == item.CategoryID {
// // // 				category = &seller.Categories[i]
// // // 				break
// // // 			}
// // // 		}
// // // 		if category == nil {
// // // 			seller.Categories = append(seller.Categories, GroupedCheckoutItemCategory{
// // // 				CategoryID:   item.CategoryID,
// // // 				CategoryName: item.CategoryName,
// // // 				Products:     []GroupedCheckoutItemProduct{},
// // // 			})
// // // 			category = &seller.Categories[len(seller.Categories)-1]
// // // 		}

// // // 		// 🧱 Product level
// // // 		var product *GroupedCheckoutItemProduct
// // // 		for i := range category.Products {
// // // 			if category.Products[i].ProductID == item.ProductID {
// // // 				product = &category.Products[i]
// // // 				break
// // // 			}
// // // 		}
// // // 		if product == nil {
// // // 			category.Products = append(category.Products, GroupedCheckoutItemProduct{
// // // 				ProductID:              item.ProductID,
// // // 				ProductTitle:           item.ProductTitle,
// // // 				ProductDescription:     item.ProductDescription,
// // // 				ProductPrimaryImageURL: item.ProductPrimaryImageUrl,
// // // 				Variants:               []GroupedCheckoutItemVariant{},
// // // 			})
// // // 			product = &category.Products[len(category.Products)-1]
// // // 		}

// // // 		// 🧱 Variant level
// // // 		product.Variants = append(product.Variants, GroupedCheckoutItemVariant{
// // // 			VariantID:        item.VariantID,
// // // 			Color:            item.Color,
// // // 			Size:             item.Size,
// // // 			BuyingMode:       item.BuyingMode,
// // // 			UnitPrice:        item.UnitPrice,
// // // 			HasDiscount:      item.HasDiscount,
// // // 			DiscountType:     item.DiscountType,
// // // 			DiscountValue:    item.DiscountValue,
// // // 			RequiredQuantity: item.RequiredQuantity,
// // // 			WeightGrams:      item.WeightGrams,
// // // 			CreatedAt:        item.CreatedAt,
// // // 		})
// // // 	}

// // // 	// Convert map to slice
// // // 	var grouped []GroupedCheckoutItemSeller
// // // 	for _, s := range sellerMap {
// // // 		grouped = append(grouped, *s)
// // // 	}
// // // 	return grouped
// // // }
