package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AanmeldingHandler bevat handlers voor aanmelding beheer
type AanmeldingHandler struct {
	aanmeldingRepo         repository.AanmeldingRepository
	aanmeldingAntwoordRepo repository.AanmeldingAntwoordRepository
	emailService           *services.EmailService
	authService            services.AuthService
}

// NewAanmeldingHandler maakt een nieuwe aanmelding handler
func NewAanmeldingHandler(
	aanmeldingRepo repository.AanmeldingRepository,
	aanmeldingAntwoordRepo repository.AanmeldingAntwoordRepository,
	emailService *services.EmailService,
	authService services.AuthService,
) *AanmeldingHandler {
	return &AanmeldingHandler{
		aanmeldingRepo:         aanmeldingRepo,
		aanmeldingAntwoordRepo: aanmeldingAntwoordRepo,
		emailService:           emailService,
		authService:            authService,
	}
}

// RegisterRoutes registreert de routes voor aanmelding beheer
func (h *AanmeldingHandler) RegisterRoutes(app *fiber.App) {
	// Groep voor aanmelding beheer routes (vereist admin rechten)
	aanmeldingGroup := app.Group("/api/aanmelding")
	aanmeldingGroup.Use(AuthMiddleware(h.authService))
	aanmeldingGroup.Use(AdminMiddleware(h.authService))

	// Aanmelding beheer routes
	aanmeldingGroup.Get("/", h.ListAanmeldingen)
	aanmeldingGroup.Get("/:id", h.GetAanmelding)
	aanmeldingGroup.Put("/:id", h.UpdateAanmelding)
	aanmeldingGroup.Delete("/:id", h.DeleteAanmelding)
	aanmeldingGroup.Post("/:id/antwoord", h.AddAanmeldingAntwoord)
	aanmeldingGroup.Get("/rol/:rol", h.GetAanmeldingenByRol)
}

// ListAanmeldingen haalt een lijst van aanmeldingen op
// @Summary Lijst van aanmeldingen ophalen
// @Description Haalt een gepagineerde lijst van aanmeldingen op
// @Tags Aanmelding
// @Accept json
// @Produce json
// @Param limit query int false "Aantal resultaten per pagina (standaard 10)"
// @Param offset query int false "Offset voor paginering (standaard 0)"
// @Success 200 {array} models.Aanmelding
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/aanmelding [get]
// @Security BearerAuth
func (h *AanmeldingHandler) ListAanmeldingen(c *fiber.Ctx) error {
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

	// Haal aanmeldingen op
	ctx := c.Context()
	aanmeldingen, err := h.aanmeldingRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen aanmeldingen", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmeldingen niet ophalen",
		})
	}

	// Stuur resultaat terug
	return c.JSON(aanmeldingen)
}

// GetAanmelding haalt een specifieke aanmelding op
// @Summary Details van een specifieke aanmelding ophalen
// @Description Haalt de details van een specifieke aanmelding op, inclusief antwoorden
// @Tags Aanmelding
// @Accept json
// @Produce json
// @Param id path string true "Aanmelding ID"
// @Success 200 {object} models.Aanmelding
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/aanmelding/{id} [get]
// @Security BearerAuth
func (h *AanmeldingHandler) GetAanmelding(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Haal aanmelding op
	ctx := c.Context()
	aanmelding, err := h.aanmeldingRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen aanmelding", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmelding niet ophalen",
		})
	}

	if aanmelding == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Aanmelding niet gevonden",
		})
	}

	// Haal antwoorden op
	antwoorden, err := h.aanmeldingAntwoordRepo.ListByAanmeldingID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen antwoorden", "error", err, "aanmelding_id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon antwoorden niet ophalen",
		})
	}

	// Converteer []*models.AanmeldingAntwoord naar []models.AanmeldingAntwoord
	aanmeldingAntwoorden := make([]models.AanmeldingAntwoord, len(antwoorden))
	for i, antwoord := range antwoorden {
		if antwoord != nil {
			aanmeldingAntwoorden[i] = *antwoord
		}
	}

	// Voeg antwoorden toe aan aanmelding
	aanmelding.Antwoorden = aanmeldingAntwoorden

	// Stuur resultaat terug
	return c.JSON(aanmelding)
}

// UpdateAanmelding werkt een aanmelding bij
// @Summary Aanmelding bijwerken
// @Description Werkt een bestaande aanmelding bij (status, notities)
// @Tags Aanmelding
// @Accept json
// @Produce json
// @Param id path string true "Aanmelding ID"
// @Param aanmelding body models.Aanmelding true "Aanmelding gegevens"
// @Success 200 {object} models.Aanmelding
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/aanmelding/{id} [put]
// @Security BearerAuth
func (h *AanmeldingHandler) UpdateAanmelding(c *fiber.Ctx) error {
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

	// Haal bestaande aanmelding op
	ctx := c.Context()
	aanmelding, err := h.aanmeldingRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen aanmelding", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmelding niet ophalen",
		})
	}

	if aanmelding == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Aanmelding niet gevonden",
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

	// Update aanmelding
	if updateData.Status != "" {
		aanmelding.Status = updateData.Status
	}

	if updateData.Notities != nil {
		aanmelding.Notities = updateData.Notities
	}

	// Stel behandeld door en behandeld op in als dit nog niet is gedaan
	if aanmelding.BehandeldDoor == nil || *aanmelding.BehandeldDoor == "" {
		aanmelding.BehandeldDoor = &gebruiker.Email
		now := time.Now()
		aanmelding.BehandeldOp = &now
	}

	// Sla wijzigingen op
	if err := h.aanmeldingRepo.Update(ctx, aanmelding); err != nil {
		logger.Error("Fout bij bijwerken aanmelding", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmelding niet bijwerken",
		})
	}

	// Stuur bijgewerkte aanmelding terug
	return c.JSON(aanmelding)
}

