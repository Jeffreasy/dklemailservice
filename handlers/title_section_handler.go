package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// TitleSectionHandler handles title section-related HTTP requests
type TitleSectionHandler struct {
	titleSectionRepo  repository.TitleSectionRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewTitleSectionHandler creates a new title section handler
func NewTitleSectionHandler(
	titleSectionRepo repository.TitleSectionRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *TitleSectionHandler {
	return &TitleSectionHandler{
		titleSectionRepo:  titleSectionRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the title section routes
func (h *TitleSectionHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/title-sections")
	public.Get("/", h.GetTitleSection)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/title-sections", AuthMiddleware(h.authService))

	// Read routes (require title_section read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "title_section", "read"))
	readGroup.Get("/admin", h.GetTitleSectionAdmin)

	// Write routes (require title_section write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "title_section", "write"))
	writeGroup.Post("/", h.CreateTitleSection)
	writeGroup.Put("/", h.UpdateTitleSection)

	// Delete routes (require title_section delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "title_section", "delete"))
	deleteGroup.Delete("/:id", h.DeleteTitleSection)
}

// GetTitleSection returns the title section content for public display
// @Summary Get title section
// @Description Returns the title section content
// @Tags Title Sections
// @Accept json
// @Produce json
// @Success 200 {object} models.TitleSection
// @Router /api/title-sections [get]
func (h *TitleSectionHandler) GetTitleSection(c *fiber.Ctx) error {
	ctx := c.Context()
	titleSection, err := h.titleSectionRepo.Get(ctx)
	if err != nil {
		logger.Error("Failed to fetch title section", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch title section",
		})
	}

	return c.JSON(titleSection)
}

// GetTitleSectionAdmin returns the title section content for admin management
// @Summary Get title section for admin
// @Description Returns the title section content for admin management
// @Tags Title Sections
// @Accept json
// @Produce json
// @Success 200 {object} models.TitleSection
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/title-sections/admin [get]
// @Security BearerAuth
func (h *TitleSectionHandler) GetTitleSectionAdmin(c *fiber.Ctx) error {
	ctx := c.Context()
	titleSection, err := h.titleSectionRepo.Get(ctx)
	if err != nil {
		logger.Error("Failed to fetch title section", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch title section",
		})
	}

	return c.JSON(titleSection)
}

// CreateTitleSection creates a new title section
// @Summary Create title section
// @Description Creates a new title section
// @Tags Title Sections
// @Accept json
// @Produce json
// @Param titleSection body models.TitleSection true "Title section data"
// @Success 201 {object} models.TitleSection
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/title-sections [post]
// @Security BearerAuth
func (h *TitleSectionHandler) CreateTitleSection(c *fiber.Ctx) error {
	var titleSection models.TitleSection
	if err := c.BodyParser(&titleSection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if titleSection.EventTitle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event title is required",
		})
	}

	ctx := c.Context()
	if err := h.titleSectionRepo.Create(ctx, &titleSection); err != nil {
		logger.Error("Failed to create title section", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create title section",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(titleSection)
}

// UpdateTitleSection updates the existing title section
// @Summary Update title section
// @Description Updates the existing title section
// @Tags Title Sections
// @Accept json
// @Produce json
// @Param titleSection body models.TitleSection true "Title section data"
// @Success 200 {object} models.TitleSection
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/title-sections [put]
// @Security BearerAuth
func (h *TitleSectionHandler) UpdateTitleSection(c *fiber.Ctx) error {
	// Get existing title section
	ctx := c.Context()
	existing, err := h.titleSectionRepo.Get(ctx)
	if err != nil {
		logger.Error("Failed to fetch title section", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch title section",
		})
	}

	// Parse update data
	var updateData models.TitleSection
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.EventTitle = updateData.EventTitle
	existing.EventSubtitle = updateData.EventSubtitle
	existing.ImageURL = updateData.ImageURL
	existing.ImageAlt = updateData.ImageAlt
	existing.Detail1Title = updateData.Detail1Title
	existing.Detail1Description = updateData.Detail1Description
	existing.Detail2Title = updateData.Detail2Title
	existing.Detail2Description = updateData.Detail2Description
	existing.Detail3Title = updateData.Detail3Title
	existing.Detail3Description = updateData.Detail3Description
	existing.ParticipantCount = updateData.ParticipantCount

	if err := h.titleSectionRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update title section", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update title section",
		})
	}

	return c.JSON(existing)
}

// DeleteTitleSection deletes a title section
// @Summary Delete title section
// @Description Deletes a title section by ID
// @Tags Title Sections
// @Accept json
// @Produce json
// @Param id path string true "Title section ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/title-sections/{id} [delete]
// @Security BearerAuth
func (h *TitleSectionHandler) DeleteTitleSection(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title section ID is required",
		})
	}

	// Check if title section exists
	ctx := c.Context()
	titleSection, err := h.titleSectionRepo.Get(ctx)
	if err != nil {
		logger.Error("Failed to fetch title section", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch title section",
		})
	}

	if titleSection == nil || titleSection.ID != id {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Title section not found",
		})
	}

	if err := h.titleSectionRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete title section", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete title section",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Title section deleted successfully",
	})
}
