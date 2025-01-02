package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string
	Storage    Storage
	HTTPServer HTTPServer
}

type HTTPServer struct {
	Address     string
	TimeOut     time.Duration
	IdleTimeOut time.Duration
}

type Storage struct {
	Path     string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	AsString string
}

// MustLoad loads the configuration from environment variables
func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	env := getEnv("ENV", "local")
	storagePort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("invalid DB_PORT value: %v", err)
	}

	httpTimeout, err := time.ParseDuration(getEnv("HTTP_SERVER_TIMEOUT", "5s"))
	if err != nil {
		log.Fatalf("invalid HTTP_SERVER_TIMEOUT value: %v", err)
	}

	httpIdleTimeout, err := time.ParseDuration(getEnv("HTTP_SERVER_IDLE_TIMEOUT", "60s"))
	if err != nil {
		log.Fatalf("invalid HTTP_SERVER_IDLE_TIMEOUT value: %v", err)
	}

	cfg := &Config{
		Env: env,
		Storage: Storage{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     storagePort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "storage"),
		},
		HTTPServer: HTTPServer{
			Address:     getEnv("HTTP_SERVER_ADDRESS", "localhost:8082"),
			TimeOut:     httpTimeout,
			IdleTimeOut: httpIdleTimeout,
		},
	}

	cfg.Storage.AsString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.DBName)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
