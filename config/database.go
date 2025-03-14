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

		// Probeer verschillende mogelijke hostnamen voor Render PostgreSQL
		// Volgens Render documentatie kan de interne URL verschillende vormen hebben
		possibleHosts := []string{
			"postgres",                            // Standaard hostname in Render
			"dklautomatie-db",                     // Service naam
			"postgresql-dklautomatie-db",          // Conventie: postgresql-<service_name>
			"internal-postgresql-dklautomatie-db", // Conventie: internal-postgresql-<service_name>
			"dklautomatie-db.internal",            // Conventie: <service_name>.internal
			"postgresql.render.com",               // Externe hostname
		}

		dkllogger.Info("Mogelijke PostgreSQL hostnamen", "hosts", possibleHosts)

		// Lees de overige configuratie uit de omgeving
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		// Gebruik de waarden uit de Render dashboard
		user := os.Getenv("DB_USER")
		if user == "" {
			// Fallback naar een standaard gebruikersnaam voor Render
			user = "postgres"
		}

		password := os.Getenv("DB_PASSWORD")

		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "dklemailservice"
		}

		sslMode := os.Getenv("DB_SSL_MODE")
		if sslMode == "" {
			// Render vereist SSL voor externe verbindingen
			sslMode = "require"
		}

		// Gebruik de eerste hostname in de lijst als standaard
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

	// Als we in productie draaien, probeer alle mogelijke hostnamen
	if os.Getenv("APP_ENV") == "prod" {
		// Lijst van mogelijke hostnamen voor Render PostgreSQL
		possibleHosts := []string{
			"postgres",                            // Standaard hostname in Render
			"dklautomatie-db",                     // Service naam
			"postgresql-dklautomatie-db",          // Conventie: postgresql-<service_name>
			"internal-postgresql-dklautomatie-db", // Conventie: internal-postgresql-<service_name>
			"dklautomatie-db.internal",            // Conventie: <service_name>.internal
			"postgresql.render.com",               // Externe hostname
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
