package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// SocialEmbedHandler handles social embed-related HTTP requests
type SocialEmbedHandler struct {
	socialEmbedRepo   repository.SocialEmbedRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewSocialEmbedHandler creates a new social embed handler
func NewSocialEmbedHandler(
	socialEmbedRepo repository.SocialEmbedRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *SocialEmbedHandler {
	return &SocialEmbedHandler{
		socialEmbedRepo:   socialEmbedRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the social embed routes
func (h *SocialEmbedHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/social-embeds")
	public.Get("/", h.ListVisibleSocialEmbeds)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/social-embeds", AuthMiddleware(h.authService))

	// Read routes (require social_embed read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "social_embed", "read"))
	readGroup.Get("/admin", h.ListSocialEmbeds)
	readGroup.Get("/:id", h.GetSocialEmbed)

	// Write routes (require social_embed write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "social_embed", "write"))
	writeGroup.Post("/", h.CreateSocialEmbed)
	writeGroup.Put("/:id", h.UpdateSocialEmbed)

	// Delete routes (require social_embed delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "social_embed", "delete"))
	deleteGroup.Delete("/:id", h.DeleteSocialEmbed)
}

// ListVisibleSocialEmbeds returns all visible social embeds for public display
// @Summary Get visible social embeds
// @Description Returns all visible social embeds ordered by order_number
// @Tags Social Embeds
// @Accept json
// @Produce json
// @Success 200 {array} models.SocialEmbed
// @Router /api/social-embeds [get]
func (h *SocialEmbedHandler) ListVisibleSocialEmbeds(c *fiber.Ctx) error {
	ctx := c.Context()
	embeds, err := h.socialEmbedRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible social embeds", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social embeds",
		})
	}

	return c.JSON(embeds)
}

// ListSocialEmbeds returns all social embeds for admin management
// @Summary List all social embeds
// @Description Returns a paginated list of all social embeds
// @Tags Social Embeds
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.SocialEmbed
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/social-embeds/admin [get]
// @Security BearerAuth
func (h *SocialEmbedHandler) ListSocialEmbeds(c *fiber.Ctx) error {
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
	embeds, err := h.socialEmbedRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch social embeds", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social embeds",
		})
	}

	return c.JSON(embeds)
}

// GetSocialEmbed returns a specific social embed
// @Summary Get social embed by ID
// @Description Returns a specific social embed by ID
// @Tags Social Embeds
// @Accept json
// @Produce json
// @Param id path string true "Social Embed ID"
// @Success 200 {object} models.SocialEmbed
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/social-embeds/{id} [get]
// @Security BearerAuth
func (h *SocialEmbedHandler) GetSocialEmbed(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Social Embed ID is required",
		})
	}

	ctx := c.Context()
	embed, err := h.socialEmbedRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch social embed", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social embed",
		})
	}

	if embed == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Social embed not found",
		})
	}

	return c.JSON(embed)
}

// CreateSocialEmbed creates a new social embed
// @Summary Create social embed
// @Description Creates a new social embed
// @Tags Social Embeds
// @Accept json
// @Produce json
// @Param embed body models.SocialEmbed true "Social Embed data"
// @Success 201 {object} models.SocialEmbed
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/social-embeds [post]
// @Security BearerAuth
func (h *SocialEmbedHandler) CreateSocialEmbed(c *fiber.Ctx) error {
	var embed models.SocialEmbed
	if err := c.BodyParser(&embed); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if embed.Platform == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Platform is required",
		})
	}

	if embed.EmbedCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Embed code is required",
		})
	}

	ctx := c.Context()
	if err := h.socialEmbedRepo.Create(ctx, &embed); err != nil {
		logger.Error("Failed to create social embed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create social embed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(embed)
}

// UpdateSocialEmbed updates an existing social embed
// @Summary Update social embed
// @Description Updates an existing social embed
// @Tags Social Embeds
// @Accept json
// @Produce json
// @Param id path string true "Social Embed ID"
// @Param embed body models.SocialEmbed true "Social Embed data"
// @Success 200 {object} models.SocialEmbed
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/social-embeds/{id} [put]
// @Security BearerAuth
func (h *SocialEmbedHandler) UpdateSocialEmbed(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Social Embed ID is required",
		})
	}

	// Get existing social embed
	ctx := c.Context()
	existing, err := h.socialEmbedRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch social embed", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social embed",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Social embed not found",
		})
	}

	// Parse update data
	var updateData models.SocialEmbed
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.Platform = updateData.Platform
	existing.EmbedCode = updateData.EmbedCode
	existing.OrderNumber = updateData.OrderNumber
	existing.Visible = updateData.Visible

	if err := h.socialEmbedRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update social embed", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update social embed",
		})
	}

	return c.JSON(existing)
}

// DeleteSocialEmbed deletes a social embed
// @Summary Delete social embed
// @Description Deletes a social embed
// @Tags Social Embeds
// @Accept json
// @Produce json
// @Param id path string true "Social Embed ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/social-embeds/{id} [delete]
// @Security BearerAuth
func (h *SocialEmbedHandler) DeleteSocialEmbed(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Social Embed ID is required",
		})
	}

	// Check if social embed exists
	ctx := c.Context()
	embed, err := h.socialEmbedRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch social embed", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch social embed",
		})
	}

	if embed == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Social embed not found",
		})
	}

	if err := h.socialEmbedRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete social embed", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete social embed",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Social embed deleted successfully",
	})
}
