package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EmailServiceInterface interface {
	SendContactEmail(data *models.ContactEmailData) error
	SendAanmeldingEmail(data *models.AanmeldingEmailData) error
}

type EmailHandler struct {
	emailService EmailServiceInterface
}

func NewEmailHandler(emailService EmailServiceInterface) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

func (h *EmailHandler) HandleContactEmail(c *fiber.Ctx) error {
	var request models.ContactFormulier
	start := time.Now()

	if err := c.BodyParser(&request); err != nil {
		logger.Error("Fout bij parsen van contact formulier",
			"error", err,
			"remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Ongeldig verzoek: " + err.Error(),
		})
	}

	// Detecteer test modus
	var testMode bool
	if testModeValue := c.Get("X-Test-Mode"); testModeValue == "true" {
		testMode = true
		logger.Info("Test modus gedetecteerd via header", "remote_ip", c.IP())
	}

	if c.Locals("test_mode") != nil {
		testMode = true
		logger.Info("Test modus gedetecteerd via locals", "remote_ip", c.IP())
	}

	var requestMap map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestMap); err == nil {
		if val, ok := requestMap["test_mode"]; ok && val.(bool) {
			testMode = true
			logger.Info("Test modus gedetecteerd via body parameter", "remote_ip", c.IP())
		}
	}

	// Log the incoming request
	logger.Info("Contact formulier ontvangen",
		"naam", request.Naam,
		"email", request.Email,
		"remote_ip", c.IP(),
		"test_mode", testMode)

	// Validate the request
	if request.Naam == "" || request.Email == "" || request.Bericht == "" {
		logger.Warn("Onvolledig contact formulier",
			"naam", request.Naam,
			"email", request.Email,
			"bericht_empty", request.Bericht == "",
			"remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Naam, email en bericht zijn verplicht",
		})
	}

	if !request.PrivacyAkkoord {
		logger.Warn("Privacy niet geaccepteerd",
			"naam", request.Naam,
			"email", request.Email,
			"remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Je moet akkoord gaan met het privacybeleid",
		})
	}

	// Set CreatedAt if not provided
	if request.CreatedAt.IsZero() {
		request.CreatedAt = time.Now()
	}

	// Send email to admin
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "info@dekoninklijkeloop.nl" // Default admin email
		logger.Warn("ADMIN_EMAIL niet geconfigureerd, gebruik standaardwaarde",
			"default", adminEmail)
	}

	adminEmailData := &models.ContactEmailData{
		Contact:    &request,
		AdminEmail: adminEmail,
		ToAdmin:    true,
	}

	// In testmodus sturen we geen echte emails
	if testMode {
		logger.Info("Test modus: Geen admin email verzonden", "admin_email", adminEmail)
	} else {
		logger.Info("Admin email wordt verzonden",
			"admin_email", adminEmail,
			"contact_naam", request.Naam)
		if err := h.emailService.SendContactEmail(adminEmailData); err != nil {
			logger.Error("Fout bij verzenden admin email",
				"error", err,
				"admin_email", adminEmail,
				"elapsed", time.Since(start))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Fout bij het verzenden van de email: " + err.Error(),
			})
		}
		logger.Info("Admin email verzonden",
			"admin_email", adminEmail,
			"elapsed", time.Since(start))
	}

	// Send confirmation email to user
	userEmailData := &models.ContactEmailData{
		Contact:    &request,
		AdminEmail: adminEmail,
		ToAdmin:    false,
	}

	// In testmodus sturen we geen echte emails
	if testMode {
		logger.Info("Test modus: Geen gebruiker email verzonden", "user_email", request.Email)
	} else {
		logger.Info("Bevestigingsemail wordt verzonden",
			"user_email", request.Email,
			"naam", request.Naam)
		if err := h.emailService.SendContactEmail(userEmailData); err != nil {
			logger.Error("Fout bij verzenden bevestigingsemail",
				"error", err,
				"user_email", request.Email,
				"elapsed", time.Since(start))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Fout bij het verzenden van de bevestigingsemail: " + err.Error(),
			})
		}
		logger.Info("Bevestigingsemail verzonden",
			"user_email", request.Email,
			"elapsed", time.Since(start))
	}

	// Return success
	if testMode {
		logger.Info("Contact formulier succesvol verwerkt in test modus",
			"naam", request.Naam,
			"email", request.Email,
			"total_elapsed", time.Since(start))
		return c.JSON(fiber.Map{
			"success":   true,
			"message":   "[TEST MODE] Je bericht is verwerkt (geen echte email verzonden).",
			"test_mode": true,
		})
	} else {
		logger.Info("Contact formulier succesvol verwerkt",
			"naam", request.Naam,
			"email", request.Email,
			"total_elapsed", time.Since(start))
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Je bericht is verzonden! Je ontvangt ook een bevestiging per email.",
		})
	}
}

