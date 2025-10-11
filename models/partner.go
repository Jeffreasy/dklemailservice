package models

import "time"

type Partner struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Logo        string    `json:"logo"`
	Website     string    `json:"website"`
	Tier        string    `json:"tier"`
	Since       string    `json:"since"`
	Visible     bool      `json:"visible" gorm:"not null;default:true"`
	OrderNumber int       `json:"order_number"`
}

// TableName specificeert de tabelnaam voor GORM
func (Partner) TableName() string {
	return "partners"
}
