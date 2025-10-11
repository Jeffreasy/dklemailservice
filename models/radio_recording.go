package models

import "time"

type RadioRecording struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text"`
	Date         string    `json:"date"`
	AudioURL     string    `json:"audio_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Visible      bool      `json:"visible" gorm:"not null;default:true"`
	OrderNumber  int       `json:"order_number"`
}

// TableName specificeert de tabelnaam voor GORM
func (RadioRecording) TableName() string {
	return "radio_recordings"
}
