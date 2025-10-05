package models

import "time"

// ChatUserPresence represents the presence status of a user
type ChatUserPresence struct {
	UserID    string    `gorm:"type:uuid;primaryKey" json:"user_id"`
	Status    string    `gorm:"type:text;default:'offline';check:status IN ('online', 'away', 'busy', 'offline')" json:"status"`
	LastSeen  time.Time `gorm:"default:now()" json:"last_seen"`
	UpdatedAt time.Time `gorm:"default:now()" json:"updated_at"`
}

func (ChatUserPresence) TableName() string {
	return "chat_user_presence"
}

// OnlineUser represents a simplified user for online list
type OnlineUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
