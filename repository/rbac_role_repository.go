package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// RBACRoleRepositoryImpl implements RBACRoleRepository
type RBACRoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRBACRoleRepository creates a new RBACRoleRepository
func NewRBACRoleRepository(db *gorm.DB) RBACRoleRepository {
	return &RBACRoleRepositoryImpl{db: db}
}

// Create creates a new role
func (r *RBACRoleRepositoryImpl) Create(ctx context.Context, role *models.RBACRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// GetByID retrieves a role by ID
func (r *RBACRoleRepositoryImpl) GetByID(ctx context.Context, id string) (*models.RBACRole, error) {
	var role models.RBACRole
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName retrieves a role by name
func (r *RBACRoleRepositoryImpl) GetByName(ctx context.Context, name string) (*models.RBACRole, error) {
	var role models.RBACRole
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List retrieves all roles with pagination
func (r *RBACRoleRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
	var roles []*models.RBACRole
	query := r.db.WithContext(ctx).Order("name ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&roles).Error
	return roles, err
}

// Update updates an existing role
func (r *RBACRoleRepositoryImpl) Update(ctx context.Context, role *models.RBACRole) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete deletes a role
func (r *RBACRoleRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.RBACRole{}, "id = ?", id).Error
}

// GetSystemRoles retrieves all system roles
func (r *RBACRoleRepositoryImpl) GetSystemRoles(ctx context.Context) ([]*models.RBACRole, error) {
	var roles []*models.RBACRole
	err := r.db.WithContext(ctx).Where("is_system_role = ?", true).Order("name ASC").Find(&roles).Error
	return roles, err
}
