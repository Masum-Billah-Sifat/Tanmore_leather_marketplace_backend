// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/disable_variant_wholesale_mode_handler.go
// 🧠 Handles DELETE /api/seller/products/:product_id/variants/:variant_id/wholesale-mode
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Calls service layer (no request body)
//     - Returns variant_id and wholesale disabled status

package product_variant

import (
	"net/http"

	service "tanmore_backend/internal/services/product/product_variant"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// 📦 Handler struct
type DisableVariantWholesaleModeHandler struct {
	Service *service.DisableWholesaleModeService
}

// 🛠️ Constructor
func NewDisableVariantWholesaleModeHandler(service *service.DisableWholesaleModeService) *DisableVariantWholesaleModeHandler {
	return &DisableVariantWholesaleModeHandler{Service: service}
}

// ❌ DELETE /api/seller/products/:product_id/variants/:variant_id/wholesale-mode
func (h *DisableVariantWholesaleModeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract seller user_id from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Parse product_id from URL path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 3️⃣ Parse variant_id from URL path
	variantIDParam := chi.URLParam(r, "variant_id")
	variantID, err := uuid.Parse(variantIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid variant ID"))
		return
	}

	// 4️⃣ Build service input
	input := service.DisableWholesaleModeInput{
		UserID:    userID,
		ProductID: productID,
		VariantID: variantID,
	}

	// 5️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Return success
	response.OK(w, "Wholesale mode disabled successfully", map[string]interface{}{
		"variant_id":        result.VariantID,
		"wholesale_enabled": result.WholesaleEnabled,
		"status":            result.Status,
	})
}
