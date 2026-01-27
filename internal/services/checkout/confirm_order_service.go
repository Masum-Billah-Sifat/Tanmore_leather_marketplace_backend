// ------------------------------------------------------------
// 📁 File: internal/services/checkout/confirm_order_service.go
// 🧐 Confirms a checkout session into a final COD order.

package checkout

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
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
// 🚀 Entry Point

func (s *ConfirmOrderService) Start(
	ctx context.Context,
	input ConfirmOrderInput,
) (*sqlc.Order, error) {
	var resultOrder *sqlc.Order

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1: Customer validation
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsArchived {
			return errors.NewAuthError("user is archived")
		}
		if user.IsBanned {
			return errors.NewAuthError("user is banned")
		}

		// Step 2: Checkout session fetch + validation
		session, err := q.GetCheckoutSessionByID(ctx, input.CheckoutSessionID)
		if err != nil {
			return errors.NewNotFoundError("checkout_session")
		}
		if session.UserID != input.UserID {
			return errors.NewAuthError("checkout session not owned by user")
		}
		if !session.ShippingAddressID.Valid {
			return errors.NewServerError("shipping address missing")
		}
		if session.Status != "ready_to_order" {
			return errors.NewServerError("checkout session not ready")
		}
		if session.PaymentMethod == "prepaid" {
			return errors.NewServerError("prepaid not supported yet")
		}

		// Step 3: Fetch checkout items
		items, err := q.GetCheckoutItemsBySessionID(ctx, session.ID)
		if err != nil || len(items) == 0 {
			return errors.NewServerError("no checkout items")
		}

		// Step 3.5: Validate using snapshots
		variantIDs := make([]uuid.UUID, 0, len(items))
		for _, item := range items {
			variantIDs = append(variantIDs, item.VariantID)
		}
		snaps, err := q.GetProductVariantSnapshotsByVariantIDs(ctx, variantIDs)
		if err != nil {
			return errors.NewServerError("variant enrichment failed")
		}
		snapMap := make(map[uuid.UUID]sqlc.GetProductVariantSnapshotsByVariantIDsRow)
		for _, snap := range snaps {
			snapMap[snap.Variantid] = snap
		}
		for _, item := range items {
			snap, ok := snapMap[item.VariantID]
			if !ok || snap.Issellerbanned || snap.Issellerarchived || !snap.Issellerapproved ||
				snap.Isproductbanned || snap.Isproductarchived || !snap.Isproductapproved ||
				snap.Iscategoryarchived || snap.Isvariantarchived ||
				!snap.Isvariantinstock || snap.Stockamount < item.RequiredQuantity {
				return errors.NewServerError("checkout contains unavailable item")
			}
		}

		// Step 4: Fetch shipping address
		shipping, err := q.GetShippingAddressByIDAndCheckoutID(ctx, sqlc.GetShippingAddressByIDAndCheckoutIDParams{
			ID:                session.ShippingAddressID.UUID,
			CheckoutSessionID: session.ID,
		})
		if err != nil {
			return errors.NewNotFoundError("shipping_address")
		}

		// Step 5: Create order
		orderID := uuid.New()
		orderCode := GenerateOrderCode()
		now := time.Now().UTC()
		order, err := q.InsertOrderRow(ctx, sqlc.InsertOrderRowParams{
			ID:                orderID,
			UserID:            input.UserID,
			CheckoutSessionID: session.ID,
			OrderCode:         orderCode,
			// Subtotal:                      session.Subtotal,
			// ShippingFee:                   sqlnull.String(session.DeliveryCharge.String),
			// TotalAmount:                   sqlnull.String(session.TotalPayable.String),

			Subtotal:    session.Subtotal,
			ShippingFee: session.DeliveryCharge, // already sql.NullString
			TotalAmount: session.TotalPayable,   // plain string

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
			return err
		}
		resultOrder = &order

		// Step 6: Create order items
		for _, item := range items {

			unitPrice, _ := strconv.Atoi(item.UnitPrice)
			totalPrice := unitPrice * int(item.RequiredQuantity)

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
				// TotalPrice:             item.UnitPrice * int32(item.RequiredQuantity),
				TotalPrice: strconv.Itoa(totalPrice),

				HasDiscount:        item.HasDiscount,
				DiscountType:       item.DiscountType,
				DiscountValue:      item.DiscountValue,
				WeightGramsPerUnit: item.WeightGrams,
				TotalWeightGrams:   item.WeightGrams * int32(item.RequiredQuantity),
				SellerStoreName:    item.SellerStoreName,
				CreatedAt:          now,
			})
			if err != nil {
				return err
			}
		}

		// Step 7: Update checkout session status
		err = q.UpdateCheckoutSessionStatusToOrderCreated(ctx, sqlc.UpdateCheckoutSessionStatusToOrderCreatedParams{
			ID:     session.ID,
			Status: "order_created",
		})
		if err != nil {
			return err
		}

		// Step 8: Insert event
		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuid.New(),
			Userid:       input.UserID,
			EventType:    "order_created",
			EventPayload: BuildOrderCreatedPayload(order, shipping, items),
			DispatchedAt: sqlnull.TimePtr(nil),
			CreatedAt:    now,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return resultOrder, nil
}

func GenerateOrderCode() string {
	// Example: TNMR-274192
	return "TNMR-" + strconv.Itoa(rand.Intn(900000)+100000)
}

func BuildOrderCreatedPayload(
	order sqlc.Order,
	shipping sqlc.ShippingAddress,
	items []sqlc.GetCheckoutItemsBySessionIDRow,
) json.RawMessage {

	payload := map[string]interface{}{
		"order": map[string]interface{}{
			"order_id":            order.ID,
			"order_code":          order.OrderCode,
			"checkout_session_id": order.CheckoutSessionID,
			"user_id":             order.UserID,
			"subtotal":            order.Subtotal,
			"shipping_fee":        order.ShippingFee.String,
			"total_amount":        order.TotalAmount,
			"currency":            order.Currency.String,
			"payment_method":      order.PaymentMethod,
			"status":              order.Status.String,
			"created_at":          order.CreatedAt.Time,
		},
		"shipping_address": map[string]interface{}{
			"recipient_name":  shipping.RecipientName,
			"recipient_phone": shipping.RecipientPhone,
			"address_line":    shipping.AddressLine,
			"city_id":         shipping.CityID,
			"zone_id":         shipping.ZoneID,
			"area_id":         shipping.AreaID,
		},
		"items": items,
	}

	b, _ := json.Marshal(payload)
	return b
}
