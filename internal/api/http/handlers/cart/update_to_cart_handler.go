// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/update_cart_quantity_handler.go
// 🧠 Handles PUT /api/cart/update
//     - Supports either user_id (from access token) or guest_user_id (from header)
//     - Parses variant_id and required_quantity from JSON body
//     - Validates input
//     - Calls service layer
//     - Returns updated quantity and status

package cart

import (
	"encoding/json"
	"fmt"
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

	// 1️⃣ Extract user ID from context (access token)
	var userID *uuid.UUID
	rawUserID := ctx.Value(token.CtxUserIDKey)
	if rawUserID != nil {
		parsed, err := uuid.Parse(rawUserID.(string))
		if err != nil {
			fmt.Println("❌ Failed to parse user ID:", err)
			response.Unauthorized(w, err)
			return
		}
		userID = &parsed
		fmt.Println("✅ Parsed authenticated user ID:", *userID)
	}

	// 2️⃣ Extract guest_user_id from header (optional)
	var guestUserID *uuid.UUID
	guestHeader := r.Header.Get("X-Tanmore-Guest-Id")
	if guestHeader != "" {
		parsed, err := uuid.Parse(guestHeader)
		if err != nil {
			fmt.Println("❌ Invalid guest_user_id UUID:", err)
			response.BadRequest(w, errors.NewValidationError("guest_user_id", "invalid UUID format in header"))
			return
		}
		guestUserID = &parsed
		fmt.Println("✅ Parsed guest user ID from header:", *guestUserID)
	}

	// 3️⃣ Enforce: exactly one of user_id or guest_user_id must be present
	if (userID == nil && guestUserID == nil) || (userID != nil && guestUserID != nil) {
		fmt.Println("❌ Must provide either user_id (token) or guest_user_id (header), but not both")
		response.BadRequest(w, errors.NewValidationError("authorization", "either user_id or guest_user_id must be provided, not both"))
		return
	}

	// 4️⃣ Parse request JSON body
	var req struct {
		VariantID        string `json:"variant_id"`
		RequiredQuantity int64  `json:"required_quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ Validate fields
	if req.VariantID == "" {
		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
		return
	}
	if req.RequiredQuantity <= 0 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "quantity must be ≥ 1"))
		return
	}
	if req.RequiredQuantity > math.MaxInt32 || req.RequiredQuantity < math.MinInt32 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "value out of int32 range"))
		return
	}

	// 6️⃣ Parse variant UUID
	variantID, err := uuid.Parse(req.VariantID)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
		return
	}

	// 7️⃣ Build service input
	input := service.UpdateCartQuantityInput{
		UserID:           userID,
		GuestUserID:      guestUserID,
		VariantID:        variantID,
		RequiredQuantity: int32(req.RequiredQuantity),
	}

	// 8️⃣ Call service layer
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 9️⃣ Return success
	response.OK(w, "Cart item updated", map[string]interface{}{
		"variant_id":       result.VariantID,
		"updated_quantity": result.UpdatedQuantity,
		"status":           result.Status,
	})
}

// // ------------------------------------------------------------
// // 📁 File: internal/api/http/handlers/cart/update_cart_quantity_handler.go
// // 🧠 Handles PUT /api/cart/update
// //     - Extracts customer user_id from context
// //     - Parses variant_id and required_quantity from JSON body
// //     - Validates input
// //     - Calls service layer to update quantity
// //     - Returns updated quantity and status

// package cart

// import (
// 	"encoding/json"
// 	"math"
// 	"net/http"

// 	service "tanmore_backend/internal/services/cart"
// 	"tanmore_backend/pkg/errors"
// 	"tanmore_backend/pkg/response"
// 	"tanmore_backend/pkg/token"

// 	"github.com/google/uuid"
// )

// // 📦 Handler struct
// type UpdateCartQuantityHandler struct {
// 	Service *service.UpdateCartQuantityService
// }

// // 🏗️ Constructor
// func NewUpdateCartQuantityHandler(service *service.UpdateCartQuantityService) *UpdateCartQuantityHandler {
// 	return &UpdateCartQuantityHandler{Service: service}
// }

// // 🔁 PUT /api/cart/update
// func (h *UpdateCartQuantityHandler) Handle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// 1️⃣ Extract customer user ID from context
// 	rawUserID := ctx.Value(token.CtxUserIDKey)
// 	userID, err := uuid.Parse(rawUserID.(string))
// 	if err != nil {
// 		response.Unauthorized(w, err)
// 		return
// 	}

// 	// 2️⃣ Parse request JSON body
// 	var req struct {
// 		VariantID        string `json:"variant_id"`
// 		RequiredQuantity int64  `json:"required_quantity"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
// 		return
// 	}

// 	// 3️⃣ Validate fields
// 	if req.VariantID == "" {
// 		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
// 		return
// 	}
// 	if req.RequiredQuantity <= 0 {
// 		response.BadRequest(w, errors.NewValidationError("required_quantity", "quantity must be ≥ 1"))
// 		return
// 	}

// 	// ✅ Ensure value fits in int32 range
// 	// if req.RequiredQuantity > int64(^int32(0)) {
// 	// 	response.BadRequest(w, errors.NewValidationError("required_quantity", "value exceeds int32 limit"))
// 	// 	return
// 	// }
// 	if req.RequiredQuantity > math.MaxInt32 || req.RequiredQuantity < math.MinInt32 {
// 		response.BadRequest(w, errors.NewValidationError("required_quantity", "value out of int32 range"))
// 		return
// 	}

// 	// 4️⃣ Parse UUID
// 	variantID, err := uuid.Parse(req.VariantID)
// 	if err != nil {
// 		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
// 		return
// 	}

// 	// 5️⃣ Build service input
// 	input := service.UpdateCartQuantityInput{
// 		UserID:           userID,
// 		VariantID:        variantID,
// 		RequiredQuantity: int32(req.RequiredQuantity),
// 	}

// 	// 6️⃣ Call service
// 	result, err := h.Service.Start(ctx, input)
// 	if err != nil {
// 		response.ServerError(w, err)
// 		return
// 	}

// 	// 7️⃣ Return success
// 	response.OK(w, "Cart item updated", map[string]interface{}{
// 		"variant_id":       result.VariantID,
// 		"updated_quantity": result.UpdatedQuantity,
// 		"status":           result.Status,
// 	})
// }
