package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"dklautomationgo/logger"
)

// LoggerConfig bevat alle logger configuratie
type LoggerConfig struct {
	Level            string        // Log niveau (debug, info, warn, error)
	ELKEnabled       bool          // Of ELK logging is ingeschakeld
	ELKEndpoint      string        // ELK endpoint URL
	ELKBatchSize     int           // Batch grootte voor ELK
	ELKFlushInterval time.Duration // Flush interval voor ELK
	ELKAppName       string        // Applicatie naam voor ELK
	ELKEnvironment   string        // Omgeving voor ELK
}

// LoadLoggerConfig laadt logger configuratie uit environment variables
func LoadLoggerConfig() *LoggerConfig {
	c := &LoggerConfig{}

	c.Level = getEnv("LOG_LEVEL", logger.InfoLevel)
	c.ELKEndpoint = os.Getenv("ELK_ENDPOINT")
	c.ELKEnabled = c.ELKEndpoint != ""

	if c.ELKEnabled {
		c.ELKBatchSize = getEnvInt("ELK_BATCH_SIZE", 100)
		c.ELKFlushInterval = getEnvDuration("ELK_FLUSH_INTERVAL", 5*time.Second)

		// Default naam generiek gemaakt voor herbruikbaarheid
		c.ELKAppName = getEnv("ELK_APP_NAME", "unknown-service")

		// Slimme fallback voor environment
		env := getEnv("ELK_ENVIRONMENT", os.Getenv("ENVIRONMENT"))
		if env == "" {
			env = "development"
		}
		c.ELKEnvironment = env
	}

	return c
}

// SetupLogger configureert de logger op basis van de configuratie
func (c *LoggerConfig) SetupLogger() error {
	// Initialiseer de logger met het opgegeven niveau
	logger.Setup(c.Level)

	// Setup ELK integratie als ingeschakeld
	if c.ELKEnabled {
		logger.SetupELK(logger.ELKConfig{
			Endpoint:      c.ELKEndpoint,
			BatchSize:     c.ELKBatchSize,
			FlushInterval: c.ELKFlushInterval,
			AppName:       c.ELKAppName,
			Environment:   c.ELKEnvironment,
		})
		logger.Info("ELK logging enabled", "endpoint", c.ELKEndpoint)
	}

	return nil
}

// --- Helper functies ---

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			return i
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return fallback
}

// ValidateLoggerConfig valideert de logger configuratie
func (c *LoggerConfig) ValidateLoggerConfig() error {
	// Valideer log niveau (switch is sneller en schoner dan een map allocatie)
	switch c.Level {
	case logger.DebugLevel, logger.InfoLevel, logger.WarnLevel, logger.ErrorLevel:
		// OK
	default:
		return fmt.Errorf("ongeldig log niveau: %s (geldige waarden: debug, info, warn, error)", c.Level)
	}

	// Valideer ELK configuratie als ingeschakeld
	if c.ELKEnabled {
		if c.ELKEndpoint == "" {
			return fmt.Errorf("ELK_ENDPOINT is verplicht als ELK logging is ingeschakeld")
		}
		if c.ELKBatchSize <= 0 {
			return fmt.Errorf("ELK_BATCH_SIZE moet groter zijn dan 0")
		}
		if c.ELKFlushInterval <= 0 {
			return fmt.Errorf("ELK_FLUSH_INTERVAL moet groter zijn dan 0")
		}
	}

	return nil
}
