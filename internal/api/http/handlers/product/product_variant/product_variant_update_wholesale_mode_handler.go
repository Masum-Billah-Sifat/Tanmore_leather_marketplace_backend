// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/product_variant/edit_variant_wholesale_info_handler.go
// 🧠 Handles: PUT /api/seller/products/:product_id/variants/:variant_id/wholesale-info
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
type EditVariantWholesaleInfoHandler struct {
	Service *service.EditWholesaleInfoService
}

// 🏗️ Constructor
func NewEditVariantWholesaleInfoHandler(service *service.EditWholesaleInfoService) *EditVariantWholesaleInfoHandler {
	return &EditVariantWholesaleInfoHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/variants/:variant_id/wholesale-info
func (h *EditVariantWholesaleInfoHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
		WholesalePrice  *int64 `json:"wholesale_price"`   // Optional
		MinQtyWholesale *int32 `json:"min_qty_wholesale"` // Optional
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON"))
		return
	}

	// 5️⃣ Require at least one field
	if req.WholesalePrice == nil && req.MinQtyWholesale == nil {
		response.BadRequest(w, errors.NewValidationError("wholesale_price|min_qty_wholesale", "at least one field must be provided"))
		return
	}

	// 6️⃣ Build service input
	input := service.EditWholesaleInfoInput{
		UserID:          userID,
		ProductID:       productID,
		VariantID:       variantID,
		WholesalePrice:  req.WholesalePrice,
		MinQtyWholesale: req.MinQtyWholesale,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Return success
	response.OK(w, "Wholesale info updated successfully", map[string]interface{}{
		"variant_id":     result.VariantID,
		"updated_fields": result.UpdatedFields,
		"status":         result.Status,
	})
}
