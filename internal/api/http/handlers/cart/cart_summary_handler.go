// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/cart_summary_handler.go
// 🧠 Handles POST /api/cart/summary
//     - Supports both user_id and guest_user_id
//     - Parses variant_ids from JSON body
//     - Validates ownership and item availability
//     - Returns total price and any invalid items

package cart

import (
	"encoding/json"
	"net/http"

	service "tanmore_backend/internal/services/cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct
type CartSummaryHandler struct {
	Service *service.CartSummaryService
}

// 🏗️ Constructor
func NewCartSummaryHandler(service *service.CartSummaryService) *CartSummaryHandler {
	return &CartSummaryHandler{Service: service}
}

// 🔁 POST /api/cart/summary
func (h *CartSummaryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// // 1️⃣ Extract user ID from token (if present)
	// var userID *uuid.UUID
	// rawUserID := ctx.Value(token.CtxUserIDKey)
	// if rawUserID != nil {
	// 	parsed, err := uuid.Parse(rawUserID.(string))
	// 	if err != nil {
	// 		fmt.Println("❌ Failed to parse user ID:", err)
	// 		response.Unauthorized(w, err)
	// 		return
	// 	}
	// 	userID = &parsed
	// 	fmt.Println("✅ Parsed authenticated user ID:", *userID)
	// }

	// if userID == nil {
	// 	response.BadRequest(w, errors.NewValidationError("authorization", "user_id must be present in token"))
	// 	return
	// }

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

	// 4️⃣ Parse request body
	var req struct {
		VariantIDs []string `json:"variant_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON"))
		return
	}

	if len(req.VariantIDs) == 0 {
		response.BadRequest(w, errors.NewValidationError("variant_ids", "must be a non-empty array"))
		return
	}

	var variantUUIDs []uuid.UUID
	for _, idStr := range req.VariantIDs {
		parsed, err := uuid.Parse(idStr)
		if err != nil {
			response.BadRequest(w, errors.NewValidationError("variant_ids", "invalid UUID: "+idStr))
			return
		}
		variantUUIDs = append(variantUUIDs, parsed)
	}

	input := service.CartSummaryInput{
		UserID:     userID, // now non-pointer
		VariantIDs: variantUUIDs,
	}

	// 6️⃣ Call service layer
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 7️⃣ Return success
	response.OK(w, "Order summary calculated", map[string]interface{}{
		"total_price":   result.TotalPrice,
		"invalid_items": result.InvalidItems,
	})
}
