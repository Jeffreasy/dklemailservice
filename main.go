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

	_ "dklautomationgo/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title DKL Email Service API
// @version 1.0
// @description API voor de De Koninklijke Loop email service - inclusief registratie, deelnemersbeheer, CMS en meer
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.dekoninklijkeloop.nl
// @contact.email info@dekoninklijkeloop.nl

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer token voor geauthenticeerde endpoints

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

	// Laad en valideer logger configuratie
	loggerConfig := config.LoadLoggerConfig()
	if err := loggerConfig.ValidateLoggerConfig(); err != nil {
		// Kan nog niet loggen omdat logger nog niet is geïnitialiseerd
		fmt.Printf("Logger configuratie fout: %v\n", err)
		os.Exit(1)
	}

	// Configureer de logger
	if err := loggerConfig.SetupLogger(); err != nil {
		fmt.Printf("Logger setup fout: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Debug: Print alle omgevingsvariabelen alleen bij DEBUG logniveau
	if strings.ToUpper(loggerConfig.Level) == logger.DebugLevel {
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
				// Verberg wachtwoorden en secrets in logs
				if strings.Contains(env, "PASSWORD") || strings.Contains(env, "SECRET") || strings.Contains(env, "KEY") {
					logger.Debug("Omgevingsvariabele gevonden", "key", env, "value", "********")
				} else {
					logger.Debug("Omgevingsvariabele gevonden", "key", env, "value", value)
				}
			}
		}
	} else {
		logger.Info("Omgevingsvariabelen debug overgeslagen (alleen beschikbaar in DEBUG modus)")
	}

	logger.Info("🚀 DKL Email Service starting", "version", handlers.Version)

	// Controleer omgevingsvariabelen
	if err := ValidateEnv(); err != nil {
		logger.Fatal("Configuratiefout", "error", err)
	}

	// Initialiseer database
	dbConfig := config.LoadDatabaseConfig()

	logger.Info("📊 Database configured", "host", dbConfig.Host, "port", dbConfig.Port, "db", dbConfig.DBName)

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

	// Initialiseer service factory
	serviceFactory := services.NewServiceFactory(repoFactory)

	// Initialiseer steps service
	stepsService := services.NewStepsService(db, repoFactory.Participant, repoFactory.Distance)

	// ✨ NIEUWE: Initialize StepsHub voor WebSocket real-time updates
	stepsHub := services.NewStepsHub(stepsService, serviceFactory.GamificationService)

	// ✨ NIEUWE: Link hub to service voor broadcasts
	stepsService.SetStepsHub(stepsHub)

	// ✨ NIEUWE: Start hub in background goroutine
	go stepsHub.Run()
	logger.Info("⚡ StepsHub initialized with WebSocket support")

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
	emailHandler := handlers.NewEmailHandler(
		serviceFactory.EmailService,
		serviceFactory.NotificationService,
		repoFactory.Participant,
		repoFactory.EventRegistration,
		repoFactory.Event,
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

	// Initialiseer ParticipantHandler
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
		stepsService,
	)

	// V30: Initialiseer de PublicRegistrationHandler voor publieke registratie
	publicRegistrationHandler := handlers.NewPublicRegistrationHandler(
		repoFactory.Participant,
		repoFactory.EventRegistration,
		repoFactory.Event,
		repoFactory.Gebruiker,
		serviceFactory.EmailService,
		serviceFactory.NotificationService,
		serviceFactory.PermissionService,
		repoFactory.RBACRole,
		repoFactory.UserRole,
	)

	// Initialiseer steps handler
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

	// Voeg SecurityHeadersMiddleware toe voor beveiligingsheaders
	app.Use(handlers.SecurityHeadersMiddleware())

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
				{"path": "/api/auth/login", "method": "POST", "description": "User login"},
				{"path": "/swagger/*", "method": "GET", "description": "OpenAPI/Swagger documentation"},
				// ... (overige routes zijn gelijk gebleven voor leesbaarheid) ...
				{"path": "/metrics", "method": "GET", "description": "Prometheus metrics"},
			},
		})
	})

	// Swagger documentation route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// API routes group
	api := app.Group("/api")

	// Health check endpoint
	api.Get("/health", handlers.HealthHandler)

	// Email routes
	api.Post("/contact-email", emailHandler.HandleContactEmail)
	api.Post("/register", emailHandler.HandleRegistrationEmail)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handlers.RateLimitMiddleware(rateLimiter, "login"), authHandler.HandleLogin)
	auth.Post("/logout", authHandler.HandleLogout)
	auth.Post("/refresh", authHandler.HandleRefreshToken)

	// Public auth routes
	auth.Post("/forgot-password", handlers.RateLimitMiddleware(rateLimiter, "forgot_password"), authHandler.HandleForgotPassword)
	auth.Post("/reset-password-with-token", authHandler.HandleResetPasswordWithToken)
	auth.Post("/send-verification", handlers.RateLimitMiddleware(rateLimiter, "email_verification"), authHandler.HandleSendEmailVerification)
	auth.Post("/verify-email", authHandler.HandleVerifyEmail)
	auth.Post("/resend-verification", handlers.RateLimitMiddleware(rateLimiter, "resend_email_verification"), authHandler.HandleResendEmailVerification)

	// Beveiligde auth routes
	authProtected := auth.Group("/", handlers.AuthMiddleware(serviceFactory.AuthService))
	authProtected.Get("/profile", authHandler.HandleGetProfile)
	authProtected.Post("/reset-password", authHandler.HandleResetPassword)
	authProtected.Delete("/account", handlers.RateLimitMiddleware(rateLimiter, "delete_account"), authHandler.HandleDeleteAccount)

	// Session management routes
	authProtected.Get("/sessions", authHandler.HandleListSessions)
	authProtected.Delete("/sessions/:sessionId", authHandler.HandleRevokeSession)
	authProtected.Post("/sessions/revoke-others", authHandler.HandleRevokeOtherSessions)

	// Metrics endpoints
	api.Get("/metrics/email", metricsHandler.HandleGetEmailMetrics)
	api.Get("/metrics/rate-limits", metricsHandler.HandleGetRateLimits)

	// Registreer handlers
	contactHandler.RegisterRoutes(app)
	stepsHandler.RegisterRoutes(app) // Must be before participant
	participantHandler.RegisterRoutes(app)
	eventRegistrationHandler.RegisterRoutes(app)
	publicRegistrationHandler.RegisterRoutes(app)
	logger.Info("Public registration routes registered - /api/public/aanmelden endpoint active")

	// WebSocket handlers
	stepsWsHandler := handlers.NewStepsWebSocketHandler(stepsHub, serviceFactory.AuthService)
	stepsWsHandler.RegisterRoutes(app)
	logger.Info("WebSocket routes registered - /ws/steps endpoint active")

	// WebSocket stats endpoint (admin only)
	app.Get("/api/ws/stats",
		handlers.AuthMiddleware(serviceFactory.AuthService),
		handlers.PermissionMiddleware(serviceFactory.PermissionService, "admin", "read"),
		stepsWsHandler.GetStats,
	)

	// Overige handlers registreren
	newsletterHandler.RegisterRoutes(app)
	notificationHandler.RegisterRoutes(app)
	mailHandler.RegisterRoutes(app)
	handlers.RegisterWFCOrderRoutes(app, serviceFactory.EmailService)

	// Telegram Bot routes (indien actief)
	if serviceFactory.TelegramBotService != nil {
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

	// Admin/Chat/RBAC/CMS Handlers
	adminMailHandler := handlers.NewAdminMailHandler(serviceFactory.EmailService, serviceFactory.AuthService, serviceFactory.PermissionService, repoFactory.IncomingEmail)
	adminMailHandler.RegisterRoutes(app)

	chatHandler := handlers.NewChatHandler(serviceFactory.ChatService, serviceFactory.AuthService, serviceFactory.PermissionService, serviceFactory.ImageService, serviceFactory.Hub)
	chatHandler.RegisterRoutes(app)
	chatHandler.SetChannelHubCallback()

	permissionHandler := handlers.NewPermissionHandler(
		repoFactory.Permission,
		repoFactory.RBACRole,
		repoFactory.RolePermission,
		repoFactory.UserRole,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	permissionHandler.RegisterRoutes(app)

	dashboardHandler := handlers.NewDashboardHandler(serviceFactory.DashboardService, serviceFactory.AuthService, serviceFactory.PermissionService)
	dashboardHandler.RegisterRoutes(app)

	userHandler := handlers.NewUserHandler(serviceFactory.AuthService, serviceFactory.PermissionService, repoFactory.UserRole, repoFactory.RBACRole)
	userHandler.RegisterRoutes(app)

	imageHandler := handlers.NewImageHandler(serviceFactory.ImageService, serviceFactory.AuthService)
	imageHandler.RegisterRoutes(app)

	// CMS handlers
	handlers.NewPartnerHandler(repoFactory.Partner, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)
	handlers.NewRadioRecordingHandler(repoFactory.RadioRecording, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)
	handlers.NewPhotoHandler(repoFactory.Photo, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)
	handlers.NewAlbumHandler(repoFactory.Album, repoFactory.Photo, repoFactory.AlbumPhoto, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)
	handlers.NewVideoHandler(repoFactory.Video, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)
	handlers.NewSponsorHandler(repoFactory.Sponsor, serviceFactory.AuthService, serviceFactory.PermissionService, serviceFactory.ImageService).RegisterRoutes(app)
	handlers.NewProgramScheduleHandler(repoFactory.ProgramSchedule, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)
	handlers.NewSocialEmbedHandler(repoFactory.SocialEmbed, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)
	handlers.NewSocialLinkHandler(repoFactory.SocialLink, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)

	underConstructionHandler := handlers.NewUnderConstructionHandler(
		repoFactory.UnderConstruction,
		serviceFactory.AuthService,
		serviceFactory.PermissionService,
	)
	// Register public check route
	api.Get("/under-construction/active", underConstructionHandler.GetActiveUnderConstruction)
	api.Get("/under-construction", underConstructionHandler.GetActiveUnderConstruction)
	// Register admin routes
	underConstructionHandler.RegisterRoutes(app)

	autoResponseHandler := handlers.NewAutoResponseHandler(repoFactory.AutoResponse, serviceFactory.AuthService, serviceFactory.PermissionService)
	autoResponseHandler.RegisterRoutes(app)

	titleSectionHandler := handlers.NewTitleSectionHandler(repoFactory.TitleSection, serviceFactory.AuthService, serviceFactory.PermissionService)
	titleSectionHandler.RegisterRoutes(app)

	gamificationHandler := handlers.NewGamificationHandler(serviceFactory.GamificationService, serviceFactory.AuthService, serviceFactory.PermissionService)
	gamificationHandler.RegisterRoutes(app)

	// Legacy/Alias routes
	api.Get("/title-sections", func(c *fiber.Ctx) error { return titleSectionHandler.GetTitleSection(c) })
	api.Get("/achievements", func(c *fiber.Ctx) error { return gamificationHandler.GetBadges(c) })
	api.Get("/notifications", func(c *fiber.Ctx) error { return c.Redirect("/api/v1/notifications", fiber.StatusMovedPermanently) })

	api.Get("/roles",
		handlers.AuthMiddleware(serviceFactory.AuthService),
		handlers.AdminPermissionMiddleware(serviceFactory.PermissionService),
		permissionHandler.ListRoles,
	)
	api.Get("/permissions",
		handlers.AuthMiddleware(serviceFactory.AuthService),
		handlers.AdminPermissionMiddleware(serviceFactory.PermissionService),
		permissionHandler.ListPermissions,
	)

	handlers.NewEventHandler(repoFactory.Event, serviceFactory.AuthService, serviceFactory.PermissionService).RegisterRoutes(app)

	notulenHandler := handlers.NewNotulenHandler(*serviceFactory.NotulenService, serviceFactory.AuthService, serviceFactory.PermissionService)
	notulenHandler.RegisterRoutes(app)
	handlers.NewNotulenWebSocketHandler(serviceFactory.NotulenService.Hub(), serviceFactory.AuthService).RegisterRoutes(app)
	logger.Info("Notulen WebSocket routes registered")

	// Alias routes for backwards compatibility (after all handlers are declared)
	// Auto response alias routes for backwards compatibility
	app.Get("/mail/autoresponse", func(c *fiber.Ctx) error {
		return autoResponseHandler.ListAutoResponses(c)
	})
	app.Get("/mail/autoresponse/:id", func(c *fiber.Ctx) error {
		return autoResponseHandler.GetAutoResponse(c)
	})
	app.Post("/mail/autoresponse", func(c *fiber.Ctx) error {
		return autoResponseHandler.CreateAutoResponse(c)
	})
	app.Put("/mail/autoresponse/:id", func(c *fiber.Ctx) error {
		return autoResponseHandler.UpdateAutoResponse(c)
	})
	app.Delete("/mail/autoresponse/:id", func(c *fiber.Ctx) error {
		return autoResponseHandler.DeleteAutoResponse(c)
	})

	// Mail unprocessed alias route for backwards compatibility
	app.Get("/mail/unprocessed", func(c *fiber.Ctx) error {
		return mailHandler.ListUnprocessedEmails(c)
	})

	// Mail account alias route for backwards compatibility
	app.Get("/mail/account/:type", func(c *fiber.Ctx) error {
		return mailHandler.ListEmailsByAccountType(c)
	})

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint niet gevonden",
		})
	})

	// =========================================================================
	// GRACEFUL SHUTDOWN IMPLEMENTATIE
	// =========================================================================

	// 1. Maak een kanaal voor OS signalen
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Bepaal poort
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Start de server in een aparte goroutine
	go func() {
		logger.Info("Server wordt gestart...", "port", port)
		if err := app.Listen(":" + port); err != nil {
			logger.Fatal("Server kon niet starten", "error", err)
		}
	}()

	// 3. Blokkeer tot signaal
	<-quit
	logger.Info("Shutdown signaal ontvangen, bezig met afsluiten...")

	// 4. Cleanup Sequentie

	// A. Stop Fiber server
	if err := app.Shutdown(); err != nil {
		logger.Error("Fout bij afsluiten server", "error", err)
	} else {
		logger.Info("HTTP server succesvol gestopt")
	}

	// B. Stop services
	if serviceFactory.EmailBatcher != nil {
		serviceFactory.EmailBatcher.Shutdown()
	}

	if serviceFactory.EmailAutoFetcher != nil && serviceFactory.EmailAutoFetcher.IsRunning() {
		logger.Info("Email auto fetcher stoppen...")
		serviceFactory.EmailAutoFetcher.Stop()
		logger.Info("Email auto fetcher gestopt")
	}

	if serviceFactory.NewsletterService != nil {
		serviceFactory.NewsletterService.Stop()
	}

	if rateLimiter != nil {
		rateLimiter.Shutdown()
	}

	// C. Sluit Database
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Kon onderliggende SQL DB niet ophalen voor sluiten", "error", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			logger.Error("Fout bij sluiten database verbinding", "error", err)
		} else {
			logger.Info("Database verbinding gesloten")
		}
	}

	serviceFactory.EmailMetrics.LogMetrics()
	logger.CloseWriters()

	logger.Info("Applicatie volledig afgesloten. Tot ziens!")
}
