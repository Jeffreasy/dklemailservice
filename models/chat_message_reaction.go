package models

import "time"

// ChatMessageReaction represents a reaction to a chat message
type ChatMessageReaction struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MessageID string    `gorm:"type:uuid;index;not null"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	Emoji     string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"default:now()"`
}

func (ChatMessageReaction) TableName() string {
	return "chat_message_reactions"
}
