package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// PostgresUploadedImageRepository implements UploadedImageRepository with PostgreSQL
type PostgresUploadedImageRepository struct {
	*PostgresRepository
}

// NewPostgresUploadedImageRepository creates a new PostgreSQL uploaded image repository
func NewPostgresUploadedImageRepository(base *PostgresRepository) *PostgresUploadedImageRepository {
	return &PostgresUploadedImageRepository{
		PostgresRepository: base,
	}
}

// Create saves a new uploaded image record
func (r *PostgresUploadedImageRepository) Create(ctx context.Context, image *models.UploadedImage) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(image)
	return r.handleError("Create", result.Error)
}

// GetByID retrieves an uploaded image by ID
func (r *PostgresUploadedImageRepository) GetByID(ctx context.Context, id string) (*models.UploadedImage, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var image models.UploadedImage
	result := r.DB().WithContext(ctx).First(&image, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &image, nil
}

// GetByPublicID retrieves an uploaded image by Cloudinary public ID
func (r *PostgresUploadedImageRepository) GetByPublicID(ctx context.Context, publicID string) (*models.UploadedImage, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var image models.UploadedImage
	result := r.DB().WithContext(ctx).Where("public_id = ?", publicID).First(&image)
	if err := r.handleError("GetByPublicID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &image, nil
}

// GetByUserID retrieves uploaded images for a user with pagination
func (r *PostgresUploadedImageRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.UploadedImage, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var images []*models.UploadedImage
	result := r.DB().WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&images)

	if err := r.handleError("GetByUserID", result.Error); err != nil {
		return nil, err
	}

	return images, nil
}

// List retrieves a paginated list of all uploaded images
func (r *PostgresUploadedImageRepository) List(ctx context.Context, limit, offset int) ([]*models.UploadedImage, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var images []*models.UploadedImage
	result := r.DB().WithContext(ctx).
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&images)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return images, nil
}

// Update updates an existing uploaded image record
func (r *PostgresUploadedImageRepository) Update(ctx context.Context, image *models.UploadedImage) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(image)
	return r.handleError("Update", result.Error)
}

// Delete removes an uploaded image record
func (r *PostgresUploadedImageRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.UploadedImage{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// SoftDelete marks an image as deleted (for GDPR compliance)
func (r *PostgresUploadedImageRepository) SoftDelete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.UploadedImage{}).
		Where("id = ?", id).
		Update("deleted_at", now)

	return r.handleError("SoftDelete", result.Error)
}

// GetByFolder retrieves images by folder
func (r *PostgresUploadedImageRepository) GetByFolder(ctx context.Context, folder string, limit, offset int) ([]*models.UploadedImage, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var images []*models.UploadedImage
	result := r.DB().WithContext(ctx).
		Where("folder = ? AND deleted_at IS NULL", folder).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&images)

	if err := r.handleError("GetByFolder", result.Error); err != nil {
		return nil, err
	}

	return images, nil
}
