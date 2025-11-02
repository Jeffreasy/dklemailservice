package tests

import (
	"dklautomationgo/services"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWebSocketConn is a mock for websocket.Conn
type MockWebSocketConn struct {
	mock.Mock
	messages [][]byte
}

func (m *MockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	m.messages = append(m.messages, data)
	return m.Called(messageType, data).Error(0)
}

func (m *MockWebSocketConn) ReadMessage() (int, []byte, error) {
	args := m.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

func (m *MockWebSocketConn) Close() error {
	return m.Called().Error(0)
}

// TestStepsHub_NewStepsHub tests the constructor
func TestStepsHub_NewStepsHub(t *testing.T) {
	hub := services.NewStepsHub(nil, nil)

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.Clients)
	assert.NotNil(t, hub.StepUpdate)
	assert.NotNil(t, hub.TotalUpdate)
	assert.NotNil(t, hub.LeaderboardUpdate)
	assert.NotNil(t, hub.BadgeEarned)
	assert.NotNil(t, hub.Register)
	assert.NotNil(t, hub.Unregister)
}

// TestStepsHub_ClientRegistration tests client registration
func TestStepsHub_ClientRegistration(t *testing.T) {
	var gamificationService services.GamificationService
	hub := services.NewStepsHub(nil, gamificationService)
	go hub.Run()

	client := &services.StepsClient{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		UserID:        "test-user",
		ParticipantID: "test-participant",
		Subscriptions: make(map[string]bool),
	}

	// Register client
	hub.Register <- client

	// Give hub time to process
	time.Sleep(10 * time.Millisecond)

	// Check client is registered
	assert.Equal(t, 1, hub.GetClientCount())

	// Unregister client
	hub.Unregister <- client

	// Give hub time to process
	time.Sleep(10 * time.Millisecond)

	// Check client is unregistered
	assert.Equal(t, 0, hub.GetClientCount())
}

// TestStepsHub_StepUpdateBroadcast tests step update broadcasting
func TestStepsHub_StepUpdateBroadcast(t *testing.T) {
	var gamificationService services.GamificationService
	hub := services.NewStepsHub(nil, gamificationService)
	go hub.Run()

	// Create client with step_updates subscription
	client := &services.StepsClient{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		UserID:        "test-user",
		ParticipantID: "participant-1",
		Subscriptions: map[string]bool{
			"step_updates": true,
		},
	}

	// Register client
	hub.Register <- client
	time.Sleep(10 * time.Millisecond)

	// Send step update
	testUpdate := &services.StepUpdateMessage{
		Type:           services.MessageTypeStepUpdate,
		ParticipantID:  "participant-1",
		Naam:           "Test User",
		Steps:          1000,
		Delta:          500,
		Route:          "10 KM",
		AllocatedFunds: 75,
		Timestamp:      time.Now().Unix(),
	}

	hub.StepUpdate <- testUpdate

	// Wait for broadcast
	select {
	case msg := <-client.Send:
		var received services.StepUpdateMessage
		err := json.Unmarshal(msg, &received)
		assert.NoError(t, err)
		assert.Equal(t, services.MessageTypeStepUpdate, received.Type)
		assert.Equal(t, "participant-1", received.ParticipantID)
		assert.Equal(t, 1000, received.Steps)
		assert.Equal(t, 500, received.Delta)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

// TestStepsHub_TotalUpdateBroadcast tests total update broadcasting
func TestStepsHub_TotalUpdateBroadcast(t *testing.T) {
	var gamificationService services.GamificationService
	hub := services.NewStepsHub(nil, gamificationService)
	go hub.Run()

	// Create client with total_updates subscription
	client := &services.StepsClient{
		Hub:  hub,
		Send: make(chan []byte, 256),
		Subscriptions: map[string]bool{
			"total_updates": true,
		},
	}

	hub.Register <- client
	time.Sleep(10 * time.Millisecond)

	// Send total update
	testUpdate := &services.TotalUpdateMessage{
		Type:       services.MessageTypeTotalUpdate,
		TotalSteps: 250000,
		Year:       2025,
		Timestamp:  time.Now().Unix(),
	}

	hub.TotalUpdate <- testUpdate

	// Wait for broadcast
	select {
	case msg := <-client.Send:
		var received services.TotalUpdateMessage
		err := json.Unmarshal(msg, &received)
		assert.NoError(t, err)
		assert.Equal(t, services.MessageTypeTotalUpdate, received.Type)
		assert.Equal(t, 250000, received.TotalSteps)
		assert.Equal(t, 2025, received.Year)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

// TestStepsHub_SubscriptionFiltering tests that only subscribed clients receive messages
func TestStepsHub_SubscriptionFiltering(t *testing.T) {
	var gamificationService services.GamificationService
	hub := services.NewStepsHub(nil, gamificationService)
	go hub.Run()

	// Client 1: subscribed to step_updates
	client1 := &services.StepsClient{
		Hub:  hub,
		Send: make(chan []byte, 256),
		Subscriptions: map[string]bool{
			"step_updates": true,
		},
	}

	// Client 2: NOT subscribed to step_updates
	client2 := &services.StepsClient{
		Hub:  hub,
		Send: make(chan []byte, 256),
		Subscriptions: map[string]bool{
			"total_updates": true, // Different subscription
		},
	}

	hub.Register <- client1
	hub.Register <- client2
	time.Sleep(10 * time.Millisecond)

	// Send step update
	testUpdate := &services.StepUpdateMessage{
		Type:          services.MessageTypeStepUpdate,
		ParticipantID: "test-participant",
		Steps:         1000,
		Timestamp:     time.Now().Unix(),
	}

	hub.StepUpdate <- testUpdate

	// Client 1 should receive message
	select {
	case <-client1.Send:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 1 should have received message")
	}

	// Client 2 should NOT receive message
	select {
	case <-client2.Send:
		t.Fatal("Client 2 should not have received message")
	case <-time.After(100 * time.Millisecond):
		// Success - no message received
	}
}

// TestStepsHub_BadgeEarnedTargeting tests badge earned is sent only to specific participant
func TestStepsHub_BadgeEarnedTargeting(t *testing.T) {
	var gamificationService services.GamificationService
	hub := services.NewStepsHub(nil, gamificationService)
	go hub.Run()

	// Client 1: target participant
	client1 := &services.StepsClient{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "participant-1",
		Subscriptions: make(map[string]bool),
	}

	// Client 2: different participant
	client2 := &services.StepsClient{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "participant-2",
		Subscriptions: make(map[string]bool),
	}

	hub.Register <- client1
	hub.Register <- client2
	time.Sleep(10 * time.Millisecond)

	// Send badge earned for participant-1
	testBadge := &services.BadgeEarnedMessage{
		Type:          services.MessageTypeBadgeEarned,
		ParticipantID: "participant-1",
		BadgeName:     "First Steps",
		Points:        10,
		Timestamp:     time.Now().Unix(),
	}

	hub.BadgeEarned <- testBadge

	// Only client 1 should receive message
	select {
	case msg := <-client1.Send:
		var received services.BadgeEarnedMessage
		err := json.Unmarshal(msg, &received)
		assert.NoError(t, err)
		assert.Equal(t, "participant-1", received.ParticipantID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 1 should have received badge message")
	}

	// Client 2 should NOT receive message
	select {
	case <-client2.Send:
		t.Fatal("Client 2 should not have received badge message")
	case <-time.After(100 * time.Millisecond):
		// Success - no message received
	}
}

// TestStepsHub_GetSubscriptionCount tests subscription counting
func TestStepsHub_GetSubscriptionCount(t *testing.T) {
	var gamificationService services.GamificationService
	hub := services.NewStepsHub(nil, gamificationService)
	go hub.Run()

	// Create clients with different subscriptions
	client1 := &services.StepsClient{
		Hub:  hub,
		Send: make(chan []byte, 256),
		Subscriptions: map[string]bool{
			"step_updates":  true,
			"total_updates": true,
		},
	}

	client2 := &services.StepsClient{
		Hub:  hub,
		Send: make(chan []byte, 256),
		Subscriptions: map[string]bool{
			"step_updates":        true,
			"leaderboard_updates": true,
		},
	}

	hub.Register <- client1
	hub.Register <- client2
	time.Sleep(10 * time.Millisecond)

	counts := hub.GetSubscriptionCount()

	assert.Equal(t, 2, counts["step_updates"])
	assert.Equal(t, 1, counts["total_updates"])
	assert.Equal(t, 1, counts["leaderboard_updates"])
}

// TestStepsHub_MultipleClients tests broadcasting to multiple clients
func TestStepsHub_MultipleClients(t *testing.T) {
	var gamificationService services.GamificationService
	hub := services.NewStepsHub(nil, gamificationService)
	go hub.Run()

	// Create 5 clients all subscribed to step_updates
	clients := make([]*services.StepsClient, 5)
	for i := 0; i < 5; i++ {
		clients[i] = &services.StepsClient{
			Hub:  hub,
			Send: make(chan []byte, 256),
			Subscriptions: map[string]bool{
				"step_updates": true,
			},
		}
		hub.Register <- clients[i]
	}

	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 5, hub.GetClientCount())

	// Send one update
	testUpdate := &services.StepUpdateMessage{
		Type:      services.MessageTypeStepUpdate,
		Steps:     1000,
		Timestamp: time.Now().Unix(),
	}

	hub.StepUpdate <- testUpdate

	// All clients should receive the message
	for i, client := range clients {
		select {
		case <-client.Send:
			// Success
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Client %d did not receive message", i)
		}
	}
}
