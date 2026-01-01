// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/category/get_category_tree_handler.go
// 🧠 Handles GET /api/categories/tree
//     - No authentication required (public endpoint)
//     - No request body or params
//     - Calls service to get full category tree
//     - Returns nested category data in JSON format

package category

import (
	"net/http"

	service "tanmore_backend/internal/services/category"
	"tanmore_backend/pkg/response"
)

// 📦 Handler struct
type GetCategoryTreeHandler struct {
	Service *service.GetCategoryTreeService
}

// 🏗️ Constructor
func NewGetCategoryTreeHandler(service *service.GetCategoryTreeService) *GetCategoryTreeHandler {
	return &GetCategoryTreeHandler{Service: service}
}

// 🔁 GET /api/categories/tree
func (h *GetCategoryTreeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Call service
	result, err := h.Service.Start(ctx)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 2️⃣ Return success response
	response.OK(w, "Category tree loaded", result.Tree)
}
