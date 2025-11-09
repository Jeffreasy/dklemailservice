package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ParticipantHandler bevat handlers voor participant (persoon) beheer.
// LET OP: Logica voor Status, Rol, Notities, etc. is VERPLAATST
// naar de (nog te maken) EventRegistrationHandler.
type ParticipantHandler struct {
	participantRepo         repository.ParticipantRepository
	participantAntwoordRepo repository.ParticipantAntwoordRepository
	emailService            *services.EmailService
	authService             services.AuthService
	permissionService       services.PermissionService

	// EventRegistration repo voor status updates bij antwoorden
	eventRegRepo repository.EventRegistrationRepository
}

// NewParticipantHandler maakt een nieuwe participant handler
func NewParticipantHandler(
	participantRepo repository.ParticipantRepository,
	participantAntwoordRepo repository.ParticipantAntwoordRepository,
	emailService *services.EmailService,
	authService services.AuthService,
	permissionService services.PermissionService,
	eventRegRepo repository.EventRegistrationRepository,
) *ParticipantHandler {
	return &ParticipantHandler{
		participantRepo:         participantRepo,
		participantAntwoordRepo: participantAntwoordRepo,
		emailService:            emailService,
		authService:             authService,
		permissionService:       permissionService,
		eventRegRepo:            eventRegRepo,
	}
}

// RegisterRoutes registreert de routes voor participant beheer
func (h *ParticipantHandler) RegisterRoutes(app *fiber.App) {
	// API-groep voor /api/participant
	// De 'resource' naam in permissies is nu 'participant'
	participantApi := app.Group("/api/participant", AuthMiddleware(h.authService))

	participantApi.Get("/",
		PermissionMiddleware(h.permissionService, "participant", "read"),
		h.ListParticipants)

	// GET /api/participant/emails - Specifieke route VOOR generieke /:id route
	// Dit voorkomt dat /emails wordt gematcht als /:id met id="emails"
	participantApi.Get("/emails",
		PermissionMiddleware(h.permissionService, "participant", "read"),
		h.GetParticipantEmails)

	participantApi.Get("/:id",
		PermissionMiddleware(h.permissionService, "participant", "read"),
		h.GetParticipant)

	participantApi.Post("/:id/antwoord",
		PermissionMiddleware(h.permissionService, "participant", "write"),
		h.AddParticipantAntwoord)

	participantApi.Delete("/:id",
		PermissionMiddleware(h.permissionService, "participant", "delete"),
		h.DeleteParticipant)

	// --- VEROUDERDE ROUTES ZIJN VERWIJDERD ---
	// GET /api/participant/rol/:rol - Is verplaatst naar EventRegistrationHandler
	// PUT /api/participant/:id - Is verplaatst naar EventRegistrationHandler

	// Plural routes (/api/participants) - aliases
	participantsApi := app.Group("/api/participants", AuthMiddleware(h.authService))

	participantsApi.Get("/",
		PermissionMiddleware(h.permissionService, "participant", "read"),
		h.ListParticipants)

	// Specifieke /emails route ook voor plural alias
	participantsApi.Get("/emails",
		PermissionMiddleware(h.permissionService, "participant", "read"),
		h.GetParticipantEmails)

	participantsApi.Get("/:id",
		PermissionMiddleware(h.permissionService, "participant", "read"),
		h.GetParticipant)

	participantsApi.Post("/:id/antwoord",
		PermissionMiddleware(h.permissionService, "participant", "write"),
		h.AddParticipantAntwoord)

	participantsApi.Delete("/:id",
		PermissionMiddleware(h.permissionService, "participant", "delete"),
		h.DeleteParticipant)
}

// ListParticipants haalt een lijst van participants op
// @Summary Lijst van participants ophalen
// @Description Haalt een gepagineerde lijst van participants (personen) op
// @Tags Participant
// @Accept json
// @Produce json
// @Param limit query int false "Aantal resultaten per pagina (standaard 10)"
// @Param offset query int false "Offset voor paginering (standaard 0)"
// @Success 200 {array} models.Participant
// @Router /api/participant [get]
// @Security BearerAuth
func (h *ParticipantHandler) ListParticipants(c *fiber.Ctx) error {
	// Haal query parameters op
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Valideer parameters
	if limit < 1 || limit > 10000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Limit moet tussen 1 en 10000 liggen",
			"code":  "INVALID_LIMIT",
		})
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Offset mag niet negatief zijn",
			"code":  "INVALID_OFFSET",
		})
	}

	// Haal participants op
	ctx := c.Context()
	participants, err := h.participantRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen participants", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon participants niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	// Stuur resultaat terug
	return c.JSON(participants)
}

