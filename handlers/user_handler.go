package handlers

import (
	"dklautomationgo/models"
	"dklautomationgo/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	authService       services.AuthService
	permissionService services.PermissionService
}

func NewUserHandler(authService services.AuthService, permissionService services.PermissionService) *UserHandler {
	return &UserHandler{
		authService:       authService,
		permissionService: permissionService,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/users", AuthMiddleware(h.authService), StaffPermissionMiddleware(h.permissionService), h.ListUsers)
	app.Get("/api/users/:id", AuthMiddleware(h.authService), StaffPermissionMiddleware(h.permissionService), h.GetUser)
	app.Post("/api/users", AuthMiddleware(h.authService), AdminPermissionMiddleware(h.permissionService), h.CreateUser)
	app.Put("/api/users/:id", AuthMiddleware(h.authService), AdminPermissionMiddleware(h.permissionService), h.UpdateUser)
	app.Delete("/api/users/:id", AuthMiddleware(h.authService), AdminPermissionMiddleware(h.permissionService), h.DeleteUser)
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
