package models

import (
	"time"
)

// RBACRole is imported here to avoid circular imports
// The actual definition is in role_rbac.go

// Gebruiker representeert een gebruiker van het systeem
type Gebruiker struct {
	ID                   string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Naam                 string     `json:"naam" gorm:"not null"`
	Email                string     `json:"email" gorm:"not null;uniqueIndex"`
	WachtwoordHash       string     `json:"-" gorm:"not null"`                    // Niet zichtbaar in JSON
	Rol                  string     `json:"rol" gorm:"default:'gebruiker';index"` // Legacy field for backward compatibility
	RoleID               *string    `json:"role_id,omitempty" gorm:"type:uuid"`   // New RBAC role reference
	IsActief             bool       `json:"is_actief" gorm:"default:true"`
	NewsletterSubscribed bool       `json:"newsletter_subscribed" gorm:"default:false;index"`
	LaatsteLogin         *time.Time `json:"laatste_login"`
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Roles []RBACRole `gorm:"many2many:user_roles;" json:"roles,omitempty"` // Many-to-many with RBAC roles
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
