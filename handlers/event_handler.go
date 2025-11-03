package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// EventHandler bevat handlers voor event beheer
type EventHandler struct {
	eventRepo         repository.EventRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewEventHandler maakt een nieuwe event handler
func NewEventHandler(
	eventRepo repository.EventRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *EventHandler {
	return &EventHandler{
		eventRepo:         eventRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registreert de routes voor event beheer
func (h *EventHandler) RegisterRoutes(app *fiber.App) {
	// Public routes - iedereen kan events zien
	publicEvents := app.Group("/api/events")
	publicEvents.Get("/", h.ListEvents)
	publicEvents.Get("/active", h.GetActiveEvent)
	publicEvents.Get("/:id", h.GetEvent)

	// Protected routes - vereisen authenticatie en permissions
	adminEvents := app.Group("/api/events",
		AuthMiddleware(h.authService),
		PermissionMiddleware(h.permissionService, "events", "write"),
	)
	adminEvents.Post("/", h.CreateEvent)
	adminEvents.Put("/:id", h.UpdateEvent)
	adminEvents.Delete("/:id", h.DeleteEvent)

	// Event participant routes
	participants := app.Group("/api/events/:id/participants",
		AuthMiddleware(h.authService),
		PermissionMiddleware(h.permissionService, "events", "read"),
	)
	participants.Get("/", h.GetEventParticipants)
}

// GetEvent haalt een event op basis van ID
// @Summary Event Details ophalen
// @Description Haalt gedetailleerde informatie op voor een specifiek event
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} models.EventResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/{id} [get]
func (h *EventHandler) GetEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event ID is verplicht",
			"code":  "MISSING_ID",
		})
	}

	event, err := h.eventRepo.GetByID(c.Context(), id)
	if err != nil {
		logger.Error("Fout bij ophalen event", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon event niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	if event == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event niet gevonden",
			"code":  "NOT_FOUND",
		})
	}

	return c.JSON(event.ToResponse())
}

// GetActiveEvent haalt het actieve event op
// @Summary Actief Event ophalen
// @Description Haalt het eerstvolgende actieve event op
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {object} models.EventResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/active [get]
func (h *EventHandler) GetActiveEvent(c *fiber.Ctx) error {
	event, err := h.eventRepo.GetActiveEvent(c.Context())
	if err != nil {
		logger.Error("Fout bij ophalen actief event", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon actief event niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	if event == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Geen actief event gevonden",
			"code":  "NOT_FOUND",
		})
	}

	return c.JSON(event.ToResponse())
}

// ListEvents haalt lijst van events op
// @Summary Events lijst ophalen
// @Description Haalt lijst van alle events op
// @Tags Events
// @Accept json
// @Produce json
// @Param limit query int false "Aantal resultaten" default(50)
// @Param offset query int false "Offset voor paginatie" default(0)
// @Param active_only query bool false "Alleen actieve events" default(false)
// @Success 200 {array} models.EventResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/events [get]
func (h *EventHandler) ListEvents(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)
	activeOnly := c.QueryBool("active_only", false)

	var events []*models.Event
	var err error

	if activeOnly {
		events, err = h.eventRepo.ListActive(c.Context())
	} else {
		events, err = h.eventRepo.List(c.Context(), limit, offset)
	}

	if err != nil {
		logger.Error("Fout bij ophalen events", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon events niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	// Convert to response format
	responses := make([]*models.EventResponse, len(events))
	for i, event := range events {
		responses[i] = event.ToResponse()
	}

	return c.JSON(responses)
}

// CreateEvent maakt een nieuw event aan
// @Summary Event aanmaken
// @Description Maakt een nieuw event aan (admin only)
// @Tags Events
// @Accept json
// @Produce json
// @Param event body models.EventCreateRequest true "Event data"
// @Success 201 {object} models.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events [post]
// @Security BearerAuth
func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var req models.EventCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
			"code":  "INVALID_REQUEST",
		})
	}

	// Validatie
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event naam is verplicht",
			"code":  "MISSING_NAME",
		})
	}

	if req.StartTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Start tijd is verplicht",
			"code":  "MISSING_START_TIME",
		})
	}

	if len(req.Geofences) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Minimaal één geofence is verplicht",
			"code":  "MISSING_GEOFENCES",
		})
	}

	// Parse start time
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige start tijd format (gebruik ISO 8601)",
			"code":  "INVALID_TIME_FORMAT",
		})
	}

	// Parse end time if provided
	var endTime *time.Time
	if req.EndTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ongeldige eind tijd format (gebruik ISO 8601)",
				"code":  "INVALID_TIME_FORMAT",
			})
		}
		endTime = &parsed
	}

	// Get user ID van context
	userID, _ := c.Locals("userID").(string)

	// Create event
	event := &models.Event{
		Name:        req.Name,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      req.Status,
		Geofences:   req.Geofences,
		EventConfig: req.EventConfig,
		IsActive:    req.IsActive,
	}

	if userID != "" {
		event.CreatedBy = &userID
	}

	// Set default status if not provided
	if event.Status == "" {
		event.Status = models.EventStatusUpcoming
	}

	if err := h.eventRepo.Create(c.Context(), event); err != nil {
		logger.Error("Fout bij aanmaken event", "error", err, "name", req.Name)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon event niet aanmaken",
			"code":  "INTERNAL_ERROR",
		})
	}

	logger.Info("Event aangemaakt", "id", event.ID, "name", event.Name, "created_by", userID)

	return c.Status(fiber.StatusCreated).JSON(event.ToResponse())
}

