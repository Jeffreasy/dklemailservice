package handlers

import (
	"dklautomationgo/models"
	"dklautomationgo/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// NotulenHandler handles HTTP requests for notulen
type NotulenHandler struct {
	service           *services.NotulenService
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewNotulenHandler creates a new notulen handler
func NewNotulenHandler(service *services.NotulenService, authService services.AuthService, permissionService services.PermissionService) *NotulenHandler {
	return &NotulenHandler{
		service:           service,
		authService:       authService,
		permissionService: permissionService,
	}
}

// CreateNotulen creates a new notulen document
func (h *NotulenHandler) CreateNotulen(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req models.NotulenCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.service.ValidateNotulen(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	notulen, err := h.service.CreateNotulen(c.Context(), userUUID, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(notulen)
}

// GetNotulen retrieves a notulen by ID
func (h *NotulenHandler) GetNotulen(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notulen ID",
		})
	}

	notulenResponse, err := h.service.GetNotulen(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if notulenResponse == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Notulen not found",
		})
	}

	// Check if markdown format is requested
	if c.Query("format") == "markdown" {
		// Convert response back to model for markdown rendering
		notulen := h.convertResponseToNotulen(notulenResponse)
		markdown, err := h.service.RenderMarkdown(notulen)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to render markdown",
			})
		}
		c.Set("Content-Type", "text/markdown")
		return c.SendString(markdown)
	}

	return c.JSON(notulenResponse)
}

// UpdateNotulen updates an existing notulen
func (h *NotulenHandler) UpdateNotulen(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notulen ID",
		})
	}

	var req models.NotulenUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	notulen, err := h.service.UpdateNotulen(c.Context(), userUUID, id, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(notulen)
}

// FinalizeNotulen finalizes a notulen document
func (h *NotulenHandler) FinalizeNotulen(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notulen ID",
		})
	}

	var req models.NotulenFinalizeRequest
	if err := c.BodyParser(&req); err != nil {
		req = models.NotulenFinalizeRequest{} // Empty request is valid
	}

	if err := h.service.FinalizeNotulen(c.Context(), userUUID, id, &req); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notulen finalized successfully",
	})
}

// ArchiveNotulen archives a notulen document
func (h *NotulenHandler) ArchiveNotulen(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notulen ID",
		})
	}

	if err := h.service.ArchiveNotulen(c.Context(), userUUID, id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notulen archived successfully",
	})
}

// DeleteNotulen deletes a notulen document
func (h *NotulenHandler) DeleteNotulen(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notulen ID",
		})
	}

	if err := h.service.DeleteNotulen(c.Context(), userUUID, id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notulen deleted successfully",
	})
}

// ListNotulen lists notulen with filtering and pagination
func (h *NotulenHandler) ListNotulen(c *fiber.Ctx) error {
	filters := &models.NotulenSearchFilters{}

	// Parse query parameters
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &parsed
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &parsed
		}
	}

	if status := c.Query("status"); status != "" {
		filters.Status = status
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	notulen, total, err := h.service.ListNotulen(c.Context(), filters)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(models.NotulenListResponse{
		Notulen: notulen,
		Total:   total,
		Limit:   filters.Limit,
		Offset:  filters.Offset,
	})
}

// convertResponseToNotulen converts a NotulenResponse back to Notulen model for compatibility
func (h *NotulenHandler) convertResponseToNotulen(response *models.NotulenResponse) *models.Notulen {
	return &models.Notulen{
		ID:                   response.ID,
		Titel:                response.Titel,
		VergaderingDatum:     response.VergaderingDatum,
		Locatie:              response.Locatie,
		Voorzitter:           response.Voorzitter,
		Notulist:             response.Notulist,
		Aanwezigen:           response.Aanwezigen,
		Afwezigen:            response.Afwezigen,
		AanwezigenGebruikers: response.AanwezigenGebruikers,
		AfwezigenGebruikers:  response.AfwezigenGebruikers,
		AanwezigenGasten:     response.AanwezigenGasten,
		AfwezigenGasten:      response.AfwezigenGasten,
		AgendaItems:          response.AgendaItems,
		Besluiten:            response.Besluiten,
		Actiepunten:          response.Actiepunten,
		Notities:             response.Notities,
		Status:               response.Status,
		Versie:               response.Versie,
		CreatedBy:            response.CreatedBy,
		CreatedAt:            response.CreatedAt,
		UpdatedAt:            response.UpdatedAt,
		UpdatedBy:            response.UpdatedBy,
		FinalizedAt:          response.FinalizedAt,
		FinalizedBy:          response.FinalizedBy,
	}
}

