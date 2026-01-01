// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/reply_to_product_review_handler.go
// 🧠 Handles POST /api/products/:product_id/reviews/:review_id/reply
//     - Extracts seller_id from context
//     - Extracts product_id and review_id from path
//     - Parses reply_text and optional reply_image_url
//     - Calls service layer
//     - Returns reply_id and review_id

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
type ReplyToReviewHandler struct {
	Service *service.ReplyToReviewService
}

// 🏗️ Constructor
func NewReplyToReviewHandler(service *service.ReplyToReviewService) *ReplyToReviewHandler {
	return &ReplyToReviewHandler{Service: service}
}

// 🔁 POST /api/products/:product_id/reviews/:review_id/reply
func (h *ReplyToReviewHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Get seller ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	sellerID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Extract product_id and review_id from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid UUID"))
		return
	}

	reviewIDParam := chi.URLParam(r, "review_id")
	reviewID, err := uuid.Parse(reviewIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("review_id", "invalid UUID"))
		return
	}

	// 3️⃣ Parse request body
	var req struct {
		ReplyText     string `json:"reply_text"`
		ReplyImageURL string `json:"reply_image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON"))
		return
	}

	// 4️⃣ Validate reply_text
	if req.ReplyText == "" {
		response.BadRequest(w, errors.NewValidationError("reply_text", "reply_text is required"))
		return
	}

	// 5️⃣ Build input for service
	input := service.ReplyToReviewInput{
		SellerID:      sellerID,
		ProductID:     productID,
		ReviewID:      reviewID,
		ReplyText:     req.ReplyText,
		ReplyImageURL: req.ReplyImageURL,
	}

	// 6️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 7️⃣ Return response
	response.Created(w, "Reply submitted successfully", map[string]interface{}{
		"review_id": result.ReviewID,
		"reply_id":  result.ReplyID,
		"status":    result.Status,
	})
}
