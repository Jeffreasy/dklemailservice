package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// DashboardHandler handles dashboard-related endpoints
type DashboardHandler struct {
	dashboardService  *services.DashboardService
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(
	dashboardService *services.DashboardService,
	authService services.AuthService,
	permissionService services.PermissionService,
) *DashboardHandler {
	return &DashboardHandler{
		dashboardService:  dashboardService,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the dashboard routes
func (h *DashboardHandler) RegisterRoutes(app *fiber.App) {
	// Admin dashboard routes (require admin permissions)
	adminGroup := app.Group("/api/admin")
	adminGroup.Use(AuthMiddleware(h.authService))
	adminGroup.Use(AdminPermissionMiddleware(h.permissionService))

	adminGroup.Get("/dashboard", h.GetDashboard)
}

// GetDashboard returns aggregated dashboard data
func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get dashboard data from service
	data, err := h.dashboardService.GetDashboardData(ctx)
	if err != nil {
		logger.Error("Fout bij ophalen dashboard data", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon dashboard data niet ophalen",
		})
	}

	return c.JSON(data)
}
