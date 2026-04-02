// all the env related configuration in this file

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string
	JWTSecret      string
	JWTExpiryHours int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	expiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	if err != nil {
		expiryHours = 24
	}

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBPath:         getEnv("DB_PATH", "./finance.db"),
		JWTSecret:      getEnv("JWT_SECRET", "change-with-real-jwt-secret"),
		JWTExpiryHours: expiryHours,
	}
}

// get env function

func getEnv(key, fallback string) string {

	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
