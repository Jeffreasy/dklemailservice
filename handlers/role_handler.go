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

	// Role permission assignments - individueel beheer
	roleGroup.Post("/:id/permissions/:permissionId", h.AddPermissionToRole)
	roleGroup.Delete("/:id/permissions/:permissionId", h.RemovePermissionFromRole)
}

// AddPermissionToRole voegt één permission toe aan een role
func (h *RoleHandler) AddPermissionToRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	permissionID := c.Params("permissionId")

	if roleID == "" || permissionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID en Permission ID zijn verplicht",
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

	// Check if permission is already assigned
	hasPermission, err := h.rolePermissionRepo.HasPermission(ctx, roleID, permissionID)
	if err != nil {
		logger.Error("Fout bij controleren bestaande permission", "error", err, "role_id", roleID, "permission_id", permissionID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permission niet controleren",
		})
	}

	if hasPermission {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Permission is al toegewezen aan deze role",
		})
	}

	// Voeg permission toe
	rp := &models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		AssignedBy:   &userID,
	}

	if err := h.rolePermissionRepo.Create(ctx, rp); err != nil {
		logger.Error("Fout bij toewijzen permission aan role", "error", err, "role_id", roleID, "permission_id", permissionID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permission niet toewijzen aan role",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Permission toegevoegd aan role",
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
