// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/checkout/get_checkout_session_details_handler.go
// 🧠 Handles GET /api/checkout/{checkout_session_id}
//     - Extracts user_id from access token (middleware)
//     - Parses checkout_session_id from URL
//     - Calls service to fetch session, items, snapshots, validations
//     - Returns full session + shipping address + valid/invalid items

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
type GetCheckoutSessionDetailsHandler struct {
	Service *service.GetCheckoutSessionDetailsService
}

// 🏗️ Constructor
func NewGetCheckoutSessionDetailsHandler(service *service.GetCheckoutSessionDetailsService) *GetCheckoutSessionDetailsHandler {
	return &GetCheckoutSessionDetailsHandler{Service: service}
}

// 🔁 GET /api/checkout/{checkout_session_id}
func (h *GetCheckoutSessionDetailsHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
	result, err := h.Service.Start(ctx, service.GetCheckoutSessionDetailsInput{
		UserID:            userID,
		CheckoutSessionID: checkoutSessionID,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return response
	response.OK(w, "Checkout session details fetched successfully", map[string]interface{}{
		"checkout_session":             result.CheckoutSession,
		"has_shipping_address_details": result.HasShippingAddressDetails,
		"shipping_address":             result.ShippingAddress,
		"valid_items":                  result.ValidItems,
		"invalid_items":                result.InvalidItems,
	})

}
