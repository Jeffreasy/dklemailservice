package services

import (
	"encoding/json"
	"time"

	"github.com/gofiber/websocket/v2"
)

// MessageType definieert het type WebSocket bericht
type MessageType string

const (
	MessageTypeStepUpdate        MessageType = "step_update"
	MessageTypeTotalUpdate       MessageType = "total_update"
	MessageTypeLeaderboardUpdate MessageType = "leaderboard_update"
	MessageTypeBadgeEarned       MessageType = "badge_earned"
	MessageTypeSubscribe         MessageType = "subscribe"
	MessageTypeUnsubscribe       MessageType = "unsubscribe"
	MessageTypePing              MessageType = "ping"
	MessageTypePong              MessageType = "pong"
)

// StepUpdateMessage - individuele deelnemer stappen update
type StepUpdateMessage struct {
	Type           MessageType `json:"type"`
	ParticipantID  string      `json:"participant_id"`
	Naam           string      `json:"naam"`
	Steps          int         `json:"steps"`
	Delta          int         `json:"delta"`
	Route          string      `json:"route"`
	AllocatedFunds int         `json:"allocated_funds"`
	Timestamp      int64       `json:"timestamp"`
}

// TotalUpdateMessage - totaal stappen update
type TotalUpdateMessage struct {
	Type       MessageType `json:"type"`
	TotalSteps int         `json:"total_steps"`
	Year       int         `json:"year"`
	Timestamp  int64       `json:"timestamp"`
}

// LeaderboardUpdateMessage - leaderboard update
type LeaderboardUpdateMessage struct {
	Type      MessageType        `json:"type"`
	TopN      int                `json:"top_n"`
	Entries   []LeaderboardEntry `json:"entries"`
	Timestamp int64              `json:"timestamp"`
}

// LeaderboardEntry - een entry in de leaderboard
type LeaderboardEntry struct {
	Rank              int    `json:"rank"`
	ParticipantID     string `json:"participant_id"`
	Naam              string `json:"naam"`
	Steps             int    `json:"steps"`
	AchievementPoints int    `json:"achievement_points"`
	TotalScore        int    `json:"total_score"`
	Route             string `json:"route"`
	BadgeCount        int    `json:"badge_count"`
}

// BadgeEarnedMessage - badge verdiend notificatie
type BadgeEarnedMessage struct {
	Type          MessageType `json:"type"`
	ParticipantID string      `json:"participant_id"`
	BadgeName     string      `json:"badge_name"`
	BadgeIcon     string      `json:"badge_icon"`
	Points        int         `json:"points"`
	Timestamp     int64       `json:"timestamp"`
}

// StepsHub beheert WebSocket connecties voor stappen updates
type StepsHub struct {
	// Registered clients
	Clients map[*StepsClient]bool

	// Broadcast channels
	StepUpdate        chan *StepUpdateMessage
	TotalUpdate       chan *TotalUpdateMessage
	LeaderboardUpdate chan *LeaderboardUpdateMessage
	BadgeEarned       chan *BadgeEarnedMessage

	// Registration channels
	Register   chan *StepsClient
	Unregister chan *StepsClient

	// Services (optional, voor future use)
	StepsService        *StepsService
	GamificationService interface{}
}

// StepsClient represents een WebSocket client
type StepsClient struct {
	Hub           *StepsHub
	Conn          *websocket.Conn
	Send          chan []byte
	UserID        string
	ParticipantID string
	Subscriptions map[string]bool // Subscription types
}

// NewStepsHub creates een nieuwe StepsHub
func NewStepsHub(stepsService *StepsService, gamificationService interface{}) *StepsHub {
	return &StepsHub{
		Clients:             make(map[*StepsClient]bool),
		StepUpdate:          make(chan *StepUpdateMessage, 256),
		TotalUpdate:         make(chan *TotalUpdateMessage, 256),
		LeaderboardUpdate:   make(chan *LeaderboardUpdateMessage, 256),
		BadgeEarned:         make(chan *BadgeEarnedMessage, 256),
		Register:            make(chan *StepsClient),
		Unregister:          make(chan *StepsClient),
		StepsService:        stepsService,
		GamificationService: gamificationService,
	}
}

