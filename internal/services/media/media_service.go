package media

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"tanmore_backend/internal/storage"
	"tanmore_backend/pkg/errors"

	"github.com/google/uuid"
)

type PresignUploadInput struct {
	UserID    uuid.UUID
	MediaType string
	FileExt   string
}

type PresignUploadResult struct {
	UploadURL string `json:"upload_url"`
	MediaURL  string `json:"media_url"`
}

type MediaService struct{}

func NewMediaService() *MediaService {
	return &MediaService{}
}

func (s *MediaService) GeneratePresignedUploadURL(
	ctx context.Context,
	input PresignUploadInput,
) (*PresignUploadResult, error) {

	log.Println("🧠 [MEDIA] service called")
	log.Printf("🧠 [MEDIA] input: user_id=%s media_type=%s file_ext=%s\n",
		input.UserID, input.MediaType, input.FileExt)

	// 🔒 Step 1: Validate
	mediaType := strings.ToLower(input.MediaType)
	fileExt := strings.ToLower(input.FileExt)

	allowedImageExts := map[string]bool{"jpg": true, "jpeg": true, "png": true}
	allowedVideoExts := map[string]bool{"mp4": true}

	switch mediaType {
	case "image":
		if !allowedImageExts[fileExt] {
			log.Println("❌ [MEDIA] unsupported image format:", fileExt)
			return nil, errors.NewValidationError("file_extension", "unsupported image format")
		}
	case "video":
		if !allowedVideoExts[fileExt] {
			log.Println("❌ [MEDIA] unsupported video format:", fileExt)
			return nil, errors.NewValidationError("file_extension", "unsupported video format")
		}
	default:
		log.Println("❌ [MEDIA] invalid media_type:", mediaType)
		return nil, errors.NewValidationError("media_type", "must be 'image' or 'video'")
	}

	log.Println("✅ [MEDIA] service validation passed")

	// 📦 Step 2: Generate object name
	uuidPart := uuid.New()
	timestamp := time.Now().Unix()
	objectName := fmt.Sprintf("media/%s_%d.%s", uuidPart.String(), timestamp, fileExt)

	log.Println("📦 [MEDIA] object name:", objectName)

	// 🧠 Step 3: Call storage (MinIO)
	uploadURL, err := storage.GeneratePresignedUploadURL(objectName)
	if err != nil {
		log.Println("❌ [MEDIA] storage presign failed:", err)
		return nil, errors.NewServerError("failed to generate presigned URL")
	}

	log.Println("✅ [MEDIA] storage presign success")

	// 🌍 Step 4: Build public URL
	mediaURL := fmt.Sprintf("https://cdn.tanmore.com/%s", objectName)

	log.Println("🌍 [MEDIA] public media URL:", mediaURL)

	return &PresignUploadResult{
		UploadURL: uploadURL,
		MediaURL:  mediaURL,
	}, nil
}