// GetParticipant haalt een specifieke participant op
// @Summary Details van een specifieke participant ophalen
// @Description Haalt de details van een specifieke participant op, inclusief antwoorden
// @Tags Participant
// @Accept json
// @Produce json
// @Param id path string true "Participant ID"
// @Success 200 {object} models.Participant
// @Router /api/participant/{id} [get]
// @Security BearerAuth
func (h *ParticipantHandler) GetParticipant(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
			"code":  "MISSING_ID",
		})
	}

	// Haal participant op
	ctx := c.Context()
	participant, err := h.participantRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen participant", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon participant niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	if participant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Participant niet gevonden",
			"code":  "NOT_FOUND",
		})
	}

	// Haal antwoorden op
	antwoorden, err := h.participantAntwoordRepo.ListByParticipantID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen antwoorden", "error", err, "participant_id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon antwoorden niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	// Converteer []*models.ParticipantAntwoord naar []models.ParticipantAntwoord
	participantAntwoorden := make([]models.ParticipantAntwoord, len(antwoorden))
	for i, antwoord := range antwoorden {
		if antwoord != nil {
			participantAntwoorden[i] = *antwoord
		}
	}

	// Voeg antwoorden toe aan participant
	participant.Antwoorden = participantAntwoorden

	// Haal ook de EventRegistrations van de participant op en mee te sturen
	registraties, err := h.eventRegRepo.ListByParticipantID(ctx, id)
	if err != nil {
		logger.Warn("Kon event registrations niet ophalen", "participant_id", id, "error", err)
		// Niet fatale fout - participant data is nog steeds geldig
	}

	// Voeg event registrations toe aan response
	response := fiber.Map{
		"participant": participant,
	}
	if registraties != nil {
		response["event_registrations"] = registraties
	}

	// Stuur resultaat terug
	return c.JSON(response)
}

// DeleteParticipant verwijdert een participant
// @Summary Participant verwijderen
// @Description Verwijdert een participant en bijbehorende antwoorden (en via CASCADE ook EventRegistrations)
// @Tags Participant
// @Accept json
// @Produce json
// @Param id path string true "Participant ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/participant/{id} [delete]
// @Security BearerAuth
func (h *ParticipantHandler) DeleteParticipant(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
			"code":  "MISSING_ID",
		})
	}

	// Controleer of participant bestaat
	ctx := c.Context()
	participant, err := h.participantRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen participant", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon participant niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	if participant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Participant niet gevonden",
			"code":  "NOT_FOUND",
		})
	}

	// Verwijder participant
	if err := h.participantRepo.Delete(ctx, id); err != nil {
		logger.Error("Fout bij verwijderen participant", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon participant niet verwijderen",
			"code":  "DELETION_FAILED",
		})
	}

	// Stuur bevestiging terug
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Participant succesvol verwijderd",
	})
}

