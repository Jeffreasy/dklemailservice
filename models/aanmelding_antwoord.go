package models

import (
	"time"
)

// AanmeldingAntwoord representeert een antwoord op een aanmeldingsformulier
type AanmeldingAntwoord struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AanmeldingID   string    `json:"aanmelding_id" gorm:"not null;index"`
	Tekst          string    `json:"tekst" gorm:"type:text;not null"`
	VerzondOp      time.Time `json:"verzonden_op" gorm:"autoCreateTime"`
	VerzondDoor    string    `json:"verzonden_door" gorm:"not null"`
	EmailVerzonden bool      `json:"email_verzonden" gorm:"default:false"`

	// Relatie met Aanmelding
	Aanmelding Aanmelding `json:"-" gorm:"foreignKey:AanmeldingID"`
}

// TableName specificeert de tabelnaam voor GORM
func (AanmeldingAntwoord) TableName() string {
	return "aanmelding_antwoorden"
}
