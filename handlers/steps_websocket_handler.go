package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/services"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// StepsWebSocketHandler handles WebSocket connecties voor stappen updates
type StepsWebSocketHandler struct {
	stepsHub    *services.StepsHub
	authService services.AuthService
}

// NewStepsWebSocketHandler creates een nieuwe StepsWebSocketHandler
func NewStepsWebSocketHandler(
	stepsHub *services.StepsHub,
	authService services.AuthService,
) *StepsWebSocketHandler {
	return &StepsWebSocketHandler{
		stepsHub:    stepsHub,
		authService: authService,
	}
}

// RegisterRoutes registreert WebSocket routes
func (h *StepsWebSocketHandler) RegisterRoutes(app *fiber.App) {
	// WebSocket upgrade check middleware
	app.Use("/api/ws/steps", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			// Optioneel: JWT validatie voor WebSocket
			// Token kan komen van query parameter of header
			token := c.Query("token")
			if token == "" {
				authHeader := c.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			// Als token aanwezig is, valideer het (maar sta ook anonymous toe)
			if token != "" {
				userID, err := h.authService.ValidateToken(token)
				if err != nil {
					logger.Error("WebSocket auth failed", "error", err, "token_length", len(token))
					// Don't reject - allow anonymous connections voor public leaderboard
					// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					// 	"error": "Invalid token",
					// })
				} else {
					// Store user info in context voor gebruik in handler
					c.Locals("userID", userID)
				}
			}

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint
	app.Get("/api/ws/steps", websocket.New(h.HandleWebSocket))

	logger.Info("WebSocket endpoint registered", "path", "/api/ws/steps")
}

// HandleWebSocket handles de WebSocket connectie
func (h *StepsWebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	// Extract user info from query params
	userID := c.Query("user_id")
	participantID := c.Query("participant_id")

	// Als userID niet in query, probeer uit context (van auth middleware)
	if userID == "" {
		if uid := c.Locals("userID"); uid != nil {
			if uidStr, ok := uid.(string); ok {
				userID = uidStr
			}
		}
	}

	logger.Info("WebSocket client connecting",
		"user_id", userID,
		"participant_id", participantID,
		"remote_addr", c.RemoteAddr().String(),
	)

	// Create client
	client := &services.StepsClient{
		Hub:           h.stepsHub,
		Conn:          c,
		Send:          make(chan []byte, 256),
		UserID:        userID,
		ParticipantID: participantID,
		Subscriptions: make(map[string]bool),
	}

	// Register client
	h.stepsHub.Register <- client

	// Log connection count
	logger.Info("WebSocket client connected",
		"user_id", userID,
		"participant_id", participantID,
		"total_clients", h.stepsHub.GetClientCount(),
	)

	// Send welcome message met instructies
	welcomeMsg := map[string]interface{}{
		"type":               "welcome",
		"message":            "Connected to StepsHub! Send {\"type\":\"subscribe\",\"channels\":[\"step_updates\",\"total_updates\",\"leaderboard_updates\"]} to receive updates",
		"available_channels": []string{"step_updates", "total_updates", "leaderboard_updates", "badge_earned"},
		"timestamp":          time.Now().Unix(),
	}
	if welcomeBytes, err := json.Marshal(welcomeMsg); err == nil {
		client.Send <- welcomeBytes
	}

	// Start pumps
	go client.WritePump()
	client.ReadPump()

	// Client disconnected
	logger.Info("WebSocket client disconnected",
		"user_id", userID,
		"total_clients", h.stepsHub.GetClientCount(),
	)
}

// GetStats returns WebSocket statistics (admin only)
func (h *StepsWebSocketHandler) GetStats(c *fiber.Ctx) error {
	stats := fiber.Map{
		"total_clients": h.stepsHub.GetClientCount(),
		"subscriptions": h.stepsHub.GetSubscriptionCount(),
	}
	return c.JSON(stats)
}
