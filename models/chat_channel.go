package models

import "time"

// ChatChannel represents a chat channel
type ChatChannel struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text"`
	Type        string    `gorm:"type:text;not null;check:type IN ('public', 'private', 'direct')"`
	CreatedBy   string    `gorm:"type:uuid"`
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`
	IsActive    bool      `gorm:"default:true"`
}

func (ChatChannel) TableName() string {
	return "chat_channels"
}
