// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/get_product_full_detail_handler.go
// 🧠 Handles GET /api/seller/products/:product_id
//     - Extracts product_id from path
//     - Extracts user_id from context (string → uuid)
//     - Calls service layer
//     - Returns full product detail

package product

import (
	"net/http"

	service "tanmore_backend/internal/services/product"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// 📦 Handler struct
type GetProductFullDetailHandler struct {
	Service *service.GetProductFullDetailService
}

// 🏗️ Constructor
func NewGetProductFullDetailHandler(
	service *service.GetProductFullDetailService,
) *GetProductFullDetailHandler {
	return &GetProductFullDetailHandler{Service: service}
}

// 🔁 GET /api/seller/products/:product_id
func (h *GetProductFullDetailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract product_id from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 2️⃣ Extract user_id from context (stored as string)
	userIDStr, ok := ctx.Value(token.CtxUserIDKey).(string)
	if !ok || userIDStr == "" {
		response.Unauthorized(w, errors.NewAuthError("missing or invalid access token"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(w, errors.NewAuthError("invalid user id in token"))
		return
	}

	// 3️⃣ Build service input
	input := service.GetProductFullDetailInput{
		UserID:    userID,
		ProductID: productID,
	}

	// 4️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 5️⃣ Return success
	response.OK(w, "Product detail fetched successfully", result)
}
