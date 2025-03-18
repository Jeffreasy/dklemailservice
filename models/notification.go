package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationPriority represents the priority level of a notification
type NotificationPriority string

const (
	// NotificationPriorityLow represents low priority notifications
	NotificationPriorityLow NotificationPriority = "low"

	// NotificationPriorityMedium represents medium priority notifications
	NotificationPriorityMedium NotificationPriority = "medium"

	// NotificationPriorityHigh represents high priority notifications
	NotificationPriorityHigh NotificationPriority = "high"

	// NotificationPriorityCritical represents critical priority notifications
	NotificationPriorityCritical NotificationPriority = "critical"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	// NotificationTypeContact represents contact form notifications
	NotificationTypeContact NotificationType = "contact"

	// NotificationTypeAanmelding represents registration notifications
	NotificationTypeAanmelding NotificationType = "aanmelding"

	// NotificationTypeAuth represents authentication notifications
	NotificationTypeAuth NotificationType = "auth"

	// NotificationTypeSystem represents system notifications
	NotificationTypeSystem NotificationType = "system"

	// NotificationTypeHealth represents health check notifications
	NotificationTypeHealth NotificationType = "health"
)

// Notification represents a notification to be sent via Telegram
type Notification struct {
	ID        string               `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Type      NotificationType     `json:"type" gorm:"type:varchar(50);not null"`
	Priority  NotificationPriority `json:"priority" gorm:"type:varchar(20);not null"`
	Title     string               `json:"title" gorm:"type:varchar(255);not null"`
	Message   string               `json:"message" gorm:"type:text;not null"`
	Sent      bool                 `json:"sent" gorm:"default:false"`
	SentAt    *time.Time           `json:"sent_at" gorm:"type:timestamptz"`
	CreatedAt time.Time            `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time            `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
}

// BeforeCreate sets the ID if it's not already set
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
