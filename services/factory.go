// Package services bevat alle services voor de applicatie
package services

import (
	"dklautomationgo/config"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// ServiceFactory bevat alle services
type ServiceFactory struct {
	EmailService        *EmailService
	SMTPClient          SMTPClient
	RateLimiter         RateLimiterInterface
	EmailMetrics        *EmailMetrics
	EmailBatcher        *EmailBatcher
	AuthService         AuthService
	EmailAutoFetcher    EmailAutoFetcherInterface
	NotificationService NotificationService
	TelegramBotService  *TelegramBotService
	ChatService         ChatService
	Hub                 *Hub
	NewsletterService   *NewsletterService
	NewsletterSender    *NewsletterSender
	PermissionService   PermissionService
	ImageService        *ImageService
	RedisClient         *redis.Client
}

// GetRateLimiter retourneert de RateLimiter als het concrete type
// Dit helpt om onveilige type assertions in de code te vermijden
func (sf *ServiceFactory) GetRateLimiter() *RateLimiter {
	// Veilige type assertion met error checking
	rateLimiter, ok := sf.RateLimiter.(*RateLimiter)
	if !ok {
		logger.Fatal("Kon RateLimiter niet casten naar juiste type")
		return nil // Voor compilatie, wordt nooit bereikt na Fatal
	}
	return rateLimiter
}

// NewServiceFactory maakt een nieuwe service factory
func NewServiceFactory(repoFactory *repository.Repository) *ServiceFactory {
	logger.Info("Initialiseren service factory")

	// Initialiseer Redis client
	redisConfig := config.LoadRedisConfig()
	redisClient := config.NewRedisClient(redisConfig)
	if redisClient != nil {
		logger.Info("Redis client geïnitialiseerd")
	} else {
		logger.Warn("Redis client kon niet worden geïnitialiseerd, sommige features werken mogelijk niet")
	}

	// Initialiseer Prometheus metrics
	prometheusMetrics := GetPrometheusMetrics()

	// Initialiseer rate limiter met Redis ondersteuning
	rateLimiter := createRateLimiter(prometheusMetrics, redisClient)

	// Initialiseer email metrics
	emailMetrics := createEmailMetrics()

	// Initialiseer SMTP client
	smtpClient := createSMTPClient()

	// Initialiseer email service
	emailService := NewEmailService(smtpClient, emailMetrics, rateLimiter, prometheusMetrics)

	// Initialiseer email batcher
	emailBatcher := createEmailBatcher(emailService)

	// Initialiseer auth service
	authService := NewAuthService(repoFactory.Gebruiker, repoFactory.RefreshToken)

	// Initialiseer permission service met Redis caching
	permissionService := NewPermissionServiceWithRedis(
		repoFactory.RBACRole,
		repoFactory.Permission,
		repoFactory.RolePermission,
		repoFactory.UserRole,
		redisClient,
	)

	// Initialiseer notification service
	notificationService := createNotificationService(repoFactory.Notification)

	// Initialiseer telegram bot service
	telegramBotService := createTelegramBotService(repoFactory.Contact, repoFactory.Aanmelding)

	chatService := NewChatService(repoFactory.ChatChannel, repoFactory.ChatChannelParticipant, repoFactory.ChatMessage, repoFactory.ChatMessageReaction, repoFactory.ChatUserPresence)

	hub := NewHub(chatService)
	go hub.Run()

	// Maak een EmailAutoFetcher aan
	// Nog niet geinitialiseerd omdat MailFetcher buiten de ServiceFactory wordt aangemaakt in main.go

	// Newsletter components
	fetcher := NewNewsletterFetcher()
	processor := NewNewsletterProcessor()
	formatter := NewNewsletterFormatter(emailService)
	sender := NewNewsletterSender(emailService, emailBatcher, repoFactory.Gebruiker, repoFactory.Newsletter, notificationService)
	var newsletterSvc *NewsletterService
	if getEnvWithDefault("ENABLE_NEWSLETTER", "false") == "true" {
		newsletterSvc = NewNewsletterService(fetcher, processor, formatter, sender)
	}

	// Initialize ImageService if Cloudinary is configured
	var imageService *ImageService
	if cloudinaryConfig := config.LoadCloudinaryConfig(); cloudinaryConfig != nil {
		var err error
		imageService, err = NewImageService(cloudinaryConfig, repoFactory.UploadedImage)
		if err != nil {
			logger.Warn("Failed to initialize Cloudinary service", "error", err)
		} else {
			logger.Info("Cloudinary ImageService initialized")
		}
	}

	return &ServiceFactory{
		EmailService:        emailService,
		SMTPClient:          smtpClient,
		RateLimiter:         rateLimiter,
		EmailMetrics:        emailMetrics,
		EmailBatcher:        emailBatcher,
		AuthService:         authService,
		EmailAutoFetcher:    nil, // Dit wordt later in main.go ingesteld
		NotificationService: notificationService,
		TelegramBotService:  telegramBotService,
		ChatService:         chatService,
		Hub:                 hub,
		NewsletterService:   newsletterSvc,
		NewsletterSender:    sender,
		PermissionService:   permissionService,
		ImageService:        imageService,
		RedisClient:         redisClient,
	}
}

// createRateLimiter maakt een nieuwe rate limiter
func createRateLimiter(prometheusMetrics *PrometheusMetrics, redisClient *redis.Client) *RateLimiter {
	rateLimiter := NewRateLimiterWithRedis(prometheusMetrics, redisClient)

	// Configureer limieten uit omgevingsvariabelen
	contactLimitCount, _ := strconv.Atoi(getEnvWithDefault("CONTACT_LIMIT_COUNT", "5"))
	contactLimitPeriod, _ := strconv.Atoi(getEnvWithDefault("CONTACT_LIMIT_PERIOD", "3600"))
	contactLimitPerIP := getEnvWithDefault("CONTACT_LIMIT_PER_IP", "true") == "true"

	aanmeldingLimitCount, _ := strconv.Atoi(getEnvWithDefault("AANMELDING_LIMIT_COUNT", "3"))
	aanmeldingLimitPeriod, _ := strconv.Atoi(getEnvWithDefault("AANMELDING_LIMIT_PERIOD", "86400"))
	aanmeldingLimitPerIP := getEnvWithDefault("AANMELDING_LIMIT_PER_IP", "true") == "true"

	// Login rate limiting
	loginLimitCount, _ := strconv.Atoi(getEnvWithDefault("LOGIN_LIMIT_COUNT", "5"))
	loginLimitPeriod, _ := strconv.Atoi(getEnvWithDefault("LOGIN_LIMIT_PERIOD", "300"))
	loginLimitPerIP := getEnvWithDefault("LOGIN_LIMIT_PER_IP", "true") == "true"

	// Voeg limieten toe
	rateLimiter.AddLimit("contact", contactLimitCount, time.Duration(contactLimitPeriod)*time.Second, contactLimitPerIP)
	rateLimiter.AddLimit("aanmelding", aanmeldingLimitCount, time.Duration(aanmeldingLimitPeriod)*time.Second, aanmeldingLimitPerIP)
	rateLimiter.AddLimit("login", loginLimitCount, time.Duration(loginLimitPeriod)*time.Second, loginLimitPerIP)

	return rateLimiter
}

// createEmailMetrics maakt een nieuwe email metrics tracker
func createEmailMetrics() *EmailMetrics {
	resetInterval, _ := strconv.Atoi(getEnvWithDefault("EMAIL_METRICS_RESET_INTERVAL", "86400"))
	return NewEmailMetrics(time.Duration(resetInterval) * time.Second)
}

// createSMTPClient maakt een nieuwe SMTP client
func createSMTPClient() SMTPClient {
	// Standaard SMTP configuratie
	host := getEnvWithDefault("SMTP_HOST", "smtp.gmail.com")
	port := getEnvWithDefault("SMTP_PORT", "587")
	username := getEnvWithDefault("SMTP_USERNAME", "info@dekoninklijkeloop.nl")
	password := getEnvWithDefault("SMTP_PASSWORD", "")
	from := getEnvWithDefault("SMTP_FROM", "info@dekoninklijkeloop.nl")

	// Registratie SMTP configuratie (kan anders zijn)
	regHost := getEnvWithDefault("REG_SMTP_HOST", host)
	regPort := getEnvWithDefault("REG_SMTP_PORT", port)
	regUsername := getEnvWithDefault("REG_SMTP_USERNAME", username)
	regPassword := getEnvWithDefault("REG_SMTP_PASSWORD", password)
	regFrom := getEnvWithDefault("REG_SMTP_FROM", from)

	// Whisky for Charity SMTP configuratie
	wfcHost := getEnvWithDefault("WFC_SMTP_HOST", "")
	wfcPort := getEnvWithDefault("WFC_SMTP_PORT", "465")
	wfcUsername := getEnvWithDefault("WFC_SMTP_USER", "")
	wfcPassword := getEnvWithDefault("WFC_SMTP_PASSWORD", "")
	wfcFrom := getEnvWithDefault("WFC_SMTP_FROM", "")
	wfcUseSSL := getEnvWithDefault("WFC_SMTP_SSL", "true") == "true"

	// Als WFC configuratie aanwezig is, gebruik die
	if wfcHost != "" && wfcUsername != "" {
		return NewRealSMTPClientWithWFC(
			host, port, username, password, from,
			regHost, regPort, regUsername, regPassword, regFrom,
			wfcHost, wfcPort, wfcUsername, wfcPassword, wfcFrom, wfcUseSSL)
	}

	// Anders, gebruik standaard client
	return NewRealSMTPClient(host, port, username, password, from, regHost, regPort, regUsername, regPassword, regFrom)
}

// createEmailBatcher maakt een nieuwe email batcher
func createEmailBatcher(emailService *EmailService) *EmailBatcher {
	batchSize, _ := strconv.Atoi(getEnvWithDefault("EMAIL_BATCH_SIZE", "10"))
	batchWindow, _ := strconv.Atoi(getEnvWithDefault("EMAIL_BATCH_WINDOW", "300"))
	return NewEmailBatcher(emailService, batchSize, time.Duration(batchWindow)*time.Second)
}

// createNotificationService maakt een nieuwe notification service
func createNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	// Check of notificaties zijn ingeschakeld
	enabled := getEnvWithDefault("ENABLE_NOTIFICATIONS", "false") == "true"
	if !enabled {
		logger.Info("Notificaties zijn uitgeschakeld")
		return nil
	}

	// Haal Telegram configuratie op
	botToken := getEnvWithDefault("TELEGRAM_BOT_TOKEN", "")
	chatID := getEnvWithDefault("TELEGRAM_CHAT_ID", "")

	// Als Telegram niet geconfigureerd is, return nil
	if botToken == "" || chatID == "" {
		logger.Warn("Telegram configuratie ontbreekt, notificaties worden niet verzonden",
			"bot_token_provided", botToken != "",
			"chat_id_provided", chatID != "")
		return nil
	}

	// Maak een nieuwe Telegram client
	telegramClient := NewTelegramClient(botToken, chatID)

	// Parseer throttle duration
	throttleDurationStr := getEnvWithDefault("NOTIFICATION_THROTTLE", "15m")
	throttleDuration, err := time.ParseDuration(throttleDurationStr)
	if err != nil {
		logger.Warn("Ongeldige throttle duur, gebruik standaard 15 minuten",
			"duration", throttleDurationStr,
			"error", err)
		throttleDuration = 15 * time.Minute
	}

	// Parseer minimale prioriteit
	minPriorityStr := getEnvWithDefault("NOTIFICATION_MIN_PRIORITY", "medium")
	var minPriority models.NotificationPriority
	switch minPriorityStr {
	case "low":
		minPriority = models.NotificationPriorityLow
	case "medium":
		minPriority = models.NotificationPriorityMedium
	case "high":
		minPriority = models.NotificationPriorityHigh
	case "critical":
		minPriority = models.NotificationPriorityCritical
	default:
		minPriority = models.NotificationPriorityMedium
	}

	// Maak een nieuwe notification service
	notificationService := NewNotificationService(
		notificationRepo,
		telegramClient,
		throttleDuration,
		minPriority,
	)

	// Start de notification service
	go notificationService.Start()

	logger.Info("Notificatie service geïnitialiseerd",
		"throttle", throttleDuration.String(),
		"min_priority", minPriority)

	return notificationService
}

// createTelegramBotService maakt een nieuwe Telegram bot service
func createTelegramBotService(contactRepo repository.ContactRepository, aanmeldingRepo repository.AanmeldingRepository) *TelegramBotService {
	// Check of bot enabled is in omgevingsvariabelen
	enabled := getEnvWithDefault("ENABLE_TELEGRAM_BOT", "false") == "true"
	if !enabled {
		logger.Info("Telegram bot is uitgeschakeld")
		return nil
	}

	// Maak een nieuwe Telegram bot service
	telegramBotService := NewTelegramBotService(contactRepo, aanmeldingRepo)

	// Als de service succesvol is aangemaakt, start polling
	if telegramBotService != nil {
		logger.Info("Telegram bot service geïnitialiseerd")
		// Polling temporarily disabled to avoid conflicts on Render
		// telegramBotService.StartPolling()
	}

	return telegramBotService
}

// getEnvWithDefault haalt een omgevingsvariabele op met een standaardwaarde
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
