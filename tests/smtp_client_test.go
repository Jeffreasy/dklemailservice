package tests

import (
	"os"
	"testing"

	"dklautomationgo/services"

	"github.com/stretchr/testify/assert"
)

// createTestClientWithConfig maakt een test SMTP client met aangepaste configuratie
func createTestClientWithConfig(host, port, user, pass, from string) *services.RealSMTPClient {
	// Voor tests gebruiken we dezelfde configuratie voor zowel standaard als registratie
	return services.NewRealSMTPClient(
		host, // host
		port, // port
		user, // user
		pass, // password
		from, // from
		host, // regHost
		port, // regPort
		user, // regUser
		pass, // regPassword
		from, // regFrom
	)
}

func createTestClient() *services.RealSMTPClient {
	return createTestClientWithConfig(
		"smtp.test.com",    // host
		"587",              // port
		"test@test.com",    // user
		"testpass",         // password
		"noreply@test.com", // from
	)
}

func TestSMTPClient(t *testing.T) {
	client := createTestClient()

	t.Run("Send email with default config", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client.SetDialer(mockSMTP)

		msg := &services.EmailMessage{
			To:      "recipient@test.com",
			Subject: "Test Subject",
			Body:    "Test Body",
		}

		err := client.Send(msg)
		assert.NoError(t, err)
		assert.Equal(t, msg, mockSMTP.GetLastEmail())
	})

	t.Run("Send email with registration config", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client.SetDialer(mockSMTP)

		msg := &services.EmailMessage{
			To:      "recipient@test.com",
			Subject: "Registration Test",
			Body:    "Test Body",
		}

		err := client.SendRegistration(msg)
		assert.NoError(t, err)
		assert.Equal(t, msg, mockSMTP.GetLastEmail())
	})

	t.Run("Send email with helper function", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client.SetDialer(mockSMTP)
		err := client.SendEmail("recipient@test.com", "Helper Test", "Test Body")
		assert.NoError(t, err)
	})

	t.Run("Verbindingsfout", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		mockSMTP.SetShouldFail(true)
		client.SetDialer(mockSMTP)

		err := client.SendEmail("test@example.com", "Test Subject", "Test Body")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mock dialer error")
	})

	t.Run("Ongeldige email configuratie", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client.SetDialer(mockSMTP)

		err := client.SendEmail("", "Test Subject", "Test Body")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid recipient")
	})

	t.Run("HTML email verzending", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client.SetDialer(mockSMTP)

		htmlBody := "<html><body><h1>Test</h1></body></html>"
		err := client.SendEmail("test@example.com", "HTML Test", htmlBody)
		assert.NoError(t, err)
	})
}

func TestSMTPClientBatch(t *testing.T) {
	// Setup test SMTP configuratie
	os.Setenv("SMTP_HOST", "test.smtp.com")
	os.Setenv("SMTP_USER", "test")
	os.Setenv("SMTP_PASSWORD", "test")
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_USER")
		os.Unsetenv("SMTP_PASSWORD")
	}()

	t.Run("Batch email verzending", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client := createTestClient()
		client.SetDialer(mockSMTP)

		messages := []*services.EmailMessage{
			{To: "test1@example.com", Subject: "Test 1", Body: "Body 1"},
			{To: "test2@example.com", Subject: "Test 2", Body: "Body 2"},
			{To: "test3@example.com", Subject: "Test 3", Body: "Body 3"},
		}

		for _, msg := range messages {
			err := client.Send(msg)
			assert.NoError(t, err)
		}
	})

	t.Run("Batch met enkele fout", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client := createTestClient()
		client.SetDialer(mockSMTP)

		messages := []*services.EmailMessage{
			{To: "test1@example.com", Subject: "Test 1", Body: "Body 1"},
			{To: "", Subject: "Test 2", Body: "Body 2"}, // Deze zou moeten falen
			{To: "test3@example.com", Subject: "Test 3", Body: "Body 3"},
		}

		successCount := 0
		for _, msg := range messages {
			err := client.Send(msg)
			if err == nil {
				successCount++
			}
		}

		assert.Equal(t, 2, successCount)
	})
}

