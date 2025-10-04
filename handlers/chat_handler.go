package handlers

import (
	"dklautomationgo/models"
	"dklautomationgo/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	chatService services.ChatService
	authService services.AuthService
}

// NewChatHandler creates a new ChatHandler
func NewChatHandler(chatService services.ChatService, authService services.AuthService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		authService: authService,
	}
}

// RegisterRoutes registers the chat routes
func (h *ChatHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/chat", h.AuthMiddleware)

	// Channels
	api.Get("/channels", h.ListChannels)
	api.Post("/channels", h.CreateChannel)
	api.Post("/channels/:id/join", h.JoinChannel)
	api.Post("/channels/:id/leave", h.LeaveChannel)

	// Messages
	api.Get("/channels/:channel_id/messages", h.GetMessages)
	api.Post("/channels/:channel_id/messages", h.SendMessage)
	api.Put("/messages/:id", h.EditMessage)
	api.Delete("/messages/:id", h.DeleteMessage)
	api.Post("/messages/:id/reactions", h.AddReaction)
	api.Delete("/messages/:id/reactions/:emoji", h.RemoveReaction)

	// Presence
	api.Put("/presence", h.UpdatePresence)
	api.Get("/online-users", h.ListOnlineUsers)

	// Typing
	api.Post("/channels/:channel_id/typing/start", h.StartTyping)
	api.Post("/channels/:channel_id/typing/stop", h.StopTyping)
	api.Get("/channels/:channel_id/typing", h.GetTypingUsers)

	// Utility
	api.Post("/channels/:id/read", h.MarkAsRead)
	api.Get("/unread", h.GetUnreadCount)
}

// AuthMiddleware checks for valid JWT
func (h *ChatHandler) AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization token"})
	}

	userID, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	c.Locals("userID", userID)
	return c.Next()
}

// ListChannels lists the user's channels
func (h *ChatHandler) ListChannels(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	channels, err := h.chatService.ListChannels(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(channels)
}

// CreateChannel creates a new channel
func (h *ChatHandler) CreateChannel(c *fiber.Ctx) error {
	var channel models.ChatChannel
	if err := c.BodyParser(&channel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userID := c.Locals("userID").(string)
	channel.CreatedBy = userID

	err := h.chatService.CreateChannel(c.Context(), &channel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(channel)
}

// JoinChannel adds the user to a channel
func (h *ChatHandler) JoinChannel(c *fiber.Ctx) error {
	channelID := c.Params("id")
	userID := c.Locals("userID").(string)

	participant := &models.ChatChannelParticipant{
		ChannelID: channelID,
		UserID:    userID,
		Role:      "member",
	}

	err := h.chatService.AddParticipant(c.Context(), participant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

// LeaveChannel removes the user from a channel
func (h *ChatHandler) LeaveChannel(c *fiber.Ctx) error {
	channelID := c.Params("id")
	userID := c.Locals("userID").(string)

	// Find the participant ID
	participants, err := h.chatService.ListParticipantsByChannel(c.Context(), channelID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var participantID string
	for _, p := range participants {
		if p.UserID == userID {
			participantID = p.ID
			break
		}
	}

	if participantID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Participant not found"})
	}

	err = h.chatService.DeleteParticipant(c.Context(), participantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

// GetMessages gets messages for a channel
func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	channelID := c.Params("channel_id")
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	messages, err := h.chatService.ListMessagesByChannel(c.Context(), channelID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(messages)
}

// SendMessage sends a new message
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	channelID := c.Params("channel_id")
	userID := c.Locals("userID").(string)

	var message models.ChatMessage
	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	message.ChannelID = channelID
	message.UserID = userID

	err := h.chatService.CreateMessage(c.Context(), &message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(message)
}

// EditMessage edits a message
func (h *ChatHandler) EditMessage(c *fiber.Ctx) error {
	id := c.Params("id")

	var update struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	message, err := h.chatService.GetMessage(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Message not found"})
	}

	message.Content = update.Content
	message.EditedAt = time.Now()

	err = h.chatService.UpdateMessage(c.Context(), message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(message)
}

// DeleteMessage deletes a message
func (h *ChatHandler) DeleteMessage(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.chatService.DeleteMessage(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

// AddReaction adds a reaction to a message
func (h *ChatHandler) AddReaction(c *fiber.Ctx) error {
	messageID := c.Params("id")
	userID := c.Locals("userID").(string)

	var update struct {
		Emoji string `json:"emoji"`
	}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	reaction := &models.ChatMessageReaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     update.Emoji,
	}

	err := h.chatService.AddReaction(c.Context(), reaction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(reaction)
}

// RemoveReaction removes a reaction from a message
func (h *ChatHandler) RemoveReaction(c *fiber.Ctx) error {
	messageID := c.Params("id")
	emoji := c.Params("emoji")
	userID := c.Locals("userID").(string)

	reactions, err := h.chatService.ListReactionsByMessage(c.Context(), messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var reactionID string
	for _, r := range reactions {
		if r.UserID == userID && r.Emoji == emoji {
			reactionID = r.ID
			break
		}
	}

	if reactionID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reaction not found"})
	}

	err = h.chatService.DeleteReaction(c.Context(), reactionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

// UpdatePresence updates user presence status
func (h *ChatHandler) UpdatePresence(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var update struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	presence := &models.ChatUserPresence{
		UserID:   userID,
		Status:   update.Status,
		LastSeen: time.Now(),
	}

	err := h.chatService.UpdatePresence(c.Context(), presence)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presence)
}

// ListOnlineUsers lists online users
func (h *ChatHandler) ListOnlineUsers(c *fiber.Ctx) error {
	// This requires listing all presences with status 'online'
	// But the repository doesn't have List method, so add if needed or implement in service
	// For now, return empty
	return c.JSON([]string{})
}

// StartTyping, StopTyping, GetTypingUsers would require additional logic for typing indicators, perhaps using memory or DB
// MarkAsRead, GetUnreadCount similarly need tracking of read status, which is not in schema, so would need additional tables or logic

// For now, implement stubs or leave as TODO
func (h *ChatHandler) StartTyping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true})
}

func (h *ChatHandler) StopTyping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true})
}

func (h *ChatHandler) GetTypingUsers(c *fiber.Ctx) error {
	return c.JSON([]string{})
}

func (h *ChatHandler) MarkAsRead(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true})
}

func (h *ChatHandler) GetUnreadCount(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"unread": 0})
}
