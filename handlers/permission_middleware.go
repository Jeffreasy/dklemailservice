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
			logger.Warn("No user ID found in context for permission check", "resource", resource, "action", action)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
			})
		}

		// Check permission
		if !permissionService.HasPermission(c.Context(), userID, resource, action) {
			logger.Warn("Permission denied",
				"user_id", userID,
				"resource", resource,
				"action", action,
				"path", c.Path(),
				"method", c.Method())
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Geen toegang",
			})
		}

		logger.Debug("Permission granted",
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

// ResourcePermissionMiddleware creates middleware for specific resource permissions
// This allows checking multiple permissions for a single route
func ResourcePermissionMiddleware(permissionService services.PermissionService, permissions ...models.PermissionCheck) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context
		userID, ok := c.Locals("userID").(string)
		if !ok || userID == "" {
			logger.Warn("No user ID found in context for resource permission check")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
			})
		}

		// Check all required permissions
		for _, perm := range permissions {
			if !permissionService.HasPermission(c.Context(), userID, perm.Resource, perm.Action) {
				logger.Warn("Resource permission denied",
					"user_id", userID,
					"resource", perm.Resource,
					"action", perm.Action,
					"path", c.Path(),
					"method", c.Method())
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Geen toegang",
				})
			}
		}

		logger.Debug("All resource permissions granted",
			"user_id", userID,
			"permissions_count", len(permissions),
			"path", c.Path())

		return c.Next()
	}
}
