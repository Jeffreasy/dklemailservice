package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// AlbumHandler handles album-related HTTP requests
type AlbumHandler struct {
	albumRepo         repository.AlbumRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewAlbumHandler creates a new album handler
func NewAlbumHandler(
	albumRepo repository.AlbumRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *AlbumHandler {
	return &AlbumHandler{
		albumRepo:         albumRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the album routes
func (h *AlbumHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/albums")
	public.Get("/", h.ListVisibleAlbums)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/albums", AuthMiddleware(h.authService))

	// Read routes (require album read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "album", "read"))
	readGroup.Get("/admin", h.ListAlbums)
	readGroup.Get("/:id", h.GetAlbum)

	// Write routes (require album write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "album", "write"))
	writeGroup.Post("/", h.CreateAlbum)
	writeGroup.Put("/:id", h.UpdateAlbum)

	// Delete routes (require album delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "album", "delete"))
	deleteGroup.Delete("/:id", h.DeleteAlbum)
}

// ListVisibleAlbums returns all visible albums for public display
// @Summary Get visible albums
// @Description Returns all visible albums ordered by order_number
// @Tags Albums
// @Accept json
// @Produce json
// @Success 200 {array} models.Album
// @Router /api/albums [get]
func (h *AlbumHandler) ListVisibleAlbums(c *fiber.Ctx) error {
	ctx := c.Context()
	albums, err := h.albumRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible albums", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch albums",
		})
	}

	return c.JSON(albums)
}

// ListAlbums returns all albums for admin management
// @Summary List all albums
// @Description Returns a paginated list of all albums
// @Tags Albums
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.Album
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/albums/admin [get]
// @Security BearerAuth
func (h *AlbumHandler) ListAlbums(c *fiber.Ctx) error {
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
	albums, err := h.albumRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch albums", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch albums",
		})
	}

	return c.JSON(albums)
}

// GetAlbum returns a specific album
// @Summary Get album by ID
// @Description Returns a specific album by ID
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Success 200 {object} models.Album
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/albums/{id} [get]
// @Security BearerAuth
func (h *AlbumHandler) GetAlbum(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Album ID is required",
		})
	}

	ctx := c.Context()
	album, err := h.albumRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch album", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch album",
		})
	}

	if album == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Album not found",
		})
	}

	return c.JSON(album)
}

// CreateAlbum creates a new album
// @Summary Create album
// @Description Creates a new album
// @Tags Albums
// @Accept json
// @Produce json
// @Param album body models.Album true "Album data"
// @Success 201 {object} models.Album
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/albums [post]
// @Security BearerAuth
func (h *AlbumHandler) CreateAlbum(c *fiber.Ctx) error {
	var album models.Album
	if err := c.BodyParser(&album); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if album.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	ctx := c.Context()
	if err := h.albumRepo.Create(ctx, &album); err != nil {
		logger.Error("Failed to create album", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create album",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(album)
}

// UpdateAlbum updates an existing album
// @Summary Update album
// @Description Updates an existing album
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Param album body models.Album true "Album data"
// @Success 200 {object} models.Album
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/albums/{id} [put]
// @Security BearerAuth
func (h *AlbumHandler) UpdateAlbum(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Album ID is required",
		})
	}

	// Get existing album
	ctx := c.Context()
	existing, err := h.albumRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch album", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch album",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Album not found",
		})
	}

	// Parse update data
	var updateData models.Album
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.Title = updateData.Title
	existing.Description = updateData.Description
	existing.CoverPhotoID = updateData.CoverPhotoID
	existing.Visible = updateData.Visible
	existing.OrderNumber = updateData.OrderNumber

	if err := h.albumRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update album", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update album",
		})
	}

	return c.JSON(existing)
}

// DeleteAlbum deletes an album
// @Summary Delete album
// @Description Deletes an album
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/albums/{id} [delete]
// @Security BearerAuth
func (h *AlbumHandler) DeleteAlbum(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Album ID is required",
		})
	}

	// Check if album exists
	ctx := c.Context()
	album, err := h.albumRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch album", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch album",
		})
	}

	if album == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Album not found",
		})
	}

	if err := h.albumRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete album", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete album",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Album deleted successfully",
	})
}
