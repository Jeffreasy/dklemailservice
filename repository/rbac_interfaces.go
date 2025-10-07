package repository

import (
	"context"
	"dklautomationgo/models"
)

// RBACRoleRepository defines the interface for RBAC role operations
type RBACRoleRepository interface {
	// Create creates a new role
	Create(ctx context.Context, role *models.RBACRole) error

	// GetByID retrieves a role by ID
	GetByID(ctx context.Context, id string) (*models.RBACRole, error)

	// GetByName retrieves a role by name
	GetByName(ctx context.Context, name string) (*models.RBACRole, error)

	// List retrieves all roles with pagination
	List(ctx context.Context, limit, offset int) ([]*models.RBACRole, error)

	// Update updates an existing role
	Update(ctx context.Context, role *models.RBACRole) error

	// Delete deletes a role
	Delete(ctx context.Context, id string) error

	// GetSystemRoles retrieves all system roles
	GetSystemRoles(ctx context.Context) ([]*models.RBACRole, error)
}

// PermissionRepository defines the interface for permission operations
type PermissionRepository interface {
	// Create creates a new permission
	Create(ctx context.Context, permission *models.Permission) error

	// GetByID retrieves a permission by ID
	GetByID(ctx context.Context, id string) (*models.Permission, error)

	// GetByResourceAction retrieves a permission by resource and action
	GetByResourceAction(ctx context.Context, resource, action string) (*models.Permission, error)

	// List retrieves all permissions with pagination
	List(ctx context.Context, limit, offset int) ([]*models.Permission, error)

	// ListByResource retrieves permissions by resource
	ListByResource(ctx context.Context, resource string) ([]*models.Permission, error)

	// Update updates an existing permission
	Update(ctx context.Context, permission *models.Permission) error

	// Delete deletes a permission
	Delete(ctx context.Context, id string) error

	// GetSystemPermissions retrieves all system permissions
	GetSystemPermissions(ctx context.Context) ([]*models.Permission, error)
}

// RolePermissionRepository defines the interface for role-permission relationship operations
type RolePermissionRepository interface {
	// Create creates a role-permission relationship
	Create(ctx context.Context, rp *models.RolePermission) error

	// Delete removes a role-permission relationship
	Delete(ctx context.Context, roleID, permissionID string) error

	// DeleteByRoleID removes all permissions for a role
	DeleteByRoleID(ctx context.Context, roleID string) error

	// DeleteByPermissionID removes a permission from all roles
	DeleteByPermissionID(ctx context.Context, permissionID string) error

	// GetPermissionsByRole retrieves all permissions for a role
	GetPermissionsByRole(ctx context.Context, roleID string) ([]*models.Permission, error)

	// GetRolesByPermission retrieves all roles that have a permission
	GetRolesByPermission(ctx context.Context, permissionID string) ([]*models.RBACRole, error)

	// HasPermission checks if a role has a specific permission
	HasPermission(ctx context.Context, roleID, permissionID string) (bool, error)
}

// UserRoleRepository defines the interface for user-role relationship operations
type UserRoleRepository interface {
	// Create creates a user-role relationship
	Create(ctx context.Context, ur *models.UserRole) error

	// GetByID retrieves a user-role relationship by ID
	GetByID(ctx context.Context, id string) (*models.UserRole, error)

	// GetByUserAndRole retrieves a user-role relationship by user and role IDs
	GetByUserAndRole(ctx context.Context, userID, roleID string) (*models.UserRole, error)

	// ListByUser retrieves all roles for a user
	ListByUser(ctx context.Context, userID string) ([]*models.UserRole, error)

	// ListByRole retrieves all users for a role
	ListByRole(ctx context.Context, roleID string) ([]*models.UserRole, error)

	// ListActiveByUser retrieves all active roles for a user
	ListActiveByUser(ctx context.Context, userID string) ([]*models.UserRole, error)

	// Update updates a user-role relationship
	Update(ctx context.Context, ur *models.UserRole) error

	// Delete removes a user-role relationship
	Delete(ctx context.Context, id string) error

	// DeleteByUser removes all roles for a user
	DeleteByUser(ctx context.Context, userID string) error

	// DeleteByRole removes a role from all users
	DeleteByRole(ctx context.Context, roleID string) error

	// Deactivate deactivates a user-role relationship
	Deactivate(ctx context.Context, id string) error

	// GetUserPermissions retrieves all permissions for a user
	GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error)
}
