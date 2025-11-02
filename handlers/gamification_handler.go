package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GamificationHandler bevat handlers voor gamification functionaliteit
type GamificationHandler struct {
	gamificationService *services.GamificationService
	authService         services.AuthService
	permissionService   services.PermissionService
}

// NewGamificationHandler maakt een nieuwe gamification handler
func NewGamificationHandler(
	gamificationService *services.GamificationService,
	authService services.AuthService,
	permissionService services.PermissionService,
) *GamificationHandler {
	return &GamificationHandler{
		gamificationService: gamificationService,
		authService:         authService,
		permissionService:   permissionService,
	}
}

// RegisterRoutes registreert de routes voor gamification
func (h *GamificationHandler) RegisterRoutes(app *fiber.App) {
	// Public endpoints
	app.Get("/api/leaderboard", h.GetLeaderboard)
	app.Get("/api/badges", h.GetBadges)

	// Participant endpoints (authenticated)
	participantGroup := app.Group("/api/participant", AuthMiddleware(h.authService))
	participantGroup.Get("/achievements", h.GetMyAchievements)
	participantGroup.Get("/rank", h.GetMyRank)
	participantGroup.Get("/:id/achievements", PermissionMiddleware(h.permissionService, "achievements", "read"), h.GetParticipantAchievements)

	// Admin endpoints voor badge management
	adminGroup := app.Group("/api/admin/badges", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "badges", "write"))
	adminGroup.Post("/", h.CreateBadge)
	adminGroup.Get("/", h.GetAllBadgesAdmin)
	adminGroup.Get("/:id", h.GetBadge)
	adminGroup.Put("/:id", h.UpdateBadge)
	adminGroup.Delete("/:id", h.DeleteBadge)

	// Admin endpoints voor achievement management
	achievementAdminGroup := app.Group("/api/admin/achievements", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "achievements", "write"))
	achievementAdminGroup.Post("/award", h.AwardBadge)
	achievementAdminGroup.Delete("/remove", h.RemoveBadge)
}

// ==================== PUBLIC ENDPOINTS ====================

// GetLeaderboard haalt de leaderboard op
func (h *GamificationHandler) GetLeaderboard(c *fiber.Ctx) error {
	ctx := context.Background()

	filters := models.LeaderboardFilters{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 50),
	}

	// Optionele filters
	if route := c.Query("route"); route != "" {
		filters.Route = &route
	}

	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			filters.Year = &year
		}
	}

	if minStepsStr := c.Query("min_steps"); minStepsStr != "" {
		minSteps, err := strconv.Atoi(minStepsStr)
		if err == nil {
			filters.MinSteps = &minSteps
		}
	}

	if topNStr := c.Query("top_n"); topNStr != "" {
		topN, err := strconv.Atoi(topNStr)
		if err == nil {
			filters.TopN = &topN
		}
	}

	leaderboard, err := h.gamificationService.GetLeaderboard(ctx, filters)
	if err != nil {
		logger.Error("Fout bij ophalen leaderboard", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon leaderboard niet ophalen",
		})
	}

	return c.JSON(leaderboard)
}

// GetBadges haalt alle actieve badges op (public)
func (h *GamificationHandler) GetBadges(c *fiber.Ctx) error {
	ctx := context.Background()

	// Check if user is authenticated to show personalized stats
	var participantID *string
	userID, ok := c.Locals("userID").(string)
	if ok && userID != "" {
		participantID = &userID
	}

	badges, err := h.gamificationService.GetBadgesWithStats(ctx, participantID)
	if err != nil {
		logger.Error("Fout bij ophalen badges", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon badges niet ophalen",
		})
	}

	return c.JSON(badges)
}

// ==================== PARTICIPANT ENDPOINTS ====================

// GetMyAchievements haalt achievements op van ingelogde gebruiker
func (h *GamificationHandler) GetMyAchievements(c *fiber.Ctx) error {
	ctx := context.Background()
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	summary, err := h.gamificationService.GetParticipantSummary(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen achievements", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon achievements niet ophalen",
		})
	}

	return c.JSON(summary)
}

// GetMyRank haalt rank informatie op van ingelogde gebruiker
func (h *GamificationHandler) GetMyRank(c *fiber.Ctx) error {
	ctx := context.Background()
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	rank, err := h.gamificationService.GetParticipantRank(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen rank", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon rank informatie niet ophalen",
		})
	}

	return c.JSON(rank)
}

