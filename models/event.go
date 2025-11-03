package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Event representeert een loopwedstrijd event
type Event struct {
	ID          string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string                 `json:"name" gorm:"not null"`
	Description string                 `json:"description,omitempty" gorm:"type:text"`
	StartTime   time.Time              `json:"start_time" gorm:"not null"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	Status      string                 `json:"status" gorm:"default:'upcoming'"`
	Geofences   Geofences              `json:"geofences" gorm:"type:jsonb;default:'[]'"`
	EventConfig map[string]interface{} `json:"event_config,omitempty" gorm:"type:jsonb;default:'{}'"`
	IsActive    bool                   `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy   *string                `json:"created_by,omitempty" gorm:"type:uuid"`
}

// TableName specificeert de tabelnaam voor GORM
func (Event) TableName() string {
	return "events"
}

// Geofence representeert een geografische fence voor event tracking
type Geofence struct {
	Type   string  `json:"type"`   // "start", "checkpoint", "finish"
	Lat    float64 `json:"lat"`    // Latitude
	Long   float64 `json:"long"`   // Longitude
	Radius float64 `json:"radius"` // Radius in meters
	Name   string  `json:"name,omitempty"`
}

// Geofences is een array van Geofence voor database storage
type Geofences []Geofence

// Scan implementeert sql.Scanner interface voor database reading
func (g *Geofences) Scan(value interface{}) error {
	if value == nil {
		*g = []Geofence{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, g)
}

// Value implementeert driver.Valuer interface voor database writing
func (g Geofences) Value() (driver.Value, error) {
	if len(g) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(g)
}

// EventParticipant representeert de koppeling tussen event en participant
type EventParticipant struct {
	ID                 string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EventID            string     `json:"event_id" gorm:"type:uuid;not null"`
	ParticipantID      string     `json:"participant_id" gorm:"type:uuid;not null"`
	RegisteredAt       time.Time  `json:"registered_at" gorm:"autoCreateTime"`
	CheckInTime        *time.Time `json:"check_in_time,omitempty"`
	StartTime          *time.Time `json:"start_time,omitempty"`
	FinishTime         *time.Time `json:"finish_time,omitempty"`
	TrackingStatus     string     `json:"tracking_status" gorm:"default:'registered'"`
	LastLocationUpdate *time.Time `json:"last_location_update,omitempty"`
	TotalDistance      float64    `json:"total_distance" gorm:"type:decimal(10,2);default:0"`
	CurrentSteps       int        `json:"current_steps" gorm:"default:0"`

	// Relations
	Event       *Event      `json:"event,omitempty" gorm:"foreignKey:EventID"`
	Participant *Aanmelding `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
}

// TableName specificeert de tabelnaam voor GORM
func (EventParticipant) TableName() string {
	return "event_participants"
}

// EventResponse is de response structuur voor API endpoints
type EventResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	StartTime   string                 `json:"start_time"` // ISO 8601 format
	EndTime     string                 `json:"end_time,omitempty"`
	Status      string                 `json:"status"`
	Geofences   []Geofence             `json:"geofences"`
	EventConfig map[string]interface{} `json:"event_config,omitempty"`
	IsActive    bool                   `json:"is_active"`
}

// ToResponse converteert Event naar EventResponse
func (e *Event) ToResponse() *EventResponse {
	resp := &EventResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   e.StartTime.Format(time.RFC3339),
		Status:      e.Status,
		Geofences:   e.Geofences,
		EventConfig: e.EventConfig,
		IsActive:    e.IsActive,
	}

	if e.EndTime != nil {
		resp.EndTime = e.EndTime.Format(time.RFC3339)
	}

	return resp
}

// EventCreateRequest is de request structuur voor het aanmaken van events
type EventCreateRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description,omitempty"`
	StartTime   string                 `json:"start_time" validate:"required"` // ISO 8601 format
	EndTime     string                 `json:"end_time,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Geofences   []Geofence             `json:"geofences" validate:"required"`
	EventConfig map[string]interface{} `json:"event_config,omitempty"`
	IsActive    bool                   `json:"is_active"`
}

// EventUpdateRequest is de request structuur voor het bijwerken van events
type EventUpdateRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	StartTime   string                 `json:"start_time,omitempty"`
	EndTime     string                 `json:"end_time,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Geofences   []Geofence             `json:"geofences,omitempty"`
	EventConfig map[string]interface{} `json:"event_config,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
}

// EventParticipantResponse is de response voor event participant data
type EventParticipantResponse struct {
	ID                 string  `json:"id"`
	EventID            string  `json:"event_id"`
	EventName          string  `json:"event_name,omitempty"`
	ParticipantID      string  `json:"participant_id"`
	ParticipantName    string  `json:"participant_name,omitempty"`
	TrackingStatus     string  `json:"tracking_status"`
	RegisteredAt       string  `json:"registered_at"`
	CheckInTime        string  `json:"check_in_time,omitempty"`
	StartTime          string  `json:"start_time,omitempty"`
	FinishTime         string  `json:"finish_time,omitempty"`
	TotalDistance      float64 `json:"total_distance"`
	CurrentSteps       int     `json:"current_steps"`
	LastLocationUpdate string  `json:"last_location_update,omitempty"`
}

// ToResponse converteert EventParticipant naar EventParticipantResponse
func (ep *EventParticipant) ToResponse() *EventParticipantResponse {
	resp := &EventParticipantResponse{
		ID:             ep.ID,
		EventID:        ep.EventID,
		ParticipantID:  ep.ParticipantID,
		TrackingStatus: ep.TrackingStatus,
		RegisteredAt:   ep.RegisteredAt.Format(time.RFC3339),
		TotalDistance:  ep.TotalDistance,
		CurrentSteps:   ep.CurrentSteps,
	}

	if ep.Event != nil {
		resp.EventName = ep.Event.Name
	}

	if ep.Participant != nil {
		resp.ParticipantName = ep.Participant.Naam
	}

	if ep.CheckInTime != nil {
		resp.CheckInTime = ep.CheckInTime.Format(time.RFC3339)
	}

	if ep.StartTime != nil {
		resp.StartTime = ep.StartTime.Format(time.RFC3339)
	}

	if ep.FinishTime != nil {
		resp.FinishTime = ep.FinishTime.Format(time.RFC3339)
	}

	if ep.LastLocationUpdate != nil {
		resp.LastLocationUpdate = ep.LastLocationUpdate.Format(time.RFC3339)
	}

	return resp
}

// LocationUpdateRequest is de request voor location updates tijdens event
type LocationUpdateRequest struct {
	Lat       float64 `json:"lat" validate:"required"`
	Long      float64 `json:"long" validate:"required"`
	Accuracy  float64 `json:"accuracy,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"` // ISO 8601 format
}

// Constanten voor event status
const (
	EventStatusUpcoming  = "upcoming"
	EventStatusActive    = "active"
	EventStatusCompleted = "completed"
	EventStatusCancelled = "cancelled"
)

// Constanten voor tracking status
const (
	TrackingStatusRegistered = "registered"
	TrackingStatusCheckedIn  = "checked_in"
	TrackingStatusStarted    = "started"
	TrackingStatusInProgress = "in_progress"
	TrackingStatusFinished   = "finished"
	TrackingStatusDNF        = "dnf" // Did Not Finish
)

// Constanten voor geofence types
const (
	GeofenceTypeStart      = "start"
	GeofenceTypeCheckpoint = "checkpoint"
	GeofenceTypeFinish     = "finish"
)
