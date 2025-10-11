package models

import "time"

type SocialEmbed struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Platform    string    `json:"platform"`
	EmbedCode   string    `json:"embed_code" gorm:"type:text"`
	OrderNumber int       `json:"order_number"`
	Visible     bool      `json:"visible" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (SocialEmbed) TableName() string {
	return "social_embeds"
}
