// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/update_variant_wholesale_discount_handler.go
// 🧠 Handles: PUT /api/seller/products/:product_id/variants/:variant_id/wholesale-discount
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
type UpdateVariantWholesaleDiscountHandler struct {
	Service *service.UpdateWholesaleDiscountService
}

// 🏗️ Constructor
func NewUpdateVariantWholesaleDiscountHandler(service *service.UpdateWholesaleDiscountService) *UpdateVariantWholesaleDiscountHandler {
	return &UpdateVariantWholesaleDiscountHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/variants/:variant_id/wholesale-discount
func (h *UpdateVariantWholesaleDiscountHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract seller user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
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
		WholesaleDiscount     *int64  `json:"wholesale_discount"`      // optional
		WholesaleDiscountType *string `json:"wholesale_discount_type"` // optional
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON"))
		return
	}

	// 5️⃣ Require at least one field
	if req.WholesaleDiscount == nil && req.WholesaleDiscountType == nil {
		response.BadRequest(w, errors.NewValidationError("wholesale_discount|wholesale_discount_type", "at least one field must be provided"))
		return
	}

	// 6️⃣ Build service input
	input := service.UpdateWholesaleDiscountInput{
		UserID:                userID,
		ProductID:             productID,
		VariantID:             variantID,
		WholesaleDiscount:     req.WholesaleDiscount,
		WholesaleDiscountType: req.WholesaleDiscountType,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Return success
	response.OK(w, "Wholesale discount updated successfully", map[string]interface{}{
		"variant_id":     result.VariantID,
		"updated_fields": result.UpdatedFields,
		"status":         result.Status,
	})
}
