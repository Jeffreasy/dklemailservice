package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/services" // Assuming your service interfaces are here
	"net/http"

	"github.com/go-playground/validator/v10" // For request validation
	"github.com/gofiber/fiber/v2"
)

// AdminMailHandler handles requests related to sending emails by admins.
type AdminMailHandler struct {
	emailService services.EmailSender // Use the existing EmailSender interface
	authService  services.AuthService // Use the existing AuthService interface
	validate     *validator.Validate  // Keep validator instance
}

// SendMailRequest defines the expected JSON body for the send mail endpoint.
type SendMailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required,min=1"`
	Body    string `json:"body" validate:"required,min=1"`
}

// NewAdminMailHandler creates a new AdminMailHandler instance.
func NewAdminMailHandler(emailSvc services.EmailSender, authSvc services.AuthService) *AdminMailHandler {
	return &AdminMailHandler{
		emailService: emailSvc,
		authService:  authSvc,
		validate:     validator.New(),
	}
}

// HandleSendMail handles the POST /api/admin/mail/send request.
func (h *AdminMailHandler) HandleSendMail(c *fiber.Ctx) error {
	var req SendMailRequest
	isTestMode := false // Default to false

	// Check for test mode from middleware
	if testModeVal := c.Locals("test_mode"); testModeVal != nil {
		if val, ok := testModeVal.(bool); ok && val {
			isTestMode = true
		}
	}

	// Log test mode status
	if isTestMode {
		logger.Info("Admin mail send request received in TEST MODE", "ip", c.IP())
	} else {
		logger.Info("Admin mail send request received", "ip", c.IP())
	}

	// Parse the request body
	if err := c.BodyParser(&req); err != nil {
		logger.Warn("Invalid request body for /admin/mail/send", "error", err, "ip", c.IP())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Ongeldige request body: " + err.Error()})
	}

	// Validate the request body using the validator
	if err := h.validate.Struct(req); err != nil {
		errors := FormatValidationErrors(err) // Helper function to format errors nicely
		logger.Warn("Validation failed for /admin/mail/send", "errors", errors, "ip", c.IP())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Validatie mislukt", "details": errors})
	}

	// Log the attempt details (including admin user ID from context)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		logger.Error("UserID not found in context after AuthMiddleware", "path", c.Path())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Interne serverfout: Gebruikerscontext ontbreekt."})
	}

	logger.Info("Processing admin mail send request", "adminUserID", userID, "to", req.To, "subject", req.Subject, "test_mode", isTestMode, "ip", c.IP())

	// --- Test Mode Check ---
	if isTestMode {
		logger.Info("[TEST MODE] Admin mail send skipped.", "adminUserID", userID, "to", req.To, "subject", req.Subject, "ip", c.IP())
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"success":   true,
			"message":   "[TEST MODE] Email verzoek verwerkt, geen echte email verzonden.",
			"test_mode": true,
		})
	}
	// --- End Test Mode Check ---

	// Send the email using the existing generic SendEmail method (only if not in test mode)
	err := h.emailService.SendEmail(req.To, req.Subject, req.Body)
	if err != nil {
		logger.Error("Error sending admin email", "error", err, "to", req.To, "adminUserID", userID, "ip", c.IP())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Kon email niet verzenden: " + err.Error()})
	}

	logger.Info("Admin mail sent successfully", "adminUserID", userID, "to", req.To, "subject", req.Subject, "ip", c.IP())
	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true, "message": "Email succesvol verzonden."})
}

// RegisterRoutes registers the admin mail routes, protected by authentication and authorization middleware.
func (h *AdminMailHandler) RegisterRoutes(app *fiber.App) {
	// Create a group for admin mail actions, protected by AuthMiddleware and AdminMiddleware
	adminMailGroup := app.Group("/api/admin/mail", AuthMiddleware(h.authService), AdminMiddleware(h.authService))

	// Register the POST route for sending mail
	adminMailGroup.Post("/send", h.HandleSendMail)

	logger.Info("Admin mail routes registered under /api/admin/mail")
}

// FormatValidationErrors is a helper to make validation errors more readable (optional but nice).
// You might already have a similar utility function elsewhere.
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := fieldErr.Field()
			tag := fieldErr.Tag()
			errors[field] = "Validation failed on tag '" + tag + "'"
		}
	} else {
		errors["general"] = err.Error()
	}
	return errors
}
