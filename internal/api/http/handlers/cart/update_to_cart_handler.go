// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/update_cart_quantity_handler.go
// 🧠 Handles PUT /api/cart/update
//     - Extracts customer user_id from context
//     - Parses variant_id and required_quantity from JSON body
//     - Validates input
//     - Calls service layer to update quantity
//     - Returns updated quantity and status

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

// 📦 Handler struct
type UpdateCartQuantityHandler struct {
	Service *service.UpdateCartQuantityService
}

// 🏗️ Constructor
func NewUpdateCartQuantityHandler(service *service.UpdateCartQuantityService) *UpdateCartQuantityHandler {
	return &UpdateCartQuantityHandler{Service: service}
}

// 🔁 PUT /api/cart/update
func (h *UpdateCartQuantityHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract customer user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Parse request JSON body
	var req struct {
		VariantID        string `json:"variant_id"`
		RequiredQuantity int64  `json:"required_quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 3️⃣ Validate fields
	if req.VariantID == "" {
		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
		return
	}
	if req.RequiredQuantity <= 0 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "quantity must be ≥ 1"))
		return
	}

	// ✅ Ensure value fits in int32 range
	// if req.RequiredQuantity > int64(^int32(0)) {
	// 	response.BadRequest(w, errors.NewValidationError("required_quantity", "value exceeds int32 limit"))
	// 	return
	// }
	if req.RequiredQuantity > math.MaxInt32 || req.RequiredQuantity < math.MinInt32 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "value out of int32 range"))
		return
	}

	// 4️⃣ Parse UUID
	variantID, err := uuid.Parse(req.VariantID)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
		return
	}

	// 5️⃣ Build service input
	input := service.UpdateCartQuantityInput{
		UserID:           userID,
		VariantID:        variantID,
		RequiredQuantity: int32(req.RequiredQuantity),
	}

	// 6️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 7️⃣ Return success
	response.OK(w, "Cart item updated", map[string]interface{}{
		"variant_id":       result.VariantID,
		"updated_quantity": result.UpdatedQuantity,
		"status":           result.Status,
	})
}