// AddParticipantAntwoord voegt een antwoord toe aan een participant
// @Summary Antwoord toevoegen aan participant
// @Description Voegt een nieuw antwoord toe aan een participant
// @Tags Participant
// @Accept json
// @Produce json
// @Param id path string true "Participant ID"
// @Param antwoord body models.ParticipantAntwoord true "Antwoord gegevens"
// @Success 200 {object} models.ParticipantAntwoord
// @Router /api/participant/{id}/antwoord [post]
// @Security BearerAuth
func (h *ParticipantHandler) AddParticipantAntwoord(c *fiber.Ctx) error {
	// Haal ID op uit URL
	participantID := c.Params("id")
	if participantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Participant ID is verplicht",
			"code":  "MISSING_ID",
		})
	}

	// Haal gebruiker op uit context
	gebruiker, ok := c.Locals("gebruiker").(*models.Gebruiker)
	if !ok || gebruiker == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon gebruiker niet ophalen uit context",
			"code":  "INTERNAL_ERROR",
		})
	}

	// Controleer of participant bestaat
	ctx := c.Context()
	participant, err := h.participantRepo.GetByID(ctx, participantID)
	if err != nil {
		logger.Error("Fout bij ophalen participant", "error", err, "id", participantID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon participant niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	if participant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Participant niet gevonden",
			"code":  "NOT_FOUND",
		})
	}

	// Haal antwoord gegevens op uit request body
	var antwoordData struct {
		Tekst string `json:"tekst"`
	}

	if err := c.BodyParser(&antwoordData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
			"code":  "INVALID_INPUT",
		})
	}

	// Valideer gegevens
	if antwoordData.Tekst == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tekst is verplicht",
			"code":  "MISSING_TEXT",
		})
	}

	// Maak nieuw antwoord
	antwoord := &models.ParticipantAntwoord{
		ParticipantID:  participantID,
		Tekst:          antwoordData.Tekst,
		VerzondDoor:    gebruiker.Email,
		EmailVerzonden: false,
	}

	// Sla antwoord op
	if err := h.participantAntwoordRepo.Create(ctx, antwoord); err != nil {
		logger.Error("Fout bij opslaan antwoord", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon antwoord niet opslaan",
			"code":  "SAVE_FAILED",
		})
	}

	// --- LOGICA GEIMPLEMENTEERD ---
	// De status "beantwoord" is nu onderdeel van de EventRegistration, niet de Participant.
	// We zoeken de actieve registratie van deze participant en werken die bij.
	reg, err := h.eventRegRepo.GetActiveRegistrationForParticipant(ctx, participantID)
	if err != nil {
		logger.Error("Kon actieve registratie niet vinden om status bij te werken", "participant_id", participantID, "error", err)
	} else if reg != nil {
		reg.Status = "beantwoord" // V27: Direct field (database column is 'status')
		reg.BehandeldDoor = &gebruiker.Email
		now := time.Now()
		reg.BehandeldOp = &now

		if err := h.eventRegRepo.Update(ctx, reg); err != nil {
			logger.Error("Kon status van event registratie niet bijwerken", "reg_id", reg.ID, "error", err)
		}
	} else {
		logger.Warn("Geen actieve registratie gevonden om status bij te werken na antwoord", "participant_id", participantID)
	}

	// Stuur e-mail met antwoord (in de achtergrond)
	go func() {
		if err := h.emailService.SendEmail(participant.Email, "Antwoord op uw registratie", antwoordData.Tekst); err != nil {
			logger.Error("Fout bij verzenden antwoord e-mail", "error", err, "participant_id", participantID)
		} else {
			// Update e-mail verzonden status
			antwoord.EmailVerzonden = true
			bgCtx := context.Background()
			if err := h.participantAntwoordRepo.Update(bgCtx, antwoord); err != nil {
				logger.Error("Fout bij bijwerken antwoord e-mail status", "error", err, "antwoord_id", antwoord.ID)
			}
		}
	}()

	// Stuur antwoord terug
	return c.JSON(antwoord)
}

// GetParticipantEmails haalt alle participant emails op inclusief systeem emails
// @Summary Alle emails ophalen (participants + systeem)
// @Description Haalt een lijst van alle participant emails en systeem emails op voor admin gebruik
// @Tags Participant
// @Accept json
// @Produce json
// @Success 200 {object} object{participant_emails=array, system_emails=array, all_emails=array}
// @Router /api/participant/emails [get]
// @Security BearerAuth
func (h *ParticipantHandler) GetParticipantEmails(c *fiber.Ctx) error {
	// Haal alle participants op
	ctx := c.Context()
	participants, err := h.participantRepo.List(ctx, 10000, 0) // Max limit om alle te krijgen
	if err != nil {
		logger.Error("Fout bij ophalen participants voor emails", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon participant emails niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	// Extract participant emails
	participantEmails := make([]string, len(participants))
	for i, participant := range participants {
		participantEmails[i] = participant.Email
	}

	// Haal systeem emails op uit environment variables
	systemEmails := []string{}

	if infoEmail := os.Getenv("INFO_EMAIL"); infoEmail != "" {
		systemEmails = append(systemEmails, infoEmail)
	}

	if inschrijvingEmail := os.Getenv("INSCHRIJVING_EMAIL"); inschrijvingEmail != "" {
		systemEmails = append(systemEmails, inschrijvingEmail)
	}

	// Combineer alle emails
	allEmails := append(participantEmails, systemEmails...)

	// Stuur uitgebreide response terug
	return c.JSON(fiber.Map{
		"participant_emails": participantEmails,
		"system_emails":      systemEmails,
		"all_emails":         allEmails,
		"counts": fiber.Map{
			"participants": len(participantEmails),
			"system":       len(systemEmails),
			"total":        len(allEmails),
		},
	})
}

// GetParticipantsByRol is VEROUDERD
// @Summary [VEROUDERD] Participants filteren op rol
// @Tags Participant
// @Router /api/participant/rol/{rol} [get]
func (h *ParticipantHandler) GetParticipantsByRol(c *fiber.Ctx) error {
	rol := c.Params("rol")
	if rol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rol is verplicht",
			"code":  "MISSING_ROLE",
		})
	}

	logger.Error("GetParticipantsByRol handler is aangeroepen, maar de logica is verouderd door V28 refactor.", "rol", rol)
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "Deze functie is verouderd (V28 Refactor). Filteren op rol moet via EventRegistrationHandler.",
		"code":  "DEPRECATED_LOGIC",
	})
}
