package models

import "time"

// ChatChannelParticipant represents a participant in a chat channel
type ChatChannelParticipant struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ChannelID  string    `gorm:"type:uuid;index;not null" json:"channel_id"`
	UserID     string    `gorm:"type:uuid;index;not null" json:"user_id"`
	Role       string    `gorm:"type:text;default:'member';check:role IN ('owner', 'admin', 'member')" json:"role"`
	JoinedAt   time.Time `gorm:"default:now()" json:"joined_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
	LastReadAt time.Time `json:"last_read_at"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
}

func (ChatChannelParticipant) TableName() string {
	return "chat_channel_participants"
}
