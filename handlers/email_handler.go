package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// EmailServiceInterface definieert de interface voor email operaties
type EmailServiceInterface interface {
	SendContactEmail(data *models.ContactEmailData) error
	// GEWIJZIGD: Gebruikt nu de nieuwe struct
	SendRegistrationEmail(data *models.RegistrationEmailData) error
}

// EmailHandler verzorgt de afhandeling van email verzoeken
type EmailHandler struct {
	emailService        EmailServiceInterface
	notificationService services.NotificationService
	participantRepo     repository.ParticipantRepository
	eventRegRepo        repository.EventRegistrationRepository
	eventRepo           repository.EventRepository
}

// NewEmailHandler maakt een nieuwe EmailHandler
func NewEmailHandler(
	emailService EmailServiceInterface,
	notificationService services.NotificationService,
	participantRepo repository.ParticipantRepository,
	eventRegRepo repository.EventRegistrationRepository,
	eventRepo repository.EventRepository,
) *EmailHandler {
	return &EmailHandler{
		emailService:        emailService,
		notificationService: notificationService,
		participantRepo:     participantRepo,
		eventRegRepo:        eventRegRepo,
		eventRepo:           eventRepo,
	}
}

// HandleContactEmail is ongewijzigd gebleven.
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
		if val, ok := requestMap["test_mode"]; ok {
			if bVal, ok := val.(bool); ok && bVal {
				testMode = true
				logger.Info("Test modus gedetecteerd via body parameter", "remote_ip", c.IP())
			}
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

// HandleRegistrationEmail (voorheen HandleAanmeldingEmail)
// @Summary Nieuwe participant en event registratie aanmaken
// @Tags Public
// @Router /api/register [post]
func (h *EmailHandler) HandleRegistrationEmail(c *fiber.Ctx) error {
	var form models.AanmeldingFormulier // We hergebruiken de DTO, dat is prima
	start := time.Now()
	ctx := c.Context()

	if err := c.BodyParser(&form); err != nil {
		logger.Error("Fout bij parsen van registratie formulier", "error", err, "remote_ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Ongeldig verzoek: " + err.Error()})
	}

	// --- Detecteer test modus (zelfde als hierboven) ---
	var testMode bool
	if c.Get("X-Test-Mode") == "true" {
		testMode = true
	}
	if c.Locals("test_mode") != nil {
		testMode = true
	}
	var requestMap map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestMap); err == nil {
		if val, ok := requestMap["test_mode"]; ok {
			if bVal, ok := val.(bool); ok && bVal {
				testMode = true
			}
		}
	}
	// --- Einde test modus detectie ---

	logger.Info("Registratie formulier ontvangen", "naam", form.Naam, "email", form.Email, "remote_ip", c.IP(), "test_mode", testMode)

	// --- Validatie (grotendeels ongewijzigd) ---
	if form.Naam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Naam is verplicht"})
	}
	if form.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Email is verplicht"})
	}
	if !strings.Contains(form.Email, "@") || !strings.Contains(form.Email, ".") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Ongeldig email adres"})
	}
	if !form.Terms {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Je moet akkoord gaan met de voorwaarden"})
	}

	// --- NIEUWE LOGICA (V28 REFACTOR) ---

	// 1. Zoek of maak de "Persoon" (Participant)
	var participant *models.Participant
	var err error

	// Probeer participant op email te vinden
	participants, err := h.participantRepo.FindByEmail(ctx, form.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Fout bij zoeken naar bestaande participant", "error", err, "email", form.Email)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Databasefout bij zoeken participant"})
	}

	if len(participants) > 0 {
		// Participant bestaat al, we gebruiken deze
		participant = participants[0]
		logger.Info("Bestaande participant gevonden", "id", participant.ID, "email", participant.Email)
	} else {
		// Participant bestaat niet, maak een nieuwe aan
		logger.Info("Nieuwe participant aanmaken", "email", form.Email)
		participant = &models.Participant{
			Naam:     form.Naam,
			Email:    form.Email,
			Telefoon: form.Telefoon,
			Terms:    form.Terms,
		}

		if !testMode {
			if err := h.participantRepo.Create(ctx, participant); err != nil {
				logger.Error("Fout bij opslaan nieuwe participant", "error", err, "email", form.Email)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Kon participant niet opslaan"})
			}
			logger.Info("Nieuwe participant succesvol opgeslagen", "id", participant.ID)
		} else {
			logger.Info("Test modus: Participant niet opgeslagen", "email", form.Email)
		}
	}

	// 2. Haal het actieve evenement op
	activeEvent, err := h.eventRepo.GetActiveEvent(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Geen actief evenement gevonden. Registratie is gesloten.", "error", err)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"success": false, "error": "De inschrijving is momenteel niet geopend."})
		}
		logger.Error("Fout bij ophalen actief evenement", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Fout bij ophalen event data"})
	}

	// 3. Maak de "Deelname" (EventRegistration) aan
	distanceRoute := form.Afstand
	roleName := form.Rol

	nieuweRegistratie := &models.EventRegistration{
		EventID:             activeEvent.ID,
		ParticipantID:       participant.ID,
		TestMode:            testMode,
		Ondersteuning:       form.Ondersteuning,
		Bijzonderheden:      form.Bijzonderheden,
		Status:              "nieuw", // V27: Direct field (database column is 'status')
		DistanceRoute:       &distanceRoute,
		ParticipantRoleName: &roleName,
		// Standaardwaarden (Steps, TotalDistance etc.) worden door DB default gezet
	}

	// Sla de registratie op
	if !testMode {
		// DIT IS DE FIX VOOR FOUT 2 (h.eventRegRepo.Create undefined)
		// We roepen nu de 'Create' methode aan die we zojuist hebben gedefinieerd.
		if err := h.eventRegRepo.Create(ctx, nieuweRegistratie); err != nil {
			logger.Error("Fout bij opslaan event registratie", "error", err, "participant_id", participant.ID, "event_id", activeEvent.ID)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Kon registratie niet opslaan"})
		}
		logger.Info("Event registratie succesvol opgeslagen", "id", nieuweRegistratie.ID, "participant_id", participant.ID)
	} else {
		logger.Info("Test modus: Event registratie niet opgeslagen", "participant_email", form.Email)
	}

	// --- Einde V28 Refactor Logica ---

	// 4. Verzend e-mails
	adminEmail := os.Getenv("REGISTRATION_EMAIL")
	if adminEmail == "" {
		adminEmail = "inschrijving@dekoninklijkeloop.nl"
		logger.Warn("REGISTRATION_EMAIL niet geconfigureerd, gebruik standaardwaarde", "default", adminEmail)
	}

	// DIT IS DE FIX VOOR FOUT 3 (undefined: models.RegistrationEmailData)
	// We gebruiken nu de correct gedefinieerde struct
	emailData := &models.RegistrationEmailData{
		Participant:  participant,
		Registration: nieuweRegistratie,
		AdminEmail:   adminEmail,
	}

	// Stuur email naar admin
	emailData.ToAdmin = true
	if testMode {
		logger.Info("Test modus: Geen admin email verzonden", "admin_email", adminEmail)
	} else {
		logger.Info("Admin email wordt verzonden", "admin_email", adminEmail, "naam", form.Naam)
		if err := h.emailService.SendRegistrationEmail(emailData); err != nil {
			logger.Error("Fout bij verzenden admin email", "error", err, "elapsed", time.Since(start))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Fout bij het verzenden van de email"})
		}
		logger.Info("Admin email verzonden", "admin_email", adminEmail, "elapsed", time.Since(start))
	}

	// Stuur bevestigingsemail naar gebruiker
	emailData.ToAdmin = false
	if testMode {
		logger.Info("Test modus: Geen gebruiker email verzonden", "user_email", form.Email)
	} else {
		logger.Info("Bevestigingsemail wordt verzonden", "user_email", form.Email, "naam", form.Naam)
		if err := h.emailService.SendRegistrationEmail(emailData); err != nil {
			logger.Error("Fout bij verzenden bevestigingsemail", "error", err, "elapsed", time.Since(start))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Fout bij het verzenden van de bevestigingsemail"})
		}
		logger.Info("Bevestigingsemail verzonden", "user_email", form.Email, "elapsed", time.Since(start))

		// 5. Update de *registratie* (niet de participant) met verzendstatus
		now := time.Now()
		nieuweRegistratie.EmailVerzonden = true
		nieuweRegistratie.EmailVerzondenOp = &now
		if err := h.eventRegRepo.Update(ctx, nieuweRegistratie); err != nil {
			logger.Error("Fout bij bijwerken registratie na email verzending", "error", err, "reg_id", nieuweRegistratie.ID)
		} else {
			logger.Info("Registratie bijgewerkt met email verzendstatus", "reg_id", nieuweRegistratie.ID)
		}
	}

	// 6. Stuur een notificatie
	h.sendRegistrationNotification(&form, testMode) // De DTO is hier prima

	// 7. Return success
	if testMode {
		logger.Info("Registratie formulier succesvol verwerkt in test modus", "naam", form.Naam, "total_elapsed", time.Since(start))
		return c.JSON(fiber.Map{
			"success":   true,
			"message":   "[TEST MODE] Je aanmelding is verwerkt (geen echte email verzonden).",
			"test_mode": true,
		})
	} else {
		logger.Info("Registratie formulier succesvol verwerkt", "naam", form.Naam, "total_elapsed", time.Since(start))
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

// sendContactNotification - GEWIJZIGD om string key te gebruiken
func (h *EmailHandler) sendContactNotification(contact *models.ContactFormulier, isTestMode bool) {
	if h.notificationService == nil || isTestMode {
		return
	}
	priority := models.NotificationPriorityMedium
	title := "Nieuw Contactverzoek"
	message := "<b>" + contact.Naam + "</b> heeft contact opgenomen via het contactformulier.\n\n" +
		"<b>Email:</b> " + contact.Email + "\n\n" +
		"<b>Bericht:</b>\n" + contact.Bericht

	// GEWIJZIGD: Gebruik V27 string key
	_, err := h.notificationService.CreateNotification(
		context.Background(),
		"contact", // voorheen models.NotificationTypeContact
		priority,
		title,
		message,
	)
	if err != nil {
		logger.Error("Fout bij aanmaken contact notificatie", "error", err, "contact_naam", contact.Naam)
	}
}

// sendRegistrationNotification - GEWIJZIGD om string key te gebruiken
func (h *EmailHandler) sendRegistrationNotification(aanmelding *models.AanmeldingFormulier, isTestMode bool) {
	if h.notificationService == nil || isTestMode {
		return
	}

	priority := models.NotificationPriorityMedium
	title := "Nieuwe Aanmelding"
	message := "<b>" + aanmelding.Naam + "</b> heeft zich aangemeld.\n\n" +
		"<b>Email:</b> " + aanmelding.Email + "\n\n" +
		"<b>Rol:</b> " + aanmelding.Rol + "\n" +
		"<b>Afstand:</b> " + aanmelding.Afstand + "\n"

	if aanmelding.Telefoon != "" {
		message += "<b>Telefoon:</b> " + aanmelding.Telefoon + "\n\n"
	}
	if aanmelding.Bijzonderheden != "" {
		message += "<b>Bijzonderheden:</b>\n" + aanmelding.Bijzonderheden
	}

	// DIT IS DE FIX VOOR FOUT 4 (undefined: models.NotificationTypeRegistration)
	_, err := h.notificationService.CreateNotification(
		context.Background(),
		"registration", // voorheen models.NotificationTypeRegistration
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
