package models

import (
	"time"
)

// RBACRole represents a role in the RBAC system
type RBACRole struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	IsSystemRole bool      `gorm:"default:false" json:"is_system_role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy    *string   `gorm:"type:uuid" json:"created_by,omitempty"`

	// Relations
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []Gebruiker  `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

func (RBACRole) TableName() string {
	return "roles"
}

// Permission represents a granular permission in the RBAC system
type Permission struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Resource           string    `gorm:"type:varchar(100);not null" json:"resource"`
	Action             string    `gorm:"type:varchar(50);not null" json:"action"`
	Description        string    `gorm:"type:text" json:"description"`
	IsSystemPermission bool      `gorm:"default:false" json:"is_system_permission"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Roles []RBACRole `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RoleID       string    `gorm:"type:uuid;not null" json:"role_id"`
	PermissionID string    `gorm:"type:uuid;not null" json:"permission_id"`
	AssignedAt   time.Time `gorm:"autoCreateTime" json:"assigned_at"`
	AssignedBy   *string   `gorm:"type:uuid" json:"assigned_by,omitempty"`

	// Relations
	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID         string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     string     `gorm:"type:uuid;not null" json:"user_id"`
	RoleID     string     `gorm:"type:uuid;not null" json:"role_id"`
	AssignedAt time.Time  `gorm:"autoCreateTime" json:"assigned_at"`
	AssignedBy *string    `gorm:"type:uuid" json:"assigned_by,omitempty"`
	ExpiresAt  *time.Time `gorm:"type:timestamp" json:"expires_at,omitempty"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`

	// Relations
	User Gebruiker `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role RBACRole  `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// UserPermission represents a flattened view of user permissions (used by the view)
type UserPermission struct {
	UserID               string    `json:"user_id"`
	Email                string    `json:"email"`
	RoleName             string    `json:"role_name"`
	Resource             string    `json:"resource"`
	Action               string    `json:"action"`
	PermissionAssignedAt time.Time `json:"permission_assigned_at"`
	RoleAssignedAt       time.Time `json:"role_assigned_at"`
}

// PermissionCheck represents a permission check request
type PermissionCheck struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// RoleAssignment represents a role assignment request
type RoleAssignment struct {
	UserID     string     `json:"user_id"`
	RoleID     string     `json:"role_id"`
	AssignedBy *string    `json:"assigned_by,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// PermissionAssignment represents a permission assignment to a role
type PermissionAssignment struct {
	RoleID       string  `json:"role_id"`
	PermissionID string  `json:"permission_id"`
	AssignedBy   *string `json:"assigned_by,omitempty"`
}
