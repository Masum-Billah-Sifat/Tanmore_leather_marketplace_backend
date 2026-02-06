// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/checkout/review_checkout_handler.go
// 🧠 Handles GET /api/checkout/{checkout_session_id}/review
//     - Extracts user_id from access token (middleware)
//     - Parses checkout_session_id from URL
//     - Calls service to review checkout summary
//     - Returns checkout summary with shipping, pricing, valid/invalid items

package checkout

import (
	"fmt"
	"net/http"

	service "tanmore_backend/internal/services/checkout"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// 📦 Handler struct
type ReviewCheckoutHandler struct {
	Service *service.ReviewCheckoutService
}

// 🏗️ Constructor
func NewReviewCheckoutHandler(service *service.ReviewCheckoutService) *ReviewCheckoutHandler {
	return &ReviewCheckoutHandler{Service: service}
}

// 🔁 GET /api/checkout/{checkout_session_id}/review
func (h *ReviewCheckoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 3️⃣ Call service layer
	result, err := h.Service.Start(ctx, service.ReviewCheckoutInput{
		UserID:            userID,
		CheckoutSessionID: checkoutSessionID,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return response
	response.OK(w, "Checkout summary fetched successfully", map[string]interface{}{
		"checkout_session_id": checkoutSessionID,
		"shipping_address":    result.ShippingAddress,
		"payment_method":      result.CheckoutSession.PaymentMethod,
		"subtotal":            result.CheckoutSession.Subtotal,
		"delivery_charge":     result.CheckoutSession.DeliveryCharge,
		"total_weight_grams":  result.CheckoutSession.TotalWeightGrams,
		"platform_discount": map[string]interface{}{
			"type":           result.CheckoutSession.PlatformDiscountType,
			"value":          result.CheckoutSession.PlatformDiscountValue,
			"amount_applied": result.CheckoutSession.PlatformDiscountAmountApplied,
		},
		"total_payable": result.CheckoutSession.TotalPayable,
		"valid_items":   result.ValidItems,
		"invalid_items": result.InvalidItems,
	})
}
