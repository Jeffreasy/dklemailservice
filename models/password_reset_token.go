package models

import "time"

// PasswordResetToken representeert een token voor wachtwoord reset
type PasswordResetToken struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string     `json:"email" gorm:"not null;index"`
	Token     string     `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	IsUsed    bool       `json:"is_used" gorm:"default:false;index"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// TableName specificeert de tabelnaam voor GORM
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsValid controleert of de password reset token nog geldig is
func (prt *PasswordResetToken) IsValid() bool {
	return !prt.IsUsed && prt.ExpiresAt.After(time.Now())
}
