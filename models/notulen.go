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
	ID                   uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Titel                string          `json:"titel" gorm:"type:varchar(255);not null"`
	VergaderingDatum     time.Time       `json:"vergadering_datum" gorm:"type:date;not null"`
	Locatie              string          `json:"locatie,omitempty" gorm:"type:varchar(255)"`
	Voorzitter           string          `json:"voorzitter,omitempty" gorm:"type:varchar(255)"`
	Notulist             string          `json:"notulist,omitempty" gorm:"type:varchar(255)"`
	Aanwezigen           pq.StringArray  `json:"aanwezigen,omitempty" gorm:"type:text[]"`            // Legacy field - combined names for backwards compatibility
	Afwezigen            pq.StringArray  `json:"afwezigen,omitempty" gorm:"type:text[]"`             // Legacy field - combined names for backwards compatibility
	AanwezigenGebruikers pq.GenericArray `json:"aanwezigen_gebruikers,omitempty" gorm:"type:uuid[]"` // UUIDs of registered users present
	AfwezigenGebruikers  pq.GenericArray `json:"afwezigen_gebruikers,omitempty" gorm:"type:uuid[]"`  // UUIDs of registered users absent
	AanwezigenGasten     pq.StringArray  `json:"aanwezigen_gasten,omitempty" gorm:"type:text[]"`     // Names of non-registered guests present
	AfwezigenGasten      pq.StringArray  `json:"afwezigen_gasten,omitempty" gorm:"type:text[]"`      // Names of non-registered guests absent
	AgendaItems          AgendaItems     `json:"agenda_items,omitempty" gorm:"type:jsonb"`
	Besluiten            Besluiten       `json:"besluiten,omitempty" gorm:"type:jsonb"`
	Actiepunten          Actiepunten     `json:"actiepunten,omitempty" gorm:"type:jsonb"`
	Notities             string          `json:"notities,omitempty" gorm:"type:text"`
	Status               string          `json:"status" gorm:"type:varchar(50);not null;default:'draft'"`
	Versie               int             `json:"versie" gorm:"type:integer;not null;default:1"`
	CreatedBy            uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt            time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy            uuid.UUID       `json:"updated_by,omitempty" gorm:"type:uuid"`
	FinalizedAt          *time.Time      `json:"finalized_at,omitempty" gorm:"type:timestamp"`
	FinalizedBy          *uuid.UUID      `json:"finalized_by,omitempty" gorm:"type:uuid"`
}

// TableName returns the table name for the Notulen model
func (Notulen) TableName() string {
	return "notulen"
}

// AgendaItem represents an item on the meeting agenda
type AgendaItem struct {
	Title   string `json:"title"`
	Details string `json:"details,omitempty"`
}

// Besluit represents a decision made during the meeting
type Besluit struct {
	Besluit string                 `json:"besluit"`
	Teams   map[string]interface{} `json:"teams,omitempty"`
}

// Actiepunt represents an action item from the meeting
type Actiepunt struct {
	Actie             string      `json:"actie"`
	Verantwoordelijke interface{} `json:"verantwoordelijke"` // Can be string or []string
}

// NotulenVersie represents a version snapshot of notulen
type NotulenVersie struct {
	ID                   uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NotulenID            uuid.UUID       `json:"notulen_id" gorm:"type:uuid;not null"`
	Versie               int             `json:"versie" gorm:"type:integer;not null"`
	Titel                string          `json:"titel" gorm:"type:varchar(255);not null"`
	VergaderingDatum     time.Time       `json:"vergadering_datum" gorm:"type:date;not null"`
	Locatie              string          `json:"locatie,omitempty" gorm:"type:varchar(255)"`
	Voorzitter           string          `json:"voorzitter,omitempty" gorm:"type:varchar(255)"`
	Notulist             string          `json:"notulist,omitempty" gorm:"type:varchar(255)"`
	Aanwezigen           pq.StringArray  `json:"aanwezigen,omitempty" gorm:"type:text[]"`            // Legacy field - combined names for backwards compatibility
	Afwezigen            pq.StringArray  `json:"afwezigen,omitempty" gorm:"type:text[]"`             // Legacy field - combined names for backwards compatibility
	AanwezigenGebruikers pq.GenericArray `json:"aanwezigen_gebruikers,omitempty" gorm:"type:uuid[]"` // UUIDs of registered users present
	AfwezigenGebruikers  pq.GenericArray `json:"afwezigen_gebruikers,omitempty" gorm:"type:uuid[]"`  // UUIDs of registered users absent
	AanwezigenGasten     pq.StringArray  `json:"aanwezigen_gasten,omitempty" gorm:"type:text[]"`     // Names of non-registered guests present
	AfwezigenGasten      pq.StringArray  `json:"afwezigen_gasten,omitempty" gorm:"type:text[]"`      // Names of non-registered guests absent
	AgendaItems          AgendaItems     `json:"agenda_items,omitempty" gorm:"type:jsonb"`
	Besluiten            Besluiten       `json:"besluiten,omitempty" gorm:"type:jsonb"`
	Actiepunten          Actiepunten     `json:"actiepunten,omitempty" gorm:"type:jsonb"`
	Notities             string          `json:"notities,omitempty" gorm:"type:text"`
	Status               string          `json:"status" gorm:"type:varchar(50);not null"`
	GewijzigdDoor        uuid.UUID       `json:"gewijzigd_door" gorm:"type:uuid;not null"`
	GewijzigdOp          time.Time       `json:"gewijzigd_op" gorm:"autoCreateTime"`
	WijzigingReden       string          `json:"wijziging_reden,omitempty" gorm:"type:text"`
}

