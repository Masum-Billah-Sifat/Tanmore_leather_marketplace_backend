// // minio.go — Connects to MinIO and provides upload & presigned URL functions

// package storage

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"tanmore_backend/internal/config"

// 	"github.com/minio/minio-go/v7"
// 	"github.com/minio/minio-go/v7/pkg/credentials"
// )

// var MinioClient *minio.Client

// func ConnectMinIO(cfg *config.Config) {
// 	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
// 		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
// 		Secure: false, // Local, no SSL
// 	})
// 	if err != nil {
// 		log.Fatalf("❌ Failed to connect to MinIO: %v", err)
// 	}

// 	MinioClient = client
// 	fmt.Println("✅ Connected to MinIO successfully")

// 	// Create bucket if not exists
// 	ctx := context.Background()
// 	exists, err := client.BucketExists(ctx, cfg.MinIOBucketName)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to check bucket: %v", err)
// 	}
// 	if !exists {
// 		err = client.MakeBucket(ctx, cfg.MinIOBucketName, minio.MakeBucketOptions{})
// 		if err != nil {
// 			log.Fatalf("❌ Failed to create bucket: %v", err)
// 		}
// 		fmt.Println("📦 MinIO bucket created:", cfg.MinIOBucketName)
// 	} else {
// 		fmt.Println("📦 MinIO bucket already exists:", cfg.MinIOBucketName)
// 	}
// }

// // GeneratePresignedUploadURL returns a time-limited URL for uploading
// func GeneratePresignedUploadURL(objectName string) (string, error) {
// 	// reqParams := make(map[string]string)
// 	ctx := context.Background()
// 	url, err := MinioClient.PresignedPutObject(ctx, "tanmoremedia", objectName, time.Minute*15)
// 	if err != nil {
// 		return "", err
// 	}
// 	return url.String(), nil
// }

// minio.go — Connects to MinIO (S3-compatible) and provides presigned URL + public URL builder

package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"tanmore_backend/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client
var BucketName string
var PublicBaseURL string

func ConnectMinIO(cfg *config.Config) {
	BucketName = cfg.MinIOBucketName
	PublicBaseURL = cfg.MinIOPublicBaseURL

	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure:       cfg.MinIOSecure,
		BucketLookup: minio.BucketLookupPath, // consistent URL format: /bucket/object (works well for MinIO + S3)
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to MinIO/S3: %v", err)
	}

	MinioClient = client
	fmt.Println("✅ Connected to MinIO/S3 successfully")

	// Create bucket if not exists (MinIO local). On AWS S3, this can fail if bucket already exists globally.
	// You can keep this for dev, but consider guarding it by environment.
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, BucketName)
	if err != nil {
		log.Fatalf("❌ Failed to check bucket: %v", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("❌ Failed to create bucket: %v", err)
		}
		fmt.Println("📦 Bucket created:", BucketName)
	} else {
		fmt.Println("📦 Bucket already exists:", BucketName)
	}
}

// GeneratePresignedUploadURL returns a time-limited URL for uploading (PUT)
func GeneratePresignedUploadURL(objectName string) (string, error) {
	ctx := context.Background()
	url, err := MinioClient.PresignedPutObject(ctx, BucketName, objectName, time.Minute*15)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// BuildPublicObjectURL builds a stable URL for viewing (requires bucket/object to be publicly readable)
func BuildPublicObjectURL(objectName string) string {
	base := strings.TrimRight(PublicBaseURL, "/")
	obj := strings.TrimLeft(objectName, "/")
	return fmt.Sprintf("%s/%s/%s", base, BucketName, obj)
}
