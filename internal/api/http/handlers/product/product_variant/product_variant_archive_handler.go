// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/remove_product_variant_handler.go
// 🧠 Handles DELETE /api/seller/products/:product_id/variants/:variant_id
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Calls service layer
//     - Returns product_id and variant_id

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
type RemoveProductVariantHandler struct {
	Service *service.RemoveProductVariantService
}

// 🏗️ Constructor
func NewRemoveProductVariantHandler(service *service.RemoveProductVariantService) *RemoveProductVariantHandler {
	return &RemoveProductVariantHandler{Service: service}
}

// 🔁 DELETE /api/seller/products/:product_id/variants/:variant_id
func (h *RemoveProductVariantHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
	input := service.RemoveProductVariantInput{
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
	response.OK(w, "Variant archived successfully", map[string]interface{}{
		"product_id": result.ProductID,
		"variant_id": result.VariantID,
		"status":     result.Status,
	})
}
