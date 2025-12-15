// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/add_variant_wholesale_discount_handler.go
// 🧠 Handles POST /api/seller/products/:product_id/variants/:variant_id/wholesale-discount
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Parses wholesale_discount and wholesale_discount_type from JSON body
//     - Calls service layer
//     - Returns variant_id and discount metadata

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
type AddVariantWholesaleDiscountHandler struct {
	Service *service.AddWholesaleDiscountService
}

// 🛠️ Constructor
func NewAddVariantWholesaleDiscountHandler(service *service.AddWholesaleDiscountService) *AddVariantWholesaleDiscountHandler {
	return &AddVariantWholesaleDiscountHandler{Service: service}
}

// ➕ POST /api/seller/products/:product_id/variants/:variant_id/wholesale-discount
func (h *AddVariantWholesaleDiscountHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
		WholesaleDiscount     int64  `json:"wholesale_discount"`
		WholesaleDiscountType string `json:"wholesale_discount_type"` // "percentage" or "flat"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ Validate input
	if req.WholesaleDiscount <= 0 {
		response.BadRequest(w, errors.NewValidationError("wholesale_discount", "must be a positive number"))
		return
	}
	if req.WholesaleDiscountType != "percentage" && req.WholesaleDiscountType != "flat" {
		response.BadRequest(w, errors.NewValidationError("wholesale_discount_type", "must be 'percentage' or 'flat'"))
		return
	}

	// 6️⃣ Build service input
	input := service.AddWholesaleDiscountInput{
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
	response.OK(w, "Wholesale discount added successfully", map[string]interface{}{
		"variant_id":              result.VariantID,
		"wholesale_discount":      result.WholesaleDiscount,
		"wholesale_discount_type": result.WholesaleDiscountType,
		"status":                  result.Status,
	})
}
