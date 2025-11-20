package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PermissionHandler bevat handlers voor permission en role beheer
type PermissionHandler struct {
	permissionRepo     repository.PermissionRepository
	roleRepo           repository.RBACRoleRepository
	rolePermissionRepo repository.RolePermissionRepository
	userRoleRepo       repository.UserRoleRepository
	authService        services.AuthService
	permissionService  services.PermissionService
}

// NewPermissionHandler maakt een nieuwe permission handler
func NewPermissionHandler(
	permissionRepo repository.PermissionRepository,
	roleRepo repository.RBACRoleRepository,
	rolePermissionRepo repository.RolePermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *PermissionHandler {
	return &PermissionHandler{
		permissionRepo:     permissionRepo,
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
		userRoleRepo:       userRoleRepo,
		authService:        authService,
		permissionService:  permissionService,
	}
}

// RegisterRoutes registreert de routes voor permission en role beheer
func (h *PermissionHandler) RegisterRoutes(app *fiber.App) {
	// RBAC routes (vereist admin of staff rechten)
	rbacGroup := app.Group("/api/rbac")
	rbacGroup.Use(AuthMiddleware(h.authService))
	rbacGroup.Use(AdminOrStaffPermissionMiddleware(h.permissionService))

	// Permission routes
	rbacGroup.Get("/permissions", h.ListPermissions)
	rbacGroup.Post("/permissions", h.CreatePermission)
	rbacGroup.Put("/permissions/:id", h.UpdatePermission)
	rbacGroup.Delete("/permissions/:id", h.DeletePermission)

	// Role routes
	rbacGroup.Get("/roles", h.ListRoles)
	rbacGroup.Post("/roles", h.CreateRole)

	// Menu management routes (specifiek voor admin interface) - MOETEN VOOR /roles/:id komen!
	rbacGroup.Get("/menu/permissions", h.GetMenuPermissions)                  // Alle menu permissies
	rbacGroup.Get("/menu/matrix", h.GetMenuPermissionMatrix)                  // Matrix view: rollen vs menu permissies
	rbacGroup.Get("/roles/:id/menu-permissions", h.GetRoleMenuPermissions)    // Menu permissies voor specifieke rol
	rbacGroup.Put("/roles/:id/menu-permissions", h.UpdateRoleMenuPermissions) // Update menu permissies voor rol

	// Role routes (generieke routes NA specifieke routes)
	rbacGroup.Get("/roles/:id", h.GetRole)                                               // Get role details
	rbacGroup.Put("/roles/:id", h.UpdateRole)                                            // Update role details
	rbacGroup.Delete("/roles/:id", h.DeleteRole)                                         // Delete role
	rbacGroup.Put("/roles/:id/permissions", h.UpdateRolePermissions)                     // Voor bulk updates (frontend compatibiliteit)
	rbacGroup.Post("/roles/:id/permissions/:permissionId", h.AddPermissionToRole)        // Voor individuele toevoeging
	rbacGroup.Delete("/roles/:id/permissions/:permissionId", h.RemovePermissionFromRole) // Voor individuele verwijdering

	// User-Role management routes
	rbacGroup.Get("/users", h.ListUsersWithRoles)                      // List users with their roles
	rbacGroup.Get("/users/:id/roles", h.GetUserRoles)                  // Get roles for specific user
	rbacGroup.Post("/users/:id/roles", h.AssignRoleToUser)             // Assign role to user
	rbacGroup.Delete("/users/:id/roles/:roleId", h.RevokeRoleFromUser) // Revoke role from user
	rbacGroup.Put("/users/:id/roles/:roleId", h.UpdateUserRole)        // Update user role (activate/deactivate)
}

// ListPermissions haalt een lijst van permissions op, gegroepeerd per resource
func (h *PermissionHandler) ListPermissions(c *fiber.Ctx) error {
	// Query parameters voor filtering
	resourceFilter := c.Query("resource", "")
	actionFilter := c.Query("action", "")
	search := c.Query("search", "")
	groupByResource := c.QueryBool("group_by_resource", true) // Default to true for grouped display

	// Pagination (alleen gebruikt als niet gegroepeerd)
	limit := c.QueryInt("limit", 1000) // Hogere limit voor gegroepeerde weergave
	offset := c.QueryInt("offset", 0)

	if limit < 1 || limit > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Limit moet tussen 1 en 1000 liggen",
		})
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Offset mag niet negatief zijn",
		})
	}

	ctx := c.Context()

	// Haal alle permissions op (we filteren client-side voor betere performance)
	permissions, err := h.permissionRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen permissions", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permissions niet ophalen",
		})
	}

	// Filter permissions gebaseerd op query parameters
	var filteredPermissions []*models.Permission
	for _, perm := range permissions {
		// Resource filter
		if resourceFilter != "" && perm.Resource != resourceFilter {
			continue
		}

		// Action filter
		if actionFilter != "" && perm.Action != actionFilter {
			continue
		}

		// Search filter (in resource, action of description)
		if search != "" {
			searchLower := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(perm.Resource), searchLower) &&
				!strings.Contains(strings.ToLower(perm.Action), searchLower) &&
				!strings.Contains(strings.ToLower(perm.Description), searchLower) {
				continue
			}
		}

		filteredPermissions = append(filteredPermissions, perm)
	}

	// Als gegroepeerd per resource gewenst
	if groupByResource {
		grouped := make(map[string][]*models.Permission)
		for _, perm := range filteredPermissions {
			grouped[perm.Resource] = append(grouped[perm.Resource], perm)
		}

		// Converteer naar frontend format
		result := make([]map[string]interface{}, 0, len(grouped))
		for resource, perms := range grouped {
			result = append(result, map[string]interface{}{
				"resource":    resource,
				"permissions": perms,
				"count":       len(perms),
			})
		}

		// Sorteer op resource naam
		sort.Slice(result, func(i, j int) bool {
			return result[i]["resource"].(string) < result[j]["resource"].(string)
		})

		return c.JSON(fiber.Map{
			"groups": result,
			"total":  len(filteredPermissions),
		})
	}

	// Backward compatible response - return array directly
	return c.JSON(filteredPermissions)
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

