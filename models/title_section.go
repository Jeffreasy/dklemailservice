package models

import (
	"time"
)

// TitleSection represents the title section content for the website
type TitleSection struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title        string    `json:"title" gorm:"not null"`
	Subtitle     string    `json:"subtitle"`
	CtaText      string    `json:"cta_text"`
	ImageURL     string    `json:"image_url"`
	EventDetails string    `json:"event_details" gorm:"type:jsonb"`
	Styling      string    `json:"styling" gorm:"type:jsonb"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (TitleSection) TableName() string {
	return "title_sections"
}
