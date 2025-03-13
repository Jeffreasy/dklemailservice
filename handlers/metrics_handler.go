package handlers

import (
	"dklautomationgo/services"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MetricsHandler verwerkt verzoeken voor statistieken
type MetricsHandler struct {
	emailMetrics *services.EmailMetrics
	rateLimiter  *services.RateLimiter
}

// NewMetricsHandler maakt een nieuwe metrics handler
func NewMetricsHandler(metrics *services.EmailMetrics, limiter *services.RateLimiter) *MetricsHandler {
	return &MetricsHandler{
		emailMetrics: metrics,
		rateLimiter:  limiter,
	}
}

// HandleGetEmailMetrics geeft email statistieken terug
func (h *MetricsHandler) HandleGetEmailMetrics(c *fiber.Ctx) error {
	// Alleen toegankelijk voor admins/monitoring
	apiKey := c.Get("X-API-Key")
	if apiKey != getAdminAPIKey() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Ongeautoriseerd",
		})
	}

	total := h.emailMetrics.GetTotalEmails()
	successRate := h.emailMetrics.GetSuccessRate()
	byType := h.emailMetrics.GetEmailsByType()

	return c.JSON(fiber.Map{
		"total_emails":   total,
		"success_rate":   successRate,
		"emails_by_type": byType,
		"generated_at":   time.Now(),
	})
}

// HandleGetRateLimits geeft informatie over rate limiting
func (h *MetricsHandler) HandleGetRateLimits(c *fiber.Ctx) error {
	// Alleen toegankelijk voor admins/monitoring
	apiKey := c.Get("X-API-Key")
	if apiKey != getAdminAPIKey() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Ongeautoriseerd",
		})
	}

	// Voorbeeld operaties die we willen controleren
	operationTypes := []string{"contact_email", "aanmelding_email"}

	result := make(map[string]interface{})

	for _, opType := range operationTypes {
		// Globale telling
		globalCount := h.rateLimiter.GetCurrentCount(opType, "")
		result[opType] = fiber.Map{
			"global_count": globalCount,
		}
	}

	return c.JSON(fiber.Map{
		"rate_limits":  result,
		"generated_at": time.Now(),
	})
}

// Helper functie voor de admin API key
func getAdminAPIKey() string {
	return os.Getenv("ADMIN_API_KEY")
}

func (h *MetricsHandler) GetEmailMetrics(w http.ResponseWriter, r *http.Request) {
	metricsData := map[string]interface{}{
		"total_emails":   h.emailMetrics.GetTotalEmails(),
		"success_rate":   h.emailMetrics.GetSuccessRate(),
		"emails_by_type": h.emailMetrics.GetEmailsByType(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metricsData)
}
