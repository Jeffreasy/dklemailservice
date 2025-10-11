package models

import "time"

type Video struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VideoID      string    `json:"video_id"`
	URL          string    `json:"url"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Visible      bool      `json:"visible" gorm:"not null;default:true"`
	OrderNumber  int       `json:"order_number"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (Video) TableName() string {
	return "videos"
}
