package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	rbacGroup.Put("/roles/:id/permissions", h.UpdateRolePermissions)                     // Voor bulk updates (frontend compatibiliteit)
	rbacGroup.Post("/roles/:id/permissions/:permissionId", h.AddPermissionToRole)        // Voor individuele toevoeging
	rbacGroup.Delete("/roles/:id/permissions/:permissionId", h.RemovePermissionFromRole) // Voor individuele verwijdering
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
	if err != nil && err != gorm.ErrRecordNotFound {
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
	if err != nil && err != gorm.ErrRecordNotFound {
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

// AddPermissionToRole voegt één permission toe aan een role
func (h *PermissionHandler) AddPermissionToRole(c *fiber.Ctx) error {
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

// UpdateRolePermissions werkt permissions bij voor een role (bulk update voor frontend compatibiliteit)
func (h *PermissionHandler) UpdateRolePermissions(c *fiber.Ctx) error {
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

	// Haal huidige permissions voor deze role op
	currentPermissions, err := h.rolePermissionRepo.GetPermissionsByRole(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen huidige permissions", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon huidige permissions niet ophalen",
		})
	}

	// Maak sets voor vergelijking
	currentPermissionIDs := make(map[string]bool)
	for _, perm := range currentPermissions {
		currentPermissionIDs[perm.ID] = true
	}

	requestedPermissionIDs := make(map[string]bool)
	for _, id := range req.PermissionIDs {
		requestedPermissionIDs[id] = true
	}

	// Bepaal welke permissions toegevoegd/verwijderd moeten worden
	var toAdd []string
	var toRemove []string

	// Permissions die toegevoegd moeten worden (in request maar niet current)
	for _, permID := range req.PermissionIDs {
		if !currentPermissionIDs[permID] {
			toAdd = append(toAdd, permID)
		}
	}

	// Permissions die verwijderd moeten worden (current maar niet in request)
	for _, perm := range currentPermissions {
		if !requestedPermissionIDs[perm.ID] {
			toRemove = append(toRemove, perm.ID)
		}
	}

	// Voer toevoegingen uit
	addedCount := 0
	for _, permissionID := range toAdd {
		rp := &models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
			AssignedBy:   &userID,
		}

		if err := h.rolePermissionRepo.Create(ctx, rp); err != nil {
			logger.Error("Fout bij toevoegen permission", "error", err, "role_id", roleID, "permission_id", permissionID)
			continue
		}
		addedCount++
	}

	// Voer verwijderingen uit
	removedCount := 0
	for _, permissionID := range toRemove {
		if err := h.rolePermissionRepo.Delete(ctx, roleID, permissionID); err != nil {
			logger.Error("Fout bij verwijderen permission", "error", err, "role_id", roleID, "permission_id", permissionID)
			continue
		}
		removedCount++
	}

	return c.JSON(fiber.Map{
		"success":         true,
		"message":         "Role permissions bijgewerkt",
		"added_count":     addedCount,
		"removed_count":   removedCount,
		"total_requested": len(req.PermissionIDs),
	})
}
