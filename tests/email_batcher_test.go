package tests

import (
	"dklautomationgo/services"
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
	mockSMTP := newMockSMTP()
	emailService := services.NewEmailService(mockSMTP, emailMetrics, rateLimiter, prometheusMetrics)

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
		sentEmails := mockSMTP.GetSentEmails()
		assert.Equal(t, 3, len(sentEmails))
		if len(sentEmails) == 3 {
			assert.Equal(t, "test1@example.com", sentEmails[0].To)
			assert.Equal(t, "test2@example.com", sentEmails[1].To)
			assert.Equal(t, "test3@example.com", sentEmails[2].To)
		}
	})

	t.Run("Batch timeout", func(t *testing.T) {
		// Reset mock
		mockSMTP = newMockSMTP()
		emailService = services.NewEmailService(mockSMTP, emailMetrics, rateLimiter, prometheusMetrics)

		batcher := services.NewEmailBatcher(emailService, 5, 100*time.Millisecond)
		defer batcher.Shutdown()

		// Voeg 2 emails toe (minder dan batch size)
		for i := 0; i < 2; i++ {
			err := emailService.SendEmail(fmt.Sprintf("test%d@example.com", i+1), "Test Subject", "Test Body")
			assert.NoError(t, err)
		}

		// Wacht op timeout verwerking
		time.Sleep(200 * time.Millisecond)

		// Controleer of emails zijn verzonden ondanks incomplete batch
		sentEmails := mockSMTP.GetSentEmails()
		assert.Equal(t, 2, len(sentEmails))
		if len(sentEmails) == 2 {
			assert.Equal(t, "test1@example.com", sentEmails[0].To)
			assert.Equal(t, "test2@example.com", sentEmails[1].To)
		}
	})

	t.Run("Batch met fouten", func(t *testing.T) {
		mock := newMockSMTP()
		rateLimiter := newMockRateLimiter()
		metrics := newMockPrometheusMetrics()

		// Configureer de rateLimiter om alle emails toe te staan
		rateLimiter.AddLimit("email_generic", 10, time.Hour, false)

		emailService := services.NewEmailService(mock, nil, rateLimiter, metrics)
		batcher := services.NewEmailBatcher(emailService, 2, time.Millisecond*200)
		defer batcher.Shutdown()

		// Configureer de mock om alleen de eerste email te laten falen
		mock.SetFailFirst(true)

		batchKey := "test-batch"
		recipients := []string{"test1@example.com", "test2@example.com"}

		for _, recipient := range recipients {
			batcher.AddToBatch(batchKey, recipient, "Test Email", "test-template", map[string]interface{}{
				"body": "Test Body",
			})
		}

		time.Sleep(time.Millisecond * 300)
		sentEmails := mock.GetSentEmails()
		assert.Equal(t, 1, len(sentEmails), "Er moet precies één email succesvol verzonden zijn")
		if len(sentEmails) > 0 {
			assert.Equal(t, "test2@example.com", sentEmails[0].To, "De tweede email moet succesvol verzonden zijn")
		}
	})

	t.Run("Graceful shutdown", func(t *testing.T) {
		// Reset mock
		mockSMTP = newMockSMTP()
		emailService = services.NewEmailService(mockSMTP, emailMetrics, rateLimiter, prometheusMetrics)

		batcher := services.NewEmailBatcher(emailService, 3, 1*time.Second)

		// Voeg emails toe
		for i := 0; i < 2; i++ {
			err := emailService.SendEmail(fmt.Sprintf("test%d@example.com", i+1), "Test Subject", "Test Body")
			assert.NoError(t, err)
		}

		// Direct shutdown
		batcher.Shutdown()

		// Controleer of emails nog steeds zijn verwerkt
		sentEmails := mockSMTP.GetSentEmails()
		assert.Equal(t, 2, len(sentEmails))
		if len(sentEmails) == 2 {
			assert.Equal(t, "test1@example.com", sentEmails[0].To)
			assert.Equal(t, "test2@example.com", sentEmails[1].To)
		}
	})
}
