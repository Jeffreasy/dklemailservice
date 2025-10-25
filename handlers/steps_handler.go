package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// StepsHandler bevat handlers voor stappen beheer
type StepsHandler struct {
	stepsService      *services.StepsService
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewStepsHandler maakt een nieuwe steps handler
func NewStepsHandler(
	stepsService *services.StepsService,
	authService services.AuthService,
	permissionService services.PermissionService,
) *StepsHandler {
	return &StepsHandler{
		stepsService:      stepsService,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registreert de routes voor stappen beheer
func (h *StepsHandler) RegisterRoutes(app *fiber.App) {
	// Groep voor stappen routes
	stepsGroup := app.Group("/api")

	// POST /api/steps - Update stappen voor ingelogde deelnemer (geen ID nodig!)
	// POST /api/steps/:id - Update stappen voor specifieke deelnemer (admin/staff)
	stepsGroup.Post("/steps", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "write"), h.UpdateSteps)
	stepsGroup.Post("/steps/:id", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "write"), h.UpdateSteps)

	// GET /api/participant/dashboard - Dashboard voor ingelogde deelnemer (geen ID nodig!)
	// GET /api/participant/:id/dashboard - Dashboard voor specifieke deelnemer (admin/staff)
	stepsGroup.Get("/participant/dashboard", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "read"), h.GetParticipantDashboard)
	stepsGroup.Get("/participant/:id/dashboard", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "read"), h.GetParticipantDashboard)

	// GET /api/total-steps - Totaal aantal stappen (alle deelnemers mogen dit zien)
	stepsGroup.Get("/total-steps", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "read_total"), h.GetTotalSteps)

	// GET /api/funds-distribution - Fondsverdeling (admin)
	stepsGroup.Get("/funds-distribution", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "read"), h.GetFundsDistribution)

	// Admin endpoints voor route fund beheer
	adminGroup := stepsGroup.Group("/admin", PermissionMiddleware(h.permissionService, "steps", "write"))
	adminGroup.Get("/route-funds", h.GetRouteFunds)
	adminGroup.Post("/route-funds", h.CreateRouteFund)
	adminGroup.Put("/route-funds/:route", h.UpdateRouteFund)
	adminGroup.Delete("/route-funds/:route", h.DeleteRouteFund)
}

