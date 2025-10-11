package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// PartnerHandler handles partner-related HTTP requests
type PartnerHandler struct {
	partnerRepo       repository.PartnerRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewPartnerHandler creates a new partner handler
func NewPartnerHandler(
	partnerRepo repository.PartnerRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *PartnerHandler {
	return &PartnerHandler{
		partnerRepo:       partnerRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the partner routes
func (h *PartnerHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/partners")
	public.Get("/", h.ListVisiblePartners)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/partners", AuthMiddleware(h.authService))

	// Read routes (require partner read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "partner", "read"))
	readGroup.Get("/admin", h.ListPartners)
	readGroup.Get("/:id", h.GetPartner)

	// Write routes (require partner write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "partner", "write"))
	writeGroup.Post("/", h.CreatePartner)
	writeGroup.Put("/:id", h.UpdatePartner)

	// Delete routes (require partner delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "partner", "delete"))
	deleteGroup.Delete("/:id", h.DeletePartner)
}

// ListVisiblePartners returns all visible partners for public display
// @Summary Get visible partners
// @Description Returns all visible partners ordered by order_number
// @Tags Partners
// @Accept json
// @Produce json
// @Success 200 {array} models.Partner
// @Router /api/partners [get]
func (h *PartnerHandler) ListVisiblePartners(c *fiber.Ctx) error {
	ctx := c.Context()
	partners, err := h.partnerRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible partners", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch partners",
		})
	}

	return c.JSON(partners)
}

// ListPartners returns all partners for admin management
// @Summary List all partners
// @Description Returns a paginated list of all partners
// @Tags Partners
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.Partner
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/partners/admin [get]
// @Security BearerAuth
func (h *PartnerHandler) ListPartners(c *fiber.Ctx) error {
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
	partners, err := h.partnerRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch partners", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch partners",
		})
	}

	return c.JSON(partners)
}

// GetPartner returns a specific partner
// @Summary Get partner by ID
// @Description Returns a specific partner by ID
// @Tags Partners
// @Accept json
// @Produce json
// @Param id path string true "Partner ID"
// @Success 200 {object} models.Partner
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/partners/{id} [get]
// @Security BearerAuth
func (h *PartnerHandler) GetPartner(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Partner ID is required",
		})
	}

	ctx := c.Context()
	partner, err := h.partnerRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch partner", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch partner",
		})
	}

	if partner == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Partner not found",
		})
	}

	return c.JSON(partner)
}

// CreatePartner creates a new partner
// @Summary Create partner
// @Description Creates a new partner
// @Tags Partners
// @Accept json
// @Produce json
// @Param partner body models.Partner true "Partner data"
// @Success 201 {object} models.Partner
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/partners [post]
// @Security BearerAuth
func (h *PartnerHandler) CreatePartner(c *fiber.Ctx) error {
	var partner models.Partner
	if err := c.BodyParser(&partner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if partner.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	ctx := c.Context()
	if err := h.partnerRepo.Create(ctx, &partner); err != nil {
		logger.Error("Failed to create partner", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create partner",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(partner)
}

// UpdatePartner updates an existing partner
// @Summary Update partner
// @Description Updates an existing partner
// @Tags Partners
// @Accept json
// @Produce json
// @Param id path string true "Partner ID"
// @Param partner body models.Partner true "Partner data"
// @Success 200 {object} models.Partner
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/partners/{id} [put]
// @Security BearerAuth
func (h *PartnerHandler) UpdatePartner(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Partner ID is required",
		})
	}

	// Get existing partner
	ctx := c.Context()
	existing, err := h.partnerRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch partner", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch partner",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Partner not found",
		})
	}

	// Parse update data
	var updateData models.Partner
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.Name = updateData.Name
	existing.Description = updateData.Description
	existing.Logo = updateData.Logo
	existing.Website = updateData.Website
	existing.Tier = updateData.Tier
	existing.Since = updateData.Since
	existing.Visible = updateData.Visible
	existing.OrderNumber = updateData.OrderNumber

	if err := h.partnerRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update partner", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update partner",
		})
	}

	return c.JSON(existing)
}

// DeletePartner deletes a partner
// @Summary Delete partner
// @Description Deletes a partner
// @Tags Partners
// @Accept json
// @Produce json
// @Param id path string true "Partner ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/partners/{id} [delete]
// @Security BearerAuth
func (h *PartnerHandler) DeletePartner(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Partner ID is required",
		})
	}

	// Check if partner exists
	ctx := c.Context()
	partner, err := h.partnerRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch partner", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch partner",
		})
	}

	if partner == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Partner not found",
		})
	}

	if err := h.partnerRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete partner", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete partner",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Partner deleted successfully",
	})
}
