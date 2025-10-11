package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresAlbumRepository implements AlbumRepository
type PostgresAlbumRepository struct {
	db *gorm.DB
}

// NewPostgresAlbumRepository creates a new album repository
func NewPostgresAlbumRepository(db *gorm.DB) *PostgresAlbumRepository {
	return &PostgresAlbumRepository{db: db}
}

// Create saves a new album
func (r *PostgresAlbumRepository) Create(ctx context.Context, album *models.Album) error {
	return r.db.WithContext(ctx).Create(album).Error
}

// GetByID retrieves an album by ID
func (r *PostgresAlbumRepository) GetByID(ctx context.Context, id string) (*models.Album, error) {
	var album models.Album
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&album).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

// List retrieves a paginated list of albums
func (r *PostgresAlbumRepository) List(ctx context.Context, limit, offset int) ([]*models.Album, error) {
	var albums []*models.Album
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("order_number ASC, created_at DESC").Find(&albums).Error
	return albums, err
}

// ListVisible retrieves only visible albums ordered by order_number
func (r *PostgresAlbumRepository) ListVisible(ctx context.Context) ([]*models.Album, error) {
	var albums []*models.Album
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("order_number ASC, created_at DESC").Find(&albums).Error
	return albums, err
}

// ListVisibleWithCovers retrieves visible albums with cover photo information
func (r *PostgresAlbumRepository) ListVisibleWithCovers(ctx context.Context) ([]*models.AlbumWithCover, error) {
	var albums []*models.AlbumWithCover
	err := r.db.WithContext(ctx).
		Where("albums.visible = ?", true).
		Preload("CoverPhoto").
		Order("albums.order_number ASC, albums.created_at DESC").
		Find(&albums).Error
	return albums, err
}

// Update updates an existing album
func (r *PostgresAlbumRepository) Update(ctx context.Context, album *models.Album) error {
	return r.db.WithContext(ctx).Save(album).Error
}

// Delete removes an album
func (r *PostgresAlbumRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Album{}, "id = ?", id).Error
}
