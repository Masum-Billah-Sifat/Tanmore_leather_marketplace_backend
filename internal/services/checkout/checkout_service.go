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
	if err != nil || user.IsBanned || user.IsArchived {
		return nil, errors.NewAuthError("unauthorized user")
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

// --------------------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------------------
// ────────────────────────────────────────
// A) WHEN IT COMES FROM PRODUCT PAGE (“Buy Now”) → `FromProduct()`
// ────────────────────────────────────────

// Input you receive:

// * `UserID`
// * `VariantID` (the specific size/color variant)
// * `Quantity`

// Goal:

// * Create checkout session + one checkout item (based on product snapshot), without touching cart.

// Step 1: Make a list of variants

// * Code: `variantIDs := []uuid.UUID{input.VariantID}`
// * Why: The processing function is written to handle “a list”, even if it’s just one.

// Step 2: Check the user is allowed (User moderation check)

// * Code calls: `Repo.GetUserByID(ctx, input.UserID)`
// * DB table: `users`
// * Validation:

//   * if user not found → `ErrAuthUserNotFound()`
//   * if archived → `ErrAuthArchivedUser()`
//   * if banned → `ErrAuthBannedUser()`
// * If any fails, stop here.

// Step 3: Fetch the “snapshot” of the product variant

// * Code calls: `Repo.GetVariantSnapshotsByVariantIDs(ctx, variantIDs)`
// * DB table: `product_variant_snapshots`
// * What is snapshot? A frozen/cached record of seller/product/variant info (title, price, weight, moderation flags, stock status, etc.) that checkout can rely on.
// * IMPORTANT: This step does NOT look at `cart_items` at all.

// Step 4: Prepare quantity mapping

// * Code:

//   * `qtyByVariant := map[uuid.UUID]int32{ input.VariantID: input.Quantity }`
// * Why mapping? Because snapshots list may contain multiple variants, and you need to know quantity per variant.

// Step 5: Validate snapshot + build checkout items (the processor)

// * Code calls: `processSnapshotsForCheckout(...)`
// * This function does the “core logic”:

//   * (a) Check quantity > 0
//   * (b) Check moderation flags:

//     * seller banned/archived/not approved
//     * product banned/archived/not approved
//     * category archived
//     * variant archived / out of stock flag false
//   * (c) Stock check (strict):

//     * if quantity > stockamount → invalid
//   * (d) Compute:

//     * line price = unitPrice * qty
//     * subtotal += line price
//     * totalWeight += qty * weight
//   * (e) Build `InsertCheckoutItemParams` struct:

//     * seller id, product title, variant id, unit price, required qty, etc.

// At the end of processor:

// * If items list is empty, you return validation error:

//   * `variant_id: variant unavailable` or some reason like “insufficient stock”, “moderated”, etc.

// Step 6: Choose platform promotion (optional discount)

// * Code calls: `selectPlatformPromotion(ctx, subtotal)`
// * DB table: promotion table (whatever `GetActivePlatformPromotions` reads)
// * It:

//   * loads promotions
//   * checks time window (start/end)
//   * checks min/max cart value
//   * picks the lowest priority promo
//   * calculates promoAmount (flat or percent)
// * Result:

//   * promo (chosen promo row or empty)
//   * promoAmount (int64 cents)

// Step 7: Compute total payable

// * `totalPayable = max(subtotal - promoAmount, 0)`

// Step 8: Create checkout session + checkout items inside a transaction

// * Code: `Repo.WithTx(ctx, func(q *sqlc.Queries) error { ... })`
// * Inside transaction:

// 8.1 Insert one row in `checkout_sessions`

// * Query: `InsertCheckoutSession`
// * Table: `checkout_sessions`
// * Contains:

//   * subtotal, totalWeight
//   * promo metadata saved
//   * status = `awaiting_shipping_info`
//   * paymentMethod = `cod`
//   * shipping address empty for now

// 8.2 Insert one row per item in `checkout_items`

// * Query: `InsertCheckoutItem` (loop)
// * Table: `checkout_items`
// * Each row stores the snapshot-based info (so later price/title changes don’t break the checkout history)

// Step 9: Return minimal response

// * `checkout_session_id` + message

// So product page flow summary (kid version):

// 1. Check user ok
// 2. Read product snapshot from `product_variant_snapshots`
// 3. Validate stock/moderation
// 4. Compute subtotal + weight
// 5. Pick promo
// 6. Insert session + items (transaction)
// 7. Return session id

// ────────────────────────────────────────
// B) WHEN IT COMES FROM CART PAGE → `FromCart()`
// ────────────────────────────────────────

// Input you receive:

// * `UserID`
// * `VariantIDs []uuid.UUID` (selected items from cart, or all cart items)

// Goal:

// * Only allow checkout for items that are actually in the user’s cart (and active), and use the cart quantities.

// Step 1: Validate user

// * Same logic: `GetUserByID`, block archived/banned.

// Step 2: Fetch cart rows joined with snapshot rows

// * Code calls: `GetActiveCartVariantSnapshotsByUserAndVariantIDs(...)`
// * DB tables involved:

//   * `cart_items` (ci)
//   * `product_variant_snapshots` (pvs)
// * Query logic:

//   * Only active cart items: `ci.is_active = true`
//   * Only this user: `ci.user_id = $userID`
//   * Only requested variants: `ci.variant_id = ANY($variant_ids)`
// * It returns rows that contain:

//   * cart info (required_quantity)
//   * snapshot info (seller/product/variant details)

// Why join for cart flow?
// Because cart flow wants:

// * proof that the item is really in cart
// * plus the snapshot info needed to make checkout items

// Step 3: Process + validate variants (cart processor)

// * Code calls: `processVariants(userID, requestedVariantIDs, rows)`
// * In this processor:

//   * Reads qty from `cart_required_quantity`
//   * Checks moderation flags + variant instock + not archived
//   * (NOTE: your current cart processor does NOT enforce `stockamount` cap; it only checks `isvariantinstock`. You might want to add stockamount check here too, same as product flow.)
//   * Computes subtotal + weight
//   * Builds checkout item structs

// Also it checks missing variants:

// * If you requested a variant ID but it didn’t appear in returned rows, it means:

//   * it’s not in cart OR not active
// * Then it adds invalid reason: “variant not found”

// Step 4: Choose platform promo

// * Same promo selection logic

// Step 5: Insert checkout session + items in transaction

// * Same as product flow:

//   * insert into `checkout_sessions`
//   * insert into `checkout_items`

// Step 6: Return minimal response

// * session id + message

// Cart page flow summary (kid version):

// 1. Check user ok
// 2. Read cart items + snapshots using join (must be in cart)
// 3. Validate each item + quantity
// 4. Compute subtotal + weight
// 5. Pick promo
// 6. Insert session + items (transaction)
// 7. Return session id

// ────────────────────────────────────────
// C) “ALGORITHMIC PICTURE” (easy mental model)
// ────────────────────────────────────────

// Two inputs → one output.

// Input A: Product Page

// * variantID + quantity
// * snapshots fetched directly from `product_variant_snapshots`
// * cart is ignored

// Input B: Cart Page

// * variantIDs (cart selections)
// * snapshots fetched through `cart_items JOIN product_variant_snapshots`
// * cart membership is enforced

// Both produce:

// * `checkout_sessions` row (1)
// * `checkout_items` rows (N)

// ────────────────────────────────────────
// D) What changed “before vs after” (important for future you)
// ────────────────────────────────────────

// Before (buggy):

// * Product page was also calling the CART JOIN query.
// * If product isn’t in cart → JOIN returns no rows → “variant not found”.

// After (fixed):

// * Product page uses snapshots-only query.
// * Cart page still uses cart join query (good because it enforces “must be in cart”).

// ────────────────────────────────────────
// E) Two small “future-safety” notes worth saving
// ────────────────────────────────────────

// 1. Stock enforcement inconsistency

// * Product flow checks `qty <= stockamount`.
// * Cart flow currently does NOT check `stockamount`, only checks `isvariantinstock`.
//   You should decide: do you want strict stock in cart flow too? Usually yes.

// 2. Discount metadata vs price calculation

// * You’re storing discount metadata in checkout items, but not modifying the unit price/subtotal with retail discount.
//   That’s okay if “retail discount” is only informational right now, but later you’ll want a single source of truth for pricing rules.

// If you want, I can write a short comment block you can paste directly above `FromProduct()` and `FromCart()` that documents this in 10–15 lines, so future dev sees it instantly.
