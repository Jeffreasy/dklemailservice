package models

import "time"

// ChatMessage represents a message in a chat channel
type ChatMessage struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ChannelID    string    `gorm:"type:uuid;index;not null" json:"channel_id"`
	UserID       string    `gorm:"type:uuid;index" json:"user_id"`
	Content      string    `gorm:"type:text" json:"content"`
	MessageType  string    `gorm:"type:text;default:'text';check:message_type IN ('text', 'image', 'file', 'system')" json:"message_type"`
	FileURL      string    `gorm:"type:text" json:"file_url"`
	FileName     string    `gorm:"type:text" json:"file_name"`
	FileSize     int       `gorm:"type:integer" json:"file_size"`
	ThumbnailURL *string   `gorm:"type:text" json:"thumbnail_url,omitempty"`
	ReplyToID    *string   `gorm:"type:uuid" json:"reply_to_id"`
	EditedAt     time.Time `json:"edited_at"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:now()" json:"updated_at"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}

// MessageWithUser extends ChatMessage with user name for responses
type MessageWithUser struct {
	ChatMessage
	UserName string `json:"user_name"`
}
