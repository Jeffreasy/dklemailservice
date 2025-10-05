package models

import "time"

// ChatMessageReaction represents a reaction to a chat message
type ChatMessageReaction struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MessageID string    `gorm:"type:uuid;index;not null" json:"message_id"`
	UserID    string    `gorm:"type:uuid;index;not null" json:"user_id"`
	Emoji     string    `gorm:"type:text;not null" json:"emoji"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
}

func (ChatMessageReaction) TableName() string {
	return "chat_message_reactions"
}
