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

// NotulenWebSocketHandler handles WebSocket connections for notulen real-time updates
type NotulenWebSocketHandler struct {
	notulenHub  *services.NotulenHub
	authService services.AuthService
}

// NewNotulenWebSocketHandler creates a new NotulenWebSocketHandler
func NewNotulenWebSocketHandler(
	notulenHub *services.NotulenHub,
	authService services.AuthService,
) *NotulenWebSocketHandler {
	return &NotulenWebSocketHandler{
		notulenHub:  notulenHub,
		authService: authService,
	}
}

// RegisterRoutes registers WebSocket routes
func (h *NotulenWebSocketHandler) RegisterRoutes(app *fiber.App) {
	// WebSocket upgrade check middleware
	app.Use("/api/ws/notulen", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			// Optional: JWT validation for WebSocket
			// Token can come from query parameter or header
			token := c.Query("token")
			if token == "" {
				authHeader := c.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			// If token is present, validate it (but allow anonymous connections for public notulen)
			if token != "" {
				userID, err := h.authService.ValidateToken(token)
				if err != nil {
					logger.Error("WebSocket auth failed", "error", err, "token_length", len(token))
					// Don't reject - allow anonymous connections for public notulen
					// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					// 	"error": "Invalid token",
					// })
				} else {
					// Store user info in context for use in handler
					c.Locals("userID", userID)
				}
			}

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint
	app.Get("/api/ws/notulen", websocket.New(h.HandleWebSocket))

	logger.Info("Notulen WebSocket endpoint registered", "path", "/api/ws/notulen")
}

// HandleWebSocket handles the WebSocket connection
func (h *NotulenWebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	// Extract user info from query params
	userID := c.Query("user_id")
	notulenID := c.Query("notulen_id")

	// If userID not in query, try from context (from auth middleware)
	if userID == "" {
		if uid := c.Locals("userID"); uid != nil {
			if uidStr, ok := uid.(string); ok {
				userID = uidStr
			}
		}
	}

	logger.Info("Notulen WebSocket client connecting",
		"user_id", userID,
		"notulen_id", notulenID,
		"remote_addr", c.RemoteAddr().String(),
	)

	// Create client
	client := &services.NotulenWebSocketClient{
		ID:        generateClientID(),
		UserID:    userID,
		NotulenID: notulenID,
		Conn:      c,
		Send:      make(chan []byte, 256),
	}

	// Register client
	h.notulenHub.RegisterClient(client)

	// Log connection count
	logger.Info("Notulen WebSocket client connected",
		"user_id", userID,
		"notulen_id", notulenID,
		"total_clients", h.notulenHub.GetClientCount(),
	)

	// Send welcome message with instructions
	welcomeMsg := map[string]interface{}{
		"type":       "welcome",
		"message":    "Connected to NotulenHub! You will receive real-time updates for notulen changes.",
		"user_id":    userID,
		"notulen_id": notulenID,
		"timestamp":  time.Now().Unix(),
	}
	if welcomeBytes, err := json.Marshal(welcomeMsg); err == nil {
		client.Send <- welcomeBytes
	}

	// Start pumps
	go client.WritePump()
	client.ReadPump()

	// Client disconnected
	logger.Info("Notulen WebSocket client disconnected",
		"user_id", userID,
		"notulen_id", notulenID,
		"total_clients", h.notulenHub.GetClientCount(),
	)
}

// GetStats returns WebSocket statistics (admin only)
func (h *NotulenWebSocketHandler) GetStats(c *fiber.Ctx) error {
	stats := fiber.Map{
		"total_clients": h.notulenHub.GetClientCount(),
	}
	return c.JSON(stats)
}

// generateClientID generates a unique client ID
func generateClientID() string {
	return "notulen_ws_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// randomString generates a random string of given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
