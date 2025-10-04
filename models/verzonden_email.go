package models

import (
	"time"
)

// VerzondEmail representeert een verzonden email
type VerzondEmail struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Ontvanger   string    `json:"ontvanger" gorm:"not null;index"`
	Onderwerp   string    `json:"onderwerp" gorm:"not null"`
	Inhoud      string    `json:"inhoud" gorm:"type:text;not null"`
	VerzondOp   time.Time `json:"verzonden_op" gorm:"autoCreateTime;index"`
	Status      string    `json:"status" gorm:"default:'verzonden';index"`
	FoutBericht string    `json:"fout_bericht" gorm:"type:text"`

	// Optionele relaties
	ContactID    *string `json:"contact_id" gorm:"index"`
	AanmeldingID *string `json:"aanmelding_id" gorm:"index"`
	TemplateID   *string `json:"template_id"`

	// Relaties
	Contact    *ContactFormulier `json:"-" gorm:"foreignKey:ContactID"`
	Aanmelding *Aanmelding       `json:"-" gorm:"foreignKey:AanmeldingID"`
	Template   *EmailTemplate    `json:"-" gorm:"foreignKey:TemplateID"`
}

// TableName specificeert de tabelnaam voor GORM
func (VerzondEmail) TableName() string {
	return "verzonden_emails"
}
