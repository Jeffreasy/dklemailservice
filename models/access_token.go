package models

import "time"

// AccessToken representeert een access token voor JWT authenticatie met server-side storage
// Hiermee kunnen tokens worden ingetrokken voordat ze verlopen
type AccessToken struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OwnerID   string     `json:"owner_id" gorm:"column:owner_id;not null;type:uuid;index"`      // Kan gebruiker of participant ID zijn
	SessionID *string    `json:"session_id,omitempty" gorm:"column:session_id;type:uuid;index"` // Koppelt token aan een sessie voor session management
	Token     string     `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	IsRevoked bool       `json:"is_revoked" gorm:"default:false;index"`
}

// TableName specificeert de tabelnaam voor GORM
func (AccessToken) TableName() string {
	return "access_tokens"
}

// IsValid controleert of de access token nog geldig is
func (at *AccessToken) IsValid() bool {
	return !at.IsRevoked && at.ExpiresAt.After(time.Now())
}
