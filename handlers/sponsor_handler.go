package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// SponsorHandler handles sponsor-related HTTP requests
type SponsorHandler struct {
	sponsorRepo       repository.SponsorRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewSponsorHandler creates a new sponsor handler
func NewSponsorHandler(
	sponsorRepo repository.SponsorRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *SponsorHandler {
	return &SponsorHandler{
		sponsorRepo:       sponsorRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the sponsor routes
func (h *SponsorHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/sponsors")
	public.Get("/", h.ListVisibleSponsors)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/sponsors", AuthMiddleware(h.authService))

	// Read routes (require sponsor read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "sponsor", "read"))
	readGroup.Get("/admin", h.ListSponsors)
	readGroup.Get("/:id", h.GetSponsor)

	// Write routes (require sponsor write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "sponsor", "write"))
	writeGroup.Post("/", h.CreateSponsor)
	writeGroup.Put("/:id", h.UpdateSponsor)

	// Delete routes (require sponsor delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "sponsor", "delete"))
	deleteGroup.Delete("/:id", h.DeleteSponsor)
}

// ListVisibleSponsors returns all visible sponsors for public display
// @Summary Get visible sponsors
// @Description Returns all visible sponsors ordered by order_number
// @Tags Sponsors
// @Accept json
// @Produce json
// @Success 200 {array} models.Sponsor
// @Router /api/sponsors [get]
func (h *SponsorHandler) ListVisibleSponsors(c *fiber.Ctx) error {
	ctx := c.Context()
	sponsors, err := h.sponsorRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible sponsors", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sponsors",
		})
	}

	return c.JSON(sponsors)
}

// ListSponsors returns all sponsors for admin management
// @Summary List all sponsors
// @Description Returns a paginated list of all sponsors
// @Tags Sponsors
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.Sponsor
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/sponsors/admin [get]
// @Security BearerAuth
func (h *SponsorHandler) ListSponsors(c *fiber.Ctx) error {
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
	sponsors, err := h.sponsorRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch sponsors", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sponsors",
		})
	}

	return c.JSON(sponsors)
}

// GetSponsor returns a specific sponsor
// @Summary Get sponsor by ID
// @Description Returns a specific sponsor by ID
// @Tags Sponsors
// @Accept json
// @Produce json
// @Param id path string true "Sponsor ID"
// @Success 200 {object} models.Sponsor
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/sponsors/{id} [get]
// @Security BearerAuth
func (h *SponsorHandler) GetSponsor(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Sponsor ID is required",
		})
	}

	ctx := c.Context()
	sponsor, err := h.sponsorRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch sponsor", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sponsor",
		})
	}

	if sponsor == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sponsor not found",
		})
	}

	return c.JSON(sponsor)
}

// CreateSponsor creates a new sponsor
// @Summary Create sponsor
// @Description Creates a new sponsor
// @Tags Sponsors
// @Accept json
// @Produce json
// @Param sponsor body models.Sponsor true "Sponsor data"
// @Success 201 {object} models.Sponsor
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/sponsors [post]
// @Security BearerAuth
func (h *SponsorHandler) CreateSponsor(c *fiber.Ctx) error {
	var sponsor models.Sponsor
	if err := c.BodyParser(&sponsor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if sponsor.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	ctx := c.Context()
	if err := h.sponsorRepo.Create(ctx, &sponsor); err != nil {
		logger.Error("Failed to create sponsor", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create sponsor",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(sponsor)
}

// UpdateSponsor updates an existing sponsor
// @Summary Update sponsor
// @Description Updates an existing sponsor
// @Tags Sponsors
// @Accept json
// @Produce json
// @Param id path string true "Sponsor ID"
// @Param sponsor body models.Sponsor true "Sponsor data"
// @Success 200 {object} models.Sponsor
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/sponsors/{id} [put]
// @Security BearerAuth
func (h *SponsorHandler) UpdateSponsor(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Sponsor ID is required",
		})
	}

	// Get existing sponsor
	ctx := c.Context()
	existing, err := h.sponsorRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch sponsor", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sponsor",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sponsor not found",
		})
	}

	// Parse update data
	var updateData models.Sponsor
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.Name = updateData.Name
	existing.Description = updateData.Description
	existing.LogoURL = updateData.LogoURL
	existing.WebsiteURL = updateData.WebsiteURL
	existing.OrderNumber = updateData.OrderNumber
	existing.IsActive = updateData.IsActive
	existing.Visible = updateData.Visible

	if err := h.sponsorRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update sponsor", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update sponsor",
		})
	}

	return c.JSON(existing)
}

// DeleteSponsor deletes a sponsor
// @Summary Delete sponsor
// @Description Deletes a sponsor
// @Tags Sponsors
// @Accept json
// @Produce json
// @Param id path string true "Sponsor ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/sponsors/{id} [delete]
// @Security BearerAuth
func (h *SponsorHandler) DeleteSponsor(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Sponsor ID is required",
		})
	}

	// Check if sponsor exists
	ctx := c.Context()
	sponsor, err := h.sponsorRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch sponsor", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sponsor",
		})
	}

	if sponsor == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sponsor not found",
		})
	}

	if err := h.sponsorRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete sponsor", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete sponsor",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Sponsor deleted successfully",
	})
}
