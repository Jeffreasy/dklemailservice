package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// distanceRepository implementeert DistanceRepository
type distanceRepository struct {
	db *gorm.DB
}

// NewDistanceRepository maakt een nieuwe distance repository
func NewDistanceRepository(db *gorm.DB) DistanceRepository {
	return &distanceRepository{db: db}
}

// Create slaat een nieuwe distance op
func (r *distanceRepository) Create(ctx context.Context, distance *models.Distance) error {
	return r.db.WithContext(ctx).Create(distance).Error
}

// GetByRoute haalt een distance op basis van route naam
func (r *distanceRepository) GetByRoute(ctx context.Context, route string) (*models.Distance, error) {
	var distance models.Distance
	err := r.db.WithContext(ctx).Where("route = ?", route).First(&distance).Error
	if err != nil {
		return nil, err
	}
	return &distance, nil
}

// GetAll haalt alle distances op
func (r *distanceRepository) GetAll(ctx context.Context) ([]*models.Distance, error) {
	var distances []*models.Distance
	err := r.db.WithContext(ctx).Order("route ASC").Find(&distances).Error
	return distances, err
}

// List haalt alle distances op (alias voor GetAll voor backward compatibility)
func (r *distanceRepository) List(ctx context.Context) ([]*models.Distance, error) {
	return r.GetAll(ctx)
}

// Update werkt een distance bij
func (r *distanceRepository) Update(ctx context.Context, distance *models.Distance) error {
	return r.db.WithContext(ctx).Save(distance).Error
}

// Delete verwijdert een distance
func (r *distanceRepository) Delete(ctx context.Context, route string) error {
	return r.db.WithContext(ctx).Where("route = ?", route).Delete(&models.Distance{}).Error
}
