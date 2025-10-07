package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// RolePermissionRepositoryImpl implements RolePermissionRepository
type RolePermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewRolePermissionRepository creates a new RolePermissionRepository
func NewRolePermissionRepository(db *gorm.DB) RolePermissionRepository {
	return &RolePermissionRepositoryImpl{db: db}
}

// Create creates a role-permission relationship
func (r *RolePermissionRepositoryImpl) Create(ctx context.Context, rp *models.RolePermission) error {
	return r.db.WithContext(ctx).Create(rp).Error
}

// Delete removes a role-permission relationship
func (r *RolePermissionRepositoryImpl) Delete(ctx context.Context, roleID, permissionID string) error {
	return r.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&models.RolePermission{}).Error
}

// DeleteByRoleID removes all permissions for a role
func (r *RolePermissionRepositoryImpl) DeleteByRoleID(ctx context.Context, roleID string) error {
	return r.db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error
}

// DeleteByPermissionID removes a permission from all roles
func (r *RolePermissionRepositoryImpl) DeleteByPermissionID(ctx context.Context, permissionID string) error {
	return r.db.WithContext(ctx).Where("permission_id = ?", permissionID).Delete(&models.RolePermission{}).Error
}

// GetPermissionsByRole retrieves all permissions for a role
func (r *RolePermissionRepositoryImpl) GetPermissionsByRole(ctx context.Context, roleID string) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Order("permissions.resource ASC, permissions.action ASC").
		Find(&permissions).Error
	return permissions, err
}

// GetRolesByPermission retrieves all roles that have a permission
func (r *RolePermissionRepositoryImpl) GetRolesByPermission(ctx context.Context, permissionID string) ([]*models.RBACRole, error) {
	var roles []*models.RBACRole
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON roles.id = rp.role_id").
		Where("rp.permission_id = ?", permissionID).
		Order("roles.name ASC").
		Find(&roles).Error
	return roles, err
}

// HasPermission checks if a role has a specific permission
func (r *RolePermissionRepositoryImpl) HasPermission(ctx context.Context, roleID, permissionID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error
	return count > 0, err
}
