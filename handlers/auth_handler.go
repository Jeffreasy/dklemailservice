package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler bevat handlers voor authenticatie
type AuthHandler struct {
	authService       services.AuthService
	permissionService services.PermissionService
	rateLimiter       services.RateLimiterService
}

// NewAuthHandler maakt een nieuwe AuthHandler
func NewAuthHandler(authService services.AuthService, permissionService services.PermissionService, rateLimiter services.RateLimiterService) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		permissionService: permissionService,
		rateLimiter:       rateLimiter,
	}
}

// HandleLogin handelt login verzoeken af
func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	// Parse request body
	var loginData models.GebruikerLogin
	if err := c.BodyParser(&loginData); err != nil {
		logger.Error("Fout bij parsen login data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige login data",
		})
	}

	// Valideer input
	if loginData.Email == "" || loginData.Wachtwoord == "" {
		logger.Warn("Ontbrekende login gegevens")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email en wachtwoord zijn verplicht",
		})
	}

	// Rate limiting voor login pogingen
	rateLimitKey := "login:" + loginData.Email
	if !h.rateLimiter.Allow(rateLimitKey) {
		logger.Warn("Rate limit overschreden voor login", "email", loginData.Email)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Te veel login pogingen, probeer het later opnieuw",
		})
	}

	// Authenticeer gebruiker
	token, err := h.authService.Login(c.Context(), loginData.Email, loginData.Wachtwoord)
	if err != nil {
		// Specifieke foutafhandeling
		switch err {
		case services.ErrInvalidCredentials:
			logger.Warn("Ongeldige inloggegevens", "email", loginData.Email)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Ongeldige inloggegevens",
			})
		case services.ErrUserInactive:
			logger.Warn("Inactieve gebruiker", "email", loginData.Email)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Gebruiker is inactief",
			})
		default:
			logger.Error("Fout bij login", "email", loginData.Email, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Er is een fout opgetreden bij het inloggen",
			})
		}
	}

	// Stel cookie in met token
	cookie := fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Strict",
	}
	c.Cookie(&cookie)

	// Stuur token terug in response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   token,
		"message": "Login succesvol",
	})
}

// HandleLogout handelt logout verzoeken af
func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	// Verwijder cookie
	c.ClearCookie("auth_token")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout succesvol",
	})
}

// HandleResetPassword handelt wachtwoord reset verzoeken af
func (h *AuthHandler) HandleResetPassword(c *fiber.Ctx) error {
	// Alleen toegankelijk voor ingelogde gebruikers
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		logger.Warn("Geen gebruiker ID gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	// Parse request body
	var resetData struct {
		HuidigWachtwoord string `json:"huidig_wachtwoord"`
		NieuwWachtwoord  string `json:"nieuw_wachtwoord"`
	}
	if err := c.BodyParser(&resetData); err != nil {
		logger.Error("Fout bij parsen wachtwoord reset data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige wachtwoord reset data",
		})
	}

	// Valideer input
	if resetData.HuidigWachtwoord == "" || resetData.NieuwWachtwoord == "" {
		logger.Warn("Ontbrekende wachtwoord reset gegevens")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Huidig wachtwoord en nieuw wachtwoord zijn verplicht",
		})
	}

	// Haal gebruiker op
	gebruiker, err := h.authService.GetUserFromToken(c.Context(), c.Locals("token").(string))
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het ophalen van de gebruiker",
		})
	}

	// Verifieer huidig wachtwoord
	if !h.authService.VerifyPassword(gebruiker.WachtwoordHash, resetData.HuidigWachtwoord) {
		logger.Warn("Ongeldig huidig wachtwoord", "user_id", userID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Ongeldig huidig wachtwoord",
		})
	}

	// Reset wachtwoord
	if err := h.authService.ResetPassword(c.Context(), gebruiker.Email, resetData.NieuwWachtwoord); err != nil {
		logger.Error("Fout bij resetten wachtwoord", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het resetten van het wachtwoord",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Wachtwoord succesvol gewijzigd",
	})
}

// HandleGetProfile handelt verzoeken af om het gebruikersprofiel op te halen
func (h *AuthHandler) HandleGetProfile(c *fiber.Ctx) error {
	// Haal user ID op uit context (gezet door AuthMiddleware)
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		logger.Warn("Geen user ID gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	// Haal gebruiker op
	gebruiker, err := h.authService.GetUser(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon gebruiker niet ophalen",
		})
	}

	if gebruiker == nil {
		logger.Warn("Gebruiker niet gevonden", "user_id", userID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Gebruiker niet gevonden",
		})
	}

	// Haal permissies op via RBAC systeem
	permissions, err := h.permissionService.GetUserPermissions(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen permissies", "user_id", userID, "error", err)
		// Fallback naar lege array als permissies niet opgehaald kunnen worden
		permissions = []*models.UserPermission{}
	}

	// Converteer permissies naar frontend format
	permissionList := make([]map[string]string, len(permissions))
	for i, perm := range permissions {
		permissionList[i] = map[string]string{
			"resource": perm.Resource,
			"action":   perm.Action,
		}
	}

	// Haal rollen op
	userRoles, err := h.permissionService.GetUserRoles(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen rollen", "user_id", userID, "error", err)
		userRoles = []*models.UserRole{}
	}

	// Converteer rollen naar frontend format
	roleList := make([]map[string]interface{}, len(userRoles))
	for i, userRole := range userRoles {
		roleList[i] = map[string]interface{}{
			"id":          userRole.Role.ID,
			"name":        userRole.Role.Name,
			"description": userRole.Role.Description,
			"assigned_at": userRole.AssignedAt,
			"is_active":   userRole.IsActive,
		}
	}

	// Stuur gebruikersprofiel terug met permissies en rollen
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            gebruiker.ID,
		"naam":          gebruiker.Naam,
		"email":         gebruiker.Email,
		"rol":           gebruiker.Rol, // Legacy field voor backward compatibility
		"permissions":   permissionList,
		"roles":         roleList,
		"is_actief":     gebruiker.IsActief,
		"laatste_login": gebruiker.LaatsteLogin,
		"created_at":    gebruiker.CreatedAt,
	})
}
