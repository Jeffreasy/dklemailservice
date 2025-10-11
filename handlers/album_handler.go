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
	photoRepo         repository.PhotoRepository
	albumPhotoRepo    repository.AlbumPhotoRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewAlbumHandler creates a new album handler
func NewAlbumHandler(
	albumRepo repository.AlbumRepository,
	photoRepo repository.PhotoRepository,
	albumPhotoRepo repository.AlbumPhotoRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *AlbumHandler {
	return &AlbumHandler{
		albumRepo:         albumRepo,
		photoRepo:         photoRepo,
		albumPhotoRepo:    albumPhotoRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the album routes
func (h *AlbumHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/albums")
	public.Get("/", h.ListVisibleAlbums)
	public.Get("/:id/photos", h.GetAlbumPhotos)

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
	writeGroup.Put("/reorder", h.ReorderAlbums)
	writeGroup.Post("/:id/photos", h.AddPhotoToAlbum)
	writeGroup.Put("/:id/photos/reorder", h.ReorderAlbumPhotos)

	// Delete routes (require album delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "album", "delete"))
	deleteGroup.Delete("/:id", h.DeleteAlbum)
	deleteGroup.Delete("/:id/photos/:photoId", h.RemovePhotoFromAlbum)
}

// ListVisibleAlbums returns all visible albums for public display
// @Summary Get visible albums
// @Description Returns all visible albums ordered by order_number. Use ?include_covers=true to include cover photo information.
// @Tags Albums
// @Accept json
// @Produce json
// @Param include_covers query bool false "Include cover photo information"
// @Success 200 {array} models.Album
// @Success 200 {array} models.AlbumWithCover
// @Router /api/albums [get]
func (h *AlbumHandler) ListVisibleAlbums(c *fiber.Ctx) error {
	includeCovers := c.QueryBool("include_covers", false)

	ctx := c.Context()

	if includeCovers {
		albums, err := h.albumRepo.ListVisibleWithCovers(ctx)
		if err != nil {
			logger.Error("Failed to fetch visible albums with covers", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch albums",
			})
		}
		return c.JSON(albums)
	} else {
		albums, err := h.albumRepo.ListVisible(ctx)
		if err != nil {
			logger.Error("Failed to fetch visible albums", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch albums",
			})
		}
		return c.JSON(albums)
	}
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

// GetAlbumPhotos returns photos for a specific album
// @Summary Get photos for album
// @Description Returns all visible photos for a specific album
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Success 200 {array} models.Photo
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/albums/{id}/photos [get]
func (h *AlbumHandler) GetAlbumPhotos(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Album ID is required",
		})
	}

	// Check if album exists and is visible
	ctx := c.Context()
	album, err := h.albumRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch album", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch album",
		})
	}

	if album == nil || !album.Visible {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Album not found",
		})
	}

	photos, err := h.photoRepo.ListByAlbumID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch album photos", "error", err, "album_id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch album photos",
		})
	}

	return c.JSON(photos)
}

