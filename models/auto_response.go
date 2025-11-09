package models

import (
	"time"
)

// AutoResponse represents an email auto-response configuration
type AutoResponse struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Email     string     `json:"email" gorm:"type:varchar(255);not null"`
	IsActive  bool       `json:"is_active" gorm:"default:false"`
	Subject   string     `json:"subject" gorm:"type:varchar(255)"`
	Message   string     `json:"message" gorm:"type:text"`
	StartDate *time.Time `json:"start_date" gorm:"type:timestamp"`
	EndDate   *time.Time `json:"end_date" gorm:"type:timestamp"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for AutoResponse
func (AutoResponse) TableName() string {
	return "auto_responses"
}
