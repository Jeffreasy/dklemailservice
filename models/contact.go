package models

import "time"

type ContactFormulier struct {
	ID               string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Naam             string     `json:"naam" gorm:"not null"`
	Email            string     `json:"email" gorm:"not null;index"`
	Bericht          string     `json:"bericht" gorm:"type:text;not null"`
	EmailVerzonden   bool       `json:"email_verzonden" gorm:"default:false"`
	EmailVerzondenOp *time.Time `json:"email_verzonden_op"`
	PrivacyAkkoord   bool       `json:"privacy_akkoord" gorm:"not null"`
	Status           string     `json:"status" gorm:"default:'nieuw';index"`
	BehandeldDoor    *string    `json:"behandeld_door"`
	BehandeldOp      *time.Time `json:"behandeld_op"`
	Notities         *string    `json:"notities" gorm:"type:text"`

	// Nieuwe velden voor antwoorden
	Beantwoord    bool       `json:"beantwoord" gorm:"default:false"`
	AntwoordTekst string     `json:"antwoord_tekst" gorm:"type:text"`
	AntwoordDatum *time.Time `json:"antwoord_datum"`
	AntwoordDoor  string     `json:"antwoord_door"`

	// Relatie met antwoorden
	Antwoorden []ContactAntwoord `json:"antwoorden,omitempty" gorm:"foreignKey:ContactID"`
}

// TableName specificeert de tabelnaam voor GORM
func (ContactFormulier) TableName() string {
	return "contact_formulieren"
}
