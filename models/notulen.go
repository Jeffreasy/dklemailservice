package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Notulen represents a meeting minutes document
type Notulen struct {
	ID                     uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Titel                  string         `json:"titel" gorm:"type:varchar(255);not null"`
	VergaderingDatum       time.Time      `json:"vergadering_datum" gorm:"type:date;not null"`
	Locatie                string         `json:"locatie,omitempty" gorm:"type:varchar(255)"`
	Voorzitter             string         `json:"voorzitter,omitempty" gorm:"type:varchar(255)"`
	Notulist               string         `json:"notulist,omitempty" gorm:"type:varchar(255)"`
	Aanwezigen             pq.StringArray `json:"aanwezigen,omitempty" gorm:"type:text[]"`               // Legacy field
	Afwezigen              pq.StringArray `json:"afwezigen,omitempty" gorm:"type:text[]"`                // Legacy field
	AanwezigenGebruikerIDs UUIDArray      `json:"aanwezigen_gebruiker_ids,omitempty" gorm:"type:uuid[]"` // New UUID array field
	AfwezigenGebruikerIDs  UUIDArray      `json:"afwezigen_gebruiker_ids,omitempty" gorm:"type:uuid[]"`  // New UUID array field
	AanwezigenGasten       pq.StringArray `json:"aanwezigen_gasten,omitempty" gorm:"type:text[]"`        // New guest array field
	AfwezigenGasten        pq.StringArray `json:"afwezigen_gasten,omitempty" gorm:"type:text[]"`         // New guest array field
	AgendaItems            AgendaItems    `json:"agenda_items,omitempty" gorm:"type:jsonb"`
	Besluiten              Besluiten      `json:"besluiten,omitempty" gorm:"type:jsonb"`
	Actiepunten            Actiepunten    `json:"actiepunten,omitempty" gorm:"type:jsonb"`
	Notities               string         `json:"notities,omitempty" gorm:"type:text"`
	Status                 string         `json:"status" gorm:"type:varchar(50);not null;default:'draft';index"`
	Versie                 int            `json:"versie" gorm:"type:integer;not null;default:1"`
	CreatedBy              uuid.UUID      `json:"created_by" gorm:"type:uuid;not null;index"`
	CreatedAt              time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt              time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedByID            uuid.UUID      `json:"updated_by_id,omitempty" gorm:"type:uuid"` // New updated_by field
	FinalizedAt            *time.Time     `json:"finalized_at,omitempty" gorm:"type:timestamp"`
	FinalizedBy            *uuid.UUID     `json:"finalized_by,omitempty" gorm:"type:uuid"`
}

// TableName returns the table name for the Notulen model
func (Notulen) TableName() string {
	return "notulen"
}

// AgendaItem represents an item on the meeting agenda
type AgendaItem struct {
	Titel        string `json:"titel"`
	Beschrijving string `json:"beschrijving,omitempty"`
	Spreker      string `json:"spreker,omitempty"`
	Tijdslot     string `json:"tijdslot,omitempty"`
}

// AgendaItems is a custom type for JSONB scanning of agenda items
type AgendaItems []AgendaItem

// Scan implements sql.Scanner interface for reading JSONB agenda items from database
func (a *AgendaItems) Scan(value interface{}) error {
	if value == nil {
		*a = AgendaItems{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		// For unsupported types, default to empty array
		*a = AgendaItems{}
		return nil
	}

	if len(bytes) == 0 {
		*a = AgendaItems{}
		return nil
	}

	// *** HIER ZIT DE FIX ***
	// Probeer te unmarshal-en. Als het faalt (omdat het een object is),
	// negeer de fout en zet 'a' op een lege slice.
	if err := json.Unmarshal(bytes, a); err != nil {
		*a = AgendaItems{}
	}
	return nil // Altijd nil teruggeven om de SQL-scan niet te laten falen
}

// Value implements driver.Valuer interface for writing JSONB agenda items to database
func (a AgendaItems) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}
	return json.Marshal(a)
}

// Besluit represents a decision made during the meeting
type Besluit struct {
	Beschrijving      string     `json:"beschrijving"`
	Verantwoordelijke string     `json:"verantwoordelijke,omitempty"`
	Deadline          *time.Time `json:"deadline,omitempty"`
}

// Besluiten is a custom type for JSONB scanning of besluiten
type Besluiten []Besluit

