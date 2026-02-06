// ------------------------------------------------------------
// 📁 internal/api/http/handlers/cart/remove_cart_item_handler.go
// 🔒 Authenticated only: DELETE /api/cart/remove/{variant_id}

package cart

import (
	"net/http"

	service "tanmore_backend/internal/services/cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RemoveCartItemHandler struct {
	Service *service.RemoveCartItemService
}

func NewRemoveCartItemHandler(service *service.RemoveCartItemService) *RemoveCartItemHandler {
	return &RemoveCartItemHandler{Service: service}
}

func (h *RemoveCartItemHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 2️⃣ Extract variant_id from URL
	variantIDStr := chi.URLParam(r, "variant_id")
	if variantIDStr == "" {
		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
		return
	}
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
		return
	}

	// 3️⃣ Call service
	result, err := h.Service.Start(ctx, service.RemoveCartItemInput{
		UserID:    userID,
		VariantID: variantID,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Success
	response.OK(w, "Cart item removed", map[string]interface{}{
		"variant_id": result.VariantID,
		"status":     result.Status,
	})
}
