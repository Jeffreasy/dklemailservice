package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// UserRoleRepositoryImpl implements UserRoleRepository
type UserRoleRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRoleRepository creates a new UserRoleRepository
func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &UserRoleRepositoryImpl{db: db}
}

// Create creates a user-role relationship
func (r *UserRoleRepositoryImpl) Create(ctx context.Context, ur *models.UserRole) error {
	return r.db.WithContext(ctx).Create(ur).Error
}

// GetByID retrieves a user-role relationship by ID
func (r *UserRoleRepositoryImpl) GetByID(ctx context.Context, id string) (*models.UserRole, error) {
	var ur models.UserRole
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ur).Error
	if err != nil {
		return nil, err
	}
	return &ur, nil
}

// GetByUserAndRole retrieves a user-role relationship by user and role IDs
func (r *UserRoleRepositoryImpl) GetByUserAndRole(ctx context.Context, userID, roleID string) (*models.UserRole, error) {
	var ur models.UserRole
	err := r.db.WithContext(ctx).Where("user_id = ? AND role_id = ? AND is_active = ?", userID, roleID, true).First(&ur).Error
	if err != nil {
		return nil, err
	}
	return &ur, nil
}

// ListByUser retrieves all roles for a user
func (r *UserRoleRepositoryImpl) ListByUser(ctx context.Context, userID string) ([]*models.UserRole, error) {
	var userRoles []*models.UserRole
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Role").Order("assigned_at DESC").Find(&userRoles).Error
	return userRoles, err
}

// ListByRole retrieves all users for a role
func (r *UserRoleRepositoryImpl) ListByRole(ctx context.Context, roleID string) ([]*models.UserRole, error) {
	var userRoles []*models.UserRole
	err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Preload("User").Order("assigned_at DESC").Find(&userRoles).Error
	return userRoles, err
}

// ListActiveByUser retrieves all active roles for a user
func (r *UserRoleRepositoryImpl) ListActiveByUser(ctx context.Context, userID string) ([]*models.UserRole, error) {
	var userRoles []*models.UserRole
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > NOW())", userID, true).
		Preload("Role").
		Order("assigned_at DESC").
		Find(&userRoles).Error
	return userRoles, err
}

// Update updates a user-role relationship
func (r *UserRoleRepositoryImpl) Update(ctx context.Context, ur *models.UserRole) error {
	return r.db.WithContext(ctx).Save(ur).Error
}

// Delete removes a user-role relationship
func (r *UserRoleRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.UserRole{}, "id = ?", id).Error
}

// DeleteByUser removes all roles for a user
func (r *UserRoleRepositoryImpl) DeleteByUser(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&models.UserRole{}, "user_id = ?", userID).Error
}

// DeleteByRole removes a role from all users
func (r *UserRoleRepositoryImpl) DeleteByRole(ctx context.Context, roleID string) error {
	return r.db.WithContext(ctx).Delete(&models.UserRole{}, "role_id = ?", roleID).Error
}

// Deactivate deactivates a user-role relationship
func (r *UserRoleRepositoryImpl) Deactivate(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.UserRole{}).Where("id = ?", id).Update("is_active", false).Error
}

// GetUserPermissions retrieves all permissions for a user through their roles
func (r *UserRoleRepositoryImpl) GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error) {
	var permissions []*models.UserPermission
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			ur.user_id,
			u.email,
			r.name as role_name,
			p.resource,
			p.action,
			rp.assigned_at as permission_assigned_at,
			ur.assigned_at as role_assigned_at
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		JOIN role_permissions rp ON r.id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		JOIN gebruikers u ON ur.user_id = u.id
		WHERE ur.user_id = ? AND ur.is_active = true
			AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		ORDER BY ur.user_id, r.name, p.resource, p.action
	`, userID).Scan(&permissions).Error
	return permissions, err
}
