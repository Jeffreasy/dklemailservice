// Package services bevat alle services voor de applicatie
package services

import (
	"dklautomationgo/logger"
	"dklautomationgo/repository"
	"os"
	"strconv"
	"time"
)

// ServiceFactory bevat alle services
type ServiceFactory struct {
	EmailService *EmailService
	SMTPClient   SMTPClient
	RateLimiter  RateLimiterInterface
	EmailMetrics *EmailMetrics
	EmailBatcher *EmailBatcher
	AuthService  AuthService
}

// NewServiceFactory maakt een nieuwe service factory
func NewServiceFactory(repoFactory *repository.RepositoryFactory) *ServiceFactory {
	logger.Info("Initialiseren service factory")

	// Initialiseer Prometheus metrics
	prometheusMetrics := GetPrometheusMetrics()

	// Initialiseer rate limiter
	rateLimiter := createRateLimiter(prometheusMetrics)

	// Initialiseer email metrics
	emailMetrics := createEmailMetrics()

	// Initialiseer SMTP client
	smtpClient := createSMTPClient()

	// Initialiseer email service
	emailService := NewEmailService(smtpClient, emailMetrics, rateLimiter, prometheusMetrics)

	// Initialiseer email batcher
	emailBatcher := createEmailBatcher(emailService)

	// Initialiseer auth service
	authService := NewAuthService(repoFactory.Gebruiker)

	return &ServiceFactory{
		EmailService: emailService,
		SMTPClient:   smtpClient,
		RateLimiter:  rateLimiter,
		EmailMetrics: emailMetrics,
		EmailBatcher: emailBatcher,
		AuthService:  authService,
	}
}

// createRateLimiter maakt een nieuwe rate limiter
func createRateLimiter(prometheusMetrics *PrometheusMetrics) *RateLimiter {
	rateLimiter := NewRateLimiter(prometheusMetrics)

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

	return NewRealSMTPClient(host, port, username, password, from, regHost, regPort, regUsername, regPassword, regFrom)
}

// createEmailBatcher maakt een nieuwe email batcher
func createEmailBatcher(emailService *EmailService) *EmailBatcher {
	batchSize, _ := strconv.Atoi(getEnvWithDefault("EMAIL_BATCH_SIZE", "10"))
	batchWindow, _ := strconv.Atoi(getEnvWithDefault("EMAIL_BATCH_WINDOW", "300"))
	return NewEmailBatcher(emailService, batchSize, time.Duration(batchWindow)*time.Second)
}

// getEnvWithDefault haalt een omgevingsvariabele op met een standaardwaarde
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
