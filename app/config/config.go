package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config struct to hold configuration parameters
type Server struct {
	Port         int
	Env          string
	DatabasePath string
}

// Config func to get env value from key
func Env(key string) string {
	// Get the absolute path to the .env file
	envPath, err := filepath.Abs(".env")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return ""
	}

	// Load .env file
	err = godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return ""
	}

	return os.Getenv(key)
}
