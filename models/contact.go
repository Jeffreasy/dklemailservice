package models

import "time"

type ContactFormulier struct {
	ID               string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Naam             string     `json:"naam"`
	Email            string     `json:"email"`
	Bericht          string     `json:"bericht"`
	EmailVerzonden   bool       `json:"email_verzonden"`
	EmailVerzondenOp *time.Time `json:"email_verzonden_op"`
	PrivacyAkkoord   bool       `json:"privacy_akkoord"`
	Status           string     `json:"status"`
	BehandeldDoor    *string    `json:"behandeld_door"`
	BehandeldOp      *time.Time `json:"behandeld_op"`
	Notities         *string    `json:"notities"`
}
