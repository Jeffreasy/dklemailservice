package models

import (
	"time"
)

// Gebruiker representeert een gebruiker van het systeem
type Gebruiker struct {
	ID             string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Naam           string     `json:"naam" gorm:"not null"`
	Email          string     `json:"email" gorm:"not null;uniqueIndex"`
	WachtwoordHash string     `json:"-" gorm:"not null"` // Niet zichtbaar in JSON
	Rol            string     `json:"rol" gorm:"default:'gebruiker';index"`
	IsActief       bool       `json:"is_actief" gorm:"default:true"`
	LaatsteLogin   *time.Time `json:"laatste_login"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specificeert de tabelnaam voor GORM
func (Gebruiker) TableName() string {
	return "gebruikers"
}

// GebruikerLogin representeert de login gegevens voor een gebruiker
type GebruikerLogin struct {
	Email      string `json:"email" binding:"required"`
	Wachtwoord string `json:"wachtwoord" binding:"required"`
}
