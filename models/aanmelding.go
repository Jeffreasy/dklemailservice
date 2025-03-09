package models

import (
	"time"
)

type Aanmelding struct {
	ID             string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Naam           string     `json:"naam"`
	Email          string     `json:"email"`
	Telefoonnummer string     `json:"telefoon"`
	Rol            string     `json:"rol"`
	Afstand        string     `json:"afstand"`
	Ondersteuning  string     `json:"ondersteuning"`
	Bijzonderheden string     `json:"bijzonderheden"`
	Terms          bool       `json:"terms"`
	EmailVerzonden bool       `json:"email_verzonden"`
	EmailVerzondOp *time.Time `json:"email_verzonden_op"`
}
