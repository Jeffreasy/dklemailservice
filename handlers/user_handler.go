package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	authService       services.AuthService
	permissionService services.PermissionService
	userRoleRepo      repository.UserRoleRepository
}

func NewUserHandler(authService services.AuthService, permissionService services.PermissionService, userRoleRepo repository.UserRoleRepository) *UserHandler {
	return &UserHandler{
		authService:       authService,
		permissionService: permissionService,
		userRoleRepo:      userRoleRepo,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/users", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "read"), h.ListUsers)
	app.Get("/api/users/:id", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "read"), h.GetUser)
	app.Post("/api/users", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "write"), h.CreateUser)
	app.Put("/api/users/:id", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "write"), h.UpdateUser)
	app.Put("/api/users/:id/roles", AuthMiddleware(h.authService), AdminPermissionMiddleware(h.permissionService), h.AssignRolesToUser)
	app.Delete("/api/users/:id", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "delete"), h.DeleteUser)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, err := h.authService.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(gebruiker)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.authService.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.authService.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.authService.DeleteUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *UserHandler) AssignRolesToUser(c *fiber.Ctx) error {
	targetUserID := c.Params("id")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is verplicht",
		})
	}

	var req struct {
		RoleIDs []string `json:"role_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	if len(req.RoleIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ten minste één role ID is verplicht",
		})
	}

	// Haal userID op uit context voor assigned_by
	currentUserID, ok := c.Locals("userID").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon userID niet ophalen uit context",
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
