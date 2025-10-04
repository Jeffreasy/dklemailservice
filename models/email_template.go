package models

import (
	"time"
)

// EmailTemplate representeert een email template die in de database is opgeslagen
type EmailTemplate struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Naam         string    `json:"naam" gorm:"not null;uniqueIndex"`
	Onderwerp    string    `json:"onderwerp" gorm:"not null"`
	Inhoud       string    `json:"inhoud" gorm:"type:text;not null"`
	Beschrijving string    `json:"beschrijving" gorm:"type:text"`
	IsActief     bool      `json:"is_actief" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy    string    `json:"created_by"`
	UpdatedBy    string    `json:"updated_by"`
}

// TableName specificeert de tabelnaam voor GORM
func (EmailTemplate) TableName() string {
	return "email_templates"
}
