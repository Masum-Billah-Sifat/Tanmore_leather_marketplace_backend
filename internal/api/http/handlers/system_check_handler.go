package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"tanmore_backend/internal/cache"
	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/internal/storage"
)

type SystemCheckHandler struct {
	queries *sqlc.Queries
}

func NewSystemCheckHandler(q *sqlc.Queries) *SystemCheckHandler {
	return &SystemCheckHandler{queries: q}
}

func (h *SystemCheckHandler) HandleCheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// 1️⃣ Insert into Postgres
	inserted, err := h.queries.InsertSystemCheckLog(ctx, sql.NullString{String: "check_from_endpoint", Valid: true})

	// this one had errro
	// inserted, err := h.queries.InsertSystemCheckLog(ctx, "check_from_endpoint")

	if err != nil {
		http.Error(w, "Postgres insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2️⃣ Write to Redis
	err = cache.RedisClient.Set(ctx, "healthcheck:redis", "OK", 10*time.Minute).Err()
	if err != nil {
		http.Error(w, "Redis set failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3️⃣ Generate MinIO presigned URL
	url, err := storage.GeneratePresignedUploadURL("test-check.png")
	if err != nil {
		http.Error(w, "MinIO presign failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return all three results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"postgres_log_id":     inserted.ID,
		"redis_status":        "OK",
		"minio_presigned_url": url,
	})
}
