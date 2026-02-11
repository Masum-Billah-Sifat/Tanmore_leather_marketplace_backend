// ------------------------------------------------------------
// 📁 File: internal/services/checkout/confirm_order_service.go
// 🧐 Confirms a checkout session into a final COD order.

package checkout

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/checkout/confirm_order"
	"tanmore_backend/pkg/errors"
	sqlnull "tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input

type ConfirmOrderInput struct {
	UserID            uuid.UUID
	CheckoutSessionID uuid.UUID
}

// ------------------------------------------------------------
// 🛠️ Service Setup

type ConfirmOrderServiceDeps struct {
	Repo repo.ConfirmOrderRepoInterface
}

type ConfirmOrderService struct {
	Deps ConfirmOrderServiceDeps
}

func NewConfirmOrderService(deps ConfirmOrderServiceDeps) *ConfirmOrderService {
	return &ConfirmOrderService{Deps: deps}
}

// ------------------------------------------------------------
// ✅ Payload DTOs (clean JSON)

type OrderCreatedPayload struct {
	Order           OrderCreatedOrder           `json:"order"`
	ShippingAddress OrderCreatedShippingAddress `json:"shipping_address"`
	Items           []OrderCreatedItem          `json:"items"`
}

type OrderCreatedOrder struct {
	OrderID           uuid.UUID `json:"order_id"`
	OrderCode         string    `json:"order_code"`
	CheckoutSessionID uuid.UUID `json:"checkout_session_id"`
	UserID            uuid.UUID `json:"user_id"`

	Subtotal      string    `json:"subtotal"`
	ShippingFee   *string   `json:"shipping_fee"`
	TotalAmount   string    `json:"total_amount"`
	Currency      *string   `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Status        *string   `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type OrderCreatedShippingAddress struct {
	RecipientName  string `json:"recipient_name"`
	RecipientPhone string `json:"recipient_phone"`
	AddressLine    string `json:"address_line"`
	CityID         int32  `json:"city_id"`
	ZoneID         int32  `json:"zone_id"`
	AreaID         int32  `json:"area_id"`
}

type OrderCreatedItem struct {
	ID                     uuid.UUID `json:"id"`
	CheckoutSessionID      uuid.UUID `json:"checkout_session_id"`
	UserID                 uuid.UUID `json:"user_id"`
	SellerID               uuid.UUID `json:"seller_id"`
	CategoryID             uuid.UUID `json:"category_id"`
	CategoryName           string    `json:"category_name"`
	ProductID              uuid.UUID `json:"product_id"`
	ProductTitle           string    `json:"product_title"`
	ProductDescription     string    `json:"product_description"`
	ProductPrimaryImageUrl string    `json:"product_primary_image_url"`
	VariantID              uuid.UUID `json:"variant_id"`
	Color                  string    `json:"color"`
	Size                   string    `json:"size"`
	BuyingMode             string    `json:"buying_mode"`

	UnitPrice     string  `json:"unit_price"`
	HasDiscount   bool    `json:"has_discount"`
	DiscountType  *string `json:"discount_type"`
	DiscountValue *string `json:"discount_value"`

	RequiredQuantity int32     `json:"required_quantity"`
	WeightGrams      int32     `json:"weight_grams"`
	SellerStoreName  string    `json:"seller_store_name"`
	CreatedAt        time.Time `json:"created_at"`
}

// ------------------------------------------------------------
// 🚀 Entry Point

func (s *ConfirmOrderService) Start(
	ctx context.Context,
	input ConfirmOrderInput,
) (*sqlc.Order, error) {
	var resultOrder *sqlc.Order
	start := time.Now()

	log.Printf("[ConfirmOrder] START user_id=%s checkout_session_id=%s", input.UserID, input.CheckoutSessionID)

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1: Customer validation
		log.Printf("[ConfirmOrder] Step1: validate user user_id=%s", input.UserID)
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			log.Printf("[ConfirmOrder] Step1 FAIL: user not found user_id=%s err=%v", input.UserID, err)
			return errors.ErrAuthUserNotFound()
		}
		if user.IsArchived {
			log.Printf("[ConfirmOrder] Step1 FAIL: user archived user_id=%s", input.UserID)
			return errors.ErrAuthArchivedUser()
		}
		if user.IsBanned {
			log.Printf("[ConfirmOrder] Step1 FAIL: user banned user_id=%s", input.UserID)
			return errors.ErrAuthBannedUser()
		}

		// Step 2: Checkout session fetch + validation
		log.Printf("[ConfirmOrder] Step2: fetch session session_id=%s", input.CheckoutSessionID)
		session, err := q.GetCheckoutSessionByID(ctx, input.CheckoutSessionID)
		if err != nil {
			log.Printf("[ConfirmOrder] Step2 FAIL: session not found session_id=%s err=%v", input.CheckoutSessionID, err)
			return errors.NewNotFoundError("checkout_session")
		}

		log.Printf("[ConfirmOrder] Step2: session loaded id=%s user_id=%s status=%s payment_method=%s shipping_address_valid=%v",
			session.ID, session.UserID, session.Status, session.PaymentMethod, session.ShippingAddressID.Valid)

		if session.UserID != input.UserID {
			log.Printf("[ConfirmOrder] Step2 FAIL: session not owned by user session_user_id=%s req_user_id=%s",
				session.UserID, input.UserID)
			return errors.NewAuthError("checkout session not owned by user")
		}
		if !session.ShippingAddressID.Valid {
			log.Printf("[ConfirmOrder] Step2 FAIL: shipping address missing session_id=%s", session.ID)
			return errors.NewServerError("shipping address missing")
		}
		if session.Status != "ready_to_order" {
			log.Printf("[ConfirmOrder] Step2 FAIL: session not ready status=%s session_id=%s", session.Status, session.ID)
			return errors.NewServerError("checkout session not ready")
		}
		if session.PaymentMethod == "prepaid" {
			log.Printf("[ConfirmOrder] Step2 FAIL: prepaid not supported session_id=%s", session.ID)
			return errors.NewServerError("prepaid not supported yet")
		}

		// Step 3: Fetch checkout items
		log.Printf("[ConfirmOrder] Step3: fetch checkout items session_id=%s", session.ID)
		items, err := q.GetCheckoutItemsBySessionID(ctx, session.ID)
		if err != nil || len(items) == 0 {
			log.Printf("[ConfirmOrder] Step3 FAIL: no checkout items session_id=%s err=%v items_len=%d", session.ID, err, len(items))
			return errors.NewServerError("no checkout items")
		}
		log.Printf("[ConfirmOrder] Step3: items loaded count=%d session_id=%s", len(items), session.ID)

		// Step 3.5: Validate using snapshots
		log.Printf("[ConfirmOrder] Step3.5: validate snapshots")
		variantIDs := make([]uuid.UUID, 0, len(items))
		for _, item := range items {
			variantIDs = append(variantIDs, item.VariantID)
		}
		snaps, err := q.GetProductVariantSnapshotsByVariantIDs(ctx, variantIDs)
		if err != nil {
			log.Printf("[ConfirmOrder] Step3.5 FAIL: variant enrichment failed err=%v", err)
			return errors.NewServerError("variant enrichment failed")
		}
		snapMap := make(map[uuid.UUID]sqlc.GetProductVariantSnapshotsByVariantIDsRow, len(snaps))
		for _, snap := range snaps {
			snapMap[snap.Variantid] = snap
		}

		for _, item := range items {
			snap, ok := snapMap[item.VariantID]
			if !ok || snap.Issellerbanned || snap.Issellerarchived || !snap.Issellerapproved ||
				snap.Isproductbanned || snap.Isproductarchived || !snap.Isproductapproved ||
				snap.Iscategoryarchived || snap.Isvariantarchived ||
				!snap.Isvariantinstock || snap.Stockamount < item.RequiredQuantity {

				log.Printf("[ConfirmOrder] Step3.5 FAIL: unavailable item variant_id=%s product_id=%s qty=%d",
					item.VariantID, item.ProductID, item.RequiredQuantity)
				return errors.NewServerError("checkout contains unavailable item")
			}
		}

		// Step 4: Fetch shipping address
		log.Printf("[ConfirmOrder] Step4: fetch shipping address shipping_id=%s session_id=%s",
			session.ShippingAddressID.UUID, session.ID)

		shipping, err := q.GetShippingAddressByIDAndCheckoutID(ctx, sqlc.GetShippingAddressByIDAndCheckoutIDParams{
			ID:                session.ShippingAddressID.UUID,
			CheckoutSessionID: session.ID,
		})
		if err != nil {
			log.Printf("[ConfirmOrder] Step4 FAIL: shipping not found err=%v", err)
			return errors.NewNotFoundError("shipping_address")
		}

		// Step 5: Create order
		orderID := uuid.New()
		orderCode := GenerateOrderCode()
		now := time.Now().UTC()

		log.Printf("[ConfirmOrder] Step5: create order order_id=%s order_code=%s now=%s", orderID, orderCode, now.Format(time.RFC3339Nano))

		order, err := q.InsertOrderRow(ctx, sqlc.InsertOrderRowParams{
			ID:                orderID,
			UserID:            input.UserID,
			CheckoutSessionID: session.ID,
			OrderCode:         orderCode,

			Subtotal:    session.Subtotal,
			ShippingFee: session.DeliveryCharge,
			TotalAmount: session.TotalPayable,

			Currency:                      sqlnull.String("BDT"),
			PaymentMethod:                 session.PaymentMethod,
			PaymentStatus:                 sqlnull.String(""),
			LatestPaymentID:               uuid.NullUUID{},
			Status:                        sqlnull.String("processing"),
			PlatformDiscountType:          session.PlatformDiscountType,
			PlatformDiscountValue:         session.PlatformDiscountValue,
			PlatformDiscountAmountApplied: session.PlatformDiscountAmountApplied,
			IsPlatformDiscountApplied:     session.IsPlatformDiscountApplied,
			CreatedAt:                     sqlnull.Time(now),
			UpdatedAt:                     sqlnull.Time(now),
		})
		if err != nil {
			log.Printf("[ConfirmOrder] Step5 FAIL: insert order err=%v", err)
			return err
		}
		resultOrder = &order

		// Step 6: Create order items
		log.Printf("[ConfirmOrder] Step6: insert order items count=%d", len(items))
		for _, item := range items {
			// FIX: your unit_price is "15.00" -> Atoi fails -> totalPrice becomes 0 silently.
			unitCents, err := moneyToCentsOne(item.UnitPrice)
			if err != nil {
				log.Printf("[ConfirmOrder] Step6 FAIL: invalid unit_price unit_price=%q item_id=%s err=%v", item.UnitPrice, item.ID, err)
				return errors.NewServerError("invalid unit_price")
			}
			totalCents := unitCents * int64(item.RequiredQuantity)
			totalPriceStr := centsToMoneyOne(totalCents) // "45.00" style

			err = q.InsertOrderItem(ctx, sqlc.InsertOrderItemParams{
				ID:                     uuid.New(),
				OrderID:                order.ID,
				CustomerID:             input.UserID,
				SellerID:               item.SellerID,
				CategoryID:             item.CategoryID,
				CategoryName:           item.CategoryName,
				ProductID:              item.ProductID,
				ProductTitle:           item.ProductTitle,
				ProductDescription:     item.ProductDescription,
				ProductPrimaryImageUrl: item.ProductPrimaryImageUrl,
				VariantID:              item.VariantID,
				Color:                  item.Color,
				Size:                   item.Size,
				BuyingMode:             item.BuyingMode,
				UnitPrice:              item.UnitPrice,
				Quantity:               item.RequiredQuantity,

				TotalPrice: totalPriceStr,

				HasDiscount:        item.HasDiscount,
				DiscountType:       item.DiscountType,
				DiscountValue:      item.DiscountValue,
				WeightGramsPerUnit: item.WeightGrams,
				TotalWeightGrams:   item.WeightGrams * int32(item.RequiredQuantity),
				SellerStoreName:    item.SellerStoreName,
				CreatedAt:          now,
			})
			if err != nil {
				log.Printf("[ConfirmOrder] Step6 FAIL: insert order item err=%v product_id=%s variant_id=%s", err, item.ProductID, item.VariantID)
				return err
			}
		}

		// Step 7: Update checkout session status
		log.Printf("[ConfirmOrder] Step7: update checkout session status -> order_created session_id=%s", session.ID)
		err = q.UpdateCheckoutSessionStatusToOrderCreated(ctx, sqlc.UpdateCheckoutSessionStatusToOrderCreatedParams{
			ID:     session.ID,
			Status: "order_created",
		})
		if err != nil {
			log.Printf("[ConfirmOrder] Step7 FAIL: update session err=%v", err)
			return err
		}

		// Step 8: Insert event
		eventID := uuid.New()
		payload := BuildOrderCreatedPayload(order, shipping, items)

		log.Printf("[ConfirmOrder] Step8: insert event event_id=%s type=order_created", eventID)
		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           eventID,
			Userid:       input.UserID,
			EventType:    "order_created",
			EventPayload: payload,
			DispatchedAt: sqlnull.TimePtr(nil),
			CreatedAt:    now,
		})
		if err != nil {
			log.Printf("[ConfirmOrder] Step8 FAIL: insert event err=%v", err)
			return err
		}

		log.Printf("[ConfirmOrder] DONE tx success order_id=%s order_code=%s", order.ID, order.OrderCode)
		return nil
	})

	if err != nil {
		log.Printf("[ConfirmOrder] END with error user_id=%s session_id=%s dur=%s err=%v",
			input.UserID, input.CheckoutSessionID, time.Since(start), err)
		return nil, err
	}

	log.Printf("[ConfirmOrder] END success user_id=%s session_id=%s dur=%s order_id=%s",
		input.UserID, input.CheckoutSessionID, time.Since(start), resultOrder.ID)

	return resultOrder, nil
}

