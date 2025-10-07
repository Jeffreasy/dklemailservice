package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PermissionRepositoryImpl implements PermissionRepository
type PermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new PermissionRepository
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

// Create creates a new permission
func (r *PermissionRepositoryImpl) Create(ctx context.Context, permission *models.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// GetByID retrieves a permission by ID
func (r *PermissionRepositoryImpl) GetByID(ctx context.Context, id string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByResourceAction retrieves a permission by resource and action
func (r *PermissionRepositoryImpl) GetByResourceAction(ctx context.Context, resource, action string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.WithContext(ctx).Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// List retrieves all permissions with pagination
func (r *PermissionRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	var permissions []*models.Permission
	query := r.db.WithContext(ctx).Order("resource ASC, action ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&permissions).Error
	return permissions, err
}

// ListByResource retrieves permissions by resource
func (r *PermissionRepositoryImpl) ListByResource(ctx context.Context, resource string) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.WithContext(ctx).Where("resource = ?", resource).Order("action ASC").Find(&permissions).Error
	return permissions, err
}

// Update updates an existing permission
func (r *PermissionRepositoryImpl) Update(ctx context.Context, permission *models.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

// Delete deletes a permission
func (r *PermissionRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Permission{}, "id = ?", id).Error
}

// GetSystemPermissions retrieves all system permissions
func (r *PermissionRepositoryImpl) GetSystemPermissions(ctx context.Context) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.WithContext(ctx).Where("is_system_permission = ?", true).Order("resource ASC, action ASC").Find(&permissions).Error
	return permissions, err
}
