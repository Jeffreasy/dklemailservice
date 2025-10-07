package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// NewsletterHandler bevat handlers voor nieuwsbrief beheer
type NewsletterHandler struct {
	newsletterRepo    repository.NewsletterRepository
	newsletterSvc     *services.NewsletterSender
	authService       services.AuthService
	permissionService services.PermissionService
}

// NewNewsletterHandler maakt een nieuwe newsletter handler
func NewNewsletterHandler(
	newsletterRepo repository.NewsletterRepository,
	newsletterSvc *services.NewsletterSender,
	authService services.AuthService,
	permissionService services.PermissionService,
) *NewsletterHandler {
	return &NewsletterHandler{
		newsletterRepo:    newsletterRepo,
		newsletterSvc:     newsletterSvc,
		authService:       authService,
		permissionService: permissionService,
	}
}

// RegisterRoutes registreert de routes voor newsletter beheer
func (h *NewsletterHandler) RegisterRoutes(app *fiber.App) {
	// Groep voor newsletter beheer routes
	newsletterGroup := app.Group("/api/newsletter")
	newsletterGroup.Use(AuthMiddleware(h.authService))

	// Read-only routes (require newsletter read)
	readGroup := newsletterGroup.Group("", PermissionMiddleware(h.permissionService, "newsletter", "read"))
	readGroup.Get("/", h.ListNewsletters)
	readGroup.Get("/:id", h.GetNewsletter)

	// Write routes (require newsletter write)
	writeGroup := newsletterGroup.Group("", PermissionMiddleware(h.permissionService, "newsletter", "write"))
	writeGroup.Post("/", h.CreateNewsletter)
	writeGroup.Put("/:id", h.UpdateNewsletter)

	// Delete routes (require newsletter delete)
	deleteGroup := newsletterGroup.Group("", PermissionMiddleware(h.permissionService, "newsletter", "delete"))
	deleteGroup.Delete("/:id", h.DeleteNewsletter)

	// Send routes (require newsletter send)
	sendGroup := newsletterGroup.Group("", PermissionMiddleware(h.permissionService, "newsletter", "send"))
	sendGroup.Post("/:id/send", h.SendNewsletter)
}

// ListNewsletters haalt een lijst van nieuwsbrieven op
func (h *NewsletterHandler) ListNewsletters(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
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
	newsletters, err := h.newsletterRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Fout bij ophalen nieuwsbrieven", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwsbrieven niet ophalen",
		})
	}

	return c.JSON(newsletters)
}

// CreateNewsletter maakt een nieuwe nieuwsbrief aan
func (h *NewsletterHandler) CreateNewsletter(c *fiber.Ctx) error {
	var req struct {
		Subject string `json:"subject"`
		Content string `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	if req.Subject == "" || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Subject en content zijn verplicht",
		})
	}

	nl := &models.Newsletter{
		Subject: req.Subject,
		Content: req.Content,
	}

	ctx := c.Context()
	if err := h.newsletterRepo.Create(ctx, nl); err != nil {
		logger.Error("Fout bij aanmaken nieuwsbrief", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwsbrief niet aanmaken",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(nl)
}

// GetNewsletter haalt een specifieke nieuwsbrief op
func (h *NewsletterHandler) GetNewsletter(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	ctx := c.Context()
	nl, err := h.newsletterRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen nieuwsbrief", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwsbrief niet ophalen",
		})
	}

	if nl == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Nieuwsbrief niet gevonden",
		})
	}

	return c.JSON(nl)
}

// UpdateNewsletter werkt een nieuwsbrief bij
func (h *NewsletterHandler) UpdateNewsletter(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	var req struct {
		Subject string `json:"subject"`
		Content string `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
		})
	}

	ctx := c.Context()
	nl, err := h.newsletterRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen nieuwsbrief", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwsbrief niet ophalen",
		})
	}

	if nl == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Nieuwsbrief niet gevonden",
		})
	}

	if nl.SentAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nieuwsbrief is al verzonden en kan niet meer worden bijgewerkt",
		})
	}

	if req.Subject != "" {
		nl.Subject = req.Subject
	}
	if req.Content != "" {
		nl.Content = req.Content
	}

	if err := h.newsletterRepo.Update(ctx, nl); err != nil {
		logger.Error("Fout bij bijwerken nieuwsbrief", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwsbrief niet bijwerken",
		})
	}

	return c.JSON(nl)
}

// DeleteNewsletter verwijdert een nieuwsbrief
func (h *NewsletterHandler) DeleteNewsletter(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	ctx := c.Context()
	nl, err := h.newsletterRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen nieuwsbrief", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwsbrief niet ophalen",
		})
	}

	if nl == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Nieuwsbrief niet gevonden",
		})
	}

	if nl.SentAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nieuwsbrief is al verzonden en kan niet meer worden verwijderd",
		})
	}

	if err := h.newsletterRepo.Delete(ctx, id); err != nil {
		logger.Error("Fout bij verwijderen nieuwsbrief", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwsbrief niet verwijderen",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Nieuwsbrief succesvol verwijderd",
	})
}

// SendNewsletter verzendt een nieuwsbrief naar subscribers
func (h *NewsletterHandler) SendNewsletter(c *fiber.Ctx) error {
	logger.Info("SendNewsletter handler called", "path", c.Path(), "method", c.Method())

	id := c.Params("id")
	if id == "" {
		logger.Warn("SendNewsletter: ID is verplicht")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	logger.Info("SendNewsletter: Starting send for newsletter", "id", id)

	ctx := c.Context()
	if err := h.newsletterSvc.SendManual(ctx, id); err != nil {
		logger.Error("SendNewsletter: Fout bij verzenden nieuwsbrief", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon nieuwsbrief niet verzenden",
		})
	}

	logger.Info("SendNewsletter: Successfully initiated newsletter send", "id", id)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Nieuwsbrief wordt verzonden naar subscribers",
	})
}