// ------------------------------------------------------------
// Order code (fix: deterministic math/rand problem)
// Still returns "TNMR-xxxxxx"

func GenerateOrderCode() string {
	n := cryptoRandInt(900000) + 100000 // 100000..999999
	return "TNMR-" + strconv.Itoa(n)
}

func cryptoRandInt(max int) int {
	if max <= 0 {
		return 0
	}
	var b [4]byte
	_, err := rand.Read(b[:])
	if err != nil {
		// fallback (should be extremely rare)
		return int(time.Now().UnixNano() % int64(max))
	}
	v := int(uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3]))
	if v < 0 {
		v = -v
	}
	return v % max
}

// ------------------------------------------------------------
// Payload builder (fix: avoid sql.NullString JSON structs)

func BuildOrderCreatedPayload(
	order sqlc.Order,
	shipping sqlc.ShippingAddress,
	items []sqlc.GetCheckoutItemsBySessionIDRow,
) json.RawMessage {

	payload := OrderCreatedPayload{
		Order: OrderCreatedOrder{
			OrderID:           order.ID,
			OrderCode:         order.OrderCode,
			CheckoutSessionID: order.CheckoutSessionID,
			UserID:            order.UserID,
			Subtotal:          order.Subtotal,
			ShippingFee:       nullStringPtr(order.ShippingFee),
			TotalAmount:       order.TotalAmount,
			Currency:          nullStringPtr(order.Currency),
			PaymentMethod:     order.PaymentMethod,
			Status:            nullStringPtr(order.Status),
			CreatedAt:         order.CreatedAt.Time,
		},
		ShippingAddress: OrderCreatedShippingAddress{
			RecipientName:  shipping.RecipientName,
			RecipientPhone: shipping.RecipientPhone,
			AddressLine:    shipping.AddressLine,
			CityID:         shipping.CityID,
			ZoneID:         shipping.ZoneID,
			AreaID:         shipping.AreaID,
		},
		Items: make([]OrderCreatedItem, 0, len(items)),
	}

	for _, it := range items {
		payload.Items = append(payload.Items, OrderCreatedItem{
			ID:                     it.ID,
			CheckoutSessionID:      it.CheckoutSessionID,
			UserID:                 it.UserID,
			SellerID:               it.SellerID,
			CategoryID:             it.CategoryID,
			CategoryName:           it.CategoryName,
			ProductID:              it.ProductID,
			ProductTitle:           it.ProductTitle,
			ProductDescription:     it.ProductDescription,
			ProductPrimaryImageUrl: it.ProductPrimaryImageUrl,
			VariantID:              it.VariantID,
			Color:                  it.Color,
			Size:                   it.Size,
			BuyingMode:             it.BuyingMode,
			UnitPrice:              it.UnitPrice,
			HasDiscount:            it.HasDiscount,
			DiscountType:           nullStringPtr(it.DiscountType),
			DiscountValue:          nullStringPtr(it.DiscountValue),
			RequiredQuantity:       it.RequiredQuantity,
			WeightGrams:            it.WeightGrams,
			SellerStoreName:        it.SellerStoreName,
			CreatedAt:              it.CreatedAt,
		})
	}

	b, err := json.Marshal(payload)
	if err != nil {
		// keep it safe: return a minimal payload if marshal fails
		fallback := fmt.Sprintf(`{"error":"failed to build payload","order_id":"%s"}`, order.ID.String())
		return json.RawMessage(fallback)
	}
	return b
}

func nullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

// ------------------------------------------------------------
// Money helpers (avoid float, fix "15.00" parsing)

func moneyToCentsOne(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty money string")
	}
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = strings.TrimPrefix(s, "-")
	}
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid money format: %q", s)
	}
	whole := parts[0]
	frac := ""
	if len(parts) == 2 {
		frac = parts[1]
	}

	if whole == "" {
		whole = "0"
	}
	w, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid whole part: %q", whole)
	}

	// normalize fractional to 2 digits
	if len(frac) > 2 {
		// do NOT silently truncate/round in payments
		return 0, fmt.Errorf("too many decimal places: %q", frac)
	}
	for len(frac) < 2 {
		frac += "0"
	}
	f := int64(0)
	if frac != "" {
		ff, err := strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid frac part: %q", frac)
		}
		f = ff
	}
	cents := w*100 + f
	if neg {
		cents = -cents
	}
	return cents, nil
}

func centsToMoneyOne(cents int64) string {
	neg := cents < 0
	if neg {
		cents = -cents
	}
	whole := cents / 100
	frac := cents % 100
	out := fmt.Sprintf("%d.%02d", whole, frac)
	if neg {
		out = "-" + out
	}
	return out
}
