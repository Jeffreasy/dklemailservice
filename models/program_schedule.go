package models

import "time"

type ProgramSchedule struct {
	ID               string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Time             string    `json:"time"`
	EventDescription string    `json:"event_description" gorm:"type:text"`
	Category         string    `json:"category"`
	IconName         string    `json:"icon_name"`
	OrderNumber      int       `json:"order_number"`
	Visible          bool      `json:"visible" gorm:"not null;default:true"`
	Latitude         *float64  `json:"latitude"`
	Longitude        *float64  `json:"longitude"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (ProgramSchedule) TableName() string {
	return "program_schedule"
}
