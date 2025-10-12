package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SponsorHandler handles sponsor-related HTTP requests
type SponsorHandler struct {
	sponsorRepo       repository.SponsorRepository
	authService       services.AuthService
	permissionService services.PermissionService
	imageService      *services.ImageService
}

// NewSponsorHandler creates a new sponsor handler
func NewSponsorHandler(
	sponsorRepo repository.SponsorRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
	imageService *services.ImageService,
) *SponsorHandler {
	return &SponsorHandler{
		sponsorRepo:       sponsorRepo,
		authService:       authService,
		permissionService: permissionService,
		imageService:      imageService,
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
// @Description Creates a new sponsor. Supports both JSON and multipart/form-data for logo upload.
// @Tags Sponsors
// @Accept json,multipartFormData
// @Produce json
// @Param sponsor body models.Sponsor true "Sponsor data (JSON)"
// @Param logo formData file false "Logo image file"
// @Param name formData string false "Sponsor name (multipart)"
// @Param description formData string false "Sponsor description (multipart)"
// @Param website_url formData string false "Website URL (multipart)"
// @Param order_number formData int false "Order number (multipart)"
// @Param is_active formData bool false "Is active (multipart)"
// @Param visible formData bool false "Is visible (multipart)"
// @Success 201 {object} models.Sponsor
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/sponsors [post]
// @Security BearerAuth
func (h *SponsorHandler) CreateSponsor(c *fiber.Ctx) error {
	var sponsor models.Sponsor

	// Check if request is multipart form data
	contentType := c.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse form data",
			})
		}

		// Parse form fields
		if name := form.Value["name"]; len(name) > 0 {
			sponsor.Name = name[0]
		}
		if description := form.Value["description"]; len(description) > 0 {
			sponsor.Description = description[0]
		}
		if websiteURL := form.Value["website_url"]; len(websiteURL) > 0 {
			sponsor.WebsiteURL = websiteURL[0]
		}
		if orderNumber := form.Value["order_number"]; len(orderNumber) > 0 {
			if orderNum, err := strconv.Atoi(orderNumber[0]); err == nil {
				sponsor.OrderNumber = orderNum
			}
		}
		if isActive := form.Value["is_active"]; len(isActive) > 0 {
			if active, err := strconv.ParseBool(isActive[0]); err == nil {
				sponsor.IsActive = active
			}
		}
		if visible := form.Value["visible"]; len(visible) > 0 {
			if vis, err := strconv.ParseBool(visible[0]); err == nil {
				sponsor.Visible = vis
			}
		}

		// Handle logo upload if present
		if files := form.File["logo"]; len(files) > 0 {
			file := files[0]

			// Validate file type
			if !h.isValidImageType(file.Filename) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid logo file type. Only JPEG, PNG, GIF, and WebP are allowed.",
				})
			}

			// Open file
			src, err := file.Open()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to open uploaded logo file",
				})
			}
			defer src.Close()

			// Upload to Cloudinary
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := h.imageService.UploadImage(ctx, src, file.Filename, "sponsor_logos", "system")
			if err != nil {
				logger.Error("Logo upload failed", "error", err, "filename", file.Filename)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to upload logo",
				})
			}

			sponsor.LogoURL = result.URL
		}
	} else {
		// Handle JSON request
		if err := c.BodyParser(&sponsor); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
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

// isValidImageType controleert of het bestandstype geldig is
func (h *SponsorHandler) isValidImageType(filename string) bool {
	validTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	filename = strings.ToLower(filename)

	for _, ext := range validTypes {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
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
