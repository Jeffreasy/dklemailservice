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
	authService services.AuthService
	rateLimiter services.RateLimiterService
}

// NewAuthHandler maakt een nieuwe AuthHandler
func NewAuthHandler(authService services.AuthService, rateLimiter services.RateLimiterService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		rateLimiter: rateLimiter,
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
	// Haal token op uit context
	token, ok := c.Locals("token").(string)
	if !ok || token == "" {
		logger.Warn("Geen token gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	// Haal gebruiker op uit token
	gebruiker, err := h.authService.GetUserFromToken(c.Context(), token)
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker uit token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	// Stuur gebruikersprofiel terug
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            gebruiker.ID,
		"naam":          gebruiker.Naam,
		"email":         gebruiker.Email,
		"rol":           gebruiker.Rol,
		"is_actief":     gebruiker.IsActief,
		"laatste_login": gebruiker.LaatsteLogin,
		"created_at":    gebruiker.CreatedAt,
	})
}