// UpdatePermission werkt een permission bij
func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	permissionID := c.Params("id")
	if permissionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Permission ID is verplicht",
		})
	}

	var req struct {
		Resource    *string `json:"resource,omitempty"`
		Action      *string `json:"action,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	ctx := c.Context()

	// Haal huidige permission op
	permission, err := h.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		logger.Error("Fout bij ophalen permission", "error", err, "permission_id", permissionID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Permission niet gevonden",
		})
	}

	// Controleer of het een systeempermission is
	if permission.IsSystemPermission {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Kan systeempermission niet bewerken",
		})
	}

	// Update velden indien opgegeven
	if req.Resource != nil {
		permission.Resource = *req.Resource
	}
	if req.Action != nil {
		permission.Action = *req.Action
	}
	if req.Description != nil {
		permission.Description = *req.Description
	}

	if err := h.permissionRepo.Update(ctx, permission); err != nil {
		logger.Error("Fout bij bijwerken permission", "error", err, "permission_id", permissionID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permission niet bijwerken",
		})
	}

	return c.JSON(permission)
}

// DeletePermission verwijdert een permission
func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	permissionID := c.Params("id")
	if permissionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Permission ID is verplicht",
		})
	}

	ctx := c.Context()

	// Haal permission op om te controleren of het een systeempermission is
	permission, err := h.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		logger.Error("Fout bij ophalen permission", "error", err, "permission_id", permissionID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Permission niet gevonden",
		})
	}

	// Controleer of het een systeempermission is
	if permission.IsSystemPermission {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Kan systeempermission niet verwijderen",
		})
	}

	// Verwijder eerst alle role-permission relaties
	if err := h.rolePermissionRepo.DeleteByPermissionID(ctx, permissionID); err != nil {
		logger.Error("Fout bij verwijderen role-permission relaties", "error", err, "permission_id", permissionID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon role-permission relaties niet verwijderen",
		})
	}

	// Verwijder de permission
	if err := h.permissionRepo.Delete(ctx, permissionID); err != nil {
		logger.Error("Fout bij verwijderen permission", "error", err, "permission_id", permissionID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permission niet verwijderen",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Permission verwijderd",
	})
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

// GetRole haalt details van een specifieke rol op
func (h *PermissionHandler) GetRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID is verplicht",
		})
	}

	ctx := c.Context()
	role, err := h.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen rol", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rol niet gevonden",
		})
	}

	return c.JSON(role)
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

// UpdateRole werkt een rol bij
func (h *PermissionHandler) UpdateRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID is verplicht",
		})
	}

	var req struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	ctx := c.Context()

	// Haal huidige rol op
	role, err := h.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen rol", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rol niet gevonden",
		})
	}

	// Controleer of het een systeemrol is
	if role.IsSystemRole {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Kan systeemrol niet bewerken",
		})
	}

	// Update velden indien opgegeven
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := h.roleRepo.Update(ctx, role); err != nil {
		logger.Error("Fout bij bijwerken rol", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon rol niet bijwerken",
		})
	}

	return c.JSON(role)
}

// DeleteRole verwijdert een rol
func (h *PermissionHandler) DeleteRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID is verplicht",
		})
	}

	ctx := c.Context()

	// Haal rol op om te controleren of het een systeemrol is
	role, err := h.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen rol", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rol niet gevonden",
		})
	}

	// Controleer of het een systeemrol is
	if role.IsSystemRole {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Kan systeemrol niet verwijderen",
		})
	}

	// Verwijder eerst alle role-permission relaties
	if err := h.rolePermissionRepo.DeleteByRoleID(ctx, roleID); err != nil {
		logger.Error("Fout bij verwijderen role-permission relaties", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon role-permission relaties niet verwijderen",
		})
	}

	// Verwijder user-role relaties
	if err := h.userRoleRepo.DeleteByRole(ctx, roleID); err != nil {
		logger.Error("Fout bij verwijderen user-role relaties", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon user-role relaties niet verwijderen",
		})
	}

	// Verwijder de rol
	if err := h.roleRepo.Delete(ctx, roleID); err != nil {
		logger.Error("Fout bij verwijderen rol", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon rol niet verwijderen",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Rol verwijderd",
	})
}

// GetMenuPermissions haalt alle menu permissies op
func (h *PermissionHandler) GetMenuPermissions(c *fiber.Ctx) error {
	ctx := c.Context()

	// Haal alle menu permissies op
	permissions, err := h.permissionRepo.List(ctx, 1000, 0) // Hoge limit voor alle menu items
	if err != nil {
		logger.Error("Fout bij ophalen menu permissions", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon menu permissions niet ophalen",
		})
	}

	// Filter alleen menu permissies
	var menuPermissions []*models.Permission
	for _, perm := range permissions {
		if perm.Resource == "menu" {
			menuPermissions = append(menuPermissions, perm)
		}
	}

	// Sorteer op action (menu item naam)
	sort.Slice(menuPermissions, func(i, j int) bool {
		return menuPermissions[i].Action < menuPermissions[j].Action
	})

	return c.JSON(fiber.Map{
		"permissions": menuPermissions,
		"total":       len(menuPermissions),
	})
}

// GetMenuPermissionMatrix geeft een matrix view van rollen vs menu permissies
func (h *PermissionHandler) GetMenuPermissionMatrix(c *fiber.Ctx) error {
	ctx := c.Context()

	// Haal alle rollen op
	roles, err := h.roleRepo.ListWithPermissions(ctx, 100, 0)
	if err != nil {
		logger.Error("Fout bij ophalen roles voor matrix", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon roles niet ophalen",
		})
	}

	// Haal alle menu permissies op
	allPermissions, err := h.permissionRepo.List(ctx, 1000, 0)
	if err != nil {
		logger.Error("Fout bij ophalen permissions voor matrix", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permissions niet ophalen",
		})
	}

	// Filter menu permissies
	var menuPermissions []*models.Permission
	for _, perm := range allPermissions {
		if perm.Resource == "menu" {
			menuPermissions = append(menuPermissions, perm)
		}
	}

	// Sorteer menu permissies
	sort.Slice(menuPermissions, func(i, j int) bool {
		return menuPermissions[i].Action < menuPermissions[j].Action
	})

	// Bouw matrix: rol -> menu_permission_id -> heeft_permission
	matrix := make(map[string]map[string]bool)
	for _, role := range roles {
		matrix[role.ID] = make(map[string]bool)
		for _, perm := range menuPermissions {
			matrix[role.ID][perm.ID] = false // Default: geen toegang
		}

		// Markeer welke permissies deze rol heeft
		for _, rolePerm := range role.Permissions {
			if rolePerm.Resource == "menu" {
				matrix[role.ID][rolePerm.ID] = true
			}
		}
	}

	return c.JSON(fiber.Map{
		"roles":       roles,
		"permissions": menuPermissions,
		"matrix":      matrix,
	})
}

// GetRoleMenuPermissions haalt menu permissies voor een specifieke rol op
func (h *PermissionHandler) GetRoleMenuPermissions(c *fiber.Ctx) error {
	roleID := c.Params("id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID is verplicht",
		})
	}

	ctx := c.Context()

	// Haal rol op
	role, err := h.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen rol", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rol niet gevonden",
		})
	}

	// Haal alle menu permissies op
	allPermissions, err := h.permissionRepo.List(ctx, 1000, 0)
	if err != nil {
		logger.Error("Fout bij ophalen permissions", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permissions niet ophalen",
		})
	}

	// Haal permissies voor deze rol op
	rolePermissions, err := h.rolePermissionRepo.GetPermissionsByRole(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen role permissions", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon role permissions niet ophalen",
		})
	}

	// Filter menu permissies en markeer welke de rol heeft
	var menuPermissions []map[string]interface{}
	rolePermissionIDs := make(map[string]bool)
	for _, rp := range rolePermissions {
		rolePermissionIDs[rp.ID] = true
	}

	for _, perm := range allPermissions {
		if perm.Resource == "menu" {
			menuPermissions = append(menuPermissions, map[string]interface{}{
				"id":          perm.ID,
				"action":      perm.Action,
				"description": perm.Description,
				"has_access":  rolePermissionIDs[perm.ID],
			})
		}
	}

	// Sorteer op action
	sort.Slice(menuPermissions, func(i, j int) bool {
		return menuPermissions[i]["action"].(string) < menuPermissions[j]["action"].(string)
	})

	return c.JSON(fiber.Map{
		"role":        role,
		"permissions": menuPermissions,
		"total":       len(menuPermissions),
	})
}

// UpdateRoleMenuPermissions werkt menu permissies bij voor een rol
func (h *PermissionHandler) UpdateRoleMenuPermissions(c *fiber.Ctx) error {
	roleID := c.Params("id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID is verplicht",
		})
	}

	var req struct {
		PermissionActions []string `json:"permission_actions"` // Array van menu actions (bijv. ["dashboard", "emails"])
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

	// Haal alle menu permissies op om IDs te vinden
	allPermissions, err := h.permissionRepo.List(ctx, 1000, 0)
	if err != nil {
		logger.Error("Fout bij ophalen permissions", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon permissions niet ophalen",
		})
	}

	// Maak mapping van action naar permission ID voor menu items
	menuPermissionMap := make(map[string]string)
	for _, perm := range allPermissions {
		if perm.Resource == "menu" {
			menuPermissionMap[perm.Action] = perm.ID
		}
	}

	// Converteer actions naar permission IDs
	var permissionIDs []string
	for _, action := range req.PermissionActions {
		if id, exists := menuPermissionMap[action]; exists {
			permissionIDs = append(permissionIDs, id)
		} else {
			logger.Warn("Onbekende menu action", "action", action)
		}
	}

	// Haal huidige menu permissions voor deze role op
	currentPermissions, err := h.rolePermissionRepo.GetPermissionsByRole(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen huidige permissions", "error", err, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon huidige permissions niet ophalen",
		})
	}

	// Filter alleen huidige menu permissions
	currentMenuPermissionIDs := make(map[string]bool)
	for _, perm := range currentPermissions {
		if perm.Resource == "menu" {
			currentMenuPermissionIDs[perm.ID] = true
		}
	}

	requestedPermissionIDs := make(map[string]bool)
	for _, id := range permissionIDs {
		requestedPermissionIDs[id] = true
	}

	// Bepaal welke permissions toegevoegd/verwijderd moeten worden
	var toAdd []string
	var toRemove []string

	// Permissions die toegevoegd moeten worden (in request maar niet current)
	for _, permID := range permissionIDs {
		if !currentMenuPermissionIDs[permID] {
			toAdd = append(toAdd, permID)
		}
	}

	// Permissions die verwijderd moeten worden (current maar niet in request)
	for permID := range currentMenuPermissionIDs {
		if !requestedPermissionIDs[permID] {
			toRemove = append(toRemove, permID)
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
			logger.Error("Fout bij toevoegen menu permission", "error", err, "role_id", roleID, "permission_id", permissionID)
			continue
		}
		addedCount++
	}

	// Voer verwijderingen uit
	removedCount := 0
	for _, permissionID := range toRemove {
		if err := h.rolePermissionRepo.Delete(ctx, roleID, permissionID); err != nil {
			logger.Error("Fout bij verwijderen menu permission", "error", err, "role_id", roleID, "permission_id", permissionID)
			continue
		}
		removedCount++
	}

	return c.JSON(fiber.Map{
		"success":         true,
		"message":         "Menu permissions bijgewerkt",
		"added_count":     addedCount,
		"removed_count":   removedCount,
		"total_requested": len(req.PermissionActions),
	})
}

// ListUsersWithRoles haalt een lijst van users op met hun roles
func (h *PermissionHandler) ListUsersWithRoles(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)
	search := c.Query("search", "")

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

	// Haal users op via auth service
	users, err := h.authService.ListUsers(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen users", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon users niet ophalen",
		})
	}

	// Filter users gebaseerd op search parameter
	var filteredUsers []*models.Gebruiker
	for _, user := range users {
		if search != "" {
			searchLower := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(user.Email), searchLower) &&
				!strings.Contains(strings.ToLower(user.Naam), searchLower) {
				continue
			}
		}
		filteredUsers = append(filteredUsers, user)
	}

	// Voor elke user, haal de roles op
	var result []map[string]interface{}
	for _, user := range filteredUsers {
		userRoles, err := h.userRoleRepo.ListActiveByUser(ctx, user.ID)
		if err != nil {
			logger.Error("Fout bij ophalen user roles", "error", err, "user_id", user.ID)
			continue
		}

		// Converteer naar response format
		roles := make([]map[string]interface{}, 0, len(userRoles))
		for _, ur := range userRoles {
			roles = append(roles, map[string]interface{}{
				"id":          ur.Role.ID,
				"name":        ur.Role.Name,
				"description": ur.Role.Description,
				"assigned_at": ur.AssignedAt,
				"expires_at":  ur.ExpiresAt,
				"is_active":   ur.IsActive,
			})
		}

		result = append(result, map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"naam":  user.Naam,
			"roles": roles,
		})
	}

	return c.JSON(fiber.Map{
		"users": result,
		"total": len(result),
	})
}

// GetUserRoles haalt alle roles op voor een specifieke user
func (h *PermissionHandler) GetUserRoles(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is verplicht",
		})
	}

	ctx := c.Context()

	// Controleer of user bestaat
	user, err := h.authService.GetUser(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen user", "error", err, "user_id", userID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Gebruiker niet gevonden",
		})
	}

	// Haal user roles op
	userRoles, err := h.userRoleRepo.ListActiveByUser(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen user roles", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon user roles niet ophalen",
		})
	}

	// Converteer naar response format
	roles := make([]map[string]interface{}, 0, len(userRoles))
	for _, ur := range userRoles {
		roles = append(roles, map[string]interface{}{
			"id":          ur.Role.ID,
			"name":        ur.Role.Name,
			"description": ur.Role.Description,
			"assigned_at": ur.AssignedAt,
			"expires_at":  ur.ExpiresAt,
			"is_active":   ur.IsActive,
		})
	}

	return c.JSON(fiber.Map{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"naam":  user.Naam,
		},
		"roles": roles,
		"total": len(roles),
	})
}

// AssignRoleToUser kent een role toe aan een user
func (h *PermissionHandler) AssignRoleToUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is verplicht",
		})
	}

	var req struct {
		RoleID    string  `json:"role_id"`
		ExpiresAt *string `json:"expires_at,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	if req.RoleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID is verplicht",
		})
	}

	// Haal assignedBy op uit context
	assignedBy, ok := c.Locals("userID").(string)
	if !ok || assignedBy == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon assignedBy niet ophalen uit context",
		})
	}

	ctx := c.Context()

	// Controleer of user bestaat
	_, err := h.authService.GetUser(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen user", "error", err, "user_id", userID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Gebruiker niet gevonden",
		})
	}

	// Controleer of role bestaat
	role, err := h.roleRepo.GetByID(ctx, req.RoleID)
	if err != nil {
		logger.Error("Fout bij ophalen role", "error", err, "role_id", req.RoleID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rol niet gevonden",
		})
	}

	// Controleer of user deze role al heeft
	existing, err := h.userRoleRepo.GetByUserAndRole(ctx, userID, req.RoleID)
	if err != nil && err.Error() != "record not found" {
		logger.Error("Fout bij controleren bestaande user-role relatie", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon user-role relatie niet controleren",
		})
	}

	if existing != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Gebruiker heeft deze rol al",
		})
	}

	// Parse expires_at indien opgegeven
	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ongeldig expires_at format (gebruik RFC3339)",
			})
		}
		expiresAt = &parsed
	}

	// Maak user-role relatie aan
	userRole := &models.UserRole{
		UserID:     userID,
		RoleID:     req.RoleID,
		AssignedBy: &assignedBy,
		ExpiresAt:  expiresAt,
		IsActive:   true,
	}

	if err := h.userRoleRepo.Create(ctx, userRole); err != nil {
		logger.Error("Fout bij toewijzen role aan user", "error", err, "user_id", userID, "role_id", req.RoleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon role niet toewijzen aan gebruiker",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Rol toegewezen aan gebruiker",
		"user_role": map[string]interface{}{
			"id":          userRole.ID,
			"user_id":     userRole.UserID,
			"role_id":     userRole.RoleID,
			"role_name":   role.Name,
			"assigned_at": userRole.AssignedAt,
			"expires_at":  userRole.ExpiresAt,
			"is_active":   userRole.IsActive,
		},
	})
}