// TableName returns the table name for the NotulenVersie model
func (NotulenVersie) TableName() string {
	return "notulen_versies"
}

// NotulenCreateRequest represents the request to create new notulen
type NotulenCreateRequest struct {
	Titel                string       `json:"titel" validate:"required,min=3,max=255"`
	VergaderingDatum     string       `json:"vergadering_datum" validate:"required"` // ISO date string
	Locatie              string       `json:"locatie,omitempty"`
	Voorzitter           string       `json:"voorzitter,omitempty"`
	Notulist             string       `json:"notulist,omitempty"`
	Aanwezigen           []string     `json:"aanwezigen,omitempty"`            // Legacy field - combined participant names
	Afwezigen            []string     `json:"afwezigen,omitempty"`             // Legacy field - combined participant names
	AanwezigenGebruikers []string     `json:"aanwezigen_gebruikers,omitempty"` // UUIDs of registered users present
	AfwezigenGebruikers  []string     `json:"afwezigen_gebruikers,omitempty"`  // UUIDs of registered users absent
	AanwezigenGasten     []string     `json:"aanwezigen_gasten,omitempty"`     // Names of non-registered guests present
	AfwezigenGasten      []string     `json:"afwezigen_gasten,omitempty"`      // Names of non-registered guests absent
	AgendaItems          []AgendaItem `json:"agenda_items,omitempty"`
	Besluiten            []Besluit    `json:"besluiten,omitempty"`
	Actiepunten          []Actiepunt  `json:"actiepunten,omitempty"`
	Notities             string       `json:"notities,omitempty"`
}

// NotulenUpdateRequest represents the request to update existing notulen
type NotulenUpdateRequest struct {
	Titel                string       `json:"titel,omitempty" validate:"omitempty,min=3,max=255"`
	Locatie              string       `json:"locatie,omitempty"`
	Voorzitter           string       `json:"voorzitter,omitempty"`
	Notulist             string       `json:"notulist,omitempty"`
	Aanwezigen           []string     `json:"aanwezigen,omitempty"`            // Legacy field - combined participant names
	Afwezigen            []string     `json:"afwezigen,omitempty"`             // Legacy field - combined participant names
	AanwezigenGebruikers []string     `json:"aanwezigen_gebruikers,omitempty"` // UUIDs of registered users present
	AfwezigenGebruikers  []string     `json:"afwezigen_gebruikers,omitempty"`  // UUIDs of registered users absent
	AanwezigenGasten     []string     `json:"aanwezigen_gasten,omitempty"`     // Names of non-registered guests present
	AfwezigenGasten      []string     `json:"afwezigen_gasten,omitempty"`      // Names of non-registered guests absent
	AgendaItems          []AgendaItem `json:"agenda_items,omitempty"`
	Besluiten            []Besluit    `json:"besluiten,omitempty"`
	Actiepunten          []Actiepunt  `json:"actiepunten,omitempty"`
	Notities             string       `json:"notities,omitempty"`
}

