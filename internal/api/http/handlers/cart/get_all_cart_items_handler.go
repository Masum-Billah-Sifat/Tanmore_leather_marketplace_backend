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

	result, err := h.Service.Start(ctx, service.GetAllCartItemsInput{
		UserID: userID,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.OK(w, "Cart items fetched successfully", map[string]any{
		"valid_items":   result.ValidItems,
		"invalid_items": result.InvalidItems,
	})
}
