package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// SocialLinkHandler handles social link-related HTTP requests
type SocialLinkHandler struct {
	socialLinkRepo    repository.SocialLinkRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewSocialLinkHandler creates a new social link handler
func NewSocialLinkHandler(
	socialLinkRepo repository.SocialLinkRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *SocialLinkHandler {
	return &SocialLinkHandler{
		socialLinkRepo:    socialLinkRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the social link routes
func (h *SocialLinkHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/social-links")
	public.Get("/", h.ListVisibleSocialLinks)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/social-links", AuthMiddleware(h.authService))

	// Read routes (require social_link read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "social_link", "read"))
	readGroup.Get("/admin", h.ListSocialLinks)
	readGroup.Get("/:id", h.GetSocialLink)

	// Write routes (require social_link write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "social_link", "write"))
	writeGroup.Post("/", h.CreateSocialLink)
	writeGroup.Put("/:id", h.UpdateSocialLink)

	// Delete routes (require social_link delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "social_link", "delete"))
	deleteGroup.Delete("/:id", h.DeleteSocialLink)
}

// ListVisibleSocialLinks returns all visible social links for public display
// @Summary Get visible social links
// @Description Returns all visible social links ordered by order_number
// @Tags Social Links
// @Accept json
// @Produce json
// @Success 200 {array} models.SocialLink
// @Router /api/social-links [get]
func (h *SocialLinkHandler) ListVisibleSocialLinks(c *fiber.Ctx) error {
	ctx := c.Context()
	links, err := h.socialLinkRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible social links", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social links",
		})
	}

	return c.JSON(links)
}

// ListSocialLinks returns all social links for admin management
// @Summary List all social links
// @Description Returns a paginated list of all social links
// @Tags Social Links
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.SocialLink
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/social-links/admin [get]
// @Security BearerAuth
func (h *SocialLinkHandler) ListSocialLinks(c *fiber.Ctx) error {
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
	links, err := h.socialLinkRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch social links", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social links",
		})
	}

	return c.JSON(links)
}

// GetSocialLink returns a specific social link
// @Summary Get social link by ID
// @Description Returns a specific social link by ID
// @Tags Social Links
// @Accept json
// @Produce json
// @Param id path string true "Social Link ID"
// @Success 200 {object} models.SocialLink
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/social-links/{id} [get]
// @Security BearerAuth
func (h *SocialLinkHandler) GetSocialLink(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Social Link ID is required",
		})
	}

	ctx := c.Context()
	link, err := h.socialLinkRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch social link", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social link",
		})
	}

	if link == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Social link not found",
		})
	}

	return c.JSON(link)
}

// CreateSocialLink creates a new social link
// @Summary Create social link
// @Description Creates a new social link
// @Tags Social Links
// @Accept json
// @Produce json
// @Param link body models.SocialLink true "Social Link data"
// @Success 201 {object} models.SocialLink
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/social-links [post]
// @Security BearerAuth
func (h *SocialLinkHandler) CreateSocialLink(c *fiber.Ctx) error {
	var link models.SocialLink
	if err := c.BodyParser(&link); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if link.Platform == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Platform is required",
		})
	}

	if link.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}

	ctx := c.Context()
	if err := h.socialLinkRepo.Create(ctx, &link); err != nil {
		logger.Error("Failed to create social link", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create social link",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(link)
}

// UpdateSocialLink updates an existing social link
// @Summary Update social link
// @Description Updates an existing social link
// @Tags Social Links
// @Accept json
// @Produce json
// @Param id path string true "Social Link ID"
// @Param link body models.SocialLink true "Social Link data"
// @Success 200 {object} models.SocialLink
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/social-links/{id} [put]
// @Security BearerAuth
func (h *SocialLinkHandler) UpdateSocialLink(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Social Link ID is required",
		})
	}

	// Get existing social link
	ctx := c.Context()
	existing, err := h.socialLinkRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch social link", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social link",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Social link not found",
		})
	}

	// Parse update data
	var updateData models.SocialLink
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.Platform = updateData.Platform
	existing.URL = updateData.URL
	existing.BgColorClass = updateData.BgColorClass
	existing.IconColorClass = updateData.IconColorClass
	existing.OrderNumber = updateData.OrderNumber
	existing.Visible = updateData.Visible

	if err := h.socialLinkRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update social link", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update social link",
		})
	}

	return c.JSON(existing)
}

// DeleteSocialLink deletes a social link
// @Summary Delete social link
// @Description Deletes a social link
// @Tags Social Links
// @Accept json
// @Produce json
// @Param id path string true "Social Link ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/social-links/{id} [delete]
// @Security BearerAuth
func (h *SocialLinkHandler) DeleteSocialLink(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Social Link ID is required",
		})
	}

	// Check if social link exists
	ctx := c.Context()
	link, err := h.socialLinkRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch social link", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social link",
		})
	}

	if link == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Social link not found",
		})
	}

	if err := h.socialLinkRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete social link", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete social link",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Social link deleted successfully",
	})
}
