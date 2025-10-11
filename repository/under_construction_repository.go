package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresUnderConstructionRepository implements UnderConstructionRepository
type PostgresUnderConstructionRepository struct {
	db *gorm.DB
}

// NewPostgresUnderConstructionRepository creates a new under construction repository
func NewPostgresUnderConstructionRepository(db *gorm.DB) *PostgresUnderConstructionRepository {
	return &PostgresUnderConstructionRepository{db: db}
}

// Create saves a new under construction record
func (r *PostgresUnderConstructionRepository) Create(ctx context.Context, uc *models.UnderConstruction) error {
	return r.db.WithContext(ctx).Create(uc).Error
}

// GetByID retrieves an under construction record by ID
func (r *PostgresUnderConstructionRepository) GetByID(ctx context.Context, id int) (*models.UnderConstruction, error) {
	var uc models.UnderConstruction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&uc).Error
	if err != nil {
		return nil, err
	}
	return &uc, nil
}

// List retrieves a paginated list of under construction records
func (r *PostgresUnderConstructionRepository) List(ctx context.Context, limit, offset int) ([]*models.UnderConstruction, error) {
	var ucs []*models.UnderConstruction
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&ucs).Error
	return ucs, err
}

// GetActive retrieves the active under construction record
func (r *PostgresUnderConstructionRepository) GetActive(ctx context.Context) (*models.UnderConstruction, error) {
	var uc models.UnderConstruction
	err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&uc).Error
	if err != nil {
		return nil, err
	}
	return &uc, nil
}

// Update updates an existing under construction record
func (r *PostgresUnderConstructionRepository) Update(ctx context.Context, uc *models.UnderConstruction) error {
	return r.db.WithContext(ctx).Save(uc).Error
}

// Delete removes an under construction record
func (r *PostgresUnderConstructionRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.UnderConstruction{}, "id = ?", id).Error
}
