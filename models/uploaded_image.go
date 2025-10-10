package models

import "time"

// UploadedImage tracks metadata for images uploaded to Cloudinary
type UploadedImage struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID       string     `gorm:"type:uuid;index;not null;foreignKey:UserID;references:Gebruikers(ID)" json:"user_id"`
	PublicID     string     `gorm:"type:text;unique;not null" json:"public_id"`
	URL          string     `gorm:"type:text;not null" json:"url"`
	SecureURL    string     `gorm:"type:text;not null" json:"secure_url"`
	Filename     string     `gorm:"type:text;not null" json:"filename"`
	Size         int64      `gorm:"type:bigint;not null" json:"size"`
	MimeType     string     `gorm:"type:text;not null" json:"mime_type"`
	Width        int        `gorm:"type:integer" json:"width"`
	Height       int        `gorm:"type:integer" json:"height"`
	Folder       string     `gorm:"type:text;index;not null" json:"folder"`
	ThumbnailURL *string    `gorm:"type:text" json:"thumbnail_url,omitempty"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
	CreatedAt    time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:now()" json:"updated_at"`
}

func (UploadedImage) TableName() string {
	return "uploaded_images"
}

// UploadedImageWithUser extends UploadedImage with user name for responses
type UploadedImageWithUser struct {
	UploadedImage
	UserName string `json:"user_name"`
}
