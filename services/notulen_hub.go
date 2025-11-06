package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// NotulenWebSocketEvent represents a WebSocket event for notulen
type NotulenWebSocketEvent struct {
	Type      string      `json:"type"`
	NotulenID string      `json:"notulenId,omitempty"`
	UserID    string      `json:"userId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// NotulenWebSocketClient represents a connected WebSocket client
type NotulenWebSocketClient struct {
	ID        string
	UserID    string
	NotulenID string // Specific notulen document they're viewing
	Conn      *websocket.Conn
	Send      chan []byte
}

// NotulenHub manages WebSocket connections for notulen real-time updates
type NotulenHub struct {
	// Registered clients
	clients map[*NotulenWebSocketClient]bool

	// Inbound messages from clients
	broadcast chan NotulenWebSocketEvent

	// Register requests from clients
	register chan *NotulenWebSocketClient

	// Unregister requests from clients
	unregister chan *NotulenWebSocketClient

	// Mutex for thread safety
	mu sync.RWMutex
}

// NewNotulenHub creates a new WebSocket hub for notulen
func NewNotulenHub() *NotulenHub {
	return &NotulenHub{
		clients:    make(map[*NotulenWebSocketClient]bool),
		broadcast:  make(chan NotulenWebSocketEvent),
		register:   make(chan *NotulenWebSocketClient),
		unregister: make(chan *NotulenWebSocketClient),
	}
}

// Run starts the hub and handles client connections
func (h *NotulenHub) Run(ctx context.Context) {
	logger.Info("Starting Notulen WebSocket Hub")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logger.Info("Notulen WebSocket client registered",
				"client_id", client.ID,
				"user_id", client.UserID,
				"notulen_id", client.NotulenID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			logger.Info("Notulen WebSocket client unregistered",
				"client_id", client.ID,
				"user_id", client.UserID)

		case event := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				// Send to all clients viewing this specific notulen, or all clients if no specific notulen
				if client.NotulenID == "" || client.NotulenID == event.NotulenID {
					select {
					case client.Send <- h.marshalEvent(event):
					default:
						// Client send channel is full, close connection
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()

		case <-ctx.Done():
			logger.Info("Stopping Notulen WebSocket Hub")
			h.mu.Lock()
			for client := range h.clients {
				close(client.Send)
			}
			h.mu.Unlock()
			return
		}
	}
}

// BroadcastNotulenUpdate broadcasts a notulen update to connected clients
func (h *NotulenHub) BroadcastNotulenUpdate(notulenID uuid.UUID, userID uuid.UUID, notulen *models.NotulenResponse) {
	event := NotulenWebSocketEvent{
		Type:      "notulen_updated",
		NotulenID: notulenID.String(),
		UserID:    userID.String(),
		Data:      notulen,
		Timestamp: time.Now(),
	}
	h.broadcast <- event
}

// BroadcastNotulenFinalized broadcasts a notulen finalization to connected clients
func (h *NotulenHub) BroadcastNotulenFinalized(notulenID uuid.UUID, userID uuid.UUID) {
	event := NotulenWebSocketEvent{
		Type:      "notulen_finalized",
		NotulenID: notulenID.String(),
		UserID:    userID.String(),
		Timestamp: time.Now(),
	}
	h.broadcast <- event
}

// BroadcastNotulenArchived broadcasts a notulen archiving to connected clients
func (h *NotulenHub) BroadcastNotulenArchived(notulenID uuid.UUID, userID uuid.UUID) {
	event := NotulenWebSocketEvent{
		Type:      "notulen_archived",
		NotulenID: notulenID.String(),
		UserID:    userID.String(),
		Timestamp: time.Now(),
	}
	h.broadcast <- event
}

// BroadcastNotulenDeleted broadcasts a notulen deletion to connected clients
func (h *NotulenHub) BroadcastNotulenDeleted(notulenID uuid.UUID, userID uuid.UUID) {
	event := NotulenWebSocketEvent{
		Type:      "notulen_deleted",
		NotulenID: notulenID.String(),
		UserID:    userID.String(),
		Timestamp: time.Now(),
	}
	h.broadcast <- event
}

// RegisterClient registers a new WebSocket client
func (h *NotulenHub) RegisterClient(client *NotulenWebSocketClient) {
	h.register <- client
}

// UnregisterClient unregisters a WebSocket client
func (h *NotulenHub) UnregisterClient(client *NotulenWebSocketClient) {
	h.unregister <- client
}

// GetClientCount returns the number of connected clients
func (h *NotulenHub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// marshalEvent converts an event to JSON bytes
func (h *NotulenHub) marshalEvent(event NotulenWebSocketEvent) []byte {
	data, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal WebSocket event", "error", err)
		return []byte{}
	}
	return data
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *NotulenWebSocketClient) ReadPump() {
	defer func() {
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

		switch msgType {
		case "ping":
			c.handlePing()
		case "subscribe":
			c.handleSubscribe(msg)
		case "unsubscribe":
			c.handleUnsubscribe(msg)
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *NotulenWebSocketClient) WritePump() {
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
				"type":      "ping",
				"timestamp": time.Now().Unix(),
			}
			message, _ := json.Marshal(ping)
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

// handlePing responds to ping with pong
func (c *NotulenWebSocketClient) handlePing() {
	pong := map[string]interface{}{
		"type":      "pong",
		"timestamp": time.Now().Unix(),
	}
	message, _ := json.Marshal(pong)
	c.Send <- message
}

// handleSubscribe handles subscription requests from client
func (c *NotulenWebSocketClient) handleSubscribe(msg map[string]interface{}) {
	// For now, notulen clients are subscribed to all updates for their document
	// Future enhancement could allow selective subscriptions
}

// handleUnsubscribe handles unsubscription requests from client
func (c *NotulenWebSocketClient) handleUnsubscribe(msg map[string]interface{}) {
	// For now, notulen clients stay subscribed until disconnect
	// Future enhancement could allow selective unsubscriptions
}
