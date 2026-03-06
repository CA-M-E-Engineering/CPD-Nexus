package config

import (
	"fmt"
	"os"
	"sgbuildex/internal/pkg/logger"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPass     string
	DBHost     string
	DBName     string
	DBDSN      string
	APIPort    string
	IngressURL string
	PitstopURL string

	JWTSecret      string
	AllowedOrigins string

	WorkerIntervalMinutes int
}

func LoadConfig() *Config {
	// Try loading .env from common locations
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	cfg := &Config{
		DBUser:         getEnvRequired("DB_USER"),
		DBPass:         getEnvRequired("DB_PASS"),
		DBHost:         getEnv("DB_HOST", "127.0.0.1:3306"),
		DBName:         getEnv("DB_NAME", "bas_mvp"),
		APIPort:        getEnv("API_PORT", "3000"),
		IngressURL:     getEnv("INGRESS_URL", "https://specs-api.uat.dextech.ai/sgbuildex"),
		PitstopURL:     getEnv("PITSTOP_URL", "https://ca-me-sgbuildex.pitstop.uat.dextech.ai"),
		JWTSecret:      getEnvRequired("JWT_SECRET"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", ""),

		WorkerIntervalMinutes: getEnvInt("WORKER_INTERVAL_MINUTES", 5),
	}

	// Enforce strong JWT secret (#11)
	if len(cfg.JWTSecret) < 32 {
		logger.Fatalf("[CONFIG] FATAL: JWT_SECRET must be at least 32 characters for production security.")
	}

	cfg.DBDSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvRequired returns the env var value or fatally exits if not set.
func getEnvRequired(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		logger.Fatalf("[CONFIG] FATAL: Required environment variable %q is not set. Please configure your .env file.", key)
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	str := getEnv(key, "")
	if str == "" {
		return fallback
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		logger.Infof("Warning: Invalid value for %s: %v. Using fallback: %d", key, err, fallback)
		return fallback
	}
	return val
}
