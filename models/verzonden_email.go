package models

import (
	"time"
)

// VerzondEmail representeert een verzonden email
type VerzondEmail struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Ontvanger string    `json:"ontvanger" gorm:"not null;index"`
	Onderwerp string    `json:"onderwerp" gorm:"not null"`
	Inhoud    string    `json:"inhoud" gorm:"type:text;not null"`
	VerzondOp time.Time `json:"verzonden_op" gorm:"autoCreateTime;index"`

	// V27: Foreign key to lookup table (database column is 'status')
	Status     string          `json:"status" gorm:"type:text;default:'verzonden';index"`
	StatusType EmailStatusType `json:"status_type,omitempty" gorm:"foreignKey:Status;references:Status"`

	FoutBericht string `json:"fout_bericht" gorm:"type:text"`

	// Optionele relaties
	ContactID     *string `json:"contact_id" gorm:"index"`
	ParticipantID *string `json:"participant_id" gorm:"index"`
	TemplateID    *string `json:"template_id"`

	// Relaties
	Contact     *ContactFormulier `json:"-" gorm:"foreignKey:ContactID"`
	Participant *Participant      `json:"-" gorm:"foreignKey:ParticipantID"`
	Template    *EmailTemplate    `json:"-" gorm:"foreignKey:TemplateID"`
}

// TableName specificeert de tabelnaam voor GORM
func (VerzondEmail) TableName() string {
	return "verzonden_emails"
}
