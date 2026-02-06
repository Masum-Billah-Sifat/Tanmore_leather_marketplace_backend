// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/archive_product_handler.go
// 🧠 Handles PUT /api/seller/products/:product_id/archive
//     - Extracts seller user_id from context
//     - Extracts product_id from path
//     - Calls service layer
//     - Returns product_id and archive status

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
type ArchiveProductHandler struct {
	Service *service.ArchiveProductService
}

// 🏗️ Constructor
func NewArchiveProductHandler(service *service.ArchiveProductService) *ArchiveProductHandler {
	return &ArchiveProductHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/archive
func (h *ArchiveProductHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 3️⃣ Build input
	input := service.ArchiveProductInput{
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
	response.OK(w, "Product archived successfully", map[string]interface{}{
		"product_id": result.ProductID,
		"status":     result.Status,
	})
}
