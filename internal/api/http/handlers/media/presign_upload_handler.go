package media

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	mediasvc "tanmore_backend/internal/services/media"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

type Handler struct {
	Service *mediasvc.MediaService
}

func NewHandler(service *mediasvc.MediaService) *Handler {
	return &Handler{Service: service}
}

type PresignUploadRequest struct {
	MediaType     string `json:"media_type"`
	FileExtension string `json:"file_extension"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("📥 [MEDIA] Presign upload handler hit")

	ctx := r.Context()

	// 🔐 Step 1: Extract user_id
	userIDStr := ctx.Value(token.CtxUserIDKey)
	if userIDStr == nil {
		log.Println("❌ [MEDIA] user_id missing in context")
		response.Unauthorized(w, errors.NewAuthError("missing user_id in token"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		log.Println("❌ [MEDIA] invalid user_id format:", err)
		response.Unauthorized(w, errors.NewAuthError("invalid user_id format"))
		return
	}

	log.Println("✅ [MEDIA] user_id:", userID.String())

	// 📥 Step 2: Parse request body
	var body PresignUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("❌ [MEDIA] failed to decode request body:", err)
		response.BadRequest(w, errors.NewServerError("invalid JSON body"))
		return
	}

	log.Println("📦 [MEDIA] request body:", body)

	mediaType := strings.ToLower(body.MediaType)
	fileExt := strings.ToLower(body.FileExtension)

	// 🔍 Step 3: Validate inputs
	allowedImageExts := map[string]bool{"jpg": true, "jpeg": true, "png": true}
	allowedVideoExts := map[string]bool{"mp4": true}

	switch mediaType {
	case "image":
		if !allowedImageExts[fileExt] {
			log.Println("❌ [MEDIA] invalid image extension:", fileExt)
			response.BadRequest(w, errors.NewValidationError("file_extension", "only jpg, jpeg, or png allowed for images"))
			return
		}
	case "video":
		if !allowedVideoExts[fileExt] {
			log.Println("❌ [MEDIA] invalid video extension:", fileExt)
			response.BadRequest(w, errors.NewValidationError("file_extension", "only mp4 allowed for videos"))
			return
		}
	default:
		log.Println("❌ [MEDIA] invalid media_type:", mediaType)
		response.BadRequest(w, errors.NewValidationError("media_type", "must be 'image' or 'video'"))
		return
	}

	log.Println("✅ [MEDIA] validation passed")

	// 🚀 Step 4: Call service
	result, err := h.Service.GeneratePresignedUploadURL(ctx, mediasvc.PresignUploadInput{
		UserID:    userID,
		MediaType: mediaType,
		FileExt:   fileExt,
	})

	if err != nil {
		log.Println("❌ [MEDIA] service error:", err)
		response.ServerError(w, err)
		return
	}

	log.Println("✅ [MEDIA] presigned URL generated successfully")

	response.OK(w, "Presigned upload URL generated", result)
}

// // ------------------------------------------------------------
// // 📁 File: internal/api/http/handler/media/presign_upload_handler.go
// // 🧠 Handles POST /api/media/presign-upload
// //     - Validates file type based on media type
// //     - Enforces restrictions for product image and video uploads

// package media

// import (
// 	"encoding/json"
// 	"net/http"
// 	"strings"

// 	mediasvc "tanmore_backend/internal/services/media"
// 	"tanmore_backend/pkg/errors"
// 	"tanmore_backend/pkg/response"
// 	"tanmore_backend/pkg/token"

// 	"github.com/google/uuid"
// )

// // 📦 Handler struct holds dependencies
// type Handler struct {
// 	Service *mediasvc.MediaService
// }

// // 🛠️ Constructor
// func NewHandler(service *mediasvc.MediaService) *Handler {
// 	return &Handler{Service: service}
// }

// // 📩 Request body structure
// type PresignUploadRequest struct {
// 	MediaType     string `json:"media_type"`     // "image" or "video"
// 	FileExtension string `json:"file_extension"` // e.g. "jpg", "mp4"
// }

// // 🚀 POST /api/media/presign-upload
// func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// 🔐 Extract user_id from context
// 	userIDStr := ctx.Value(token.CtxUserIDKey)
// 	if userIDStr == nil {
// 		response.Unauthorized(w, errors.NewAuthError("missing user_id in token"))
// 		return
// 	}

// 	userID, err := uuid.Parse(userIDStr.(string))
// 	if err != nil {
// 		response.Unauthorized(w, errors.NewAuthError("invalid user_id format"))
// 		return
// 	}

// 	// 📥 Parse request body
// 	var body PresignUploadRequest
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		response.BadRequest(w, errors.NewServerError("invalid JSON body"))
// 		return
// 	}

// 	mediaType := strings.ToLower(body.MediaType)
// 	fileExt := strings.ToLower(body.FileExtension)

// 	// ✅ Allowed formats based on type
// 	allowedImageExts := map[string]bool{"jpg": true, "jpeg": true, "png": true}
// 	allowedVideoExts := map[string]bool{"mp4": true}

// 	// 🔍 Validate media type and extension
// 	switch mediaType {
// 	case "image":
// 		if !allowedImageExts[fileExt] {
// 			response.BadRequest(w, errors.NewValidationError("file_extension", "only jpg, jpeg, or png allowed for images"))
// 			return
// 		}
// 	case "video":
// 		if !allowedVideoExts[fileExt] {
// 			response.BadRequest(w, errors.NewValidationError("file_extension", "only mp4 allowed for videos"))
// 			return
// 		}
// 	default:
// 		response.BadRequest(w, errors.NewValidationError("media_type", "must be 'image' or 'video'"))
// 		return
// 	}

// 	// 🚀 Call service layer
// 	result, err := h.Service.GeneratePresignedUploadURL(ctx, mediasvc.PresignUploadInput{
// 		UserID:    userID,
// 		MediaType: mediaType,
// 		FileExt:   fileExt,
// 	})
// 	if err != nil {
// 		response.ServerError(w, err)
// 		return
// 	}

// 	// ✅ Return presigned + media URL
// 	response.OK(w, "Presigned upload URL generated", result)
// }
