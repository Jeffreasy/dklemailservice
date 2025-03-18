package models

import (
	"time"
)

// IncomingEmail vertegenwoordigt een inkomende e-mail in het systeem
type IncomingEmail struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	MessageID   string     `json:"message_id" gorm:"index"`
	From        string     `json:"from"`
	To          string     `json:"to"`
	Subject     string     `json:"subject"`
	Body        string     `json:"body" gorm:"type:text"`
	ContentType string     `json:"content_type"`
	ReceivedAt  time.Time  `json:"received_at"`
	UID         string     `json:"uid" gorm:"uniqueIndex"`
	AccountType string     `json:"account_type" gorm:"index"` // "info" of "inschrijving"
	IsProcessed bool       `json:"is_processed" gorm:"index"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
