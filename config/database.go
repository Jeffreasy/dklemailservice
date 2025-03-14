package config

import (
	"fmt"
	"os"
	"time"

	dkllogger "dklautomationgo/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DatabaseConfig bevat alle database configuratie
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadDatabaseConfig laadt database configuratie uit environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "dklemailservice"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

// ConnectionString genereert een database connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// InitDatabase initialiseert de database verbinding
func InitDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := config.ConnectionString()

	// Configureer GORM logger
	gormLog := gormlogger.New(
		&dbLogger{},
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLog,
	})

	if err != nil {
		return nil, err
	}

	// Configureer connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Helper functie om environment variables te lezen met fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		dkllogger.Debug("Environment variable gelezen", "key", key, "value", value)
		return value
	}
	dkllogger.Warn("Environment variable niet gevonden, fallback gebruikt", "key", key, "fallback", fallback)
	return fallback
}

// dbLogger implementeert de GORM logger interface en gebruikt onze eigen logger
type dbLogger struct{}

func (l *dbLogger) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	dkllogger.Debug("Database", "message", msg)
}
