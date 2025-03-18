package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MailHandler bevat handlers voor e-mail beheer
type MailHandler struct {
	mailFetcher       *services.MailFetcher
	incomingEmailRepo repository.IncomingEmailRepository
	authService       services.AuthService
	lastRun           time.Time
}

// NewMailHandler maakt een nieuwe MailHandler
func NewMailHandler(
	mailFetcher *services.MailFetcher,
	incomingEmailRepo repository.IncomingEmailRepository,
	authService services.AuthService,
) *MailHandler {
	return &MailHandler{
		mailFetcher:       mailFetcher,
		incomingEmailRepo: incomingEmailRepo,
		authService:       authService,
		lastRun:           time.Now().Add(-24 * time.Hour),
	}
}

// RegisterRoutes registreert de routes voor mail beheer
func (h *MailHandler) RegisterRoutes(app *fiber.App) {
	// Groep voor mail beheer routes (vereist admin rechten of API key)
	mailGroup := app.Group("/api/mail")

	// Custom auth middleware die zowel API key als JWT token accepteert
	mailGroup.Use(func(c *fiber.Ctx) error {
		// Haal token op uit Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Geen Authorization header gevonden")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Niet geautoriseerd",
			})
		}

		// Check voor API key als Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Ongeldige Authorization header", "header", authHeader)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Ongeldige Authorization header",
			})
		}

		token := parts[1]

		// Controleer eerst of het een API key is
		adminAPIKey := os.Getenv("ADMIN_API_KEY")
		if token == adminAPIKey {
			// Geldige API key, ga door
			logger.Info("Mail API toegang verleend via API key")
			return c.Next()
		}

		// Geen geldige API key, probeer JWT token
		userID, err := h.authService.ValidateToken(token)
		if err != nil {
			logger.Error("Fout bij valideren token", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Ongeldig token",
			})
		}

		// JWT token is geldig, sla gebruiker ID op in context
		c.Locals("userID", userID)
		c.Locals("token", token)

		// Controleer nog of gebruiker admin is
		ctx := c.Context()
		gebruiker, err := h.authService.GetUserFromToken(ctx, token)
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
		logger.Info("Mail API toegang verleend via JWT token", "user_id", gebruiker.ID)

		return c.Next()
	})

	// Mail beheer routes
	mailGroup.Get("/", h.ListEmails)
	mailGroup.Get("/:id", h.GetEmail)
	mailGroup.Put("/:id/processed", h.MarkAsProcessed)
	mailGroup.Delete("/:id", h.DeleteEmail)
	mailGroup.Post("/fetch", h.FetchEmails)
	mailGroup.Get("/unprocessed", h.ListUnprocessedEmails)
	mailGroup.Get("/account/:type", h.ListEmailsByAccountType)
}

// ListEmails haalt een lijst van emails op
// @Summary Lijst van emails ophalen
// @Description Haalt een gepagineerde lijst van emails op
// @Tags Mail
// @Accept json
// @Produce json
// @Param limit query int false "Aantal resultaten per pagina (standaard 10)"
// @Param offset query int false "Offset voor paginering (standaard 0)"
// @Success 200 {array} models.IncomingEmail
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail [get]
// @Security BearerAuth
func (h *MailHandler) ListEmails(c *fiber.Ctx) error {
	// Haal query parameters op
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Valideer parameters
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

	// Haal emails op
	ctx := c.Context()
	emails, err := h.incomingEmailRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen emails", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon emails niet ophalen",
		})
	}

	// Stuur resultaat terug
	return c.JSON(emails)
}

// GetEmail haalt een specifieke email op
// @Summary Details van een specifieke email ophalen
// @Description Haalt de details van een specifieke email op
// @Tags Mail
// @Accept json
// @Produce json
// @Param id path string true "Email ID"
// @Success 200 {object} models.IncomingEmail
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/{id} [get]
// @Security BearerAuth
func (h *MailHandler) GetEmail(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Haal email op
	ctx := c.Context()
	email, err := h.incomingEmailRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen email", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon email niet ophalen",
		})
	}

	if email == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Email niet gevonden",
		})
	}

	// Stuur resultaat terug
	return c.JSON(email)
}

// MarkAsProcessed markeert een email als verwerkt
// @Summary Email als verwerkt markeren
// @Description Markeert een email als verwerkt om aan te geven dat deze is afgehandeld
// @Tags Mail
// @Accept json
// @Produce json
// @Param id path string true "Email ID"
// @Success 200 {object} models.IncomingEmail
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/{id}/processed [put]
// @Security BearerAuth
func (h *MailHandler) MarkAsProcessed(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Haal email op
	ctx := c.Context()
	email, err := h.incomingEmailRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen email", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon email niet ophalen",
		})
	}

	if email == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Email niet gevonden",
		})
	}

	// Markeer als verwerkt
	email.IsProcessed = true
	now := time.Now()
	email.ProcessedAt = &now

	// Sla wijzigingen op
	if err := h.incomingEmailRepo.Update(ctx, email); err != nil {
		logger.Error("Fout bij markeren email als verwerkt", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon email niet markeren als verwerkt",
		})
	}

	// Stuur bijgewerkte email terug
	return c.JSON(email)
}

