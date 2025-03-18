package models

import (
	"time"
)

type Aanmelding struct {
	ID               string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Naam             string     `json:"naam" gorm:"not null"`
	Email            string     `json:"email" gorm:"not null;index"`
	Telefoon         string     `json:"telefoon"`
	Rol              string     `json:"rol"`
	Afstand          string     `json:"afstand"`
	Ondersteuning    string     `json:"ondersteuning"`
	Bijzonderheden   string     `json:"bijzonderheden" gorm:"type:text"`
	Terms            bool       `json:"terms" gorm:"not null"`
	EmailVerzonden   bool       `json:"email_verzonden" gorm:"default:false"`
	EmailVerzondenOp *time.Time `json:"email_verzonden_op"`

	// Nieuwe velden voor status en behandeling
	Status        string     `json:"status" gorm:"default:'nieuw';index"`
	BehandeldDoor *string    `json:"behandeld_door"`
	BehandeldOp   *time.Time `json:"behandeld_op"`
	Notities      *string    `json:"notities" gorm:"type:text"`

	// Test mode veld
	TestMode bool `json:"test_mode" gorm:"type:boolean;not null;default:false"`

	// Relatie met antwoorden
	Antwoorden []AanmeldingAntwoord `json:"antwoorden,omitempty" gorm:"foreignKey:AanmeldingID"`
}

// TableName specificeert de tabelnaam voor GORM
func (Aanmelding) TableName() string {
	return "aanmeldingen"
}

// AanmeldingFormulier represents the registration form data from the frontend
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
