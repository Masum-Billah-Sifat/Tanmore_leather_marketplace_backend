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