func TestSMTPClientWithInvalidConfig(t *testing.T) {
	client := createTestClientWithConfig(
		"",                 // Invalid host
		"587",              // port
		"test@test.com",    // user
		"testpass",         // password
		"noreply@test.com", // from
	)

	mockSMTP := newMockSMTP()
	mockSMTP.SetShouldFail(true)
	client.SetDialer(mockSMTP)

	msg := &services.EmailMessage{
		To:      "recipient@test.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err := client.Send(msg)
	assert.Error(t, err)
}

func TestSMTPClientWithInvalidPort(t *testing.T) {
	client := createTestClientWithConfig(
		"smtp.test.com",    // host
		"invalid",          // Invalid port
		"test@test.com",    // user
		"testpass",         // password
		"noreply@test.com", // from
	)

	mockSMTP := newMockSMTP()
	mockSMTP.SetShouldFail(true)
	client.SetDialer(mockSMTP)

	msg := &services.EmailMessage{
		To:      "recipient@test.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err := client.Send(msg)
	assert.Error(t, err)
}

func TestSMTPClientWithEmptyRecipient(t *testing.T) {
	client := createTestClient()

	msg := &services.EmailMessage{
		To:      "", // Empty recipient
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err := client.Send(msg)
	assert.Error(t, err)
}

func TestSMTPClientWithMockDialer(t *testing.T) {
	client := createTestClient()

	t.Run("Successful send", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client.SetDialer(mockSMTP)

		msg := &services.EmailMessage{
			To:      "recipient@test.com",
			Subject: "Test Subject",
			Body:    "Test Body",
		}

		err := client.Send(msg)
		assert.NoError(t, err)
		assert.Equal(t, msg, mockSMTP.GetLastEmail())
	})

	t.Run("Failed send", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		mockSMTP.SetShouldFail(true)
		client.SetDialer(mockSMTP)

		msg := &services.EmailMessage{
			To:      "recipient@test.com",
			Subject: "Test Subject",
			Body:    "Test Body",
		}

		err := client.Send(msg)
		assert.Error(t, err)
	})
}

func TestSMTPClientWithMissingConfig(t *testing.T) {
	client := createTestClient()

	mockSMTP := newMockSMTP()
	mockSMTP.SetShouldFail(true)
	client.SetDialer(mockSMTP)

	msg := &services.EmailMessage{
		To:      "recipient@test.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err := client.Send(msg)
	assert.Error(t, err)
}

func TestSMTPClientWithInvalidCredentials(t *testing.T) {
	client := createTestClient()

	mockSMTP := newMockSMTP()
	mockSMTP.SetShouldFail(true)
	client.SetDialer(mockSMTP)

	msg := &services.EmailMessage{
		To:      "recipient@test.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err := client.Send(msg)
	assert.Error(t, err)
}

func TestSMTPClientWithEmptyMessage(t *testing.T) {
	client := createTestClient()

	mockSMTP := newMockSMTP()
	client.SetDialer(mockSMTP)

	msg := &services.EmailMessage{
		To:      "recipient@test.com",
		Subject: "", // Empty subject
		Body:    "", // Empty body
	}

	err := client.Send(msg)
	assert.NoError(t, err) // Empty subject and body are allowed
}

func TestSMTPClientWithRegistrationConfig(t *testing.T) {
	client := createTestClient()

	mockSMTP := newMockSMTP()
	client.SetDialer(mockSMTP)

	msg := &services.EmailMessage{
		To:      "recipient@test.com",
		Subject: "Registration Test",
		Body:    "Test Body",
	}

	err := client.SendRegistration(msg)
	assert.NoError(t, err)
}

func TestSMTPClientValidation(t *testing.T) {
	client := createTestClient()

	t.Run("Empty recipient", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client.SetDialer(mockSMTP)

		msg := &services.EmailMessage{
			To:      "",
			Subject: "Test Subject",
			Body:    "Test Body",
		}

		err := client.Send(msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid recipient")
	})

	t.Run("Empty subject and body", func(t *testing.T) {
		mockSMTP := newMockSMTP()
		client.SetDialer(mockSMTP)

		msg := &services.EmailMessage{
			To:      "recipient@test.com",
			Subject: "",
			Body:    "",
		}

		err := client.Send(msg)
		assert.NoError(t, err)
		assert.Equal(t, msg, mockSMTP.GetLastEmail())
	})
}
