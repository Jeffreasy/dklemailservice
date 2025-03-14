package tests

import (
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRepositoryFactory creates a mock repository factory for testing
func mockRepositoryFactory() *repository.RepositoryFactory {
	// Create a mock DB - we don't actually need a real DB for these tests
	// since we're just testing that the service factory creates services
	return &repository.RepositoryFactory{}
}

func TestServiceFactory(t *testing.T) {
	// Setup
	repoFactory := mockRepositoryFactory()

	// Test factory creation
	factory := services.NewServiceFactory(repoFactory)

	// Verify all services are created
	assert.NotNil(t, factory, "Factory zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailService, "EmailService zou niet nil moeten zijn")
	assert.NotNil(t, factory.SMTPClient, "SMTPClient zou niet nil moeten zijn")
	assert.NotNil(t, factory.RateLimiter, "RateLimiter zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailMetrics, "EmailMetrics zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailBatcher, "EmailBatcher zou niet nil moeten zijn")

	// Verify the RateLimiter implements the RateLimiterInterface
	assert.Implements(t, (*services.RateLimiterInterface)(nil), factory.RateLimiter,
		"RateLimiter zou de RateLimiterInterface moeten implementeren")
}

func TestServiceFactory_WithCustomConfig(t *testing.T) {
	// Setup
	repoFactory := mockRepositoryFactory()

	// Sla originele omgevingsvariabelen op
	originalContactLimit := os.Getenv("CONTACT_LIMIT_COUNT")
	originalAanmeldingLimit := os.Getenv("AANMELDING_LIMIT_COUNT")
	originalEmailMetricsReset := os.Getenv("EMAIL_METRICS_RESET_INTERVAL")
	originalEmailBatchSize := os.Getenv("EMAIL_BATCH_SIZE")

	// Herstel omgevingsvariabelen na de test
	defer func() {
		os.Setenv("CONTACT_LIMIT_COUNT", originalContactLimit)
		os.Setenv("AANMELDING_LIMIT_COUNT", originalAanmeldingLimit)
		os.Setenv("EMAIL_METRICS_RESET_INTERVAL", originalEmailMetricsReset)
		os.Setenv("EMAIL_BATCH_SIZE", originalEmailBatchSize)
	}()

	// Stel aangepaste configuratie in
	os.Setenv("CONTACT_LIMIT_COUNT", "10")
	os.Setenv("AANMELDING_LIMIT_COUNT", "5")
	os.Setenv("EMAIL_METRICS_RESET_INTERVAL", "3600")
	os.Setenv("EMAIL_BATCH_SIZE", "20")

	// Test factory creation met aangepaste configuratie
	factory := services.NewServiceFactory(repoFactory)

	// Verify all services are created
	assert.NotNil(t, factory, "Factory zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailService, "EmailService zou niet nil moeten zijn")
	assert.NotNil(t, factory.SMTPClient, "SMTPClient zou niet nil moeten zijn")
	assert.NotNil(t, factory.RateLimiter, "RateLimiter zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailMetrics, "EmailMetrics zou niet nil moeten zijn")
	assert.NotNil(t, factory.EmailBatcher, "EmailBatcher zou niet nil moeten zijn")

	// Verify the RateLimiter implements the RateLimiterInterface
	assert.Implements(t, (*services.RateLimiterInterface)(nil), factory.RateLimiter,
		"RateLimiter zou de RateLimiterInterface moeten implementeren")

	// Get the limits from the RateLimiter
	limits := factory.RateLimiter.GetLimits()

	contactLimit, exists := limits["contact"]
	assert.True(t, exists, "Contact limiet zou moeten bestaan")
	assert.Equal(t, 10, contactLimit.Count, "Contact limiet zou 10 moeten zijn")

	aanmeldingLimit, exists := limits["aanmelding"]
	assert.True(t, exists, "Aanmelding limiet zou moeten bestaan")
	assert.Equal(t, 5, aanmeldingLimit.Count, "Aanmelding limiet zou 5 moeten zijn")
}
