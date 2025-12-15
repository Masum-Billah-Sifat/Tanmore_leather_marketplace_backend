// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/add_variant_discount_handler.go
// 🧠 Handles POST /api/seller/products/:product_id/variants/:variant_id/retail-discount
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Parses retail_discount and retail_discount_type from JSON body
//     - Calls service layer
//     - Returns variant_id and applied discount info

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
type AddVariantRetailDiscountHandler struct {
	Service *service.AddVariantRetailDiscountService
}

// 🛠️ Constructor
func NewAddVariantRetailDiscountHandler(service *service.AddVariantRetailDiscountService) *AddVariantRetailDiscountHandler {
	return &AddVariantRetailDiscountHandler{Service: service}
}

// 🔁 POST /api/seller/products/:product_id/variants/:variant_id/retail-discount
func (h *AddVariantRetailDiscountHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
		RetailDiscount     int64  `json:"retail_discount"`
		RetailDiscountType string `json:"retail_discount_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ Validate input
	if req.RetailDiscount <= 0 {
		response.BadRequest(w, errors.NewValidationError("retail_discount", "must be a positive number"))
		return
	}
	if req.RetailDiscountType != "flat" && req.RetailDiscountType != "percentage" {
		response.BadRequest(w, errors.NewValidationError("retail_discount_type", "must be 'flat' or 'percentage'"))
		return
	}

	// 6️⃣ Build service input
	input := service.AddVariantRetailDiscountInput{
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
	response.OK(w, "Retail discount applied successfully", map[string]interface{}{
		"variant_id":           result.VariantID,
		"retail_discount":      result.RetailDiscount,
		"retail_discount_type": result.RetailDiscountType,
		"status":               result.Status,
	})
}
