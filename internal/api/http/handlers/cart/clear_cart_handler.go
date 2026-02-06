// 📁 File: internal/api/http/handlers/cart/clear_cart_handler.go

package cart

import (
	"net/http"

	service "tanmore_backend/internal/services/cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

type ClearCartHandler struct {
	Service *service.ClearCartService
}

func NewClearCartHandler(service *service.ClearCartService) *ClearCartHandler {
	return &ClearCartHandler{Service: service}
}

// 🔁 DELETE /api/cart/clear
func (h *ClearCartHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// fmt.Println("✅ Parsed authenticated user ID:", userID)

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

	// 4️⃣ Return success
	response.OK(w, "Cart cleared", map[string]string{
		"status": result.Status,
	})
}
