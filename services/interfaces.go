package services

import (
	"context"
	"dklautomationgo/models"
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
}

type TemplateRenderer interface {
	RenderTemplate(name string, data interface{}) (string, error)
}
