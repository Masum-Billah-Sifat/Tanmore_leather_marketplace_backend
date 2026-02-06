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
type ArchiveProductReviewHandler struct {
	Service *service.ArchiveProductReviewService
}

// 🏗️ Constructor
func NewArchiveProductReviewHandler(service *service.ArchiveProductReviewService) *ArchiveProductReviewHandler {
	return &ArchiveProductReviewHandler{Service: service}
}

// 🔁 PUT /api/products/:product_id/reviews/:review_id/archive
func (h *ArchiveProductReviewHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 2️⃣ Extract product ID from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 3️⃣ Extract review ID from path
	reviewIDParam := chi.URLParam(r, "review_id")
	reviewID, err := uuid.Parse(reviewIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("review_id", "invalid review ID"))
		return
	}

	// 4️⃣ Build input
	input := service.ArchiveProductReviewInput{
		UserID:    userID,
		ProductID: productID,
		ReviewID:  reviewID,
	}

	// 5️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Return success
	response.Created(w, "Review archived successfully", map[string]interface{}{
		"review_id":  result.ReviewID,
		"product_id": result.ProductID,
		"status":     result.Status,
	})
}
