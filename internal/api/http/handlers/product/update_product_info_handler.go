// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/update_product_info_handler.go
// 🧠 Handles PUT /api/seller/products/:product_id
//     - Extracts seller user_id from context
//     - Extracts product_id from path
//     - Parses optional title and description from JSON body
//     - Validates: at least one field is required
//     - Calls service layer and returns updated fields

package product

import (
	"encoding/json"
	"net/http"

	service "tanmore_backend/internal/services/product"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// 📦 Handler struct
type UpdateProductInfoHandler struct {
	Service *service.UpdateProductInfoService
}

// 🏗️ Constructor
func NewUpdateProductInfoHandler(service *service.UpdateProductInfoService) *UpdateProductInfoHandler {
	return &UpdateProductInfoHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id
func (h *UpdateProductInfoHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 3️⃣ Parse request JSON body
	var req struct {
		Title       *string `json:"title"`       // Optional
		Description *string `json:"description"` // Optional
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 4️⃣ Validate: at least one field must be present
	if req.Title == nil && req.Description == nil {
		response.BadRequest(w, errors.NewValidationError("title|description", "at least one of 'title' or 'description' must be provided"))
		return
	}

	// 5️⃣ Build service input
	input := service.UpdateProductInfoInput{
		UserID:      userID,
		ProductID:   productID,
		Title:       req.Title,
		Description: req.Description,
	}

	// 6️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 7️⃣ Return success
	response.OK(w, "Product info updated successfully", map[string]interface{}{
		"product_id":     result.ProductID,
		"updated_fields": result.UpdatedFields,
		"updated":        result.Updated,
	})
}
