package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresSocialLinkRepository implements SocialLinkRepository
type PostgresSocialLinkRepository struct {
	db *gorm.DB
}

// NewPostgresSocialLinkRepository creates a new social link repository
func NewPostgresSocialLinkRepository(db *gorm.DB) *PostgresSocialLinkRepository {
	return &PostgresSocialLinkRepository{db: db}
}

// Create saves a new social link
func (r *PostgresSocialLinkRepository) Create(ctx context.Context, link *models.SocialLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

// GetByID retrieves a social link by ID
func (r *PostgresSocialLinkRepository) GetByID(ctx context.Context, id string) (*models.SocialLink, error) {
	var link models.SocialLink
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// List retrieves a paginated list of social links
func (r *PostgresSocialLinkRepository) List(ctx context.Context, limit, offset int) ([]*models.SocialLink, error) {
	var links []*models.SocialLink
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("order_number ASC, created_at DESC").Find(&links).Error
	return links, err
}

// ListVisible retrieves only visible social links ordered by order_number
func (r *PostgresSocialLinkRepository) ListVisible(ctx context.Context) ([]*models.SocialLink, error) {
	var links []*models.SocialLink
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("order_number ASC, created_at DESC").Find(&links).Error
	return links, err
}

// Update updates an existing social link
func (r *PostgresSocialLinkRepository) Update(ctx context.Context, link *models.SocialLink) error {
	return r.db.WithContext(ctx).Save(link).Error
}

// Delete removes a social link
func (r *PostgresSocialLinkRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.SocialLink{}, "id = ?", id).Error
}
