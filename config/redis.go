package config

import (
	"context"
	"dklautomationgo/logger"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig bevat de Redis configuratie
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool
}

// LoadRedisConfig laadt de Redis configuratie uit omgevingsvariabelen
func LoadRedisConfig() *RedisConfig {
	config := &RedisConfig{
		Host:     getEnvWithDefault("REDIS_HOST", "localhost"),
		Port:     getEnvWithDefault("REDIS_PORT", "6379"),
		Password: os.Getenv("REDIS_PASSWORD"), // Leeg als niet ingesteld
		Enabled:  getEnvWithDefault("REDIS_ENABLED", "false") == "true",
	}

	// Parse DB nummer
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			config.DB = db
		} else {
			logger.Warn("Invalid REDIS_DB value, using default", "value", dbStr, "default", config.DB)
		}
	}

	logger.Info("Redis configuratie geladen",
		"enabled", config.Enabled,
		"host", config.Host,
		"port", config.Port,
		"db", config.DB,
		"has_password", config.Password != "")

	return config
}

// NewRedisClient maakt een nieuwe Redis client aan
func NewRedisClient(config *RedisConfig) *redis.Client {
	if !config.Enabled {
		logger.Info("Redis is uitgeschakeld, geen client aangemaakt")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Host + ":" + config.Port,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	// Test de verbinding
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Error("Redis client kon niet worden geïnitialiseerd", "error", err)
		return nil
	}

	logger.Info("Redis client succesvol geïnitialiseerd")
	return client
}

// getEnvWithDefault haalt een omgevingsvariabele op met een standaardwaarde
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
