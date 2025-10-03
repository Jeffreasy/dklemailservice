package services

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// IEmailService definieert de interface voor email operaties
// Deze interface lijkt specifiek voor de EmailHandler te zijn en gebruikt Contact/Aanmelding modellen
// Het is mogelijk overbodig als EmailSender interface algemener is.
type IEmailService interface {
	// SendContactEmail verzendt een email voor een contactformulier
	// SendContactEmail(ctx context.Context, contact *models.ContactFormulier) error // Mogelijk vervangen door EmailSender?

	// SendAanmeldingEmail verzendt een email voor een aanmelding
	// SendAanmeldingEmail(ctx context.Context, aanmelding *models.Aanmelding) error // Mogelijk vervangen door EmailSender?
}

// RateLimiterService definieert de interface voor rate limiting
type RateLimiterService interface {
	// Allow controleert of een verzoek is toegestaan
	Allow(key string) bool

	// GetLimits haalt de huidige limieten op
	GetLimits() map[string]RateLimit

	// GetCurrentValues haalt de huidige waarden op
	GetCurrentValues() map[string]int

	// Shutdown cleans up the rate limiter resources
	Shutdown()
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

// EmailSender definieert de generieke interface voor het versturen van e-mails.
// Deze wordt gebruikt door AdminMailHandler.
type EmailSender interface {
	// SendEmail verstuurt een platte tekst/HTML e-mail, met optioneel 'From' adres.
	SendEmail(to, subject, body string, fromAddress ...string) error
	// SendTemplateEmail stuurt een e-mail met behulp van een template, met optioneel 'From' adres.
	SendTemplateEmail(recipient, subject, templateName string, templateData map[string]interface{}, fromAddress ...string) error

	// Methoden specifiek gebruikt door EmailHandler (Contact/Aanmelding).
	// Behoud originele signature als ze altijd de geconfigureerde afzender moeten gebruiken.
	SendContactEmail(data *models.ContactEmailData) error
	SendAanmeldingEmail(data *models.AanmeldingEmailData) error

	// Methode specifiek voor WFC
	SendWFCEmail(to, subject, body string) error
	SendWFCOrderEmail(data *models.WFCOrderEmailData) error
}

// SMTPClient definieert de laag-niveau interface voor SMTP interactie.
type SMTPClient interface {
	// Send stuurt een bericht via de standaard SMTP configuratie (SMTP_FROM als afzender).
	Send(msg *EmailMessage) error
	// SendRegistration stuurt een bericht via de registratie SMTP configuratie (REGISTRATION_SMTP_FROM als afzender).
	SendRegistration(msg *EmailMessage) error
	// SendWFC stuurt een bericht via de WFC SMTP configuratie (WFC_SMTP_FROM als afzender).
	SendWFC(msg *EmailMessage) error
	// SendWithFrom stuurt een bericht via de standaard SMTP configuratie, maar met een expliciet opgegeven 'From' adres.
	SendWithFrom(from string, msg *EmailMessage) error

	// Oudere helper methoden, mogelijk overbodig maken of aanpassen?
	// SendEmail(to, subject, body string, fromAddress ...string) error
	// SendWFCEmail(to, subject, body string) error
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
