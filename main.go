package main

import (
	"dklautomationgo/handlers"
	"dklautomationgo/logger"
	"dklautomationgo/services"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// Initialiseer Prometheus metrics
	prometheusMetrics := services.GetPrometheusMetrics()

	// SMTP client setup
	smtpClient := services.NewRealSMTPClient(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_FROM"),
		os.Getenv("REGISTRATION_SMTP_HOST"),
		os.Getenv("REGISTRATION_SMTP_PORT"),
		os.Getenv("REGISTRATION_SMTP_USER"),
		os.Getenv("REGISTRATION_SMTP_PASSWORD"),
		os.Getenv("REGISTRATION_SMTP_FROM"),
	)

	// Initialize services
	emailMetrics := services.NewEmailMetrics(24 * time.Hour) // Reset elke 24 uur
	rateLimiter := services.NewRateLimiter(prometheusMetrics)
	rateLimiter.AddLimit("contact_email", 100, time.Hour, false)    // 100 contact emails per uur globaal
	rateLimiter.AddLimit("contact_email", 5, time.Hour, true)       // 5 contact emails per uur per IP
	rateLimiter.AddLimit("aanmelding_email", 200, time.Hour, false) // 200 aanmelding emails per uur globaal
	rateLimiter.AddLimit("aanmelding_email", 10, time.Hour, true)   // 10 aanmelding emails per uur per IP
	emailService := services.NewEmailService(smtpClient, emailMetrics, rateLimiter, prometheusMetrics)
	emailBatcher := services.NewEmailBatcher(emailService, 50, 15*time.Minute)
	metricsHandler := handlers.NewMetricsHandler(emailMetrics, rateLimiter)

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)

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
		allowedOrigins = []string{"https://www.dekoninklijkeloop.nl", "https://dekoninklijkeloop.nl"}
	}

	logger.Info("CORS geconfigureerd", "origins", allowedOrigins)

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(allowedOrigins, ","),
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	// Serve static files from public directory
	app.Static("/", "./public")

	// Specific route for favicon.ico
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendFile("./public/favicon.ico")
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

	// Metrics endpoint toevoegen
	api.Get("/metrics/email", metricsHandler.HandleGetEmailMetrics)
	api.Get("/metrics/rate-limits", metricsHandler.HandleGetRateLimits)

	// Voeg Prometheus metrics endpoint toe aan server
	http.Handle("/metrics", promhttp.Handler())

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
	if emailBatcher != nil {
		emailBatcher.Shutdown()
	}

	// Sluit rate limiter af
	rateLimiter.Shutdown()

	// Log laatste metrics
	emailMetrics.LogMetrics()

	// Sluit alle log writers
	logger.CloseWriters()

	// Graceful shutdown met Fiber
	if err := app.Shutdown(); err != nil {
		logger.Fatal("Server shutdown fout", "error", err)
	}

	logger.Info("Server succesvol afgesloten")
}
