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

	// Get admin email from environment
	adminEmail := os.Getenv("WFC_ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "info@whiskyforcharity.com" // Fallback
	}

	// Get site URL from environment
	siteURL := os.Getenv("WFC_SITE_URL")

	// Track success
	var emailErrors []error
	var customerEmailSent bool
	var adminEmailSent bool

	// 1. Always send customer confirmation email
	customerEmailData := &models.WFCOrderEmailData{
		Order:      order,
		ToAdmin:    false,      // This is for customer
		AdminEmail: adminEmail, // Include admin email even though not directly used
		SiteURL:    siteURL,
	}

	// Send customer email
	if err := h.emailService.SendWFCOrderEmail(customerEmailData); err != nil {
		logger.Error("Failed to send WFC customer order email", "error", err, "order_id", req.OrderID)
		emailErrors = append(emailErrors, err)
	} else {
		customerEmailSent = true
		logger.Info("WFC customer order email sent successfully", "recipient", order.CustomerEmail, "order_id", req.OrderID)
	}

	// 2. Always send admin notification email
	adminEmailData := &models.WFCOrderEmailData{
		Order:      order,
		ToAdmin:    true, // This is for admin
		AdminEmail: adminEmail,
		SiteURL:    siteURL,
	}

	// Send admin email
	if err := h.emailService.SendWFCOrderEmail(adminEmailData); err != nil {
		logger.Error("Failed to send WFC admin order email", "error", err, "order_id", req.OrderID)
		emailErrors = append(emailErrors, err)
	} else {
		adminEmailSent = true
		logger.Info("WFC admin order email sent successfully", "recipient", adminEmail, "order_id", req.OrderID)
	}

	// Check if at least one email was sent successfully
	if !customerEmailSent && !adminEmailSent {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":        "Failed to send all emails",
			"error_detail": emailErrors[0].Error(),
		})
	}

	// Respond with success status for each email
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success":             true,
		"customer_email_sent": customerEmailSent,
		"admin_email_sent":    adminEmailSent,
		"order_id":            req.OrderID,
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
