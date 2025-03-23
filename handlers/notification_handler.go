package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"

	"github.com/gofiber/fiber/v2"
)

// NotificationHandler struct voor het afhandelen van notificatie-gerelateerde routes
type NotificationHandler struct {
	notificationRepo    repository.NotificationRepository
	notificationService services.NotificationService
	authService         services.AuthService
}

// NewNotificationHandler maakt een nieuwe NotificationHandler
func NewNotificationHandler(
	notificationRepo repository.NotificationRepository,
	notificationService services.NotificationService,
	authService services.AuthService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo:    notificationRepo,
		notificationService: notificationService,
		authService:         authService,
	}
}

// RegisterRoutes registreert de routes voor de NotificationHandler
func (h *NotificationHandler) RegisterRoutes(app *fiber.App) {
	// Groepeer routes onder /api/v1/notifications met auth middleware
	notificationGroup := app.Group("/api/v1/notifications", h.authMiddleware)

	// Route definities
	notificationGroup.Get("/", h.ListNotifications)
	notificationGroup.Post("/", h.CreateNotification)
	notificationGroup.Get("/:id", h.GetNotification)
	notificationGroup.Delete("/:id", h.DeleteNotification)
	notificationGroup.Post("/reprocess-all", h.ReprocessAllNotifications)
}

// authMiddleware controleert of de gebruiker geauthenticeerd is
func (h *NotificationHandler) authMiddleware(c *fiber.Ctx) error {
	// Haal token op uit Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Strip "Bearer " prefix
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}

	// Valideer token
	if _, err := h.authService.ValidateToken(token); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	return c.Next()
}

// ListNotifications haalt alle notificaties op
func (h *NotificationHandler) ListNotifications(c *fiber.Ctx) error {
	// Parse query parameters voor filters
	notificationType := c.Query("type")
	priority := c.Query("priority")

	ctx := c.Context()
	var notifications []*models.Notification
	var err error

	// Filter op basis van query parameters
	if notificationType != "" {
		notifications, err = h.notificationRepo.ListByType(ctx, models.NotificationType(notificationType))
	} else if priority != "" {
		notifications, err = h.notificationRepo.ListByPriority(ctx, models.NotificationPriority(priority))
	} else {
		// Standaard alle niet verzonden notificaties
		notifications, err = h.notificationRepo.ListUnsent(ctx)
	}

	if err != nil {
		logger.Error("Fout bij ophalen notificaties", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het ophalen van notificaties",
		})
	}

	return c.Status(fiber.StatusOK).JSON(notifications)
}

// CreateNotification maakt een nieuwe notificatie aan
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	// Parse request body
	var request struct {
		Type     string `json:"type"`
		Priority string `json:"priority"`
		Title    string `json:"title"`
		Message  string `json:"message"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request body",
		})
	}

	// Valideer request
	if request.Type == "" || request.Priority == "" || request.Title == "" || request.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Type, priority, title en message zijn verplicht",
		})
	}

	// Controleer of notificationService beschikbaar is
	if h.notificationService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Notification service is niet beschikbaar",
		})
	}

	// Maak een nieuwe notificatie aan
	notification, err := h.notificationService.CreateNotification(
		c.Context(),
		models.NotificationType(request.Type),
		models.NotificationPriority(request.Priority),
		request.Title,
		request.Message,
	)

	if err != nil {
		logger.Error("Fout bij aanmaken notificatie", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het aanmaken van de notificatie",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(notification)
}

// GetNotification haalt een specifieke notificatie op
func (h *NotificationHandler) GetNotification(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	notification, err := h.notificationRepo.GetByID(c.Context(), id)
	if err != nil {
		logger.Error("Fout bij ophalen notificatie", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het ophalen van de notificatie",
		})
	}

	if notification == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notificatie niet gevonden",
		})
	}

	return c.Status(fiber.StatusOK).JSON(notification)
}

// DeleteNotification verwijdert een notificatie
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is verplicht",
		})
	}

	// Controleer of de notificatie bestaat
	notification, err := h.notificationRepo.GetByID(c.Context(), id)
	if err != nil {
		logger.Error("Fout bij ophalen notificatie", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het ophalen van de notificatie",
		})
	}

	if notification == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notificatie niet gevonden",
		})
	}

	// Verwijder de notificatie
	if err := h.notificationRepo.Delete(c.Context(), id); err != nil {
		logger.Error("Fout bij verwijderen notificatie", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het verwijderen van de notificatie",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Notificatie succesvol verwijderd",
	})
}

// ReprocessAllNotifications markeert alle notificaties als niet-verzonden en verwerkt ze opnieuw
func (h *NotificationHandler) ReprocessAllNotifications(c *fiber.Ctx) error {
	// Controleer of notificationService beschikbaar is
	if h.notificationService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Notification service is niet beschikbaar",
		})
	}

	// Haal alle notificaties op uit de database (een nieuwe repository functie nodig)
	ctx := c.Context()

	// Haal notificaties op per type om de database niet te overbelasten
	types := []models.NotificationType{
		models.NotificationTypeContact,
		models.NotificationTypeAanmelding,
		models.NotificationTypeAuth,
		models.NotificationTypeSystem,
		models.NotificationTypeHealth,
	}

	totalCount := 0
	reprocessedCount := 0

	for _, notificationType := range types {
		notifications, err := h.notificationRepo.ListByType(ctx, notificationType)
		if err != nil {
			logger.Error("Fout bij ophalen notificaties voor herverwerking",
				"error", err,
				"type", notificationType)
			continue
		}

		totalCount += len(notifications)

		// Markeer elke notificatie als niet verzonden en update in de database
		for _, notification := range notifications {
			// Alleen verwerk notificaties die prioriteit medium of hoger hebben
			if notification.Priority == models.NotificationPriorityLow {
				continue
			}

			oldSentStatus := notification.Sent
			notification.Sent = false
			notification.SentAt = nil

			if err := h.notificationRepo.Update(ctx, notification); err != nil {
				logger.Error("Fout bij markeren notificatie als niet verzonden",
					"error", err,
					"id", notification.ID)
				continue
			}

			// Verwerk de notificatie direct
			if err := h.notificationService.SendNotification(ctx, notification); err != nil {
				logger.Error("Fout bij opnieuw verzenden notificatie",
					"error", err,
					"id", notification.ID)
				continue
			}

			if oldSentStatus {
				reprocessedCount++
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":     true,
		"message":     "Notificaties opnieuw verwerkt",
		"total":       totalCount,
		"reprocessed": reprocessedCount,
	})
}
