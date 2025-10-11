package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// VideoHandler handles video-related HTTP requests
type VideoHandler struct {
	videoRepo         repository.VideoRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewVideoHandler creates a new video handler
func NewVideoHandler(
	videoRepo repository.VideoRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *VideoHandler {
	return &VideoHandler{
		videoRepo:         videoRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the video routes
func (h *VideoHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/videos")
	public.Get("/", h.ListVisibleVideos)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/videos", AuthMiddleware(h.authService))

	// Read routes (require video read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "video", "read"))
	readGroup.Get("/admin", h.ListVideos)
	readGroup.Get("/:id", h.GetVideo)

	// Write routes (require video write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "video", "write"))
	writeGroup.Post("/", h.CreateVideo)
	writeGroup.Put("/:id", h.UpdateVideo)

	// Delete routes (require video delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "video", "delete"))
	deleteGroup.Delete("/:id", h.DeleteVideo)
}

// ListVisibleVideos returns all visible videos for public display
// @Summary Get visible videos
// @Description Returns all visible videos ordered by order_number
// @Tags Videos
// @Accept json
// @Produce json
// @Success 200 {array} models.Video
// @Router /api/videos [get]
func (h *VideoHandler) ListVisibleVideos(c *fiber.Ctx) error {
	ctx := c.Context()
	videos, err := h.videoRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible videos", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch videos",
		})
	}

	return c.JSON(videos)
}

// ListVideos returns all videos for admin management
// @Summary List all videos
// @Description Returns a paginated list of all videos
// @Tags Videos
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.Video
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/videos/admin [get]
// @Security BearerAuth
func (h *VideoHandler) ListVideos(c *fiber.Ctx) error {
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
	videos, err := h.videoRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch videos", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch videos",
		})
	}

	return c.JSON(videos)
}

// GetVideo returns a specific video
// @Summary Get video by ID
// @Description Returns a specific video by ID
// @Tags Videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} models.Video
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/videos/{id} [get]
// @Security BearerAuth
func (h *VideoHandler) GetVideo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video ID is required",
		})
	}

	ctx := c.Context()
	video, err := h.videoRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch video", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch video",
		})
	}

	if video == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Video not found",
		})
	}

	return c.JSON(video)
}

// CreateVideo creates a new video
// @Summary Create video
// @Description Creates a new video
// @Tags Videos
// @Accept json
// @Produce json
// @Param video body models.Video true "Video data"
// @Success 201 {object} models.Video
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/videos [post]
// @Security BearerAuth
func (h *VideoHandler) CreateVideo(c *fiber.Ctx) error {
	var video models.Video
	if err := c.BodyParser(&video); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if video.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	ctx := c.Context()
	if err := h.videoRepo.Create(ctx, &video); err != nil {
		logger.Error("Failed to create video", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create video",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(video)
}

// UpdateVideo updates an existing video
// @Summary Update video
// @Description Updates an existing video
// @Tags Videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param video body models.Video true "Video data"
// @Success 200 {object} models.Video
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/videos/{id} [put]
// @Security BearerAuth
func (h *VideoHandler) UpdateVideo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video ID is required",
		})
	}

	// Get existing video
	ctx := c.Context()
	existing, err := h.videoRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch video", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch video",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Video not found",
		})
	}

	// Parse update data
	var updateData models.Video
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.VideoID = updateData.VideoID
	existing.URL = updateData.URL
	existing.Title = updateData.Title
	existing.Description = updateData.Description
	existing.ThumbnailURL = updateData.ThumbnailURL
	existing.Visible = updateData.Visible
	existing.OrderNumber = updateData.OrderNumber

	if err := h.videoRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update video", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update video",
		})
	}

	return c.JSON(existing)
}

// DeleteVideo deletes a video
// @Summary Delete video
// @Description Deletes a video
// @Tags Videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/videos/{id} [delete]
// @Security BearerAuth
func (h *VideoHandler) DeleteVideo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video ID is required",
		})
	}

	// Check if video exists
	ctx := c.Context()
	video, err := h.videoRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch video", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch video",
		})
	}

	if video == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Video not found",
		})
	}

	if err := h.videoRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete video", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete video",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Video deleted successfully",
	})
}
