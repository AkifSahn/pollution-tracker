package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ServerHost string
	ServerPort string
}

var cfg *Config

func LoadConfig() *Config {

	if cfg != nil {
		return cfg
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	cfg = &Config{
		DBUser:     getEnv("DB_USER", "timescale"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "pollution"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerHost: getEnv("SERVER_HOST", "127.0.0.1"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
