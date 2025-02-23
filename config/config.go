// config/config.go
package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds application configuration values.
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr string
	RedisPass string
	RedisDB   int
}

// LoadConfig reads environment variables and returns a Config struct.
func LoadConfig() *Config {
	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "123456"),
		DBName:     getEnv("DB_NAME", "chat_websocket"),

		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: getEnv("REDIS_PASSWORD", ""),
		RedisDB:   getEnvAsInt("REDIS_DB", 0),
	}
	log.Printf("[CONFIG] Loaded: %+v\n", cfg)
	return cfg
}

// getEnv retrieves the value of the environment variable or returns defaultVal if not set.
func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvAsInt retrieves the integer value of the environment variable or returns defaultVal if not set/invalid.
func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
