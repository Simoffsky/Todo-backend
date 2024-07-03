package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func NewEnvConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING:", err) // not returning error because we can work without .env file
	}
	return &Config{
		LoggerLevel: getEnv("LOGGER_LEVEL", "DEBUG"),
		HTTPTimeout: time.Duration(getEnvAsInt("HTTP_TIMEOUT", 5)) * time.Second,
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}