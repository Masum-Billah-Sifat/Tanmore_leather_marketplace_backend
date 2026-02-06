// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/set_primary_image_handler.go
// 🧠 Handles PUT /api/seller/products/:product_id/images/:media_id/set-primary
//     - Extracts seller user_id from context
//     - Extracts product_id and media_id from path
//     - Calls service layer to set the image as primary
//     - Returns product_id and primary_image_id in response

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
type SetPrimaryImageHandler struct {
	Service *service.SetPrimaryImageService
}

// 🏗️ Constructor
func NewSetPrimaryImageHandler(service *service.SetPrimaryImageService) *SetPrimaryImageHandler {
	return &SetPrimaryImageHandler{Service: service}
}

// 🔁 PUT /api/seller/products/:product_id/images/:media_id/set-primary
func (h *SetPrimaryImageHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 3️⃣ Get media_id from path
	mediaIDParam := chi.URLParam(r, "media_id")
	mediaID, err := uuid.Parse(mediaIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("media_id", "invalid media ID"))
		return
	}

	// 4️⃣ Build service input
	input := service.SetPrimaryImageInput{
		UserID:    userID,
		ProductID: productID,
		MediaID:   mediaID,
	}

	// 5️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Return success response
	response.OK(w, "Primary image set successfully", map[string]interface{}{
		"product_id":       result.ProductID,
		"primary_image_id": result.PrimaryImageID,
		"status":           result.Status,
	})
}