// NotulenResponse represents a notulen with resolved user names for API responses
type NotulenResponse struct {
	ID                   uuid.UUID       `json:"id"`
	Titel                string          `json:"titel"`
	VergaderingDatum     time.Time       `json:"vergadering_datum"`
	Locatie              string          `json:"locatie,omitempty"`
	Voorzitter           string          `json:"voorzitter,omitempty"`
	Notulist             string          `json:"notulist,omitempty"`
	Aanwezigen           pq.StringArray  `json:"aanwezigen,omitempty"`
	Afwezigen            pq.StringArray  `json:"afwezigen,omitempty"`
	AanwezigenGebruikers pq.GenericArray `json:"aanwezigen_gebruikers,omitempty"`
	AfwezigenGebruikers  pq.GenericArray `json:"afwezigen_gebruikers,omitempty"`
	AanwezigenGasten     pq.StringArray  `json:"aanwezigen_gasten,omitempty"`
	AfwezigenGasten      pq.StringArray  `json:"afwezigen_gasten,omitempty"`
	AgendaItems          AgendaItems     `json:"agenda_items,omitempty"`
	Besluiten            Besluiten       `json:"besluiten,omitempty"`
	Actiepunten          Actiepunten     `json:"actiepunten,omitempty"`
	Notities             string          `json:"notities,omitempty"`
	Status               string          `json:"status"`
	Versie               int             `json:"versie"`
	CreatedBy            uuid.UUID       `json:"created_by"`
	CreatedByName        string          `json:"created_by_name,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            uuid.UUID       `json:"updated_by,omitempty"`
	UpdatedByName        string          `json:"updated_by_name,omitempty"`
	FinalizedAt          *time.Time      `json:"finalized_at,omitempty"`
	FinalizedBy          *uuid.UUID      `json:"finalized_by,omitempty"`
	FinalizedByName      string          `json:"finalized_by_name,omitempty"`
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

// NotulenFinalizeRequest represents the request to finalize notulen
type NotulenFinalizeRequest struct {
	WijzigingReden string `json:"wijziging_reden,omitempty"`
}

// AgendaItemsWrapper represents the JSONB structure for agenda items
type AgendaItemsWrapper struct {
	Items []AgendaItem `json:"items"`
}

// AgendaItems is a custom type for JSONB storage of agenda items
type AgendaItems AgendaItemsWrapper

// Scan implements sql.Scanner interface for database reading
func (a *AgendaItems) Scan(value interface{}) error {
	if value == nil {
		*a = AgendaItems(AgendaItemsWrapper{Items: []AgendaItem{}})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	var wrapper AgendaItemsWrapper
	if err := json.Unmarshal(bytes, &wrapper); err != nil {
		return err
	}

	*a = AgendaItems(wrapper)
	return nil
}

// Value implements driver.Valuer interface for database writing
func (a AgendaItems) Value() (driver.Value, error) {
	wrapper := AgendaItemsWrapper(a)
	return json.Marshal(wrapper)
}

// BesluitenWrapper represents the JSONB structure for decisions
type BesluitenWrapper struct {
	Besluiten []Besluit `json:"besluiten"`
}

// Besluiten is a custom type for JSONB storage of decisions
type Besluiten BesluitenWrapper

// Scan implements sql.Scanner interface for database reading
func (b *Besluiten) Scan(value interface{}) error {
	if value == nil {
		*b = Besluiten(BesluitenWrapper{Besluiten: []Besluit{}})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	var wrapper BesluitenWrapper
	if err := json.Unmarshal(bytes, &wrapper); err != nil {
		return err
	}

	*b = Besluiten(wrapper)
	return nil
}

// Value implements driver.Valuer interface for database writing
func (b Besluiten) Value() (driver.Value, error) {
	wrapper := BesluitenWrapper(b)
	return json.Marshal(wrapper)
}

// ActiepuntenWrapper represents the JSONB structure for action items
type ActiepuntenWrapper struct {
	Acties []Actiepunt `json:"acties"`
}

// Actiepunten is a custom type for JSONB storage of action items
type Actiepunten ActiepuntenWrapper

// Scan implements sql.Scanner interface for database reading
func (a *Actiepunten) Scan(value interface{}) error {
	if value == nil {
		*a = Actiepunten(ActiepuntenWrapper{Acties: []Actiepunt{}})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	var wrapper ActiepuntenWrapper
	if err := json.Unmarshal(bytes, &wrapper); err != nil {
		return err
	}

	*a = Actiepunten(wrapper)
	return nil
}

// Value implements driver.Valuer interface for database writing
func (a Actiepunten) Value() (driver.Value, error) {
	wrapper := ActiepuntenWrapper(a)
	return json.Marshal(wrapper)
}
