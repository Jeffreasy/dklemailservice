package models

import "time"

type AlbumPhoto struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AlbumID     string    `json:"album_id" gorm:"not null"`
	PhotoID     string    `json:"photo_id" gorm:"not null"`
	OrderNumber int       `json:"order_number"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName specifies the table name for GORM
func (AlbumPhoto) TableName() string {
	return "album_photos"
}

// AddPhotoToAlbumRequest represents the request to add a photo to an album
type AddPhotoToAlbumRequest struct {
	PhotoID     string `json:"photo_id" validate:"required"`
	OrderNumber int    `json:"order_number"`
}

// PhotoOrder represents a photo with its order number
type PhotoOrder struct {
	PhotoID     string `json:"photo_id" validate:"required"`
	OrderNumber int    `json:"order_number" validate:"min=0"`
}

// ReorderPhotosRequest represents the request to reorder photos in an album
type ReorderPhotosRequest struct {
	PhotoOrder []PhotoOrder `json:"photo_order" validate:"required,min=1"`
}

// AlbumOrder represents an album with its order number
type AlbumOrder struct {
	ID          string `json:"id" validate:"required"`
	OrderNumber int    `json:"order_number" validate:"min=0"`
}

// ReorderAlbumsRequest represents the request to reorder albums
type ReorderAlbumsRequest struct {
	AlbumOrder []AlbumOrder `json:"album_order" validate:"required,min=1"`
}

// PhotoWithAlbumInfo extends Photo with album relationship information
type PhotoWithAlbumInfo struct {
	Photo
	AlbumID     string `json:"album_id"`
	OrderNumber int    `json:"order_number"`
}