func (h *EmailHandler) HandleAanmeldingEmail(c *fiber.Ctx) error {
	var aanmelding models.AanmeldingFormulier
	start := time.Now()

	if err := c.BodyParser(&aanmelding); err != nil {
		logger.Error("Fout bij parsen van aanmelding formulier",
			"error", err,
			"remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Ongeldig verzoek: " + err.Error(),
		})
	}

	// Detecteer test modus
	var testMode bool
	if testModeValue := c.Get("X-Test-Mode"); testModeValue == "true" {
		testMode = true
		logger.Info("Test modus gedetecteerd via header", "remote_ip", c.IP())
	}

	if c.Locals("test_mode") != nil {
		testMode = true
		logger.Info("Test modus gedetecteerd via locals", "remote_ip", c.IP())
	}

	var requestMap map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestMap); err == nil {
		if val, ok := requestMap["test_mode"]; ok && val.(bool) {
			testMode = true
			logger.Info("Test modus gedetecteerd via body parameter", "remote_ip", c.IP())
		}
	}

	// Log the incoming request
	logger.Info("Aanmelding formulier ontvangen",
		"naam", aanmelding.Naam,
		"email", aanmelding.Email,
		"remote_ip", c.IP(),
		"test_mode", testMode)

	// Validate required fields
	if aanmelding.Naam == "" {
		logger.Warn("Ontbrekende naam in aanmelding",
			"email", aanmelding.Email,
			"remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Naam is verplicht",
		})
	}

	if aanmelding.Email == "" {
		logger.Warn("Ontbrekende email in aanmelding",
			"naam", aanmelding.Naam,
			"remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Email is verplicht",
		})
	}

	// Validate email format
	if !strings.Contains(aanmelding.Email, "@") || !strings.Contains(aanmelding.Email, ".") {
		logger.Warn("Ongeldig email formaat",
			"email", aanmelding.Email,
			"remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Ongeldig email adres",
		})
	}

	// Validate terms acceptance
	if !aanmelding.Terms {
		logger.Warn("Terms niet geaccepteerd",
			"naam", aanmelding.Naam,
			"email", aanmelding.Email,
			"remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Je moet akkoord gaan met de voorwaarden",
		})
	}

	// Send email to admin
	adminEmail := os.Getenv("REGISTRATION_EMAIL")
	if adminEmail == "" {
		adminEmail = "inschrijving@dekoninklijkeloop.nl" // Default registration email
		logger.Warn("REGISTRATION_EMAIL niet geconfigureerd, gebruik standaardwaarde",
			"default", adminEmail)
	}

	// Stuur email naar admin
	adminEmailData := &models.AanmeldingEmailData{
		ToAdmin:    true,
		Aanmelding: &aanmelding,
		AdminEmail: adminEmail,
	}

	// In testmodus sturen we geen echte emails
	if testMode {
		logger.Info("Test modus: Geen admin email verzonden", "admin_email", adminEmail)
	} else {
		logger.Info("Admin email wordt verzonden",
			"admin_email", adminEmail,
			"aanmelding_naam", aanmelding.Naam)
		if err := h.emailService.SendAanmeldingEmail(adminEmailData); err != nil {
			logger.Error("Fout bij verzenden admin email",
				"error", err,
				"admin_email", adminEmail,
				"elapsed", time.Since(start))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Fout bij het verzenden van de email: " + err.Error(),
			})
		}
		logger.Info("Admin email verzonden",
			"admin_email", adminEmail,
			"elapsed", time.Since(start))
	}

	// Stuur bevestigingsemail naar gebruiker
	userEmailData := &models.AanmeldingEmailData{
		ToAdmin:    false,
		Aanmelding: &aanmelding,
	}

	// In testmodus sturen we geen echte emails
	if testMode {
		logger.Info("Test modus: Geen gebruiker email verzonden", "user_email", aanmelding.Email)
	} else {
		logger.Info("Bevestigingsemail wordt verzonden",
			"user_email", aanmelding.Email,
			"naam", aanmelding.Naam)
		if err := h.emailService.SendAanmeldingEmail(userEmailData); err != nil {
			logger.Error("Fout bij verzenden bevestigingsemail",
				"error", err,
				"user_email", aanmelding.Email,
				"elapsed", time.Since(start))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Fout bij het verzenden van de bevestigingsemail: " + err.Error(),
			})
		}
		logger.Info("Bevestigingsemail verzonden",
			"user_email", aanmelding.Email,
			"elapsed", time.Since(start))
	}

	// Return success
	if testMode {
		logger.Info("Aanmelding formulier succesvol verwerkt in test modus",
			"naam", aanmelding.Naam,
			"email", aanmelding.Email,
			"total_elapsed", time.Since(start))
		return c.JSON(fiber.Map{
			"success":   true,
			"message":   "[TEST MODE] Je aanmelding is verwerkt (geen echte email verzonden).",
			"test_mode": true,
		})
	} else {
		logger.Info("Aanmelding formulier succesvol verwerkt",
			"naam", aanmelding.Naam,
			"email", aanmelding.Email,
			"total_elapsed", time.Since(start))
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Je aanmelding is verzonden! Je ontvangt ook een bevestiging per email.",
		})
	}
}

// LogUserActivity logt gebruikersactiviteit (helper voor tests)
func (h *EmailHandler) LogUserActivity(email, activity, ip string) {
	logger.Info("Gebruikersactiviteit",
		"email", email,
		"activiteit", activity,
		"ip", ip)
}
