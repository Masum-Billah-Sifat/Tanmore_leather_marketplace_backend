// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/get_products_by_category_handler.go
// 🧠 Handles GET /api/category-products
//     - Public endpoint (no auth)
//     - Extracts category_id from query param
//     - Validates UUID format
//     - Calls service layer
//     - Returns list of products grouped by product ID

package product

import (
	"net/http"

	service "tanmore_backend/internal/services/product"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"

	"github.com/google/uuid"
)

// 📦 Handler struct
type GetProductsByCategoryHandler struct {
	Service *service.GetProductsByCategoryService
}

// 🏗️ Constructor
func NewGetProductsByCategoryHandler(service *service.GetProductsByCategoryService) *GetProductsByCategoryHandler {
	return &GetProductsByCategoryHandler{Service: service}
}

// 🔁 GET /api/category-products?category_id=...
func (h *GetProductsByCategoryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract category_id from query
	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr == "" {
		response.BadRequest(w, errors.NewValidationError("category_id", "is required"))
		return
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("category_id", "must be a valid UUID"))
		return
	}

	// 2️⃣ Build input
	input := service.GetProductsByCategoryInput{
		CategoryID: categoryID,
	}

	// 3️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return success
	response.OK(w, "Products fetched successfully", result.Data)
}
