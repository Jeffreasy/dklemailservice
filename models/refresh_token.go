package models

import "time"

// RefreshToken representeert een refresh token voor JWT authenticatie
// V34: OwnerID kan zowel gebruikers.id als participants.id bevatten
type RefreshToken struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OwnerID   string     `json:"owner_id" gorm:"column:owner_id;not null;type:uuid;index"` // V34: Was UserID, nu OwnerID (kan gebruiker of participant zijn)
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