// SearchNotulen performs full-text search on notulen
func (h *NotulenHandler) SearchNotulen(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	filters := &models.NotulenSearchFilters{}

	// Parse additional filters
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &parsed
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &parsed
		}
	}

	if status := c.Query("status"); status != "" {
		filters.Status = status
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	notulen, total, err := h.service.SearchNotulen(c.Context(), query, filters)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(models.NotulenListResponse{
		Notulen: notulen,
		Total:   total,
		Limit:   filters.Limit,
		Offset:  filters.Offset,
	})
}

// GetNotulenVersions retrieves all versions of a notulen
func (h *NotulenHandler) GetNotulenVersions(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notulen ID",
		})
	}

	versions, err := h.service.GetNotulenVersions(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"versions": versions,
	})
}

// GetNotulenVersion retrieves a specific version of a notulen
func (h *NotulenHandler) GetNotulenVersion(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notulen ID",
		})
	}

	versionStr := c.Params("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid version number",
		})
	}

	notulenVersion, err := h.service.GetNotulenVersion(c.Context(), id, version)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if notulenVersion == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Version not found",
		})
	}

	return c.JSON(notulenVersion)
}

// RegisterRoutes registers all notulen routes
func (h *NotulenHandler) RegisterRoutes(app *fiber.App) {
	// Notulen routes group
	notulen := app.Group("/api/notulen")

	// Public routes (no auth required for reading finalized notulen)
	notulen.Get("/public", h.ListPublicNotulen)

	// Protected routes (authentication required)
	protected := notulen.Group("/", AuthMiddleware(h.authService))
	protected.Use(PermissionMiddleware(h.permissionService, "notulen", "read"))

	// CRUD operations
	protected.Post("/", PermissionMiddleware(h.permissionService, "notulen", "write"), h.CreateNotulen)
	protected.Get("/", h.ListNotulen)
	protected.Get("/:id", h.GetNotulen)
	protected.Put("/:id", PermissionMiddleware(h.permissionService, "notulen", "write"), h.UpdateNotulen)
	protected.Delete("/:id", PermissionMiddleware(h.permissionService, "notulen", "delete"), h.DeleteNotulen)

	// Special operations
	protected.Post("/:id/finalize", PermissionMiddleware(h.permissionService, "notulen", "write"), h.FinalizeNotulen)
	protected.Post("/:id/archive", PermissionMiddleware(h.permissionService, "notulen", "write"), h.ArchiveNotulen)

	// Search and versioning
	protected.Get("/search", h.SearchNotulen)
	protected.Get("/:id/versions", h.GetNotulenVersions)
	protected.Get("/:id/versions/:version", h.GetNotulenVersion)
}

// ListPublicNotulen lists only finalized notulen (public access)
func (h *NotulenHandler) ListPublicNotulen(c *fiber.Ctx) error {
	filters := &models.NotulenSearchFilters{
		Status: "finalized", // Only show finalized notulen publicly
	}

	// Parse query parameters
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &parsed
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &parsed
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 50 {
			filters.Limit = limit
		}
	}

	notulen, total, err := h.service.ListNotulen(c.Context(), filters)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(models.NotulenListResponse{
		Notulen: notulen,
		Total:   total,
		Limit:   filters.Limit,
		Offset:  filters.Offset,
	})
}
