package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// MailResponse definieert de structuur voor API responses
// Dit zorgt ervoor dat de HTML in het 'html' veld zit voor de frontend.
type MailResponse struct {
	ID          string     `json:"id"`
	MessageID   string     `json:"message_id"`
	From        string     `json:"sender"`
	To          string     `json:"to"`
	Subject     string     `json:"subject"`
	HTML        string     `json:"html"` // Gedecodeerde/gesanitized HTML body
	ContentType string     `json:"content_type"`
	ReceivedAt  time.Time  `json:"received_at"`
	UID         string     `json:"uid"`
	AccountType string     `json:"account_type"`
	IsProcessed bool       `json:"read"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Helper functie om IncomingEmail naar MailResponse te mappen
func mapEmailToResponse(email *models.IncomingEmail) *MailResponse {
	if email == nil {
		return nil
	}
	return &MailResponse{
		ID:          email.ID,
		MessageID:   email.MessageID,
		From:        email.From,
		To:          email.To,
		Subject:     email.Subject,
		HTML:        email.Body, // Map Body -> HTML
		ContentType: email.ContentType,
		ReceivedAt:  email.ReceivedAt,
		UID:         email.UID,
		AccountType: email.AccountType,
		IsProcessed: email.IsProcessed,
		ProcessedAt: email.ProcessedAt,
		CreatedAt:   email.CreatedAt,
		UpdatedAt:   email.UpdatedAt,
	}
}

// Helper functie om een slice van IncomingEmail naar MailResponse te mappen
func mapEmailsToResponse(emails []*models.IncomingEmail) []*MailResponse {
	responses := make([]*MailResponse, len(emails))
	for i, email := range emails {
		responses[i] = mapEmailToResponse(email)
	}
	return responses
}

// PaginatedEmailsResponse definieert de structuur voor gepagineerde email responses
// Nu met *MailResponse type pointer
type PaginatedMailResponse struct {
	Emails     []*MailResponse `json:"emails"`
	TotalCount int64           `json:"totalCount"`
}

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
// @Success 200 {object} PaginatedMailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail [get]
// @Security BearerAuth
func (h *MailHandler) ListEmails(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	if limit < 1 || limit > 100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Limit moet tussen 1 en 100 liggen"})
	}
	if offset < 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Offset mag niet negatief zijn"})
	}

	ctx := c.Context()
	emails, totalCount, err := h.incomingEmailRepo.ListByAccountTypePaginated(ctx, "", limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen van gepagineerde e-mails", "error", err, "limit", limit, "offset", offset)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Kon emails niet ophalen"})
	}

	responseEmails := mapEmailsToResponse(emails)
	response := PaginatedMailResponse{
		Emails:     responseEmails,
		TotalCount: totalCount,
	}

	return c.JSON(response)
}

// GetEmail haalt een specifieke email op
// @Summary Details van een specifieke email ophalen
// @Description Haalt de details van een specifieke email op
// @Tags Mail
// @Accept json
// @Produce json
// @Param id path string true "Email ID"
// @Success 200 {object} MailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/{id} [get]
// @Security BearerAuth
func (h *MailHandler) GetEmail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID is verplicht"})
	}

	ctx := c.Context()
	email, err := h.incomingEmailRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("E-mail niet gevonden", "id", id)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "E-mail niet gevonden"})
		}
		logger.Error("Fout bij ophalen e-mail", "error", err, "id", id)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Interne serverfout bij ophalen e-mail"})
	}
	if email == nil {
		logger.Warn("E-mail niet gevonden (nil result)", "id", id)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "E-mail niet gevonden"})
	}

	response := mapEmailToResponse(email)
	// Log the response just before sending
	logger.Info("GetEmail response wordt verzonden", "email_id", id, "response_html_preview", getFirstNChars(response.HTML, 100)) // Log first 100 chars of HTML

	return c.JSON(response)
}

// MarkAsProcessed markeert een email als verwerkt
// @Summary Email als verwerkt markeren
// @Description Markeert een email als verwerkt om aan te geven dat deze is afgehandeld
// @Tags Mail
// @Accept json
// @Produce json
// @Param id path string true "Email ID"
// @Success 200 {object} MailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/{id}/processed [put]
// @Security BearerAuth
func (h *MailHandler) MarkAsProcessed(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID is verplicht"})
	}

	ctx := c.Context()
	email, err := h.incomingEmailRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen email voor markeren", "error", err, "id", id)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Kon email niet ophalen"})
	}
	if email == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Email niet gevonden"})
	}

	now := time.Now()
	email.IsProcessed = true
	email.ProcessedAt = &now

	err = h.incomingEmailRepo.Update(ctx, email)
	if err != nil {
		logger.Error("Fout bij markeren email als verwerkt", "error", err, "id", id)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Kon email niet bijwerken"})
	}

	logger.Info("Email gemarkeerd als verwerkt", "id", id)
	return c.JSON(fiber.Map{"message": fmt.Sprintf("Email %s gemarkeerd als verwerkt", id)})
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
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID is verplicht"})
	}

	ctx := c.Context()
	if err := h.incomingEmailRepo.Delete(ctx, id); err != nil {
		logger.Error("Fout bij verwijderen email", "error", err, "id", id)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Kon email niet verwijderen"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true, "message": "Email succesvol verwijderd"})
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
		logger.Error("Fout bij handmatig ophalen mails", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Fout bij ophalen emails: " + err.Error()})
	}

	// Sla opgehaalde emails op in de database
	var savedCount int
	for _, email := range emails {
		if err := h.incomingEmailRepo.Create(c.Context(), email); err != nil {
			logger.Error("Fout bij opslaan opgehaalde email", "error", err, "messageID", email.MessageID)
			// Ga door met de volgende email
		} else {
			savedCount++
		}
	}

	return c.JSON(fiber.Map{
		"message":   fmt.Sprintf("%d emails opgehaald, %d succesvol opgeslagen", len(emails), savedCount),
		"fetchTime": h.mailFetcher.GetLastFetchTime(),
	})
}

// ListUnprocessedEmails haalt alle onverwerkte emails op
// @Summary Onverwerkte emails ophalen
// @Description Haalt een lijst van alle onverwerkte emails op
// @Tags Mail
// @Accept json
// @Produce json
// @Success 200 {object} PaginatedMailResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/unprocessed [get]
// @Security BearerAuth
func (h *MailHandler) ListUnprocessedEmails(c *fiber.Ctx) error {
	logger.Info("Request ontvangen voor ListUnprocessedEmails")

	// FindUnprocessed haalt alle onverwerkte e-mails op (geen paginering ondersteund door repo methode)
	emails, err := h.incomingEmailRepo.FindUnprocessed(c.Context())
	if err != nil {
		logger.Error("Fout bij ophalen onverwerkte e-mails", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Interne serverfout bij ophalen onverwerkte e-mails"})
	}

	// Map de e-mails naar response DTOs
	responses := make([]*MailResponse, 0, len(emails)) // Initialize as slice of pointers
	for _, email := range emails {
		responses = append(responses, mapEmailToResponse(email)) // Pass the pointer
	}

	logger.Info("Onverwerkte e-mails opgehaald", "count", len(responses))

	return c.JSON(responses)
}

// ListEmailsByAccountType haalt een lijst van emails op per account type met paginering
// @Summary Lijst van emails ophalen per account type
// @Description Haalt een gepagineerde lijst van emails op voor een specifiek account type
// @Tags Mail
// @Accept json
// @Produce json
// @Param type path string true "Account Type (bv. 'info', 'inschrijving')"
// @Param limit query int false "Aantal resultaten per pagina (standaard 10)"
// @Param offset query int false "Offset voor paginering (standaard 0)"
// @Success 200 {object} PaginatedMailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mail/account/{type} [get]
// @Security BearerAuth
func (h *MailHandler) ListEmailsByAccountType(c *fiber.Ctx) error {
	accountType := c.Params("type")
	if accountType == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Account type is verplicht"})
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	if limit < 1 || limit > 100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Limit moet tussen 1 en 100 liggen"})
	}
	if offset < 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Offset mag niet negatief zijn"})
	}

	ctx := c.Context()
	emails, totalCount, err := h.incomingEmailRepo.ListByAccountTypePaginated(ctx, accountType, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen emails per account type", "error", err, "account_type", accountType, "limit", limit, "offset", offset)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Interne serverfout bij ophalen e-mails"})
	}

	responseEmails := mapEmailsToResponse(emails)
	response := PaginatedMailResponse{
		Emails:     responseEmails,
		TotalCount: totalCount,
	}

	// Log the response just before sending
	logger.Info("ListEmailsByAccountType response wordt verzonden", "account_type", accountType, "count", len(responseEmails), "total_count", totalCount)
	// Optional: Log HTML preview for the first email if needed
	// if len(responseEmails) > 0 {
	//  logger.Info("First email HTML preview", "email_id", responseEmails[0].ID, "response_html_preview", getFirstNChars(responseEmails[0].HTML, 100))
	// }

	return c.JSON(response)
}

// Helper function to get first N characters of a string
func getFirstNChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
