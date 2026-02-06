// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/update_variant_in_stock_handler.go
// 🧠 Handles PUT /api/seller/products/:product_id/variants/:variant_id/in-stock
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Parses in_stock from JSON body
//     - Calls service layer
//     - Returns variant_id and new in_stock status

package product_variant

import (
	"encoding/json"
	"net/http"

	service "tanmore_backend/internal/services/product/product_variant"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// 📦 Handler struct
type UpdateVariantInStockHandler struct {
	Service *service.UpdateVariantInStockService
}

// 🏗️ Constructor
func NewUpdateVariantInStockHandler(service *service.UpdateVariantInStockService) *UpdateVariantInStockHandler {
	return &UpdateVariantInStockHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/variants/:variant_id/in-stock
func (h *UpdateVariantInStockHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 2️⃣ Get product_id from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 3️⃣ Get variant_id from path
	variantIDParam := chi.URLParam(r, "variant_id")
	variantID, err := uuid.Parse(variantIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid variant ID"))
		return
	}

	// 4️⃣ Parse request JSON body
	var req struct {
		InStock *bool `json:"in_stock"` // Must be a pointer to detect missing field
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ Validate that in_stock is provided
	if req.InStock == nil {
		response.BadRequest(w, errors.NewValidationError("in_stock", "field is required and must be true or false"))
		return
	}

	// 6️⃣ Build service input
	input := service.UpdateVariantInStockInput{
		UserID:    userID,
		ProductID: productID,
		VariantID: variantID,
		InStock:   *req.InStock,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Return success
	response.OK(w, "In-stock status updated successfully", map[string]interface{}{
		"variant_id": result.VariantID,
		"in_stock":   result.InStock,
		"status":     result.Status,
	})
}
