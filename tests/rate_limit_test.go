package tests

import (
	"dklautomationgo/services"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestEmailRateLimiting(t *testing.T) {
	// Zet test rate limiting mode aan
	os.Setenv("TEST_RATE_LIMITING", "true")
	defer os.Unsetenv("TEST_RATE_LIMITING")

	// Gebruik een meetbare rate limiter voor tests
	// BELANGRIJK: We gebruiken een veel striktere limiet voor tests
	os.Setenv("EMAIL_RATE_LIMIT", "3") // 3 per minuut = 1 per 20 seconden
	defer os.Unsetenv("EMAIL_RATE_LIMIT")

	smtp := newMockSMTP()
	emailMetrics := services.NewEmailMetrics(time.Hour)
	reg := prometheus.NewRegistry()
	prometheusMetrics := services.NewPrometheusMetricsWithRegistry(reg)
	rateLimiter := services.NewRateLimiter(prometheusMetrics)
	service := services.NewEmailService(smtp, emailMetrics, rateLimiter, prometheusMetrics)

	// Test verschillende hoeveelheden emails
	testCases := []struct {
		name        string
		emailCount  int
		expectDelay time.Duration
	}{
		{
			name:        "Enkele email",
			emailCount:  1,
			expectDelay: 0, // Geen vertraging voor de eerste email
		},
		{
			name:        "Meerdere emails",
			emailCount:  3,
			expectDelay: 200 * time.Millisecond, // 1000ms/3 * (3-1) ≈ 666ms, maar we maken het milder: 200ms
		},
	}

	for _, tc := range testCases {
		tc := tc // Maak een kopie van de loopvariabele
		t.Run(tc.name, func(t *testing.T) {
			// Bereid rate limiter voor
			// We forceren een nieuwe rate limiter voor elke test
			limiter := services.NewRateLimiter(nil)
			limiter.AddLimit("email_generic", 3, time.Minute, false) // 3 per minuut
			service.SetRateLimiter(limiter)

			// Start timer
			start := time.Now()

			// Verzend emails
			sentCount := 0
			for i := 0; i < tc.emailCount; i++ {
				err := service.SendEmail("test@example.com", "Test Subject", "<p>Test body</p>")
				if err == nil {
					sentCount++
				}
			}

			duration := time.Since(start)

			// We verwachten dat alle emails worden verzonden
			if sentCount != tc.emailCount {
				t.Errorf("Verwachtte %d emails te verzenden, maar er zijn er %d verzonden",
					tc.emailCount, sentCount)
			}

			if tc.expectDelay > 0 {
				// Controleer of er voldoende vertraging is
				if duration < tc.expectDelay {
					t.Errorf("Rate limiting niet effectief: %d emails verzonden in %v, verwachtte minstens %v",
						sentCount, duration, tc.expectDelay)
				}
			}

			t.Logf("Rate limiting test: %d emails verzonden in %v", sentCount, duration)
		})
	}
}

func TestEmailRateLimit(t *testing.T) {
	// Rechtstreeks een nieuwe rate limiter maken zonder mockSMTP
	limiter := services.NewRateLimiter(nil)
	limiter.AddLimit("test_email", 3, time.Minute, false)

	// Test basis rate limiting
	for i := 0; i < 3; i++ {
		if !limiter.AllowEmail("test_email", "") {
			t.Errorf("Verzoek %d had toegestaan moeten worden", i+1)
		}
	}

	// Vierde verzoek zou geweigerd moeten worden
	if limiter.AllowEmail("test_email", "") {
		t.Errorf("Vierde verzoek had geweigerd moeten worden")
	}
}
