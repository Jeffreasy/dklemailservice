package models

import "time"

// ChatUserPresence represents the presence status of a user
type ChatUserPresence struct {
	UserID    string    `gorm:"type:uuid;primaryKey"`
	Status    string    `gorm:"type:text;default:'offline';check:status IN ('online', 'away', 'busy', 'offline')"`
	LastSeen  time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

func (ChatUserPresence) TableName() string {
	return "chat_user_presence"
}
