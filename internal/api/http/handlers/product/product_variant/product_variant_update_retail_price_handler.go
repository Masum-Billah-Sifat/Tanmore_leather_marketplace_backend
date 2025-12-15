// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/update_variant_retail_price_handler.go
// 🧠 Handles PUT /api/seller/products/:product_id/variants/:variant_id/retail-price
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Parses retail_price from JSON body
//     - Calls service layer
//     - Returns variant_id and new retail_price

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
type UpdateVariantRetailPriceHandler struct {
	Service *service.UpdateVariantRetailPriceService
}

// 🏗️ Constructor
func NewUpdateVariantRetailPriceHandler(service *service.UpdateVariantRetailPriceService) *UpdateVariantRetailPriceHandler {
	return &UpdateVariantRetailPriceHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/variants/:variant_id/retail-price
func (h *UpdateVariantRetailPriceHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
		RetailPrice int64 `json:"retail_price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ Validate that retail_price is positive
	if req.RetailPrice <= 0 {
		response.BadRequest(w, errors.NewValidationError("retail_price", "must be a positive number"))
		return
	}

	// 6️⃣ Build service input
	input := service.UpdateVariantRetailPriceInput{
		UserID:      userID,
		ProductID:   productID,
		VariantID:   variantID,
		RetailPrice: req.RetailPrice,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Return success
	response.OK(w, "Retail price updated successfully", map[string]interface{}{
		"variant_id":   result.VariantID,
		"retail_price": result.RetailPrice,
		"status":       result.Status,
	})
}
