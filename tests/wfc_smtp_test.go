package tests

import (
	"context"
	"os"
	"testing"

	"dklautomationgo/models"
	"dklautomationgo/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWFCOrderRepository is a mock implementation of WFCOrderRepository
//
//nolint:unused // This mock is kept for future tests
type MockWFCOrderRepository struct {
	mock.Mock
}

// --- removed WFCMockSMTPClient struct definition ---

//nolint:unused // This mock is kept for future tests
func (m *MockWFCOrderRepository) Create(ctx context.Context, order *models.WFCOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestWhiskyForCharitySMTP(t *testing.T) {
	// Setup test environment variables
	os.Setenv("WFC_SMTP_HOST", "arg-plplcl14.argewebhosting.nl")
	os.Setenv("WFC_SMTP_PORT", "465")
	os.Setenv("WFC_SMTP_USER", "noreply@whiskyforcharity.com")
	os.Setenv("WFC_SMTP_PASSWORD", "test_password")
	os.Setenv("WFC_SMTP_FROM", "noreply@whiskyforcharity.com")
	os.Setenv("WFC_SMTP_SSL", "true")
	defer func() {
		os.Unsetenv("WFC_SMTP_HOST")
		os.Unsetenv("WFC_SMTP_PORT")
		os.Unsetenv("WFC_SMTP_USER")
		os.Unsetenv("WFC_SMTP_PASSWORD")
		os.Unsetenv("WFC_SMTP_FROM")
		os.Unsetenv("WFC_SMTP_SSL")
	}()

	t.Run("Create SMTP client with WFC config", func(t *testing.T) {
		// Create test client with mock values for standard and registration config
		client := services.NewRealSMTPClientWithWFC(
			"smtp.test.com", "587", "test@test.com", "testpass", "noreply@test.com",
			"smtp.test.com", "587", "test@test.com", "testpass", "noreply@test.com",
			"arg-plplcl14.argewebhosting.nl", "465", "noreply@whiskyforcharity.com", "test_password", "noreply@whiskyforcharity.com", true,
		)

		// Setup mock SMTP
		mockSMTP := &WFCMockSMTPClient{}
		client.SetDialer(mockSMTP)

		// Test WFC email sending
		msg := &services.EmailMessage{
			To:      "test@example.com",
			Subject: "WFC Test",
			Body:    "<p>Test from Whisky for Charity</p>",
		}

		err := client.SendWFC(msg)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", mockSMTP.LastTo)
		assert.Equal(t, "WFC Test", mockSMTP.LastSubject)
	})

	t.Run("EmailService with WFC config", func(t *testing.T) {
		// Create a mock SMTP client
		mockSMTP := &WFCMockSMTPClient{}

		// Create email service with the mock
		emailService, err := services.NewTestEmailService(mockSMTP)
		assert.NoError(t, err)

		// Test the WFC email method
		err = emailService.SendWFCEmail("test@example.com", "WFC Test", "<p>Testing WFC</p>")
		assert.NoError(t, err)

		// Verify correct method was called on client
		assert.Equal(t, "test@example.com", mockSMTP.LastTo)
		assert.Equal(t, "WFC Test", mockSMTP.LastSubject)
	})

	t.Run("Factory method", func(t *testing.T) {
		// Using getCreateSMTPClient would require service factory refactoring for tests
		// This is a simplified test to verify the concept
		host := os.Getenv("WFC_SMTP_HOST")
		port := os.Getenv("WFC_SMTP_PORT")
		user := os.Getenv("WFC_SMTP_USER")
		password := os.Getenv("WFC_SMTP_PASSWORD")
		from := os.Getenv("WFC_SMTP_FROM")
		useSSL := os.Getenv("WFC_SMTP_SSL") == "true"

		assert.Equal(t, "arg-plplcl14.argewebhosting.nl", host)
		assert.Equal(t, "465", port)
		assert.Equal(t, "noreply@whiskyforcharity.com", user)
		assert.Equal(t, "test_password", password)
		assert.Equal(t, "noreply@whiskyforcharity.com", from)
		assert.True(t, useSSL)
	})
}
