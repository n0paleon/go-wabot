package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

// GetEnv returns the value of an environment variable
func GetEnv(key string) string {
    return os.Getenv(key)
}

// Inisialisasi env variable saat package di-load
func init() {
    LoadEnv()
}