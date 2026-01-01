// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/archive_review_reply_handler.go
// 🧠 Handles PUT /api/products/:product_id/reviews/:review_id/reply/archive
//     - Extracts seller_user_id from context
//     - Extracts product_id and review_id from URL
//     - Calls service layer
//     - Returns review_id on success

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
type ArchiveReviewReplyHandler struct {
	Service *service.ArchiveReviewReplyService
}

// 🏗️ Constructor
func NewArchiveReviewReplyHandler(service *service.ArchiveReviewReplyService) *ArchiveReviewReplyHandler {
	return &ArchiveReviewReplyHandler{Service: service}
}

// 🔁 PUT /api/products/:product_id/reviews/:review_id/reply/archive
func (h *ArchiveReviewReplyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract seller_user_id from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	sellerUserID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Extract product_id and review_id from URL
	productIDParam := chi.URLParam(r, "product_id")
	reviewIDParam := chi.URLParam(r, "review_id")

	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid UUID"))
		return
	}

	reviewID, err := uuid.Parse(reviewIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("review_id", "invalid UUID"))
		return
	}

	// 3️⃣ Call service layer
	input := service.ArchiveReviewReplyInput{
		SellerUserID: sellerUserID,
		ProductID:    productID,
		ReviewID:     reviewID,
	}

	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return success response
	response.OK(w, "Reply archived successfully", map[string]interface{}{
		"review_id": result.ReviewID,
		"status":    result.Status,
	})
}
