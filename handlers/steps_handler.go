package handlers

import (
	"dklautomationgo/logger"
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

	// POST /api/steps/:id - Update stappen voor deelnemer
	stepsGroup.Post("/steps/:id", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "write"), h.UpdateSteps)

	// GET /api/participant/:id/dashboard - Dashboard voor deelnemer
	stepsGroup.Get("/participant/:id/dashboard", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "read"), h.GetParticipantDashboard)

	// GET /api/total-steps - Totaal aantal stappen (admin)
	stepsGroup.Get("/total-steps", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "read"), h.GetTotalSteps)

	// GET /api/funds-distribution - Fondsverdeling (admin)
	stepsGroup.Get("/funds-distribution", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "steps", "read"), h.GetFundsDistribution)
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
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Parse request body
	var req struct {
		Steps int `json:"steps"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Update stappen
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
	// Haal ID op uit URL
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
