package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresVideoRepository implements VideoRepository
type PostgresVideoRepository struct {
	db *gorm.DB
}

// NewPostgresVideoRepository creates a new video repository
func NewPostgresVideoRepository(db *gorm.DB) *PostgresVideoRepository {
	return &PostgresVideoRepository{db: db}
}

// Create saves a new video
func (r *PostgresVideoRepository) Create(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

// GetByID retrieves a video by ID
func (r *PostgresVideoRepository) GetByID(ctx context.Context, id string) (*models.Video, error) {
	var video models.Video
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// List retrieves a paginated list of videos
func (r *PostgresVideoRepository) List(ctx context.Context, limit, offset int) ([]*models.Video, error) {
	var videos []*models.Video
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("order_number ASC, created_at DESC").Find(&videos).Error
	return videos, err
}

// ListVisible retrieves only visible videos ordered by order_number
func (r *PostgresVideoRepository) ListVisible(ctx context.Context) ([]*models.Video, error) {
	var videos []*models.Video
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("order_number ASC, created_at DESC").Find(&videos).Error
	return videos, err
}

// Update updates an existing video
func (r *PostgresVideoRepository) Update(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Save(video).Error
}

// Delete removes a video
func (r *PostgresVideoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Video{}, "id = ?", id).Error
}
