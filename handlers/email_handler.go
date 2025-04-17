package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// EmailServiceInterface definieert de interface voor email operaties
type EmailServiceInterface interface {
	SendContactEmail(data *models.ContactEmailData) error
	SendAanmeldingEmail(data *models.AanmeldingEmailData) error
}

// EmailHandler verzorgt de afhandeling van email verzoeken
type EmailHandler struct {
	emailService        EmailServiceInterface
	notificationService services.NotificationService
	aanmeldingRepo      repository.AanmeldingRepository
}

// NewEmailHandler maakt een nieuwe EmailHandler
func NewEmailHandler(
	emailService EmailServiceInterface,
	notificationService services.NotificationService,
	aanmeldingRepo repository.AanmeldingRepository,
) *EmailHandler {
	return &EmailHandler{
		emailService:        emailService,
		notificationService: notificationService,
		aanmeldingRepo:      aanmeldingRepo,
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

	// Stuur een notificatie over het nieuwe contactformulier
	// We doen dit alleen in productie modus
	h.sendContactNotification(&request, testMode)

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

	// Maak een Aanmelding object aan voor de database
	// We gebruiken pointers voor optionele velden zoals Bijzonderheden
	nieuweAanmelding := &models.Aanmelding{
		Naam:           aanmelding.Naam,
		Email:          aanmelding.Email,
		Telefoon:       aanmelding.Telefoon,
		Rol:            aanmelding.Rol,
		Afstand:        aanmelding.Afstand,
		Ondersteuning:  aanmelding.Ondersteuning,
		Bijzonderheden: aanmelding.Bijzonderheden,
		Terms:          aanmelding.Terms,
		Status:         "nieuw",  // Standaard status
		TestMode:       testMode, // Neem test mode over
	}

	// Sla de aanmelding op in de database (niet in test modus)
	if !testMode {
		logger.Info("Aanmelding opslaan in database",
			"naam", nieuweAanmelding.Naam,
			"email", nieuweAanmelding.Email)
		ctx := c.Context()
		if err := h.aanmeldingRepo.Create(ctx, nieuweAanmelding); err != nil {
			logger.Error("Fout bij opslaan aanmelding in database",
				"error", err,
				"naam", nieuweAanmelding.Naam,
				"email", nieuweAanmelding.Email,
				"elapsed", time.Since(start))
			// Geef een fout terug als opslaan mislukt
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Kon aanmelding niet opslaan: " + err.Error(),
			})
		}
		logger.Info("Aanmelding succesvol opgeslagen in database",
			"id", nieuweAanmelding.ID,
			"naam", nieuweAanmelding.Naam)
	} else {
		logger.Info("Test modus: Aanmelding niet opgeslagen in database",
			"naam", nieuweAanmelding.Naam)
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

	// Stuur een notificatie voor een nieuwe aanmelding
	h.sendAanmeldingNotification(&aanmelding, testMode)

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

// Stuur een notificatie voor een nieuw contactformulier
func (h *EmailHandler) sendContactNotification(contact *models.ContactFormulier, isTestMode bool) {
	// Skip als de notification service niet beschikbaar is of als we in test mode zijn
	if h.notificationService == nil || isTestMode {
		return
	}

	priority := models.NotificationPriorityMedium
	title := "Nieuw Contactverzoek"
	message := ""

	if contact.Naam != "" && contact.Email != "" {
		message = "<b>" + contact.Naam + "</b> heeft contact opgenomen via het contactformulier.\n\n" +
			"<b>Email:</b> " + contact.Email + "\n\n" +
			"<b>Bericht:</b>\n" + contact.Bericht
	} else {
		// Fallback als naam of email ontbreekt
		message = "Er is een nieuw contactverzoek ontvangen.\n\n" +
			"<b>Bericht:</b>\n" + contact.Bericht
	}

	// Maak een notificatie aan
	_, err := h.notificationService.CreateNotification(
		context.Background(),
		models.NotificationTypeContact,
		priority,
		title,
		message,
	)

	if err != nil {
		logger.Error("Fout bij aanmaken contact notificatie",
			"error", err,
			"contact_naam", contact.Naam,
			"contact_email", contact.Email)
	}
}

// Stuur een notificatie voor een nieuwe aanmelding
func (h *EmailHandler) sendAanmeldingNotification(aanmelding *models.AanmeldingFormulier, isTestMode bool) {
	// Skip als de notification service niet beschikbaar is of als we in test mode zijn
	if h.notificationService == nil || isTestMode {
		return
	}

	priority := models.NotificationPriorityMedium
	title := "Nieuwe Aanmelding"
	message := ""

	if aanmelding.Naam != "" && aanmelding.Email != "" {
		message = "<b>" + aanmelding.Naam + "</b> heeft zich aangemeld.\n\n" +
			"<b>Email:</b> " + aanmelding.Email + "\n\n" +
			"<b>Rol:</b> " + aanmelding.Rol + "\n" +
			"<b>Afstand:</b> " + aanmelding.Afstand + "\n"

		if aanmelding.Telefoon != "" {
			message += "<b>Telefoon:</b> " + aanmelding.Telefoon + "\n\n"
		}

		if aanmelding.Bijzonderheden != "" {
			message += "<b>Bijzonderheden:</b>\n" + aanmelding.Bijzonderheden
		}
	} else {
		// Fallback als naam of email ontbreekt
		message = "Er is een nieuwe aanmelding ontvangen."
	}

	// Maak een notificatie aan
	_, err := h.notificationService.CreateNotification(
		context.Background(),
		models.NotificationTypeAanmelding,
		priority,
		title,
		message,
	)

	if err != nil {
		logger.Error("Fout bij aanmaken aanmelding notificatie",
			"error", err,
			"aanmelding_naam", aanmelding.Naam,
			"aanmelding_email", aanmelding.Email)
	}
}
