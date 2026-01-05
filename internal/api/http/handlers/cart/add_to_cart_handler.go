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

	// 1️⃣ Extract user ID from context (set by access token middleware)
	rawUserID := ctx.Value(token.CtxUserIDKey)
	fmt.Println("👀 Raw user ID from context:", rawUserID)

	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		fmt.Println("❌ Failed to parse user ID:", err)
		response.Unauthorized(w, err)
		return
	}
	fmt.Println("✅ Parsed user ID:", userID)

	// 2️⃣ Decode request JSON body
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

	// 3️⃣ Validate body fields
	if req.ProductID == "" {
		fmt.Println("❌ Missing product_id")
		response.BadRequest(w, errors.NewValidationError("product_id", "product_id is required"))
		return
	}
	if req.VariantID == "" {
		fmt.Println("❌ Missing variant_id")
		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
		return
	}
	if req.RequiredQuantity <= 0 {
		fmt.Println("❌ Invalid quantity:", req.RequiredQuantity)
		response.BadRequest(w, errors.NewValidationError("required_quantity", "quantity must be greater than 0"))
		return
	}

	// 4️⃣ Parse UUIDs
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		fmt.Println("❌ Invalid product_id UUID:", err)
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid UUID format"))
		return
	}
	variantID, err := uuid.Parse(req.VariantID)
	if err != nil {
		fmt.Println("❌ Invalid variant_id UUID:", err)
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
		return
	}

	// 5️⃣ Validate quantity fits in int32 range
	if req.RequiredQuantity > math.MaxInt32 || req.RequiredQuantity < math.MinInt32 {
		fmt.Println("❌ Quantity out of int32 range:", req.RequiredQuantity)
		response.BadRequest(w, errors.NewValidationError("required_quantity", "value out of int32 range"))
		return
	}

	// 6️⃣ Build service input
	input := service.AddToCartInput{
		UserID:           userID,
		ProductID:        productID,
		VariantID:        variantID,
		RequiredQuantity: int32(req.RequiredQuantity),
	}

	fmt.Println("🚀 Calling service with input:", input)

	// 7️⃣ Call service layer
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		fmt.Println("❌ Service returned error:", err)
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Send response
	fmt.Println("✅ Cart item processed:", result)
	response.Created(w, "Cart item processed", map[string]interface{}{
		"variant_id": result.VariantID,
		"status":     result.Status,
	})
}

// // 🔁 POST /api/cart/add
// func (h *AddToCartHandler) Handle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// 1️⃣ Get customer user ID from context
// 	rawUserID := ctx.Value(token.CtxUserIDKey)
// 	fmt.Println("👀 Raw user ID from context:", rawUserID)

// 	userID, err := uuid.Parse(rawUserID.(string))
// 	if err != nil {
// 		fmt.Println("❌ Failed to parse user ID:", err)
// 		response.Unauthorized(w, err)
// 		return
// 	}
// 	fmt.Println("✅ Parsed user ID:", userID)

// 	// 2️⃣ Parse request JSON body
// 	var req struct {
// 		ProductID        string `json:"product_id"`
// 		VariantID        string `json:"variant_id"`
// 		RequiredQuantity int64  `json:"required_quantity"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
// 		return
// 	}

// 	// 3️⃣ Validate fields
// 	if req.ProductID == "" {
// 		response.BadRequest(w, errors.NewValidationError("product_id", "product_id is required"))
// 		return
// 	}
// 	if req.VariantID == "" {
// 		response.BadRequest(w, errors.NewValidationError("variant_id", "variant_id is required"))
// 		return
// 	}
// 	if req.RequiredQuantity <= 0 {
// 		response.BadRequest(w, errors.NewValidationError("required_quantity", "quantity must be greater than 0"))
// 		return
// 	}

// 	// 4️⃣ Parse UUIDs
// 	productID, err := uuid.Parse(req.ProductID)
// 	if err != nil {
// 		response.BadRequest(w, errors.NewValidationError("product_id", "invalid UUID format"))
// 		return
// 	}
// 	variantID, err := uuid.Parse(req.VariantID)
// 	if err != nil {
// 		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid UUID format"))
// 		return
// 	}

// 	// ✅ Validate quantity range before casting
// 	// if req.RequiredQuantity > int64(^int32(0)) || req.RequiredQuantity < int64(-1<<31) {
// 	if req.RequiredQuantity > math.MaxInt32 || req.RequiredQuantity < math.MinInt32 {

// 		response.BadRequest(w, errors.NewValidationError("required_quantity", "value out of int32 range"))
// 		return
// 	}

// 	// 5️⃣ Build service input
// 	input := service.AddToCartInput{
// 		UserID:           userID,
// 		ProductID:        productID,
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
// 	response.Created(w, "Cart item processed", map[string]interface{}{
// 		"variant_id": result.VariantID,
// 		"status":     result.Status,
// 	})
// }
