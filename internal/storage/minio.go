// minio.go — Connects to MinIO and provides upload & presigned URL functions

package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"tanmore_backend/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinIO(cfg *config.Config) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: false, // Local, no SSL
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to MinIO: %v", err)
	}

	MinioClient = client
	fmt.Println("✅ Connected to MinIO successfully")

	// Create bucket if not exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinIOBucketName)
	if err != nil {
		log.Fatalf("❌ Failed to check bucket: %v", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.MinIOBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("❌ Failed to create bucket: %v", err)
		}
		fmt.Println("📦 MinIO bucket created:", cfg.MinIOBucketName)
	} else {
		fmt.Println("📦 MinIO bucket already exists:", cfg.MinIOBucketName)
	}
}

// GeneratePresignedUploadURL returns a time-limited URL for uploading
func GeneratePresignedUploadURL(objectName string) (string, error) {
	// reqParams := make(map[string]string)
	ctx := context.Background()
	url, err := MinioClient.PresignedPutObject(ctx, "tanmoremedia", objectName, time.Minute*15)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
