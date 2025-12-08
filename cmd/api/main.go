package main

import (
	"fmt"
	"log"
	"net/http"

	"tanmore_backend/internal/config"
	"tanmore_backend/internal/db"

	"github.com/go-chi/chi/v5"

	// update with new imports
	"tanmore_backend/internal/api/http/handlers"
	"tanmore_backend/internal/db/sqlc"
	service "tanmore_backend/internal/services"

	// add more imports

	"tanmore_backend/internal/cache"
	"tanmore_backend/internal/storage"
	// add more import
	// "tanmore_backend/internal/api/http/handlers"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to DB
	db.ConnectDB(cfg)

	cache.ConnectRedis(cfg)
	storage.ConnectMinIO(cfg)

	queries := sqlc.New(db.DB)
	userService := service.NewUserService(queries)

	// Set up router
	r := chi.NewRouter()

	systemHandler := handlers.NewSystemCheckHandler(queries)
	r.Get("/debug/system-check", systemHandler.HandleCheck)

	r.Post("/users", handlers.CreateUserHandler(userService))

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	// Start server
	fmt.Println("✅ Server starting on port", cfg.AppPort)
	err := http.ListenAndServe(":"+cfg.AppPort, r)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}

// // main.go — Entry point for starting the HTTP server

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	// Load environment variables from .env
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Read port from env
// 	port := os.Getenv("APP_PORT")
// 	if port == "" {
// 		log.Fatal("APP_PORT not set in .env")
// 	}

// 	// Set up HTTP router
// 	router := chi.NewRouter()

// 	// Health check route
// 	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Server is running"))
// 	})

// 	// Start server
// 	fmt.Println("✅ Server starting on port", port)
// 	err = http.ListenAndServe(":"+port, router)
// 	if err != nil {
// 		log.Fatal("Server failed to start:", err)
// 	}
// }
