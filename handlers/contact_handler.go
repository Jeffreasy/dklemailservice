package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ContactHandler bevat handlers voor contact formulier beheer
type ContactHandler struct {
	contactRepo         repository.ContactRepository
	contactAntwoordRepo repository.ContactAntwoordRepository
	emailService        *services.EmailService
	authService         services.AuthService
	notificationService services.NotificationService
}

// NewContactHandler maakt een nieuwe contact handler
func NewContactHandler(
	contactRepo repository.ContactRepository,
	contactAntwoordRepo repository.ContactAntwoordRepository,
	emailService *services.EmailService,
	authService services.AuthService,
	notificationService services.NotificationService,
) *ContactHandler {
	return &ContactHandler{
		contactRepo:         contactRepo,
		contactAntwoordRepo: contactAntwoordRepo,
		emailService:        emailService,
		authService:         authService,
		notificationService: notificationService,
	}
}

// RegisterRoutes registreert de routes voor contact beheer
func (h *ContactHandler) RegisterRoutes(app *fiber.App) {
	// Groep voor contact beheer routes (vereist admin rechten)
	contactGroup := app.Group("/api/contact")
	contactGroup.Use(AuthMiddleware(h.authService))
	contactGroup.Use(AdminMiddleware(h.authService))

	// Contact beheer routes
	contactGroup.Get("/", h.ListContactFormulieren)
	contactGroup.Get("/:id", h.GetContactFormulier)
	contactGroup.Put("/:id", h.UpdateContactFormulier)
	contactGroup.Delete("/:id", h.DeleteContactFormulier)
	contactGroup.Post("/:id/antwoord", h.AddContactAntwoord)
	contactGroup.Get("/status/:status", h.GetContactFormulierenByStatus)
}

// ListContactFormulieren haalt een lijst van contactformulieren op
// @Summary Lijst van contactformulieren ophalen
// @Description Haalt een gepagineerde lijst van contactformulieren op
// @Tags Contact
// @Accept json
// @Produce json
// @Param limit query int false "Aantal resultaten per pagina (standaard 10)"
// @Param offset query int false "Offset voor paginering (standaard 0)"
// @Success 200 {array} models.ContactFormulier
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/contact [get]
// @Security BearerAuth
func (h *ContactHandler) ListContactFormulieren(c *fiber.Ctx) error {
	// Haal query parameters op
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Valideer parameters
	if limit < 1 || limit > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Limit moet tussen 1 en 100 liggen",
		})
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Offset mag niet negatief zijn",
		})
	}

	// Haal contactformulieren op
	ctx := c.Context()
	contacts, err := h.contactRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen contactformulieren", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulieren niet ophalen",
		})
	}

	// Stuur resultaat terug
	return c.JSON(contacts)
}

// GetContactFormulier haalt een specifiek contactformulier op
// @Summary Details van een specifiek contactformulier ophalen
// @Description Haalt de details van een specifiek contactformulier op, inclusief antwoorden
// @Tags Contact
// @Accept json
// @Produce json
// @Param id path string true "Contact ID"
// @Success 200 {object} models.ContactFormulier
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/contact/{id} [get]
// @Security BearerAuth
func (h *ContactHandler) GetContactFormulier(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Haal contactformulier op
	ctx := c.Context()
	contact, err := h.contactRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen contactformulier", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulier niet ophalen",
		})
	}

	if contact == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Contactformulier niet gevonden",
		})
	}

	// Haal antwoorden op
	antwoorden, err := h.contactAntwoordRepo.ListByContactID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen antwoorden", "error", err, "contact_id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon antwoorden niet ophalen",
		})
	}

	// Converteer []*models.ContactAntwoord naar []models.ContactAntwoord
	contactAntwoorden := make([]models.ContactAntwoord, len(antwoorden))
	for i, antwoord := range antwoorden {
		if antwoord != nil {
			contactAntwoorden[i] = *antwoord
		}
	}

	// Voeg antwoorden toe aan contactformulier
	contact.Antwoorden = contactAntwoorden

	// Stuur resultaat terug
	return c.JSON(contact)
}

