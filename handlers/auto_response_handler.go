package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// AutoResponseHandler handles auto response-related HTTP requests
type AutoResponseHandler struct {
	autoResponseRepo  repository.AutoResponseRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewAutoResponseHandler creates a new auto response handler
func NewAutoResponseHandler(
	autoResponseRepo repository.AutoResponseRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *AutoResponseHandler {
	return &AutoResponseHandler{
		autoResponseRepo:  autoResponseRepo,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registers the auto response routes
func (h *AutoResponseHandler) RegisterRoutes(app *fiber.App) {
	// Admin routes (require authentication)
	admin := app.Group("/api/mail/autoresponse", AuthMiddleware(h.authService))

	// Read routes
	admin.Get("", h.ListAutoResponses)
	admin.Get("/:id", h.GetAutoResponse)

	// Write routes
	admin.Post("", h.CreateAutoResponse)
	admin.Put("/:id", h.UpdateAutoResponse)

	// Delete routes
	admin.Delete("/:id", h.DeleteAutoResponse)
}

// ListAutoResponses returns all auto responses
func (h *AutoResponseHandler) ListAutoResponses(c *fiber.Ctx) error {
	ctx := c.Context()
	responses, err := h.autoResponseRepo.GetAll(ctx)
	if err != nil {
		logger.Error("Failed to fetch auto responses", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch auto responses",
		})
	}

	// Return empty array if no responses found
	if responses == nil {
		responses = []*models.AutoResponse{}
	}

	return c.JSON(responses)
}

// GetAutoResponse returns a specific auto response
func (h *AutoResponseHandler) GetAutoResponse(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Auto Response ID is required",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Auto Response ID",
		})
	}

	ctx := c.Context()
	response, err := h.autoResponseRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch auto response", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch auto response",
		})
	}

	if response == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Auto response not found",
		})
	}

	return c.JSON(response)
}

// CreateAutoResponse creates a new auto response
func (h *AutoResponseHandler) CreateAutoResponse(c *fiber.Ctx) error {
	var response models.AutoResponse
	if err := c.BodyParser(&response); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if response.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	ctx := c.Context()
	if err := h.autoResponseRepo.Create(ctx, &response); err != nil {
		logger.Error("Failed to create auto response", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create auto response",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateAutoResponse updates an existing auto response
func (h *AutoResponseHandler) UpdateAutoResponse(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Auto Response ID is required",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Auto Response ID",
		})
	}

	ctx := c.Context()
	existing, err := h.autoResponseRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch auto response", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch auto response",
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Auto response not found",
		})
	}

	var updateData models.AutoResponse
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existing.Email = updateData.Email
	existing.IsActive = updateData.IsActive
	existing.Subject = updateData.Subject
	existing.Message = updateData.Message
	existing.StartDate = updateData.StartDate
	existing.EndDate = updateData.EndDate

	if err := h.autoResponseRepo.Update(ctx, existing); err != nil {
		logger.Error("Failed to update auto response", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update auto response",
		})
	}

	return c.JSON(existing)
}

// DeleteAutoResponse deletes an auto response
func (h *AutoResponseHandler) DeleteAutoResponse(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Auto Response ID is required",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Auto Response ID",
		})
	}

	ctx := c.Context()
	response, err := h.autoResponseRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch auto response", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch auto response",
		})
	}

	if response == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Auto response not found",
		})
	}

	if err := h.autoResponseRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete auto response", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete auto response",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Auto response deleted successfully",
	})
}
