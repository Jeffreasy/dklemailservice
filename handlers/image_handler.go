package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/services"
	"mime/multipart"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ImageHandler bevat handlers voor image upload en beheer
type ImageHandler struct {
	imageService *services.ImageService
	authService  services.AuthService
}

// NewImageHandler maakt een nieuwe image handler
func NewImageHandler(imageService *services.ImageService, authService services.AuthService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
		authService:  authService,
	}
}

// ValidateImageUpload is een middleware die image uploads valideert
func (h *ImageHandler) ValidateImageUpload(c *fiber.Ctx) error {
	// Check file size (max 10MB)
	if c.Get("Content-Length") != "" {
		contentLength := c.Get("Content-Length")
		if len(contentLength) > 7 { // 10MB = 10485760 (8 digits), so >7 digits means >10MB
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": "File too large. Maximum size is 10MB.",
			})
		}
	}

	// Check content type
	contentType := c.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid content type. Expected multipart/form-data.",
		})
	}

	return c.Next()
}

// UploadImage handelt single image uploads af
// @Summary Upload een enkele image
// @Description Upload een image naar Cloudinary
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {object} services.UploadResult
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 413 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/images/upload [post]
// @Security BearerAuth
func (h *ImageHandler) UploadImage(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	files := form.File["image"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
		})
	}

	file := files[0]

	// Validate file type
	if !h.isValidImageType(file.Filename) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed.",
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Upload to Cloudinary
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := h.imageService.UploadImage(ctx, src, file.Filename, "user_uploads", userID)
	if err != nil {
		logger.Error("Image upload failed", "error", err, "user_id", userID, "filename", file.Filename)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload image",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// UploadBatchImages handelt batch image uploads af
// @Summary Upload meerdere images
// @Description Upload meerdere images tegelijk naar Cloudinary (max 10). Ondersteunt parallel en sequentiële modes.
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param images formData file true "Image files (max 10)"
// @Param mode query string false "Upload mode: 'parallel' (default) or 'sequential'" Enums(parallel,sequential)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 413 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/images/batch-upload [post]
// @Security BearerAuth
func (h *ImageHandler) UploadBatchImages(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get upload mode from query parameter (default: "parallel")
	mode := c.Query("mode", "parallel")

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	files := form.File["images"]
	if len(files) == 0 || len(files) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Provide 1-10 image files",
		})
	}

	// Validate all files first
	validFiles := make([]*multipart.FileHeader, 0, len(files))
	for _, file := range files {
		if h.isValidImageType(file.Filename) {
			validFiles = append(validFiles, file)
		}
	}

	if len(validFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid image files provided",
		})
	}

	var results []*services.UploadResult
	var uploadErrors []error

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second) // Extended timeout for batch
	defer cancel()

	if mode == "sequential" {
		// Sequential upload: process files one by one, more reliable for large files
		results, uploadErrors = h.imageService.UploadBatchImagesSequential(ctx, validFiles, "user_uploads", userID)
	} else {
		// Parallel upload: process all files simultaneously (default, faster for small files)
		results = make([]*services.UploadResult, 0, len(validFiles))

		for _, file := range validFiles {
			src, err := file.Open()
			if err != nil {
				continue
			}

			result, err := h.imageService.UploadImage(ctx, src, file.Filename, "user_uploads", userID)
			src.Close()

			if err == nil {
				results = append(results, result)
			}
		}
	}

	// Return results with error information if any
	response := fiber.Map{
		"success":        true,
		"data":           results,
		"uploaded_count": len(results),
		"total_files":    len(validFiles),
		"mode":           mode,
	}

	if len(uploadErrors) > 0 {
		errorMessages := make([]string, len(uploadErrors))
		for i, err := range uploadErrors {
			errorMessages[i] = err.Error()
		}
		response["errors"] = errorMessages
		response["errors_count"] = len(uploadErrors)
	}

	return c.JSON(response)
}

// GetImageMetadata haalt metadata van een image op
// @Summary Haal image metadata op
// @Description Haalt metadata van een geüploade image op
// @Tags Images
// @Accept json
// @Produce json
// @Param public_id path string true "Cloudinary public ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/images/{public_id} [get]
// @Security BearerAuth
func (h *ImageHandler) GetImageMetadata(c *fiber.Ctx) error {
	publicID := c.Params("public_id")
	if publicID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Public ID required",
		})
	}

	// Get image URL without transformations for metadata
	url := h.imageService.GetImageURL(publicID, map[string]string{})
	if url == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"public_id": publicID,
			"url":       url,
		},
	})
}

// DeleteImage verwijdert een image
// @Summary Verwijder een image
// @Description Verwijdert een image van Cloudinary
// @Tags Images
// @Accept json
// @Produce json
// @Param public_id path string true "Cloudinary public ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/images/{public_id} [delete]
// @Security BearerAuth
func (h *ImageHandler) DeleteImage(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	publicID := c.Params("public_id")
	if publicID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Public ID required",
		})
	}

	// TODO: Check ownership via database if tracking images
	// For now, allow authenticated users to delete their own images

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.imageService.DeleteImage(ctx, publicID)
	if err != nil {
		logger.Error("Image deletion failed", "error", err, "public_id", publicID, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete image",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Image deleted successfully",
	})
}

// RegisterRoutes registreert de image routes
func (h *ImageHandler) RegisterRoutes(app *fiber.App) {
	// Skip registration if image service is not available
	if h.imageService == nil {
		logger.Info("Image service not available, skipping image routes registration")
		return
	}

	api := app.Group("/api")

	// Protected routes
	imagesGroup := api.Group("/images", AuthMiddleware(h.authService))
	imagesGroup.Post("/upload", h.ValidateImageUpload, h.UploadImage)
	imagesGroup.Post("/batch-upload", h.ValidateImageUpload, h.UploadBatchImages)
	imagesGroup.Get("/:public_id", h.GetImageMetadata)
	imagesGroup.Delete("/:public_id", h.DeleteImage)

	logger.Info("Image routes registered")
}

// isValidImageType controleert of het bestandstype geldig is
func (h *ImageHandler) isValidImageType(filename string) bool {
	validTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	filename = strings.ToLower(filename)

	for _, ext := range validTypes {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
