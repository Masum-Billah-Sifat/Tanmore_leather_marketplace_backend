// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/checkout/confirm_order_handler.go
// 🧠 Handles POST /api/checkout/{checkout_session_id}/confirm
//     - Extracts user_id from access token (middleware)
//     - Parses checkout_session_id from URL
//     - Calls service to validate + create order (COD only)
//     - Returns created order object

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
type ConfirmOrderHandler struct {
	Service *service.ConfirmOrderService
}

// 🏗️ Constructor
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

	// 3️⃣ Call service layer
	order, err := h.Service.Start(ctx, service.ConfirmOrderInput{
		UserID:            userID,
		CheckoutSessionID: checkoutSessionID,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return response
	response.OK(w, "Order confirmed successfully", map[string]interface{}{
		"order": order,
	})
}
