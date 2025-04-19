package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/services" // Assuming your service interfaces are here
	"net/http"
	"os"
	"strings"

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
	To                string                 `json:"to" validate:"required,email"`
	Subject           string                 `json:"subject" validate:"required,min=1"`
	From              string                 `json:"from,omitempty" validate:"omitempty,email"` // OPTIONEEL: Gewenst afzenderadres
	Body              string                 `json:"body,omitempty"`                            // Body is now optional if template is used
	TemplateName      string                 `json:"template_name,omitempty"`
	TemplateVariables map[string]interface{} `json:"template_variables,omitempty"`
}

// ValidateSendMailRequest performs custom validation logic
func ValidateSendMailRequest(sl validator.StructLevel) {
	req := sl.Current().Interface().(SendMailRequest)

	if req.Body == "" && req.TemplateName == "" {
		sl.ReportError(req.Body, "body", "Body", "required_without", "TemplateName")
		sl.ReportError(req.TemplateName, "template_name", "TemplateName", "required_without", "Body")
	}

	if req.TemplateName != "" && req.Body != "" {
		sl.ReportError(req.TemplateName, "template_name", "TemplateName", "excluded_with", "Body")
		sl.ReportError(req.Body, "body", "Body", "excluded_with", "TemplateName")
	}

	if req.TemplateName != "" && (req.TemplateVariables == nil || len(req.TemplateVariables) == 0) {
		// Optionally require variables if template is used, or allow empty maps
		// sl.ReportError(req.TemplateVariables, "template_variables", "TemplateVariables", "required_with", "TemplateName")
	}
}

// NewAdminMailHandler creates a new AdminMailHandler instance.
func NewAdminMailHandler(emailSvc services.EmailSender, authSvc services.AuthService) *AdminMailHandler {
	v := validator.New()
	// Register custom validation
	v.RegisterStructValidation(ValidateSendMailRequest, SendMailRequest{})

	return &AdminMailHandler{
		emailService: emailSvc,
		authService:  authSvc,
		validate:     v,
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

	// Validate the request body using the validator (including custom validation)
	if err := h.validate.Struct(req); err != nil {
		errors := FormatValidationErrors(err) // Helper function to format errors nicely
		logger.Warn("Validation failed for /admin/mail/send", "errors", errors, "ip", c.IP())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Validatie mislukt", "details": errors})
	}

	// --- AFZENDER VALIDATIE ---
	var actualFromAddress string                            // Dit wordt doorgegeven aan de service
	allowedSendersEnv := os.Getenv("ALLOWED_SENDER_EMAILS") // Bv: "info@...,jeffrey@...,marieke@..."
	allowedSenders := []string{}
	if allowedSendersEnv != "" {
		allowedSenders = strings.Split(allowedSendersEnv, ",")
	}
	isAllowed := false
	requestedFrom := strings.TrimSpace(req.From) // Trim spaties

	if requestedFrom != "" {
		if len(allowedSenders) == 0 {
			logger.Warn("Requested sender address provided, but ALLOWED_SENDER_EMAILS is not configured. Falling back to default.", "requested_from", requestedFrom)
		} else {
			for _, allowed := range allowedSenders {
				if strings.TrimSpace(allowed) == requestedFrom {
					isAllowed = true
					break
				}
			}
		}

		if isAllowed {
			actualFromAddress = requestedFrom
			logger.Info("Using requested sender address", "from", actualFromAddress)
		} else {
			logger.Warn("Requested sender address not allowed or ALLOWED_SENDER_EMAILS not configured, falling back to default", "requested_from", requestedFrom, "allowed_list", allowedSendersEnv)
			// Laat actualFromAddress leeg, zodat de service de default pakt
		}
	} else {
		logger.Info("No sender address requested, using default")
		// Laat actualFromAddress leeg, zodat de service de default pakt
	}
	// --- EINDE AFZENDER VALIDATIE ---

	// Log the attempt details (including admin user ID from context)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		logger.Error("UserID not found in context after AuthMiddleware", "path", c.Path())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Interne serverfout: Gebruikerscontext ontbreekt."})
	}

	logger.Info("Processing admin mail send request", "adminUserID", userID, "to", req.To, "subject", req.Subject, "from", actualFromAddress, "template", req.TemplateName, "has_body", req.Body != "", "test_mode", isTestMode, "ip", c.IP())

	// --- Test Mode Check ---
	if isTestMode {
		logger.Info("[TEST MODE] Admin mail send skipped.", "adminUserID", userID, "to", req.To, "subject", req.Subject, "from", actualFromAddress, "template", req.TemplateName, "ip", c.IP())
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"success":   true,
			"message":   "[TEST MODE] Email verzoek verwerkt, geen echte email verzonden.",
			"test_mode": true,
		})
	}
	// --- End Test Mode Check ---

	var err error
	if req.TemplateName != "" {
		// Send using template, pass the validated 'from' address (can be empty)
		err = h.emailService.SendTemplateEmail(req.To, req.Subject, req.TemplateName, req.TemplateVariables, actualFromAddress)
	} else {
		// Send using plain body, pass the validated 'from' address (can be empty)
		err = h.emailService.SendEmail(req.To, req.Subject, req.Body, actualFromAddress)
	}

	if err != nil {
		logger.Error("Error sending admin email", "error", err, "to", req.To, "from", actualFromAddress, "template", req.TemplateName, "adminUserID", userID, "ip", c.IP())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Kon email niet verzenden: " + err.Error()})
	}

	logger.Info("Admin mail sent successfully", "adminUserID", userID, "to", req.To, "from", actualFromAddress, "subject", req.Subject, "template", req.TemplateName, "ip", c.IP())
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
