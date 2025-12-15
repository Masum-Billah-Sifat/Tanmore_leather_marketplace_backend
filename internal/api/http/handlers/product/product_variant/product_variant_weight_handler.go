// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/update_variant_weight_handler.go
// 🧠 Handles PUT /api/seller/products/:product_id/variants/:variant_id/weight
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Parses weight_grams from JSON body
//     - Calls service layer
//     - Returns variant_id and updated weight in grams

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
type UpdateVariantWeightHandler struct {
	Service *service.UpdateVariantWeightService
}

// 🏗️ Constructor
func NewUpdateVariantWeightHandler(service *service.UpdateVariantWeightService) *UpdateVariantWeightHandler {
	return &UpdateVariantWeightHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/variants/:variant_id/weight
func (h *UpdateVariantWeightHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Get seller user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
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
		WeightGrams int64 `json:"weight_grams"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ Validate that weight_grams is positive
	if req.WeightGrams <= 0 {
		response.BadRequest(w, errors.NewValidationError("weight_grams", "must be a positive number"))
		return
	}

	// 6️⃣ Build service input
	input := service.UpdateVariantWeightInput{
		UserID:      userID,
		ProductID:   productID,
		VariantID:   variantID,
		WeightGrams: req.WeightGrams,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Return success
	response.OK(w, "Variant weight updated successfully", map[string]interface{}{
		"variant_id":   result.VariantID,
		"weight_grams": result.WeightGrams,
		"status":       result.Status,
	})
}
