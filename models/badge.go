package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Badge representeert een badge die deelnemers kunnen verdienen
type Badge struct {
	ID           string        `json:"id" db:"id"`
	Name         string        `json:"name" db:"name"`
	Description  string        `json:"description" db:"description"`
	IconURL      *string       `json:"icon_url,omitempty" db:"icon_url"`
	Criteria     BadgeCriteria `json:"criteria" db:"criteria"`
	Points       int           `json:"points" db:"points"`
	IsActive     bool          `json:"is_active" db:"is_active"`
	DisplayOrder int           `json:"display_order" db:"display_order"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}

// BadgeCriteria bevat de criteria waaraan voldaan moet worden om badge te verdienen
type BadgeCriteria struct {
	MinSteps         *int    `json:"min_steps,omitempty"`
	MinDays          *int    `json:"min_days,omitempty"`
	ConsecutiveDays  *int    `json:"consecutive_days,omitempty"`
	Route            *string `json:"route,omitempty"`
	EarlyParticipant *bool   `json:"early_participant,omitempty"`
	HasTeam          *bool   `json:"has_team,omitempty"`
}

// Value implementeert driver.Valuer voor database opslag
func (c BadgeCriteria) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implementeert sql.Scanner voor database ophalen
func (c *BadgeCriteria) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, c)
}

// BadgeRequest representeert een request voor het aanmaken/bijwerken van een badge
type BadgeRequest struct {
	Name         string        `json:"name" validate:"required,max=100"`
	Description  string        `json:"description" validate:"required"`
	IconURL      *string       `json:"icon_url,omitempty" validate:"omitempty,url,max=500"`
	Criteria     BadgeCriteria `json:"criteria" validate:"required"`
	Points       int           `json:"points" validate:"gte=0"`
	IsActive     *bool         `json:"is_active,omitempty"`
	DisplayOrder *int          `json:"display_order,omitempty" validate:"omitempty,gte=0"`
}

// BadgeWithStats bevat badge informatie met statistieken
type BadgeWithStats struct {
	Badge
	EarnedCount         int        `json:"earned_count"`
	LastEarnedAt        *time.Time `json:"last_earned_at,omitempty"`
	EarnedByCurrentUser bool       `json:"earned_by_current_user"`
}
