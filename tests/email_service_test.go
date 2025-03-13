package tests

import (
	"dklautomationgo/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmailService_SendEmail(t *testing.T) {
	// Maak een nieuwe registry voor deze test
	reg := NewTestRegistry()

	// Maak een mock SMTP client
	mockSMTP := &mockSMTP{}

	// Maak metrics instanties
	emailMetrics := services.NewEmailMetrics(time.Hour)
	prometheusMetrics := services.NewPrometheusMetricsWithRegistry(reg)
	rateLimiter := services.NewRateLimiter(prometheusMetrics)

	// Maak de email service met alle benodigde dependencies
	emailService := services.NewEmailService(mockSMTP, emailMetrics, rateLimiter, prometheusMetrics)

	// Test cases
	t.Run("Succesvolle verzending", func(t *testing.T) {
		err := emailService.SendEmail("test@example.com", "Test", "Test body")
		assert.NoError(t, err)
		sentEmails := mockSMTP.GetSentEmails()
		assert.Equal(t, 1, len(sentEmails))
		if len(sentEmails) > 0 {
			assert.Equal(t, "test@example.com", sentEmails[0].To)
		}
	})

	t.Run("SMTP permanente fout", func(t *testing.T) {
		mockSMTP.SetShouldFail(true)
		err := emailService.SendEmail("test@example.com", "Test", "Test body")
		assert.Error(t, err)
		mockSMTP.SetShouldFail(false)
	})

	t.Run("Rate limit overschreden", func(t *testing.T) {
		// Voeg een strenge rate limit toe
		rateLimiter.AddLimit("email_generic", 1, time.Hour, false)

		// Eerste email moet lukken
		err := emailService.SendEmail("test@example.com", "Test", "Test body")
		assert.NoError(t, err)

		// Tweede email moet falen vanwege rate limit
		err = emailService.SendEmail("test@example.com", "Test", "Test body")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit")
	})
}
