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

		// Gebruik de exacte verbindingsgegevens voor de Render PostgreSQL database
		// Interne hostname: dpg-cva4c01c1ekc738q6q0g-a
		// Externe hostname: dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com
		possibleHosts := []string{
			"dpg-cva4c01c1ekc738q6q0g-a",                            // Interne hostname
			"dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com", // Externe hostname
		}

		dkllogger.Info("Exacte PostgreSQL hostnamen", "hosts", possibleHosts)

		// Gebruik de exacte verbindingsgegevens
		port := "5432"
		user := "dekoninklijkeloopdatabase_user"
		password := "I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB"
		dbName := "dekoninklijkeloopdatabase"
		sslMode := "require" // Render vereist SSL voor externe verbindingen

		// Gebruik de interne hostname als standaard
		host := possibleHosts[0]

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

	// Als we in productie draaien, probeer de exacte hostnamen
	if os.Getenv("APP_ENV") == "prod" {
		// Lijst van exacte hostnamen voor Render PostgreSQL
		possibleHosts := []string{
			"dpg-cva4c01c1ekc738q6q0g-a",                            // Interne hostname
			"dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com", // Externe hostname
		}

		// Probeer elke hostname
		var lastErr error
		for _, host := range possibleHosts {
			// Maak een kopie van de configuratie met de nieuwe hostname
			testConfig := *config
			testConfig.Host = host
			dsn := testConfig.ConnectionString()

			dkllogger.Info("Probeer database verbinding", "host", host, "dsn", dsn)

			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: gormLog,
			})

			if err == nil {
				dkllogger.Info("Database verbinding succesvol", "host", host)

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

			dkllogger.Warn("Database verbinding mislukt", "host", host, "error", err)
			lastErr = err
		}

		// Als alle verbindingen mislukken, geef de laatste fout terug
		return nil, fmt.Errorf("alle database verbindingen mislukt, laatste fout: %w", lastErr)
	}

	// Voor niet-productie omgevingen, gebruik de normale verbinding
	dsn := config.ConnectionString()
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
