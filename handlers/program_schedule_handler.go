package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// ProgramScheduleHandler handles program schedule-related HTTP requests
type ProgramScheduleHandler struct {
	programScheduleRepo repository.ProgramScheduleRepository
	authService         services.AuthService
	permissionService   services.PermissionService
}

// NewProgramScheduleHandler creates a new program schedule handler
func NewProgramScheduleHandler(
	programScheduleRepo repository.ProgramScheduleRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *ProgramScheduleHandler {
	return &ProgramScheduleHandler{
		programScheduleRepo: programScheduleRepo,
		authService:         authService,
		permissionService:   permissionService,
	}
}

// RegisterRoutes registers the program schedule routes
func (h *ProgramScheduleHandler) RegisterRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	public := app.Group("/api/program-schedule")
	public.Get("/", h.ListVisibleProgramSchedules)

	// Admin routes (require authentication and permissions)
	admin := app.Group("/api/program-schedule", AuthMiddleware(h.authService))

	// Read routes (require program_schedule read permission)
	readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "program_schedule", "read"))
	readGroup.Get("/admin", h.ListProgramSchedules)
	readGroup.Get("/:id", h.GetProgramSchedule)

	// Write routes (require program_schedule write permission)
	writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "program_schedule", "write"))
	writeGroup.Post("/", h.CreateProgramSchedule)
	writeGroup.Put("/:id", h.UpdateProgramSchedule)

	// Delete routes (require program_schedule delete permission)
	deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "program_schedule", "delete"))
	deleteGroup.Delete("/:id", h.DeleteProgramSchedule)
}

// ListVisibleProgramSchedules returns all visible program schedules for public display
// @Summary Get visible program schedules
// @Description Returns all visible program schedules ordered by order_number
// @Tags Program Schedule
// @Accept json
// @Produce json
// @Success 200 {array} models.ProgramSchedule
// @Router /api/program-schedule [get]
func (h *ProgramScheduleHandler) ListVisibleProgramSchedules(c *fiber.Ctx) error {
	ctx := c.Context()
	schedules, err := h.programScheduleRepo.ListVisible(ctx)
	if err != nil {
		logger.Error("Failed to fetch visible program schedules", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch program schedules",
		})
	}

	return c.JSON(schedules)
}

// ListProgramSchedules returns all program schedules for admin management
// @Summary List all program schedules
// @Description Returns a paginated list of all program schedules
// @Tags Program Schedule
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} models.ProgramSchedule
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/program-schedule/admin [get]
// @Security BearerAuth
func (h *ProgramScheduleHandler) ListProgramSchedules(c *fiber.Ctx) error {
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
	schedules, err := h.programScheduleRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch program schedules", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch program schedules",
		})
	}

	return c.JSON(schedules)
}

// GetProgramSchedule returns a specific program schedule
// @Summary Get program schedule by ID
// @Description Returns a specific program schedule by ID
// @Tags Program Schedule
// @Accept json
// @Produce json
// @Param id path string true "Program Schedule ID"
// @Success 200 {object} models.ProgramSchedule
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/program-schedule/{id} [get]
// @Security BearerAuth
func (h *ProgramScheduleHandler) GetProgramSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program Schedule ID is required",
		})
	}

	ctx := c.Context()
	schedule, err := h.programScheduleRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch program schedule", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch program schedule",
		})
	}

	if schedule == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Program schedule not found",
		})
	}

	return c.JSON(schedule)
}

// CreateProgramSchedule creates a new program schedule
// @Summary Create program schedule
// @Description Creates a new program schedule
// @Tags Program Schedule
// @Accept json
// @Produce json
// @Param schedule body models.ProgramSchedule true "Program Schedule data"
// @Success 201 {object} models.ProgramSchedule
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/program-schedule [post]
// @Security BearerAuth
func (h *ProgramScheduleHandler) CreateProgramSchedule(c *fiber.Ctx) error {
	var schedule models.ProgramSchedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if schedule.Time == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Time is required",
		})
	}

	if schedule.EventDescription == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event description is required",
		})
	}

	ctx := c.Context()
	if err := h.programScheduleRepo.Create(ctx, &schedule); err != nil {
		logger.Error("Failed to create program schedule", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create program schedule",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(schedule)
}

// UpdateProgramSchedule updates an existing program schedule
// @Summary Update program schedule
// @Description Updates an existing program schedule
// @Tags Program Schedule
// @Accept json
// @Produce json
// @Param id path string true "Program Schedule ID"
// @Param schedule body models.ProgramSchedule true "Program Schedule data"
// @Success 200 {object} models.ProgramSchedule
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/program-schedule/{id} [put]
// @Security BearerAuth
func (h *ProgramScheduleHandler) UpdateProgramSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program Schedule ID is required",
		})
	}

	// Get existing program schedule
	ctx := c.Context()
	existing, err := h.programScheduleRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch program schedule", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch program schedule",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Program schedule not found",
		})
	}

	// Parse update data
	var updateData models.ProgramSchedule
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.Time = updateData.Time
	existing.EventDescription = updateData.EventDescription
	existing.Category = updateData.Category
	existing.IconName = updateData.IconName
	existing.OrderNumber = updateData.OrderNumber
	existing.Visible = updateData.Visible
	existing.Latitude = updateData.Latitude
	existing.Longitude = updateData.Longitude

	if err := h.programScheduleRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update program schedule", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update program schedule",
		})
	}

	return c.JSON(existing)
}

// DeleteProgramSchedule deletes a program schedule
// @Summary Delete program schedule
// @Description Deletes a program schedule
// @Tags Program Schedule
// @Accept json
// @Produce json
// @Param id path string true "Program Schedule ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/program-schedule/{id} [delete]
// @Security BearerAuth
func (h *ProgramScheduleHandler) DeleteProgramSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program Schedule ID is required",
		})
	}

	// Check if program schedule exists
	ctx := c.Context()
	schedule, err := h.programScheduleRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch program schedule", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch program schedule",
		})
	}

	if schedule == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Program schedule not found",
		})
	}

	if err := h.programScheduleRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete program schedule", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete program schedule",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Program schedule deleted successfully",
	})
}
