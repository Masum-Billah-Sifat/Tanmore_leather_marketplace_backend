// db.go — Uses standard sql.DB instead of pgxpool

package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"tanmore_backend/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB(cfg *config.Config) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Unable to open DB:", err)
	}

	DB.SetConnMaxLifetime(time.Minute * 5)
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(10)

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Unable to ping DB:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL successfully")
}

// // db.go — Connects to PostgreSQL using pgxpool

// package db

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"tanmore_backend/internal/config"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// var DB *pgxpool.Pool

// func ConnectDB(cfg *config.Config) {

// 	// before it was
// 	// dbURL := fmt.Sprintf(
// 	// 	"postgres://%s:%s@%s:%s/%s",
// 	// 	cfg.DBUser,
// 	// 	cfg.DBPassword,
// 	// 	cfg.DBHost,
// 	// 	cfg.DBPort,
// 	// 	cfg.DBName,
// 	// )

// 	// now it is
// 	dbURL := fmt.Sprintf(
// 		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
// 		cfg.DBUser,
// 		cfg.DBPassword,
// 		cfg.DBHost,
// 		cfg.DBPort,
// 		cfg.DBName,
// 	)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	var err error
// 	DB, err = pgxpool.New(ctx, dbURL)
// 	if err != nil {
// 		log.Fatalf("Unable to create DB pool: %v", err)
// 	}

// 	err = DB.Ping(ctx)
// 	if err != nil {
// 		log.Fatalf("Unable to connect to DB: %v", err)
// 	}

// 	fmt.Println("✅ Connected to PostgreSQL successfully")
// }