// Run starts de hub event loop
func (h *StepsHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

		case update := <-h.StepUpdate:
			h.broadcastStepUpdate(update)

		case update := <-h.TotalUpdate:
			h.broadcastTotalUpdate(update)

		case update := <-h.LeaderboardUpdate:
			h.broadcastLeaderboardUpdate(update)

		case update := <-h.BadgeEarned:
			h.broadcastBadgeEarned(update)
		}
	}
}

// broadcastStepUpdate broadcast een stappen update naar geabonneerde clients
func (h *StepsHub) broadcastStepUpdate(update *StepUpdateMessage) {
	message, err := json.Marshal(update)
	if err != nil {
		return
	}

	for client := range h.Clients {
		// Check if client is subscribed to step_updates
		if client.Subscriptions["step_updates"] {
			// Ook checken of dit de eigen stappen zijn voor participant
			if client.ParticipantID == "" || client.ParticipantID == update.ParticipantID {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

// broadcastTotalUpdate broadcast totaal stappen update
func (h *StepsHub) broadcastTotalUpdate(update *TotalUpdateMessage) {
	message, err := json.Marshal(update)
	if err != nil {
		return
	}

	for client := range h.Clients {
		if client.Subscriptions["total_updates"] {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
}

// broadcastLeaderboardUpdate broadcast leaderboard update
func (h *StepsHub) broadcastLeaderboardUpdate(update *LeaderboardUpdateMessage) {
	message, err := json.Marshal(update)
	if err != nil {
		return
	}

	for client := range h.Clients {
		if client.Subscriptions["leaderboard_updates"] {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
}

// broadcastBadgeEarned broadcast badge earned notificatie
func (h *StepsHub) broadcastBadgeEarned(update *BadgeEarnedMessage) {
	message, err := json.Marshal(update)
	if err != nil {
		return
	}

	for client := range h.Clients {
		// Stuur badge earned alleen naar de specifieke participant
		if client.ParticipantID == update.ParticipantID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
}

// ReadPump pumps messages van de websocket connectie naar de hub
func (c *StepsClient) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		msgType, ok := msg["type"].(string)
		if !ok {
			continue
		}

		switch MessageType(msgType) {
		case MessageTypeSubscribe:
			c.handleSubscribe(msg)
		case MessageTypeUnsubscribe:
			c.handleUnsubscribe(msg)
		case MessageTypePing:
			c.handlePing()
		}
	}
}

// WritePump pumps messages van de hub naar de websocket connectie
func (c *StepsClient) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// Hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			ping := map[string]interface{}{
				"type":      MessageTypePing,
				"timestamp": time.Now().Unix(),
			}
			message, _ := json.Marshal(ping)
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

// handleSubscribe handles subscription requests van client
func (c *StepsClient) handleSubscribe(msg map[string]interface{}) {
	channels, ok := msg["channels"].([]interface{})
	if !ok {
		return
	}

	for _, ch := range channels {
		if channel, ok := ch.(string); ok {
			c.Subscriptions[channel] = true
		}
	}
}

// handleUnsubscribe handles unsubscription requests van client
func (c *StepsClient) handleUnsubscribe(msg map[string]interface{}) {
	channels, ok := msg["channels"].([]interface{})
	if !ok {
		return
	}

	for _, ch := range channels {
		if channel, ok := ch.(string); ok {
			delete(c.Subscriptions, channel)
		}
	}
}

// handlePing responds to ping with pong
func (c *StepsClient) handlePing() {
	pong := map[string]interface{}{
		"type":      MessageTypePong,
		"timestamp": time.Now().Unix(),
	}
	message, _ := json.Marshal(pong)
	c.Send <- message
}

// GetClientCount returns het aantal actieve clients
func (h *StepsHub) GetClientCount() int {
	return len(h.Clients)
}

// GetSubscriptionCount returns aantal clients per subscription type
func (h *StepsHub) GetSubscriptionCount() map[string]int {
	counts := make(map[string]int)
	for client := range h.Clients {
		for subscription := range client.Subscriptions {
			counts[subscription]++
		}
	}
	return counts
}
