package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresPhotoRepository implements PhotoRepository
type PostgresPhotoRepository struct {
	db *gorm.DB
}

// NewPostgresPhotoRepository creates a new photo repository
func NewPostgresPhotoRepository(db *gorm.DB) *PostgresPhotoRepository {
	return &PostgresPhotoRepository{db: db}
}

// Create saves a new photo
func (r *PostgresPhotoRepository) Create(ctx context.Context, photo *models.Photo) error {
	return r.db.WithContext(ctx).Create(photo).Error
}

// GetByID retrieves a photo by ID
func (r *PostgresPhotoRepository) GetByID(ctx context.Context, id string) (*models.Photo, error) {
	var photo models.Photo
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&photo).Error
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

// List retrieves a paginated list of photos
func (r *PostgresPhotoRepository) List(ctx context.Context, limit, offset int) ([]*models.Photo, error) {
	var photos []*models.Photo
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&photos).Error
	return photos, err
}

// ListVisible retrieves only visible photos
func (r *PostgresPhotoRepository) ListVisible(ctx context.Context) ([]*models.Photo, error) {
	var photos []*models.Photo
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("created_at DESC").Find(&photos).Error
	return photos, err
}

// ListByAlbumID retrieves photos for a specific album
func (r *PostgresPhotoRepository) ListByAlbumID(ctx context.Context, albumID string) ([]*models.Photo, error) {
	var photos []*models.Photo
	err := r.db.WithContext(ctx).
		Joins("JOIN album_photos ON photos.id = album_photos.photo_id").
		Where("album_photos.album_id = ? AND photos.visible = ?", albumID, true).
		Order("album_photos.order_number ASC, photos.created_at DESC").
		Find(&photos).Error
	return photos, err
}

// Update updates an existing photo
func (r *PostgresPhotoRepository) Update(ctx context.Context, photo *models.Photo) error {
	return r.db.WithContext(ctx).Save(photo).Error
}

// Delete removes a photo
func (r *PostgresPhotoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Photo{}, "id = ?", id).Error
}
