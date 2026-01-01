// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/edit_product_review_handler.go
// 🧠 Handles PUT /api/products/:product_id/reviews/:review_id
//     - Extracts user_id (customer) from context
//     - Extracts product_id and review_id from path
//     - Parses review_text
//     - Calls service layer
//     - Returns review_id and product_id

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
type EditProductReviewHandler struct {
	Service *service.EditProductReviewService
}

// 🏗️ Constructor
func NewEditProductReviewHandler(service *service.EditProductReviewService) *EditProductReviewHandler {
	return &EditProductReviewHandler{Service: service}
}

// 🔁 PUT /api/products/:product_id/reviews/:review_id
func (h *EditProductReviewHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Get user ID from context
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

	// 3️⃣ Get review_id from path
	reviewIDParam := chi.URLParam(r, "review_id")
	reviewID, err := uuid.Parse(reviewIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("review_id", "invalid review ID"))
		return
	}

	// 4️⃣ Parse request body
	var req struct {
		ReviewText string `json:"review_text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ Validate review_text
	if req.ReviewText == "" {
		response.BadRequest(w, errors.NewValidationError("review_text", "review_text is required"))
		return
	}

	// 6️⃣ Build service input
	input := service.EditProductReviewInput{
		UserID:     userID,
		ProductID:  productID,
		ReviewID:   reviewID,
		ReviewText: req.ReviewText,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 8️⃣ Return success
	response.Created(w, "Review updated successfully", map[string]interface{}{
		"review_id":  result.ReviewID,
		"product_id": result.ProductID,
		"status":     result.Status,
	})
}
