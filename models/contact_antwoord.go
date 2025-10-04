package models

import (
	"time"
)

// ContactAntwoord representeert een antwoord op een contactformulier
type ContactAntwoord struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ContactID      string    `json:"contact_id" gorm:"not null;index"`
	Tekst          string    `json:"tekst" gorm:"type:text;not null"`
	VerzondOp      time.Time `json:"verzonden_op" gorm:"autoCreateTime"`
	VerzondDoor    string    `json:"verzonden_door" gorm:"not null"`
	EmailVerzonden bool      `json:"email_verzonden" gorm:"default:false"`

	// Relatie met ContactFormulier
	Contact ContactFormulier `json:"-" gorm:"foreignKey:ContactID"`
}

// TableName specificeert de tabelnaam voor GORM
func (ContactAntwoord) TableName() string {
	return "contact_antwoorden"
}
