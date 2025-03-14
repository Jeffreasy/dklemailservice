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
	// Controleer of we in productie draaien
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "prod" {
		// In productie, gebruik de Render PostgreSQL database
		dkllogger.Info("Productieomgeving gedetecteerd, gebruik Render PostgreSQL configuratie")

		// Haal de database service naam op uit de omgeving
		dbServiceName := os.Getenv("RENDER_DB_SERVICE_NAME")
		if dbServiceName == "" {
			dbServiceName = "dklautomatie-db"
			dkllogger.Info("RENDER_DB_SERVICE_NAME niet gevonden, gebruik standaard naam", "default", dbServiceName)
		}

		// Bouw de hostname op basis van de Render conventies
		// Render gebruikt de vorm: postgresql-<service_name>
		host := fmt.Sprintf("postgresql-%s", dbServiceName)
		dkllogger.Info("Render PostgreSQL hostname opgebouwd", "host", host)

		// Lees de overige configuratie uit de omgeving
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "dklemailservice"
		}

		sslMode := os.Getenv("DB_SSL_MODE")
		if sslMode == "" {
			sslMode = "require"
		}

		dkllogger.Info("Render PostgreSQL configuratie geladen",
			"host", host,
			"port", port,
			"user", user,
			"dbname", dbName,
			"sslmode", sslMode)

		return &DatabaseConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			DBName:   dbName,
			SSLMode:  sslMode,
		}
	}

	// Lees direct de omgevingsvariabelen
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")

	// Log de waarden (verberg wachtwoord)
	dkllogger.Info("Database configuratie direct uit omgevingsvariabelen:",
		"DB_HOST", host,
		"DB_PORT", port,
		"DB_USER", user,
		"DB_NAME", dbName,
		"DB_SSL_MODE", sslMode)

	// Gebruik fallback waarden alleen als de omgevingsvariabelen leeg zijn
	if host == "" {
		host = "localhost"
		dkllogger.Warn("DB_HOST omgevingsvariabele niet gevonden, fallback gebruikt", "fallback", host)
	}
	if port == "" {
		port = "5432"
		dkllogger.Warn("DB_PORT omgevingsvariabele niet gevonden, fallback gebruikt", "fallback", port)
	}
	if user == "" {
		user = "postgres"
		dkllogger.Warn("DB_USER omgevingsvariabele niet gevonden, fallback gebruikt", "fallback", user)
	}
	if dbName == "" {
		dbName = "dklemailservice"
		dkllogger.Warn("DB_NAME omgevingsvariabele niet gevonden, fallback gebruikt", "fallback", dbName)
	}
	if sslMode == "" {
		sslMode = "disable"
		dkllogger.Warn("DB_SSL_MODE omgevingsvariabele niet gevonden, fallback gebruikt", "fallback", sslMode)
	}

	return &DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
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

	// Probeer eerst de normale verbinding
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLog,
	})

	// Als dat niet lukt en we zijn in productie, probeer een alternatieve verbinding
	if err != nil && os.Getenv("APP_ENV") == "prod" {
		dkllogger.Warn("Normale database verbinding mislukt, probeer alternatieve verbinding", "error", err)

		// Probeer een alternatieve verbinding met de interne Render hostname
		altConfig := *config
		altConfig.Host = "internal-postgresql-" + os.Getenv("RENDER_DB_SERVICE_NAME")
		altDsn := altConfig.ConnectionString()

		dkllogger.Info("Probeer alternatieve database verbinding", "host", altConfig.Host)
		db, err = gorm.Open(postgres.Open(altDsn), &gorm.Config{
			Logger: gormLog,
		})

		if err != nil {
			dkllogger.Error("Alternatieve database verbinding mislukt", "error", err)
			return nil, err
		}

		dkllogger.Info("Alternatieve database verbinding succesvol")
	} else if err != nil {
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
