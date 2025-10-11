package models

import "time"

type Photo struct {
	ID               string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	URL              string    `json:"url" gorm:"not null"`
	AltText          string    `json:"alt_text"`
	Visible          bool      `json:"visible" gorm:"not null;default:true"`
	ThumbnailURL     string    `json:"thumbnail_url"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Year             int       `json:"year"`
	CloudinaryFolder string    `json:"cloudinary_folder"`
}

// TableName specifies the table name for GORM
func (Photo) TableName() string {
	return "photos"
}
