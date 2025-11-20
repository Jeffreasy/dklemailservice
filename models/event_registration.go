package models

import (
	"time"
)

// EventRegistration representeert de "Deelname" aan een evenement.
// Deze tabel was voorheen 'event_participants' (V23) en bevat nu
// data die voorheen op 'aanmeldingen' stond (V16, V27).
type EventRegistration struct {
	ID            string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EventID       string     `json:"event_id" gorm:"type:uuid;not null;index"`
	ParticipantID string     `json:"participant_id" gorm:"type:uuid;not null;index"`
	RegisteredAt  time.Time  `json:"registered_at" gorm:"autoCreateTime"`
	CheckInTime   *time.Time `json:"check_in_time"`
	StartTime     *time.Time `json:"start_time"`
	FinishTime    *time.Time `json:"finish_time"`

	// GPS Tracking velden (V23 event_participants)
	TrackingStatus     *string    `json:"tracking_status" gorm:"default:'registered'"`
	LastLocationUpdate *time.Time `json:"last_location_update"`
	TotalDistance      float64    `json:"total_distance" gorm:"type:decimal(10,2);default:0"`

	// --- Velden verplaatst van 'Aanmelding' (nu Participant) ---
	Steps          int    `json:"steps" gorm:"default:0"` // V16
	TestMode       bool   `json:"test_mode" gorm:"type:boolean;not null;default:false"`
	Ondersteuning  string `json:"ondersteuning"`
	Bijzonderheden string `json:"bijzonderheden" gorm:"type:text"`

	// V37: Terms veld behouden in event_registrations
	Terms bool `json:"terms" gorm:"not null;default:false"`

	// Antwoorden count voor tracking
	AntwoordenCount int `json:"antwoorden_count" gorm:"default:0"`

	// V37: Transport vraag
	HeeftVervoer *bool `json:"heeft_vervoer" gorm:"type:boolean"`

	Notities         *string    `json:"notities" gorm:"type:text"`
	BehandeldDoor    *string    `json:"behandeld_door"`
	BehandeldOp      *time.Time `json:"behandeld_op"`
	EmailVerzonden   bool       `json:"email_verzonden" gorm:"default:false"`
	EmailVerzondenOp *time.Time `json:"email_verzonden_op"`

	// --- Foreign Keys (V26 & V27) ---
	// V27: Database column is 'status' (renamed from status_key)
	Status             string                 `json:"status" gorm:"type:text;index;default:'registered'"`
	RegistrationStatus RegistrationStatusType `json:"registration_status,omitempty" gorm:"foreignKey:Status;references:Status"`

	// V26: Database column is 'distance_route'
	DistanceRoute *string  `json:"distance_route" gorm:"type:text;index"`
	Distance      Distance `json:"distance,omitempty" gorm:"foreignKey:DistanceRoute;references:Route"`

	// V26: Database column is 'participant_role_name'
	ParticipantRoleName *string         `json:"participant_role_name" gorm:"type:text;index"`
	ParticipantRole     ParticipantRole `json:"role,omitempty" gorm:"foreignKey:ParticipantRoleName;references:Name"`

	// --- GORM Relaties ---
	Participant Participant `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	Event       Event       `json:"event,omitempty" gorm:"foreignKey:EventID"`
}

// UpdateEventRegistrationRequest represents the request body for updating an event registration
type UpdateEventRegistrationRequest struct {
	Status   string  `json:"status" example:"confirmed"`
	Notities *string `json:"notities" example:"Some notes about the registration"`
}

// TableName specificeert de tabelnaam voor GORM
func (EventRegistration) TableName() string {
	return "event_registrations"
}
