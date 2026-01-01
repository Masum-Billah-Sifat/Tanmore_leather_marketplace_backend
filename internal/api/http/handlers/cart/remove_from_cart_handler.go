// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/remove_cart_item_handler.go
// 🧠 Handles DELETE /api/cart/remove/{variant_id}

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

// 📦 Handler struct
type RemoveCartItemHandler struct {
	Service *service.RemoveCartItemService
}

// 🏗️ Constructor
func NewRemoveCartItemHandler(service *service.RemoveCartItemService) *RemoveCartItemHandler {
	return &RemoveCartItemHandler{Service: service}
}

// 🔁 DELETE /api/cart/remove/{variant_id}
func (h *RemoveCartItemHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Get customer user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Extract variant_id from path param
	variantIDStr := chi.URLParam(r, "variant_id")
	if variantIDStr == "" {
		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
		return
	}

	// 3️⃣ Parse UUID
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
		return
	}

	// 4️⃣ Build service input
	input := service.RemoveCartItemInput{
		UserID:    userID,
		VariantID: variantID,
	}

	// 5️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		// Optional: smarter error handling
		// if errors.IsNotFound(err) {
		// 	response.NotFound(w, err)
		// 	return
		// }
		// if errors.IsValidation(err) || errors.IsConflict(err) {
		// 	response.Conflict(w, err)
		// 	return
		// }
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Return success
	response.OK(w, "Cart item removed", map[string]interface{}{
		"variant_id": result.VariantID,
		"status":     result.Status,
	})
}

// // ------------------------------------------------------------
// // 📁 File: internal/api/http/handlers/cart/remove_cart_item_handler.go
// // 🧠 Handles DELETE /api/cart/remove
// //     - Extracts customer user_id from context
// //     - Parses variant_id from JSON body
// //     - Calls service layer
// //     - Returns variant_id and cart status

// package cart

// import (
// 	"encoding/json"
// 	"net/http"

// 	service "tanmore_backend/internal/services/cart"
// 	"tanmore_backend/pkg/errors"
// 	"tanmore_backend/pkg/response"
// 	"tanmore_backend/pkg/token"

// 	"github.com/google/uuid"
// )

// // 📦 Handler struct
// type RemoveCartItemHandler struct {
// 	Service *service.RemoveCartItemService
// }

// // 🏗️ Constructor
// func NewRemoveCartItemHandler(service *service.RemoveCartItemService) *RemoveCartItemHandler {
// 	return &RemoveCartItemHandler{Service: service}
// }

// // 🔁 DELETE /api/cart/remove
// func (h *RemoveCartItemHandler) Handle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// 1️⃣ Get customer user ID from context
// 	rawUserID := ctx.Value(token.CtxUserIDKey)
// 	userID, err := uuid.Parse(rawUserID.(string))
// 	if err != nil {
// 		response.Unauthorized(w, err)
// 		return
// 	}

// 	// 2️⃣ Parse request JSON body
// 	var req struct {
// 		VariantID string `json:"variant_id"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
// 		return
// 	}

// 	// 3️⃣ Validate input
// 	if req.VariantID == "" {
// 		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
// 		return
// 	}

// 	// 4️⃣ Parse UUID
// 	variantID, err := uuid.Parse(req.VariantID)
// 	if err != nil {
// 		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
// 		return
// 	}

// 	// 5️⃣ Build service input
// 	input := service.RemoveCartItemInput{
// 		UserID:    userID,
// 		VariantID: variantID,
// 	}

// 	// 6️⃣ Call service
// 	result, err := h.Service.Start(ctx, input)
// 	if err != nil {
// 		response.ServerError(w, err)
// 		return
// 	}

// 	// 7️⃣ Return success
// 	response.OK(w, "Cart item removed", map[string]interface{}{
// 		"variant_id": result.VariantID,
// 		"status":     result.Status,
// 	})
// }
