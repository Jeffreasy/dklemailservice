package models

import "time"

type Album struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text"`
	CoverPhotoID string    `json:"cover_photo_id"`
	Visible      bool      `json:"visible" gorm:"not null;default:true"`
	OrderNumber  int       `json:"order_number"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (Album) TableName() string {
	return "albums"
}
