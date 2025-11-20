package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// PermissionMiddleware creates middleware that checks for specific permissions
func PermissionMiddleware(permissionService services.PermissionService, resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by AuthMiddleware)
		userID, ok := c.Locals("userID").(string)
		if !ok || userID == "" {
			logger.Warn("Gebruiker ID niet gevonden in context voor permission check", "resource", resource, "action", action)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
				"code":  "UNAUTHORIZED",
			})
		}

		// Check permission
		if !permissionService.HasPermission(c.Context(), userID, resource, action) {
			logger.Warn("Toegang geweigerd - onvoldoende rechten",
				"user_id", userID,
				"resource", resource,
				"action", action,
				"path", c.Path(),
				"method", c.Method())
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Geen toegang tot deze resource",
				"code":  "FORBIDDEN",
			})
		}

		logger.Debug("Toegang verleend",
			"user_id", userID,
			"resource", resource,
			"action", action,
			"path", c.Path())

		return c.Next()
	}
}

// AdminPermissionMiddleware is a convenience middleware for admin-only actions
func AdminPermissionMiddleware(permissionService services.PermissionService) fiber.Handler {
	return PermissionMiddleware(permissionService, "admin", "access")
}

// StaffPermissionMiddleware is a convenience middleware for staff-level actions
func StaffPermissionMiddleware(permissionService services.PermissionService) fiber.Handler {
	return PermissionMiddleware(permissionService, "staff", "access")
}

// AdminOrStaffPermissionMiddleware allows access for users with either admin OR staff permissions
func AdminOrStaffPermissionMiddleware(permissionService services.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by AuthMiddleware)
		userID, ok := c.Locals("userID").(string)
		if !ok || userID == "" {
			logger.Warn("Gebruiker ID niet gevonden in context voor admin/staff permission check")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
				"code":  "UNAUTHORIZED",
			})
		}

		// Check for admin OR staff access
		hasAdminAccess := permissionService.HasPermission(c.Context(), userID, "admin", "access")
		hasStaffAccess := permissionService.HasPermission(c.Context(), userID, "staff", "access")

		logger.Debug("AdminOrStaffPermissionMiddleware check",
			"user_id", userID,
			"has_admin_access", hasAdminAccess,
			"has_staff_access", hasStaffAccess,
			"path", c.Path())

		if !hasAdminAccess && !hasStaffAccess {
			logger.Warn("Toegang geweigerd - geen admin of staff rechten",
				"user_id", userID,
				"has_admin", hasAdminAccess,
				"has_staff", hasStaffAccess,
				"path", c.Path(),
				"method", c.Method())
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin of staff rechten vereist",
				"code":  "FORBIDDEN",
			})
		}

		logger.Debug("Admin/staff toegang verleend",
			"user_id", userID,
			"has_admin", hasAdminAccess,
			"has_staff", hasStaffAccess,
			"path", c.Path())

		return c.Next()
	}
}

// ResourcePermissionMiddleware creates middleware for specific resource permissions
// This allows checking multiple permissions for a single route
func ResourcePermissionMiddleware(permissionService services.PermissionService, permissions ...models.PermissionCheck) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context
		userID, ok := c.Locals("userID").(string)
		if !ok || userID == "" {
			logger.Warn("Gebruiker ID niet gevonden in context voor resource permission check")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
				"code":  "UNAUTHORIZED",
			})
		}

		// Check all required permissions
		for _, perm := range permissions {
			if !permissionService.HasPermission(c.Context(), userID, perm.Resource, perm.Action) {
				logger.Warn("Resource toegang geweigerd",
					"user_id", userID,
					"resource", perm.Resource,
					"action", perm.Action,
					"path", c.Path(),
					"method", c.Method())
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Geen toegang tot deze resources",
					"code":  "FORBIDDEN",
				})
			}
		}

		logger.Debug("Alle resource permissions verleend",
			"user_id", userID,
			"permissions_count", len(permissions),
			"path", c.Path())

		return c.Next()
	}
}
