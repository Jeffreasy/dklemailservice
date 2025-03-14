package models

import (
	"time"
)

// Migratie representeert een database migratie
type Migratie struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Versie    string    `json:"versie" gorm:"not null;uniqueIndex"`
	Naam      string    `json:"naam" gorm:"not null"`
	Toegepast time.Time `json:"toegepast" gorm:"autoCreateTime"`
}

// TableName specificeert de tabelnaam voor GORM
func (Migratie) TableName() string {
	return "migraties"
}
