package tests

import (
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"dklautomationgo/tests/mocks"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRepository creates a mock repository for testing
func mockRepository() *repository.Repository {
	mockDB := mocks.NewMockDB()
	return &repository.Repository{
		Contact:   mocks.NewMockContactRepository(mockDB),
		Migratie:  mocks.NewMockMigratieRepository(mockDB),
		Gebruiker: mocks.NewMockGebruikerRepository(mockDB),
	}
}

func TestServiceFactory(t *testing.T) {
	// Setup
	repo := mockRepository()

	// Test factory creation
	factory := services.NewServiceFactory(repo)

	// Verify all services are created
	assert.NotNil(t, factory, "ServiceFactory zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailService, "EmailService zou niet nil moeten zijn")
	assert.NotNil(t, factory.SMTPClient, "SMTPClient zou niet nil moeten zijn")
	assert.NotNil(t, factory.RateLimiter, "RateLimiter zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailMetrics, "EmailMetrics zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailBatcher, "EmailBatcher zou niet nil moeten zijn")
	assert.NotNil(t, factory.AuthService, "AuthService zou niet nil moeten zijn")

	// Test RateLimiter interface conformance - remove unnecessary type assertion
	assert.Implements(t, (*services.RateLimiterInterface)(nil), factory.RateLimiter)
}

func TestServiceFactory_WithCustomConfig(t *testing.T) {
	// Setup
	repo := mockRepository()

	// Set custom env values
	os.Setenv("RATE_LIMIT_MAX", "20")
	os.Setenv("RATE_LIMIT_DURATION", "2")
	defer func() {
		os.Unsetenv("RATE_LIMIT_MAX")
		os.Unsetenv("RATE_LIMIT_DURATION")
	}()

	// Test factory creation with custom config
	factory := services.NewServiceFactory(repo)

	// Verify all services are created
	assert.NotNil(t, factory, "ServiceFactory zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailService, "EmailService zou niet nil moeten zijn")
	assert.NotNil(t, factory.SMTPClient, "SMTPClient zou niet nil moeten zijn")
	assert.NotNil(t, factory.RateLimiter, "RateLimiter zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailMetrics, "EmailMetrics zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailBatcher, "EmailBatcher zou niet nil moeten zijn")
	assert.NotNil(t, factory.AuthService, "AuthService zou niet nil moeten zijn")

	// Test RateLimiter interface conformance - remove unnecessary type assertion
	assert.Implements(t, (*services.RateLimiterInterface)(nil), factory.RateLimiter)
}
