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
	permissionRepo    repository.PermissionRepository
	roleRepo          repository.RBACRoleRepository
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewPermissionHandler maakt een nieuwe permission handler
func NewPermissionHandler(
	permissionRepo repository.PermissionRepository,
	roleRepo repository.RBACRoleRepository,
	authService services.AuthService,
	permissionService services.PermissionService,
) *PermissionHandler {
	return &PermissionHandler{
		permissionRepo:    permissionRepo,
		roleRepo:          roleRepo,
		authService:       authService,
		permissionService: permissionService,
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

// ListRoles haalt een lijst van roles op
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
	roles, err := h.roleRepo.List(ctx, limit, offset)
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
