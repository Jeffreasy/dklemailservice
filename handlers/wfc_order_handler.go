package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// WFCOrderHandler handles requests for Whisky for Charity order emails
type WFCOrderHandler struct {
	emailService *services.EmailService
}

// NewWFCOrderHandler creates a new WFC order handler
func NewWFCOrderHandler(emailService *services.EmailService) *WFCOrderHandler {
	return &WFCOrderHandler{
		emailService: emailService,
	}
}

// WFCAPIKeyMiddleware creates a middleware for API key authentication
func WFCAPIKeyMiddleware(keyEnvVar string) fiber.Handler {
	if keyEnvVar == "" {
		keyEnvVar = "WFC_API_KEY"
	}

	return func(c *fiber.Ctx) error {
		apiKey := os.Getenv(keyEnvVar)
		if apiKey == "" {
			logger.Warn("API key not configured", "keyEnvVar", keyEnvVar)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "API authentication not configured",
			})
		}

		providedKey := c.Get("X-API-Key")
		if providedKey == "" {
			logger.Warn("No API key provided in request")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key required",
			})
		}

		if providedKey != apiKey {
			logger.Warn("Invalid API key provided")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		return c.Next()
	}
}

// HandleWFCOrderEmail processes order email requests
func (h *WFCOrderHandler) HandleWFCOrderEmail(c *fiber.Ctx) error {
	var req models.WFCOrderRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Invalid WFC order request", "error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Check for required fields
	if req.OrderID == "" || req.CustomerEmail == "" || req.CustomerName == "" {
		logger.Error("Missing required fields in WFC order request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	// Create order data
	order := &models.WFCOrder{
		ID:            req.OrderID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		TotalAmount:   req.TotalAmount,
		Items:         req.Items,
		Status:        "pending", // Default status
		CreatedAt:     time.Now(),
	}

	// Create email data
	emailData := &models.WFCOrderEmailData{
		Order:   order,
		ToAdmin: req.NotifyAdmin,
		SiteURL: os.Getenv("WFC_SITE_URL"),
	}

	// If this is an admin notification, set admin email
	if req.NotifyAdmin {
		emailData.AdminEmail = os.Getenv("WFC_ADMIN_EMAIL")
		if emailData.AdminEmail == "" {
			emailData.AdminEmail = "info@whiskyforcharity.com" // Fallback
		}
	}

	// Send the email
	if err := h.emailService.SendWFCOrderEmail(emailData); err != nil {
		logger.Error("Failed to send WFC order email", "error", err, "order_id", req.OrderID)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	// Respond with success
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success":  true,
		"message":  "Order email sent successfully",
		"order_id": req.OrderID,
	})
}

// RegisterWFCOrderRoutes registers the WFC order routes
func RegisterWFCOrderRoutes(app *fiber.App, emailService *services.EmailService) {
	handler := NewWFCOrderHandler(emailService)

	// Create a group for WFC endpoints with API key authentication
	wfcGroup := app.Group("/api/wfc")
	wfcGroup.Use(WFCAPIKeyMiddleware("WFC_API_KEY"))

	// Register routes
	wfcGroup.Post("/order-email", handler.HandleWFCOrderEmail)
}
