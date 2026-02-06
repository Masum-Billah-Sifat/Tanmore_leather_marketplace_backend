// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/remove_variant_wholesale_discount_handler.go
// 🧠 Handles DELETE /api/seller/products/:product_id/variants/:variant_id/wholesale-discount
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Calls service layer (no JSON body required)
//     - Returns variant_id and status message

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
type RemoveVariantWholesaleDiscountHandler struct {
	Service *service.RemoveVariantWholesaleDiscountService
}

// 🛠️ Constructor
func NewRemoveVariantWholesaleDiscountHandler(service *service.RemoveVariantWholesaleDiscountService) *RemoveVariantWholesaleDiscountHandler {
	return &RemoveVariantWholesaleDiscountHandler{Service: service}
}

// 🔁 DELETE /api/seller/products/:product_id/variants/:variant_id/wholesale-discount
func (h *RemoveVariantWholesaleDiscountHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 4️⃣ Build service input
	input := service.RemoveVariantWholesaleDiscountInput{
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

	// 6️⃣ Return success response
	response.OK(w, "Wholesale discount removed successfully", map[string]interface{}{
		"variant_id": result.VariantID,
		"status":     result.Status,
	})
}
