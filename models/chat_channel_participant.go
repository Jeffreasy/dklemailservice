package models

import "time"

// ChatChannelParticipant represents a participant in a chat channel
type ChatChannelParticipant struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChannelID  string    `gorm:"type:uuid;index;not null"`
	UserID     string    `gorm:"type:uuid;index;not null"`
	Role       string    `gorm:"type:text;default:'member';check:role IN ('owner', 'admin', 'member')"`
	JoinedAt   time.Time `gorm:"default:now()"`
	LastSeenAt time.Time
	IsActive   bool `gorm:"default:true"`
}

func (ChatChannelParticipant) TableName() string {
	return "chat_channel_participants"
}
