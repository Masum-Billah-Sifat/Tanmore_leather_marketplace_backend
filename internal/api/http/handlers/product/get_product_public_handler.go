// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/get_product_public_detail_handler.go
// 🧠 Handles GET /api/products/:product_id (Public Endpoint)
//     - Extracts product_id from path
//     - Calls service layer
//     - Returns product detail (with variants)

package product

import (
	"net/http"

	service "tanmore_backend/internal/services/product"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// 📦 Handler struct
type GetProductPublicDetailHandler struct {
	Service *service.GetProductPublicDetailService
}

// 🏗️ Constructor
func NewGetProductPublicDetailHandler(
	service *service.GetProductPublicDetailService,
) *GetProductPublicDetailHandler {
	return &GetProductPublicDetailHandler{Service: service}
}

// 🔁 GET /api/products/:product_id
func (h *GetProductPublicDetailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract product_id from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 2️⃣ Build input (no userID needed for public view)
	input := service.GetProductPublicDetailInput{
		ProductID: productID,
	}

	// 3️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return success response
	response.OK(w, "Product detail fetched successfully", result)
}
