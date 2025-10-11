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

// ListVisibleFiltered retrieves visible photos with filtering
func (r *PostgresPhotoRepository) ListVisibleFiltered(ctx context.Context, filters map[string]interface{}) ([]*models.Photo, error) {
	var photos []*models.Photo
	query := r.db.WithContext(ctx).Where("visible = ?", true)

	// Apply filters
	if year, ok := filters["year"].(int); ok && year > 0 {
		query = query.Where("year = ?", year)
	}

	if title, ok := filters["title"].(string); ok && title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if description, ok := filters["description"].(string); ok && description != "" {
		query = query.Where("description ILIKE ?", "%"+description+"%")
	}

	if cloudinaryFolder, ok := filters["cloudinary_folder"].(string); ok && cloudinaryFolder != "" {
		query = query.Where("cloudinary_folder = ?", cloudinaryFolder)
	}

	err := query.Order("created_at DESC").Find(&photos).Error
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

// ListByAlbumIDWithInfo retrieves photos for a specific album with relationship info
func (r *PostgresPhotoRepository) ListByAlbumIDWithInfo(ctx context.Context, albumID string) ([]*models.PhotoWithAlbumInfo, error) {
	var photos []*models.PhotoWithAlbumInfo
	err := r.db.WithContext(ctx).
		Table("photos").
		Select("photos.*, album_photos.album_id, album_photos.order_number").
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
