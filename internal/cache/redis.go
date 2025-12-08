// redis.go — Connects to Redis and exposes the client

package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"tanmore_backend/internal/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr, // from .env like localhost:6379
		Password: "",            // no password for now
		DB:       0,             // default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	fmt.Println("✅ Connected to Redis successfully")
}
