package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// PhotoHandler handles photo-related HTTP requests
type PhotoHandler struct {
	photoRepo         repository.PhotoRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewPhotoHandler creates a new photo handler
func NewPhotoHandler(
	photoRepo repository.PhotoRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *PhotoHandler {
	return &PhotoHandler{
		photoRepo:         photoRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the photo routes
func (h *PhotoHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/photos")
	public.Get("/", h.ListVisiblePhotos)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/photos", AuthMiddleware(h.authService))

	// Read routes (require photo read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "photo", "read"))
	readGroup.Get("/admin", h.ListPhotos)
	readGroup.Get("/:id", h.GetPhoto)

	// Write routes (require photo write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "photo", "write"))
	writeGroup.Post("/", h.CreatePhoto)
	writeGroup.Put("/:id", h.UpdatePhoto)

	// Delete routes (require photo delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "photo", "delete"))
	deleteGroup.Delete("/:id", h.DeletePhoto)
}

// ListVisiblePhotos returns all visible photos for public display
// @Summary Get visible photos
// @Description Returns all visible photos
// @Tags Photos
// @Accept json
// @Produce json
// @Success 200 {array} models.Photo
// @Router /api/photos [get]
func (h *PhotoHandler) ListVisiblePhotos(c *fiber.Ctx) error {
	ctx := c.Context()
	photos, err := h.photoRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible photos", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch photos",
		})
	}

	return c.JSON(photos)
}

// ListPhotos returns all photos for admin management
// @Summary List all photos
// @Description Returns a paginated list of all photos
// @Tags Photos
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.Photo
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/photos/admin [get]
// @Security BearerAuth
func (h *PhotoHandler) ListPhotos(c *fiber.Ctx) error {
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
	photos, err := h.photoRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch photos", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch photos",
		})
	}

	return c.JSON(photos)
}

// GetPhoto returns a specific photo
// @Summary Get photo by ID
// @Description Returns a specific photo by ID
// @Tags Photos
// @Accept json
// @Produce json
// @Param id path string true "Photo ID"
// @Success 200 {object} models.Photo
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/photos/{id} [get]
// @Security BearerAuth
func (h *PhotoHandler) GetPhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	ctx := c.Context()
	photo, err := h.photoRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch photo", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch photo",
		})
	}

	if photo == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found",
		})
	}

	return c.JSON(photo)
}

// CreatePhoto creates a new photo
// @Summary Create photo
// @Description Creates a new photo
// @Tags Photos
// @Accept json
// @Produce json
// @Param photo body models.Photo true "Photo data"
// @Success 201 {object} models.Photo
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/photos [post]
// @Security BearerAuth
func (h *PhotoHandler) CreatePhoto(c *fiber.Ctx) error {
	var photo models.Photo
	if err := c.BodyParser(&photo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if photo.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}

	ctx := c.Context()
	if err := h.photoRepo.Create(ctx, &photo); err != nil {
		logger.Error("Failed to create photo", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create photo",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(photo)
}

// UpdatePhoto updates an existing photo
// @Summary Update photo
// @Description Updates an existing photo
// @Tags Photos
// @Accept json
// @Produce json
// @Param id path string true "Photo ID"
// @Param photo body models.Photo true "Photo data"
// @Success 200 {object} models.Photo
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/photos/{id} [put]
// @Security BearerAuth
func (h *PhotoHandler) UpdatePhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Get existing photo
	ctx := c.Context()
	existing, err := h.photoRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch photo", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch photo",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found",
		})
	}

	// Parse update data
	var updateData models.Photo
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.URL = updateData.URL
	existing.AltText = updateData.AltText
	existing.Visible = updateData.Visible
	existing.ThumbnailURL = updateData.ThumbnailURL
	existing.Title = updateData.Title
	existing.Description = updateData.Description
	existing.Year = updateData.Year
	existing.CloudinaryFolder = updateData.CloudinaryFolder

	if err := h.photoRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update photo", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update photo",
		})
	}

	return c.JSON(existing)
}

// DeletePhoto deletes a photo
// @Summary Delete photo
// @Description Deletes a photo
// @Tags Photos
// @Accept json
// @Produce json
// @Param id path string true "Photo ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/photos/{id} [delete]
// @Security BearerAuth
func (h *PhotoHandler) DeletePhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Check if photo exists
	ctx := c.Context()
	photo, err := h.photoRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch photo", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch photo",
		})
	}

	if photo == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found",
		})
	}

	if err := h.photoRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete photo", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete photo",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Photo deleted successfully",
	})
}
