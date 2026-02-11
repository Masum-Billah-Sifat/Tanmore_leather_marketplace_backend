// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/checkout/confirm_order_handler.go
// 🧠 Handles POST /api/checkout/{checkout_session_id}/confirm

package checkout

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	service "tanmore_backend/internal/services/checkout"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ------------------------------------------------------------
// ✅ Response DTO (clean JSON: no {String, Valid} / {Time, Valid})

type OrderResponse struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	CheckoutSessionID uuid.UUID `json:"checkout_session_id"`
	OrderCode         string    `json:"order_code"`

	Subtotal    string  `json:"subtotal"`
	ShippingFee *string `json:"shipping_fee"`
	TotalAmount string  `json:"total_amount"`
	Currency    *string `json:"currency"`

	PaymentMethod   string     `json:"payment_method"`
	PaymentStatus   *string    `json:"payment_status"`
	LatestPaymentID *uuid.UUID `json:"latest_payment_id"`

	Status *string `json:"status"`

	PlatformDiscountType          *string `json:"platform_discount_type"`
	PlatformDiscountValue         *string `json:"platform_discount_value"`
	PlatformDiscountAmountApplied *string `json:"platform_discount_amount_applied"`
	IsPlatformDiscountApplied     bool    `json:"is_platform_discount_applied"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func nullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func nullTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	t := nt.Time
	return &t
}

func nullUUIDPtr(n uuid.NullUUID) *uuid.UUID {
	if !n.Valid {
		return nil
	}
	u := n.UUID
	return &u
}

// ------------------------------------------------------------
// 📦 Handler struct

type ConfirmOrderHandler struct {
	Service *service.ConfirmOrderService
}

func NewConfirmOrderHandler(service *service.ConfirmOrderService) *ConfirmOrderHandler {
	return &ConfirmOrderHandler{Service: service}
}

// 🔁 POST /api/checkout/{checkout_session_id}/confirm
func (h *ConfirmOrderHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Step 1️⃣: Extract user ID from access token context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	if rawUserID == nil {
		response.Unauthorized(w, errors.ErrAuthMissingToken())
		return
	}

	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, errors.ErrAuthInvalidUserID())
		return
	}

	// 2️⃣ Parse checkout_session_id from URL
	sessionIDStr := chi.URLParam(r, "checkout_session_id")
	checkoutSessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		fmt.Println("❌ Invalid checkout_session_id format:", err)
		response.BadRequest(w, errors.NewValidationError("checkout_session_id", "invalid UUID"))
		return
	}

	log.Printf("[ConfirmOrderHandler] HIT user_ctx=%v url_session_id=%s", rawUserID, sessionIDStr)

	// 3️⃣ Call service layer
	order, err := h.Service.Start(ctx, service.ConfirmOrderInput{
		UserID:            userID,
		CheckoutSessionID: checkoutSessionID,
	})
	if err != nil {
		log.Printf("[ConfirmOrderHandler] FAIL user_id=%s session_id=%s err=%v", userID, checkoutSessionID, err)
		response.ServerError(w, err)
		return
	}

	log.Printf("[ConfirmOrderHandler] OK user_id=%s session_id=%s order_id=%s order_code=%s",
		userID, checkoutSessionID, order.ID, order.OrderCode)

	// 4️⃣ Map sqlc.Order -> clean DTO
	orderRes := OrderResponse{
		ID:                order.ID,
		UserID:            order.UserID,
		CheckoutSessionID: order.CheckoutSessionID,
		OrderCode:         order.OrderCode,

		Subtotal:    order.Subtotal,
		ShippingFee: nullStringPtr(order.ShippingFee),
		TotalAmount: order.TotalAmount,
		Currency:    nullStringPtr(order.Currency),

		PaymentMethod:   order.PaymentMethod,
		PaymentStatus:   nullStringPtr(order.PaymentStatus),
		LatestPaymentID: nullUUIDPtr(order.LatestPaymentID),

		Status: nullStringPtr(order.Status),

		PlatformDiscountType:          nullStringPtr(order.PlatformDiscountType),
		PlatformDiscountValue:         nullStringPtr(order.PlatformDiscountValue),
		PlatformDiscountAmountApplied: nullStringPtr(order.PlatformDiscountAmountApplied),
		IsPlatformDiscountApplied:     order.IsPlatformDiscountApplied,

		CreatedAt: nullTimePtr(order.CreatedAt),
		UpdatedAt: nullTimePtr(order.UpdatedAt),
	}

	// 5️⃣ Return response (same envelope, clean "order")
	response.OK(w, "Order confirmed successfully", map[string]interface{}{
		"order": orderRes,
	})
}