// DeleteEmail verwijdert een email
// @Summary Email verwijderen
// @Description Verwijdert een email uit het systeem
// @Tags Mail
// @Accept json
// @Produce json
// @Param id path string true "Email ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/{id} [delete]
// @Security BearerAuth
func (h *MailHandler) DeleteEmail(c *fiber.Ctx) error {
	// Haal ID op uit URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Controleer of email bestaat
	ctx := c.Context()
	email, err := h.incomingEmailRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen email", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon email niet ophalen",
		})
	}

	if email == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Email niet gevonden",
		})
	}

	// Verwijder email
	if err := h.incomingEmailRepo.Delete(ctx, id); err != nil {
		logger.Error("Fout bij verwijderen email", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon email niet verwijderen",
		})
	}

	// Stuur bevestiging terug
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Email succesvol verwijderd",
	})
}

// FetchEmails haalt nieuwe emails op van de mailserver
// @Summary Nieuwe emails ophalen
// @Description Haalt nieuwe emails op van de mailserver en slaat deze op in de database
// @Tags Mail
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/fetch [post]
// @Security BearerAuth
func (h *MailHandler) FetchEmails(c *fiber.Ctx) error {
	// Haal nieuwe emails op
	emails, err := h.mailFetcher.FetchMails()
	if err != nil {
		logger.Error("Fout bij ophalen nieuwe emails", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwe emails niet ophalen: " + err.Error(),
		})
	}

	// Sla emails op in database (als ze nog niet bestaan)
	saved := 0
	ctx := c.Context()

	for _, email := range emails {
		// Controleer of email al bestaat (op basis van UID)
		existing, err := h.incomingEmailRepo.FindByUID(ctx, email.UID)
		if err != nil {
			logger.Error("Fout bij zoeken bestaande email", "error", err, "uid", email.UID)
			continue
		}

		// Als email nog niet bestaat, sla op
		if existing == nil {
			if err := h.incomingEmailRepo.Create(ctx, email); err != nil {
				logger.Error("Fout bij opslaan nieuwe email", "error", err, "uid", email.UID)
				continue
			}
			saved++
		}
	}

	// Bijwerken laatste uitvoering
	h.lastRun = time.Now()

	// Stuur resultaat terug
	return c.JSON(fiber.Map{
		"success":      true,
		"emails_found": len(emails),
		"emails_saved": saved,
		"last_run":     h.lastRun,
		"message":      "Emails succesvol opgehaald",
	})
}

// ListUnprocessedEmails haalt alle onverwerkte emails op
// @Summary Onverwerkte emails ophalen
// @Description Haalt een lijst van alle onverwerkte emails op
// @Tags Mail
// @Accept json
// @Produce json
// @Success 200 {array} models.IncomingEmail
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/unprocessed [get]
// @Security BearerAuth
func (h *MailHandler) ListUnprocessedEmails(c *fiber.Ctx) error {
	// Haal onverwerkte emails op
	ctx := c.Context()
	emails, err := h.incomingEmailRepo.FindUnprocessed(ctx)
	if err != nil {
		logger.Error("Fout bij ophalen onverwerkte emails", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon onverwerkte emails niet ophalen",
		})
	}

	// Stuur resultaat terug
	return c.JSON(emails)
}

// ListEmailsByAccountType haalt emails op basis van account type op
// @Summary Emails filteren op account type
// @Description Haalt een lijst van emails op gefilterd op account type
// @Tags Mail
// @Accept json
// @Produce json
// @Param type path string true "Account type (info, inschrijving)"
// @Success 200 {array} models.IncomingEmail
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/account/{type} [get]
// @Security BearerAuth
func (h *MailHandler) ListEmailsByAccountType(c *fiber.Ctx) error {
	// Haal account type op uit URL
	accountType := c.Params("type")
	if accountType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account type is verplicht",
		})
	}

	// Valideer account type
	if accountType != "info" && accountType != "inschrijving" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldig account type (gebruik 'info' of 'inschrijving')",
		})
	}

	// Haal emails op
	ctx := c.Context()
	emails, err := h.incomingEmailRepo.FindByAccountType(ctx, accountType)
	if err != nil {
		logger.Error("Fout bij ophalen emails op account type", "error", err, "type", accountType)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon emails niet ophalen",
		})
	}

	// Stuur resultaat terug
	return c.JSON(emails)
}

// GetEmails haalt alle inkomende emails op
func (h *MailHandler) GetEmails(c *fiber.Ctx) error {
	// Admin check
	user, err := h.authService.GetUserFromToken(c.Context(), string(c.Request().Header.Peek("Authorization")))
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	// Controleer of gebruiker admin is (direct op rol)
	if user.Rol != "admin" {
		logger.Warn("Toegang geweigerd: gebruiker is geen admin", "user_id", user.ID, "role", user.Rol)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Alleen beheerders hebben toegang tot deze functie",
		})
	}

	// Debug logging voor database issues
	logger.Info("GetEmails API aangeroepen door admin", "user_id", user.ID, "user_email", user.Email)

	// Haal emails op uit de repository
	emails, err := h.incomingEmailRepo.List(c.Context(), 100, 0)
	if err != nil {
		logger.Error("Fout bij ophalen inkomende e-mails", "error", err)
		logger.Error("Fout bij ophalen emails", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon emails niet ophalen",
		})
	}

	return c.JSON(fiber.Map{
		"emails": emails,
	})
}
