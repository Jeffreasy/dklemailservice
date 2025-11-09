package handlers_test

import (
	"bytes"
	"context"
	"dklautomationgo/handlers"
	"dklautomationgo/models"
	"dklautomationgo/tests"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository for testing
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRepository) List(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRepository) Update(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) GetActiveEvent(ctx context.Context) (*models.Event, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRepository) ListActive(ctx context.Context) ([]*models.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRepository) GetEventParticipants(ctx context.Context, eventID string) ([]*models.EventParticipant, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).([]*models.EventParticipant), args.Error(1)
}

func (m *MockEventRepository) GetParticipantEvents(ctx context.Context, participantID string) ([]*models.EventParticipant, error) {
	args := m.Called(ctx, participantID)
	return args.Get(0).([]*models.EventParticipant), args.Error(1)
}

func (m *MockEventRepository) GetEventParticipant(ctx context.Context, eventID, participantID string) (*models.EventParticipant, error) {
	args := m.Called(ctx, eventID, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventParticipant), args.Error(1)
}

func (m *MockEventRepository) RegisterParticipant(ctx context.Context, eventID, participantID string) (*models.EventParticipant, error) {
	args := m.Called(ctx, eventID, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventParticipant), args.Error(1)
}

func (m *MockEventRepository) UpdateEventParticipant(ctx context.Context, ep *models.EventParticipant) error {
	args := m.Called(ctx, ep)
	return args.Error(0)
}

// TestEventHandler_CreateEvent_V27_StatusField tests event creation with V27 status field
func TestEventHandler_CreateEvent_V27_StatusField(t *testing.T) {
	mockEventRepo := new(MockEventRepository)
	mockAuthService := new(tests.MockAuthService)
	mockPermService := tests.NewMockPermissionService()

	handler := handlers.NewEventHandler(mockEventRepo, mockAuthService, mockPermService)

	app := fiber.New()
	app.Post("/api/events", func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user")
		return handler.CreateEvent(c)
	})

	// Test with explicit status
	requestBody := models.EventCreateRequest{
		Name:      "Test Event",
		StartTime: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		Status:    models.EventStatusUpcoming,
		Geofences: []models.Geofence{
			{Type: "start", Lat: 52.0, Long: 5.0, Radius: 100},
		},
		IsActive: true,
	}

	bodyBytes, _ := json.Marshal(requestBody)

	mockEventRepo.On("Create", mock.Anything, mock.MatchedBy(func(event *models.Event) bool {
		// Verify Status field is set correctly (string, not pointer)
		return event.Status == models.EventStatusUpcoming &&
			event.Name == "Test Event"
	})).Return(nil)

	req := httptest.NewRequest("POST", "/api/events", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	mockEventRepo.AssertExpectations(t)
}

// TestEventHandler_CreateEvent_V27_DefaultStatus tests default status assignment
func TestEventHandler_CreateEvent_V27_DefaultStatus(t *testing.T) {
	mockEventRepo := new(MockEventRepository)
	mockAuthService := new(tests.MockAuthService)
	mockPermService := tests.NewMockPermissionService()

	handler := handlers.NewEventHandler(mockEventRepo, mockAuthService, mockPermService)

	app := fiber.New()
	app.Post("/api/events", func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user")
		return handler.CreateEvent(c)
	})

	// Test without explicit status (should default to 'upcoming')
	requestBody := models.EventCreateRequest{
		Name:      "Test Event No Status",
		StartTime: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		// Status NOT provided
		Geofences: []models.Geofence{
			{Type: "start", Lat: 52.0, Long: 5.0, Radius: 100},
		},
		IsActive: true,
	}

	bodyBytes, _ := json.Marshal(requestBody)

	mockEventRepo.On("Create", mock.Anything, mock.MatchedBy(func(event *models.Event) bool {
		// Should default to EventStatusUpcoming
		return event.Status == models.EventStatusUpcoming
	})).Return(nil)

	req := httptest.NewRequest("POST", "/api/events", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	mockEventRepo.AssertExpectations(t)
}

// TestEventHandler_UpdateEvent_V27_StatusField tests event update with status
func TestEventHandler_UpdateEvent_V27_StatusField(t *testing.T) {
	mockEventRepo := new(MockEventRepository)
	mockAuthService := new(tests.MockAuthService)
	mockPermService := tests.NewMockPermissionService()

	handler := handlers.NewEventHandler(mockEventRepo, mockAuthService, mockPermService)

	app := fiber.New()
	app.Put("/api/events/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user")
		return handler.UpdateEvent(c)
	})

	existingEvent := &models.Event{
		ID:        "test-event-id",
		Name:      "Existing Event",
		StartTime: time.Now(),
		Status:    models.EventStatusUpcoming,
	}

	requestBody := models.EventUpdateRequest{
		Status: models.EventStatusActive, // Update status to active
	}

	bodyBytes, _ := json.Marshal(requestBody)

	mockEventRepo.On("GetByID", mock.Anything, "test-event-id").Return(existingEvent, nil)
	mockEventRepo.On("Update", mock.Anything, mock.MatchedBy(func(event *models.Event) bool {
		// Verify status was updated correctly (direct field, not pointer)
		return event.Status == models.EventStatusActive
	})).Return(nil)

	req := httptest.NewRequest("PUT", "/api/events/test-event-id", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockEventRepo.AssertExpectations(t)
}

// TestEventResponse_V27_StatusKey tests that response includes status
func TestEventResponse_V27_StatusKey(t *testing.T) {
	event := &models.Event{
		ID:        "test-event",
		Name:      "Test Event",
		StartTime: time.Now(),
		Status:    models.EventStatusActive,
		Geofences: []models.Geofence{},
	}

	resp := event.ToResponse()

	// Verify response contains correct status
	assert.Equal(t, models.EventStatusActive, resp.Status)
	assert.IsType(t, "", resp.Status)
}
