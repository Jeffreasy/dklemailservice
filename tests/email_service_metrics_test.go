package tests

import (
	"dklautomationgo/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmailService_WithMetrics(t *testing.T) {
	// Maak een mock SMTP client
	mockSMTP := newMockSMTP()

	// Maak metrics instanties
	emailMetrics := services.NewEmailMetrics(time.Hour)
	mockPrometheus := newMockPrometheusMetrics()
	rateLimiter := newMockRateLimiter()

	// Maak de email service met alle benodigde dependencies
	emailService := services.NewEmailService(mockSMTP, emailMetrics, rateLimiter, mockPrometheus)

	// Voer de test uit
	err := emailService.SendEmail("test@example.com", "Test Subject", "Test Body")
	assert.NoError(t, err)

	// Wacht even om zeker te zijn dat de metrics zijn bijgewerkt
	time.Sleep(time.Millisecond * 100)

	// Controleer metrics
	assert.Equal(t, int64(1), emailMetrics.GetTotalEmails(), "Totaal aantal emails moet 1 zijn")
	assert.Equal(t, 100.0, emailMetrics.GetSuccessRate(), "Success rate moet 100% zijn")
	assert.Equal(t, 1, mockPrometheus.emailsSent, "Prometheus metrics moet 1 verzonden email registreren")
	assert.Equal(t, 0, mockPrometheus.emailsFailed, "Prometheus metrics moet 0 gefaalde emails registreren")

	// Controleer of de email is verzonden
	sentEmails := mockSMTP.GetSentEmails()
	assert.Equal(t, 1, len(sentEmails))
	if len(sentEmails) > 0 {
		assert.Equal(t, "test@example.com", sentEmails[0].To)
	}
}
