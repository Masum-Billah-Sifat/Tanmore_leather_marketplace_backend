// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/edit_review_reply_handler.go
// 🧠 Handles PUT /api/products/:product_id/reviews/:review_id/reply
//     - Extracts seller_user_id from context
//     - Extracts product_id and review_id from URL
//     - Parses reply_text from body
//     - Calls service layer
//     - Returns review_id on success

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
type EditReviewReplyHandler struct {
	Service *service.EditReviewReplyService
}

// 🏗️ Constructor
func NewEditReviewReplyHandler(service *service.EditReviewReplyService) *EditReviewReplyHandler {
	return &EditReviewReplyHandler{Service: service}
}

// 🔁 PUT /api/products/:product_id/reviews/:review_id/reply
func (h *EditReviewReplyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract seller_user_id from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	sellerUserID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Extract product_id and review_id from path
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

	// 3️⃣ Parse JSON body
	var req struct {
		ReplyText string `json:"reply_text"`
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

	// 5️⃣ Call service
	input := service.EditReviewReplyInput{
		SellerUserID: sellerUserID,
		ProductID:    productID,
		ReviewID:     reviewID,
		ReplyText:    req.ReplyText,
	}

	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Respond
	response.OK(w, "Reply updated successfully", map[string]interface{}{
		"review_id": result.ReviewID,
		"status":    result.Status,
	})
}
