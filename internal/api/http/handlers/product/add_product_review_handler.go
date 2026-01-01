// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/add_product_review_handler.go
// 🧠 Handles POST /api/products/:product_id/reviews
//     - Extracts user_id (customer) from context
//     - Extracts product_id from path
//     - Parses review_text and optional review_image_url
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
type AddProductReviewHandler struct {
	Service *service.AddProductReviewService
}

// 🏗️ Constructor
func NewAddProductReviewHandler(service *service.AddProductReviewService) *AddProductReviewHandler {
	return &AddProductReviewHandler{Service: service}
}

// 🔁 POST /api/products/:product_id/reviews
func (h *AddProductReviewHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 3️⃣ Parse request body
	var req struct {
		ReviewText     string `json:"review_text"`
		ReviewImageURL string `json:"review_image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 4️⃣ Validate review_text
	if req.ReviewText == "" {
		response.BadRequest(w, errors.NewValidationError("review_text", "review_text is required"))
		return
	}

	// 5️⃣ Build service input
	input := service.AddProductReviewInput{
		UserID:         userID,
		ProductID:      productID,
		ReviewText:     req.ReviewText,
		ReviewImageURL: req.ReviewImageURL,
	}

	// 6️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 7️⃣ Return success
	response.Created(w, "Review submitted successfully", map[string]interface{}{
		"review_id":  result.ReviewID,
		"product_id": result.ProductID,
		"status":     result.Status,
	})
}
