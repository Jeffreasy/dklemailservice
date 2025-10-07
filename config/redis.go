package config

import (
	"context"
	"dklautomationgo/logger"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	// Als REDIS_URL is ingesteld (bijv. van Render), parse deze
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		if parsedConfig, err := parseRedisURL(redisURL); err == nil {
			config = parsedConfig
			config.Enabled = true // Enable Redis als URL is opgegeven
		} else {
			logger.Warn("Failed to parse REDIS_URL, falling back to individual settings", "error", err)
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

// parseRedisURL parseert een Redis URL in de vorm redis://username:password@host:port/db
func parseRedisURL(url string) (*RedisConfig, error) {
	// Verwacht formaat: redis://username:password@host:port/db
	if !strings.HasPrefix(url, "redis://") {
		return nil, fmt.Errorf("invalid Redis URL format")
	}

	// Verwijder "redis://" prefix
	withoutPrefix := strings.TrimPrefix(url, "redis://")

	// Split op "@" voor credentials en host
	parts := strings.Split(withoutPrefix, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid Redis URL format")
	}

	// Parse credentials (username:password)
	creds := parts[0]
	var password string
	if strings.Contains(creds, ":") {
		credParts := strings.Split(creds, ":")
		if len(credParts) == 2 {
			password = credParts[1]
		}
	}

	// Parse host:port/db
	hostPart := parts[1]
	hostPortDB := strings.Split(hostPart, "/")
	if len(hostPortDB) < 1 {
		return nil, fmt.Errorf("invalid Redis URL format")
	}

	hostPort := hostPortDB[0]
	var db int
	if len(hostPortDB) > 1 && hostPortDB[1] != "" {
		if parsedDB, err := strconv.Atoi(hostPortDB[1]); err == nil {
			db = parsedDB
		}
	}

	// Split host:port
	hostPortParts := strings.Split(hostPort, ":")
	if len(hostPortParts) != 2 {
		return nil, fmt.Errorf("invalid host:port format")
	}

	host := hostPortParts[0]
	port := hostPortParts[1]

	return &RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
		Enabled:  true,
	}, nil
}

// getEnvWithDefault haalt een omgevingsvariabele op met een standaardwaarde
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
