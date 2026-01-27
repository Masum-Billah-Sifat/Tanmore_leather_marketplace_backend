// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/add_to_cart_handler.go
// 🧠 Handles POST /api/cart/add
//     - Extracts customer user_id from context
//     - Parses product_id, variant_id, and required_quantity from JSON body
//     - Calls service layer
//     - Returns variant_id and cart status

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
type AddToCartHandler struct {
	Service *service.AddToCartService
}

// 🏗️ Constructor
func NewAddToCartHandler(service *service.AddToCartService) *AddToCartHandler {
	return &AddToCartHandler{Service: service}
}

func (h *AddToCartHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract user ID from access token middleware (if present)
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

	// 2️⃣ Extract guest_user_id from header (if present)
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

	// 3️⃣ Enforce: exactly one of userID or guestUserID must be present
	if (userID == nil && guestUserID == nil) || (userID != nil && guestUserID != nil) {
		fmt.Println("❌ Must provide either user_id (token) or guest_user_id (header), but not both")
		response.BadRequest(w, errors.NewValidationError("authorization", "either user_id or guest_user_id must be provided, not both"))
		return
	}

	// 4️⃣ Decode request JSON body
	var req struct {
		ProductID        string `json:"product_id"`
		VariantID        string `json:"variant_id"`
		RequiredQuantity int64  `json:"required_quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("❌ JSON decode error:", err)
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}
	fmt.Println("📥 Request Body:", req)

	// 5️⃣ Validate body fields
	if req.ProductID == "" {
		response.BadRequest(w, errors.NewValidationError("product_id", "product_id is required"))
		return
	}
	if req.VariantID == "" {
		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
		return
	}
	if req.RequiredQuantity <= 0 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "quantity must be greater than 0"))
		return
	}

	// 6️⃣ Parse UUIDs
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid UUID format"))
		return
	}
	variantID, err := uuid.Parse(req.VariantID)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
		return
	}
	if req.RequiredQuantity > math.MaxInt32 || req.RequiredQuantity < math.MinInt32 {
		response.BadRequest(w, errors.NewValidationError("required_quantity", "value out of int32 range"))
		return
	}

	// 7️⃣ Build service input
	input := service.AddToCartInput{
		UserID:           userID,      // may be nil
		GuestUserID:      guestUserID, // may be nil
		ProductID:        productID,
		VariantID:        variantID,
		RequiredQuantity: int32(req.RequiredQuantity),
	}
	fmt.Println("🚀 Calling service with input:", input)

	// 8️⃣ Call service layer
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		fmt.Println("❌ Service returned error:", err)
		response.ServerError(w, err)
		return
	}

	// 9️⃣ Send response
	response.Created(w, "Cart item processed", map[string]interface{}{
		"variant_id": result.VariantID,
		"status":     result.Status,
	})
}
