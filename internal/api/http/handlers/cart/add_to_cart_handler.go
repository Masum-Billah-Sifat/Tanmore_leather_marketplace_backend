// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/add_to_cart_handler.go (REFINED)

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

type AddToCartHandler struct {
	Service *service.AddToCartService
}

func NewAddToCartHandler(service *service.AddToCartService) *AddToCartHandler {
	return &AddToCartHandler{Service: service}
}

func (h *AddToCartHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// Step 2️⃣: Parse and validate body
	var req struct {
		ProductID        string `json:"product_id"`
		VariantID        string `json:"variant_id"`
		RequiredQuantity int64  `json:"required_quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	if req.ProductID == "" || req.VariantID == "" || req.RequiredQuantity <= 0 {
		response.BadRequest(w, errors.NewValidationError("fields", "missing or invalid fields"))
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid UUID"))
		return
	}
	variantID, err := uuid.Parse(req.VariantID)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID"))
		return
	}
	if req.RequiredQuantity > math.MaxInt32 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "value too large"))
		return
	}

	// Step 3️⃣: Execute service call
	input := service.AddToCartInput{
		UserID:           userID,
		ProductID:        productID,
		VariantID:        variantID,
		RequiredQuantity: int32(req.RequiredQuantity),
	}

	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// Step 4️⃣: Return response
	response.Created(w, "Cart item processed", map[string]interface{}{
		"variant_id": result.VariantID,
		"status":     result.Status,
	})
}
