package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// PermissionHandler bevat handlers voor permission en role beheer
type PermissionHandler struct {
	permissionRepo     repository.PermissionRepository
	roleRepo           repository.RBACRoleRepository
	rolePermissionRepo repository.RolePermissionRepository
	authService        services.AuthService
	permissionService  services.PermissionService
}

// NewPermissionHandler maakt een nieuwe permission handler
func NewPermissionHandler(
	permissionRepo repository.PermissionRepository,
	roleRepo repository.RBACRoleRepository,
	rolePermissionRepo repository.RolePermissionRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *PermissionHandler {
	return &PermissionHandler{
		permissionRepo:     permissionRepo,
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
		authService:        authService,
		permissionService:  permissionService,
	}
}

// RegisterRoutes registreert de routes voor permission en role beheer
func (h *PermissionHandler) RegisterRoutes(app *fiber.App) {
	// RBAC routes (vereist admin rechten)
	rbacGroup := app.Group("/api/rbac")
	rbacGroup.Use(AuthMiddleware(h.authService))
	rbacGroup.Use(AdminPermissionMiddleware(h.permissionService))

	// Permission routes
	rbacGroup.Get("/permissions", h.ListPermissions)
	rbacGroup.Post("/permissions", h.CreatePermission)

	// Role routes
	rbacGroup.Get("/roles", h.ListRoles)
	rbacGroup.Post("/roles", h.CreateRole)
	rbacGroup.Put("/roles/:id/permissions", h.AssignPermissionsToRole)
	rbacGroup.Delete("/roles/:id/permissions/:permissionId", h.RemovePermissionFromRole)
}

// ListPermissions haalt een lijst van permissions op
func (h *PermissionHandler) ListPermissions(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	if limit < 1 || limit > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Limit moet tussen 1 en 100 liggen",
		})
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Offset mag niet negatief zijn",
		})
	}

	ctx := c.Context()
	permissions, err := h.permissionRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen permissions", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permissions niet ophalen",
		})
	}

	return c.JSON(permissions)
}

// CreatePermission maakt een nieuwe permission aan
func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	var req struct {
		Resource    string `json:"resource"`
		Action      string `json:"action"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	if req.Resource == "" || req.Action == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource en action zijn verplicht",
		})
	}

	// Check if permission already exists
	ctx := c.Context()
	existing, err := h.permissionRepo.GetByResourceAction(ctx, req.Resource, req.Action)
	if err != nil {
		logger.Error("Fout bij controleren bestaande permission", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permission niet controleren",
		})
	}

	if existing != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Permission bestaat al",
		})
	}

	permission := &models.Permission{
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	if err := h.permissionRepo.Create(ctx, permission); err != nil {
		logger.Error("Fout bij aanmaken permission", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permission niet aanmaken",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(permission)
}

// ListRoles haalt een lijst van roles op met hun permissions
func (h *PermissionHandler) ListRoles(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	if limit < 1 || limit > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Limit moet tussen 1 en 100 liggen",
		})
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Offset mag niet negatief zijn",
		})
	}

	ctx := c.Context()
	roles, err := h.roleRepo.ListWithPermissions(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen roles", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon roles niet ophalen",
		})
	}

	return c.JSON(roles)
}

// CreateRole maakt een nieuwe role aan
func (h *PermissionHandler) CreateRole(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is verplicht",
		})
	}

	// Check if role already exists
	ctx := c.Context()
	existing, err := h.roleRepo.GetByName(ctx, req.Name)
	if err != nil {
		logger.Error("Fout bij controleren bestaande role", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon role niet controleren",
		})
	}

	if existing != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Role bestaat al",
		})
	}

	role := &models.RBACRole{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.roleRepo.Create(ctx, role); err != nil {
		logger.Error("Fout bij aanmaken role", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon role niet aanmaken",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(role)
}

// AssignPermissionsToRole wijst permissions toe aan een role (vervangt alle bestaande permissions)
func (h *PermissionHandler) AssignPermissionsToRole(c *fiber.Ctx) error {
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

	// Haal userID op uit context
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon userID niet ophalen uit context",
		})
	}

	ctx := c.Context()

	// Start transaction - verwijder bestaande permissions en voeg nieuwe toe
	// Dit is een "replace all" operatie zoals in de JavaScript implementatie

	// Verwijder bestaande permissions voor deze role
	if err := h.rolePermissionRepo.DeleteByRoleID(ctx, roleID); err != nil {
		logger.Error("Fout bij verwijderen bestaande permissions", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon bestaande permissions niet verwijderen",
		})
	}

	// Voeg nieuwe permissions toe
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

	// Haal de bijgewerkte role op met permissions (zoals in JavaScript implementatie)
	updatedRole, err := h.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen bijgewerkte role", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon bijgewerkte role niet ophalen",
		})
	}

	// Haal permissions voor deze role op
	permissions, err := h.rolePermissionRepo.GetPermissionsByRole(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen permissions voor role", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permissions niet ophalen",
		})
	}

	// Format response zoals in JavaScript implementatie
	permissionObjects := make([]fiber.Map, len(permissions))
	for i, perm := range permissions {
		permissionObjects[i] = fiber.Map{
			"id":          perm.ID,
			"resource":    perm.Resource,
			"action":      perm.Action,
			"description": perm.Description,
		}
	}

	return c.JSON(fiber.Map{
		"id":          updatedRole.ID,
		"name":        updatedRole.Name,
		"description": updatedRole.Description,
		"created_at":  updatedRole.CreatedAt,
		"updated_at":  updatedRole.UpdatedAt,
		"permissions": permissionObjects,
	})
}

// RemovePermissionFromRole verwijdert een permission van een role
func (h *PermissionHandler) RemovePermissionFromRole(c *fiber.Ctx) error {
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
