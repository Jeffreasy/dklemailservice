package models

import "time"

// RefreshToken representeert een refresh token voor JWT authenticatie
type RefreshToken struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string     `json:"user_id" gorm:"not null;type:uuid;index"`
	Token     string     `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	IsRevoked bool       `json:"is_revoked" gorm:"default:false;index"`
}

// TableName specificeert de tabelnaam voor GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsValid controleert of de refresh token nog geldig is
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && rt.ExpiresAt.After(time.Now())
}
