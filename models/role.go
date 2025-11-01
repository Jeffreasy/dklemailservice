package models

// Legacy role constants - DEPRECATED
// Use RBAC system via user_roles table instead
// These constants are kept for backward compatibility only
type Role string

const (
	// System roles
	RoleAdmin Role = "admin"
	RoleStaff Role = "staff"
	RoleUser  Role = "user"

	// Chat roles - Fixed naming conflicts
	RoleChatOwner  Role = "owner"
	RoleChatAdmin  Role = "chat_admin" // Changed from "admin" to avoid conflict
	RoleChatMember Role = "member"

	// Event roles
	RoleDeelnemer    Role = "deelnemer"    // Lowercase for consistency
	RoleBegeleider   Role = "begeleider"   // Lowercase for consistency
	RoleVrijwilliger Role = "vrijwilliger" // Lowercase for consistency
)
