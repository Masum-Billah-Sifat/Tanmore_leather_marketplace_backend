// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/product_variant/update_variant_retail_discount_handler.go
// 🧠 Handles: PUT /api/seller/products/:product_id/variants/:variant_id/retail-discount
//     - Parses seller token and URL path params
//     - Validates and decodes JSON request
//     - Requires at least one field present
//     - Calls service and responds with updated info

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
type UpdateVariantRetailDiscountHandler struct {
	Service *service.UpdateVariantRetailDiscountService
}

// 🏗️ Constructor
func NewUpdateVariantRetailDiscountHandler(service *service.UpdateVariantRetailDiscountService) *UpdateVariantRetailDiscountHandler {
	return &UpdateVariantRetailDiscountHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/variants/:variant_id/retail-discount
func (h *UpdateVariantRetailDiscountHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 2️⃣ Extract product_id from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 3️⃣ Extract variant_id from path
	variantIDParam := chi.URLParam(r, "variant_id")
	variantID, err := uuid.Parse(variantIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid variant ID"))
		return
	}

	// 4️⃣ Parse request body
	var req struct {
		RetailDiscount     *int64  `json:"retail_discount"`      // Optional
		RetailDiscountType *string `json:"retail_discount_type"` // Optional
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON"))
		return
	}

	// 5️⃣ At least one field must be present
	if req.RetailDiscount == nil && req.RetailDiscountType == nil {
		response.BadRequest(w, errors.NewValidationError("retail_discount|retail_discount_type", "at least one field must be provided"))
		return
	}

	// 6️⃣ Build service input
	input := service.UpdateVariantRetailDiscountInput{
		UserID:             userID,
		ProductID:          productID,
		VariantID:          variantID,
		RetailDiscount:     req.RetailDiscount,
		RetailDiscountType: req.RetailDiscountType,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Return success
	response.OK(w, "Retail discount updated successfully", map[string]interface{}{
		"variant_id":     result.VariantID,
		"updated_fields": result.UpdatedFields,
		"status":         result.Status,
	})
}
