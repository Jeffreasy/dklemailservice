package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// RadioRecordingHandler handles radio recording-related HTTP requests
type RadioRecordingHandler struct {
	radioRecordingRepo repository.RadioRecordingRepository
	authService        services.AuthService
	permissionService  services.PermissionService
}

// NewRadioRecordingHandler creates a new radio recording handler
func NewRadioRecordingHandler(
	radioRecordingRepo repository.RadioRecordingRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *RadioRecordingHandler {
	return &RadioRecordingHandler{
		radioRecordingRepo: radioRecordingRepo,
		authService:        authService,
		permissionService:  permissionService,
	}
}

// RegisterRoutes registers the radio recording routes
func (h *RadioRecordingHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/radio-recordings")
	public.Get("/", h.ListVisibleRadioRecordings)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/radio-recordings", AuthMiddleware(h.authService))

	// Read routes (require radio_recording read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "radio_recording", "read"))
	readGroup.Get("/admin", h.ListRadioRecordings)
	readGroup.Get("/:id", h.GetRadioRecording)

	// Write routes (require radio_recording write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "radio_recording", "write"))
	writeGroup.Post("/", h.CreateRadioRecording)
	writeGroup.Put("/:id", h.UpdateRadioRecording)

	// Delete routes (require radio_recording delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "radio_recording", "delete"))
	deleteGroup.Delete("/:id", h.DeleteRadioRecording)
}

// ListVisibleRadioRecordings returns all visible radio recordings for public display
// @Summary Get visible radio recordings
// @Description Returns all visible radio recordings ordered by order_number
// @Tags Radio Recordings
// @Accept json
// @Produce json
// @Success 200 {array} models.RadioRecording
// @Router /api/radio-recordings [get]
func (h *RadioRecordingHandler) ListVisibleRadioRecordings(c *fiber.Ctx) error {
	ctx := c.Context()
	recordings, err := h.radioRecordingRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible radio recordings", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch radio recordings",
		})
	}

	return c.JSON(recordings)
}

// ListRadioRecordings returns all radio recordings for admin management
// @Summary List all radio recordings
// @Description Returns a paginated list of all radio recordings
// @Tags Radio Recordings
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.RadioRecording
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/radio-recordings/admin [get]
// @Security BearerAuth
func (h *RadioRecordingHandler) ListRadioRecordings(c *fiber.Ctx) error {
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
	recordings, err := h.radioRecordingRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch radio recordings", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch radio recordings",
		})
	}

	return c.JSON(recordings)
}

// GetRadioRecording returns a specific radio recording
// @Summary Get radio recording by ID
// @Description Returns a specific radio recording by ID
// @Tags Radio Recordings
// @Accept json
// @Produce json
// @Param id path string true "Radio Recording ID"
// @Success 200 {object} models.RadioRecording
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/radio-recordings/{id} [get]
// @Security BearerAuth
func (h *RadioRecordingHandler) GetRadioRecording(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Radio recording ID is required",
		})
	}

	ctx := c.Context()
	recording, err := h.radioRecordingRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch radio recording", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch radio recording",
		})
	}

	if recording == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Radio recording not found",
		})
	}

	return c.JSON(recording)
}

// CreateRadioRecording creates a new radio recording
// @Summary Create radio recording
// @Description Creates a new radio recording
// @Tags Radio Recordings
// @Accept json
// @Produce json
// @Param recording body models.RadioRecording true "Radio recording data"
// @Success 201 {object} models.RadioRecording
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/radio-recordings [post]
// @Security BearerAuth
func (h *RadioRecordingHandler) CreateRadioRecording(c *fiber.Ctx) error {
	var recording models.RadioRecording
	if err := c.BodyParser(&recording); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if recording.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	ctx := c.Context()
	if err := h.radioRecordingRepo.Create(ctx, &recording); err != nil {
		logger.Error("Failed to create radio recording", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create radio recording",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(recording)
}

// UpdateRadioRecording updates an existing radio recording
// @Summary Update radio recording
// @Description Updates an existing radio recording
// @Tags Radio Recordings
// @Accept json
// @Produce json
// @Param id path string true "Radio Recording ID"
// @Param recording body models.RadioRecording true "Radio recording data"
// @Success 200 {object} models.RadioRecording
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/radio-recordings/{id} [put]
// @Security BearerAuth
func (h *RadioRecordingHandler) UpdateRadioRecording(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Radio recording ID is required",
		})
	}

	// Get existing recording
	ctx := c.Context()
	existing, err := h.radioRecordingRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch radio recording", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch radio recording",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Radio recording not found",
		})
	}

	// Parse update data
	var updateData models.RadioRecording
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.Title = updateData.Title
	existing.Description = updateData.Description
	existing.Date = updateData.Date
	existing.AudioURL = updateData.AudioURL
	existing.ThumbnailURL = updateData.ThumbnailURL
	existing.Visible = updateData.Visible
	existing.OrderNumber = updateData.OrderNumber

	if err := h.radioRecordingRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update radio recording", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update radio recording",
		})
	}

	return c.JSON(existing)
}

// DeleteRadioRecording deletes a radio recording
// @Summary Delete radio recording
// @Description Deletes a radio recording
// @Tags Radio Recordings
// @Accept json
// @Produce json
// @Param id path string true "Radio Recording ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/radio-recordings/{id} [delete]
// @Security BearerAuth
func (h *RadioRecordingHandler) DeleteRadioRecording(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Radio recording ID is required",
		})
	}

	// Check if recording exists
	ctx := c.Context()
	recording, err := h.radioRecordingRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch radio recording", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch radio recording",
		})
	}

	if recording == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Radio recording not found",
		})
	}

	if err := h.radioRecordingRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete radio recording", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete radio recording",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Radio recording deleted successfully",
	})
}
