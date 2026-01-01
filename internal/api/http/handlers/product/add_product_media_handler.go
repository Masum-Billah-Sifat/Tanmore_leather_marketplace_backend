// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/add_product_media_handler.go
// 🧠 Handles POST /api/seller/products/:product_id/media
//     - Extracts seller user_id from context
//     - Extracts product_id from path
//     - Parses and validates media_url and media_type
//     - Calls service layer
//     - Returns media_id and product_id

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
type AddProductMediaHandler struct {
	Service *service.AddProductMediaService
}

// 🏗️ Constructor
func NewAddProductMediaHandler(service *service.AddProductMediaService) *AddProductMediaHandler {
	return &AddProductMediaHandler{Service: service}
}

// 🔁 POST /api/seller/products/:product_id/media
func (h *AddProductMediaHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 3️⃣ Parse request body
	var req struct {
		MediaURL  string `json:"media_url"`
		MediaType string `json:"media_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 4️⃣ Validate required fields
	if req.MediaURL == "" {
		response.BadRequest(w, errors.NewValidationError("media_url", "media_url is required"))
		return
	}

	if req.MediaType != "image" && req.MediaType != "promo_video" {
		response.BadRequest(w, errors.NewValidationError("media_type", "must be 'image' or 'promo_video'"))
		return
	}

	// 5️⃣ Build service input
	input := service.AddProductMediaInput{
		UserID:    userID,
		ProductID: productID,
		MediaURL:  req.MediaURL,
		MediaType: req.MediaType,
	}

	// 6️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 7️⃣ Return success
	response.Created(w, "Product media added successfully", map[string]interface{}{
		"media_id":   result.MediaID,
		"product_id": result.ProductID,
		"status":     result.Status,
	})
}
