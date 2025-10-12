package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UnderConstructionHandler handles under construction-related HTTP requests
type UnderConstructionHandler struct {
	underConstructionRepo repository.UnderConstructionRepository
	authService           services.AuthService
	permissionService     services.PermissionService
}

// NewUnderConstructionHandler creates a new under construction handler
func NewUnderConstructionHandler(
	underConstructionRepo repository.UnderConstructionRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *UnderConstructionHandler {
	return &UnderConstructionHandler{
		underConstructionRepo: underConstructionRepo,
		authService:           authService,
		permissionService:     permissionService,
	}
}

// RegisterRoutes registers the under construction routes
func (h *UnderConstructionHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/under-construction")
	public.Get("/active", h.GetActiveUnderConstruction)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/under-construction", AuthMiddleware(h.authService))

	// Read routes (require under_construction read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "under_construction", "read"))
	readGroup.Get("/admin", h.ListUnderConstruction)
	readGroup.Get("/:id", h.GetUnderConstruction)

	// Write routes (require under_construction write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "under_construction", "write"))
	writeGroup.Post("/", h.CreateUnderConstruction)
	writeGroup.Put("/:id", h.UpdateUnderConstruction)

	// Delete routes (require under_construction delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "under_construction", "delete"))
	deleteGroup.Delete("/:id", h.DeleteUnderConstruction)
}

// GetActiveUnderConstruction returns the active under construction record for public display
// @Summary Get active under construction
// @Description Returns the active under construction record if any
// @Tags Under Construction
// @Accept json
// @Produce json
// @Success 200 {object} models.UnderConstruction
// @Failure 404 {object} map[string]interface{}
// @Router /api/under-construction/active [get]
func (h *UnderConstructionHandler) GetActiveUnderConstruction(c *fiber.Ctx) error {
	ctx := c.Context()
	uc, err := h.underConstructionRepo.GetActive(ctx)
	if err != nil {
		logger.Error("Failed to fetch active under construction", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch under construction",
		})
	}

	if uc == nil {
		// Return 404 without logging an error - this is expected when maintenance mode is disabled
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No active under construction found",
		})
	}

	return c.JSON(uc)
}

// ListUnderConstruction returns all under construction records for admin management
// @Summary List all under construction records
// @Description Returns a paginated list of all under construction records
// @Tags Under Construction
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.UnderConstruction
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/under-construction/admin [get]
// @Security BearerAuth
func (h *UnderConstructionHandler) ListUnderConstruction(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	if limit < 1 || limit > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Limit must be between 1 and 100",
		})
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Offset must be non-negative",
		})
	}

	ctx := c.Context()
	ucs, err := h.underConstructionRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch under construction records", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch under construction records",
		})
	}

	return c.JSON(ucs)
}

// GetUnderConstruction returns a specific under construction record
// @Summary Get under construction by ID
// @Description Returns a specific under construction record by ID
// @Tags Under Construction
// @Accept json
// @Produce json
// @Param id path int true "Under Construction ID"
// @Success 200 {object} models.UnderConstruction
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/under-construction/{id} [get]
// @Security BearerAuth
func (h *UnderConstructionHandler) GetUnderConstruction(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Under Construction ID is required",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Under Construction ID",
		})
	}

	ctx := c.Context()
	uc, err := h.underConstructionRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch under construction", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch under construction",
		})
	}

	if uc == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Under construction not found",
		})
	}

	return c.JSON(uc)
}

// CreateUnderConstruction creates a new under construction record
// @Summary Create under construction
// @Description Creates a new under construction record
// @Tags Under Construction
// @Accept json
// @Produce json
// @Param uc body models.UnderConstruction true "Under Construction data"
// @Success 201 {object} models.UnderConstruction
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/under-construction [post]
// @Security BearerAuth
func (h *UnderConstructionHandler) CreateUnderConstruction(c *fiber.Ctx) error {
	var uc models.UnderConstruction
	if err := c.BodyParser(&uc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if uc.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	if uc.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message is required",
		})
	}

	ctx := c.Context()
	if err := h.underConstructionRepo.Create(ctx, &uc); err != nil {
		logger.Error("Failed to create under construction", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create under construction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(uc)
}

// UpdateUnderConstruction updates an existing under construction record
// @Summary Update under construction
// @Description Updates an existing under construction record
// @Tags Under Construction
// @Accept json
// @Produce json
// @Param id path int true "Under Construction ID"
// @Param uc body models.UnderConstruction true "Under Construction data"
// @Success 200 {object} models.UnderConstruction
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/under-construction/{id} [put]
// @Security BearerAuth
func (h *UnderConstructionHandler) UpdateUnderConstruction(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Under Construction ID is required",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Under Construction ID",
		})
	}

	// Get existing under construction
	ctx := c.Context()
	existing, err := h.underConstructionRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch under construction", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch under construction",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Under construction not found",
		})
	}

	// Parse update data
	var updateData models.UnderConstruction
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.IsActive = updateData.IsActive
	existing.Title = updateData.Title
	existing.Message = updateData.Message
	existing.FooterText = updateData.FooterText
	existing.LogoURL = updateData.LogoURL
	existing.ExpectedDate = updateData.ExpectedDate
	existing.SocialLinks = updateData.SocialLinks
	existing.ProgressPercentage = updateData.ProgressPercentage
	existing.ContactEmail = updateData.ContactEmail
	existing.NewsletterEnabled = updateData.NewsletterEnabled

	if err := h.underConstructionRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update under construction", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update under construction",
		})
	}

	return c.JSON(existing)
}

// DeleteUnderConstruction deletes an under construction record
// @Summary Delete under construction
// @Description Deletes an under construction record
// @Tags Under Construction
// @Accept json
// @Produce json
// @Param id path int true "Under Construction ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/under-construction/{id} [delete]
// @Security BearerAuth
func (h *UnderConstructionHandler) DeleteUnderConstruction(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Under Construction ID is required",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Under Construction ID",
		})
	}

	// Check if under construction exists
	ctx := c.Context()
	uc, err := h.underConstructionRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch under construction", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch under construction",
		})
	}

	if uc == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Under construction not found",
		})
	}

	if err := h.underConstructionRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete under construction", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete under construction",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Under construction deleted successfully",
	})
}