// Scan implements sql.Scanner interface for reading JSONB besluiten from database
func (b *Besluiten) Scan(value interface{}) error {
	if value == nil {
		*b = Besluiten{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		// For unsupported types, default to empty array
		*b = Besluiten{}
		return nil
	}

	if len(bytes) == 0 {
		*b = Besluiten{}
		return nil
	}

	// *** HIER ZIT DE FIX ***
	if err := json.Unmarshal(bytes, b); err != nil {
		*b = Besluiten{}
	}
	return nil // Altijd nil teruggeven
}

// Value implements driver.Valuer interface for writing JSONB besluiten to database
func (b Besluiten) Value() (driver.Value, error) {
	if len(b) == 0 {
		return "[]", nil
	}
	return json.Marshal(b)
}

// Actiepunt represents an action item from the meeting
type Actiepunt struct {
	Beschrijving      string     `json:"beschrijving"`
	Verantwoordelijke string     `json:"verantwoordelijke"`
	Deadline          *time.Time `json:"deadline,omitempty"`
	Status            string     `json:"status"` // pending, in_progress, completed
}

// Actiepunten is a custom type for JSONB scanning of actiepunten
type Actiepunten []Actiepunt

// Scan implements sql.Scanner interface for reading JSONB actiepunten from database
func (a *Actiepunten) Scan(value interface{}) error {
	if value == nil {
		*a = Actiepunten{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		// For unsupported types, default to empty array
		*a = Actiepunten{}
		return nil
	}

	if len(bytes) == 0 {
		*a = Actiepunten{}
		return nil
	}

	// *** HIER ZIT DE FIX ***
	if err := json.Unmarshal(bytes, a); err != nil {
		*a = Actiepunten{}
	}
	return nil // Altijd nil teruggeven
}

// Value implements driver.Valuer interface for writing JSONB actiepunten to database
func (a Actiepunten) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}
	return json.Marshal(a)
}

// NotulenVersie represents a version snapshot of notulen
type NotulenVersie struct {
	ID                     uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NotulenID              uuid.UUID      `json:"notulen_id" gorm:"type:uuid;not null"`
	Versie                 int            `json:"versie" gorm:"type:integer;not null"`
	Titel                  string         `json:"titel" gorm:"type:varchar(255);not null"`
	VergaderingDatum       time.Time      `json:"vergadering_datum" gorm:"type:date;not null;index"`
	Locatie                string         `json:"locatie,omitempty" gorm:"type:varchar(255)"`
	Voorzitter             string         `json:"voorzitter,omitempty" gorm:"type:varchar(255)"`
	Notulist               string         `json:"notulist,omitempty" gorm:"type:varchar(255)"`
	Aanwezigen             pq.StringArray `json:"aanwezigen,omitempty" gorm:"type:text[]"`                     // Legacy field
	Afwezigen              pq.StringArray `json:"afwezigen,omitempty" gorm:"type:text[]"`                      // Legacy field
	AanwezigenGebruikerIDs UUIDArray      `json:"aanwezigen_gebruiker_ids,omitempty" gorm:"type:uuid[];index"` // New UUID array field
	AfwezigenGebruikerIDs  UUIDArray      `json:"afwezigen_gebruiker_ids,omitempty" gorm:"type:uuid[];index"`  // New UUID array field
	AanwezigenGasten       pq.StringArray `json:"aanwezigen_gasten,omitempty" gorm:"type:text[]"`              // New guest array field
	AfwezigenGasten        pq.StringArray `json:"afwezigen_gasten,omitempty" gorm:"type:text[]"`               // New guest array field
	AgendaItems            AgendaItems    `json:"agenda_items,omitempty" gorm:"type:jsonb"`
	Besluiten              Besluiten      `json:"besluiten,omitempty" gorm:"type:jsonb"`
	Actiepunten            Actiepunten    `json:"actiepunten,omitempty" gorm:"type:jsonb"`
	Notities               string         `json:"notities,omitempty" gorm:"type:text"`
	Status                 string         `json:"status" gorm:"type:varchar(50);not null"`
	GewijzigdDoor          uuid.UUID      `json:"gewijzigd_door" gorm:"type:uuid;not null"`
	GewijzigdOp            time.Time      `json:"gewijzigd_op" gorm:"autoCreateTime"`
	WijzigingReden         string         `json:"wijziging_reden,omitempty" gorm:"type:text"`
}

// TableName returns the table name for the NotulenVersie model
func (NotulenVersie) TableName() string {
	return "notulen_versies"
}

// NotulenCreateRequest represents the request to create new notulen
type NotulenCreateRequest struct {
	Titel                  string       `json:"titel" validate:"required,min=3,max=255"`
	VergaderingDatum       string       `json:"vergadering_datum" validate:"required"` // ISO date string
	Locatie                string       `json:"locatie,omitempty"`
	Voorzitter             string       `json:"voorzitter,omitempty"`
	Notulist               string       `json:"notulist,omitempty"`
	Aanwezigen             []string     `json:"aanwezigen,omitempty"`
	Afwezigen              []string     `json:"afwezigen,omitempty"`
	AanwezigenGebruikerIDs []string     `json:"aanwezigen_gebruiker_ids,omitempty"` // UUID strings for registered users
	AfwezigenGebruikerIDs  []string     `json:"afwezigen_gebruiker_ids,omitempty"`  // UUID strings for registered users
	AanwezigenGasten       []string     `json:"aanwezigen_gasten,omitempty"`        // Guest names
	AfwezigenGasten        []string     `json:"afwezigen_gasten,omitempty"`         // Guest names
	AgendaItems            []AgendaItem `json:"agenda_items,omitempty"`
	Besluiten              []Besluit    `json:"besluiten,omitempty"`
	Actiepunten            []Actiepunt  `json:"actiepunten,omitempty"`
	Notities               string       `json:"notities,omitempty"`
}

