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

	t.Run("SendRegistrationEmail Admin", func(t *testing.T) {
		participant := &models.Participant{
			ID:    "test-participant-id",
			Naam:  "Test Participant",
			Email: "testaanmelding@example.com",
		}
		registration := &models.EventRegistration{
			ID:            "test-reg-id",
			ParticipantID: participant.ID,
			TestMode:      true,
		}
		data := &models.RegistrationEmailData{
			Participant:  participant,
			Registration: registration,
			AdminEmail:   "admin@example.com",
			ToAdmin:      true,
		}
		err := emailService.SendRegistrationEmail(data)
		assert.NoError(t, err)
		mockSMTP.mutex.Lock()
		assert.True(t, mockSMTP.SendCalled, "Send should be called for admin registration")
		assert.Equal(t, "admin@example.com", mockSMTP.LastTo)
		mockSMTP.mutex.Unlock()
	})

	t.Run("SendRegistrationEmail User", func(t *testing.T) {
		participant := &models.Participant{
			ID:    "test-participant-id",
			Naam:  "Test Participant",
			Email: "user@example.com",
		}
		registration := &models.EventRegistration{
			ID:            "test-reg-id",
			ParticipantID: participant.ID,
			TestMode:      true,
		}
		data := &models.RegistrationEmailData{
			Participant:  participant,
			Registration: registration,
			ToAdmin:      false,
		}
		err := emailService.SendRegistrationEmail(data)
		assert.NoError(t, err)
		mockSMTP.mutex.Lock()
		assert.True(t, mockSMTP.SendRegCalled, "SendRegistration should be called for user registration")
		assert.Equal(t, "user@example.com", mockSMTP.LastTo)
		mockSMTP.mutex.Unlock()
	})
}
