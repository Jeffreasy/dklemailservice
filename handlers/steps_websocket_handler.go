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
	// WebSocket upgrade handler with authentication
	wsHandler := func(c *fiber.Ctx) error {
		// Check if this is a WebSocket upgrade request
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		// Token kan komen van query parameter of header
		token := c.Query("token")

		// Trim whitespace en check of token niet leeg is
		token = strings.TrimSpace(token)

		// Als token leeg is in query, probeer Authorization header
		if token == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			}
		}

		// Validate token and store userID for the handler
		var authenticatedUserID string
		if token != "" {
			userID, err := h.authService.ValidateToken(token)
			if err != nil {
				logger.Warn("WebSocket auth failed, allowing anonymous access",
					"error", err,
					"token_length", len(token),
					"remote_addr", c.IP())
				// Allow anonymous connections for public leaderboard/stats
			} else {
				authenticatedUserID = userID
				logger.Debug("WebSocket authenticated user", "user_id", userID)
			}
		} else {
			logger.Debug("WebSocket anonymous connection", "remote_addr", c.IP())
		}

		// Upgrade to WebSocket with authentication context
		return websocket.New(func(conn *websocket.Conn) {
			h.handleWebSocketConnection(conn, authenticatedUserID)
		})(c)
	}

	// WebSocket endpoint (primary route)
	app.Get("/api/ws/steps", wsHandler)

	// Public alias for frontend compatibility (without /api prefix)
	app.Get("/ws/steps", wsHandler)

	logger.Info("WebSocket endpoints registered", "paths", "/api/ws/steps, /ws/steps (both support anonymous & authenticated)")
}

// handleWebSocketConnection handles de WebSocket connectie met authenticated user context
func (h *StepsWebSocketHandler) handleWebSocketConnection(c *websocket.Conn, authenticatedUserID string) {
	// Extract user info from query params
	userID := c.Query("user_id")
	participantID := c.Query("participant_id")

	// Als userID niet in query params, gebruik authenticated user ID from token
	if userID == "" {
		userID = authenticatedUserID
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
