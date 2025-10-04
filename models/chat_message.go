package models

import "time"

// ChatMessage represents a message in a chat channel
type ChatMessage struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChannelID   string `gorm:"type:uuid;index;not null"`
	UserID      string `gorm:"type:uuid;index"`
	Content     string `gorm:"type:text"`
	MessageType string `gorm:"type:text;default:'text';check:message_type IN ('text', 'image', 'file', 'system')"`
	FileURL     string `gorm:"type:text"`
	FileName    string `gorm:"type:text"`
	FileSize    int    `gorm:"type:integer"`
	ReplyToID   string `gorm:"type:uuid"`
	EditedAt    time.Time
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
