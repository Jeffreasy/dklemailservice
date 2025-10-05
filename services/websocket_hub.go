package services

import (
	"context"
	"dklautomationgo/models"
	"time"

	"github.com/gofiber/websocket/v2"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	ChatService ChatService
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	UserID string
}

// NewHub creates a new Hub
func NewHub(chatService ChatService) *Hub {
	return &Hub{
		Broadcast:   make(chan []byte),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[*Client]bool),
		ChatService: chatService,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			_ = h.ChatService.UpdatePresence(context.Background(), &models.ChatUserPresence{
				UserID:   client.UserID,
				Status:   "online",
				LastSeen: time.Now(),
			})
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				_ = h.ChatService.UpdatePresence(context.Background(), &models.ChatUserPresence{
					UserID:   client.UserID,
					Status:   "offline",
					LastSeen: time.Now(),
				})
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
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

// ServeWs handles websocket requests from the peer.
func (h *Hub) ServeWs(conn *websocket.Conn) {
	client := &Client{Hub: h, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	client.ReadPump()
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		c.Hub.Broadcast <- message
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
	// Channel is nu gesloten (loop gestopt), dus stuur close message
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
