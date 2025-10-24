package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// RouteFundRepository interface voor route fund operaties
type RouteFundRepository interface {
	Create(ctx context.Context, routeFund *models.RouteFund) error
	GetByRoute(ctx context.Context, route string) (*models.RouteFund, error)
	GetAll(ctx context.Context) ([]*models.RouteFund, error)
	Update(ctx context.Context, routeFund *models.RouteFund) error
	Delete(ctx context.Context, route string) error
}

// routeFundRepository implementeert RouteFundRepository
type routeFundRepository struct {
	db *gorm.DB
}

// NewRouteFundRepository maakt een nieuwe route fund repository
func NewRouteFundRepository(db *gorm.DB) RouteFundRepository {
	return &routeFundRepository{db: db}
}

// Create slaat een nieuwe route fund op
func (r *routeFundRepository) Create(ctx context.Context, routeFund *models.RouteFund) error {
	return r.db.WithContext(ctx).Create(routeFund).Error
}

// GetByRoute haalt een route fund op basis van route naam
func (r *routeFundRepository) GetByRoute(ctx context.Context, route string) (*models.RouteFund, error) {
	var routeFund models.RouteFund
	err := r.db.WithContext(ctx).Where("route = ?", route).First(&routeFund).Error
	if err != nil {
		return nil, err
	}
	return &routeFund, nil
}

// GetAll haalt alle route funds op
func (r *routeFundRepository) GetAll(ctx context.Context) ([]*models.RouteFund, error) {
	var routeFunds []*models.RouteFund
	err := r.db.WithContext(ctx).Order("route ASC").Find(&routeFunds).Error
	return routeFunds, err
}

// Update werkt een route fund bij
func (r *routeFundRepository) Update(ctx context.Context, routeFund *models.RouteFund) error {
	return r.db.WithContext(ctx).Save(routeFund).Error
}

// Delete verwijdert een route fund
func (r *routeFundRepository) Delete(ctx context.Context, route string) error {
	return r.db.WithContext(ctx).Where("route = ?", route).Delete(&models.RouteFund{}).Error
}
