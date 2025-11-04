package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	authService       services.AuthService
	permissionService services.PermissionService
	userRoleRepo      repository.UserRoleRepository
	roleRepo          repository.RBACRoleRepository
}

func NewUserHandler(authService services.AuthService, permissionService services.PermissionService, userRoleRepo repository.UserRoleRepository, roleRepo repository.RBACRoleRepository) *UserHandler {
	return &UserHandler{
		authService:       authService,
		permissionService: permissionService,
		userRoleRepo:      userRoleRepo,
		roleRepo:          roleRepo,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/users", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "read"), h.ListUsers)
	app.Get("/api/users/search", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "read"), h.SearchUsers)
	app.Get("/api/users/:id", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "read"), h.GetUser)
	app.Post("/api/users", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "write"), h.CreateUser)
	app.Put("/api/users/:id", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "write"), h.UpdateUser)

	// RBAC User-Role Management
	app.Get("/api/users/:id/roles", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "read"), h.GetUserRoles)
	app.Post("/api/users/:id/roles", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "manage_roles"), h.AssignRoleToUser)
	app.Put("/api/users/:id/roles", AuthMiddleware(h.authService), AdminPermissionMiddleware(h.permissionService), h.AssignRolesToUser)
	app.Delete("/api/users/:id/roles/:roleId", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "manage_roles"), h.RemoveRoleFromUser)
	app.Get("/api/users/:id/permissions", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "read"), h.GetUserPermissions)

	app.Delete("/api/users/:id", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "delete"), h.DeleteUser)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	query := c.Query("q") // Support search query parameter

	var users []*models.Gebruiker
	var err error

	if query != "" {
		// Use search functionality if query parameter is provided
		users, err = h.authService.SearchUsers(c.Context(), query, limit)
	} else {
		// Use regular list functionality
		users, err = h.authService.ListUsers(c.Context(), limit, offset)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "INTERNAL_ERROR",
		})
	}
	return c.JSON(users)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req struct {
		Email                string `json:"email"`
		Naam                 string `json:"naam"`
		Rol                  string `json:"rol"`
		Password             string `json:"password"`
		IsActief             bool   `json:"is_actief"`
		NewsletterSubscribed bool   `json:"newsletter_subscribed"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
			"code":  "INVALID_INPUT",
		})
	}

	gebruiker := &models.Gebruiker{
		Email:                req.Email,
		Naam:                 req.Naam,
		Rol:                  req.Rol,
		IsActief:             req.IsActief,
		NewsletterSubscribed: req.NewsletterSubscribed,
	}

	err := h.authService.CreateUser(c.Context(), gebruiker, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "USER_CREATION_FAILED",
		})
	}
	return c.JSON(gebruiker)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.authService.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
			"code":  "USER_NOT_FOUND",
		})
	}
	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.authService.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
			"code":  "USER_NOT_FOUND",
		})
	}

	var req struct {
		Email                *string `json:"email,omitempty"`
		Naam                 *string `json:"naam,omitempty"`
		Rol                  *string `json:"rol,omitempty"`
		IsActief             *bool   `json:"is_actief,omitempty"`
		NewsletterSubscribed *bool   `json:"newsletter_subscribed,omitempty"`
		Password             *string `json:"password,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
			"code":  "INVALID_INPUT",
		})
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Naam != nil {
		user.Naam = *req.Naam
	}
	if req.Rol != nil {
		user.Rol = *req.Rol
	}
	if req.IsActief != nil {
		user.IsActief = *req.IsActief
	}
	if req.NewsletterSubscribed != nil {
		user.NewsletterSubscribed = *req.NewsletterSubscribed
	}

	err = h.authService.UpdateUser(c.Context(), user, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "USER_UPDATE_FAILED",
		})
	}
	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.authService.DeleteUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "USER_DELETION_FAILED",
		})
	}
	return c.JSON(fiber.Map{"success": true})
}

