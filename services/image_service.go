package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"dklautomationgo/config"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Custom errors
var (
	ErrUploadFailed    = errors.New("image upload failed")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file too large")
	ErrInvalidFile     = errors.New("invalid file")
)

type UploadResult struct {
	PublicID     string `json:"public_id"`
	URL          string `json:"url"`
	SecureURL    string `json:"secure_url"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Format       string `json:"format"`
	Bytes        int    `json:"bytes"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

type ImageService struct {
	cld    *cloudinary.Cloudinary
	config *config.CloudinaryConfig
	repo   repository.UploadedImageRepository
}

func NewImageService(config *config.CloudinaryConfig, repo repository.UploadedImageRepository) (*ImageService, error) {
	cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary client: %w", err)
	}

	return &ImageService{
		cld:    cld,
		config: config,
		repo:   repo,
	}, nil
}

func (s *ImageService) UploadImage(ctx context.Context, file multipart.File, filename string, folder string, userID string) (*UploadResult, error) {
	// Generate unique public ID
	publicID := s.generatePublicID(userID, filename)

	// Upload parameters
	uploadParams := uploader.UploadParams{
		PublicID:     publicID,
		Folder:       folder,
		ResourceType: "image",
		Format:       "auto", // Auto-detect format
	}

	// Add upload preset if configured
	if s.config.UploadPreset != "" {
		uploadParams.UploadPreset = s.config.UploadPreset
	}

	// TODO: Add eager transformations for thumbnails in Phase 2
	// uploadParams.Eager = "w_200,h_200,c_thumb,g_face"

	// Perform upload
	resp, err := s.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		logger.Error("Cloudinary upload failed",
			"error", err,
			"public_id", publicID,
			"folder", folder,
			"user_id", userID)
		return nil, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	result := &UploadResult{
		PublicID:  resp.PublicID,
		URL:       resp.URL,
		SecureURL: resp.SecureURL,
		Width:     resp.Width,
		Height:    resp.Height,
		Format:    resp.Format,
		Bytes:     resp.Bytes,
	}

	// TODO: Generate thumbnail URL if eager transformation exists (Phase 2)
	// if len(resp.Eager) > 0 {
	//     result.ThumbnailURL = resp.Eager[0].SecureURL
	// }

	// Create database record for tracking
	uploadedImage := &models.UploadedImage{
		UserID:       userID,
		PublicID:     result.PublicID,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		Filename:     filename,
		Size:         int64(result.Bytes),
		MimeType:     "image/" + result.Format, // Infer MIME type from format
		Width:        result.Width,
		Height:       result.Height,
		Folder:       folder,
		ThumbnailURL: &result.ThumbnailURL,
	}

	if err := s.repo.Create(ctx, uploadedImage); err != nil {
		logger.Error("Failed to create uploaded image record",
			"error", err,
			"public_id", result.PublicID,
			"user_id", userID)
		// Don't fail the upload if database record creation fails
		// The image is already uploaded to Cloudinary
	}

	logger.Info("Image uploaded successfully",
		"public_id", result.PublicID,
		"user_id", userID,
		"filename", filename,
		"size_bytes", result.Bytes,
		"folder", folder)

	return result, nil
}

func (s *ImageService) UploadImageAsync(ctx context.Context, file multipart.File, filename string, folder string, userID string, callbackURL string) (string, error) {
	publicID := s.generatePublicID(userID, filename)

	async := true
	uploadParams := uploader.UploadParams{
		PublicID:        publicID,
		Folder:          folder,
		ResourceType:    "image",
		Async:           &async,
		NotificationURL: callbackURL,
	}

	resp, err := s.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		logger.Error("Async Cloudinary upload failed",
			"error", err,
			"public_id", publicID,
			"folder", folder,
			"user_id", userID)
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	logger.Info("Async image upload initiated",
		"public_id", publicID,
		"user_id", userID,
		"filename", filename,
		"folder", folder)

	return resp.PublicID, nil
}

// UploadBatchImagesSequential uploads multiple images sequentially with individual error handling
func (s *ImageService) UploadBatchImagesSequential(ctx context.Context, files []*multipart.FileHeader, folder string, userID string) ([]*UploadResult, []error) {
	results := make([]*UploadResult, 0, len(files))
	errors := make([]error, 0)

	for i, fileHeader := range files {
		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to open file %d (%s): %w", i, fileHeader.Filename, err))
			continue
		}

		// Upload image
		result, err := s.UploadImage(ctx, file, fileHeader.Filename, folder, userID)
		file.Close() // Always close the file

		if err != nil {
			errors = append(errors, fmt.Errorf("failed to upload file %d (%s): %w", i, fileHeader.Filename, err))
			continue
		}

		results = append(results, result)

		logger.Info("Sequential batch upload progress",
			"file_index", i,
			"filename", fileHeader.Filename,
			"public_id", result.PublicID,
			"total_files", len(files),
			"uploaded_count", len(results))
	}

	logger.Info("Sequential batch upload completed",
		"total_files", len(files),
		"successful_uploads", len(results),
		"errors_count", len(errors),
		"user_id", userID,
		"folder", folder)

	return results, errors
}

func (s *ImageService) DeleteImage(ctx context.Context, publicID string) error {
	invalidate := true
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
		Invalidate:   &invalidate, // Clear CDN cache
	})

	if err != nil {
		logger.Error("Cloudinary delete failed",
			"error", err,
			"public_id", publicID)
		return fmt.Errorf("failed to delete image: %w", err)
	}

	// Mark database record as deleted (soft delete)
	if image, err := s.repo.GetByPublicID(ctx, publicID); err == nil && image != nil {
		if err := s.repo.SoftDelete(ctx, image.ID); err != nil {
			logger.Warn("Failed to soft delete image record",
				"error", err,
				"public_id", publicID,
				"record_id", image.ID)
			// Don't fail the operation if database update fails
		}
	} else {
		logger.Warn("Could not find image record for soft delete",
			"public_id", publicID,
			"error", err)
	}

	logger.Info("Image deleted successfully",
		"public_id", publicID)

	return nil
}

func (s *ImageService) GetImageURL(publicID string, transformations map[string]string) string {
	// Build transformation string
	var transformParts []string
	for key, value := range transformations {
		switch strings.ToLower(key) {
		case "width":
			transformParts = append(transformParts, "w_"+value)
		case "height":
			transformParts = append(transformParts, "h_"+value)
		case "crop":
			transformParts = append(transformParts, "c_"+value)
		case "quality":
			transformParts = append(transformParts, "q_"+value)
		case "format":
			transformParts = append(transformParts, "f_"+value)
		}
	}

	transformStr := ""
	if len(transformParts) > 0 {
		transformStr = strings.Join(transformParts, ",") + "/"
	}

	// Build URL manually for now
	protocol := "http"
	if s.config.Secure {
		protocol = "https"
	}

	url := fmt.Sprintf("%s://res.cloudinary.com/%s/image/upload/%s%s",
		protocol, s.config.CloudName, transformStr, publicID)

	return url
}

func (s *ImageService) generatePublicID(userID, filename string) string {
	// Generate unique ID combining user ID and random bytes
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)

	// Extract file extension
	parts := strings.Split(filename, ".")
	ext := ""
	if len(parts) > 1 {
		ext = parts[len(parts)-1]
	}

	// Create public ID with folder structure
	publicID := fmt.Sprintf("%s/%s_%s", s.config.UploadFolder, userID, randomStr)
	if ext != "" {
		publicID += "." + ext
	}

	return publicID
}
