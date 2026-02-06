package cart

import (
	"encoding/json"
	"math"
	"net/http"

	service "tanmore_backend/internal/services/cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

type UpdateCartQuantityHandler struct {
	Service *service.UpdateCartQuantityService
}

func NewUpdateCartQuantityHandler(service *service.UpdateCartQuantityService) *UpdateCartQuantityHandler {
	return &UpdateCartQuantityHandler{Service: service}
}

func (h *UpdateCartQuantityHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		VariantID        string `json:"variant_id"`
		RequiredQuantity int64  `json:"required_quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	if req.VariantID == "" {
		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
		return
	}

	if req.RequiredQuantity <= 0 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "quantity must be >= 1"))
		return
	}

	if req.RequiredQuantity > math.MaxInt32 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "quantity too large"))
		return
	}

	variantID, err := uuid.Parse(req.VariantID)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID"))
		return
	}

	input := service.UpdateCartQuantityInput{
		UserID:           userID,
		VariantID:        variantID,
		RequiredQuantity: int32(req.RequiredQuantity),
	}

	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.OK(w, "Cart item updated", map[string]interface{}{
		"variant_id":       result.VariantID,
		"updated_quantity": result.UpdatedQuantity,
		"status":           result.Status,
	})
}