// UpdateEvent werkt een bestaand event bij
// @Summary Event bijwerken
// @Description Werkt een bestaand event bij (admin only)
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param event body models.EventUpdateRequest true "Event data"
// @Success 200 {object} models.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/{id} [put]
// @Security BearerAuth
func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event ID is verplicht",
			"code":  "MISSING_ID",
		})
	}

	var req models.EventUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
			"code":  "INVALID_REQUEST",
		})
	}

	// Haal bestaand event op
	event, err := h.eventRepo.GetByID(c.Context(), id)
	if err != nil {
		logger.Error("Fout bij ophalen event", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon event niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	if event == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event niet gevonden",
			"code":  "NOT_FOUND",
		})
	}

	// Update fields
	if req.Name != "" {
		event.Name = req.Name
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Status != "" {
		event.Status = req.Status
	}
	if req.StartTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ongeldige start tijd format",
				"code":  "INVALID_TIME_FORMAT",
			})
		}
		event.StartTime = parsed
	}
	if req.EndTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ongeldige eind tijd format",
				"code":  "INVALID_TIME_FORMAT",
			})
		}
		event.EndTime = &parsed
	}
	if req.Geofences != nil {
		event.Geofences = req.Geofences
	}
	if req.EventConfig != nil {
		event.EventConfig = req.EventConfig
	}
	if req.IsActive != nil {
		event.IsActive = *req.IsActive
	}

	if err := h.eventRepo.Update(c.Context(), event); err != nil {
		logger.Error("Fout bij bijwerken event", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon event niet bijwerken",
			"code":  "INTERNAL_ERROR",
		})
	}

	logger.Info("Event bijgewerkt", "id", event.ID, "name", event.Name)

	return c.JSON(event.ToResponse())
}

// DeleteEvent verwijdert een event
// @Summary Event verwijderen
// @Description Verwijdert een event (admin only)
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/{id} [delete]
// @Security BearerAuth
func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event ID is verplicht",
			"code":  "MISSING_ID",
		})
	}

	// Check if event exists
	event, err := h.eventRepo.GetByID(c.Context(), id)
	if err != nil {
		logger.Error("Fout bij ophalen event", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon event niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	if event == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event niet gevonden",
			"code":  "NOT_FOUND",
		})
	}

	if err := h.eventRepo.Delete(c.Context(), id); err != nil {
		logger.Error("Fout bij verwijderen event", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon event niet verwijderen",
			"code":  "INTERNAL_ERROR",
		})
	}

	logger.Info("Event verwijderd", "id", id, "name", event.Name)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Event succesvol verwijderd",
	})
}

// GetEventParticipants haalt alle participants voor een event op
// @Summary Event Participants ophalen
// @Description Haalt lijst van participants voor een event op
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {array} models.EventParticipantResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/{id}/participants [get]
// @Security BearerAuth
func (h *EventHandler) GetEventParticipants(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event ID is verplicht",
			"code":  "MISSING_ID",
		})
	}

	participants, err := h.eventRepo.GetEventParticipants(c.Context(), id)
	if err != nil {
		logger.Error("Fout bij ophalen event participants", "error", err, "event_id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon participants niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	// Convert to response format
	responses := make([]*models.EventParticipantResponse, len(participants))
	for i, p := range participants {
		responses[i] = p.ToResponse()
	}

	return c.JSON(responses)
}
