package models

import "time"

type SocialLink struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Platform       string    `json:"platform" gorm:"not null"`
	URL            string    `json:"url" gorm:"not null"`
	BgColorClass   *string   `json:"bg_color_class"`
	IconColorClass *string   `json:"icon_color_class"`
	OrderNumber    int       `json:"order_number"`
	Visible        bool      `json:"visible" gorm:"not null;default:true"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (SocialLink) TableName() string {
	return "social_links"
}
