// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/archive_product_media_handler.go
// 🧠 Handles DELETE /api/seller/products/:product_id/media/:media_id?media_type=image|promo_video
//     - Extracts seller user_id from context
//     - Extracts product_id and media_id from path
//     - Extracts media_type from query param
//     - Validates media_type
//     - Calls service layer
//     - Returns media_id and product_id

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
type ArchiveProductMediaHandler struct {
	Service *service.ArchiveProductMediaService
}

// 🏗️ Constructor
func NewArchiveProductMediaHandler(service *service.ArchiveProductMediaService) *ArchiveProductMediaHandler {
	return &ArchiveProductMediaHandler{Service: service}
}

// 🔁 DELETE /api/seller/products/:product_id/media/:media_id?media_type=image|promo_video
func (h *ArchiveProductMediaHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Get seller user ID from context
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

	// 3️⃣ Get media_id from path
	mediaIDParam := chi.URLParam(r, "media_id")
	mediaID, err := uuid.Parse(mediaIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("media_id", "invalid media ID"))
		return
	}

	// 4️⃣ Get media_type from query param
	mediaType := r.URL.Query().Get("media_type")
	if mediaType != "image" && mediaType != "promo_video" {
		response.BadRequest(w, errors.NewValidationError("media_type", "must be 'image' or 'promo_video'"))
		return
	}

	// 5️⃣ Build service input
	input := service.ArchiveProductMediaInput{
		UserID:    userID,
		ProductID: productID,
		MediaID:   mediaID,
		MediaType: mediaType,
	}

	// 6️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 7️⃣ Return success
	response.OK(w, "Media archived successfully", map[string]interface{}{
		"media_id":   result.MediaID,
		"product_id": result.ProductID,
		"status":     result.Status,
	})
}
