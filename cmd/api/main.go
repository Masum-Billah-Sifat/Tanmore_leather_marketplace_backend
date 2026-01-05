package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api_http "tanmore_backend/internal/api/http"
	"tanmore_backend/internal/cache"
	"tanmore_backend/internal/config"
	"tanmore_backend/internal/db"
	"tanmore_backend/internal/jobs/engine"
	"tanmore_backend/internal/storage" // ✅ ADD THIS
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// ✅ Connect to MinIO
	storage.ConnectMinIO(cfg)

	// Connect to Postgres
	db.ConnectDB(cfg)

	// Connect to Redis
	cache.ConnectRedis(cfg)

	// Start background job processor
	ctx, cancel := context.WithCancel(context.Background())
	processor := engine.NewProcessorEngine(db.Queries)
	go func() {
		if err := processor.Start(ctx); err != nil {
			log.Println("⚠️ Processor stopped:", err)
		}
	}()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("📦 Shutting down gracefully...")
		cancel()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	// Create router
	r := api_http.NewRouter(db.DB, cache.RedisClient)

	// Start HTTP server
	fmt.Println("✅ Server starting on port", cfg.AppPort)
	err := http.ListenAndServe(":"+cfg.AppPort, r)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}

// // my updated main.go with background jobs attached to it
// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	api_http "tanmore_backend/internal/api/http" // 🟢 alias to avoid stdlib conflict
// 	"tanmore_backend/internal/cache"
// 	"tanmore_backend/internal/config"
// 	"tanmore_backend/internal/db"
// 	"tanmore_backend/internal/jobs/engine"
// 	"tanmore_backend/internal/storage"

// )

// func main() {
// 	// Load config
// 	cfg := config.LoadConfig()

// 	// Connect to Postgres
// 	db.ConnectDB(cfg)

// 	// Connect to Redis
// 	cache.ConnectRedis(cfg)

// 	// ✅ Start background job processor
// 	ctx, cancel := context.WithCancel(context.Background())
// 	processor := engine.NewProcessorEngine(db.Queries) // 🟢 use db.Queries directly
// 	go func() {
// 		if err := processor.Start(ctx); err != nil {
// 			log.Println("⚠️ Processor stopped:", err)
// 		}
// 	}()

// 	// ✅ Graceful shutdown
// 	sigs := make(chan os.Signal, 1)
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
// 	go func() {
// 		<-sigs
// 		log.Println("📦 Shutting down gracefully...")
// 		cancel()
// 		time.Sleep(1 * time.Second)
// 		os.Exit(0)
// 	}()

// 	// ✅ Create router from your API package
// 	r := api_http.NewRouter(db.DB, cache.RedisClient)

// 	// ✅ Start HTTP server
// 	fmt.Println("✅ Server starting on port", cfg.AppPort)
// 	err := http.ListenAndServe(":"+cfg.AppPort, r)
// 	if err != nil {
// 		log.Fatal("Server failed:", err)
// 	}
// }

// its the initial one where we actually tried to work with the user creation endpoint remember this is when we started

// package main

// import (
// 	"fmt"
// 	"log"
// 	net_http "net/http"

// 	"tanmore_backend/internal/api/http" // ✅ this is your centralized router
// 	"tanmore_backend/internal/cache"
// 	"tanmore_backend/internal/config"
// 	"tanmore_backend/internal/db"
// )

// func main() {
// 	// Load config
// 	cfg := config.LoadConfig()

// 	// Connect to Postgres
// 	db.ConnectDB(cfg)

// 	// Connect to Redis
// 	cache.ConnectRedis(cfg)

// 	// Connect to MinIO (for future use)
// 	// storage.ConnectMinIO(cfg) // optional unless used in router

// 	// ✅ Create chi.Router with full wiring
// 	r := http.NewRouter(db.DB, cache.RedisClient)

// 	// ✅ Start server
// 	fmt.Println("✅ Server starting on port", cfg.AppPort)
// 	err := net_http.ListenAndServe(":"+cfg.AppPort, r)
// 	if err != nil {
// 		log.Fatal("Server failed:", err)
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"tanmore_backend/internal/config"
// 	"tanmore_backend/internal/db"

// 	"github.com/go-chi/chi/v5"

// 	// update with new imports
// 	"tanmore_backend/internal/api/http/handlers"
// 	"tanmore_backend/internal/db/sqlc"
// 	service "tanmore_backend/internal/services"

// 	// add more imports

// 	"tanmore_backend/internal/cache"
// 	"tanmore_backend/internal/storage"
// 	// add more import
// 	// "tanmore_backend/internal/api/http/handlers"
// )

// func main() {
// 	// Load config
// 	cfg := config.LoadConfig()

// 	// Connect to DB
// 	db.ConnectDB(cfg)

// 	cache.ConnectRedis(cfg)
// 	storage.ConnectMinIO(cfg)

// 	queries := sqlc.New(db.DB)
// 	userService := service.NewUserService(queries)

// 	// Set up router
// 	r := chi.NewRouter()

// 	systemHandler := handlers.NewSystemCheckHandler(queries)
// 	r.Get("/debug/system-check", systemHandler.HandleCheck)

// 	r.Post("/users", handlers.CreateUserHandler(userService))

// 	// Health check
// 	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Server is running"))
// 	})

// 	// Start server
// 	fmt.Println("✅ Server starting on port", cfg.AppPort)
// 	err := http.ListenAndServe(":"+cfg.AppPort, r)
// 	if err != nil {
// 		log.Fatal("Server failed:", err)
// 	}
// }

// // // main.go — Entry point for starting the HTTP server

// // package main

// // import (
// // 	"fmt"
// // 	"log"
// // 	"net/http"
// // 	"os"

// // 	"github.com/go-chi/chi/v5"
// // 	"github.com/joho/godotenv"
// // )

// // func main() {
// // 	// Load environment variables from .env
// // 	err := godotenv.Load()
// // 	if err != nil {
// // 		log.Fatal("Error loading .env file")
// // 	}

// // 	// Read port from env
// // 	port := os.Getenv("APP_PORT")
// // 	if port == "" {
// // 		log.Fatal("APP_PORT not set in .env")
// // 	}

// // 	// Set up HTTP router
// // 	router := chi.NewRouter()

// // 	// Health check route
// // 	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
// // 		w.Write([]byte("Server is running"))
// // 	})

// // 	// Start server
// // 	fmt.Println("✅ Server starting on port", port)
// // 	err = http.ListenAndServe(":"+port, router)
// // 	if err != nil {
// // 		log.Fatal("Server failed to start:", err)
// // 	}
// // }