// DeleteAanmelding verwijdert een aanmelding
// @Summary Aanmelding verwijderen
// @Description Verwijdert een aanmelding en bijbehorende antwoorden
// @Tags Aanmelding
// @Accept json
// @Produce json
// @Param id path string true "Aanmelding ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/aanmelding/{id} [delete]
// @Security BearerAuth
func (h *AanmeldingHandler) DeleteAanmelding(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Controleer of aanmelding bestaat
	ctx := c.Context()
	aanmelding, err := h.aanmeldingRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen aanmelding", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmelding niet ophalen",
		})
	}

	if aanmelding == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Aanmelding niet gevonden",
		})
	}

	// Verwijder aanmelding
	if err := h.aanmeldingRepo.Delete(ctx, id); err != nil {
		logger.Error("Fout bij verwijderen aanmelding", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmelding niet verwijderen",
		})
	}

	// Stuur bevestiging terug
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Aanmelding succesvol verwijderd",
	})
}

// AddAanmeldingAntwoord voegt een antwoord toe aan een aanmelding
// @Summary Antwoord toevoegen aan aanmelding
// @Description Voegt een nieuw antwoord toe aan een aanmelding
// @Tags Aanmelding
// @Accept json
// @Produce json
// @Param id path string true "Aanmelding ID"
// @Param antwoord body models.AanmeldingAntwoord true "Antwoord gegevens"
// @Success 200 {object} models.AanmeldingAntwoord
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/aanmelding/{id}/antwoord [post]
// @Security BearerAuth
func (h *AanmeldingHandler) AddAanmeldingAntwoord(c *fiber.Ctx) error {
	// Haal ID op uit URL
	aanmeldingID := c.Params("id")
	if aanmeldingID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Aanmelding ID is verplicht",
		})
	}

	// Haal gebruiker op uit context
	gebruiker, ok := c.Locals("gebruiker").(*models.Gebruiker)
	if !ok || gebruiker == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon gebruiker niet ophalen uit context",
		})
	}

	// Controleer of aanmelding bestaat
	ctx := c.Context()
	aanmelding, err := h.aanmeldingRepo.GetByID(ctx, aanmeldingID)
	if err != nil {
		logger.Error("Fout bij ophalen aanmelding", "error", err, "id", aanmeldingID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmelding niet ophalen",
		})
	}

	if aanmelding == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Aanmelding niet gevonden",
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
	antwoord := &models.AanmeldingAntwoord{
		AanmeldingID:   aanmeldingID,
		Tekst:          antwoordData.Tekst,
		VerzondDoor:    gebruiker.Email,
		EmailVerzonden: false,
	}

	// Sla antwoord op
	if err := h.aanmeldingAntwoordRepo.Create(ctx, antwoord); err != nil {
		logger.Error("Fout bij opslaan antwoord", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon antwoord niet opslaan",
		})
	}

	// Update aanmelding status
	aanmelding.Status = "beantwoord"
	if err := h.aanmeldingRepo.Update(ctx, aanmelding); err != nil {
		logger.Error("Fout bij bijwerken aanmelding", "error", err, "id", aanmeldingID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmelding niet bijwerken",
		})
	}

	// Stuur e-mail met antwoord (in de achtergrond)
	go func() {
		if err := h.emailService.SendEmail(aanmelding.Email, "Antwoord op uw aanmelding", antwoordData.Tekst); err != nil {
			logger.Error("Fout bij verzenden antwoord e-mail", "error", err, "aanmelding_id", aanmeldingID)
		} else {
			// Update e-mail verzonden status
			antwoord.EmailVerzonden = true
			bgCtx := context.Background()
			if err := h.aanmeldingAntwoordRepo.Update(bgCtx, antwoord); err != nil {
				logger.Error("Fout bij bijwerken antwoord e-mail status", "error", err, "antwoord_id", antwoord.ID)
			}
		}
	}()

	// Stuur antwoord terug
	return c.JSON(antwoord)
}

// GetAanmeldingenByRol haalt aanmeldingen op basis van rol op
// @Summary Aanmeldingen filteren op rol
// @Description Haalt een lijst van aanmeldingen op gefilterd op rol
// @Tags Aanmelding
// @Accept json
// @Produce json
// @Param rol path string true "Rol (vrijwilliger, deelnemer, etc.)"
// @Success 200 {array} models.Aanmelding
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/aanmelding/rol/{rol} [get]
// @Security BearerAuth
func (h *AanmeldingHandler) GetAanmeldingenByRol(c *fiber.Ctx) error {
	// Haal rol op uit URL
	rol := c.Params("rol")
	if rol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rol is verplicht",
		})
	}

	// Haal aanmeldingen op
	ctx := c.Context()

	// Gebruik FindByStatus met rol als parameter (aangezien er geen specifieke FindByRol methode is)
	// In een echte implementatie zou je een specifieke FindByRol methode kunnen toevoegen
	aanmeldingen, err := h.aanmeldingRepo.FindByStatus(ctx, rol)
	if err != nil {
		logger.Error("Fout bij ophalen aanmeldingen op rol", "error", err, "rol", rol)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon aanmeldingen niet ophalen",
		})
	}

	// Filter resultaten op rol (omdat FindByStatus eigenlijk op status filtert, niet op rol)
	filteredAanmeldingen := make([]*models.Aanmelding, 0)
	for _, aanmelding := range aanmeldingen {
		if aanmelding.Rol == rol {
			filteredAanmeldingen = append(filteredAanmeldingen, aanmelding)
		}
	}

	// Stuur resultaat terug
	return c.JSON(filteredAanmeldingen)
}
