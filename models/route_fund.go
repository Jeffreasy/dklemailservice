package models

import (
	"time"
)

// RouteFund vertegenwoordigt de fondsallocatie per route
type RouteFund struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Route     string    `json:"route" gorm:"uniqueIndex;not null"` // Bijv. "6 KM", "10 KM", etc.
	Amount    int       `json:"amount" gorm:"not null"`            // Bedrag in euro's
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specificeert de tabelnaam voor GORM
func (RouteFund) TableName() string {
	return "route_funds"
}

// RouteFundRequest voor API requests
type RouteFundRequest struct {
	Route  string `json:"route" validate:"required"`
	Amount int    `json:"amount" validate:"required,min=0"`
}

// RouteFundResponse voor API responses
type RouteFundResponse struct {
	ID        string    `json:"id"`
	Route     string    `json:"route"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
