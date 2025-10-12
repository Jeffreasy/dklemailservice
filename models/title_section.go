package models

import (
	"time"
)

// TitleSection represents the title section content for the website
type TitleSection struct {
	ID                 string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EventTitle         string    `json:"event_title" gorm:"not null"`
	EventSubtitle      string    `json:"event_subtitle"`
	ImageURL           string    `json:"image_url"`
	ImageAlt           string    `json:"image_alt"`
	Detail1Title       string    `json:"detail_1_title"`
	Detail1Description string    `json:"detail_1_description"`
	Detail2Title       string    `json:"detail_2_title"`
	Detail2Description string    `json:"detail_2_description"`
	Detail3Title       string    `json:"detail_3_title"`
	Detail3Description string    `json:"detail_3_description"`
	ParticipantCount   int       `json:"participant_count"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (TitleSection) TableName() string {
	return "title_section_content"
}
