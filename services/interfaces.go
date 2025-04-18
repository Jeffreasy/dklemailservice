package services

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// IEmailService definieert de interface voor email operaties
type IEmailService interface {
	// SendContactEmail verzendt een email voor een contactformulier
	SendContactEmail(ctx context.Context, contact *models.ContactFormulier) error

	// SendAanmeldingEmail verzendt een email voor een aanmelding
	SendAanmeldingEmail(ctx context.Context, aanmelding *models.Aanmelding) error
}

// RateLimiterService definieert de interface voor rate limiting
type RateLimiterService interface {
	// Allow controleert of een verzoek is toegestaan
	Allow(key string) bool

	// GetLimits haalt de huidige limieten op
	GetLimits() map[string]RateLimit

	// GetCurrentValues haalt de huidige waarden op
	GetCurrentValues() map[string]int
}

// EmailMetricsService definieert de interface voor email metrics
type EmailMetricsService interface {
	// RecordEmailSent registreert een verzonden email
	RecordEmailSent(emailType string)

	// RecordEmailError registreert een email fout
	RecordEmailError(emailType string, err error)

	// GetMetrics haalt de huidige metrics op
	GetMetrics() map[string]interface{}

	// LogMetrics logt de huidige metrics
	LogMetrics()
}

// AuthService definieert de interface voor authenticatie operaties
type AuthService interface {
	// Login authenticeert een gebruiker en geeft een JWT token terug
	Login(ctx context.Context, email, wachtwoord string) (string, error)

	// ValidateToken valideert een JWT token en geeft de gebruiker ID terug
	ValidateToken(token string) (string, error)

	// GetUserFromToken haalt de gebruiker op basis van een JWT token
	GetUserFromToken(ctx context.Context, token string) (*models.Gebruiker, error)

	// HashPassword genereert een hash voor een wachtwoord
	HashPassword(wachtwoord string) (string, error)

	// VerifyPassword verifieert een wachtwoord tegen een hash
	VerifyPassword(hash, wachtwoord string) bool

	// ResetPassword reset het wachtwoord van een gebruiker
	ResetPassword(ctx context.Context, email, nieuwWachtwoord string) error
}

type EmailSender interface {
	SendEmail(to, subject, body string) error
	// SendWFCEmail stuurt een Whisky for Charity e-mail met platte tekst body
	SendWFCEmail(to, subject, body string) error
	// SendTemplateEmail stuurt een e-mail met behulp van een template en variabelen
	SendTemplateEmail(recipient, subject, templateName string, templateData map[string]interface{}) error
}

type TemplateRenderer interface {
	RenderTemplate(name string, data interface{}) (string, error)
}

// EmailAutoFetcherInterface definieert de interface voor automatisch ophalen van emails
type EmailAutoFetcherInterface interface {
	// Start begint het periodiek ophalen van emails
	Start()

	// Stop stopt het periodiek ophalen van emails
	Stop()

	// IsRunning controleert of de auto fetcher actief is
	IsRunning() bool

	// GetLastRunTime geeft de laatste keer dat emails zijn opgehaald
	GetLastRunTime() time.Time
}

// NotificationService definieert de interface voor notificaties
type NotificationService interface {
	// SendNotification verstuurt een notificatie
	SendNotification(ctx context.Context, notification *models.Notification) error

	// CreateNotification maakt een nieuwe notificatie aan
	CreateNotification(ctx context.Context, notificationType models.NotificationType,
		priority models.NotificationPriority, title, message string) (*models.Notification, error)

	// GetNotification haalt een notificatie op basis van ID
	GetNotification(ctx context.Context, id string) (*models.Notification, error)

	// ListUnsentNotifications haalt alle niet verzonden notificaties op
	ListUnsentNotifications(ctx context.Context) ([]*models.Notification, error)

	// ProcessUnsentNotifications verwerkt alle niet verzonden notificaties
	ProcessUnsentNotifications(ctx context.Context) error

	// Start begint het periodiek verzenden van notificaties
	Start()

	// Stop stopt het periodiek verzenden van notificaties
	Stop()

	// IsRunning controleert of de service actief is
	IsRunning() bool
}
