package tests

import (
	"dklautomationgo/models"
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

		// Verify email was sent using mock fields
		mockSMTP.mutex.Lock()
		assert.True(t, mockSMTP.SendCalled, "Send should have been called")
		assert.Equal(t, "test@example.com", mockSMTP.LastTo)
		assert.Equal(t, "Test", mockSMTP.LastSubject)
		mockSMTP.mutex.Unlock()
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

	t.Run("SendContactEmail Admin", func(t *testing.T) {
		data := &models.ContactEmailData{
			Contact: &models.ContactFormulier{
				Naam:  "Test Contact",
				Email: "testcontact@example.com",
			},
			ToAdmin:    true,
			AdminEmail: "admin@example.com",
		}
		err := emailService.SendContactEmail(data)
		assert.NoError(t, err)
		mockSMTP.mutex.Lock()
		assert.True(t, mockSMTP.SendCalled, "Send should be called for admin contact")
		assert.Equal(t, "admin@example.com", mockSMTP.LastTo)
		mockSMTP.mutex.Unlock()
	})

	t.Run("SendContactEmail User", func(t *testing.T) {
		data := &models.ContactEmailData{
			Contact: &models.ContactFormulier{
				Naam:  "Test Contact",
				Email: "user@example.com",
			},
			ToAdmin: false,
		}
		err := emailService.SendContactEmail(data)
		assert.NoError(t, err)
		mockSMTP.mutex.Lock()
		assert.True(t, mockSMTP.SendCalled, "Send should be called for user contact confirmation")
		assert.Equal(t, "user@example.com", mockSMTP.LastTo)
		mockSMTP.mutex.Unlock()
	})

	t.Run("SendAanmeldingEmail Admin", func(t *testing.T) {
		data := &models.AanmeldingEmailData{
			Aanmelding: &models.AanmeldingFormulier{
				Naam:  "Test Aanmelding",
				Email: "testaanmelding@example.com",
			},
			ToAdmin:    true,
			AdminEmail: "admin@example.com",
		}
		err := emailService.SendAanmeldingEmail(data)
		assert.NoError(t, err)
		mockSMTP.mutex.Lock()
		assert.True(t, mockSMTP.SendCalled, "Send should be called for admin aanmelding")
		assert.Equal(t, "admin@example.com", mockSMTP.LastTo)
		mockSMTP.mutex.Unlock()
	})

	t.Run("SendAanmeldingEmail User", func(t *testing.T) {
		data := &models.AanmeldingEmailData{
			Aanmelding: &models.AanmeldingFormulier{
				Naam:  "Test Aanmelding",
				Email: "user@example.com",
			},
			ToAdmin: false,
		}
		err := emailService.SendAanmeldingEmail(data)
		assert.NoError(t, err)
		mockSMTP.mutex.Lock()
		assert.True(t, mockSMTP.SendRegCalled, "SendRegistration should be called for user aanmelding")
		assert.Equal(t, "user@example.com", mockSMTP.LastTo)
		mockSMTP.mutex.Unlock()
	})
}
