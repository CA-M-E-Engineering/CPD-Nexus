package config

import (
	"fmt"
	"log"
	"os"
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
	BridgeURL  string
	IngressURL string

	BridgeIntervalSeconds int
	WorkerIntervalMinutes int
}

func LoadConfig() *Config {
	// Try loading .env from common locations
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	cfg := &Config{
		DBUser:     getEnv("DB_USER", "bas_user"),
		DBPass:     getEnv("DB_PASS", "new_password"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1:3306"),
		DBName:     getEnv("DB_NAME", "bas_mvp"),
		APIPort:    getEnv("API_PORT", "3000"),
		BridgeURL:  getEnv("BRIDGE_URL", "ws://localhost:8080/ws"),
		IngressURL: getEnv("INGRESS_URL", "https://specs-api.uat.dextech.ai/sgbuildex"),

		BridgeIntervalSeconds: getEnvInt("BRIDGE_INTERVAL_SECONDS", 10),
		WorkerIntervalMinutes: getEnvInt("WORKER_INTERVAL_MINUTES", 5),
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

func getEnvInt(key string, fallback int) int {
	str := getEnv(key, "")
	if str == "" {
		return fallback
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("Warning: Invalid value for %s: %v. Using fallback: %d", key, err, fallback)
		return fallback
	}
	return val
}
