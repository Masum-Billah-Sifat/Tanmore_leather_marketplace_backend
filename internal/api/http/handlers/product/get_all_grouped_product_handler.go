// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/get_all_products_by_seller_handler.go
// 🧠 Handles GET /api/seller/products
//     - Extracts user_id from context (string → uuid)
//     - Calls service layer
//     - Returns all products grouped by status with variant details

package product

import (
	"net/http"

	service "tanmore_backend/internal/services/product"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct
type GetAllProductsBySellerHandler struct {
	Service *service.GetAllProductsBySellerService
}

// 🏗️ Constructor
func NewGetAllProductsBySellerHandler(
	service *service.GetAllProductsBySellerService,
) *GetAllProductsBySellerHandler {
	return &GetAllProductsBySellerHandler{Service: service}
}

// 🔁 GET /api/seller/products
func (h *GetAllProductsBySellerHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract user_id from context (string)
	userIDStr, ok := ctx.Value(token.CtxUserIDKey).(string)
	if !ok || userIDStr == "" {
		response.Unauthorized(w, errors.NewAuthError("missing or invalid access token"))
		return
	}

	// 2️⃣ Parse UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(w, errors.NewAuthError("invalid user id in token"))
		return
	}

	// 3️⃣ Call service
	// result, err := h.Service.Start(ctx, userID)

	result, err := h.Service.Start(ctx, service.GetAllProductsBySellerInput{
		UserID: userID,
	})

	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return result
	response.OK(w, "Seller products fetched successfully", result)
}
