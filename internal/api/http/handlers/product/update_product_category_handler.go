// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/update_product_category_handler.go
// 🧠 Handles PUT /api/seller/products/:product_id/category
//     - Extracts seller user_id from context
//     - Extracts product_id from path
//     - Parses category_id from JSON body
//     - Calls service to update product's category

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
type UpdateProductCategoryHandler struct {
	Service *service.UpdateProductCategoryService
}

// 🏗️ Constructor
func NewUpdateProductCategoryHandler(service *service.UpdateProductCategoryService) *UpdateProductCategoryHandler {
	return &UpdateProductCategoryHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/category
func (h *UpdateProductCategoryHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 3️⃣ Parse category_id from JSON body
	var req struct {
		CategoryID uuid.UUID `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	if req.CategoryID == uuid.Nil {
		response.BadRequest(w, errors.NewValidationError("category_id", "category ID is required"))
		return
	}

	// 4️⃣ Build service input
	input := service.UpdateProductCategoryInput{
		UserID:     userID,
		ProductID:  productID,
		CategoryID: req.CategoryID,
	}

	// 5️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Return success
	response.OK(w, "Product category updated successfully", map[string]interface{}{
		"product_id":            result.ProductID,
		"updated_category_id":   result.UpdatedCategoryID,
		"updated_category_name": result.UpdatedCategoryName,
		"updated":               result.Updated,
	})
}
