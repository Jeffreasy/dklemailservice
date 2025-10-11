package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresSponsorRepository implements SponsorRepository
type PostgresSponsorRepository struct {
	db *gorm.DB
}

// NewPostgresSponsorRepository creates a new sponsor repository
func NewPostgresSponsorRepository(db *gorm.DB) *PostgresSponsorRepository {
	return &PostgresSponsorRepository{db: db}
}

// Create saves a new sponsor
func (r *PostgresSponsorRepository) Create(ctx context.Context, sponsor *models.Sponsor) error {
	return r.db.WithContext(ctx).Create(sponsor).Error
}

// GetByID retrieves a sponsor by ID
func (r *PostgresSponsorRepository) GetByID(ctx context.Context, id string) (*models.Sponsor, error) {
	var sponsor models.Sponsor
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sponsor).Error
	if err != nil {
		return nil, err
	}
	return &sponsor, nil
}

// List retrieves a paginated list of sponsors
func (r *PostgresSponsorRepository) List(ctx context.Context, limit, offset int) ([]*models.Sponsor, error) {
	var sponsors []*models.Sponsor
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("order_number ASC, created_at DESC").Find(&sponsors).Error
	return sponsors, err
}

// ListVisible retrieves only visible sponsors ordered by order_number
func (r *PostgresSponsorRepository) ListVisible(ctx context.Context) ([]*models.Sponsor, error) {
	var sponsors []*models.Sponsor
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("order_number ASC, created_at DESC").Find(&sponsors).Error
	return sponsors, err
}

// Update updates an existing sponsor
func (r *PostgresSponsorRepository) Update(ctx context.Context, sponsor *models.Sponsor) error {
	return r.db.WithContext(ctx).Save(sponsor).Error
}

// Delete removes a sponsor
func (r *PostgresSponsorRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Sponsor{}, "id = ?", id).Error
}