// GetUserRoles retrieves all roles for a user
func (h *UserHandler) GetUserRoles(c *fiber.Ctx) error {
	userID := c.Params("id")

	userRoles, err := h.userRoleRepo.GetByUserIDWithRoles(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen user roles", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon user roles niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	return c.JSON(userRoles)
}

// AssignRoleToUser assigns a single role to a user
func (h *UserHandler) AssignRoleToUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	var req struct {
		RoleID    string  `json:"role_id"`
		ExpiresAt *string `json:"expires_at,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
			"code":  "INVALID_INPUT",
		})
	}

	if req.RoleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role ID is verplicht",
			"code":  "MISSING_ROLE_ID",
		})
	}

	currentUserID, _ := c.Locals("userID").(string)

	ur := &models.UserRole{
		UserID:     userID,
		RoleID:     req.RoleID,
		AssignedBy: &currentUserID,
		IsActive:   true,
	}

	// Parse expires_at if provided
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ongeldige expires_at formaat (gebruik RFC3339)",
				"code":  "INVALID_DATE_FORMAT",
			})
		}
		ur.ExpiresAt = &t
	}

	// Create user-role relationship
	if err := h.userRoleRepo.Create(c.Context(), ur); err != nil {
		// Check for duplicate
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User heeft deze rol al",
				"code":  "DUPLICATE_ROLE",
			})
		}
		logger.Error("Fout bij toewijzen role", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon role niet toewijzen",
			"code":  "ASSIGNMENT_FAILED",
		})
	}

	// Load full role object for response
	role, err := h.roleRepo.GetByID(c.Context(), req.RoleID)
	if err == nil {
		ur.Role = *role
	}

	return c.Status(fiber.StatusCreated).JSON(ur)
}

// GetUserPermissions retrieves effective permissions for a user
func (h *UserHandler) GetUserPermissions(c *fiber.Ctx) error {
	userID := c.Params("id")

	permissions, err := h.permissionService.GetUserPermissions(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen user permissions", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon user permissions niet ophalen",
			"code":  "INTERNAL_ERROR",
		})
	}

	return c.JSON(permissions)
}

func (h *UserHandler) AssignRolesToUser(c *fiber.Ctx) error {
	targetUserID := c.Params("id")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is verplicht",
			"code":  "MISSING_USER_ID",
		})
	}

	var req struct {
		RoleIDs []string `json:"role_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
			"code":  "INVALID_INPUT",
		})
	}

	if len(req.RoleIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ten minste één role ID is verplicht",
			"code":  "MISSING_ROLE_ID",
		})
	}

	// Haal userID op uit context voor assigned_by
	currentUserID, ok := c.Locals("userID").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon userID niet ophalen uit context",
			"code":  "INTERNAL_ERROR",
		})
	}

	ctx := c.Context()
	assignedRoles := 0

	for _, roleID := range req.RoleIDs {
		ur := &models.UserRole{
			UserID:     targetUserID,
			RoleID:     roleID,
			AssignedBy: &currentUserID,
			IsActive:   true,
		}

		if err := h.userRoleRepo.Create(ctx, ur); err != nil {
			logger.Error("Fout bij toewijzen role aan user", "error", err, "user_id", targetUserID, "role_id", roleID)
			// Continue with other roles
			continue
		}
		assignedRoles++
	}

	return c.JSON(fiber.Map{
		"success":         true,
		"message":         "Roles toegewezen aan user",
		"assigned_roles":  assignedRoles,
		"total_requested": len(req.RoleIDs),
	})
}

func (h *UserHandler) RemoveRoleFromUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	roleID := c.Params("roleId")
	if userID == "" || roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID en Role ID zijn verplicht",
			"code":  "MISSING_PARAMETERS",
		})
	}

	ctx := c.Context()

	// Find the user-role relationship
	userRole, err := h.userRoleRepo.GetByUserAndRole(ctx, userID, roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Gebruiker heeft deze rol niet",
				"code":  "ROLE_NOT_FOUND",
			})
		}
		logger.Error("Fout bij ophalen user-role relatie", "error", err, "user_id", userID, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon user-role relatie niet controleren",
			"code":  "INTERNAL_ERROR",
		})
	}

	// Deactivate the relationship
	if err := h.userRoleRepo.Deactivate(ctx, userRole.ID); err != nil {
		logger.Error("Fout bij verwijderen rol van user", "error", err, "user_id", userID, "role_id", roleID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon rol niet verwijderen van user",
			"code":  "REMOVAL_FAILED",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Rol verwijderd van user",
	})
}

// SearchUsers searches for users by name or email
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
			"code":  "MISSING_QUERY",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit > 50 {
		limit = 50 // Max limit
	}

	users, err := h.authService.SearchUsers(c.Context(), query, limit)
	if err != nil {
		logger.Error("Fout bij zoeken naar gebruikers", "query", query, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon gebruikers niet zoeken",
			"code":  "SEARCH_FAILED",
		})
	}

	return c.JSON(users)
}
