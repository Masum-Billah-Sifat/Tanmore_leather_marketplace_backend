// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/get_all_cart_items_handler.go
// 🧠 Handles GET /api/cart/items
//     - Extracts customer user_id from context
//     - Calls service to fetch enriched + grouped cart data
//     - Returns grouped valid_items and flat invalid_items
//     - Handles moderation + stock validation in service

package cart

import (
	"net/http"

	service "tanmore_backend/internal/services/cart"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct
type GetAllCartItemsHandler struct {
	Service *service.GetAllCartItemsService
}

// 🏗️ Constructor
func NewGetAllCartItemsHandler(service *service.GetAllCartItemsService) *GetAllCartItemsHandler {
	return &GetAllCartItemsHandler{Service: service}
}

// 🔁 GET /api/cart/items
func (h *GetAllCartItemsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract customer user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Call service layer
	result, err := h.Service.Start(ctx, service.GetAllCartItemsInput{
		UserID: userID,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 3️⃣ Return grouped valid_items and flat invalid_items
	response.OK(w, "Cart items fetched successfully", map[string]interface{}{
		"valid_items":   result.ValidItems,
		"invalid_items": result.InvalidItems,
	})
}
