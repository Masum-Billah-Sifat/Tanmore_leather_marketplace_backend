// // config.go — Loads all environment variables from .env

// package config

// import (
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// type Config struct {
// 	AppPort         string
// 	DBHost          string
// 	DBPort          string
// 	DBUser          string
// 	DBPassword      string
// 	DBName          string
// 	RedisAddr       string
// 	MinIOEndpoint   string
// 	MinIOAccessKey  string
// 	MinIOSecretKey  string
// 	MinIOBucketName string
// }

// func LoadConfig() *Config {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	return &Config{
// 		AppPort:         os.Getenv("APP_PORT"),
// 		DBHost:          os.Getenv("DB_HOST"),
// 		DBPort:          os.Getenv("DB_PORT"),
// 		DBUser:          os.Getenv("DB_USER"),
// 		DBPassword:      os.Getenv("DB_PASSWORD"),
// 		DBName:          os.Getenv("DB_NAME"),
// 		RedisAddr:       os.Getenv("REDIS_ADDR"),
// 		MinIOEndpoint:   os.Getenv("MINIO_ENDPOINT"),
// 		MinIOAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
// 		MinIOSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
// 		MinIOBucketName: os.Getenv("MINIO_BUCKET_NAME"),
// 	}
// }

// config.go — Loads all environment variables from .env

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisAddr  string

	MinIOEndpoint      string
	MinIOPublicBaseURL string
	MinIOSecure        bool
	MinIOAccessKey     string
	MinIOSecretKey     string
	MinIOBucketName    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// MINIO_SECURE defaults to false if empty/invalid
	minioSecure, _ := strconv.ParseBool(os.Getenv("MINIO_SECURE"))

	return &Config{
		AppPort:    os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		RedisAddr:  os.Getenv("REDIS_ADDR"),

		MinIOEndpoint:      os.Getenv("MINIO_ENDPOINT"),
		MinIOPublicBaseURL: os.Getenv("MINIO_PUBLIC_BASE_URL"),
		MinIOSecure:        minioSecure,
		MinIOAccessKey:     os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey:     os.Getenv("MINIO_SECRET_KEY"),
		MinIOBucketName:    os.Getenv("MINIO_BUCKET_NAME"),
	}
}
