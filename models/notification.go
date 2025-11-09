package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DEPRECATED: Old string constants for notification types - use lookup tables (V27)
// Kept for backwards compatibility
const (
	// NotificationPriorityLow represents low priority notifications
	NotificationPriorityLow = "low"

	// NotificationPriorityMedium represents medium priority notifications
	NotificationPriorityMedium = "medium"

	// NotificationPriorityHigh represents high priority notifications
	NotificationPriorityHigh = "high"

	// NotificationPriorityCritical represents critical priority notifications
	NotificationPriorityCritical = "critical"

	// NotificationTypeContact represents contact form notifications
	NotificationTypeContact = "contact"

	// NotificationTypeAanmelding represents registration notifications
	NotificationTypeAanmelding = "aanmelding"

	// NotificationTypeAuth represents authentication notifications
	NotificationTypeAuth = "auth"

	// NotificationTypeSystem represents system notifications
	NotificationTypeSystem = "system"

	// NotificationTypeHealth represents health check notifications
	NotificationTypeHealth = "health"
)

// Notification represents a notification to be sent via Telegram
// Updated for V27: Uses lookup table foreign keys instead of string types
type Notification struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	// V27: Foreign keys to lookup tables (database columns are 'type' and 'priority')
	Type           string                   `json:"type" gorm:"type:text;index;not null"`
	TypeLookup     NotificationType         `json:"type_lookup,omitempty" gorm:"foreignKey:Type;references:Name"`
	Priority       string                   `json:"priority" gorm:"type:text;index;not null"`
	PriorityLookup NotificationPriorityType `json:"priority_lookup,omitempty" gorm:"foreignKey:Priority;references:Name"`

	Title     string     `json:"title" gorm:"type:varchar(255);not null"`
	Message   string     `json:"message" gorm:"type:text;not null"`
	Sent      bool       `json:"sent" gorm:"default:false"`
	SentAt    *time.Time `json:"sent_at" gorm:"type:timestamptz"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
}

// BeforeCreate sets the ID if it's not already set
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
