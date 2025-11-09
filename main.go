package main

import (
	"dklautomationgo/config"
	"dklautomationgo/database"
	"dklautomationgo/handlers"
	"dklautomationgo/logger"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"fmt"
	"net/http/httptest"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
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

	// Controleer email fetcher configuratie indien ingeschakeld
	if os.Getenv("DISABLE_AUTO_EMAIL_FETCH") != "true" {
		emailFetcherVars := []string{
			"INFO_EMAIL",
			"INFO_EMAIL_PASSWORD",
			"INSCHRIJVING_EMAIL",
			"INSCHRIJVING_EMAIL_PASSWORD",
		}

		missingVars := []string{}
		for _, env := range emailFetcherVars {
			if os.Getenv(env) == "" {
				missingVars = append(missingVars, env)
			}
		}

		if len(missingVars) > 0 {
			logger.Warn("Email fetcher credentials missing, some accounts will not be configured",
				"missing_vars", strings.Join(missingVars, ", "))
		}
	}

	// Whisky for Charity configuratie is optioneel
	wfcConfigured := os.Getenv("WFC_SMTP_HOST") != "" &&
		os.Getenv("WFC_SMTP_USER") != "" &&
		os.Getenv("WFC_SMTP_PASSWORD") != "" &&
		os.Getenv("WFC_SMTP_FROM") != ""

	if wfcConfigured {
		logger.Info("Whisky for Charity SMTP configuratie gevonden")
	} else {
		logger.Info("Whisky for Charity SMTP configuratie niet gevonden, deze functionaliteit is uitgeschakeld")
	}

	// Newsletter configuratie (optioneel)
	enableNewsletter := os.Getenv("ENABLE_NEWSLETTER") == "true"
	if enableNewsletter {
		if os.Getenv("NEWSLETTER_SOURCES") == "" {
			logger.Warn("ENABLE_NEWSLETTER is true maar NEWSLETTER_SOURCES is leeg")
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

	// Debug: Print alle omgevingsvariabelen alleen bij DEBUG logniveau
	if strings.ToUpper(logLevel) == logger.DebugLevel {
		logger.Debug("Omgevingsvariabelen debug:")
		for _, env := range []string{
			"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
			"SMTP_HOST", "SMTP_PORT", "SMTP_USER", "SMTP_PASSWORD", "SMTP_FROM",
			"REGISTRATION_SMTP_HOST", "REGISTRATION_SMTP_PORT", "REGISTRATION_SMTP_USER",
			"REGISTRATION_SMTP_PASSWORD", "REGISTRATION_SMTP_FROM",
			"WFC_SMTP_HOST", "WFC_SMTP_PORT", "WFC_SMTP_USER", "WFC_SMTP_PASSWORD", "WFC_SMTP_FROM",
			"ADMIN_EMAIL", "REGISTRATION_EMAIL",
			"JWT_SECRET",
		} {
			value := os.Getenv(env)
			if value == "" {
				logger.Debug("Omgevingsvariabele niet gevonden", "key", env)
			} else {
				// Verberg wachtwoorden in logs
				if strings.Contains(env, "PASSWORD") {
					logger.Debug("Omgevingsvariabele gevonden", "key", env, "value", "********")
				} else {
					logger.Debug("Omgevingsvariabele gevonden", "key", env, "value", value)
				}
			}
		}
	} else {
		logger.Info("Omgevingsvariabelen debug overgeslagen (alleen beschikbaar in DEBUG modus)")
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
	// BELANGRIJK: Zorg dat je 'repository/repository.go' hebt bijgewerkt
	// zodat deze 'Participant', 'ParticipantAntwoord' en 'EventRegistration' correct aanmaakt.
	repoFactory := repository.NewRepository(db)

	// Voer database migraties uit
	migrationManager := database.NewMigrationManager(db, repoFactory.Migratie)
	if err := migrationManager.MigrateDatabase(); err != nil {
		logger.Fatal("Database migratie fout", "error", err)
	}

	// Initialiseer service factory
	serviceFactory := services.NewServiceFactory(repoFactory)

	// Initialiseer steps service
	// GEWIJZIGD: Gebruikt nu EventRegistrationRepo (aangezien steps daar nu op staan)
	stepsService := services.NewStepsService(db, repoFactory.Participant, repoFactory.Distance)

	// ✨ NIEUWE: Initialize StepsHub voor WebSocket real-time updates
	stepsHub := services.NewStepsHub(stepsService, serviceFactory.GamificationService)

	// ✨ NIEUWE: Link hub to service voor broadcasts
	stepsService.SetStepsHub(stepsHub)

	// ✨ NIEUWE: Start hub in background goroutine
	go stepsHub.Run()
	logger.Info("StepsHub started successfully - WebSocket support enabled")

	// Start Newsletter service indien geconfigureerd
	if serviceFactory.NewsletterService != nil {
		serviceFactory.NewsletterService.Start()
	}

	// Gebruik de GetRateLimiter methode in de ServiceFactory
	rateLimiter := serviceFactory.GetRateLimiter()

	// Stel rate limiter en Redis client in voor health checks
	handlers.SetRateLimiter(rateLimiter)
	handlers.SetRedisClient(serviceFactory.RedisClient)

	// Initialiseer handlers

	// GEWIJZIGD: Injecteer ParticipantRepo en EventRegistrationRepo
	emailHandler := handlers.NewEmailHandler(
		serviceFactory.EmailService,
		serviceFactory.NotificationService,
		repoFactory.Participant,
		repoFactory.EventRegistration, // Nieuwe dependency
		repoFactory.Event,             // Nieuwe dependency
	)
	authHandler := handlers.NewAuthHandler(serviceFactory.AuthService, serviceFactory.PermissionService, rateLimiter)
	metricsHandler := handlers.NewMetricsHandler(serviceFactory.EmailMetrics, rateLimiter)

	// Initialiseer NotificationHandler
	notificationHandler := handlers.NewNotificationHandler(
		repoFactory.Notification,
		serviceFactory.NotificationService,
		serviceFactory.AuthService,
	)

	// Initialiseer nieuwe handlers voor contact en participant beheer
	contactHandler := handlers.NewContactHandler(
		repoFactory.Contact,
		repoFactory.ContactAntwoord,
		serviceFactory.EmailService,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
		serviceFactory.NotificationService,
	)

	// GEWIJZIGD: Hernoemd van AanmeldingHandler naar ParticipantHandler
	// GEWIJZIGD: Injecteer EventRegistrationRepo voor de status-update bij antwoorden
	participantHandler := handlers.NewParticipantHandler(
		repoFactory.Participant,
		repoFactory.ParticipantAntwoord,
		serviceFactory.EmailService,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
		repoFactory.EventRegistration,
	)

	// ✨ NIEUW: Initialiseer de EventRegistrationHandler voor de verplaatste logica
	eventRegistrationHandler := handlers.NewEventRegistrationHandler(
		repoFactory.EventRegistration,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)

	// Initialiseer steps handler
	// GEWIJZIGD: De permissies verwijzen mogelijk nog naar 'aanmelding'
	stepsHandler := handlers.NewStepsHandler(
		stepsService,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)

	// Initialiseer newsletter handler
	newsletterHandler := handlers.NewNewsletterHandler(
		repoFactory.Newsletter,
		serviceFactory.NewsletterSender,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)

	// Configureer en initialiseer de mail fetcher service
	mailFetcherTyped := services.NewMailFetcher(serviceFactory.EmailMetrics)
	mailHandler := handlers.NewMailHandler(mailFetcherTyped, repoFactory.IncomingEmail, serviceFactory.AuthService, serviceFactory.PermissionService)

	// Maak een EmailAutoFetcher aan voor automatisch ophalen van emails
	emailAutoFetcher := services.NewEmailAutoFetcher(mailFetcherTyped, repoFactory.IncomingEmail)

	// Configureer mail accounts als credentials beschikbaar zijn
	imapServer := os.Getenv("IMAP_SERVER")
	if imapServer == "" {
		imapServer = "imap.gmail.com" // Default fallback
	}
	imapPort := 993 // Default IMAP SSL port
	if portStr := os.Getenv("IMAP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			imapPort = port
		}
	}

	if infoEmail := os.Getenv("INFO_EMAIL"); infoEmail != "" {
		if infoPassword := os.Getenv("INFO_EMAIL_PASSWORD"); infoPassword != "" {
			mailFetcherTyped.AddAccount(infoEmail, infoPassword, imapServer, imapPort, "info")
			logger.Info("INFO email account configured", "email", infoEmail, "server", imapServer)
		}
	}
	if inschrijvingEmail := os.Getenv("INSCHRIJVING_EMAIL"); inschrijvingEmail != "" {
		if inschrijvingPassword := os.Getenv("INSCHRIJVING_EMAIL_PASSWORD"); inschrijvingPassword != "" {
			mailFetcherTyped.AddAccount(inschrijvingEmail, inschrijvingPassword, imapServer, imapPort, "inschrijving")
			logger.Info("INSCHRIJVING email account configured", "email", inschrijvingEmail, "server", imapServer)
		}
	}

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
		allowedOrigins = []string{"https://www.dekoninklijkeloop.nl", "https://dekoninklijkeloop.nl", "https://admin.dekoninklijkeloop.nl", "http://localhost:3000", "http://localhost:5173"}
	}

	logger.Info("CORS geconfigureerd", "origins", allowedOrigins)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Test-Mode",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Content-Type",
	}))

	// Voeg TestModeMiddleware toe als globale middleware
	app.Use(handlers.TestModeMiddleware())

	// Serve static files from public directory
	app.Static("/", "./public")

	// Specific route for favicon.ico
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
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

	// Root route - BIJGEWERKT MET NIEUWE ROUTES
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
				{"path": "/api/register", "method": "POST", "description": "Create new participant and event registration"}, // Hernoemd
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
				{"path": "/api/participant", "method": "GET", "description": "List participants (persons) (requires admin auth)"},                // Hernoemd
				{"path": "/api/participant/:id", "method": "GET", "description": "Get participant details (requires admin auth)"},                // Hernoemd
				{"path": "/api/participant/:id", "method": "DELETE", "description": "Delete participant (requires admin auth)"},                  // Hernoemd
				{"path": "/api/participant/:id/antwoord", "method": "POST", "description": "Add reply to participant (requires admin auth)"},     // Hernoemd
				{"path": "/api/registration/:id", "method": "GET", "description": "Get registration details (requires admin auth)"},              // NIEUW
				{"path": "/api/registration/:id", "method": "PUT", "description": "Update registration status/notes (requires admin auth)"},      // NIEUW (verplaatst)
				{"path": "/api/registration/rol/:rol", "method": "GET", "description": "Filter registrations by role (requires admin auth)"},     // NIEUW (verplaatst)
				{"path": "/api/registration/:id/steps", "method": "POST", "description": "Update steps for registration (requires steps write)"}, // Hernoemd
				{"path": "/api/registration/:id/dashboard", "method": "GET", "description": "Get registration dashboard (requires steps read)"},  // Hernoemd
				{"path": "/api/events/:id/registrations", "method": "GET", "description": "Get event registrations (requires events read)"},      // Hernoemd
				// ... (rest van je CMS en andere routes) ...
				{"path": "/api/total-steps", "method": "GET", "description": "Get total steps for year (requires steps read permission)"},
				{"path": "/api/funds-distribution", "method": "GET", "description": "Get funds distribution (requires steps read permission)"},
				{"path": "/api/events", "method": "GET", "description": "List events (public)"},
				{"path": "/api/events/active", "method": "GET", "description": "Get active event (public)"},
				{"path": "/api/events/:id", "method": "GET", "description": "Get event details (public)"},
				{"path": "/api/events", "method": "POST", "description": "Create event (requires events write permission)"},
				{"path": "/api/events/:id", "method": "PUT", "description": "Update event (requires events write permission)"},
				{"path": "/api/events/:id", "method": "DELETE", "description": "Delete event (requires events write permission)"},
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

	// GEWIJZIGD: Hernoemd van /aanmelding-email en HandleAanmeldingEmail
	api.Post("/register", emailHandler.HandleRegistrationEmail)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handlers.RateLimitMiddleware(rateLimiter, "login"), authHandler.HandleLogin)
	auth.Post("/logout", authHandler.HandleLogout)
	auth.Post("/refresh", authHandler.HandleRefreshToken)

	// Beveiligde auth routes (vereisen authenticatie)
	authProtected := auth.Group("/", handlers.AuthMiddleware(serviceFactory.AuthService))
	authProtected.Get("/profile", authHandler.HandleGetProfile)
	authProtected.Post("/reset-password", authHandler.HandleResetPassword)

	// Metrics endpoints
	api.Get("/metrics/email", metricsHandler.HandleGetEmailMetrics)
	api.Get("/metrics/rate-limits", metricsHandler.HandleGetRateLimits)

	// Registreer routes voor contact en participant beheer
	contactHandler.RegisterRoutes(app)
	participantHandler.RegisterRoutes(app) // Hernoemd

	// ✨ NIEUW: Registreer de routes voor de EventRegistrationHandler
	eventRegistrationHandler.RegisterRoutes(app)

	// Registreer routes voor stappen beheer
	stepsHandler.RegisterRoutes(app)

	// Initialiseer en registreer WebSocket handler voor steps
	stepsWsHandler := handlers.NewStepsWebSocketHandler(stepsHub, serviceFactory.AuthService)
	stepsWsHandler.RegisterRoutes(app)
	logger.Info("WebSocket routes registered - /ws/steps endpoint active")

	// WebSocket stats endpoint (admin only)
	app.Get("/api/ws/stats",
		handlers.AuthMiddleware(serviceFactory.AuthService),
		handlers.PermissionMiddleware(serviceFactory.PermissionService, "admin", "read"),
		stepsWsHandler.GetStats,
	)

	// Registreer routes voor newsletter beheer
	newsletterHandler.RegisterRoutes(app)

	// Registreer routes voor notificaties
	notificationHandler.RegisterRoutes(app)

	// Registreer de mailHandler
	mailHandler.RegisterRoutes(app)

	// Registreer de WFC routes
	handlers.RegisterWFCOrderRoutes(app, serviceFactory.EmailService)

	// Registreer telegram bot handler
	if serviceFactory.TelegramBotService != nil {
		// (Telegram routes blijven ongewijzigd)
		app.Get("/api/v1/telegrambot/config", func(c *fiber.Ctx) error {
			authHeader := c.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
			}
			if serviceFactory.TelegramBotService == nil {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"enabled":  false,
					"message":  "Telegram bot service is niet geactiveerd",
					"chatId":   "",
					"commands": []string{},
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"enabled": true,
				"message": "Telegram bot service is actief",
				"chatId":  serviceFactory.TelegramBotService.GetChatID(),
			})
		})
		app.Post("/api/v1/telegrambot/send", func(c *fiber.Ctx) error {
			authHeader := c.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
			}
			if serviceFactory.TelegramBotService == nil {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": false, "message": "Telegram bot service is niet geactiveerd"})
			}
			var req struct {
				Message string `json:"message"`
			}
			if err := c.BodyParser(&req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Ongeldige request"})
			}
			err := serviceFactory.TelegramBotService.SendMessage(req.Message)
			if err != nil {
				logger.Error("Fout bij verzenden Telegram bericht", "error", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Fout bij verzenden bericht: " + err.Error()})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Bericht succesvol verzonden"})
		})
		app.Get("/api/v1/telegrambot/commands", func(c *fiber.Ctx) error {
			authHeader := c.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
			}
			if serviceFactory.TelegramBotService == nil {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"success":  false,
					"message":  "Telegram bot service is niet geactiveerd",
					"commands": []interface{}{},
				})
			}
			commands := serviceFactory.TelegramBotService.GetCommands()
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success":  true,
				"message":  "Commando's succesvol opgehaald",
				"commands": commands,
			})
		})
		logger.Info("Telegram bot routes geregistreerd")
	}

	// Prometheus metrics endpoint
	app.Get("/metrics", func(c *fiber.Ctx) error {
		registry := prometheus.DefaultRegisterer.(*prometheus.Registry)
		handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/metrics", nil)
		handler.ServeHTTP(recorder, request)
		for k, v := range recorder.Header() {
			for _, val := range v {
				c.Set(k, val)
			}
		}
		return c.Status(recorder.Code).Send(recorder.Body.Bytes())
	})

	// Admin mail handler
	adminMailHandler := handlers.NewAdminMailHandler(serviceFactory.EmailService, serviceFactory.AuthService, serviceFactory.PermissionService, repoFactory.IncomingEmail)
	adminMailHandler.RegisterRoutes(app)

	// Chat handler
	chatHandler := handlers.NewChatHandler(serviceFactory.ChatService, serviceFactory.AuthService, serviceFactory.PermissionService, serviceFactory.ImageService, serviceFactory.Hub)
	chatHandler.RegisterRoutes(app)
	chatHandler.SetChannelHubCallback()

	// RBAC handlers (Permission, Role, User)
	permissionHandler := handlers.NewPermissionHandler(
		repoFactory.Permission,
		repoFactory.RBACRole,
		repoFactory.RolePermission,
		repoFactory.UserRole,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	permissionHandler.RegisterRoutes(app)

	userHandler := handlers.NewUserHandler(serviceFactory.AuthService, serviceFactory.PermissionService, repoFactory.UserRole, repoFactory.RBACRole)
	userHandler.RegisterRoutes(app)

	// Image handler
	imageHandler := handlers.NewImageHandler(serviceFactory.ImageService, serviceFactory.AuthService)
	imageHandler.RegisterRoutes(app)

	// --- CMS Handlers ---
	partnerHandler := handlers.NewPartnerHandler(
		repoFactory.Partner,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	partnerHandler.RegisterRoutes(app)

	radioRecordingHandler := handlers.NewRadioRecordingHandler(
		repoFactory.RadioRecording,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	radioRecordingHandler.RegisterRoutes(app)

	photoHandler := handlers.NewPhotoHandler(
		repoFactory.Photo,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	photoHandler.RegisterRoutes(app)

	albumHandler := handlers.NewAlbumHandler(
		repoFactory.Album,
		repoFactory.Photo,
		repoFactory.AlbumPhoto,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	albumHandler.RegisterRoutes(app)

	videoHandler := handlers.NewVideoHandler(
		repoFactory.Video,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	videoHandler.RegisterRoutes(app)

	sponsorHandler := handlers.NewSponsorHandler(
		repoFactory.Sponsor,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
		serviceFactory.ImageService,
	)
	sponsorHandler.RegisterRoutes(app)

	programScheduleHandler := handlers.NewProgramScheduleHandler(
		repoFactory.ProgramSchedule,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	programScheduleHandler.RegisterRoutes(app)

	socialEmbedHandler := handlers.NewSocialEmbedHandler(
		repoFactory.SocialEmbed,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	socialEmbedHandler.RegisterRoutes(app)

	socialLinkHandler := handlers.NewSocialLinkHandler(
		repoFactory.SocialLink,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	socialLinkHandler.RegisterRoutes(app)

	underConstructionHandler := handlers.NewUnderConstructionHandler(
		repoFactory.UnderConstruction,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	underConstructionHandler.RegisterRoutes(app)

	autoResponseHandler := handlers.NewAutoResponseHandler(
		repoFactory.AutoResponse,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	autoResponseHandler.RegisterRoutes(app)

	titleSectionHandler := handlers.NewTitleSectionHandler(
		repoFactory.TitleSection,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	titleSectionHandler.RegisterRoutes(app)

	gamificationHandler := handlers.NewGamificationHandler(
		serviceFactory.GamificationService,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	gamificationHandler.RegisterRoutes(app)

	// Public alias routes voor backwards compatibility met test endpoints
	// Deze routes redirecten naar de correcte handler endpoints
	api.Get("/title-sections", func(c *fiber.Ctx) error {
		return titleSectionHandler.GetTitleSection(c)
	})

	api.Get("/achievements", func(c *fiber.Ctx) error {
		return gamificationHandler.GetBadges(c) // Achievements zijn eigenlijk badges
	})

	// Notifications alias - wijst naar v1 endpoint
	api.Get("/notifications", func(c *fiber.Ctx) error {
		// Redirect to the actual v1 endpoint
		return c.Redirect("/api/v1/notifications", fiber.StatusMovedPermanently)
	})

	// Roles endpoint alias (vereist admin rechten via AdminPermissionMiddleware)
	api.Get("/roles",
		handlers.AuthMiddleware(serviceFactory.AuthService),
		handlers.AdminPermissionMiddleware(serviceFactory.PermissionService),
		permissionHandler.ListRoles,
	)

	// Permissions endpoint alias (vereist admin rechten via AdminPermissionMiddleware)
	api.Get("/permissions",
		handlers.AuthMiddleware(serviceFactory.AuthService),
		handlers.AdminPermissionMiddleware(serviceFactory.PermissionService),
		permissionHandler.ListPermissions,
	)

	// Event handler
	// GEWIJZIGD: Injecteer EventRegistrationRepo
	eventHandler := handlers.NewEventHandler(
		repoFactory.Event,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	eventHandler.RegisterRoutes(app)

	// Notulen handler
	notulenHandler := handlers.NewNotulenHandler(*serviceFactory.NotulenService, serviceFactory.AuthService, serviceFactory.PermissionService)
	notulenHandler.RegisterRoutes(app)

	// Notulen WebSocket handler
	notulenWsHandler := handlers.NewNotulenWebSocketHandler(serviceFactory.NotulenService.Hub(), serviceFactory.AuthService)
	notulenWsHandler.RegisterRoutes(app)
	logger.Info("Notulen WebSocket routes registered - /api/ws/notulen endpoint active")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 for web traffic
	}

	// Start server in een goroutine
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

	// Stop de Newsletter service
	if serviceFactory.NewsletterService != nil {
		serviceFactory.NewsletterService.Stop()
	}

	// Sluit rate limiter af
	if rateLimiter != nil {
		rateLimiter.Shutdown()
	}

	// Log laatste metrics
	serviceFactory.EmailMetrics.LogMetrics()

	// Sluit alle log writers
	logger.CloseWriters()

	// Graceful shutdown
} //