// AddPhotoToAlbum adds a photo to an album
// @Summary Add photo to album
// @Description Adds a photo to an album with optional order number
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Param photo body models.AddPhotoToAlbumRequest true "Photo data"
// @Success 201 {object} models.AlbumPhoto
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/albums/{id}/photos [post]
// @Security BearerAuth
func (h *AlbumHandler) AddPhotoToAlbum(c *fiber.Ctx) error {
	albumID := c.Params("id")
	if albumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Album ID is required",
		})
	}

	var req models.AddPhotoToAlbumRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.PhotoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	ctx := c.Context()

	// Check if album exists
	album, err := h.albumRepo.GetByID(ctx, albumID)
	if err != nil {
		logger.Error("Failed to fetch album", "error", err, "id", albumID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch album",
		})
	}

	if album == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Album not found",
		})
	}

	// Check if photo exists
	photo, err := h.photoRepo.GetByID(ctx, req.PhotoID)
	if err != nil {
		logger.Error("Failed to fetch photo", "error", err, "id", req.PhotoID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch photo",
		})
	}

	if photo == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found",
		})
	}

	// Check if photo is already in album
	existing, err := h.albumPhotoRepo.GetByAlbumAndPhoto(ctx, albumID, req.PhotoID)
	if err != nil {
		logger.Error("Failed to check existing album photo", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check album photo relationship",
		})
	}

	if existing != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Photo is already in this album",
		})
	}

	// Create album photo relationship
	albumPhoto := &models.AlbumPhoto{
		AlbumID:     albumID,
		PhotoID:     req.PhotoID,
		OrderNumber: req.OrderNumber,
	}

	if err := h.albumPhotoRepo.Create(ctx, albumPhoto); err != nil {
		logger.Error("Failed to add photo to album", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add photo to album",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(albumPhoto)
}

// RemovePhotoFromAlbum removes a photo from an album
// @Summary Remove photo from album
// @Description Removes a photo from an album
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Param photoId path string true "Photo ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/albums/{id}/photos/{photoId} [delete]
// @Security BearerAuth
func (h *AlbumHandler) RemovePhotoFromAlbum(c *fiber.Ctx) error {
	albumID := c.Params("id")
	photoID := c.Params("photoId")

	if albumID == "" || photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Album ID and Photo ID are required",
		})
	}

	ctx := c.Context()

	// Check if album exists
	album, err := h.albumRepo.GetByID(ctx, albumID)
	if err != nil {
		logger.Error("Failed to fetch album", "error", err, "id", albumID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch album",
		})
	}

	if album == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Album not found",
		})
	}

	// Check if photo exists
	photo, err := h.photoRepo.GetByID(ctx, photoID)
	if err != nil {
		logger.Error("Failed to fetch photo", "error", err, "id", photoID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch photo",
		})
	}

	if photo == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found",
		})
	}

	// Check if photo is in album
	existing, err := h.albumPhotoRepo.GetByAlbumAndPhoto(ctx, albumID, photoID)
	if err != nil {
		logger.Error("Failed to check existing album photo", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check album photo relationship",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo is not in this album",
		})
	}

	// Remove photo from album
	if err := h.albumPhotoRepo.Delete(ctx, albumID, photoID); err != nil {
		logger.Error("Failed to remove photo from album", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove photo from album",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Photo removed from album successfully",
	})
}

// ReorderAlbumPhotos reorders photos in an album
// @Summary Reorder photos in album
// @Description Updates the order of photos in an album
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Param order body models.ReorderPhotosRequest true "Photo order data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/albums/{id}/photos/reorder [put]
// @Security BearerAuth
func (h *AlbumHandler) ReorderAlbumPhotos(c *fiber.Ctx) error {
	albumID := c.Params("id")
	if albumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Album ID is required",
		})
	}

	var req models.ReorderPhotosRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.PhotoOrder) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo order is required",
		})
	}

	ctx := c.Context()

	// Check if album exists
	album, err := h.albumRepo.GetByID(ctx, albumID)
	if err != nil {
		logger.Error("Failed to fetch album", "error", err, "id", albumID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch album",
		})
	}

	if album == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Album not found",
		})
	}

	// Update order for each photo
	for _, photoOrder := range req.PhotoOrder {
		if photoOrder.PhotoID == "" {
			continue
		}

		if err := h.albumPhotoRepo.UpdateOrder(ctx, albumID, photoOrder.PhotoID, photoOrder.OrderNumber); err != nil {
			logger.Error("Failed to update photo order", "error", err, "album_id", albumID, "photo_id", photoOrder.PhotoID)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update photo order",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Photos reordered successfully",
	})
}

// ReorderAlbums reorders albums
// @Summary Reorder albums
// @Description Updates the order of multiple albums
// @Tags Albums
// @Accept json
// @Produce json
// @Param album_order body models.ReorderAlbumsRequest true "Album order data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/albums/reorder [put]
// @Security BearerAuth
func (h *AlbumHandler) ReorderAlbums(c *fiber.Ctx) error {
	var req models.ReorderAlbumsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.AlbumOrder) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Album order is required",
		})
	}

	ctx := c.Context()

	// Update order for each album
	for _, albumOrder := range req.AlbumOrder {
		if albumOrder.ID == "" {
			continue
		}

		if err := h.albumRepo.UpdateOrder(ctx, albumOrder.ID, albumOrder.OrderNumber); err != nil {
			logger.Error("Failed to update album order", "error", err, "album_id", albumOrder.ID)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update album order",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Albums reordered successfully",
	})
}
