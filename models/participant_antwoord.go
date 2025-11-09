package models

import (
	"time"
)

// ParticipantAntwoord representeert een antwoord op een participant (voorheen aanmelding)
type ParticipantAntwoord struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ParticipantID  string    `json:"participant_id" gorm:"not null;index"` // Hernoemd
	Tekst          string    `json:"tekst" gorm:"type:text;not null"`
	VerzondOp      time.Time `json:"verzonden_op" gorm:"autoCreateTime"`
	VerzondDoor    string    `json:"verzonden_door" gorm:"not null"`
	EmailVerzonden bool      `json:"email_verzonden" gorm:"default:false"`

	// Relatie met Participant
	Participant Participant `json:"-" gorm:"foreignKey:ParticipantID"` // Hernoemd
}

// TableName specificeert de tabelnaam voor GORM
func (ParticipantAntwoord) TableName() string {
	return "participant_antwoorden" // Hernoemd
}
