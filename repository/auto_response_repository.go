package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresAutoResponseRepository implements AutoResponseRepository using PostgreSQL
type PostgresAutoResponseRepository struct {
	db *gorm.DB
}

// NewPostgresAutoResponseRepository creates a new auto response repository
func NewPostgresAutoResponseRepository(db *gorm.DB) *PostgresAutoResponseRepository {
	return &PostgresAutoResponseRepository{db: db}
}

// GetAll retrieves all auto responses
func (r *PostgresAutoResponseRepository) GetAll(ctx context.Context) ([]*models.AutoResponse, error) {
	var responses []*models.AutoResponse
	if err := r.db.WithContext(ctx).Find(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

// GetByID retrieves an auto response by ID
func (r *PostgresAutoResponseRepository) GetByID(ctx context.Context, id int) (*models.AutoResponse, error) {
	var response models.AutoResponse
	if err := r.db.WithContext(ctx).First(&response, id).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

// Create saves a new auto response
func (r *PostgresAutoResponseRepository) Create(ctx context.Context, response *models.AutoResponse) error {
	return r.db.WithContext(ctx).Create(response).Error
}

// Update updates an existing auto response
func (r *PostgresAutoResponseRepository) Update(ctx context.Context, response *models.AutoResponse) error {
	return r.db.WithContext(ctx).Save(response).Error
}

// Delete removes an auto response
func (r *PostgresAutoResponseRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.AutoResponse{}, id).Error
}
