package handlers

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	chatService services.ChatService
	authService services.AuthService
	hub         *services.Hub // global, if needed
	mutex       sync.Mutex
	channelHubs map[string]*services.Hub
}

// NewChatHandler creates a new ChatHandler
func NewChatHandler(chatService services.ChatService, authService services.AuthService, hub *services.Hub) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		authService: authService,
		hub:         hub,
		channelHubs: make(map[string]*services.Hub),
	}
}

func (h *ChatHandler) getChannelHub(channelID string) *services.Hub {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if hub, ok := h.channelHubs[channelID]; ok {
		return hub
	}
	hub := services.NewHub(h.chatService)
	go hub.Run()
	h.channelHubs[channelID] = hub
	return hub
}

// RegisterRoutes registers the chat routes
func (h *ChatHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/chat", h.AuthMiddleware)

	// Channels
	api.Get("/channels", h.ListChannels)
	api.Get("/channels/:id/participants", h.ListParticipants)
	api.Get("/public-channels", h.ListPublicChannels)
	api.Post("/direct", h.CreateDirectChannel)
	api.Post("/channels", h.CreateChannel)
	api.Post("/channels/:id/join", h.JoinChannel)
	api.Post("/channels/:id/leave", h.LeaveChannel)

	// Users for direct chat
	api.Get("/users", h.ListUsers)

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
	api.Get("/ws/:channel_id", websocket.New(h.handleWebsocket))
	api.Get("/ws", websocket.New(h.handleWebsocket))
}

// AuthMiddleware checks for valid JWT
func (h *ChatHandler) handleWebsocket(c *websocket.Conn) {
	channelID := c.Params("channel_id")
	userID := c.Locals("userID").(string)

	// Check if user is participant
	role, err := h.chatService.GetParticipantRole(context.Background(), channelID, userID)
	if err != nil || role == "" {
		c.Close()
		return
	}

	client := &services.Client{
		Hub:    h.getChannelHub(channelID),
		Conn:   c,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	client.Hub.Register <- client

	go client.WritePump()
	client.ReadPump()
}

func (h *ChatHandler) AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		tokenString = c.Query("token")
	}
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
func (h *ChatHandler) CreateDirectChannel(c *fiber.Ctx) error {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userID := c.Locals("userID").(string)

	// Check if direct channel already exists
	channels, err := h.chatService.ListChannelsForUser(c.Context(), userID, 100, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	for _, channel := range channels {
		if channel.Type == "direct" {
			participants, err := h.chatService.ListParticipantsByChannel(c.Context(), channel.ID)
			if err != nil {
				continue
			}
			if len(participants) == 2 && (participants[0].UserID == req.UserID || participants[1].UserID == req.UserID) {
				return c.JSON(channel)
			}
		}
	}

	// Get user names
	currentUser, err := h.authService.GetUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get current user"})
	}
	targetUser, err := h.authService.GetUser(c.Context(), req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get target user"})
	}

	// Create new
	channel := &models.ChatChannel{
		Name:      fmt.Sprintf("Chat between %s and %s", currentUser.Naam, targetUser.Naam),
		Type:      "direct",
		CreatedBy: userID,
	}
	err = h.chatService.CreateChannel(c.Context(), channel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Add users
	err = h.chatService.AddParticipant(c.Context(), &models.ChatChannelParticipant{
		ChannelID: channel.ID,
		UserID:    userID,
		Role:      "member",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.chatService.AddParticipant(c.Context(), &models.ChatChannelParticipant{
		ChannelID: channel.ID,
		UserID:    req.UserID,
		Role:      "member",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(channel)
}

// ListChannels lists the user's channels
func (h *ChatHandler) ListChannels(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	channels, err := h.chatService.ListChannelsForUser(c.Context(), userID, limit, offset)
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

	participant := &models.ChatChannelParticipant{
		ChannelID: channel.ID,
		UserID:    userID,
		Role:      "owner",
	}
	err = h.chatService.AddParticipant(c.Context(), participant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add creator as participant: " + err.Error()})
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

	userID := c.Locals("userID").(string)
	role, err := h.chatService.GetParticipantRole(c.Context(), message.ChannelID, userID)
	if err != nil || (role != "owner" && role != "admin" && message.UserID != userID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized to edit this message"})
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

	message, err := h.chatService.GetMessage(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Message not found"})
	}

	userID := c.Locals("userID").(string)
	role, err := h.chatService.GetParticipantRole(c.Context(), message.ChannelID, userID)
	if err != nil || (role != "owner" && role != "admin" && message.UserID != userID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized to delete this message"})
	}

	err = h.chatService.DeleteMessage(c.Context(), id)
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
func (h *ChatHandler) ListPublicChannels(c *fiber.Ctx) error {
	channels, err := h.chatService.ListPublicChannels(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(channels)
}

func (h *ChatHandler) ListOnlineUsers(c *fiber.Ctx) error {
	users, err := h.chatService.ListOnlineUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// ListParticipants lists participants of a channel
func (h *ChatHandler) ListParticipants(c *fiber.Ctx) error {
	channelID := c.Params("id")
	userID := c.Locals("userID").(string)

	// Check if user is participant
	role, err := h.chatService.GetParticipantRole(c.Context(), channelID, userID)
	if err != nil || role == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not a participant"})
	}

	participants, err := h.chatService.ListParticipantsByChannel(c.Context(), channelID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(participants)
}

// ListUsers lists users for direct chat selection
func (h *ChatHandler) ListUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, err := h.authService.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
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
