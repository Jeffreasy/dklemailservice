package models

import "time"

type Sponsor struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	LogoURL     string    `json:"logo_url"`
	WebsiteURL  string    `json:"website_url"`
	OrderNumber int       `json:"order_number"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`
	Visible     bool      `json:"visible" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (Sponsor) TableName() string {
	return "sponsors"
}