// RevokeRoleFromUser verwijdert een role van een user
func (h *PermissionHandler) RevokeRoleFromUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	roleID := c.Params("roleId")

	if userID == "" || roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID en Role ID zijn verplicht",
		})
	}

	ctx := c.Context()

	// Controleer of user-role relatie bestaat
	userRole, err := h.userRoleRepo.GetByUserAndRole(ctx, userID, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen user-role relatie", "error", err, "user_id", userID, "role_id", roleID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User-role relatie niet gevonden",
		})
	}

	// Verwijder de relatie
	if err := h.userRoleRepo.Delete(ctx, userRole.ID); err != nil {
		logger.Error("Fout bij verwijderen user-role relatie", "error", err, "user_role_id", userRole.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon role niet verwijderen van gebruiker",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Rol verwijderd van gebruiker",
	})
}

// UpdateUserRole werkt een user-role relatie bij (activate/deactivate)
func (h *PermissionHandler) UpdateUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	roleID := c.Params("roleId")

	if userID == "" || roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID en Role ID zijn verplicht",
		})
	}

	var req struct {
		IsActive  *bool   `json:"is_active,omitempty"`
		ExpiresAt *string `json:"expires_at,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	ctx := c.Context()

	// Controleer of user-role relatie bestaat
	userRole, err := h.userRoleRepo.GetByUserAndRole(ctx, userID, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen user-role relatie", "error", err, "user_id", userID, "role_id", roleID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User-role relatie niet gevonden",
		})
	}

	// Update velden indien opgegeven
	if req.IsActive != nil {
		userRole.IsActive = *req.IsActive
	}

	if req.ExpiresAt != nil {
		if *req.ExpiresAt == "" {
			userRole.ExpiresAt = nil
		} else {
			parsed, err := time.Parse(time.RFC3339, *req.ExpiresAt)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Ongeldig expires_at format (gebruik RFC3339)",
				})
			}
			userRole.ExpiresAt = &parsed
		}
	}

	if err := h.userRoleRepo.Update(ctx, userRole); err != nil {
		logger.Error("Fout bij bijwerken user-role relatie", "error", err, "user_role_id", userRole.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon user-role relatie niet bijwerken",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User-role relatie bijgewerkt",
		"user_role": map[string]interface{}{
			"id":          userRole.ID,
			"user_id":     userRole.UserID,
			"role_id":     userRole.RoleID,
			"assigned_at": userRole.AssignedAt,
			"expires_at":  userRole.ExpiresAt,
			"is_active":   userRole.IsActive,
		},
	})
}