// NotulenUpdateRequest represents the request to update existing notulen
type NotulenUpdateRequest struct {
	Titel                  string       `json:"titel,omitempty" validate:"omitempty,min=3,max=255"`
	Locatie                string       `json:"locatie,omitempty"`
	Voorzitter             string       `json:"voorzitter,omitempty"`
	Notulist               string       `json:"notulist,omitempty"`
	Aanwezigen             []string     `json:"aanwezigen,omitempty"`
	Afwezigen              []string     `json:"afwezigen,omitempty"`
	AanwezigenGebruikerIDs []string     `json:"aanwezigen_gebruiker_ids,omitempty"` // UUID strings for registered users
	AfwezigenGebruikerIDs  []string     `json:"afwezigen_gebruiker_ids,omitempty"`  // UUID strings for registered users
	AanwezigenGasten       []string     `json:"aanwezigen_gasten,omitempty"`        // Guest names
	AfwezigenGasten        []string     `json:"afwezigen_gasten,omitempty"`         // Guest names
	AgendaItems            []AgendaItem `json:"agenda_items,omitempty"`
	Besluiten              []Besluit    `json:"besluiten,omitempty"`
	Actiepunten            []Actiepunt  `json:"actiepunten,omitempty"`
	Notities               string       `json:"notities,omitempty"`
}

// NotulenListResponse represents the response for listing notulen
type NotulenListResponse struct {
	Notulen []NotulenResponse `json:"notulen"`
	Total   int               `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
}

// NotulenSearchFilters represents search and filter options
type NotulenSearchFilters struct {
	Query     string     `json:"query,omitempty"`      // Full-text search query
	DateFrom  *time.Time `json:"date_from,omitempty"`  // Start date filter
	DateTo    *time.Time `json:"date_to,omitempty"`    // End date filter
	Status    string     `json:"status,omitempty"`     // Status filter
	CreatedBy *uuid.UUID `json:"created_by,omitempty"` // Creator filter
	Limit     int        `json:"limit,omitempty"`      // Pagination limit
	Offset    int        `json:"offset,omitempty"`     // Pagination offset
}

// NotulenResponse represents the response model for notulen with resolved user names
type NotulenResponse struct {
	ID                     uuid.UUID    `json:"id"`
	Titel                  string       `json:"titel"`
	VergaderingDatum       time.Time    `json:"vergadering_datum"`
	Locatie                string       `json:"locatie,omitempty"`
	Voorzitter             string       `json:"voorzitter,omitempty"`
	Notulist               string       `json:"notulist,omitempty"`
	Aanwezigen             []string     `json:"aanwezigen"` // Combined and resolved names
	Afwezigen              []string     `json:"afwezigen"`  // Combined and resolved names
	AanwezigenGebruikerIDs UUIDArray    `json:"aanwezigen_gebruiker_ids,omitempty"`
	AfwezigenGebruikerIDs  UUIDArray    `json:"afwezigen_gebruiker_ids,omitempty"`
	AanwezigenGasten       []string     `json:"aanwezigen_gasten,omitempty"`
	AfwezigenGasten        []string     `json:"afwezigen_gasten,omitempty"`
	AgendaItems            []AgendaItem `json:"agenda_items,omitempty"`
	Besluiten              []Besluit    `json:"besluiten,omitempty"`
	Actiepunten            []Actiepunt  `json:"actiepunten,omitempty"`
	Notities               string       `json:"notities,omitempty"`
	Status                 string       `json:"status"`
	Versie                 int          `json:"versie"`
	CreatedBy              uuid.UUID    `json:"created_by"`
	CreatedByName          string       `json:"created_by_name,omitempty"`
	CreatedAt              time.Time    `json:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at"`
	UpdatedByID            uuid.UUID    `json:"updated_by_id,omitempty"`
	UpdatedByName          string       `json:"updated_by_name,omitempty"`
	FinalizedAt            *time.Time   `json:"finalized_at,omitempty"`
	FinalizedBy            *uuid.UUID   `json:"finalized_by,omitempty"`
	FinalizedByName        string       `json:"finalized_by_name,omitempty"`
}

// NotulenFinalizeRequest represents the request to finalize notulen
type NotulenFinalizeRequest struct {
	WijzigingReden string `json:"wijziging_reden,omitempty"`
}
