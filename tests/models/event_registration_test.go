package models_test

import (
	"dklautomationgo/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEventRegistration_FieldMapping tests that EventRegistration fields map correctly
func TestEventRegistration_FieldMapping(t *testing.T) {
	now := time.Now()

	reg := &models.EventRegistration{
		ID:            "test-id",
		EventID:       "event-id",
		ParticipantID: "participant-id",
		RegisteredAt:  now,
		Steps:         1000,
		TotalDistance: 5.5,
		TestMode:      true,
	}

	assert.Equal(t, "test-id", reg.ID)
	assert.Equal(t, "event-id", reg.EventID)
	assert.Equal(t, "participant-id", reg.ParticipantID)
	assert.Equal(t, now, reg.RegisteredAt)
	assert.Equal(t, 1000, reg.Steps)
	assert.Equal(t, 5.5, reg.TotalDistance)
	assert.True(t, reg.TestMode)
}

// TestEventRegistration_ToResponse tests response conversion
func TestEventRegistration_ToResponse(t *testing.T) {
	now := time.Now()
	checkIn := now.Add(1 * time.Hour)
	startTime := now.Add(2 * time.Hour)
	finishTime := now.Add(3 * time.Hour)
	lastUpdate := now.Add(4 * time.Hour)

	trackingStatus := "in_progress"

	reg := &models.EventRegistration{
		ID:                 "test-id",
		EventID:            "event-id",
		ParticipantID:      "participant-id",
		TrackingStatus:     &trackingStatus,
		RegisteredAt:       now,
		CheckInTime:        &checkIn,
		StartTime:          &startTime,
		FinishTime:         &finishTime,
		Steps:              1500,
		TotalDistance:      7.5,
		LastLocationUpdate: &lastUpdate,
	}

	// Test with relations loaded
	reg.Event = models.Event{ID: "event-id", Name: "Test Event"}
	reg.Participant = models.Participant{ID: "participant-id", Naam: "Test Participant"}

	resp := reg.ToResponse()

	assert.Equal(t, "test-id", resp.ID)
	assert.Equal(t, "event-id", resp.EventID)
	assert.Equal(t, "Test Event", resp.EventName)
	assert.Equal(t, "participant-id", resp.ParticipantID)
	assert.Equal(t, "Test Participant", resp.ParticipantName)
	assert.Equal(t, "in_progress", resp.TrackingStatus)
	assert.Equal(t, 1500, resp.Steps)
	assert.Equal(t, 1500, resp.CurrentSteps) // Backwards compatibility
	assert.Equal(t, 7.5, resp.TotalDistance)
	assert.NotEmpty(t, resp.RegisteredAt)
	assert.NotEmpty(t, resp.CheckInTime)
	assert.NotEmpty(t, resp.StartTime)
	assert.NotEmpty(t, resp.FinishTime)
	assert.NotEmpty(t, resp.LastLocationUpdate)
}

// TestEventRegistration_BackwardsCompatibility tests type alias
func TestEventRegistration_BackwardsCompatibility(t *testing.T) {
	// EventParticipant should be an alias for EventRegistration
	var ep models.EventParticipant = models.EventRegistration{
		ID:            "test-id",
		EventID:       "event-id",
		ParticipantID: "participant-id",
	}

	assert.Equal(t, "test-id", ep.ID)
	assert.Equal(t, "event-id", ep.EventID)
	assert.Equal(t, "participant-id", ep.ParticipantID)
}

// TestEventRegistration_DefaultValues tests default field values
func TestEventRegistration_DefaultValues(t *testing.T) {
	reg := &models.EventRegistration{}

	assert.Equal(t, 0, reg.Steps)
	assert.Equal(t, 0.0, reg.TotalDistance)
	assert.Nil(t, reg.TrackingStatus)
	assert.Nil(t, reg.CheckInTime)
	assert.Nil(t, reg.StartTime)
	assert.Nil(t, reg.FinishTime)
}
