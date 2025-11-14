package models

import (
	"time"
)

// AccountType definieert de types van participant accounts
type AccountType string

const (
	AccountTypeFull      AccountType = "full"      // Volledige account met app toegang
	AccountTypeTemporary AccountType = "temporary" // Tijdelijke registratie alleen voor event
)

// Participant representeert de "Persoon" (voorheen Aanmelding)
// GEFIXTE STRUCT: Bevat alleen persoonsgegevens.
// V30: Uitgebreid met duaal registratiesysteem (full vs temporary accounts)
type Participant struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Naam        string    `json:"naam" gorm:"not null"`
	Email       string    `json:"email" gorm:"not null;index"`
	Telefoon    string    `json:"telefoon"`
	Terms       bool      `json:"terms" gorm:"not null"` // 'terms' bleef op Participant (V28)
	GebruikerID *string   `json:"gebruiker_id,omitempty" gorm:"type:uuid;index;column:gebruiker_id"`
	TestMode    bool      `json:"test_mode" gorm:"default:false"` // Dit veld ontbrak nog (uit V1)

	// V30: Duaal registratiesysteem velden
	AccountType           AccountType `json:"account_type" gorm:"type:text;not null;default:'temporary';check:account_type IN ('full', 'temporary')"`
	RegistrationYear      *int        `json:"registration_year,omitempty" gorm:"index"` // Voor temporary accounts: evenementjaar
	WachtwoordHash        *string     `json:"-" gorm:"type:text"`                       // Alleen voor full accounts
	HasAppAccess          bool        `json:"has_app_access" gorm:"default:false;index"`
	UpgradedToGebruikerID *string     `json:"upgraded_to_gebruiker_id,omitempty" gorm:"type:uuid;index"`
	UpgradedAt            *time.Time  `json:"upgraded_at,omitempty"`

	// Relatie met antwoorden
	Antwoorden []ParticipantAntwoord `json:"antwoorden,omitempty" gorm:"foreignKey:ParticipantID"`

	// Relatie met event registraties
	EventRegistrations []EventRegistration `json:"event_registrations,omitempty" gorm:"foreignKey:ParticipantID"`
}

// TableName specificeert de tabelnaam voor GORM
func (Participant) TableName() string {
	return "participants"
}

// IsFullAccount controleert of dit een volledig account is
func (p *Participant) IsFullAccount() bool {
	return p.AccountType == AccountTypeFull
}

// IsTemporaryAccount controleert of dit een tijdelijk account is
func (p *Participant) IsTemporaryAccount() bool {
	return p.AccountType == AccountTypeTemporary
}

// CanAccessApp controleert of participant toegang heeft tot de app
func (p *Participant) CanAccessApp() bool {
	return p.HasAppAccess && p.IsFullAccount()
}

// AanmeldingFormulier (DTO) - DEPRECATED - gebruik PublicRegistrationRequest
type AanmeldingFormulier struct {
	Naam           string `json:"naam"`
	Email          string `json:"email"`
	Telefoon       string `json:"telefoon"`
	Rol            string `json:"rol"`
	Afstand        string `json:"afstand"`
	Ondersteuning  string `json:"ondersteuning"`
	Bijzonderheden string `json:"bijzonderheden"`
	Terms          bool   `json:"terms"`
	TestMode       bool   `json:"test_mode"`
}

// PublicRegistrationRequest - Nieuwe DTO voor publieke registratie (V30)
type PublicRegistrationRequest struct {
	// Persoonlijke gegevens
	Naam     string `json:"naam" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Telefoon string `json:"telefoon"` // Optioneel, maar verplicht voor Begeleider/Vrijwilliger

	// Event keuzes
	Rol            string `json:"rol" binding:"required,oneof=Deelnemer Begeleider Vrijwilliger"`
	Afstand        string `json:"afstand" binding:"required,oneof='2.5 KM' '6 KM' '10 KM' '15 KM'"`
	Ondersteuning  string `json:"ondersteuning" binding:"required,oneof=Ja Nee Anders"`
	Bijzonderheden string `json:"bijzonderheden"` // Verplicht als Ondersteuning = Ja of Anders

	// Account keuze - V30: Nieuw duaal systeem
	WantAccount bool    `json:"want_account"` // true = full account, false = temporary
	Wachtwoord  *string `json:"wachtwoord"`   // Verplicht als WantAccount = true

	// Voorwaarden
	Terms bool `json:"terms" binding:"required"`

	// Event context
	EventID *string `json:"event_id"` // Optioneel - gebruikt default event als niet opgegeven

	// Test mode
	TestMode bool `json:"test_mode"`
}

// PublicRegistrationResponse - Response na succesvolle registratie
type PublicRegistrationResponse struct {
	Success        bool        `json:"success"`
	Message        string      `json:"message"`
	ParticipantID  string      `json:"participant_id"`
	RegistrationID string      `json:"registration_id,omitempty"` // EventRegistration ID
	AccountType    AccountType `json:"account_type"`
	HasAppAccess   bool        `json:"has_app_access"`

	// Extra info voor full accounts
	GebruikerID *string `json:"gebruiker_id,omitempty"`

	// Event details
	EventName string `json:"event_name,omitempty"`
	EventDate string `json:"event_date,omitempty"`
}

// UpgradeToFullAccountRequest - Request voor upgrade van temporary naar full account
type UpgradeToFullAccountRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Wachtwoord string `json:"wachtwoord" binding:"required,min=8"`
}

// UpgradeToFullAccountResponse - Response na upgrade
type UpgradeToFullAccountResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	GebruikerID string `json:"gebruiker_id"`

	// Participant info blijft hetzelfde maar heeft nu app toegang
	ParticipantID string `json:"participant_id"`
	HasAppAccess  bool   `json:"has_app_access"`
}

// ParticipantUpgrade model voor audit trail
type ParticipantUpgrade struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ParticipantID string    `json:"participant_id" gorm:"type:uuid;not null;index"`
	GebruikerID   string    `json:"gebruiker_id" gorm:"type:uuid;not null;index"`
	UpgradedAt    time.Time `json:"upgraded_at" gorm:"autoCreateTime"`
	UpgradedBy    *string   `json:"upgraded_by,omitempty" gorm:"type:uuid"`
	Notes         string    `json:"notes,omitempty" gorm:"type:text"`

	// Relations
	Participant Participant `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	Gebruiker   Gebruiker   `json:"gebruiker,omitempty" gorm:"foreignKey:GebruikerID"`
}

// TableName voor ParticipantUpgrade
func (ParticipantUpgrade) TableName() string {
	return "participant_upgrades"
}