// GetParticipantAchievements haalt achievements op voor specifieke deelnemer (admin/staff)
func (h *GamificationHandler) GetParticipantAchievements(c *fiber.Ctx) error {
	ctx := context.Background()
	participantID := c.Params("id")

	if participantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Deelnemer ID is verplicht",
		})
	}

	summary, err := h.gamificationService.GetParticipantSummary(ctx, participantID)
	if err != nil {
		logger.Error("Fout bij ophalen achievements", "error", err, "participant_id", participantID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon achievements niet ophalen",
		})
	}

	return c.JSON(summary)
}

// ==================== ADMIN BADGE MANAGEMENT ====================

// CreateBadge maakt een nieuwe badge aan (admin only)
func (h *GamificationHandler) CreateBadge(c *fiber.Ctx) error {
	ctx := context.Background()

	var req models.BadgeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Validatie
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Naam is verplicht",
		})
	}

	if req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Beschrijving is verplicht",
		})
	}

	if req.Points < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Punten moet groter of gelijk zijn aan 0",
		})
	}

	badge, err := h.gamificationService.CreateBadge(ctx, req)
	if err != nil {
		logger.Error("Fout bij aanmaken badge", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(badge)
}

// GetAllBadgesAdmin haalt alle badges op inclusief inactieve (admin only)
func (h *GamificationHandler) GetAllBadgesAdmin(c *fiber.Ctx) error {
	ctx := context.Background()

	// Check if we should include inactive badges
	activeOnly := c.QueryBool("active_only", false)

	badges, err := h.gamificationService.GetAllBadges(ctx, activeOnly)
	if err != nil {
		logger.Error("Fout bij ophalen badges", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon badges niet ophalen",
		})
	}

	return c.JSON(badges)
}

// GetBadge haalt één badge op (admin only)
func (h *GamificationHandler) GetBadge(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Badge ID is verplicht",
		})
	}

	badge, err := h.gamificationService.GetBadge(ctx, id)
	if err != nil {
		logger.Error("Fout bij ophalen badge", "error", err, "id", id)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Badge niet gevonden",
		})
	}

	return c.JSON(badge)
}

// UpdateBadge werkt een badge bij (admin only)
func (h *GamificationHandler) UpdateBadge(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Badge ID is verplicht",
		})
	}

	var req models.BadgeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Validatie
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Naam is verplicht",
		})
	}

	if req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Beschrijving is verplicht",
		})
	}

	if req.Points < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Punten moet groter of gelijk zijn aan 0",
		})
	}

	badge, err := h.gamificationService.UpdateBadge(ctx, id, req)
	if err != nil {
		logger.Error("Fout bij bijwerken badge", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(badge)
}

// DeleteBadge verwijdert een badge (admin only)
func (h *GamificationHandler) DeleteBadge(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Badge ID is verplicht",
		})
	}

	err := h.gamificationService.DeleteBadge(ctx, id)
	if err != nil {
		logger.Error("Fout bij verwijderen badge", "error", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Badge succesvol verwijderd",
	})
}

// ==================== ADMIN ACHIEVEMENT MANAGEMENT ====================

// AwardBadge kent een badge toe aan een deelnemer (admin/staff)
func (h *GamificationHandler) AwardBadge(c *fiber.Ctx) error {
	ctx := context.Background()

	var req models.AchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	if req.ParticipantID == "" || req.BadgeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Deelnemer ID en Badge ID zijn verplicht",
		})
	}

	achievement, err := h.gamificationService.AwardBadge(ctx, req.ParticipantID, req.BadgeID)
	if err != nil {
		logger.Error("Fout bij toekennen badge", "error", err, "participant_id", req.ParticipantID, "badge_id", req.BadgeID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(achievement)
}

// RemoveBadge verwijdert een badge van een deelnemer (admin/staff)
func (h *GamificationHandler) RemoveBadge(c *fiber.Ctx) error {
	ctx := context.Background()

	var req models.AchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	if req.ParticipantID == "" || req.BadgeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Deelnemer ID en Badge ID zijn verplicht",
		})
	}

	err := h.gamificationService.RemoveBadge(ctx, req.ParticipantID, req.BadgeID)
	if err != nil {
		logger.Error("Fout bij verwijderen achievement", "error", err, "participant_id", req.ParticipantID, "badge_id", req.BadgeID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Achievement succesvol verwijderd",
	})
}
