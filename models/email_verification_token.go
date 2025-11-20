package models

import "time"

// EmailVerificationToken representeert een token voor email verificatie
type EmailVerificationToken struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string     `json:"email" gorm:"not null;index"`
	Token     string     `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	IsUsed    bool       `json:"is_used" gorm:"default:false;index"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UserType  string     `json:"user_type" gorm:"not null;default:'participant'"` // 'gebruiker' or 'participant'
	UserID    string     `json:"user_id" gorm:"not null;index"`                   // ID of the user/participant
}

// TableName specificeert de tabelnaam voor GORM
func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

// IsValid controleert of de email verification token nog geldig is
func (evt *EmailVerificationToken) IsValid() bool {
	return !evt.IsUsed && evt.ExpiresAt.After(time.Now())
}
