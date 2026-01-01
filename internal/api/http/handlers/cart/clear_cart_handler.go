// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/clear_cart_handler.go
// 🧠 Handles DELETE /api/cart/clear
//     - Extracts customer user_id from context
//     - Calls the service layer to clear all active cart items
//     - Returns cart cleared or already empty status

package cart

import (
	"net/http"

	service "tanmore_backend/internal/services/cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct
type ClearCartHandler struct {
	Service *service.ClearCartService
}

// 🏗️ Constructor
func NewClearCartHandler(service *service.ClearCartService) *ClearCartHandler {
	return &ClearCartHandler{Service: service}
}

// 🔁 DELETE /api/cart/clear
func (h *ClearCartHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Get customer user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, errors.NewAuthError("invalid access token"))
		return
	}

	// 2️⃣ Build service input
	input := service.ClearCartInput{
		UserID: userID,
	}

	// 3️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return response
	response.OK(w, "Cart status", map[string]string{
		"status": result.Status,
	})
}
