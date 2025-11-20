package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/services"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MIGRATION STATUS: Legacy Role-Based Middleware Removal (V38)
// ============================================================
// The following legacy middleware functions have been REMOVED in V38:
// - StaffMiddleware: Used legacy gebruiker.Rol field instead of RBAC permissions
// - AdminMiddleware: Used legacy gebruiker.Rol field instead of RBAC permissions
//
// Migration completed:
// ✅ Added deprecation warnings with removal timeline (V37)
// ✅ Verified no active usage in current routes
// ✅ All routes use modern permission-based middleware
// ✅ REMOVED legacy middleware functions (V38)
//
// Modern alternatives:
// - StaffPermissionMiddleware: Uses permissionService.HasPermission(userID, "staff", "access")
// - AdminPermissionMiddleware: Uses permissionService.HasPermission(userID, "admin", "access")
//
// Next steps:
// - Remove legacy Rol column from database in V39
// - Update JWT tokens to remove legacy Role field
// - Clean up User API to remove rol field support

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

		// Access Token Rotation: Controleer of token niet ingetrokken is in database
		// Dit gebeurt alleen als de JWT validatie slaagt
		if err := authService.ValidateAccessToken(c.Context(), token); err != nil {
			logger.Warn("Access token validatie gefaald",
				"user_id", userID,
				"error", err,
				"path", c.Path(),
				"ip", c.IP())

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Access token ingetrokken of ongeldig",
				"code":  "TOKEN_REVOKED",
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

// SessionValidationMiddleware valideert dat de sessie gekoppeld aan de access token nog actief is
// Dit middleware moet NA AuthMiddleware worden gebruikt
func SessionValidationMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Haal user ID en token op uit context (gezet door AuthMiddleware)
		userID, ok := c.Locals("userID").(string)
		if !ok || userID == "" {
			logger.Warn("Geen user ID gevonden in context voor session validatie", "path", c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
				"code":  "NO_USER_ID",
			})
		}

		token, ok := c.Locals("token").(string)
		if !ok || token == "" {
			logger.Warn("Geen token gevonden in context voor session validatie", "path", c.Path(), "user_id", userID)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
				"code":  "NO_TOKEN",
			})
		}

		// Controleer of er een actieve sessie is voor deze gebruiker en token
		// Dit gebeurt alleen als session repository beschikbaar is
		if authService != nil {
			sessions, err := authService.ListUserSessions(c.Context(), userID)
			if err != nil {
				logger.Error("Fout bij ophalen sessies voor validatie", "user_id", userID, "error", err)
				// Graceful degradation: ga door als sessie check faalt
				logger.Warn("Session validatie overgeslagen vanwege error", "user_id", userID, "path", c.Path())
			} else {
				// FIX: Zoek specifiek naar de sessie die bij DIT token hoort.
				found := false
				for _, session := range sessions {
					if session.AccessToken == token {
						if session.IsActive && !session.IsExpired() {
							// Sessie gevonden en geldig!
							c.Locals("sessionID", session.ID) // Sla op voor Logout handler
							found = true
						}
						// We stoppen met zoeken zodra we de token match hebben,
						// ongeacht of hij geldig is (want token is uniek per sessie).
						break
					}
				}

				if !found {
					logger.Warn("Huidige sessie niet gevonden, inactief of verlopen", "user_id", userID, "path", c.Path())
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "Sessie verlopen of ingetrokken",
						"code":  "SESSION_EXPIRED",
					})
				}

				logger.Debug("Session validatie succesvol", "user_id", userID, "path", c.Path())
			}
		}

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

// SecurityHeadersMiddleware voegt beveiligingsheaders toe aan alle responses
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Basis beveiligingsheaders
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")

		// HSTS (HTTP Strict Transport Security) - alleen in productie
		if os.Getenv("ENV") == "production" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Content Security Policy - configureerbaar via environment variable
		if csp := os.Getenv("CONTENT_SECURITY_POLICY"); csp != "" {
			c.Set("Content-Security-Policy", csp)
		} else {
			// Default CSP voor development - meer permissief voor Swagger UI
			c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; script-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src 'self' https:;")
		}

		// Referrer Policy
		if rp := os.Getenv("REFERRER_POLICY"); rp != "" {
			c.Set("Referrer-Policy", rp)
		} else {
			c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		}

		// Permissions Policy (voorheen Feature Policy)
		if pp := os.Getenv("PERMISSIONS_POLICY"); pp != "" {
			c.Set("Permissions-Policy", pp)
		}

		logger.Debug("Security headers toegevoegd", "path", c.Path(), "production", os.Getenv("ENV") == "production")

		return c.Next()
	}
}
