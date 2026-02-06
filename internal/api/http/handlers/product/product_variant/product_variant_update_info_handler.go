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
type UpdateVariantInfoHandler struct {
	Service *service.UpdateVariantInfoService
}

// 🏗️ Constructor
func NewUpdateVariantInfoHandler(service *service.UpdateVariantInfoService) *UpdateVariantInfoHandler {
	return &UpdateVariantInfoHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/variants/:variant_id/info
func (h *UpdateVariantInfoHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 4️⃣ Parse request JSON body
	var req struct {
		Color *string `json:"color"` // Optional
		Size  *string `json:"size"`  // Optional
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ At least one of color or size must be provided
	if req.Color == nil && req.Size == nil {
		response.BadRequest(w, errors.NewValidationError("color|size", "at least one of 'color' or 'size' must be provided"))
		return
	}

	// 6️⃣ Build service input
	input := service.UpdateVariantInfoInput{
		UserID:    userID,
		ProductID: productID,
		VariantID: variantID,
		Color:     req.Color,
		Size:      req.Size,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Return success
	response.OK(w, "Variant info updated successfully", map[string]interface{}{
		"variant_id":     result.VariantID,
		"updated_fields": result.UpdatedFields,
		"status":         result.Status,
	})
}
