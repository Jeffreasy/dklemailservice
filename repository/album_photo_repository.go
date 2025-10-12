package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresAlbumPhotoRepository implements AlbumPhotoRepository
type PostgresAlbumPhotoRepository struct {
	db *gorm.DB
}

// NewPostgresAlbumPhotoRepository creates a new album photo repository
func NewPostgresAlbumPhotoRepository(db *gorm.DB) *PostgresAlbumPhotoRepository {
	return &PostgresAlbumPhotoRepository{db: db}
}

// Create adds a photo to an album
func (r *PostgresAlbumPhotoRepository) Create(ctx context.Context, albumPhoto *models.AlbumPhoto) error {
	return r.db.WithContext(ctx).Create(albumPhoto).Error
}

// Delete removes a photo from an album
func (r *PostgresAlbumPhotoRepository) Delete(ctx context.Context, albumID, photoID string) error {
	return r.db.WithContext(ctx).Where("album_id = ? AND photo_id = ?", albumID, photoID).Delete(&models.AlbumPhoto{}).Error
}

// GetByAlbumAndPhoto retrieves a specific album-photo relationship
func (r *PostgresAlbumPhotoRepository) GetByAlbumAndPhoto(ctx context.Context, albumID, photoID string) (*models.AlbumPhoto, error) {
	var albumPhoto models.AlbumPhoto
	err := r.db.WithContext(ctx).Where("album_id = ? AND photo_id = ?", albumID, photoID).First(&albumPhoto).Error
	if err != nil {
		return nil, err
	}
	return &albumPhoto, nil
}

// ListByAlbum retrieves all photos for an album ordered by order_number
func (r *PostgresAlbumPhotoRepository) ListByAlbum(ctx context.Context, albumID string) ([]*models.AlbumPhoto, error) {
	var albumPhotos []*models.AlbumPhoto
	err := r.db.WithContext(ctx).Where("album_id = ?", albumID).Order("order_number ASC").Find(&albumPhotos).Error
	return albumPhotos, err
}

// UpdateOrder updates the order number of a photo in an album
func (r *PostgresAlbumPhotoRepository) UpdateOrder(ctx context.Context, albumID, photoID string, orderNumber int) error {
	return r.db.WithContext(ctx).Model(&models.AlbumPhoto{}).
		Where("album_id = ? AND photo_id = ?", albumID, photoID).
		Update("order_number", orderNumber).Error
}

// DeleteByAlbum removes all photos from an album
func (r *PostgresAlbumPhotoRepository) DeleteByAlbum(ctx context.Context, albumID string) error {
	return r.db.WithContext(ctx).Where("album_id = ?", albumID).Delete(&models.AlbumPhoto{}).Error
}

// DeleteByPhoto removes a photo from all albums
func (r *PostgresAlbumPhotoRepository) DeleteByPhoto(ctx context.Context, photoID string) error {
	return r.db.WithContext(ctx).Where("photo_id = ?", photoID).Delete(&models.AlbumPhoto{}).Error
}
