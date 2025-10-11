package models

import "time"

type UnderConstruction struct {
	ID                 int        `json:"id" gorm:"primaryKey;autoIncrement"`
	IsActive           bool       `json:"is_active" gorm:"not null;default:false"`
	Title              string     `json:"title" gorm:"not null"`
	Message            string     `json:"message" gorm:"type:text;not null"`
	FooterText         string     `json:"footer_text"`
	LogoURL            string     `json:"logo_url"`
	ExpectedDate       *time.Time `json:"expected_date"`
	SocialLinks        string     `json:"social_links" gorm:"type:text"` // JSON string
	ProgressPercentage int        `json:"progress_percentage" gorm:"default:0"`
	ContactEmail       string     `json:"contact_email"`
	NewsletterEnabled  bool       `json:"newsletter_enabled" gorm:"not null;default:false"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (UnderConstruction) TableName() string {
	return "under_construction"
}
