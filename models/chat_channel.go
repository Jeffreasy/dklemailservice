package models

import "time"

// ChatChannel represents a chat channel
type ChatChannel struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"type:text;not null;check:type IN ('public', 'private', 'direct')" json:"type"`
	CreatedBy   string    `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
}

func (ChatChannel) TableName() string {
	return "chat_channels"
}
