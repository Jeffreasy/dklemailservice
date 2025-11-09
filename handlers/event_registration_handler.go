package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// EventRegistrationHandler bevat handlers voor het beheren van "deelnames" (registraties).
type EventRegistrationHandler struct {
	eventRegRepo      repository.EventRegistrationRepository // De nieuwe repository
	authService       services.AuthService
	permissionService services.PermissionService
	// Voeg hier andere services toe die je nodig hebt (bv. EmailService)
}

// NewEventRegistrationHandler maakt een nieuwe handler voor event registraties.
func NewEventRegistrationHandler(
	eventRegRepo repository.EventRegistrationRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *EventRegistrationHandler {
	return &EventRegistrationHandler{
		eventRegRepo:      eventRegRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registreert de routes voor event registratie beheer.
func (h *EventRegistrationHandler) RegisterRoutes(app *fiber.App) {
	// API-groep voor /api/registration
	regApi := app.Group("/api/registration", AuthMiddleware(h.authService))

	// Route voor het ophalen van één specifieke registratie
	regApi.Get("/:id",
		PermissionMiddleware(h.permissionService, "participant", "read"), // Gebruikt oude permissie
		h.GetRegistration)

	// Route voor het bijwerken van een registratie (status, notities, etc.)
	// Dit was voorheen: PUT /api/aanmelding/:id
	regApi.Put("/:id",
		PermissionMiddleware(h.permissionService, "participant", "write"), // Gebruikt oude permissie
		h.UpdateRegistration)

	// Route voor het filteren op rol
	// Dit was voorheen: GET /api/aanmelding/rol/:rol
	regApi.Get("/rol/:rol",
		PermissionMiddleware(h.permissionService, "participant", "read"), // Gebruikt oude permissie
		h.ListRegistrationsByRole)

	// Alias routes voor backwards compatibility
	registrationsApi := app.Group("/api/registrations", AuthMiddleware(h.authService))

	registrationsApi.Get("/:id",
		PermissionMiddleware(h.permissionService, "participant", "read"),
		h.GetRegistration)

	registrationsApi.Put("/:id",
		PermissionMiddleware(h.permissionService, "participant", "write"),
		h.UpdateRegistration)

	registrationsApi.Get("/rol/:rol",
		PermissionMiddleware(h.permissionService, "participant", "read"),
		h.ListRegistrationsByRole)
}

// GetRegistration haalt één specifieke registratie op.
// @Summary Details van een specifieke registratie ophalen
// @Tags Registration
// @Param id path string true "Event Registration ID"
// @Router /api/registration/{id} [get]
// @Security BearerAuth
func (h *EventRegistrationHandler) GetRegistration(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is verplicht"})
	}

	ctx := c.Context()
	registration, err := h.eventRegRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen event registration", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kon registratie niet ophalen"})
	}
	if registration == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Registratie niet gevonden"})
	}

	return c.JSON(registration)
}

// UpdateRegistration werkt een registratie bij (status, notities, etc.)
// Dit is de *verplaatste logica* uit de oude UpdateAanmelding handler.
// @Summary Registratie bijwerken (status, notities)
// @Tags Registration
// @Param id path string true "Event Registration ID"
// @Param body body object{ status_key=string, notities=string } true "Update data"
// @Router /api/registration/{id} [put]
// @Security BearerAuth
func (h *EventRegistrationHandler) UpdateRegistration(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is verplicht"})
	}

	// Haal gebruiker op uit context
	gebruiker, ok := c.Locals("gebruiker").(*models.Gebruiker)
	if !ok || gebruiker == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kon gebruiker niet ophalen"})
	}

	// Haal bestaande registratie op
	ctx := c.Context()
	registration, err := h.eventRegRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen registratie", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kon registratie niet ophalen"})
	}
	if registration == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Registratie niet gevonden"})
	}

	// Haal update gegevens op uit request body
	var updateData struct {
		Status   string  `json:"status"` // V27: Direct field (database column is 'status')
		Notities *string `json:"notities"`
		// Voeg hier andere update-bare velden toe (bv. Steps, TotalDistance)
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ongeldige gegevens"})
	}

	// Update registratie
	if updateData.Status != "" {
		registration.Status = updateData.Status
	}
	if updateData.Notities != nil {
		registration.Notities = updateData.Notities
	}

	// Stel 'behandeld door' en 'behandeld op' in
	registration.BehandeldDoor = &gebruiker.Email
	now := time.Now()
	registration.BehandeldOp = &now

	// Sla wijzigingen op
	if err := h.eventRegRepo.Update(ctx, registration); err != nil {
		logger.Error("Fout bij bijwerken registratie", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kon registratie niet bijwerken"})
	}

	// Stuur bijgewerkte registratie terug
	return c.JSON(registration)
}

// ListRegistrationsByRole haalt registraties op basis van rol op.
// Dit is de *verplaatste logica* uit de oude GetAanmeldingenByRol handler.
// @Summary Registraties filteren op rol
// @Tags Registration
// @Param rol path string true "Rol (vrijwilliger, deelnemer, etc.)"
// @Router /api/registration/rol/{rol} [get]
// @Security BearerAuth
func (h *EventRegistrationHandler) ListRegistrationsByRole(c *fiber.Ctx) error {
	// Haal rol op uit URL
	rol := c.Params("rol")
	if rol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rol is verplicht",
			"code":  "MISSING_ROLE",
		})
	}

	// Haal registraties op
	ctx := c.Context()
	registrations, err := h.eventRegRepo.ListByRole(ctx, rol)
	if err != nil {
		logger.Error("Fout bij ophalen registraties op rol", "error", err, "rol", rol)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon registraties niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	// Stuur resultaat terug
	return c.JSON(registrations)
}
