package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware is een middleware die controleert of de gebruiker is ingelogd
func AuthMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Haal token op uit Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Geen Authorization header gevonden", "path", c.Path(), "ip", c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
				"code":  "NO_AUTH_HEADER",
			})
		}

		// Controleer of het een Bearer token is
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Ongeldige Authorization header", "header", authHeader, "path", c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Ongeldige Authorization header",
				"code":  "INVALID_AUTH_HEADER",
			})
		}

		// Valideer token
		token := parts[1]
		userID, err := authService.ValidateToken(token)
		if err != nil {
			// Bepaal error type voor betere frontend handling
			errorCode := "INVALID_TOKEN"
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "expired") || strings.Contains(errorMsg, "exp") {
				errorCode = "TOKEN_EXPIRED"
			} else if strings.Contains(errorMsg, "malformed") {
				errorCode = "TOKEN_MALFORMED"
			} else if strings.Contains(errorMsg, "signature") {
				errorCode = "TOKEN_SIGNATURE_INVALID"
			}

			logger.Warn("Token validatie gefaald",
				"error", err,
				"code", errorCode,
				"path", c.Path(),
				"ip", c.IP())

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Ongeldig token",
				"code":  errorCode,
			})
		}

		// Sla gebruiker ID op in context
		c.Locals("userID", userID)
		c.Locals("token", token)

		logger.Debug("Authenticatie succesvol", "user_id", userID, "path", c.Path())

		// Ga door naar volgende handler
		return c.Next()
	}
}

/* StaffMiddleware allows "admin" or "staff" roles */
func StaffMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Haal token op uit context
		token, ok := c.Locals("token").(string)
		if !ok || token == "" {
			logger.Warn("Geen token gevonden in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
			})
		}

		// Haal gebruiker op uit token
		ctx := c.Context()
		gebruiker, err := authService.GetUserFromToken(ctx, token)
		if err != nil {
			logger.Warn("Kon gebruiker niet ophalen uit token", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
			})
		}

		// Controleer of gebruiker admin of staff is
		if gebruiker.Rol != "admin" && gebruiker.Rol != "staff" {
			logger.Warn("Gebruiker is geen admin of staff", "user_id", gebruiker.ID, "role", gebruiker.Rol)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Geen toegang",
			})
		}

		// Sla gebruiker op in context
		c.Locals("gebruiker", gebruiker)

		// Ga door naar volgende handler
		return c.Next()
	}
}

// AdminMiddleware is een middleware die controleert of de gebruiker een admin is
func AdminMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Haal token op uit context
		token, ok := c.Locals("token").(string)
		if !ok || token == "" {
			logger.Warn("Geen token gevonden in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
			})
		}

		// Haal gebruiker op uit token
		ctx := c.Context()
		gebruiker, err := authService.GetUserFromToken(ctx, token)
		if err != nil {
			logger.Warn("Kon gebruiker niet ophalen uit token", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
			})
		}

		// Controleer of gebruiker admin is
		if gebruiker.Rol != "admin" {
			logger.Warn("Gebruiker is geen admin", "user_id", gebruiker.ID, "role", gebruiker.Rol)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Geen toegang",
			})
		}

		// Sla gebruiker op in context
		c.Locals("gebruiker", gebruiker)

		// Ga door naar volgende handler
		return c.Next()
	}
}

// RateLimitMiddleware is een middleware die rate limiting toepast
func RateLimitMiddleware(rateLimiter services.RateLimiterService, keyPrefix string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Haal IP adres op
		ip := c.IP()
		if ip == "" {
			ip = "unknown"
		}

		// Maak rate limit key
		key := keyPrefix + ":" + ip

		// Controleer rate limit
		if !rateLimiter.Allow(key) {
			logger.Warn("Rate limit overschreden", "ip", ip, "key", key)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Te veel verzoeken, probeer het later opnieuw",
			})
		}

		// Ga door naar volgende handler
		return c.Next()
	}
}

// TestModeMiddleware controleert of de request in testmodus moet worden uitgevoerd
func TestModeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Controleer op test mode header
		if testMode := c.Get("X-Test-Mode"); testMode == "true" {
			c.Locals("test_mode", true)
			logger.Debug("Test modus geactiveerd via header", "path", c.Path(), "ip", c.IP())
		}

		// Controleer op test_mode query parameter
		if testMode := c.Query("test_mode"); testMode == "true" {
			c.Locals("test_mode", true)
			logger.Debug("Test modus geactiveerd via query parameter", "path", c.Path(), "ip", c.IP())
		}

		// Ga verder met de request
		return c.Next()
	}
}
