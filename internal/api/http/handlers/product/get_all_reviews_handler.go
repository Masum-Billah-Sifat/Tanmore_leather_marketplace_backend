// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/get_all_reviews_handler.go
// 🧠 Handles GET /api/products/:product_id/reviews
//     - Extracts product_id from path
//     - Extracts optional pagination params (page, limit)
//     - Calls service layer
//     - Returns reviews with optional replies (public endpoint)

package product

import (
	"net/http"
	"strconv"

	service "tanmore_backend/internal/services/product"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// 📦 Handler struct
type GetAllReviewsHandler struct {
	Service *service.GetAllProductReviewsService
}

// 🏗️ Constructor
func NewGetAllReviewsHandler(service *service.GetAllProductReviewsService) *GetAllReviewsHandler {
	return &GetAllReviewsHandler{Service: service}
}

// 🔁 GET /api/products/:product_id/reviews
func (h *GetAllReviewsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract product_id from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 2️⃣ Parse pagination query params (optional)
	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	// 3️⃣ Build service input
	input := service.GetAllProductReviewsInput{
		ProductID: productID,
		Page:      page,
		Limit:     limit,
	}

	// 4️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 5️⃣ Return success
	response.OK(w, "Reviews fetched successfully", result)
}
