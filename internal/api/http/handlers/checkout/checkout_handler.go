// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/checkout_handler.go
// 🧠 Handles POST /api/checkout/initiate
//     - Extracts customer user_id from context
//     - Parses and validates source type ("product" or "cart")
//     - Routes to appropriate service logic
//     - Returns created session_id and pricing summary

package checkout

import (
	"encoding/json"
	"net/http"

	service "tanmore_backend/internal/services/checkout"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct
type CheckoutHandler struct {
	Service *service.CheckoutService
}

// 🏗️ Constructor
func NewCheckoutHandler(service *service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{Service: service}
}

// 🔁 POST /api/checkout/initiate
func (h *CheckoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract customer user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Parse request body
	var req struct {
		Source     string   `json:"source"`                // "product" or "cart"
		VariantID  string   `json:"variant_id"`            // only for "product"
		Quantity   int32    `json:"quantity"`              // only for "product"
		VariantIDs []string `json:"variant_ids,omitempty"` // only for "cart"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON"))
		return
	}

	// 3️⃣ Validate source
	if req.Source != "product" && req.Source != "cart" {
		response.BadRequest(w, errors.NewValidationError("source", "must be 'product' or 'cart'"))
		return
	}

	// 4️⃣ Route to correct logic
	switch req.Source {
	case "product":
		// Validate product source fields
		if req.VariantID == "" || req.Quantity <= 0 {
			response.BadRequest(w, errors.NewValidationError("variant_id or quantity", "required for product checkout"))
			return
		}
		variantUUID, err := uuid.Parse(req.VariantID)
		if err != nil {
			response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID"))
			return
		}

		input := service.CheckoutFromProductInput{
			UserID:    userID,
			VariantID: variantUUID,
			Quantity:  req.Quantity,
		}
		result, err := h.Service.FromProduct(ctx, input)
		if err != nil {
			response.ServerError(w, err)
			return
		}
		// response.OK(w, "Checkout session created", map[string]interface{}{
		// 	"checkout_session_id": result.CheckoutSessionID,
		// 	"total_price":         result.TotalPrice,
		// 	"total_weight_grams":  result.TotalWeightGrams,
		// })

		response.OK(w, "Checkout session created", map[string]interface{}{
			"checkout_session_id": result.CheckoutSessionID,
			"valid_items":         result.ValidItems,        // just a single one wrapped in slice
			"valid_items_grouped": result.ValidItemsGrouped, // ✅ Include this!

		})

		return

	case "cart":
		// Validate cart source fields
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

		input := service.CheckoutFromCartInput{
			UserID:     userID,
			VariantIDs: variantUUIDs,
		}
		result, err := h.Service.FromCart(ctx, input)
		if err != nil {
			response.ServerError(w, err)
			return
		}
		// response.OK(w, "Checkout session created", map[string]interface{}{
		// 	"checkout_session_id": result.CheckoutSessionID,
		// 	"total_price":         result.TotalPrice,
		// 	"total_weight_grams":  result.TotalWeightGrams,
		// })

		response.OK(w, "Checkout session created", map[string]interface{}{
			"checkout_session_id": result.CheckoutSessionID,
			"valid_items":         result.ValidItems,
			"invalid_items":       result.InvalidItems,
			"valid_items_grouped": result.ValidItemsGrouped, // ✅ Include this!

		})
		return
	}
}
