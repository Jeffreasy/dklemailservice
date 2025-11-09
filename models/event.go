package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// EventConfig is een custom type voor event configuratie JSONB
type EventConfig map[string]interface{}

// Scan implementeert sql.Scanner interface voor database reading
func (ec *EventConfig) Scan(value interface{}) error {
	if value == nil {
		*ec = EventConfig{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*ec = EventConfig(result)
	return nil
}

// Value implementeert driver.Valuer interface voor database writing
func (ec EventConfig) Value() (driver.Value, error) {
	if len(ec) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(ec)
}

// Event representeert een loopwedstrijd event
// V27 Update: Status now uses foreign key to event_status_types lookup table
type Event struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description,omitempty" gorm:"type:text"`
	StartTime   time.Time  `json:"start_time" gorm:"not null"`
	EndTime     *time.Time `json:"end_time,omitempty"`

	// V27: Foreign key to lookup table (database column is 'status')
	Status     string          `json:"status" gorm:"type:text;index;default:'upcoming'"`
	StatusType EventStatusType `json:"status_type,omitempty" gorm:"foreignKey:Status;references:Status"`

	Geofences   Geofences   `json:"geofences" gorm:"type:jsonb;default:'[]'"`
	EventConfig EventConfig `json:"event_config,omitempty" gorm:"type:jsonb;default:'{}'"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy   *string     `json:"created_by,omitempty" gorm:"type:uuid"`
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

// DEPRECATED: EventParticipant is deprecated as of V28 migration
// Use EventRegistration from event_registration.go instead
// This type alias provides backwards compatibility and will be removed in v2.0
//
// Migration: Replace all EventParticipant usage with EventRegistration
// Database: Table renamed from event_participants → event_registrations (V28)
type EventParticipant = EventRegistration

// EventResponse is de response structuur voor API endpoints
type EventResponse struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Description       string      `json:"description,omitempty"`
	StartTime         string      `json:"start_time"` // ISO 8601 format
	EndTime           string      `json:"end_time,omitempty"`
	Status            string      `json:"status"`                       // V27: Direct database column
	StatusDescription string      `json:"status_description,omitempty"` // V27: Lookup table description
	Geofences         []Geofence  `json:"geofences"`
	EventConfig       EventConfig `json:"event_config,omitempty"`
	IsActive          bool        `json:"is_active"`
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

	// V27: Include status description if StatusType is preloaded
	if e.StatusType.Status != "" {
		resp.StatusDescription = e.StatusType.Description
	}

	if e.EndTime != nil {
		resp.EndTime = e.EndTime.Format(time.RFC3339)
	}

	return resp
}

// EventCreateRequest is de request structuur voor het aanmaken van events
type EventCreateRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description,omitempty"`
	StartTime   string      `json:"start_time" validate:"required"` // ISO 8601 format
	EndTime     string      `json:"end_time,omitempty"`
	Status      string      `json:"status,omitempty"`
	Geofences   []Geofence  `json:"geofences" validate:"required"`
	EventConfig EventConfig `json:"event_config,omitempty"`
	IsActive    bool        `json:"is_active"`
}

// EventUpdateRequest is de request structuur voor het bijwerken van events
type EventUpdateRequest struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	StartTime   string      `json:"start_time,omitempty"`
	EndTime     string      `json:"end_time,omitempty"`
	Status      string      `json:"status,omitempty"`
	Geofences   []Geofence  `json:"geofences,omitempty"`
	EventConfig EventConfig `json:"event_config,omitempty"`
	IsActive    *bool       `json:"is_active,omitempty"`
}

// EventParticipantResponse is DEPRECATED, use EventRegistrationResponse instead
// Alias kept for backward compatibility
type EventParticipantResponse = EventRegistrationResponse

// EventRegistrationResponse is de response voor event registration data
type EventRegistrationResponse struct {
	ID                  string  `json:"id"`
	EventID             string  `json:"event_id"`
	EventName           string  `json:"event_name,omitempty"`
	ParticipantID       string  `json:"participant_id"`
	ParticipantName     string  `json:"participant_name,omitempty"`
	Status              string  `json:"status"`                       // V27: Registration status
	StatusDescription   string  `json:"status_description,omitempty"` // V27: Status lookup description
	TrackingStatus      string  `json:"tracking_status"`
	DistanceRoute       string  `json:"distance_route,omitempty"`               // V26: Distance route
	ParticipantRole     string  `json:"participant_role,omitempty"`             // V26: Role name
	ParticipantRoleDesc string  `json:"participant_role_description,omitempty"` // V26: Role description
	RegisteredAt        string  `json:"registered_at"`
	CheckInTime         string  `json:"check_in_time,omitempty"`
	StartTime           string  `json:"start_time,omitempty"`
	FinishTime          string  `json:"finish_time,omitempty"`
	TotalDistance       float64 `json:"total_distance"`
	Steps               int     `json:"steps"`         // Changed from CurrentSteps to match EventRegistration
	CurrentSteps        int     `json:"current_steps"` // Deprecated, kept for backwards compatibility
	LastLocationUpdate  string  `json:"last_location_update,omitempty"`
}

// ToResponse converteert EventParticipant naar EventParticipantResponse
// Works with EventRegistration (EventParticipant is now an alias)
func (ep *EventParticipant) ToResponse() *EventParticipantResponse {
	trackingStatus := "registered" // default
	if ep.TrackingStatus != nil {
		trackingStatus = *ep.TrackingStatus
	}

	resp := &EventParticipantResponse{
		ID:             ep.ID,
		EventID:        ep.EventID,
		ParticipantID:  ep.ParticipantID,
		Status:         ep.Status,
		TrackingStatus: trackingStatus,
		RegisteredAt:   ep.RegisteredAt.Format(time.RFC3339),
		TotalDistance:  ep.TotalDistance,
		Steps:          ep.Steps,
		CurrentSteps:   ep.Steps, // Backwards compatibility
	}

	// V27: Include status description if RegistrationStatus is preloaded
	if ep.RegistrationStatus.Status != "" {
		resp.StatusDescription = ep.RegistrationStatus.Description
	}

	// V26: Include distance route if set
	if ep.DistanceRoute != nil {
		resp.DistanceRoute = *ep.DistanceRoute
	}

	// V26: Include role info if ParticipantRole is preloaded
	if ep.ParticipantRoleName != nil {
		resp.ParticipantRole = *ep.ParticipantRoleName
		if ep.ParticipantRole.Name != "" {
			resp.ParticipantRoleDesc = ep.ParticipantRole.Description
		}
	}

	// Check if Event relation is loaded (not zero value)
	if ep.Event.ID != "" {
		resp.EventName = ep.Event.Name
	}

	// Check if Participant relation is loaded (not zero value)
	if ep.Participant.ID != "" {
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
