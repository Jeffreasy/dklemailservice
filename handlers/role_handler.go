package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// RoleHandler bevat handlers voor role-permission assignments
type RoleHandler struct {
	rolePermissionRepo repository.RolePermissionRepository
	authService        services.AuthService
	permissionService  services.PermissionService
}

// NewRoleHandler maakt een nieuwe role handler
func NewRoleHandler(
	rolePermissionRepo repository.RolePermissionRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *RoleHandler {
	return &RoleHandler{
		rolePermissionRepo: rolePermissionRepo,
		authService:        authService,
		permissionService:  permissionService,
	}
}

// RegisterRoutes registreert de routes voor role assignments
func (h *RoleHandler) RegisterRoutes(app *fiber.App) {
	// Groep voor role assignment routes (vereist admin rechten)
	roleGroup := app.Group("/api/roles")
	roleGroup.Use(AuthMiddleware(h.authService))
	roleGroup.Use(AdminPermissionMiddleware(h.permissionService))

	// Role permission assignments
	roleGroup.Put("/:id/permissions", h.AssignPermissionsToRole)
	roleGroup.Delete("/:id/permissions/:permissionId", h.RemovePermissionFromRole)
}

// AssignPermissionsToRole wijst permissions toe aan een role
func (h *RoleHandler) AssignPermissionsToRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID is verplicht",
		})
	}

	var req struct {
		PermissionIDs []string `json:"permission_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	if len(req.PermissionIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ten minste één permission ID is verplicht",
		})
	}

	// Haal userID op uit context
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon userID niet ophalen uit context",
		})
	}

	ctx := c.Context()
	assignedPermissions := 0

	for _, permissionID := range req.PermissionIDs {
		rp := &models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
			AssignedBy:   &userID,
		}

		if err := h.rolePermissionRepo.Create(ctx, rp); err != nil {
			logger.Error("Fout bij toewijzen permission aan role", "error", err, "role_id", roleID, "permission_id", permissionID)
			// Continue with other permissions
			continue
		}
		assignedPermissions++
	}

	return c.JSON(fiber.Map{
		"success":              true,
		"message":              "Permissions toegewezen aan role",
		"assigned_permissions": assignedPermissions,
		"total_requested":      len(req.PermissionIDs),
	})
}

// RemovePermissionFromRole verwijdert een permission van een role
func (h *RoleHandler) RemovePermissionFromRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	permissionID := c.Params("permissionId")

	if roleID == "" || permissionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID en Permission ID zijn verplicht",
		})
	}

	ctx := c.Context()
	if err := h.rolePermissionRepo.Delete(ctx, roleID, permissionID); err != nil {
		logger.Error("Fout bij verwijderen permission van role", "error", err, "role_id", roleID, "permission_id", permissionID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permission niet verwijderen van role",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Permission verwijderd van role",
	})
}
