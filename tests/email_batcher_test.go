package tests

import (
	"dklautomationgo/services"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmailBatcher(t *testing.T) {
	reg := NewTestRegistry()
	emailMetrics := services.NewEmailMetrics(time.Hour)
	prometheusMetrics := services.NewPrometheusMetricsWithRegistry(reg)
	rateLimiter := services.NewRateLimiter(prometheusMetrics)

	// Maak een mock SMTP client
	smtpMock := &mockSMTP{}
	emailService := services.NewEmailService(smtpMock, emailMetrics, rateLimiter, prometheusMetrics)

	t.Run("Basis batch verwerking", func(t *testing.T) {
		batcher := services.NewEmailBatcher(emailService, 3, 100*time.Millisecond)
		defer batcher.Shutdown()

		// Voeg emails toe aan de batch
		for i := 0; i < 3; i++ {
			err := emailService.SendEmail(fmt.Sprintf("test%d@example.com", i+1), "Test Subject", "Test Body")
			assert.NoError(t, err)
		}

		// Wacht op verwerking
		time.Sleep(200 * time.Millisecond)

		// Controleer of alle emails zijn verzonden in de juiste volgorde
		smtpMock.mutex.Lock()
		assert.True(t, smtpMock.SendCalled, "Send should have been called")
		assert.Equal(t, "test3@example.com", smtpMock.LastTo)
		smtpMock.mutex.Unlock()
	})

	t.Run("Batch timeout", func(t *testing.T) {
		// Reset mock
		smtpMock = &mockSMTP{}
		emailService = services.NewEmailService(smtpMock, emailMetrics, rateLimiter, prometheusMetrics)

		batcher := services.NewEmailBatcher(emailService, 5, 100*time.Millisecond)
		defer batcher.Shutdown()

		// Voeg 2 emails toe (minder dan batch size)
		for i := 0; i < 2; i++ {
			err := emailService.SendEmail(fmt.Sprintf("test%d@example.com", i+1), "Test Subject", "Test Body")
			assert.NoError(t, err)
		}

		// Wacht op timeout verwerking
		time.Sleep(200 * time.Millisecond)

		// Controleer of emails zijn verzonden
		smtpMock.mutex.Lock()
		assert.True(t, smtpMock.SendCalled, "Send should have been called on timeout")
		assert.Equal(t, "test2@example.com", smtpMock.LastTo)
		smtpMock.mutex.Unlock()
	})

	t.Run("Batch met fouten", func(t *testing.T) {
		mock := &mockSMTP{}
		rateLimiter := newMockRateLimiter()
		metrics := newMockPrometheusMetrics()

		// Configureer de rateLimiter om alle emails toe te staan
		rateLimiter.AddLimit("email_generic", 10, time.Hour, false)

		emailService := services.NewEmailService(mock, nil, rateLimiter, metrics)
		batcher := services.NewEmailBatcher(emailService, 2, time.Millisecond*200)
		defer batcher.Shutdown()

		// Configureer de mock om te falen
		mock.SendError = errors.New("SMTP failed")

		batchKey := "test-batch"
		recipients := []string{"test1@example.com", "test2@example.com"}

		for _, recipient := range recipients {
			batcher.AddToBatch(batchKey, recipient, "Test Email", "test-template", map[string]interface{}{
				"body": "Test Body",
			})
		}

		time.Sleep(time.Millisecond * 300)

		// Controleer dat Send is aangeroepen maar gefaald (error is set)
		mock.mutex.Lock()
		assert.True(t, mock.SendCalled, "Send should have been called")
		assert.Error(t, mock.SendError, "SendError should be set")
		assert.Equal(t, "test2@example.com", mock.LastTo)
		mock.mutex.Unlock()
	})

	t.Run("Graceful shutdown", func(t *testing.T) {
		// Reset mock
		smtpMock = &mockSMTP{}
		emailService = services.NewEmailService(smtpMock, emailMetrics, rateLimiter, prometheusMetrics)

		batcher := services.NewEmailBatcher(emailService, 3, 1*time.Second)

		// Voeg emails toe
		for i := 0; i < 2; i++ {
			err := emailService.SendEmail(fmt.Sprintf("test%d@example.com", i+1), "Test Subject", "Test Body")
			assert.NoError(t, err)
		}

		// Direct shutdown
		batcher.Shutdown()

		// Controleer of emails nog steeds zijn verwerkt
		smtpMock.mutex.Lock()
		assert.True(t, smtpMock.SendCalled, "Send should have been called on shutdown")
		assert.Equal(t, "test2@example.com", smtpMock.LastTo)
		smtpMock.mutex.Unlock()
	})
}
