package main

import (
	"dklautomationgo/config"
	"dklautomationgo/database"
	"dklautomationgo/handlers"
	"dklautomationgo/logger"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// ValidateEnv controleert of alle benodigde omgevingsvariabelen zijn ingesteld
func ValidateEnv() error {
	required := []string{
		// Algemene SMTP configuratie
		"SMTP_HOST",
		"SMTP_USER",
		"SMTP_PASSWORD",
		"SMTP_FROM",

		// Registratie SMTP configuratie
		"REGISTRATION_SMTP_HOST",
		"REGISTRATION_SMTP_USER",
		"REGISTRATION_SMTP_PASSWORD",
		"REGISTRATION_SMTP_FROM",

		// Email adressen
		"ADMIN_EMAIL",
		"REGISTRATION_EMAIL",

		// Database configuratie
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSL_MODE",

		// JWT configuratie
		"JWT_SECRET",
	}

	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("ontbrekende omgevingsvariabele: %s", env)
		}
	}

	return nil
}

func main() {
	// Laad .env bestand als het bestaat
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		logger.Warn("Kon .env bestand niet laden", "error", err)
	}

	// Initialiseer de logger met niveau uit omgevingsvariabele of standaard INFO
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = logger.InfoLevel
	}
	logger.Setup(logLevel)
	defer logger.Sync()

	// Debug: Print alle omgevingsvariabelen
	logger.Info("Omgevingsvariabelen debug:")
	for _, env := range []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
		"SMTP_HOST", "SMTP_PORT", "SMTP_USER", "SMTP_PASSWORD", "SMTP_FROM",
		"REGISTRATION_SMTP_HOST", "REGISTRATION_SMTP_PORT", "REGISTRATION_SMTP_USER",
		"REGISTRATION_SMTP_PASSWORD", "REGISTRATION_SMTP_FROM",
		"ADMIN_EMAIL", "REGISTRATION_EMAIL",
		"JWT_SECRET",
	} {
		value := os.Getenv(env)
		if value == "" {
			logger.Warn("Omgevingsvariabele niet gevonden", "key", env)
		} else {
			// Verberg wachtwoorden in logs
			if strings.Contains(env, "PASSWORD") {
				logger.Info("Omgevingsvariabele gevonden", "key", env, "value", "********")
			} else {
				logger.Info("Omgevingsvariabele gevonden", "key", env, "value", value)
			}
		}
	}

	// Setup ELK integratie als omgevingsvariabele is ingesteld
	elkEndpoint := os.Getenv("ELK_ENDPOINT")
	if elkEndpoint != "" {
		logger.SetupELK(logger.ELKConfig{
			Endpoint:      elkEndpoint,
			BatchSize:     100,
			FlushInterval: 5 * time.Second,
			AppName:       "dklemailservice",
			Environment:   os.Getenv("ENVIRONMENT"),
		})
		logger.Info("ELK logging enabled", "endpoint", elkEndpoint)
	}

	logger.Info("DKL Email Service wordt gestart", "version", handlers.Version)

	// Controleer omgevingsvariabelen
	if err := ValidateEnv(); err != nil {
		logger.Fatal("Configuratiefout", "error", err)
	}

	// Initialiseer database
	dbConfig := config.LoadDatabaseConfig()

	// Log database configuratie voor debugging
	logger.Info("Database configuratie geladen",
		"host", dbConfig.Host,
		"port", dbConfig.Port,
		"user", dbConfig.User,
		"dbname", dbConfig.DBName,
		"sslmode", dbConfig.SSLMode)

	// Test database verbinding direct
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)
	logger.Info("Probeer directe database verbinding", "connection_string", connectionString)

	db, err := config.InitDatabase(dbConfig)
	if err != nil {
		logger.Fatal("Database initialisatie fout", "error", err)
	}

	// Initialiseer repository factory
	repoFactory := repository.NewRepository(db)

	// Voer database migraties uit
	migrationManager := database.NewMigrationManager(db, repoFactory.Migratie)
	if err := migrationManager.MigrateDatabase(); err != nil {
		logger.Fatal("Database migratie fout", "error", err)
	}

	// Seed database met initiële data
	if err := migrationManager.SeedDatabase(); err != nil {
		logger.Fatal("Database seeding fout", "error", err)
	}

	// Initialiseer service factory
	serviceFactory := services.NewServiceFactory(repoFactory)

	// Type assertion voor RateLimiter
	rateLimiter, ok := serviceFactory.RateLimiter.(*services.RateLimiter)
	if !ok {
		logger.Fatal("Kon RateLimiter niet casten naar juiste type")
	}

	// Initialiseer handlers
	emailHandler := handlers.NewEmailHandler(serviceFactory.EmailService)
	authHandler := handlers.NewAuthHandler(serviceFactory.AuthService, rateLimiter)
	metricsHandler := handlers.NewMetricsHandler(serviceFactory.EmailMetrics, rateLimiter)

	// Initialiseer nieuwe handlers voor contact en aanmelding beheer
	contactHandler := handlers.NewContactHandler(
		repoFactory.Contact,
		repoFactory.ContactAntwoord,
		serviceFactory.EmailService,
		serviceFactory.AuthService,
	)

	aanmeldingHandler := handlers.NewAanmeldingHandler(
		repoFactory.Aanmelding,
		repoFactory.AanmeldingAntwoord,
		serviceFactory.EmailService,
		serviceFactory.AuthService,
	)

	// Configureer en initialiseer de mail fetcher service
	mailFetcher := initializeMailFetcher(serviceFactory.EmailMetrics)
	mailHandler := handlers.NewMailHandler(mailFetcher, repoFactory.IncomingEmail, serviceFactory.AuthService)

	// Maak een EmailAutoFetcher aan voor automatisch ophalen van emails
	emailAutoFetcher := services.NewEmailAutoFetcher(mailFetcher, repoFactory.IncomingEmail)

	// Sla de emailAutoFetcher op in de serviceFactory
	serviceFactory.EmailAutoFetcher = emailAutoFetcher

	// Start de automatische email fetcher als deze niet is uitgeschakeld
	if os.Getenv("DISABLE_AUTO_EMAIL_FETCH") != "true" {
		logger.Info("Automatisch ophalen van emails starten...")
		serviceFactory.EmailAutoFetcher.Start()
		logger.Info("Automatische email fetcher gestart")
	} else {
		logger.Info("Automatisch ophalen van emails is uitgeschakeld")
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error("Request fout",
				"path", c.Path(),
				"method", c.Method(),
				"error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Er is een fout opgetreden bij het verwerken van je verzoek",
			})
		},
	})

	// Configure CORS
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "") {
		allowedOrigins = []string{"https://www.dekoninklijkeloop.nl", "https://dekoninklijkeloop.nl", "http://localhost:3000"}
	}

	logger.Info("CORS geconfigureerd", "origins", allowedOrigins)

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(allowedOrigins, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Voeg TestModeMiddleware toe als globale middleware
	app.Use(handlers.TestModeMiddleware())

	// Serve static files from public directory
	app.Static("/", "./public")

	// Specific route for favicon.ico
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		// Get the current working directory
		workDir, err := os.Getwd()
		if err != nil {
			logger.Error("Kon werkdirectory niet bepalen", "error", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		faviconPath := filepath.Join(workDir, "public", "favicon.ico")
		if _, err := os.Stat(faviconPath); os.IsNotExist(err) {
			logger.Error("Favicon niet gevonden", "path", faviconPath, "error", err)
			return c.SendStatus(fiber.StatusNotFound)
		}
		c.Set("Content-Type", "image/x-icon")
		c.Set("Cache-Control", "public, max-age=31536000") // Cache voor 1 jaar
		return c.SendFile(faviconPath, false)
	})

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":     "DKL Email Service API",
			"version":     handlers.Version,
			"status":      "running",
			"environment": os.Getenv("ENVIRONMENT"),
			"timestamp":   time.Now(),
			"endpoints": []fiber.Map{
				{"path": "/api/health", "method": "GET", "description": "Service health status"},
				{"path": "/api/contact-email", "method": "POST", "description": "Send contact form email"},
				{"path": "/api/aanmelding-email", "method": "POST", "description": "Send registration form email"},
				{"path": "/api/metrics/email", "method": "GET", "description": "Email metrics (requires API key)"},
				{"path": "/api/metrics/rate-limits", "method": "GET", "description": "Rate limit metrics (requires API key)"},
				{"path": "/api/auth/login", "method": "POST", "description": "User login"},
				{"path": "/api/auth/logout", "method": "POST", "description": "User logout"},
				{"path": "/api/auth/profile", "method": "GET", "description": "Get user profile (requires auth)"},
				{"path": "/api/auth/reset-password", "method": "POST", "description": "Reset password (requires auth)"},
				{"path": "/api/contact", "method": "GET", "description": "List contact forms (requires admin auth)"},
				{"path": "/api/contact/:id", "method": "GET", "description": "Get contact form details (requires admin auth)"},
				{"path": "/api/contact/:id", "method": "PUT", "description": "Update contact form (requires admin auth)"},
				{"path": "/api/contact/:id", "method": "DELETE", "description": "Delete contact form (requires admin auth)"},
				{"path": "/api/contact/:id/antwoord", "method": "POST", "description": "Add reply to contact form (requires admin auth)"},
				{"path": "/api/contact/status/:status", "method": "GET", "description": "Filter contact forms by status (requires admin auth)"},
				{"path": "/api/aanmelding", "method": "GET", "description": "List registrations (requires admin auth)"},
				{"path": "/api/aanmelding/:id", "method": "GET", "description": "Get registration details (requires admin auth)"},
				{"path": "/api/aanmelding/:id", "method": "PUT", "description": "Update registration (requires admin auth)"},
				{"path": "/api/aanmelding/:id", "method": "DELETE", "description": "Delete registration (requires admin auth)"},
				{"path": "/api/aanmelding/:id/antwoord", "method": "POST", "description": "Add reply to registration (requires admin auth)"},
				{"path": "/api/aanmelding/rol/:rol", "method": "GET", "description": "Filter registrations by role (requires admin auth)"},
				{"path": "/metrics", "method": "GET", "description": "Prometheus metrics"},
			},
		})
	})

	// API routes group
	api := app.Group("/api")

	// Health check endpoint
	api.Get("/health", handlers.HealthHandler)

	// Email routes
	api.Post("/contact-email", emailHandler.HandleContactEmail)
	api.Post("/aanmelding-email", emailHandler.HandleAanmeldingEmail)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handlers.RateLimitMiddleware(rateLimiter, "login"), authHandler.HandleLogin)
	auth.Post("/logout", authHandler.HandleLogout)

	// Beveiligde auth routes (vereisen authenticatie)
	authProtected := auth.Group("/", handlers.AuthMiddleware(serviceFactory.AuthService))
	authProtected.Get("/profile", authHandler.HandleGetProfile)
	authProtected.Post("/reset-password", authHandler.HandleResetPassword)

	// Admin routes (vereisen admin rol)
	// Commentaar: admin routes worden momenteel niet gebruikt, maar kunnen later worden toegevoegd
	// admin := api.Group("/admin", handlers.AuthMiddleware(serviceFactory.AuthService), handlers.AdminMiddleware(serviceFactory.AuthService))

	// Metrics endpoints direct onder /api/metrics/... (vereisen API key)
	api.Get("/metrics/email", metricsHandler.HandleGetEmailMetrics)
	api.Get("/metrics/rate-limits", metricsHandler.HandleGetRateLimits)

	// Registreer routes voor contact en aanmelding beheer
	contactHandler.RegisterRoutes(app)
	aanmeldingHandler.RegisterRoutes(app)

	// Registreer de mailHandler in de main functie na repo en authService
	mailHandler.RegisterRoutes(app)

	// Voeg Prometheus metrics endpoint toe aan Fiber app in plaats van standaard HTTP server
	app.Get("/metrics", func(c *fiber.Ctx) error {
		// Eenvoudige implementatie die een string teruggeeft
		// In een volledige implementatie zou je hier de Prometheus metrics moeten teruggeven
		return c.Status(fiber.StatusOK).SendString("Prometheus metrics endpoint")
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 for web traffic
	}

	// Start server in een goroutine met Fiber's eigen methoden
	go func() {
		logger.Info("Server gestart", "port", port)
		if err := app.Listen(":" + port); err != nil {
			logger.Fatal("Server fout", "error", err)
		}
	}()

	// Wacht op interrupt signaal (CTRL+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("Server wordt afgesloten...")

	// Graceful shutdown
	if serviceFactory.EmailBatcher != nil {
		serviceFactory.EmailBatcher.Shutdown()
	}

	// Stop de email auto fetcher
	if serviceFactory.EmailAutoFetcher != nil && serviceFactory.EmailAutoFetcher.IsRunning() {
		logger.Info("Email auto fetcher stoppen...")
		serviceFactory.EmailAutoFetcher.Stop()
		logger.Info("Email auto fetcher gestopt")
	}

	// Sluit rate limiter af
	if rateLimiter != nil {
		rateLimiter.Shutdown()
	}

	// Log laatste metrics
	serviceFactory.EmailMetrics.LogMetrics()

	// Sluit alle log writers
	logger.CloseWriters()

	// Graceful shutdown met Fiber
	if err := app.Shutdown(); err != nil {
		logger.Fatal("Server shutdown fout", "error", err)
	}

	logger.Info("Server succesvol afgesloten")
}

// Configureer en initialiseer de mail fetcher service
func initializeMailFetcher(metrics *services.EmailMetrics) *services.MailFetcher {
	mailFetcher := services.NewMailFetcher(metrics)

	// Voeg de accounts toe
	mailFetcher.AddAccount(
		"info@dekoninklijkeloop.nl",
		"X6FdeT5dTakH^Ae^f5BV",
		"mail.hostnet.nl",
		993, // IMAP SSL poort
		"info",
	)

	mailFetcher.AddAccount(
		"inschrijving@dekoninklijkeloop.nl",
		"y9tKkS&^1pbp2X6KuUbb",
		"mail.hostnet.nl",
		993, // IMAP SSL poort
		"inschrijving",
	)

	return mailFetcher
}
