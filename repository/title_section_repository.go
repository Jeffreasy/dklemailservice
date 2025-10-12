package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresTitleSectionRepository implements TitleSectionRepository
type PostgresTitleSectionRepository struct {
	db *gorm.DB
}

// NewPostgresTitleSectionRepository creates a new title section repository
func NewPostgresTitleSectionRepository(db *gorm.DB) *PostgresTitleSectionRepository {
	return &PostgresTitleSectionRepository{db: db}
}

// Get retrieves the title section content (assuming there's only one record)
func (r *PostgresTitleSectionRepository) Get(ctx context.Context) (*models.TitleSection, error) {
	var titleSection models.TitleSection
	err := r.db.WithContext(ctx).First(&titleSection).Error
	if err != nil {
		return nil, err
	}
	return &titleSection, nil
}

// Create saves a new title section
func (r *PostgresTitleSectionRepository) Create(ctx context.Context, titleSection *models.TitleSection) error {
	return r.db.WithContext(ctx).Create(titleSection).Error
}

// Update updates an existing title section
func (r *PostgresTitleSectionRepository) Update(ctx context.Context, titleSection *models.TitleSection) error {
	return r.db.WithContext(ctx).Save(titleSection).Error
}

// Delete removes a title section
func (r *PostgresTitleSectionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.TitleSection{}, "id = ?", id).Error
}