// UpdateSteps werkt stappen bij voor een deelnemer
// @Summary Stappen bijwerken voor deelnemer
// @Description Voegt stappen toe aan een deelnemer (delta)
// @Tags Steps
// @Accept json
// @Produce json
// @Param id path string true "Deelnemer ID"
// @Param request body object{steps=int} true "Stappen delta"
// @Success 200 {object} models.Aanmelding
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/steps/{id} [post]
// @Security BearerAuth
func (h *StepsHandler) UpdateSteps(c *fiber.Ctx) error {
	// Probeer eerst gebruiker ID uit de context (voor ingelogde gebruikers)
	userID, ok := c.Locals("userID").(string)

	// Parse request body
	var req struct {
		Steps int `json:"steps"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Als er een gebruiker is ingelogd, update dan hun stappen
	if ok && userID != "" {
		participant, err := h.stepsService.UpdateStepsByUserID(userID, req.Steps)
		if err != nil {
			logger.Error("Fout bij bijwerken stappen", "error", err, "user_id", userID)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Kon stappen niet bijwerken",
			})
		}
		return c.JSON(participant)
	}

	// Fallback: gebruik ID uit URL parameter (voor admin/staff toegang)
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Update stappen via aanmelding ID
	participant, err := h.stepsService.UpdateSteps(id, req.Steps)
	if err != nil {
		logger.Error("Fout bij bijwerken stappen", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon stappen niet bijwerken",
		})
	}

	return c.JSON(participant)
}

// GetParticipantDashboard haalt dashboard data op voor een deelnemer
// @Summary Dashboard data voor deelnemer
// @Description Haalt stappen en toegewezen fondsen op voor een deelnemer
// @Tags Steps
// @Accept json
// @Produce json
// @Param id path string true "Deelnemer ID"
// @Success 200 {object} object{steps=int,route=string,allocatedFunds=int}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/participant/{id}/dashboard [get]
// @Security BearerAuth
func (h *StepsHandler) GetParticipantDashboard(c *fiber.Ctx) error {
	// Probeer eerst gebruiker ID uit de context (voor ingelogde gebruikers)
	userID, ok := c.Locals("userID").(string)

	// Als er een gebruiker is ingelogd, gebruik dan hun dashboard
	if ok && userID != "" {
		participant, allocatedFunds, err := h.stepsService.GetParticipantDashboardByUserID(userID)
		if err != nil {
			logger.Error("Fout bij ophalen dashboard", "error", err, "id", userID)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Kon dashboard data niet ophalen",
			})
		}

		return c.JSON(fiber.Map{
			"steps":          participant.Steps,
			"route":          participant.Afstand,
			"allocatedFunds": allocatedFunds,
			"naam":           participant.Naam,
			"email":          participant.Email,
		})
	}

	// Fallback: gebruik ID uit URL parameter (voor admin/staff toegang)
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Haal dashboard data op
	participant, allocatedFunds, err := h.stepsService.GetParticipantDashboard(id)
	if err != nil {
		logger.Error("Fout bij ophalen dashboard", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon dashboard data niet ophalen",
		})
	}

	return c.JSON(fiber.Map{
		"steps":          participant.Steps,
		"route":          participant.Afstand,
		"allocatedFunds": allocatedFunds,
		"naam":           participant.Naam,
		"email":          participant.Email,
	})
}

// GetTotalSteps haalt totaal aantal stappen op
// @Summary Totaal aantal stappen
// @Description Haalt totaal aantal stappen op voor een jaar
// @Tags Steps
// @Accept json
// @Produce json
// @Param year query int false "Jaar (standaard huidig jaar)"
// @Success 200 {object} object{total_steps=int}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/total-steps [get]
// @Security BearerAuth
func (h *StepsHandler) GetTotalSteps(c *fiber.Ctx) error {
	// Haal jaar op uit query parameter
	yearStr := c.Query("year", "2025")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldig jaar",
		})
	}

	// Haal totaal stappen op
	totalSteps, err := h.stepsService.GetTotalSteps(year)
	if err != nil {
		logger.Error("Fout bij ophalen totaal stappen", "error", err, "year", year)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon totaal stappen niet ophalen",
		})
	}

	return c.JSON(fiber.Map{
		"total_steps": totalSteps,
	})
}

// GetFundsDistribution haalt fondsverdeling op
// @Summary Fondsverdeling
// @Description Haalt verdeling van fondsen over routes op
// @Tags Steps
// @Accept json
// @Produce json
// @Success 200 {object} object{totalX=int,routes=map[string]int}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/funds-distribution [get]
// @Security BearerAuth
func (h *StepsHandler) GetFundsDistribution(c *fiber.Ctx) error {
	// Haal fondsverdeling op (proportioneel)
	distribution, totalX, err := h.stepsService.GetFundsDistributionProportional()
	if err != nil {
		logger.Error("Fout bij ophalen fondsverdeling", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon fondsverdeling niet ophalen",
		})
	}

	return c.JSON(fiber.Map{
		"totalX": totalX,
		"routes": distribution,
	})
}

// GetRouteFunds haalt alle route fondsallocaties op (admin only)
// @Summary Route fondsallocaties ophalen
// @Description Haalt alle route fondsallocaties op voor beheer
// @Tags Steps Admin
// @Accept json
// @Produce json
// @Success 200 {array} models.RouteFund
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/steps/admin/route-funds [get]
// @Security BearerAuth
func (h *StepsHandler) GetRouteFunds(c *fiber.Ctx) error {
	routeFunds, err := h.stepsService.GetRouteFunds()
	if err != nil {
		logger.Error("Fout bij ophalen route funds", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon route funds niet ophalen",
		})
	}

	return c.JSON(routeFunds)
}

// CreateRouteFund maakt een nieuwe route fondsallocatie aan (admin only)
// @Summary Route fondsallocatie aanmaken
// @Description Maakt een nieuwe route fondsallocatie aan
// @Tags Steps Admin
// @Accept json
// @Produce json
// @Param request body models.RouteFundRequest true "Route fund data"
// @Success 201 {object} models.RouteFund
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/steps/admin/route-funds [post]
// @Security BearerAuth
func (h *StepsHandler) CreateRouteFund(c *fiber.Ctx) error {
	var req models.RouteFundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Validatie
	if req.Route == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Route is verplicht",
		})
	}

	if req.Amount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bedrag moet groter of gelijk zijn aan 0",
		})
	}

	routeFund, err := h.stepsService.CreateRouteFund(req.Route, req.Amount)
	if err != nil {
		logger.Error("Fout bij aanmaken route fund", "error", err, "route", req.Route)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon route fund niet aanmaken",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(routeFund)
}

// UpdateRouteFund werkt een route fondsallocatie bij (admin only)
// @Summary Route fondsallocatie bijwerken
// @Description Werkt een bestaande route fondsallocatie bij
// @Tags Steps Admin
// @Accept json
// @Produce json
// @Param route path string true "Route naam"
// @Param request body models.RouteFundRequest true "Route fund data"
// @Success 200 {object} models.RouteFund
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/steps/admin/route-funds/{route} [put]
// @Security BearerAuth
func (h *StepsHandler) UpdateRouteFund(c *fiber.Ctx) error {
	route := c.Params("route")
	if route == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Route parameter is verplicht",
		})
	}

	var req models.RouteFundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	if req.Amount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bedrag moet groter of gelijk zijn aan 0",
		})
	}

	routeFund, err := h.stepsService.UpdateRouteFund(route, req.Amount)
	if err != nil {
		logger.Error("Fout bij bijwerken route fund", "error", err, "route", route)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon route fund niet bijwerken",
		})
	}

	return c.JSON(routeFund)
}

// DeleteRouteFund verwijdert een route fondsallocatie (admin only)
// @Summary Route fondsallocatie verwijderen
// @Description Verwijdert een route fondsallocatie
// @Tags Steps Admin
// @Accept json
// @Produce json
// @Param route path string true "Route naam"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/steps/admin/route-funds/{route} [delete]
// @Security BearerAuth
func (h *StepsHandler) DeleteRouteFund(c *fiber.Ctx) error {
	route := c.Params("route")
	if route == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Route parameter is verplicht",
		})
	}

	if err := h.stepsService.DeleteRouteFund(route); err != nil {
		logger.Error("Fout bij verwijderen route fund", "error", err, "route", route)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon route fund niet verwijderen",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Route fund succesvol verwijderd",
	})
}