// UpdateContactFormulier werkt een contactformulier bij
// @Summary Contactformulier bijwerken
// @Description Werkt een bestaand contactformulier bij (status, notities)
// @Tags Contact
// @Accept json
// @Produce json
// @Param id path string true "Contact ID"
// @Param contact body models.ContactFormulier true "Contact gegevens"
// @Success 200 {object} models.ContactFormulier
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/contact/{id} [put]
// @Security BearerAuth
func (h *ContactHandler) UpdateContactFormulier(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Haal gebruiker op uit context
	gebruiker, ok := c.Locals("gebruiker").(*models.Gebruiker)
	if !ok || gebruiker == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon gebruiker niet ophalen uit context",
		})
	}

	// Haal bestaand contactformulier op
	ctx := c.Context()
	contact, err := h.contactRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen contactformulier", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulier niet ophalen",
		})
	}

	if contact == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Contactformulier niet gevonden",
		})
	}

	// Haal update gegevens op uit request body
	var updateData struct {
		Status   string  `json:"status"`
		Notities *string `json:"notities"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	// Update contactformulier
	if updateData.Status != "" {
		contact.Status = updateData.Status
	}

	if updateData.Notities != nil {
		contact.Notities = updateData.Notities
	}

	// Stel behandeld door en behandeld op in als dit nog niet is gedaan
	if contact.BehandeldDoor == nil || *contact.BehandeldDoor == "" {
		contact.BehandeldDoor = &gebruiker.Email
		now := time.Now()
		contact.BehandeldOp = &now
	}

	// Sla wijzigingen op
	if err := h.contactRepo.Update(ctx, contact); err != nil {
		logger.Error("Fout bij bijwerken contactformulier", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulier niet bijwerken",
		})
	}

	// Stuur bijgewerkt contactformulier terug
	return c.JSON(contact)
}

// DeleteContactFormulier verwijdert een contactformulier
// @Summary Contactformulier verwijderen
// @Description Verwijdert een contactformulier en bijbehorende antwoorden
// @Tags Contact
// @Accept json
// @Produce json
// @Param id path string true "Contact ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/contact/{id} [delete]
// @Security BearerAuth
func (h *ContactHandler) DeleteContactFormulier(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Controleer of contactformulier bestaat
	ctx := c.Context()
	contact, err := h.contactRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen contactformulier", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulier niet ophalen",
		})
	}

	if contact == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Contactformulier niet gevonden",
		})
	}

	// Verwijder contactformulier
	if err := h.contactRepo.Delete(ctx, id); err != nil {
		logger.Error("Fout bij verwijderen contactformulier", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulier niet verwijderen",
		})
	}

	// Stuur bevestiging terug
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Contactformulier succesvol verwijderd",
	})
}

// AddContactAntwoord voegt een antwoord toe aan een contactformulier
// @Summary Antwoord toevoegen aan contactformulier
// @Description Voegt een nieuw antwoord toe aan een contactformulier
// @Tags Contact
// @Accept json
// @Produce json
// @Param id path string true "Contact ID"
// @Param antwoord body models.ContactAntwoord true "Antwoord gegevens"
// @Success 200 {object} models.ContactAntwoord
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/contact/{id}/antwoord [post]
// @Security BearerAuth
func (h *ContactHandler) AddContactAntwoord(c *fiber.Ctx) error {
	// Haal ID op uit URL
	contactID := c.Params("id")
	if contactID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Contact ID is verplicht",
		})
	}

	// Haal gebruiker op uit context
	gebruiker, ok := c.Locals("gebruiker").(*models.Gebruiker)
	if !ok || gebruiker == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon gebruiker niet ophalen uit context",
		})
	}

	// Controleer of contactformulier bestaat
	ctx := c.Context()
	contact, err := h.contactRepo.GetByID(ctx, contactID)
	if err != nil {
		logger.Error("Fout bij ophalen contactformulier", "error", err, "id", contactID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulier niet ophalen",
		})
	}

	if contact == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Contactformulier niet gevonden",
		})
	}

	// Haal antwoord gegevens op uit request body
	var antwoordData struct {
		Tekst string `json:"tekst"`
	}

	if err := c.BodyParser(&antwoordData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	// Valideer gegevens
	if antwoordData.Tekst == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tekst is verplicht",
		})
	}

	// Maak nieuw antwoord
	antwoord := &models.ContactAntwoord{
		ContactID:      contactID,
		Tekst:          antwoordData.Tekst,
		VerzondDoor:    gebruiker.Email,
		EmailVerzonden: false,
	}

	// Sla antwoord op
	if err := h.contactAntwoordRepo.Create(ctx, antwoord); err != nil {
		logger.Error("Fout bij opslaan antwoord", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon antwoord niet opslaan",
		})
	}

	// Update contactformulier status naar beantwoord
	contact.Status = "beantwoord"
	contact.Beantwoord = true
	contact.AntwoordTekst = antwoordData.Tekst
	now := time.Now()
	contact.AntwoordDatum = &now
	contact.AntwoordDoor = gebruiker.Email

	if err := h.contactRepo.Update(ctx, contact); err != nil {
		logger.Error("Fout bij bijwerken contactformulier", "error", err, "id", contactID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulier niet bijwerken",
		})
	}

	// Stuur e-mail met antwoord (in de achtergrond)
	go func() {
		// Gebruik de juiste methode voor het verzenden van e-mail
		if err := h.emailService.SendEmail(contact.Email, "Antwoord op uw contactformulier", antwoordData.Tekst); err != nil {
			logger.Error("Fout bij verzenden antwoord e-mail", "error", err, "contact_id", contactID)
		} else {
			// Update e-mail verzonden status
			antwoord.EmailVerzonden = true
			bgCtx := context.Background()
			if err := h.contactAntwoordRepo.Update(bgCtx, antwoord); err != nil {
				logger.Error("Fout bij bijwerken antwoord e-mail status", "error", err, "antwoord_id", antwoord.ID)
			}
		}
	}()

	// Stuur antwoord terug
	return c.JSON(antwoord)
}

// GetContactFormulierenByStatus haalt contactformulieren op basis van status op
// @Summary Contactformulieren filteren op status
// @Description Haalt een lijst van contactformulieren op gefilterd op status
// @Tags Contact
// @Accept json
// @Produce json
// @Param status path string true "Status (nieuw, in_behandeling, beantwoord, gesloten)"
// @Success 200 {array} models.ContactFormulier
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/contact/status/{status} [get]
// @Security BearerAuth
func (h *ContactHandler) GetContactFormulierenByStatus(c *fiber.Ctx) error {
	// Haal status op uit URL
	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is verplicht",
		})
	}

	// Valideer status
	validStatuses := map[string]bool{
		"nieuw":          true,
		"in_behandeling": true,
		"beantwoord":     true,
		"gesloten":       true,
	}

	if !validStatuses[status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige status",
		})
	}

	// Haal contactformulieren op
	ctx := c.Context()
	contacts, err := h.contactRepo.FindByStatus(ctx, status)
	if err != nil {
		logger.Error("Fout bij ophalen contactformulieren op status", "error", err, "status", status)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon contactformulieren niet ophalen",
		})
	}

	// Stuur resultaat terug
	return c.JSON(contacts)
}

// handleContactNotification stuurt een notificatie voor een nieuw contactformulier
func (h *ContactHandler) handleContactNotification(ctx context.Context, contact *models.ContactFormulier) {
	// Skip als de notification service niet beschikbaar is
	if h.notificationService == nil {
		return
	}

	priority := models.NotificationPriorityMedium
	title := "Nieuw Contactverzoek"
	message := fmt.Sprintf(
		"<b>%s</b> heeft contact opgenomen.\n\n"+
			"<b>Email:</b> %s\n\n"+
			"<b>Bericht:</b>\n%s",
		contact.Naam,
		contact.Email,
		contact.Bericht,
	)

	// Maak een notificatie aan
	_, err := h.notificationService.CreateNotification(
		ctx,
		models.NotificationTypeContact,
		priority,
		title,
		message,
	)

	if err != nil {
		logger.Error("Fout bij aanmaken contact notificatie",
			"error", err,
			"contact_id", contact.ID)
	}
}
