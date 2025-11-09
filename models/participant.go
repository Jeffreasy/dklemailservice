package models

import (
	"time"
)

// Participant representeert de "Persoon" (voorheen Aanmelding)
// GEFIXTE STRUCT: Bevat alleen persoonsgegevens.
// Velden als Steps en Afstand zijn verwijderd omdat ze nu in EventRegistration horen.
type Participant struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Naam        string    `json:"naam" gorm:"not null"`
	Email       string    `json:"email" gorm:"not null;index"`
	Telefoon    string    `json:"telefoon"`
	Terms       bool      `json:"terms" gorm:"not null"` // 'terms' bleef op Participant (V28)
	GebruikerID *string   `json:"gebruiker_id,omitempty" gorm:"type:uuid;index;column:gebruiker_id"`
	TestMode    bool      `json:"test_mode" gorm:"default:false"` // Dit veld ontbrak nog (uit V1)

	// Relatie met antwoorden
	Antwoorden []ParticipantAntwoord `json:"antwoorden,omitempty" gorm:"foreignKey:ParticipantID"`

	// Relatie met event registraties
	EventRegistrations []EventRegistration `json:"event_registrations,omitempty" gorm:"foreignKey:ParticipantID"`
}

// TableName specificeert de tabelnaam voor GORM
func (Participant) TableName() string {
	return "participants"
}

// AanmeldingFormulier (DTO) blijft ongewijzigd
type AanmeldingFormulier struct {
	Naam           string `json:"naam"`
	Email          string `json:"email"`
	Telefoon       string `json:"telefoon"`
	Rol            string `json:"rol"`
	Afstand        string `json:"afstand"`
	Ondersteuning  string `json:"ondersteuning"`
	Bijzonderheden string `json:"bijzonderheden"`
	Terms          bool   `json:"terms"`
	TestMode       bool   `json:"test_mode"`
}
