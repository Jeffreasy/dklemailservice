package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresPartnerRepository implements PartnerRepository
type PostgresPartnerRepository struct {
	db *gorm.DB
}

// NewPostgresPartnerRepository creates a new partner repository
func NewPostgresPartnerRepository(db *gorm.DB) *PostgresPartnerRepository {
	return &PostgresPartnerRepository{db: db}
}

// Create saves a new partner
func (r *PostgresPartnerRepository) Create(ctx context.Context, partner *models.Partner) error {
	return r.db.WithContext(ctx).Create(partner).Error
}

// GetByID retrieves a partner by ID
func (r *PostgresPartnerRepository) GetByID(ctx context.Context, id string) (*models.Partner, error) {
	var partner models.Partner
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&partner).Error
	if err != nil {
		return nil, err
	}
	return &partner, nil
}

// List retrieves a paginated list of partners
func (r *PostgresPartnerRepository) List(ctx context.Context, limit, offset int) ([]*models.Partner, error) {
	var partners []*models.Partner
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("order_number ASC, created_at DESC").Find(&partners).Error
	return partners, err
}

// ListVisible retrieves only visible partners ordered by order_number
func (r *PostgresPartnerRepository) ListVisible(ctx context.Context) ([]*models.Partner, error) {
	var partners []*models.Partner
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("order_number ASC, created_at DESC").Find(&partners).Error
	return partners, err
}

// Update updates an existing partner
func (r *PostgresPartnerRepository) Update(ctx context.Context, partner *models.Partner) error {
	return r.db.WithContext(ctx).Save(partner).Error
}

// Delete removes a partner
func (r *PostgresPartnerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Partner{}, "id = ?", id).Error
}
