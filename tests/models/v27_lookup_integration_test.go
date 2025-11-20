package models_test

import (
	"dklautomationgo/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEvent_V27_StatusFieldMapping tests that Event.Status maps correctly to database
func TestEvent_V27_StatusFieldMapping(t *testing.T) {
	event := &models.Event{
		Status: models.EventStatusUpcoming, // Database column: 'status'
	}

	// Verify Status field is a string (not pointer)
	assert.Equal(t, models.EventStatusUpcoming, event.Status)
	assert.IsType(t, "", event.Status)

	// Test status assignment
	event.Status = models.EventStatusActive
	assert.Equal(t, models.EventStatusActive, event.Status)

	// Test default status
	newEvent := &models.Event{}
	// Status should be set via GORM default tag
	assert.Equal(t, "", newEvent.Status) // Empty until saved to DB
}

// TestEvent_V27_StatusTypePreload tests StatusType relation can be preloaded
func TestEvent_V27_StatusTypePreload(t *testing.T) {
	event := &models.Event{}

	// StatusType should be available for preloading
	assert.NotNil(t, event.StatusType)

	// Simulate preloaded StatusType
	event.StatusType = models.EventStatusType{
		Status:      models.EventStatusActive,
		Description: "Event is currently active",
	}

	assert.Equal(t, models.EventStatusActive, event.StatusType.Status)
	assert.Equal(t, "Event is currently active", event.StatusType.Description)
}

// TestNotification_V27_TypeAndPriorityFieldMapping tests Notification field mapping
func TestNotification_V27_TypeAndPriorityFieldMapping(t *testing.T) {
	notification := &models.Notification{
		Type:     models.NotificationTypeContact,  // Database column: 'type'
		Priority: models.NotificationPriorityHigh, // Database column: 'priority'
	}

	// Verify fields are strings (not pointers)
	assert.Equal(t, models.NotificationTypeContact, notification.Type)
	assert.Equal(t, models.NotificationPriorityHigh, notification.Priority)
	assert.IsType(t, "", notification.Type)
	assert.IsType(t, "", notification.Priority)

	// Test field assignment
	notification.Type = models.NotificationTypeSystem
	notification.Priority = models.NotificationPriorityCritical

	assert.Equal(t, models.NotificationTypeSystem, notification.Type)
	assert.Equal(t, models.NotificationPriorityCritical, notification.Priority)
}

// TestNotification_V27_LookupPreload tests lookup relation preloading
func TestNotification_V27_LookupPreload(t *testing.T) {
	notification := &models.Notification{}

	// Simulate preloaded lookups
	notification.TypeLookup = models.NotificationType{
		Name:        models.NotificationTypeContact,
		Description: "Contact form notification",
	}

	notification.PriorityLookup = models.NotificationPriorityType{
		Name:        models.NotificationPriorityHigh,
		Description: "High priority notification",
	}

	assert.Equal(t, "Contact form notification", notification.TypeLookup.Description)
	assert.Equal(t, "High priority notification", notification.PriorityLookup.Description)
}

// TestContactFormulier_V27_StatusFieldMapping tests Contact status field
func TestContactFormulier_V27_StatusFieldMapping(t *testing.T) {
	contact := &models.ContactFormulier{
		Status: "nieuw", // Database column: 'status'
	}

	// Verify Status field is a string (not pointer)
	assert.Equal(t, "nieuw", contact.Status)
	assert.IsType(t, "", contact.Status)

	// Test status assignment
	contact.Status = "beantwoord"
	assert.Equal(t, "beantwoord", contact.Status)
}

// TestContactFormulier_V27_StatusTypePreload tests StatusType relation
func TestContactFormulier_V27_StatusTypePreload(t *testing.T) {
	contact := &models.ContactFormulier{}

	// Simulate preloaded StatusType
	contact.StatusType = models.ContactStatusType{
		Status:      "nieuw",
		Description: "Nieuwe contactaanvraag",
	}

	assert.Equal(t, "Nieuwe contactaanvraag", contact.StatusType.Description)
}

// TestChatChannel_V27_TypeFieldMapping tests ChatChannel type field
func TestChatChannel_V27_TypeFieldMapping(t *testing.T) {
	channel := &models.ChatChannel{
		Type: "public", // Database column: 'type'
	}

	// Verify Type field is a string (not pointer)
	assert.Equal(t, "public", channel.Type)
	assert.IsType(t, "", channel.Type)

	// Test type assignment
	channel.Type = "private"
	assert.Equal(t, "private", channel.Type)
}

// TestChatChannel_V27_TypeLookupPreload tests TypeLookup relation
func TestChatChannel_V27_TypeLookupPreload(t *testing.T) {
	channel := &models.ChatChannel{}

	// Simulate preloaded TypeLookup
	channel.TypeLookup = models.ChatChannelType{
		Type:        "direct",
		Description: "Direct message channel",
	}

	assert.Equal(t, "Direct message channel", channel.TypeLookup.Description)
}

// TestV27_AllStatusConstants verifies status constants are strings
func TestV27_AllStatusConstants(t *testing.T) {
	// All V27 constants should be strings
	assert.IsType(t, "", models.EventStatusUpcoming)
	assert.IsType(t, "", models.EventStatusActive)
	assert.IsType(t, "", models.EventStatusCompleted)
	assert.IsType(t, "", models.EventStatusCancelled)

	assert.IsType(t, "", models.NotificationTypeContact)
	assert.IsType(t, "", models.NotificationTypeAanmelding)
	assert.IsType(t, "", models.NotificationTypeAuth)
	assert.IsType(t, "", models.NotificationTypeSystem)

	assert.IsType(t, "", models.NotificationPriorityLow)
	assert.IsType(t, "", models.NotificationPriorityMedium)
	assert.IsType(t, "", models.NotificationPriorityHigh)
	assert.IsType(t, "", models.NotificationPriorityCritical)
}

// TestV27_NoPointerFields verifies status/type/priority fields are not pointers
func TestV27_NoPointerFields(t *testing.T) {
	event := models.Event{}
	notification := models.Notification{}
	contact := models.ContactFormulier{}
	channel := models.ChatChannel{}

	// Event.Status should NOT be a pointer
	assert.IsType(t, "", event.Status)

	// Notification Type/Priority should NOT be pointers
	assert.IsType(t, "", notification.Type)
	assert.IsType(t, "", notification.Priority)

	// Contact.Status should NOT be a pointer
	assert.IsType(t, "", contact.Status)

	// ChatChannel.Type should NOT be a pointer
	assert.IsType(t, "", channel.Type)
}

// TestEventResponse_V27_BackwardsCompatibility tests API response structure
func TestEventResponse_V27_BackwardsCompatibility(t *testing.T) {
	event := &models.Event{
		ID:     "test-event",
		Name:   "Test Event",
		Status: models.EventStatusActive,
	}

	resp := event.ToResponse()

	// Response should contain status field
	assert.Equal(t, models.EventStatusActive, resp.Status)

	// Response should be backwards compatible
	assert.NotNil(t, resp)
	assert.Equal(t, "test-event", resp.ID)
	assert.Equal(t, "Test Event", resp.Name)
}
